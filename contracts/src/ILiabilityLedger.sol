// SPDX-License-Identifier: UNLICENSED
pragma solidity 0.8.37;

import {Epoch} from "./Epoch.sol";
import {EpochHeader} from "./EpochHeader.sol";
import {SumNode} from "./SumNode.sol";

interface ILiabilityLedger {
    event EpochAnchored(
        uint64 indexed number,
        bytes32 indexed digest,
        bytes32 previousDigest,
        bytes32 root,
        uint64 total,
        uint64 leafCount,
        uint64 asOf
    );
    event LedgerSealed(bytes32 indexed head, uint64 nextEpoch, uint64 lastAsOf);

    error InvalidGenesis(uint64 genesisEpoch, bytes32 genesisDigest);
    error ZeroAnchorer();
    error LedgerIsSealed();
    error EpochBeforeGenesis(uint64 genesisEpoch, uint64 number);
    error EpochAlreadyAnchored(uint64 number, bytes32 digest);
    error EpochConflict(uint64 number, bytes32 anchored, bytes32 given);
    error EpochOutOfSequence(uint64 expected, uint64 given);
    error PreviousDigestMismatch(bytes32 expected, bytes32 given);
    error TooManyLeaves(uint64 leafCount, uint64 maximum);
    error EmptyEpochMismatch(bytes32 root, uint64 total);
    error SnapshotNotAfterPrevious(uint64 asOf, uint64 previous);
    error SnapshotInFuture(uint64 asOf, uint64 blockTime);
    error HeadMismatch(bytes32 head, bytes32 expected);
    error UnknownEpoch(uint64 number);

    function anchor(EpochHeader calldata header) external returns (bytes32);

    function seal(bytes32 expectedHead) external;

    function verifyInclusion(
        uint64 number,
        bytes32 commitment,
        uint64 amount,
        uint64 index,
        SumNode[] calldata siblings
    ) external view returns (bool);

    function epochAt(uint64 number) external view returns (Epoch memory);

    function head() external view returns (bytes32);

    function nextEpoch() external view returns (uint64);

    function lastAsOf() external view returns (uint64);

    function isSealed() external view returns (bool);

    function genesisEpoch() external view returns (uint64);

    function genesisDigest() external view returns (bytes32);

    function scheme() external pure returns (uint16);
}
