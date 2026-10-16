// SPDX-License-Identifier: MIT
pragma solidity 0.8.37;

import {IEntryPoint} from "account-abstraction/contracts/interfaces/IEntryPoint.sol";
import {IPaymaster} from "account-abstraction/contracts/interfaces/IPaymaster.sol";
import {IStakeManager} from "account-abstraction/contracts/interfaces/IStakeManager.sol";

import {IStatementCourt} from "../court/IStatementCourt.sol";
import {IWithholdingCourt} from "../withholding/IWithholdingCourt.sol";

interface INotaryPaymaster is IPaymaster {
    event Withdrawn(address indexed to, uint256 amount);

    error NoContract(address target);
    error NotAStatementCourt(address court);
    error NotAWithholdingCourt(address withholding);
    error WithholdingOfAnotherCourt(address withholding, address court);
    error NoStatementSigner(address court);
    error DomainMismatch(address court);
    error ZeroCostCeiling();
    error UnstakeDelayTooShort(uint32 given, uint32 least);
    error NotTheEntryPoint(address caller);
    error CostTooHigh(uint256 maxCost, uint256 ceiling);
    error NotDelegated(address sender);
    error NotAnExecute(bytes4 selector);
    error ValueNotZero(uint256 value);
    error TargetNotSponsored(address target);
    error CallNotSponsored(address target, bytes4 selector);
    error NotTheHolder(address holder, address sender);
    error NotSignedByTheNotary(address recovered);
    error AnotherSignersEpoch(uint64 epoch, address signer);
    error EpochNotAnchored(uint64 epoch);
    error BeforeRegistration(uint64 epoch, uint64 fromEpoch);
    error RequestStillOpen(bytes32 registration, uint64 epoch);
    error AlreadyRequested(bytes32 registration, uint64 epoch);
    error NoOpenRequest(bytes32 registration, uint64 epoch);
    error WrongLedger(address expected, address given);
    error SiblingsMismatch(bytes32 signed, bytes32 given);
    error ProofHolds(uint64 epoch);
    error DifferentEpochs(uint64 first, uint64 second);
    error NoConflict();
    error AlreadyProven(bytes32 evidence);
    error BalanceNotSent(address to, uint256 amount);
    error NoPostOp();

    function withdrawTo(address payable to, uint256 amount) external;

    function addStake() external payable;

    function withdrawBalance() external;

    function unlockStake() external;

    function withdrawStake(address payable to) external;

    function sponsorship(
        address sender,
        bytes calldata callData,
        bytes calldata paymasterData
    ) external view returns (uint48 validAfter, uint48 validUntil);

    function entryPoint() external view returns (IEntryPoint);

    function account() external view returns (address);

    function court() external view returns (IStatementCourt);

    function withholding() external view returns (IWithholdingCourt);

    function statementSigner() external view returns (address);

    function maxCostPerOperation() external view returns (uint256);

    function unstakeDelay() external view returns (uint32);

    function getDeposit() external view returns (uint256);

    function getDepositInfo() external view returns (IStakeManager.DepositInfo memory);
}
