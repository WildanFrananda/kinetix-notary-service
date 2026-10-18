// SPDX-License-Identifier: MIT
pragma solidity 0.8.37;

import {Batch} from "../../src/evidence/Batch.sol";
import {BatchChain} from "../../src/evidence/BatchChain.sol";
import {BatchHeader} from "../../src/evidence/BatchHeader.sol";
import {StreamState} from "../../src/evidence/StreamState.sol";
import {EvidenceFixture} from "../support/EvidenceFixture.sol";
import {EvidenceHandler} from "./EvidenceHandler.sol";

contract EvidenceRegistryInvariantTest is EvidenceFixture {
    EvidenceHandler private handler;

    function setUp() public override {
        super.setUp();
        handler = new EvidenceHandler(registry, admin, anchorer);
        targetContract(address(handler));
    }

    function invariant_NoActionHadAnOutcomeTheRulesForbid() public view {
        assertEq(handler.unexpected(), 0);
    }

    function invariant_EachStreamsHeadIsItsLastAnchoredDigest() public view {
        for (uint8 stream = 1; stream <= 3; stream++) {
            uint256 count = handler.anchoredCount(stream);
            StreamState memory state = registry.stateOf(stream);
            assertEq(state.nextBatch, count);
            if (count == 0) {
                assertEq(state.head, bytes32(0));
            } else {
                assertEq(state.head, handler.digestAt(stream, count - 1));
            }
        }
    }

    function invariant_EveryBatchChainsToItsPredecessorAndNeverChanges() public view {
        for (uint8 stream = 1; stream <= 3; stream++) {
            bytes32 previous = bytes32(0);
            uint256 count = handler.anchoredCount(stream);
            for (uint256 i = 0; i < count; i++) {
                BatchHeader memory header = handler.headerAt(stream, i);
                Batch memory stored = registry.batchAt(stream, header.number);
                assertEq(header.previousDigest, previous);
                assertEq(BatchChain.digest(header), handler.digestAt(stream, i));
                assertEq(stored.digest, handler.digestAt(stream, i));
                assertEq(stored.root, header.root);
                assertEq(stored.leafCount, header.leafCount);
                previous = stored.digest;
            }
        }
    }

    function invariant_FeedPositionsAndCloseTimesStrictlyAdvance() public view {
        for (uint8 stream = 1; stream <= 3; stream++) {
            uint256 count = handler.anchoredCount(stream);
            for (uint256 i = 1; i < count; i++) {
                BatchHeader memory before = handler.headerAt(stream, i - 1);
                BatchHeader memory later = handler.headerAt(stream, i);
                assertGt(later.closedAt, before.closedAt);
                assertTrue(
                    later.throughTransaction > before.throughTransaction
                        || (later.throughTransaction == before.throughTransaction
                            && later.throughSequence > before.throughSequence)
                );
            }
        }
    }

    function invariant_SealedMeansSealed() public view {
        assertEq(registry.isSealed(), handler.sealedByAdmin());
    }
}
