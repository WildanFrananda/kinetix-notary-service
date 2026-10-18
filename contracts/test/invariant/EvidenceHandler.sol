// SPDX-License-Identifier: MIT
pragma solidity 0.8.37;

import {Test} from "forge-std/src/Test.sol";

import {BatchHeader} from "../../src/evidence/BatchHeader.sol";
import {EvidenceRegistry} from "../../src/evidence/EvidenceRegistry.sol";
import {EvidenceTree} from "../../src/evidence/EvidenceTree.sol";
import {StreamState} from "../../src/evidence/StreamState.sol";

contract EvidenceHandler is Test {
    EvidenceRegistry public immutable registry;
    address public immutable admin;
    address public immutable anchorer;

    mapping(uint8 stream => BatchHeader[]) private _headers;
    mapping(uint8 stream => bytes32[]) private _digests;
    bool public sealedByAdmin;
    uint256 public unexpected;

    constructor(EvidenceRegistry registry_, address admin_, address anchorer_) {
        registry = registry_;
        admin = admin_;
        anchorer = anchorer_;
    }

    function anchorNext(
        uint256 rawStream,
        uint256 rootSeed,
        uint64 rawLeafCount,
        uint256 rawAdvance,
        uint64 rawTransactionStep,
        uint64 rawSequenceStep
    ) external {
        uint8 stream = _stream(rawStream);
        vm.warp(vm.getBlockTimestamp() + bound(rawAdvance, 1, 2 days));
        StreamState memory state = registry.stateOf(stream);
        uint64 transactionStep = uint64(bound(rawTransactionStep, 0, 1000));
        uint64 sequenceStep = uint64(bound(rawSequenceStep, 1, 1000));
        BatchHeader memory header = BatchHeader({
            stream: stream,
            number: state.nextBatch,
            previousDigest: state.head,
            root: keccak256(abi.encode(rootSeed)),
            leafCount: uint64(bound(rawLeafCount, 1, EvidenceTree.MAX_LEAVES)),
            throughTransaction: state.throughTransaction + transactionStep,
            throughSequence: state.throughSequence + sequenceStep,
            closedAt: uint64(vm.getBlockTimestamp())
        });

        vm.prank(anchorer);
        try registry.anchor(header) returns (bytes32 digest) {
            if (sealedByAdmin) {
                unexpected++;
            }

            _headers[stream].push(header);
            _digests[stream].push(digest);
        } catch {
            if (!sealedByAdmin) {
                unexpected++;
            }
        }
    }

    function replayAnAnchoredBatch(uint256 rawStream, uint256 rawWhich, bool alter) external {
        uint8 stream = _stream(rawStream);
        uint256 count = _headers[stream].length;
        if (count == 0) {
            return;
        }

        BatchHeader memory header = _headers[stream][bound(rawWhich, 0, count - 1)];
        if (alter) {
            header.root = keccak256(abi.encode(header.root));
        }

        vm.prank(anchorer);
        try registry.anchor(header) {
            unexpected++;
        } catch {}
    }

    function seal() external {
        if (sealedByAdmin) {
            return;
        }

        uint8[] memory ids = registry.streams();
        bytes32[] memory heads = new bytes32[](ids.length);
        for (uint256 i = 0; i < ids.length; i++) {
            heads[i] = registry.stateOf(ids[i]).head;
        }

        vm.prank(admin);
        registry.seal(heads);
        sealedByAdmin = true;
    }

    function anchoredCount(uint8 stream) external view returns (uint256) {
        return _headers[stream].length;
    }

    function headerAt(uint8 stream, uint256 index) external view returns (BatchHeader memory) {
        return _headers[stream][index];
    }

    function digestAt(uint8 stream, uint256 index) external view returns (bytes32) {
        return _digests[stream][index];
    }

    function _stream(uint256 raw) private pure returns (uint8) {
        return uint8(bound(raw, 1, 3));
    }
}
