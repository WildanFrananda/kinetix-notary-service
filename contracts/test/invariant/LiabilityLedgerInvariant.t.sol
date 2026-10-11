// SPDX-License-Identifier: UNLICENSED
pragma solidity 0.8.37;

import {Epoch} from "../../src/Epoch.sol";
import {EpochChain} from "../../src/EpochChain.sol";
import {EpochHeader} from "../../src/EpochHeader.sol";
import {LedgerFixture} from "../support/LedgerFixture.sol";
import {LedgerHandler} from "./LedgerHandler.sol";

contract LiabilityLedgerInvariantTest is LedgerFixture {
    LedgerHandler private handler;

    function setUp() public override {
        super.setUp();
        handler = new LedgerHandler(ledger, admin, anchorer);
        targetContract(address(handler));
    }

    function invariant_NoActionHadAnOutcomeTheRulesForbid() public view {
        assertEq(handler.unexpected(), 0);
    }

    function invariant_TheHeadIsTheLastAnchoredDigest() public view {
        uint256 count = handler.anchoredCount();
        if (count == 0) {
            assertEq(ledger.head(), ledger.genesisDigest());
            return;
        }
        assertEq(ledger.head(), handler.digestAt(count - 1));
    }

    function invariant_TheNextEpochCountsEveryAnchor() public view {
        assertEq(ledger.nextEpoch(), ledger.genesisEpoch() + handler.anchoredCount());
    }

    function invariant_AnchoredEpochsNeverChange() public view {
        uint256 count = handler.anchoredCount();
        for (uint256 i = 0; i < count; i++) {
            EpochHeader memory header = handler.headerAt(i);
            Epoch memory stored = ledger.epochAt(header.number);
            assertEq(stored.root, header.root);
            assertEq(stored.total, header.total);
            assertEq(stored.leafCount, header.leafCount);
            assertEq(stored.asOf, header.asOf);
            assertEq(stored.digest, handler.digestAt(i));
        }
    }

    function invariant_EveryDigestRecomputesFromItsPredecessor() public view {
        uint256 count = handler.anchoredCount();
        bytes32 previous = ledger.genesisDigest();
        for (uint256 i = 0; i < count; i++) {
            EpochHeader memory header = handler.headerAt(i);
            assertEq(header.previousDigest, previous);
            assertEq(EpochChain.digest(header), handler.digestAt(i));
            previous = handler.digestAt(i);
        }
    }

    function invariant_SnapshotsStrictlyAdvance() public view {
        uint256 count = handler.anchoredCount();
        for (uint256 i = 1; i < count; i++) {
            assertGt(handler.headerAt(i).asOf, handler.headerAt(i - 1).asOf);
        }
        if (count > 0) {
            assertEq(ledger.lastAsOf(), handler.headerAt(count - 1).asOf);
        }
    }

    function invariant_SealingIsPermanent() public view {
        assertEq(ledger.isSealed(), handler.sealedByAdmin());
    }
}
