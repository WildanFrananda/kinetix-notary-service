// SPDX-License-Identifier: UNLICENSED
pragma solidity 0.8.37;

import {IAccessControl} from "@openzeppelin-contracts/access/IAccessControl.sol";
import {IERC165} from "@openzeppelin-contracts/utils/introspection/IERC165.sol";

import {Epoch} from "../src/Epoch.sol";
import {EpochChain} from "../src/EpochChain.sol";
import {EpochHeader} from "../src/EpochHeader.sol";
import {ILiabilityLedger} from "../src/ILiabilityLedger.sol";
import {LiabilityLedger} from "../src/LiabilityLedger.sol";
import {SumNode} from "../src/SumNode.sol";
import {SumTree} from "../src/SumTree.sol";
import {LedgerFixture} from "./support/LedgerFixture.sol";
import {TreeBuilder} from "./support/TreeBuilder.sol";

contract LiabilityLedgerTest is LedgerFixture {
    function test_AnchorsTheFirstEpochOntoAnEmptyChain() public {
        SumNode[][] memory levels = _tree(0, 1, 5);
        EpochHeader memory header = _header(0, bytes32(0), levels, 5, _now() - 60);

        bytes32 digest = _anchor(header);

        assertEq(digest, EpochChain.digest(header));
        assertEq(ledger.head(), digest);
        assertEq(ledger.nextEpoch(), 1);
        assertEq(ledger.lastAsOf(), header.asOf);
        Epoch memory stored = ledger.epochAt(0);
        assertEq(stored.root, header.root);
        assertEq(stored.digest, digest);
        assertEq(stored.total, header.total);
        assertEq(stored.leafCount, 5);
        assertEq(stored.asOf, header.asOf);
        assertEq(stored.anchoredAt, _now());
    }

    function test_ChainsEachEpochToThePreviousDigest() public {
        bytes32 first = _anchor(_header(0, bytes32(0), _tree(0, 1, 3), 3, _now() - 60));
        vm.warp(vm.getBlockTimestamp() + 1 hours);
        EpochHeader memory second = _header(1, first, _tree(1, 2, 4), 4, _now() - 60);

        bytes32 digest = _anchor(second);

        assertEq(digest, EpochChain.digest(second));
        assertEq(ledger.head(), digest);
        assertEq(ledger.nextEpoch(), 2);
    }

    function test_EmitsEpochAnchored() public {
        EpochHeader memory header = _header(0, bytes32(0), _tree(0, 1, 2), 2, _now() - 1);

        vm.expectEmit(address(ledger));
        emit ILiabilityLedger.EpochAnchored(
            0,
            EpochChain.digest(header),
            bytes32(0),
            header.root,
            header.total,
            header.leafCount,
            header.asOf
        );
        _anchor(header);
    }

    function test_RevertWhen_CallerIsNotTheAnchorer() public {
        EpochHeader memory header = _header(0, bytes32(0), _tree(0, 1, 2), 2, _now() - 1);

        vm.expectRevert(
            abi.encodeWithSelector(
                IAccessControl.AccessControlUnauthorizedAccount.selector, stranger, ledger.ANCHORER_ROLE()
            )
        );
        vm.prank(stranger);
        ledger.anchor(header);
    }

    function test_RevertWhen_TheSameEpochIsAnchoredTwice() public {
        EpochHeader memory header = _header(0, bytes32(0), _tree(0, 1, 2), 2, _now() - 1);
        bytes32 digest = _anchor(header);

        vm.expectRevert(abi.encodeWithSelector(ILiabilityLedger.EpochAlreadyAnchored.selector, 0, digest));
        _anchor(header);
    }

    function test_RevertWhen_AnAnchoredEpochIsRewritten() public {
        EpochHeader memory header = _header(0, bytes32(0), _tree(0, 1, 2), 2, _now() - 1);
        bytes32 digest = _anchor(header);
        EpochHeader memory rewritten = _header(0, bytes32(0), _tree(0, 9, 2), 2, _now() - 1);

        vm.expectRevert(
            abi.encodeWithSelector(ILiabilityLedger.EpochConflict.selector, 0, digest, EpochChain.digest(rewritten))
        );
        _anchor(rewritten);
    }

    function test_RevertWhen_AnEpochIsSkipped() public {
        EpochHeader memory header = _header(1, bytes32(0), _tree(1, 1, 2), 2, _now() - 1);

        vm.expectRevert(abi.encodeWithSelector(ILiabilityLedger.EpochOutOfSequence.selector, 0, 1));
        _anchor(header);
    }

    function test_RevertWhen_ThePreviousDigestIsNotTheHead() public {
        bytes32 first = _anchor(_header(0, bytes32(0), _tree(0, 1, 2), 2, _now() - 10));
        bytes32 forged = keccak256("forged");
        EpochHeader memory header = _header(1, forged, _tree(1, 2, 2), 2, _now() - 1);

        vm.expectRevert(abi.encodeWithSelector(ILiabilityLedger.PreviousDigestMismatch.selector, first, forged));
        _anchor(header);
    }

    function test_RevertWhen_TheSnapshotIsNotAfterThePreviousOne() public {
        uint64 asOf = _now() - 10;
        bytes32 first = _anchor(_header(0, bytes32(0), _tree(0, 1, 2), 2, asOf));
        EpochHeader memory header = _header(1, first, _tree(1, 2, 2), 2, asOf);

        vm.expectRevert(abi.encodeWithSelector(ILiabilityLedger.SnapshotNotAfterPrevious.selector, asOf, asOf));
        _anchor(header);
    }

    function test_RevertWhen_TheFirstSnapshotIsAtTimeZero() public {
        EpochHeader memory header = _header(0, bytes32(0), _tree(0, 1, 2), 2, 0);

        vm.expectRevert(abi.encodeWithSelector(ILiabilityLedger.SnapshotNotAfterPrevious.selector, 0, 0));
        _anchor(header);
    }

    function test_RevertWhen_TheSnapshotIsInTheFuture() public {
        EpochHeader memory header = _header(0, bytes32(0), _tree(0, 1, 2), 2, _now() + 1);

        vm.expectRevert(abi.encodeWithSelector(ILiabilityLedger.SnapshotInFuture.selector, _now() + 1, _now()));
        _anchor(header);
    }

    function test_RevertWhen_TheTreeIsDeeperThanTheScheme() public {
        EpochHeader memory header = EpochHeader({
            number: 0,
            previousDigest: bytes32(0),
            root: keccak256("root"),
            total: 1,
            leafCount: SumTree.MAX_LEAVES + 1,
            asOf: _now() - 1
        });

        vm.expectRevert(
            abi.encodeWithSelector(ILiabilityLedger.TooManyLeaves.selector, SumTree.MAX_LEAVES + 1, SumTree.MAX_LEAVES)
        );
        _anchor(header);
    }

    function test_AcceptsAnEmptyEpochWithTheCanonicalEmptyRoot() public {
        EpochHeader memory header = _header(0, bytes32(0), _tree(0, 1, 0), 0, _now() - 1);

        _anchor(header);

        assertEq(ledger.epochAt(0).root, SumTree.empty().hash);
        assertEq(ledger.epochAt(0).total, 0);
    }

    function test_RevertWhen_AnEmptyEpochCarriesARoot() public {
        EpochHeader memory header = EpochHeader({
            number: 0,
            previousDigest: bytes32(0),
            root: keccak256("root"),
            total: 0,
            leafCount: 0,
            asOf: _now() - 1
        });

        vm.expectRevert(abi.encodeWithSelector(ILiabilityLedger.EmptyEpochMismatch.selector, header.root, 0));
        _anchor(header);
    }

    function test_RevertWhen_AnEmptyEpochCarriesATotal() public {
        bytes32 emptyRoot = SumTree.empty().hash;
        EpochHeader memory header = EpochHeader({
            number: 0,
            previousDigest: bytes32(0),
            root: emptyRoot,
            total: 7,
            leafCount: 0,
            asOf: _now() - 1
        });

        vm.expectRevert(abi.encodeWithSelector(ILiabilityLedger.EmptyEpochMismatch.selector, emptyRoot, 7));
        _anchor(header);
    }

    function test_VerifiesEveryAccountInTheEpoch() public {
        uint256 count = 7;
        SumNode[][] memory levels = _tree(0, 3, count);
        _anchor(_header(0, bytes32(0), levels, count, _now() - 1));
        bytes32[] memory commitments = _commitments(3, count);
        uint64[] memory amounts = _amounts(3, count);

        for (uint256 i = 0; i < count; i++) {
            assertTrue(ledger.verifyInclusion(0, commitments[i], amounts[i], uint64(i), TreeBuilder.proof(levels, i)));
        }
    }

    function test_VerifiesTheOnlyAccountOfASingleLeafEpoch() public {
        SumNode[][] memory levels = _tree(0, 4, 1);
        _anchor(_header(0, bytes32(0), levels, 1, _now() - 1));

        assertTrue(ledger.verifyInclusion(0, _commitments(4, 1)[0], _amounts(4, 1)[0], 0, new SumNode[](0)));
    }

    function test_RejectsAnInflatedAmount() public {
        (SumNode[][] memory levels, bytes32 commitment, uint64 amount) = _anchoredTree(5, 2);

        assertFalse(ledger.verifyInclusion(0, commitment, amount + 1, 2, TreeBuilder.proof(levels, 2)));
    }

    function test_RejectsAForeignCommitment() public {
        (SumNode[][] memory levels,, uint64 amount) = _anchoredTree(5, 2);

        assertFalse(ledger.verifyInclusion(0, keccak256("someone else"), amount, 2, TreeBuilder.proof(levels, 2)));
    }

    function test_RejectsASiblingSumMovedOntoTheAmount() public {
        (SumNode[][] memory levels, bytes32 commitment, uint64 amount) = _anchoredTree(5, 2);
        SumNode[] memory siblings = TreeBuilder.proof(levels, 2);
        assertGt(siblings[0].sum, 0);
        uint64 moved = siblings[0].sum;
        siblings[0].sum = 0;

        assertFalse(ledger.verifyInclusion(0, commitment, amount + moved, 2, siblings));
    }

    function test_RejectsAProofWhoseSumsOverflowRatherThanReverting() public {
        (SumNode[][] memory levels, bytes32 commitment, uint64 amount) = _anchoredTree(5, 2);
        SumNode[] memory siblings = TreeBuilder.proof(levels, 2);
        for (uint256 i = 0; i < siblings.length; i++) {
            siblings[i].sum = type(uint64).max;
        }

        assertFalse(ledger.verifyInclusion(0, commitment, amount, 2, siblings));
    }

    function test_RejectsAForgedSiblingHash() public {
        (SumNode[][] memory levels, bytes32 commitment, uint64 amount) = _anchoredTree(5, 2);
        SumNode[] memory siblings = TreeBuilder.proof(levels, 2);
        siblings[1].hash = keccak256("forged");

        assertFalse(ledger.verifyInclusion(0, commitment, amount, 2, siblings));
    }

    function test_RejectsTheWrongIndex() public {
        (SumNode[][] memory levels, bytes32 commitment, uint64 amount) = _anchoredTree(5, 2);

        assertFalse(ledger.verifyInclusion(0, commitment, amount, 3, TreeBuilder.proof(levels, 2)));
    }

    function test_RejectsAnIndexInThePadding() public {
        (SumNode[][] memory levels,,) = _anchoredTree(5, 2);

        assertFalse(ledger.verifyInclusion(0, bytes32(0), 0, 6, TreeBuilder.proof(levels, 6)));
    }

    function test_RejectsAProofOfTheWrongLength() public {
        (SumNode[][] memory levels, bytes32 commitment, uint64 amount) = _anchoredTree(5, 2);
        SumNode[] memory full = TreeBuilder.proof(levels, 2);
        SumNode[] memory short = new SumNode[](full.length - 1);
        for (uint256 i = 0; i < short.length; i++) {
            short[i] = full[i];
        }

        assertFalse(ledger.verifyInclusion(0, commitment, amount, 2, short));
    }

    function test_RevertWhen_VerifyingAgainstAnEpochThatWasNeverAnchored() public {
        vm.expectRevert(abi.encodeWithSelector(ILiabilityLedger.UnknownEpoch.selector, 0));
        ledger.verifyInclusion(0, bytes32(0), 0, 0, new SumNode[](0));
    }

    function test_RevertWhen_ReadingAnEpochThatWasNeverAnchored() public {
        vm.expectRevert(abi.encodeWithSelector(ILiabilityLedger.UnknownEpoch.selector, 0));
        ledger.epochAt(0);
    }

    function test_SealingIsOneWayAndStopsAnchoring() public {
        bytes32 head = _anchor(_header(0, bytes32(0), _tree(0, 1, 2), 2, _now() - 10));

        vm.expectEmit(address(ledger));
        emit ILiabilityLedger.LedgerSealed(head, 1, _now() - 10);
        vm.prank(admin);
        ledger.seal(head);

        assertTrue(ledger.isSealed());
        vm.expectRevert(ILiabilityLedger.LedgerIsSealed.selector);
        _anchor(_header(1, head, _tree(1, 2, 2), 2, _now() - 1));
        vm.expectRevert(ILiabilityLedger.LedgerIsSealed.selector);
        vm.prank(admin);
        ledger.seal(head);
    }

    function test_RevertWhen_SealingAHeadThatHasMoved() public {
        bytes32 head = _anchor(_header(0, bytes32(0), _tree(0, 1, 2), 2, _now() - 10));
        bytes32 stale = keccak256("stale");

        vm.expectRevert(abi.encodeWithSelector(ILiabilityLedger.HeadMismatch.selector, head, stale));
        vm.prank(admin);
        ledger.seal(stale);
    }

    function test_RevertWhen_TheAnchorerTriesToSeal() public {
        vm.expectRevert(
            abi.encodeWithSelector(IAccessControl.AccessControlUnauthorizedAccount.selector, anchorer, bytes32(0))
        );
        vm.prank(anchorer);
        ledger.seal(bytes32(0));
    }

    function test_ASuccessorContinuesTheSealedChain() public {
        bytes32 first = _anchor(_header(0, bytes32(0), _tree(0, 1, 2), 2, _now() - 20));
        bytes32 head = _anchor(_header(1, first, _tree(1, 2, 3), 3, _now() - 10));
        vm.prank(admin);
        ledger.seal(head);

        LiabilityLedger successor =
            new LiabilityLedger(admin, ADMIN_DELAY, anchorer, ledger.nextEpoch(), ledger.head(), ledger.lastAsOf());
        EpochHeader memory next = _header(2, head, _tree(2, 3, 4), 4, _now() - 1);
        vm.prank(anchorer);
        bytes32 digest = successor.anchor(next);

        assertEq(digest, EpochChain.digest(next));
        assertEq(successor.genesisEpoch(), 2);
        assertEq(successor.genesisDigest(), head);
        vm.expectRevert(abi.encodeWithSelector(ILiabilityLedger.UnknownEpoch.selector, 1));
        successor.epochAt(1);
    }

    function test_RevertWhen_ASuccessorIsOfferedAnEpochBeforeItsGenesis() public {
        LiabilityLedger successor =
            new LiabilityLedger(admin, ADMIN_DELAY, anchorer, 5, keccak256("head"), _now() - 100);
        EpochHeader memory old = _header(4, bytes32(0), _tree(4, 1, 2), 2, _now() - 1);

        vm.expectRevert(abi.encodeWithSelector(ILiabilityLedger.EpochBeforeGenesis.selector, 5, 4));
        vm.prank(anchorer);
        successor.anchor(old);
    }

    function test_RevertWhen_AFreshChainClaimsAGenesisDigest() public {
        vm.expectRevert(abi.encodeWithSelector(ILiabilityLedger.InvalidGenesis.selector, 0, keccak256("head")));
        new LiabilityLedger(admin, ADMIN_DELAY, anchorer, 0, keccak256("head"), 0);
    }

    function test_RevertWhen_ASuccessorHasNoGenesisDigest() public {
        vm.expectRevert(abi.encodeWithSelector(ILiabilityLedger.InvalidGenesis.selector, 3, bytes32(0)));
        new LiabilityLedger(admin, ADMIN_DELAY, anchorer, 3, bytes32(0), 0);
    }

    function test_RevertWhen_ThereIsNoAnchorer() public {
        vm.expectRevert(ILiabilityLedger.ZeroAnchorer.selector);
        new LiabilityLedger(admin, ADMIN_DELAY, address(0), 0, bytes32(0), 0);
    }

    function test_TheAdminRotatesTheAnchorer() public {
        address replacement = makeAddr("replacement");
        bytes32 role = ledger.ANCHORER_ROLE();
        vm.startPrank(admin);
        ledger.revokeRole(role, anchorer);
        ledger.grantRole(role, replacement);
        vm.stopPrank();
        EpochHeader memory header = _header(0, bytes32(0), _tree(0, 1, 2), 2, _now() - 1);

        vm.expectRevert(
            abi.encodeWithSelector(IAccessControl.AccessControlUnauthorizedAccount.selector, anchorer, role)
        );
        vm.prank(anchorer);
        ledger.anchor(header);
        vm.prank(replacement);
        ledger.anchor(header);
    }

    function test_TheAdminTransferWaitsForTheDelay() public {
        address successorAdmin = makeAddr("successorAdmin");
        vm.prank(admin);
        ledger.beginDefaultAdminTransfer(successorAdmin);

        vm.expectRevert();
        vm.prank(successorAdmin);
        ledger.acceptDefaultAdminTransfer();

        vm.warp(vm.getBlockTimestamp() + ADMIN_DELAY + 1);
        vm.prank(successorAdmin);
        ledger.acceptDefaultAdminTransfer();
        assertEq(ledger.defaultAdmin(), successorAdmin);
    }

    function test_AnnouncesTheLedgerInterface() public view {
        assertTrue(ledger.supportsInterface(type(ILiabilityLedger).interfaceId));
        assertTrue(ledger.supportsInterface(type(IAccessControl).interfaceId));
        assertTrue(ledger.supportsInterface(type(IERC165).interfaceId));
        assertFalse(ledger.supportsInterface(0xffffffff));
        assertEq(ledger.scheme(), 1);
    }

    function _anchoredTree(uint256 seed, uint256 index) private returns (SumNode[][] memory, bytes32, uint64) {
        uint256 count = 5;
        SumNode[][] memory levels = _tree(0, seed, count);
        _anchor(_header(0, bytes32(0), levels, count, _now() - 1));
        return (levels, _commitments(seed, count)[index], _amounts(seed, count)[index]);
    }
}
