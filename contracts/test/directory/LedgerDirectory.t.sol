// SPDX-License-Identifier: MIT
pragma solidity 0.8.37;

import {IAccessControl} from "@openzeppelin-contracts/access/IAccessControl.sol";

import {ILiabilityLedger} from "../../src/ILiabilityLedger.sol";
import {LiabilityLedger} from "../../src/LiabilityLedger.sol";
import {IStatementCourt} from "../../src/court/IStatementCourt.sol";
import {StatementCourt} from "../../src/court/StatementCourt.sol";
import {Generation} from "../../src/directory/Generation.sol";
import {ILedgerDirectory} from "../../src/directory/ILedgerDirectory.sol";
import {LedgerDirectory} from "../../src/directory/LedgerDirectory.sol";
import {PendingGeneration} from "../../src/directory/PendingGeneration.sol";
import {IWithholdingCourt} from "../../src/withholding/IWithholdingCourt.sol";
import {WithholdingCourt} from "../../src/withholding/WithholdingCourt.sol";
import {DirectoryFixture} from "../support/DirectoryFixture.sol";

contract LedgerDirectoryTest is DirectoryFixture {
    function test_GenerationZeroIsTheLedgerAndCourtsItWasGiven() public view {
        Generation memory zero = directory.generation(0);

        assertEq(directory.generationCount(), 1);
        assertEq(zero.chainId, block.chainid);
        assertEq(zero.ledger, address(ledger));
        assertEq(zero.court, address(court));
        assertEq(zero.withholding, address(withholding));
        assertEq(zero.genesisEpoch, 0);
        assertEq(zero.genesisDigest, bytes32(0));
        assertEq(zero.scheme, 1);
        assertEq(zero.registeredAt, START);
        assertEq(abi.encode(directory.current()), abi.encode(zero));
        assertEq(directory.defaultAdmin(), admin);
        assertEq(directory.defaultAdminDelay(), ADMIN_DELAY);
        assertEq(directory.foreignSuccessorDelay(), FOREIGN_DELAY);
        assertTrue(directory.supportsInterface(type(ILedgerDirectory).interfaceId));
    }

    function test_GenerationZeroTakesItsGenesisFromTheLedger() public {
        _anchorInto(ledger, 1);
        _seal(ledger);
        (LiabilityLedger next, StatementCourt nextCourt, WithholdingCourt nextWithholding) = _successorOf(ledger);

        LedgerDirectory later = new LedgerDirectory(admin, ADMIN_DELAY, next, nextCourt, nextWithholding, FOREIGN_DELAY);

        assertEq(later.generation(0).genesisEpoch, 1);
        assertEq(later.generation(0).genesisDigest, ledger.head());
    }

    function test_RefusesAZeroForeignDelay() public {
        vm.expectRevert(ILedgerDirectory.ZeroDelay.selector);
        new LedgerDirectory(admin, ADMIN_DELAY, ledger, court, withholding, 0);
    }

    function test_RefusesWhatIsNotALedger() public {
        vm.expectRevert(abi.encodeWithSelector(ILedgerDirectory.NotALedger.selector, stranger));
        new LedgerDirectory(admin, ADMIN_DELAY, ILiabilityLedger(stranger), court, withholding, FOREIGN_DELAY);
    }

    function test_RefusesWhatIsNotAStatementCourt() public {
        IStatementCourt notACourt = IStatementCourt(address(ledger));

        vm.expectRevert(abi.encodeWithSelector(ILedgerDirectory.NotAStatementCourt.selector, address(ledger)));
        new LedgerDirectory(admin, ADMIN_DELAY, ledger, notACourt, withholding, FOREIGN_DELAY);
    }

    function test_RefusesWhatIsNotAWithholdingCourt() public {
        IWithholdingCourt notAWithholdingCourt = IWithholdingCourt(address(court));

        vm.expectRevert(abi.encodeWithSelector(ILedgerDirectory.NotAWithholdingCourt.selector, address(court)));
        new LedgerDirectory(admin, ADMIN_DELAY, ledger, court, notAWithholdingCourt, FOREIGN_DELAY);
    }

    function test_RefusesACourtThatJudgesAnotherLedger() public {
        (, StatementCourt otherCourt, WithholdingCourt otherWithholding) = _successorOf(ledger);

        vm.expectRevert(
            abi.encodeWithSelector(ILedgerDirectory.CourtOfAnotherLedger.selector, address(otherCourt), address(ledger))
        );
        new LedgerDirectory(admin, ADMIN_DELAY, ledger, otherCourt, otherWithholding, FOREIGN_DELAY);
    }

    function test_RefusesAWithholdingCourtBesideAnotherCourt() public {
        StatementCourt secondCourt = new StatementCourt(admin, ADMIN_DELAY, ledger, signer, SIGNER_DELAY);

        vm.expectRevert(
            abi.encodeWithSelector(
                ILedgerDirectory.WithholdingOfAnotherCourt.selector, address(withholding), address(secondCourt)
            )
        );
        new LedgerDirectory(admin, ADMIN_DELAY, ledger, secondCourt, withholding, FOREIGN_DELAY);
    }

    function test_RegistersASuccessorThatStartsAtTheSealedHead() public {
        _anchorInto(ledger, 1);
        _anchorInto(ledger, 2);
        _seal(ledger);
        (LiabilityLedger next, StatementCourt nextCourt, WithholdingCourt nextWithholding) = _successorOf(ledger);
        bytes32 head = ledger.head();

        vm.expectEmit(address(directory));
        emit ILedgerDirectory.GenerationRegistered(
            1,
            uint64(block.chainid),
            address(next),
            address(nextCourt),
            address(nextWithholding),
            2,
            head,
            1
        );
        vm.prank(admin);
        uint256 index = directory.registerSuccessor(next, nextCourt, nextWithholding);

        Generation memory one = directory.generation(1);
        assertEq(index, 1);
        assertEq(directory.generationCount(), 2);
        assertEq(one.chainId, block.chainid);
        assertEq(one.ledger, address(next));
        assertEq(one.court, address(nextCourt));
        assertEq(one.withholding, address(nextWithholding));
        assertEq(one.genesisEpoch, 2);
        assertEq(one.genesisDigest, head);
        assertEq(one.scheme, 1);
        assertEq(one.registeredAt, _now());
        assertEq(abi.encode(directory.current()), abi.encode(one));
    }

    function test_OnlyTheAdminRegisters() public {
        _seal(ledger);
        (LiabilityLedger next, StatementCourt nextCourt, WithholdingCourt nextWithholding) = _successorOf(ledger);
        bytes32 role = directory.DEFAULT_ADMIN_ROLE();

        vm.expectRevert(
            abi.encodeWithSelector(IAccessControl.AccessControlUnauthorizedAccount.selector, stranger, role)
        );
        vm.prank(stranger);
        directory.registerSuccessor(next, nextCourt, nextWithholding);
    }

    function test_RefusesASuccessorWhileThePredecessorStillAnchors() public {
        _anchorInto(ledger, 1);
        (LiabilityLedger next, StatementCourt nextCourt, WithholdingCourt nextWithholding) = _successorOf(ledger);

        vm.expectRevert(abi.encodeWithSelector(ILedgerDirectory.PredecessorNotSealed.selector, address(ledger)));
        vm.prank(admin);
        directory.registerSuccessor(next, nextCourt, nextWithholding);
    }

    function test_RefusesASuccessorThatIsSealedAlready() public {
        _anchorInto(ledger, 1);
        _seal(ledger);
        (LiabilityLedger next, StatementCourt nextCourt, WithholdingCourt nextWithholding) = _successorOf(ledger);
        _seal(next);

        vm.expectRevert(abi.encodeWithSelector(ILedgerDirectory.SuccessorSealed.selector, address(next)));
        vm.prank(admin);
        directory.registerSuccessor(next, nextCourt, nextWithholding);
    }

    function test_RefusesTheSealedPredecessorAsItsOwnSuccessor() public {
        _seal(ledger);

        vm.expectRevert(abi.encodeWithSelector(ILedgerDirectory.SuccessorSealed.selector, address(ledger)));
        vm.prank(admin);
        directory.registerSuccessor(ledger, court, withholding);
    }

    function test_RefusesASuccessorThatStartsAtAnotherEpoch() public {
        _anchorInto(ledger, 1);
        _seal(ledger);
        bytes32 head = ledger.head();
        LiabilityLedger skipping = new LiabilityLedger(admin, ADMIN_DELAY, anchorer, 2, head, ledger.lastAsOf());
        StatementCourt skippingCourt = new StatementCourt(admin, ADMIN_DELAY, skipping, signer, SIGNER_DELAY);
        WithholdingCourt skippingWithholding = new WithholdingCourt(skippingCourt, WINDOW, PERIOD);

        vm.expectRevert(abi.encodeWithSelector(ILedgerDirectory.GenesisMismatch.selector, 1, head, 2, head));
        vm.prank(admin);
        directory.registerSuccessor(skipping, skippingCourt, skippingWithholding);
    }

    function test_RefusesASuccessorThatStartsFromAnotherDigest() public {
        _anchorInto(ledger, 1);
        _seal(ledger);
        bytes32 head = ledger.head();
        bytes32 forged = keccak256("not the sealed head");
        LiabilityLedger forking = new LiabilityLedger(admin, ADMIN_DELAY, anchorer, 1, forged, ledger.lastAsOf());
        StatementCourt forkingCourt = new StatementCourt(admin, ADMIN_DELAY, forking, signer, SIGNER_DELAY);
        WithholdingCourt forkingWithholding = new WithholdingCourt(forkingCourt, WINDOW, PERIOD);

        vm.expectRevert(abi.encodeWithSelector(ILedgerDirectory.GenesisMismatch.selector, 1, head, 1, forged));
        vm.prank(admin);
        directory.registerSuccessor(forking, forkingCourt, forkingWithholding);
    }

    function test_RefusesASuccessorJudgedByThePredecessorsCourt() public {
        _anchorInto(ledger, 1);
        _seal(ledger);
        (LiabilityLedger next,,) = _successorOf(ledger);

        vm.expectRevert(
            abi.encodeWithSelector(ILedgerDirectory.CourtOfAnotherLedger.selector, address(court), address(next))
        );
        vm.prank(admin);
        directory.registerSuccessor(next, court, withholding);
    }

    function test_EachEpochBelongsToTheGenerationThatHoldsIt() public {
        _anchorInto(ledger, 1);
        _anchorInto(ledger, 2);
        LiabilityLedger second = _succeed(ledger);
        _anchorInto(second, 3);
        LiabilityLedger third = _succeed(second);

        assertEq(_indexOf(0), 0);
        assertEq(_indexOf(1), 0);
        assertEq(_indexOf(2), 1);
        assertEq(_indexOf(3), 2);
        assertEq(_indexOf(type(uint64).max), 2);
        (, Generation memory found) = directory.generationOf(3);
        assertEq(found.ledger, address(third));
    }

    function test_APredecessorSealedBeforeItAnchoredOwnsNoEpoch() public {
        _anchorInto(ledger, 1);
        LiabilityLedger empty = _succeed(ledger);
        _succeed(empty);

        assertEq(_indexOf(0), 0);
        assertEq(_indexOf(1), 2);
    }

    function test_AnEpochBeforeTheFirstGenerationBelongsToNone() public {
        _anchorInto(ledger, 1);
        _anchorInto(ledger, 2);
        _seal(ledger);
        (LiabilityLedger next, StatementCourt nextCourt, WithholdingCourt nextWithholding) = _successorOf(ledger);
        LedgerDirectory later = new LedgerDirectory(admin, ADMIN_DELAY, next, nextCourt, nextWithholding, FOREIGN_DELAY);

        vm.expectRevert(abi.encodeWithSelector(ILedgerDirectory.EpochBeforeDirectory.selector, 1, 2));
        later.generationOf(1);
    }

    function test_AnUnknownGenerationIsRefusedNotZero() public {
        vm.expectRevert(abi.encodeWithSelector(ILedgerDirectory.UnknownGeneration.selector, 1, 1));
        directory.generation(1);
    }

    function test_AForeignSuccessorWaitsTheDelayThenCounts() public {
        _anchorInto(ledger, 1);
        _seal(ledger);
        uint64 activatesAt = _now() + FOREIGN_DELAY;

        vm.expectEmit(address(directory));
        emit ILedgerDirectory.ForeignSuccessorProposed(
            FOREIGN_CHAIN,
            address(0xF1),
            address(0xF2),
            address(0xF3),
            1,
            activatesAt
        );
        _propose(FOREIGN_CHAIN);

        PendingGeneration memory pending = directory.pendingSuccessor();
        assertEq(pending.activatesAt, activatesAt);
        assertEq(pending.generation.chainId, FOREIGN_CHAIN);
        assertEq(pending.generation.genesisEpoch, 1);
        assertEq(pending.generation.genesisDigest, ledger.head());
        assertEq(directory.generationCount(), 1);

        vm.warp(activatesAt - 1);
        vm.expectRevert(
            abi.encodeWithSelector(ILedgerDirectory.SuccessorTooEarly.selector, activatesAt, activatesAt - 1)
        );
        directory.activateForeignSuccessor();

        vm.warp(activatesAt);
        vm.prank(stranger);
        uint256 index = directory.activateForeignSuccessor();

        Generation memory moved = directory.generation(1);
        assertEq(index, 1);
        assertEq(moved.chainId, FOREIGN_CHAIN);
        assertEq(moved.ledger, address(0xF1));
        assertEq(moved.court, address(0xF2));
        assertEq(moved.withholding, address(0xF3));
        assertEq(moved.genesisEpoch, 1);
        assertEq(moved.registeredAt, activatesAt);
        assertEq(directory.pendingSuccessor().activatesAt, 0);
        assertEq(_indexOf(1), 1);
        assertEq(_indexOf(0), 0);
    }

    function test_AfterMovingToAnotherChainThisDirectoryTakesNoMore() public {
        _seal(ledger);
        _propose(FOREIGN_CHAIN);
        vm.warp(_now() + FOREIGN_DELAY);
        directory.activateForeignSuccessor();
        (LiabilityLedger next, StatementCourt nextCourt, WithholdingCourt nextWithholding) = _successorOf(ledger);

        vm.expectRevert(abi.encodeWithSelector(ILedgerDirectory.DirectoryMoved.selector, FOREIGN_CHAIN));
        vm.prank(admin);
        directory.registerSuccessor(next, nextCourt, nextWithholding);

        vm.expectRevert(abi.encodeWithSelector(ILedgerDirectory.DirectoryMoved.selector, FOREIGN_CHAIN));
        _propose(FOREIGN_CHAIN + 1);
    }

    function test_APendingForeignSuccessorBlocksEveryOtherRegistration() public {
        _seal(ledger);
        _propose(FOREIGN_CHAIN);
        uint64 activatesAt = _now() + FOREIGN_DELAY;
        (LiabilityLedger next, StatementCourt nextCourt, WithholdingCourt nextWithholding) = _successorOf(ledger);

        vm.expectRevert(abi.encodeWithSelector(ILedgerDirectory.SuccessorPending.selector, activatesAt));
        vm.prank(admin);
        directory.registerSuccessor(next, nextCourt, nextWithholding);

        vm.expectRevert(abi.encodeWithSelector(ILedgerDirectory.SuccessorPending.selector, activatesAt));
        _propose(FOREIGN_CHAIN + 1);
    }

    function test_TheAdminCancelsAForeignSuccessorAndMayThenRegisterOnThisChain() public {
        _seal(ledger);
        _propose(FOREIGN_CHAIN);

        vm.expectEmit(address(directory));
        emit ILedgerDirectory.ForeignSuccessorCancelled(FOREIGN_CHAIN, address(0xF1));
        vm.prank(admin);
        directory.cancelForeignSuccessor();

        assertEq(directory.pendingSuccessor().generation.chainId, 0);
        vm.warp(_now() + FOREIGN_DELAY);
        vm.expectRevert(ILedgerDirectory.NoSuccessorPending.selector);
        directory.activateForeignSuccessor();

        (LiabilityLedger next, StatementCourt nextCourt, WithholdingCourt nextWithholding) = _successorOf(ledger);
        vm.prank(admin);
        assertEq(directory.registerSuccessor(next, nextCourt, nextWithholding), 1);
    }

    function test_OnlyTheAdminProposesOrCancels() public {
        _seal(ledger);
        bytes32 role = directory.DEFAULT_ADMIN_ROLE();

        vm.expectRevert(
            abi.encodeWithSelector(IAccessControl.AccessControlUnauthorizedAccount.selector, stranger, role)
        );
        vm.prank(stranger);
        directory.proposeForeignSuccessor(FOREIGN_CHAIN, address(0xF1), address(0xF2), address(0xF3), 1);

        _propose(FOREIGN_CHAIN);
        vm.expectRevert(
            abi.encodeWithSelector(IAccessControl.AccessControlUnauthorizedAccount.selector, stranger, role)
        );
        vm.prank(stranger);
        directory.cancelForeignSuccessor();
    }

    function test_AForeignSuccessorNeedsASealedPredecessor() public {
        vm.expectRevert(abi.encodeWithSelector(ILedgerDirectory.PredecessorNotSealed.selector, address(ledger)));
        _propose(FOREIGN_CHAIN);
    }

    function test_AForeignSuccessorMustBeForeignAndComplete() public {
        _seal(ledger);
        uint64 here = uint64(block.chainid);

        vm.expectRevert(abi.encodeWithSelector(ILedgerDirectory.NotForeign.selector, here));
        _propose(here);

        vm.expectRevert(abi.encodeWithSelector(ILedgerDirectory.NotForeign.selector, 0));
        _propose(0);

        vm.startPrank(admin);
        vm.expectRevert(ILedgerDirectory.ZeroAddress.selector);
        directory.proposeForeignSuccessor(FOREIGN_CHAIN, address(0), address(0xF2), address(0xF3), 1);
        vm.expectRevert(ILedgerDirectory.ZeroAddress.selector);
        directory.proposeForeignSuccessor(FOREIGN_CHAIN, address(0xF1), address(0), address(0xF3), 1);
        vm.expectRevert(ILedgerDirectory.ZeroAddress.selector);
        directory.proposeForeignSuccessor(FOREIGN_CHAIN, address(0xF1), address(0xF2), address(0), 1);
        vm.expectRevert(ILedgerDirectory.ZeroScheme.selector);
        directory.proposeForeignSuccessor(FOREIGN_CHAIN, address(0xF1), address(0xF2), address(0xF3), 0);
        vm.stopPrank();
    }

    function test_NothingToCancelOrActivateWithoutAProposal() public {
        vm.expectRevert(ILedgerDirectory.NoSuccessorPending.selector);
        vm.prank(admin);
        directory.cancelForeignSuccessor();

        vm.expectRevert(ILedgerDirectory.NoSuccessorPending.selector);
        directory.activateForeignSuccessor();
    }

    function _indexOf(uint64 epoch) private view returns (uint256) {
        (uint256 index,) = directory.generationOf(epoch);
        return index;
    }
}
