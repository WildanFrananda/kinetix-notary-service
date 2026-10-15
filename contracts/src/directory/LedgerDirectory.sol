// SPDX-License-Identifier: MIT
pragma solidity 0.8.37;

import {
    AccessControlDefaultAdminRules
} from "@openzeppelin-contracts/access/extensions/AccessControlDefaultAdminRules.sol";
import {ERC165Checker} from "@openzeppelin-contracts/utils/introspection/ERC165Checker.sol";
import {SafeCast} from "@openzeppelin-contracts/utils/math/SafeCast.sol";

import {ILiabilityLedger} from "../ILiabilityLedger.sol";
import {IStatementCourt} from "../court/IStatementCourt.sol";
import {IWithholdingCourt} from "../withholding/IWithholdingCourt.sol";
import {Generation} from "./Generation.sol";
import {ILedgerDirectory} from "./ILedgerDirectory.sol";
import {PendingGeneration} from "./PendingGeneration.sol";

contract LedgerDirectory is ILedgerDirectory, AccessControlDefaultAdminRules {
    uint64 public immutable foreignSuccessorDelay;

    Generation[] private _generations;
    PendingGeneration private _pending;

    constructor(
        address admin,
        uint48 adminTransferDelay,
        ILiabilityLedger ledger,
        IStatementCourt court,
        IWithholdingCourt withholding,
        uint64 foreignSuccessorDelay_
    ) AccessControlDefaultAdminRules(adminTransferDelay, admin) {
        if (foreignSuccessorDelay_ == 0) {
            revert ZeroDelay();
        }

        foreignSuccessorDelay = foreignSuccessorDelay_;
        _append(_local(ledger, court, withholding));
    }

    function registerSuccessor(
        ILiabilityLedger ledger,
        IStatementCourt court,
        IWithholdingCourt withholding
    ) external onlyRole(DEFAULT_ADMIN_ROLE) returns (uint256) {
        ILiabilityLedger predecessor = _sealedPredecessor();
        Generation memory successor = _local(ledger, court, withholding);
        if (_isSealed(ledger)) {
            revert SuccessorSealed(address(ledger));
        }

        uint64 expectedEpoch = _nextEpochOf(predecessor);
        bytes32 expectedDigest = _headOf(predecessor);
        if (successor.genesisEpoch != expectedEpoch || successor.genesisDigest != expectedDigest) {
            revert GenesisMismatch(expectedEpoch, expectedDigest, successor.genesisEpoch, successor.genesisDigest);
        }

        _append(successor);
        return _generations.length - 1;
    }

    function proposeForeignSuccessor(
        uint64 chainId,
        address ledger,
        address court,
        address withholding,
        uint16 scheme
    ) external onlyRole(DEFAULT_ADMIN_ROLE) {
        ILiabilityLedger predecessor = _sealedPredecessor();
        if (chainId == 0 || chainId == _chainId()) {
            revert NotForeign(chainId);
        }

        if (ledger == address(0) || court == address(0) || withholding == address(0)) {
            revert ZeroAddress();
        }

        if (scheme == 0) {
            revert ZeroScheme();
        }

        uint64 activatesAt = _now() + foreignSuccessorDelay;
        _pending.generation = Generation({
            chainId: chainId,
            ledger: ledger,
            court: court,
            withholding: withholding,
            genesisEpoch: _nextEpochOf(predecessor),
            genesisDigest: _headOf(predecessor),
            scheme: scheme,
            registeredAt: 0
        });
        _pending.activatesAt = activatesAt;
        emit ForeignSuccessorProposed(chainId, ledger, court, withholding, scheme, activatesAt);
    }

    function cancelForeignSuccessor() external onlyRole(DEFAULT_ADMIN_ROLE) {
        if (!_hasPending()) {
            revert NoSuccessorPending();
        }

        Generation memory cancelled = _pending.generation;
        delete _pending;
        emit ForeignSuccessorCancelled(cancelled.chainId, cancelled.ledger);
    }

    function activateForeignSuccessor() external returns (uint256) {
        if (!_hasPending()) {
            revert NoSuccessorPending();
        }

        PendingGeneration memory pending = _pending;
        uint64 blockTime = _now();
        if (blockTime < pending.activatesAt) {
            revert SuccessorTooEarly(pending.activatesAt, blockTime);
        }

        delete _pending;
        _append(pending.generation);
        return _generations.length - 1;
    }

    function generationCount() external view returns (uint256) {
        return _generations.length;
    }

    function generation(uint256 index) external view returns (Generation memory) {
        uint256 count = _generations.length;
        if (index >= count) {
            revert UnknownGeneration(index, count);
        }

        return _generations[index];
    }

    function current() external view returns (Generation memory) {
        return _generations[_generations.length - 1];
    }

    function generationOf(uint64 epoch) external view returns (uint256, Generation memory) {
        uint256 index = _generations.length;
        while (index > 0) {
            index--;
            if (_generations[index].genesisEpoch <= epoch) {
                return (index, _generations[index]);
            }
        }

        revert EpochBeforeDirectory(epoch, _generations[0].genesisEpoch);
    }

    function pendingSuccessor() external view returns (PendingGeneration memory) {
        return _pending;
    }

    function supportsInterface(bytes4 interfaceId) public view override returns (bool) {
        return interfaceId == type(ILedgerDirectory).interfaceId || super.supportsInterface(interfaceId);
    }

    function _append(Generation memory next) private {
        uint256 index = _generations.length;
        next.registeredAt = _now();
        _generations.push(next);
        emit GenerationRegistered(
            index,
            next.chainId,
            next.ledger,
            next.court,
            next.withholding,
            next.genesisEpoch,
            next.genesisDigest,
            next.scheme
        );
    }

    function _sealedPredecessor() private view returns (ILiabilityLedger) {
        if (_hasPending()) {
            revert SuccessorPending(_pending.activatesAt);
        }

        Generation memory last = _generations[_generations.length - 1];
        if (last.chainId != _chainId()) {
            revert DirectoryMoved(last.chainId);
        }

        ILiabilityLedger predecessor = ILiabilityLedger(last.ledger);
        if (!_isSealed(predecessor)) {
            revert PredecessorNotSealed(last.ledger);
        }

        return predecessor;
    }

    function _local(
        ILiabilityLedger ledger,
        IStatementCourt court,
        IWithholdingCourt withholding
    ) private view returns (Generation memory) {
        if (!ERC165Checker.supportsInterface(address(ledger), type(ILiabilityLedger).interfaceId)) {
            revert NotALedger(address(ledger));
        }

        if (!ERC165Checker.supportsInterface(address(court), type(IStatementCourt).interfaceId)) {
            revert NotAStatementCourt(address(court));
        }

        if (!ERC165Checker.supportsInterface(address(withholding), type(IWithholdingCourt).interfaceId)) {
            revert NotAWithholdingCourt(address(withholding));
        }

        if (address(_ledgerOf(court)) != address(ledger)) {
            revert CourtOfAnotherLedger(address(court), address(ledger));
        }

        if (address(_courtOf(withholding)) != address(court)) {
            revert WithholdingOfAnotherCourt(address(withholding), address(court));
        }

        return Generation({
            chainId: _chainId(),
            ledger: address(ledger),
            court: address(court),
            withholding: address(withholding),
            genesisEpoch: _genesisEpochOf(ledger),
            genesisDigest: _genesisDigestOf(ledger),
            scheme: ledger.scheme(),
            registeredAt: 0
        });
    }

    function _hasPending() private view returns (bool) {
        return _pending.generation.chainId != 0;
    }

    function _isSealed(ILiabilityLedger ledger) private view returns (bool) {
        return ledger.isSealed();
    }

    function _nextEpochOf(ILiabilityLedger ledger) private view returns (uint64) {
        return ledger.nextEpoch();
    }

    function _headOf(ILiabilityLedger ledger) private view returns (bytes32) {
        return ledger.head();
    }

    function _genesisEpochOf(ILiabilityLedger ledger) private view returns (uint64) {
        return ledger.genesisEpoch();
    }

    function _genesisDigestOf(ILiabilityLedger ledger) private view returns (bytes32) {
        return ledger.genesisDigest();
    }

    function _ledgerOf(IStatementCourt court) private view returns (ILiabilityLedger) {
        return court.ledger();
    }

    function _courtOf(IWithholdingCourt withholding) private view returns (IStatementCourt) {
        return withholding.court();
    }

    function _chainId() private view returns (uint64) {
        return SafeCast.toUint64(block.chainid);
    }

    function _now() private view returns (uint64) {
        return SafeCast.toUint64(block.timestamp);
    }
}
