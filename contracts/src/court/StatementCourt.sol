// SPDX-License-Identifier: MIT
pragma solidity 0.8.37;

import {
    AccessControlDefaultAdminRules
} from "@openzeppelin-contracts/access/extensions/AccessControlDefaultAdminRules.sol";
import {ECDSA} from "@openzeppelin-contracts/utils/cryptography/ECDSA.sol";
import {EIP712} from "@openzeppelin-contracts/utils/cryptography/EIP712.sol";
import {ERC165Checker} from "@openzeppelin-contracts/utils/introspection/ERC165Checker.sol";
import {SafeCast} from "@openzeppelin-contracts/utils/math/SafeCast.sol";

import {ILiabilityLedger} from "../ILiabilityLedger.sol";
import {SumNode} from "../SumNode.sol";
import {EpochStatement} from "./EpochStatement.sol";
import {Fault} from "./Fault.sol";
import {FaultKind} from "./FaultKind.sol";
import {IStatementCourt} from "./IStatementCourt.sol";
import {PendingSigner} from "./PendingSigner.sol";
import {SignerTerm} from "./SignerTerm.sol";
import {StatementHashing} from "./StatementHashing.sol";

contract StatementCourt is IStatementCourt, AccessControlDefaultAdminRules, EIP712 {
    ILiabilityLedger public immutable ledger;
    uint64 public immutable signerChangeDelay;

    SignerTerm[] private _terms;
    PendingSigner private _pending;
    mapping(uint64 epoch => Fault[]) private _faults;
    mapping(bytes32 evidence => bool) private _proven;
    uint256 private _totalFaults;

    constructor(
        address admin,
        uint48 adminTransferDelay,
        ILiabilityLedger ledger_,
        address initialSigner,
        uint64 signerChangeDelay_
    )
        AccessControlDefaultAdminRules(adminTransferDelay, admin)
        EIP712(StatementHashing.NAME, StatementHashing.VERSION)
    {
        if (!ERC165Checker.supportsInterface(address(ledger_), type(ILiabilityLedger).interfaceId)) {
            revert NotALedger(address(ledger_));
        }

        if (initialSigner == address(0)) {
            revert ZeroSigner();
        }

        if (signerChangeDelay_ == 0) {
            revert ZeroDelay();
        }

        ledger = ledger_;
        signerChangeDelay = signerChangeDelay_;
        _terms.push(SignerTerm({signer: initialSigner, fromEpoch: ledger_.genesisEpoch()}));
    }

    function proveBadProof(
        EpochStatement calldata statement,
        bytes calldata signature,
        SumNode[] calldata siblings
    ) external returns (bytes32) {
        bytes32 digest = _attributed(statement, signature);
        bytes32 given = StatementHashing.siblingsHash(siblings);
        if (given != statement.siblingsHash) {
            revert SiblingsMismatch(statement.siblingsHash, given);
        }

        bool holds =
            ledger.verifyInclusion(statement.epoch, statement.commitment, statement.amount, statement.index, siblings);
        if (holds) {
            revert ProofHolds(statement.epoch);
        }

        _record(statement.epoch, FaultKind.BadProof, digest);
        return digest;
    }

    function proveEquivocation(
        EpochStatement calldata first,
        bytes calldata firstSignature,
        EpochStatement calldata second,
        bytes calldata secondSignature
    ) external returns (bytes32) {
        if (first.epoch != second.epoch) {
            revert DifferentEpochs(first.epoch, second.epoch);
        }

        bytes32 firstDigest = _attributed(first, firstSignature);
        bytes32 secondDigest = _attributed(second, secondSignature);
        if (firstDigest == secondDigest || !StatementHashing.conflict(first, second)) {
            revert NoConflict();
        }

        bytes32 evidence = firstDigest < secondDigest
            ? keccak256(abi.encodePacked(firstDigest, secondDigest))
            : keccak256(abi.encodePacked(secondDigest, firstDigest));
        _record(first.epoch, FaultKind.Equivocation, evidence);
        return evidence;
    }

    function scheduleSigner(address signer, uint64 fromEpoch) external onlyRole(DEFAULT_ADMIN_ROLE) {
        if (signer == address(0)) {
            revert ZeroSigner();
        }

        if (_pending.signer != address(0)) {
            revert SignerChangePending(_pending.signer, _pending.fromEpoch);
        }

        SignerTerm storage last = _terms[_terms.length - 1];
        if (signer == last.signer) {
            revert SameSigner(signer);
        }

        if (last.fromEpoch == type(uint64).max) {
            revert FromEpochTooLow(fromEpoch, type(uint64).max);
        }

        uint64 least = _leastFromEpoch(last.fromEpoch);
        if (fromEpoch < least) {
            revert FromEpochTooLow(fromEpoch, least);
        }

        uint64 activatesAt = SafeCast.toUint64(block.timestamp) + signerChangeDelay;
        _pending = PendingSigner({signer: signer, fromEpoch: fromEpoch, activatesAt: activatesAt});
        emit SignerScheduled(signer, fromEpoch, activatesAt);
    }

    function cancelSigner() external onlyRole(DEFAULT_ADMIN_ROLE) {
        PendingSigner memory pending = _pending;
        if (pending.signer == address(0)) {
            revert NoSignerChangePending();
        }

        delete _pending;
        emit SignerCancelled(pending.signer, pending.fromEpoch);
    }

    function activateSigner() external {
        PendingSigner memory pending = _pending;
        if (pending.signer == address(0)) {
            revert NoSignerChangePending();
        }

        uint64 blockTime = SafeCast.toUint64(block.timestamp);
        if (blockTime < pending.activatesAt) {
            revert SignerChangeTooEarly(pending.activatesAt, blockTime);
        }

        _requireCurrent(pending.fromEpoch);
        _terms.push(SignerTerm({signer: pending.signer, fromEpoch: pending.fromEpoch}));
        delete _pending;
        emit SignerActivated(pending.signer, pending.fromEpoch);
    }

    function signerOf(uint64 epoch) public view returns (address) {
        for (uint256 i = _terms.length; i > 0; i--) {
            SignerTerm storage term = _terms[i - 1];
            if (term.fromEpoch <= epoch) {
                return term.signer;
            }
        }

        return address(0);
    }

    function signerTerms() external view returns (SignerTerm[] memory) {
        return _terms;
    }

    function pendingSigner() external view returns (PendingSigner memory) {
        return _pending;
    }

    function statementDigest(EpochStatement calldata statement) public view returns (bytes32) {
        return _hashTypedDataV4(StatementHashing.structHash(statement));
    }

    function siblingsHash(SumNode[] calldata siblings) external pure returns (bytes32) {
        return StatementHashing.siblingsHash(siblings);
    }

    function recoverSigner(
        EpochStatement calldata statement,
        bytes calldata signature
    ) external view returns (address) {
        return _recover(statementDigest(statement), signature);
    }

    function faultCount(uint64 epoch) external view returns (uint256) {
        return _faults[epoch].length;
    }

    function faultsOf(uint64 epoch) external view returns (Fault[] memory) {
        return _faults[epoch];
    }

    function totalFaults() external view returns (uint256) {
        return _totalFaults;
    }

    function isProven(bytes32 evidence) external view returns (bool) {
        return _proven[evidence];
    }

    function supportsInterface(bytes4 interfaceId) public view override returns (bool) {
        return interfaceId == type(IStatementCourt).interfaceId || super.supportsInterface(interfaceId);
    }

    function _attributed(EpochStatement calldata statement, bytes calldata signature) private view returns (bytes32) {
        if (statement.ledger != address(ledger)) {
            revert WrongLedger(address(ledger), statement.ledger);
        }

        if (statement.epoch < ledger.genesisEpoch() || statement.epoch >= ledger.nextEpoch()) {
            revert EpochNotAnchored(statement.epoch);
        }

        bytes32 digest = statementDigest(statement);
        address signer = signerOf(statement.epoch);
        address recovered = _recover(digest, signature);
        if (recovered == address(0) || recovered != signer) {
            revert NotSignedBySigner(statement.epoch, signer, recovered);
        }

        return digest;
    }

    function _leastFromEpoch(uint64 lastFromEpoch) private view returns (uint64) {
        uint64 next = ledger.nextEpoch();
        if (lastFromEpoch + 1 > next) {
            return lastFromEpoch + 1;
        }

        return next;
    }

    function _requireCurrent(uint64 fromEpoch) private view {
        uint64 next = ledger.nextEpoch();
        if (fromEpoch < next) {
            revert SignerChangeStale(fromEpoch, next);
        }
    }

    function _recover(bytes32 digest, bytes calldata signature) private pure returns (address) {
        (address recovered, ECDSA.RecoverError failure, bytes32 detail) = ECDSA.tryRecoverCalldata(digest, signature);
        if (failure != ECDSA.RecoverError.NoError || detail != bytes32(0)) {
            return address(0);
        }

        return recovered;
    }

    function _record(uint64 epoch, FaultKind kind, bytes32 evidence) private {
        if (_proven[evidence]) {
            revert AlreadyProven(evidence);
        }

        _proven[evidence] = true;
        _faults[epoch].push(
            Fault({kind: kind, evidence: evidence, prover: msg.sender, provenAt: SafeCast.toUint64(block.timestamp)})
        );
        _totalFaults += 1;
        emit FaultProven(epoch, kind, evidence, msg.sender);
    }
}
