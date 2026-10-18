// SPDX-License-Identifier: MIT
pragma solidity 0.8.37;

import {
    AccessControlDefaultAdminRules
} from "@openzeppelin-contracts/access/extensions/AccessControlDefaultAdminRules.sol";
import {SafeCast} from "@openzeppelin-contracts/utils/math/SafeCast.sol";

import {Batch} from "./Batch.sol";
import {BatchChain} from "./BatchChain.sol";
import {BatchHeader} from "./BatchHeader.sol";
import {EvidenceTree} from "./EvidenceTree.sol";
import {IEvidenceRegistry} from "./IEvidenceRegistry.sol";
import {StreamGenesis} from "./StreamGenesis.sol";
import {StreamState} from "./StreamState.sol";

contract EvidenceRegistry is IEvidenceRegistry, AccessControlDefaultAdminRules {
    bytes32 public constant ANCHORER_ROLE = keccak256("kinetix.notary.EVIDENCE_ANCHORER");
    uint256 public constant MAX_STREAMS = 32;

    bytes32 private immutable _streamIds;
    uint256 private immutable _streamCount;
    mapping(uint8 stream => StreamState) private _states;
    mapping(uint8 stream => mapping(uint64 number => Batch)) private _batches;
    bool private _sealed;

    constructor(
        address admin,
        uint48 adminTransferDelay,
        address anchorer,
        StreamGenesis[] memory genesis
    ) AccessControlDefaultAdminRules(adminTransferDelay, admin) {
        if (anchorer == address(0)) {
            revert ZeroAnchorer();
        }

        if (genesis.length == 0) {
            revert NoStreams();
        }

        if (genesis.length > MAX_STREAMS) {
            revert TooManyStreams(genesis.length, MAX_STREAMS);
        }

        uint256 invalid = _firstInvalidGenesis(genesis);
        if (invalid < genesis.length) {
            StreamGenesis memory offending = genesis[invalid];
            uint8 previous = invalid == 0 ? 0 : genesis[invalid - 1].stream;
            if (offending.stream <= previous) {
                revert StreamsNotAscending(previous, offending.stream);
            }

            revert InvalidGenesis(offending.stream, offending.batch, offending.digest);
        }

        bytes32 packed = bytes32(0);
        for (uint256 i = 0; i < genesis.length; i++) {
            StreamGenesis memory start = genesis[i];
            packed |= bytes32(uint256(start.stream) << (8 * (31 - i)));
            _states[start.stream] = StreamState({
                known: true,
                head: start.digest,
                genesisBatch: start.batch,
                nextBatch: start.batch,
                lastClosedAt: start.closedAt,
                throughTransaction: start.throughTransaction,
                throughSequence: start.throughSequence
            });
        }

        _streamIds = packed;
        _streamCount = genesis.length;
        if (!_grantRole(ANCHORER_ROLE, anchorer)) {
            revert AnchorerNotGranted(anchorer);
        }
    }

    function anchor(BatchHeader calldata header) external onlyRole(ANCHORER_ROLE) returns (bytes32) {
        if (_sealed) {
            revert RegistryIsSealed();
        }

        StreamState storage state = _known(header.stream);
        bytes32 digest = BatchChain.digest(header);
        if (header.number < state.genesisBatch) {
            revert BatchBeforeGenesis(header.stream, state.genesisBatch, header.number);
        }

        if (header.number < state.nextBatch) {
            bytes32 anchored = _batches[header.stream][header.number].digest;
            if (anchored == digest) {
                revert BatchAlreadyAnchored(header.stream, header.number, digest);
            }

            revert BatchConflict(header.stream, header.number, anchored, digest);
        }

        if (header.number > state.nextBatch) {
            revert BatchOutOfSequence(header.stream, state.nextBatch, header.number);
        }

        if (header.previousDigest != state.head) {
            revert PreviousDigestMismatch(header.stream, state.head, header.previousDigest);
        }

        _requireWellFormed(header, state);

        _batches[header.stream][header.number] = Batch({
            root: header.root,
            digest: digest,
            leafCount: header.leafCount,
            throughTransaction: header.throughTransaction,
            throughSequence: header.throughSequence,
            closedAt: header.closedAt,
            anchoredAt: SafeCast.toUint64(block.timestamp)
        });
        state.head = digest;
        state.nextBatch = header.number + 1;
        state.lastClosedAt = header.closedAt;
        state.throughTransaction = header.throughTransaction;
        state.throughSequence = header.throughSequence;

        emit BatchAnchored(header.stream, header.number, digest, header);
        return digest;
    }

    function seal(bytes32[] calldata expectedHeads) external onlyRole(DEFAULT_ADMIN_ROLE) {
        if (_sealed) {
            revert RegistryIsSealed();
        }

        uint8[] memory streamIds = _streamList();
        if (expectedHeads.length != streamIds.length) {
            revert HeadCountMismatch(streamIds.length, expectedHeads.length);
        }

        uint256 mismatch = _firstHeadMismatch(streamIds, expectedHeads);
        if (mismatch < streamIds.length) {
            revert HeadMismatch(streamIds[mismatch], _states[streamIds[mismatch]].head, expectedHeads[mismatch]);
        }

        _sealed = true;
        for (uint256 i = 0; i < streamIds.length; i++) {
            StreamState memory state = _states[streamIds[i]];
            emit StreamSealed(
                streamIds[i],
                state.head,
                state.nextBatch,
                state.lastClosedAt,
                state.throughTransaction,
                state.throughSequence
            );
        }
    }

    function verifyInclusion(
        uint8 stream,
        uint64 number,
        bytes32 salt,
        bytes32 recordDigest,
        uint64 index,
        bytes32[] calldata siblings
    ) external view returns (bool) {
        Batch storage anchored = _anchored(stream, number);
        if (index >= anchored.leafCount) {
            return false;
        }

        if (siblings.length != EvidenceTree.depthOf(anchored.leafCount)) {
            return false;
        }

        bytes32 leaf = EvidenceTree.leaf(stream, salt, recordDigest);
        return EvidenceTree.computeRoot(leaf, index, siblings) == anchored.root;
    }

    function batchAt(uint8 stream, uint64 number) external view returns (Batch memory) {
        return _anchored(stream, number);
    }

    function stateOf(uint8 stream) external view returns (StreamState memory) {
        return _known(stream);
    }

    function streams() external view returns (uint8[] memory) {
        return _streamList();
    }

    function isSealed() external view returns (bool) {
        return _sealed;
    }

    function scheme() external pure returns (uint16) {
        return EvidenceTree.SCHEME;
    }

    function supportsInterface(bytes4 interfaceId) public view override returns (bool) {
        return interfaceId == type(IEvidenceRegistry).interfaceId || super.supportsInterface(interfaceId);
    }

    function _streamList() private view returns (uint8[] memory) {
        uint8[] memory ids = new uint8[](_streamCount);
        for (uint256 i = 0; i < _streamCount; i++) {
            ids[i] = uint8(_streamIds[i]);
        }

        return ids;
    }

    function _firstInvalidGenesis(StreamGenesis[] memory genesis) private pure returns (uint256) {
        uint8 previous = 0;
        for (uint256 i = 0; i < genesis.length; i++) {
            StreamGenesis memory start = genesis[i];
            if (start.stream <= previous || (start.batch == 0) != (start.digest == bytes32(0))) {
                return i;
            }

            previous = start.stream;
        }

        return genesis.length;
    }

    function _firstHeadMismatch(
        uint8[] memory streamIds,
        bytes32[] calldata expectedHeads
    ) private view returns (uint256) {
        for (uint256 i = 0; i < streamIds.length; i++) {
            if (_states[streamIds[i]].head != expectedHeads[i]) {
                return i;
            }
        }

        return streamIds.length;
    }

    function _known(uint8 stream) private view returns (StreamState storage) {
        StreamState storage state = _states[stream];
        if (!state.known) {
            revert UnknownStream(stream);
        }

        return state;
    }

    function _anchored(uint8 stream, uint64 number) private view returns (Batch storage) {
        StreamState storage state = _known(stream);
        if (number < state.genesisBatch || number >= state.nextBatch) {
            revert UnknownBatch(stream, number);
        }

        return _batches[stream][number];
    }

    function _requireWellFormed(BatchHeader calldata header, StreamState storage state) private view {
        if (header.leafCount == 0) {
            revert EmptyBatch(header.stream);
        }

        if (header.leafCount > EvidenceTree.MAX_LEAVES) {
            revert TooManyLeaves(header.leafCount, EvidenceTree.MAX_LEAVES);
        }

        if (header.closedAt <= state.lastClosedAt) {
            revert CloseNotAfterPrevious(header.stream, header.closedAt, state.lastClosedAt);
        }

        if (header.closedAt > block.timestamp) {
            revert CloseInFuture(header.stream, header.closedAt, SafeCast.toUint64(block.timestamp));
        }

        bool laterTransaction = header.throughTransaction > state.throughTransaction;
        bool sameTransaction = header.throughTransaction == state.throughTransaction;
        if (!laterTransaction && !(sameTransaction && header.throughSequence > state.throughSequence)) {
            revert PositionNotAfterPrevious(
                header.stream,
                header.throughTransaction,
                header.throughSequence,
                state.throughTransaction,
                state.throughSequence
            );
        }
    }
}
