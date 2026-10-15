// SPDX-License-Identifier: MIT
pragma solidity 0.8.37;

import {IAccessControl} from "@openzeppelin-contracts/access/IAccessControl.sol";

import {LiabilityLedger} from "../../src/LiabilityLedger.sol";
import {StatementCourt} from "../../src/court/StatementCourt.sol";
import {Generation} from "../../src/directory/Generation.sol";
import {ILedgerDirectory} from "../../src/directory/ILedgerDirectory.sol";
import {WithholdingCourt} from "../../src/withholding/WithholdingCourt.sol";
import {DirectoryFixture} from "../support/DirectoryFixture.sol";

contract LedgerDirectoryFuzzTest is DirectoryFixture {
    function testFuzz_EveryEpochBelongsToTheLatestGenerationThatStartsAtOrBeforeIt(
        uint256 seed,
        uint256 rawGenerations
    ) public {
        uint256 generations = bound(rawGenerations, 1, 5);
        LiabilityLedger currentLedger = ledger;
        for (uint256 g = 0; g < generations; g++) {
            uint256 epochs = uint256(keccak256(abi.encode(seed, g))) % 4;
            for (uint256 e = 0; e < epochs; e++) {
                _anchorInto(currentLedger, uint256(keccak256(abi.encode(seed, g, e))));
            }

            if (g + 1 < generations) {
                currentLedger = _succeed(currentLedger);
            }
        }

        uint256 count = directory.generationCount();
        assertEq(count, generations);
        for (uint64 epoch = 0; epoch <= currentLedger.nextEpoch() + 1; epoch++) {
            (uint256 index, Generation memory found) = directory.generationOf(epoch);
            assertEq(abi.encode(found), abi.encode(directory.generation(index)));
            assertLe(found.genesisEpoch, epoch);
            if (index + 1 < count) {
                assertGt(directory.generation(index + 1).genesisEpoch, epoch);
            }
        }
    }

    function testFuzz_OnlyTheSealedHeadAndNextEpochContinueTheLedger(uint64 epoch, bytes32 digest) public {
        _anchorInto(ledger, 1);
        _seal(ledger);
        bytes32 head = ledger.head();
        vm.assume((epoch == 0) == (digest == bytes32(0)));
        vm.assume(epoch != 1 || digest != head);
        LiabilityLedger next = new LiabilityLedger(admin, ADMIN_DELAY, anchorer, epoch, digest, ledger.lastAsOf());
        StatementCourt nextCourt = new StatementCourt(admin, ADMIN_DELAY, next, signer, SIGNER_DELAY);
        WithholdingCourt nextWithholding = new WithholdingCourt(nextCourt, WINDOW, PERIOD);

        vm.expectRevert(abi.encodeWithSelector(ILedgerDirectory.GenesisMismatch.selector, 1, head, epoch, digest));
        vm.prank(admin);
        directory.registerSuccessor(next, nextCourt, nextWithholding);
    }

    function testFuzz_NobodyButTheAdminRegistersOrProposes(address caller) public {
        vm.assume(caller != admin);
        _seal(ledger);
        (LiabilityLedger next, StatementCourt nextCourt, WithholdingCourt nextWithholding) = _successorOf(ledger);
        bytes memory refusal = abi.encodeWithSelector(
            IAccessControl.AccessControlUnauthorizedAccount.selector, caller, directory.DEFAULT_ADMIN_ROLE()
        );

        vm.expectRevert(refusal);
        vm.prank(caller);
        directory.registerSuccessor(next, nextCourt, nextWithholding);

        vm.expectRevert(refusal);
        vm.prank(caller);
        directory.proposeForeignSuccessor(FOREIGN_CHAIN, address(0xF1), address(0xF2), address(0xF3), 1);
    }

    function testFuzz_AForeignSuccessorCountsExactlyOnceTheDelayHasPassed(uint256 rawWait) public {
        _seal(ledger);
        _propose(FOREIGN_CHAIN);
        uint64 activatesAt = _now() + FOREIGN_DELAY;
        uint64 wait = uint64(bound(rawWait, 0, 2 * FOREIGN_DELAY));
        vm.warp(_now() + wait);

        if (wait < FOREIGN_DELAY) {
            vm.expectRevert(abi.encodeWithSelector(ILedgerDirectory.SuccessorTooEarly.selector, activatesAt, _now()));
        }

        directory.activateForeignSuccessor();

        assertEq(directory.generationCount(), wait < FOREIGN_DELAY ? 1 : 2);
    }
}
