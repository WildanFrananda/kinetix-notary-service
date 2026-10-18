// SPDX-License-Identifier: MIT
pragma solidity 0.8.37;

import {Batch} from "./Batch.sol";
import {BatchHeader} from "./BatchHeader.sol";
import {StreamState} from "./StreamState.sol";

interface IEvidenceRegistry {
    event BatchAnchored(uint8 indexed stream, uint64 indexed number, bytes32 digest, BatchHeader header);
    event StreamSealed(
        uint8 indexed stream,
        bytes32 head,
        uint64 nextBatch,
        uint64 lastClosedAt,
        uint64 throughTransaction,
        uint64 throughSequence
    );

    error ZeroAnchorer();
    error NoStreams();
    error TooManyStreams(uint256 streams, uint256 maximum);
    error AnchorerNotGranted(address anchorer);
    error StreamsNotAscending(uint8 previous, uint8 stream);
    error InvalidGenesis(uint8 stream, uint64 batch, bytes32 digest);
    error RegistryIsSealed();
    error UnknownStream(uint8 stream);
    error UnknownBatch(uint8 stream, uint64 number);
    error BatchBeforeGenesis(uint8 stream, uint64 genesisBatch, uint64 number);
    error BatchAlreadyAnchored(uint8 stream, uint64 number, bytes32 digest);
    error BatchConflict(uint8 stream, uint64 number, bytes32 anchored, bytes32 offered);
    error BatchOutOfSequence(uint8 stream, uint64 expected, uint64 offered);
    error PreviousDigestMismatch(uint8 stream, bytes32 head, bytes32 offered);
    error EmptyBatch(uint8 stream);
    error TooManyLeaves(uint64 leafCount, uint64 maximum);
    error CloseNotAfterPrevious(uint8 stream, uint64 closedAt, uint64 lastClosedAt);
    error CloseInFuture(uint8 stream, uint64 closedAt, uint64 blockTime);
    error PositionNotAfterPrevious(
        uint8 stream,
        uint64 throughTransaction,
        uint64 throughSequence,
        uint64 previousTransaction,
        uint64 previousSequence
    );
    error HeadCountMismatch(uint256 streams, uint256 offered);
    error HeadMismatch(uint8 stream, bytes32 head, bytes32 expected);

    function anchor(BatchHeader calldata header) external returns (bytes32);

    function seal(bytes32[] calldata expectedHeads) external;

    function verifyInclusion(
        uint8 stream,
        uint64 number,
        bytes32 salt,
        bytes32 recordDigest,
        uint64 index,
        bytes32[] calldata siblings
    ) external view returns (bool);

    function batchAt(uint8 stream, uint64 number) external view returns (Batch memory);

    function stateOf(uint8 stream) external view returns (StreamState memory);

    function streams() external view returns (uint8[] memory);

    function isSealed() external view returns (bool);

    function scheme() external pure returns (uint16);
}
