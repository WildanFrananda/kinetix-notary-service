// SPDX-License-Identifier: MIT
pragma solidity 0.8.37;

import {LiabilityLedger} from "../../src/LiabilityLedger.sol";
import {StatementCourt} from "../../src/court/StatementCourt.sol";
import {Generation} from "../../src/directory/Generation.sol";
import {LedgerDirectory} from "../../src/directory/LedgerDirectory.sol";
import {WithholdingCourt} from "../../src/withholding/WithholdingCourt.sol";
import {WithholdingFixture} from "../support/WithholdingFixture.sol";
import {DirectoryHandler} from "./DirectoryHandler.sol";

contract LedgerDirectoryInvariantTest is WithholdingFixture {
    LedgerDirectory private directory;
    DirectoryHandler private handler;

    function setUp() public override {
        super.setUp();
        directory = new LedgerDirectory(admin, ADMIN_DELAY, ledger, court, withholding, 7 days);
        handler = new DirectoryHandler(directory, admin, anchorer, signer, stranger);
        targetContract(address(handler));
    }

    function invariant_NoActionHadAnOutcomeTheRulesForbid() public view {
        assertEq(handler.unexpected(), 0);
    }

    function invariant_TheListOnlyEverGrows() public view {
        uint256 count = directory.generationCount();

        assertEq(count, handler.recordedCount());
        for (uint256 i = 0; i < count; i++) {
            assertEq(keccak256(abi.encode(directory.generation(i))), handler.recorded(i));
        }
    }

    function invariant_EveryLocalSuccessorContinuesItsSealedPredecessor() public view {
        uint256 count = directory.generationCount();
        for (uint256 i = 1; i < count; i++) {
            Generation memory successor = directory.generation(i);
            LiabilityLedger predecessor = LiabilityLedger(directory.generation(i - 1).ledger);

            assertTrue(predecessor.isSealed());
            assertEq(successor.genesisEpoch, predecessor.nextEpoch());
            assertEq(successor.genesisDigest, predecessor.head());
            if (successor.chainId == block.chainid) {
                assertEq(address(StatementCourt(successor.court).ledger()), successor.ledger);
                assertEq(address(WithholdingCourt(successor.withholding).court()), successor.court);
            }
        }
    }

    function invariant_OnlyTheLastGenerationMayLiveOnAnotherChain() public view {
        uint256 count = directory.generationCount();
        for (uint256 i = 0; i + 1 < count; i++) {
            assertEq(directory.generation(i).chainId, block.chainid);
        }
    }

    function invariant_EveryEpochHasExactlyOneGeneration() public view {
        uint256 count = directory.generationCount();
        for (uint256 i = 0; i < count; i++) {
            Generation memory each = directory.generation(i);
            (uint256 found,) = directory.generationOf(each.genesisEpoch);
            assertGe(found, i);
            assertEq(directory.generation(found).genesisEpoch, each.genesisEpoch);
        }
    }
}
