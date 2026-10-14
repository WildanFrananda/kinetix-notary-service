// SPDX-License-Identifier: MIT
pragma solidity 0.8.37;

import {ILiabilityLedger} from "../ILiabilityLedger.sol";
import {SumNode} from "../SumNode.sol";
import {IStatementCourt} from "../court/IStatementCourt.sol";
import {OpenRequest} from "./OpenRequest.sol";
import {Registration} from "./Registration.sol";
import {Silence} from "./Silence.sol";
import {StatementRequest} from "./StatementRequest.sol";

interface IWithholdingCourt {
    event StatementRequested(
        bytes32 indexed registration,
        bytes32 indexed account,
        uint64 indexed epoch,
        address holder,
        uint64 dueAt
    );
    event StatementAnswered(
        bytes32 indexed registration,
        bytes32 indexed account,
        uint64 indexed epoch,
        bytes32 commitment,
        uint64 amount,
        address answerer
    );
    event SilenceRecorded(
        bytes32 indexed registration,
        bytes32 indexed account,
        uint64 indexed epoch,
        address recorder
    );

    error NotACourt(address court);
    error DurationOutOfRange(uint64 given, uint64 maximum);
    error NotTheHolder(address holder, address caller);
    error NotSignedBySigner(uint64 fromEpoch, address signer, address recovered);
    error BeforeRegistration(uint64 epoch, uint64 fromEpoch);
    error EpochNotAnchored(uint64 epoch);
    error RequestWindowClosed(uint64 epoch, uint256 closedAt);
    error RequestStillOpen(bytes32 registration, uint64 epoch);
    error AlreadyRequested(bytes32 registration, uint64 epoch);
    error NoOpenRequest(bytes32 registration, uint64 epoch);
    error AnswerTooLate(uint64 dueAt, uint64 blockTime);
    error AnswerDoesNotHold(uint64 epoch);
    error NotOverdue(uint64 dueAt, uint64 blockTime);

    function requestStatement(
        Registration calldata registration,
        bytes calldata signature,
        uint64 epoch
    ) external returns (bytes32);

    function answer(
        Registration calldata registration,
        uint64 epoch,
        bytes32 salt,
        uint64 amount,
        uint64 index,
        SumNode[] calldata siblings
    ) external returns (bytes32);

    function recordUnanswered(Registration calldata registration, uint64 epoch) external returns (bytes32);

    function ledger() external view returns (ILiabilityLedger);

    function court() external view returns (IStatementCourt);

    function requestWindow() external view returns (uint64);

    function answerPeriod() external view returns (uint64);

    function registrationDigest(Registration calldata registration) external view returns (bytes32);

    function recoverSigner(Registration calldata registration, bytes calldata signature) external view returns (address);

    function commitmentOf(bytes32 account, bytes32 salt) external pure returns (bytes32);

    function requestOf(bytes32 registration, uint64 epoch) external view returns (StatementRequest memory);

    function openRequestOf(bytes32 registration) external view returns (OpenRequest memory);

    function silenceCount(uint64 epoch) external view returns (uint256);

    function silencesOf(uint64 epoch) external view returns (Silence[] memory);

    function totalSilences() external view returns (uint256);
}
