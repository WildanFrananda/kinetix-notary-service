// SPDX-License-Identifier: MIT
pragma solidity 0.8.37;

import {Epoch} from "../../src/Epoch.sol";
import {EpochHeader} from "../../src/EpochHeader.sol";
import {ConfidentialEpochChain} from "../../src/confidential/ConfidentialEpochChain.sol";
import {ConfidentialTree} from "../../src/confidential/ConfidentialTree.sol";
import {ConfidentialLedgerFixture} from "../support/ConfidentialLedgerFixture.sol";
import {ConfidentialLedgerHandler} from "./ConfidentialLedgerHandler.sol";

contract LiabilityLedgerV2InvariantTest is ConfidentialLedgerFixture {
    ConfidentialLedgerHandler private handler;

    function setUp() public override {
        super.setUp();
        handler = new ConfidentialLedgerHandler(ledger, verifier, admin, anchorer);
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

    function invariant_EveryAnchoredEpochIsAProvenStatementWithItsTotalHidden() public view {
        uint256 count = handler.anchoredCount();
        for (uint256 i = 0; i < count; i++) {
            EpochHeader memory header = handler.headerAt(i);
            Epoch memory stored = ledger.epochAt(header.number);
            assertEq(stored.total, 0);
            assertEq(stored.root, header.root);
            assertEq(stored.leafCount, header.leafCount);
            assertTrue(ledger.isProven(header.number, header.root, header.leafCount));
            assertEq(
                ConfidentialTree.bindingRoot(ledger.treeRootOf(header.number), ledger.totalCommitmentOf(header.number)),
                header.root
            );
        }
    }

    function invariant_EveryDigestRecomputesFromItsPredecessor() public view {
        uint256 count = handler.anchoredCount();
        bytes32 previous = ledger.genesisDigest();
        for (uint256 i = 0; i < count; i++) {
            EpochHeader memory header = handler.headerAt(i);
            assertEq(header.previousDigest, previous);
            assertEq(ConfidentialEpochChain.digest(header), handler.digestAt(i));
            previous = handler.digestAt(i);
        }
    }

    function invariant_SealingIsPermanent() public view {
        assertEq(ledger.isSealed(), handler.sealedByAdmin());
    }
}
