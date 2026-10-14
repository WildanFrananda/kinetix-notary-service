// SPDX-License-Identifier: MIT
pragma solidity 0.8.37;

import {ECDSA} from "@openzeppelin-contracts/utils/cryptography/ECDSA.sol";
import {EIP712} from "@openzeppelin-contracts/utils/cryptography/EIP712.sol";
import {ERC165} from "@openzeppelin-contracts/utils/introspection/ERC165.sol";
import {ERC165Checker} from "@openzeppelin-contracts/utils/introspection/ERC165Checker.sol";
import {SafeCast} from "@openzeppelin-contracts/utils/math/SafeCast.sol";

import {ILiabilityLedger} from "../ILiabilityLedger.sol";
import {SumNode} from "../SumNode.sol";
import {IStatementCourt} from "../court/IStatementCourt.sol";
import {AccountCommitments} from "./AccountCommitments.sol";
import {IWithholdingCourt} from "./IWithholdingCourt.sol";
import {OpenRequest} from "./OpenRequest.sol";
import {Registration} from "./Registration.sol";
import {RegistrationHashing} from "./RegistrationHashing.sol";
import {RequestState} from "./RequestState.sol";
import {Silence} from "./Silence.sol";
import {StatementRequest} from "./StatementRequest.sol";

contract WithholdingCourt is IWithholdingCourt, EIP712, ERC165 {
    uint64 internal constant MAX_DURATION = 365 days;

    ILiabilityLedger public immutable ledger;
    IStatementCourt public immutable court;
    uint64 public immutable requestWindow;
    uint64 public immutable answerPeriod;

    mapping(bytes32 registration => mapping(uint64 epoch => StatementRequest)) private _requests;
    mapping(bytes32 registration => OpenRequest) private _open;
    mapping(uint64 epoch => Silence[]) private _silences;
    uint256 private _totalSilences;

    constructor(
        IStatementCourt court_,
        uint64 requestWindow_,
        uint64 answerPeriod_
    ) EIP712(RegistrationHashing.NAME, RegistrationHashing.VERSION) {
        _requireDuration(requestWindow_);
        _requireDuration(answerPeriod_);

        requestWindow = requestWindow_;
        answerPeriod = answerPeriod_;
        court = court_;
        ledger = _ledgerOf(court_);
    }

    function requestStatement(
        Registration calldata registration,
        bytes calldata signature,
        uint64 epoch
    ) external returns (bytes32) {
        if (msg.sender != registration.holder) {
            revert NotTheHolder(registration.holder, msg.sender);
        }

        bytes32 digest = registrationDigest(registration);
        address signer = _signerOf(registration.fromEpoch);
        address recovered = _recover(digest, signature);
        if (recovered == address(0) || recovered != signer) {
            revert NotSignedBySigner(registration.fromEpoch, signer, recovered);
        }

        if (epoch < registration.fromEpoch) {
            revert BeforeRegistration(epoch, registration.fromEpoch);
        }

        uint256 closesAt = uint256(_anchoredAt(epoch)) + requestWindow;
        if (block.timestamp > closesAt) {
            revert RequestWindowClosed(epoch, closesAt);
        }

        OpenRequest memory open = _open[digest];
        if (open.active) {
            revert RequestStillOpen(digest, open.epoch);
        }

        StatementRequest storage request = _requests[digest][epoch];
        if (request.state != RequestState.None) {
            revert AlreadyRequested(digest, epoch);
        }

        uint64 dueAt = SafeCast.toUint64(block.timestamp) + answerPeriod;
        request.state = RequestState.Open;
        request.dueAt = dueAt;
        _open[digest] = OpenRequest({active: true, epoch: epoch});

        emit StatementRequested(digest, registration.account, epoch, registration.holder, dueAt);
        return digest;
    }

    function answer(
        Registration calldata registration,
        uint64 epoch,
        bytes32 salt,
        uint64 amount,
        uint64 index,
        SumNode[] calldata siblings
    ) external returns (bytes32) {
        bytes32 digest = registrationDigest(registration);
        StatementRequest storage request = _openRequest(digest, epoch);

        uint64 blockTime = SafeCast.toUint64(block.timestamp);
        if (blockTime > request.dueAt) {
            revert AnswerTooLate(request.dueAt, blockTime);
        }

        bytes32 commitment = AccountCommitments.commitment(registration.account, salt);
        if (!_includes(epoch, commitment, amount, index, siblings)) {
            revert AnswerDoesNotHold(epoch);
        }

        request.state = RequestState.Answered;
        delete _open[digest];

        emit StatementAnswered(digest, registration.account, epoch, commitment, amount, msg.sender);
        return commitment;
    }

    function recordUnanswered(Registration calldata registration, uint64 epoch) external returns (bytes32) {
        bytes32 digest = registrationDigest(registration);
        StatementRequest storage request = _openRequest(digest, epoch);

        uint64 blockTime = SafeCast.toUint64(block.timestamp);
        if (blockTime <= request.dueAt) {
            revert NotOverdue(request.dueAt, blockTime);
        }

        request.state = RequestState.Unanswered;
        delete _open[digest];
        _silences[epoch].push(
            Silence({
                registration: digest,
                account: registration.account,
                holder: registration.holder,
                dueAt: request.dueAt,
                recorder: msg.sender,
                recordedAt: blockTime
            })
        );
        _totalSilences += 1;

        emit SilenceRecorded(digest, registration.account, epoch, msg.sender);
        return digest;
    }

    function registrationDigest(Registration calldata registration) public view returns (bytes32) {
        return _hashTypedDataV4(RegistrationHashing.structHash(registration));
    }

    function recoverSigner(
        Registration calldata registration,
        bytes calldata signature
    ) external view returns (address) {
        return _recover(registrationDigest(registration), signature);
    }

    function commitmentOf(bytes32 account, bytes32 salt) external pure returns (bytes32) {
        return AccountCommitments.commitment(account, salt);
    }

    function requestOf(bytes32 registration, uint64 epoch) external view returns (StatementRequest memory) {
        return _requests[registration][epoch];
    }

    function openRequestOf(bytes32 registration) external view returns (OpenRequest memory) {
        return _open[registration];
    }

    function silenceCount(uint64 epoch) external view returns (uint256) {
        return _silences[epoch].length;
    }

    function silencesOf(uint64 epoch) external view returns (Silence[] memory) {
        return _silences[epoch];
    }

    function totalSilences() external view returns (uint256) {
        return _totalSilences;
    }

    function supportsInterface(bytes4 interfaceId) public view override returns (bool) {
        return interfaceId == type(IWithholdingCourt).interfaceId || super.supportsInterface(interfaceId);
    }

    function _openRequest(bytes32 digest, uint64 epoch) private view returns (StatementRequest storage) {
        StatementRequest storage request = _requests[digest][epoch];
        if (request.state != RequestState.Open) {
            revert NoOpenRequest(digest, epoch);
        }

        return request;
    }

    function _ledgerOf(IStatementCourt court_) private view returns (ILiabilityLedger) {
        if (!ERC165Checker.supportsInterface(address(court_), type(IStatementCourt).interfaceId)) {
            revert NotACourt(address(court_));
        }

        return court_.ledger();
    }

    function _signerOf(uint64 fromEpoch) private view returns (address) {
        return court.signerOf(fromEpoch);
    }

    function _anchoredAt(uint64 epoch) private view returns (uint64) {
        if (epoch < ledger.genesisEpoch() || epoch >= ledger.nextEpoch()) {
            revert EpochNotAnchored(epoch);
        }

        return ledger.epochAt(epoch).anchoredAt;
    }

    function _includes(
        uint64 epoch,
        bytes32 commitment,
        uint64 amount,
        uint64 index,
        SumNode[] calldata siblings
    ) private view returns (bool) {
        return ledger.verifyInclusion(epoch, commitment, amount, index, siblings);
    }

    function _recover(bytes32 digest, bytes calldata signature) private pure returns (address) {
        (address recovered, ECDSA.RecoverError failure, bytes32 detail) = ECDSA.tryRecoverCalldata(digest, signature);
        if (failure != ECDSA.RecoverError.NoError || detail != bytes32(0)) {
            return address(0);
        }

        return recovered;
    }

    function _requireDuration(uint64 duration) private pure {
        if (duration == 0 || duration > MAX_DURATION) {
            revert DurationOutOfRange(duration, MAX_DURATION);
        }
    }
}
