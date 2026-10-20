// SPDX-License-Identifier: MIT
pragma solidity 0.8.37;

import {IAccessControl} from "@openzeppelin-contracts/access/IAccessControl.sol";
import {IERC165} from "@openzeppelin-contracts/utils/introspection/IERC165.sol";

import {Epoch} from "../../src/Epoch.sol";
import {EpochHeader} from "../../src/EpochHeader.sol";
import {ILiabilityLedger} from "../../src/ILiabilityLedger.sol";
import {SumNode} from "../../src/SumNode.sol";
import {ConfidentialEpochChain} from "../../src/confidential/ConfidentialEpochChain.sol";
import {ConfidentialTree} from "../../src/confidential/ConfidentialTree.sol";
import {IConfidentialLedger} from "../../src/confidential/IConfidentialLedger.sol";
import {LiabilityLedgerV2} from "../../src/confidential/LiabilityLedgerV2.sol";
import {PlonkVerifier} from "../../src/confidential/PlonkVerifier.sol";
import {Poseidon2} from "../../src/confidential/Poseidon2.sol";
import {ConfidentialLedgerFixture} from "../support/ConfidentialLedgerFixture.sol";
import {ConfidentialStatement} from "../support/ConfidentialStatement.sol";

contract LiabilityLedgerV2Test is ConfidentialLedgerFixture {
    function test_SpeaksSchemeTwoAtItsDepth() public view {
        assertEq(ledger.scheme(), 2);
        assertEq(ledger.depth(), DEPTH);
        assertEq(ledger.capacity(), 8);
        assertEq(ledger.verifier(), address(verifier));
    }

    function test_AnnouncesBothLedgerInterfaces() public view {
        assertTrue(ledger.supportsInterface(type(ILiabilityLedger).interfaceId));
        assertTrue(ledger.supportsInterface(type(IConfidentialLedger).interfaceId));
        assertTrue(ledger.supportsInterface(type(IERC165).interfaceId));
        assertFalse(ledger.supportsInterface(0xdeadbeef));
    }

    function test_RevertWhen_TheGenesisIsHalfSet() public {
        vm.expectRevert(abi.encodeWithSelector(ILiabilityLedger.InvalidGenesis.selector, 4, bytes32(0)));
        new LiabilityLedgerV2(admin, ADMIN_DELAY, anchorer, 4, bytes32(0), 0, PlonkVerifier(address(verifier)), DEPTH);
    }

    function test_RevertWhen_TheAnchorerIsZero() public {
        vm.expectRevert(ILiabilityLedger.ZeroAnchorer.selector);
        new LiabilityLedgerV2(admin, ADMIN_DELAY, address(0), 0, bytes32(0), 0, PlonkVerifier(address(verifier)), DEPTH);
    }

    function test_RevertWhen_TheVerifierHasNoCode() public {
        vm.expectRevert(abi.encodeWithSelector(IConfidentialLedger.NotAVerifier.selector, stranger));
        new LiabilityLedgerV2(admin, ADMIN_DELAY, anchorer, 0, bytes32(0), 0, PlonkVerifier(stranger), DEPTH);
    }

    function test_RevertWhen_TheDepthExceedsTheTreesMaximum() public {
        vm.expectRevert(abi.encodeWithSelector(IConfidentialLedger.DepthOutOfRange.selector, 17, 16));
        new LiabilityLedgerV2(admin, ADMIN_DELAY, anchorer, 0, bytes32(0), 0, PlonkVerifier(address(verifier)), 17);
    }

    function test_AProvenStatementIsRecordedUnderItsBindingRoot() public {
        ConfidentialStatement memory statement = _statement(0, 1, 5);
        bytes32 root = ConfidentialTree.bindingRoot(statement.treeRoot, statement.totalCommitment);
        verifier.accept(statement.proof, _inputs(statement));

        vm.expectEmit(address(ledger));
        emit IConfidentialLedger.EpochProven(0, root, statement.treeRoot, statement.totalCommitment, 5);
        vm.prank(stranger);
        bytes32 returned = ledger.proveEpoch(0, statement.treeRoot, statement.totalCommitment, 5, statement.proof);

        assertEq(returned, root);
        assertTrue(ledger.isProven(0, root, 5));
        assertFalse(ledger.isProven(0, root, 4));
        assertFalse(ledger.isProven(1, root, 5));
    }

    function test_TheVerifierIsGivenTheStatementInTheCircuitsOrder() public {
        ConfidentialStatement memory statement = _statement(0, 1, 5);
        uint256[] memory swapped = _inputs(statement);
        (swapped[1], swapped[2]) = (swapped[2], swapped[1]);
        verifier.accept(statement.proof, swapped);

        vm.expectRevert(abi.encodeWithSelector(IConfidentialLedger.ProofRejected.selector, 0));
        ledger.proveEpoch(0, statement.treeRoot, statement.totalCommitment, 5, statement.proof);
    }

    function test_RevertWhen_TheVerifierRefuses() public {
        ConfidentialStatement memory statement = _statement(0, 1, 5);

        vm.expectRevert(abi.encodeWithSelector(IConfidentialLedger.ProofRejected.selector, 0));
        ledger.proveEpoch(0, statement.treeRoot, statement.totalCommitment, 5, statement.proof);
    }

    function test_RevertWhen_TheVerifierReverts() public {
        ConfidentialStatement memory statement = _statement(0, 1, 5);
        verifier.accept(statement.proof, _inputs(statement));
        verifier.revertEverything();

        vm.expectRevert(abi.encodeWithSelector(IConfidentialLedger.ProofRejected.selector, 0));
        ledger.proveEpoch(0, statement.treeRoot, statement.totalCommitment, 5, statement.proof);
    }

    function test_RevertWhen_TheSameStatementIsProvenTwice() public {
        ConfidentialStatement memory statement = _statement(0, 1, 5);
        bytes32 root = _prove(statement);

        vm.expectRevert(abi.encodeWithSelector(IConfidentialLedger.StatementAlreadyProven.selector, 0, root, 5));
        ledger.proveEpoch(0, statement.treeRoot, statement.totalCommitment, 5, statement.proof);
    }

    function test_AStrangersStatementForTheSameEpochNeverDisplacesTheNotarys() public {
        ConfidentialStatement memory ours = _statement(0, 1, 5);
        ConfidentialStatement memory theirs = _statement(0, 99, 8);
        _prove(ours);
        vm.prank(stranger);
        _prove(theirs);

        bytes32 digest = _anchor(_header(ours, bytes32(0), _now() - 60));

        assertEq(ledger.head(), digest);
        assertEq(ledger.treeRootOf(0), ours.treeRoot);
        assertEq(ledger.totalCommitmentOf(0), ours.totalCommitment);
    }

    function test_AnEpochMayBeProvenAheadOfItsTurn() public {
        _prove(_statement(5, 1, 3));

        assertTrue(
            ledger.isProven(
                5, ConfidentialTree.bindingRoot(_statement(5, 1, 3).treeRoot, _statement(5, 1, 3).totalCommitment), 3
            )
        );
    }

    function test_RevertWhen_TheEpochIsAlreadyClosed() public {
        _proveAndAnchor(_statement(0, 1, 5), bytes32(0));
        ConfidentialStatement memory late = _statement(0, 2, 5);

        vm.expectRevert(abi.encodeWithSelector(IConfidentialLedger.EpochClosed.selector, 0, 1));
        ledger.proveEpoch(0, late.treeRoot, late.totalCommitment, 5, late.proof);
    }

    function test_RevertWhen_TheStatementHasMoreAccountsThanTheCapacity() public {
        ConfidentialStatement memory statement = _statement(0, 1, 5);

        vm.expectRevert(abi.encodeWithSelector(ILiabilityLedger.TooManyLeaves.selector, 9, 8));
        ledger.proveEpoch(0, statement.treeRoot, statement.totalCommitment, 9, statement.proof);
    }

    function test_RevertWhen_ATreeRootIsNotAFieldElement() public {
        ConfidentialStatement memory statement = _statement(0, 1, 5);
        bytes32 aliased = bytes32(uint256(statement.treeRoot) + Poseidon2.PRIME);

        vm.expectRevert(abi.encodeWithSelector(IConfidentialLedger.NotAFieldElement.selector, aliased));
        ledger.proveEpoch(0, aliased, statement.totalCommitment, 5, statement.proof);
    }

    function test_RevertWhen_ATotalCommitmentIsNotAFieldElement() public {
        ConfidentialStatement memory statement = _statement(0, 1, 5);
        bytes32 aliased = bytes32(uint256(statement.totalCommitment) + Poseidon2.PRIME);

        vm.expectRevert(abi.encodeWithSelector(IConfidentialLedger.NotAFieldElement.selector, aliased));
        ledger.proveEpoch(0, statement.treeRoot, aliased, 5, statement.proof);
    }

    function test_AnchorsAProvenEpochWithItsTotalHidden() public {
        ConfidentialStatement memory statement = _statement(0, 1, 5);
        _prove(statement);
        EpochHeader memory header = _header(statement, bytes32(0), _now() - 60);

        vm.expectEmit(address(ledger));
        emit ILiabilityLedger.EpochAnchored(
            0,
            ConfidentialEpochChain.digest(header),
            bytes32(0),
            header.root,
            0,
            5,
            header.asOf
        );
        bytes32 digest = _anchor(header);

        assertEq(digest, ConfidentialEpochChain.digest(header));
        assertEq(ledger.head(), digest);
        assertEq(ledger.nextEpoch(), 1);
        assertEq(ledger.lastAsOf(), header.asOf);
        Epoch memory stored = ledger.epochAt(0);
        assertEq(stored.root, header.root);
        assertEq(stored.digest, digest);
        assertEq(stored.total, 0);
        assertEq(stored.leafCount, 5);
        assertEq(stored.asOf, header.asOf);
        assertEq(stored.anchoredAt, _now());
        assertEq(ledger.treeRootOf(0), statement.treeRoot);
        assertEq(ledger.totalCommitmentOf(0), statement.totalCommitment);
    }

    function test_ChainsEachEpochToThePreviousDigest() public {
        bytes32 first = _proveAndAnchor(_statement(0, 1, 3), bytes32(0));
        vm.warp(vm.getBlockTimestamp() + 1 hours);

        bytes32 second = _proveAndAnchor(_statement(1, 2, 8), first);

        assertEq(ledger.head(), second);
        assertEq(ledger.nextEpoch(), 2);
    }

    function test_AnEmptyEpochIsProvenLikeAnyOther() public {
        _proveAndAnchor(_statement(0, 1, 0), bytes32(0));

        assertEq(ledger.epochAt(0).leafCount, 0);
    }

    function test_RevertWhen_TheEpochWasNeverProven() public {
        EpochHeader memory header = _header(_statement(0, 1, 5), bytes32(0), _now() - 60);

        vm.expectRevert(abi.encodeWithSelector(IConfidentialLedger.EpochNotProven.selector, 0, header.root, 5));
        _anchor(header);
    }

    function test_RevertWhen_TheHeaderDiffersFromWhatWasProven() public {
        ConfidentialStatement memory statement = _statement(0, 1, 5);
        _prove(statement);
        EpochHeader memory proven = _header(statement, bytes32(0), _now() - 60);

        EpochHeader[] memory others = new EpochHeader[](3);
        others[0] = proven;
        others[0].leafCount = 4;
        others[1] = _header(statement, bytes32(0), _now() - 60);
        others[1].root = ConfidentialTree.bindingRoot(statement.treeRoot, bytes32(uint256(1)));
        others[2] = _header(_statement(0, 2, 5), bytes32(0), _now() - 60);

        for (uint256 i = 0; i < others.length; i++) {
            vm.expectRevert(
                abi.encodeWithSelector(
                    IConfidentialLedger.EpochNotProven.selector, 0, others[i].root, others[i].leafCount
                )
            );
            this.anchorAs(anchorer, others[i]);
        }
    }

    function test_RevertWhen_AProofOfAnotherEpochIsOffered() public {
        ConfidentialStatement memory statement = _statement(1, 1, 5);
        _prove(statement);
        EpochHeader memory header = _header(statement, bytes32(0), _now() - 60);
        header.number = 0;

        vm.expectRevert(abi.encodeWithSelector(IConfidentialLedger.EpochNotProven.selector, 0, header.root, 5));
        _anchor(header);
    }

    function test_RevertWhen_TheHeaderCarriesATotal() public {
        ConfidentialStatement memory statement = _statement(0, 1, 5);
        _prove(statement);
        EpochHeader memory header = _header(statement, bytes32(0), _now() - 60);
        header.total = 7;

        vm.expectRevert(abi.encodeWithSelector(IConfidentialLedger.TotalNotHidden.selector, 7));
        _anchor(header);
    }

    function test_RevertWhen_CallerIsNotTheAnchorer() public {
        ConfidentialStatement memory statement = _statement(0, 1, 5);
        _prove(statement);
        EpochHeader memory header = _header(statement, bytes32(0), _now() - 60);

        vm.expectRevert(
            abi.encodeWithSelector(
                IAccessControl.AccessControlUnauthorizedAccount.selector, stranger, ledger.ANCHORER_ROLE()
            )
        );
        vm.prank(stranger);
        ledger.anchor(header);
    }

    function test_RevertWhen_TheSameEpochIsAnchoredTwice() public {
        ConfidentialStatement memory statement = _statement(0, 1, 5);
        _prove(statement);
        EpochHeader memory header = _header(statement, bytes32(0), _now() - 60);
        bytes32 digest = _anchor(header);

        vm.expectRevert(abi.encodeWithSelector(ILiabilityLedger.EpochAlreadyAnchored.selector, 0, digest));
        _anchor(header);
    }

    function test_RevertWhen_AnAnchoredEpochIsContradicted() public {
        bytes32 digest = _proveAndAnchor(_statement(0, 1, 5), bytes32(0));
        EpochHeader memory other = _header(_statement(0, 2, 5), bytes32(0), _now() - 30);

        vm.expectRevert(
            abi.encodeWithSelector(
                ILiabilityLedger.EpochConflict.selector, 0, digest, ConfidentialEpochChain.digest(other)
            )
        );
        _anchor(other);
    }

    function test_RevertWhen_AnEpochSkipsAhead() public {
        ConfidentialStatement memory statement = _statement(1, 1, 5);
        _prove(statement);

        vm.expectRevert(abi.encodeWithSelector(ILiabilityLedger.EpochOutOfSequence.selector, 0, 1));
        _anchor(_header(statement, bytes32(0), _now() - 60));
    }

    function test_RevertWhen_ThePreviousDigestIsNotTheHead() public {
        _proveAndAnchor(_statement(0, 1, 5), bytes32(0));
        ConfidentialStatement memory next = _statement(1, 2, 5);
        _prove(next);

        vm.expectRevert(
            abi.encodeWithSelector(ILiabilityLedger.PreviousDigestMismatch.selector, ledger.head(), bytes32(uint256(9)))
        );
        _anchor(_header(next, bytes32(uint256(9)), _now() - 30));
    }

    function test_RevertWhen_TheSnapshotIsNotAfterThePrevious() public {
        _proveAndAnchor(_statement(0, 1, 5), bytes32(0));
        uint64 previous = ledger.lastAsOf();
        ConfidentialStatement memory next = _statement(1, 2, 5);
        _prove(next);
        bytes32 head = ledger.head();

        vm.expectRevert(abi.encodeWithSelector(ILiabilityLedger.SnapshotNotAfterPrevious.selector, previous, previous));
        _anchor(_header(next, head, previous));
    }

    function test_RevertWhen_TheSnapshotIsInTheFuture() public {
        ConfidentialStatement memory statement = _statement(0, 1, 5);
        _prove(statement);
        uint64 future = _now() + 1;

        vm.expectRevert(abi.encodeWithSelector(ILiabilityLedger.SnapshotInFuture.selector, future, _now()));
        _anchor(_header(statement, bytes32(0), future));
    }

    function test_EveryAccountOfAnAnchoredEpochVerifies() public {
        for (uint256 count = 1; count <= 8; count++) {
            uint64 number = uint64(count - 1);
            bytes32 previous = ledger.head();
            vm.warp(vm.getBlockTimestamp() + 1 hours);
            _proveAndAnchor(_statement(number, count, count), previous);

            bytes32[] memory commitments = _commitments(count, count);
            uint64[] memory amounts = _amounts(count, count);
            for (uint256 i = 0; i < count; i++) {
                SumNode[] memory siblings = _siblings(number, count, count, i);
                assertTrue(ledger.verifyInclusion(number, commitments[i], amounts[i], uint64(i), siblings));
            }
        }
    }

    function test_AnInclusionIsFalseWhenAnythingDiffers() public {
        _proveAndAnchor(_statement(0, 4, 6), bytes32(0));
        bytes32[] memory commitments = _commitments(4, 6);
        uint64[] memory amounts = _amounts(4, 6);
        SumNode[] memory siblings = _siblings(0, 4, 6, 2);

        assertTrue(ledger.verifyInclusion(0, commitments[2], amounts[2], 2, siblings));
        assertFalse(ledger.verifyInclusion(0, commitments[2], amounts[2] + 1, 2, siblings));
        assertFalse(ledger.verifyInclusion(0, commitments[3], amounts[2], 2, siblings));
        assertFalse(ledger.verifyInclusion(0, commitments[2], amounts[2], 3, siblings));
        assertFalse(ledger.verifyInclusion(0, commitments[2], amounts[2], 6, _siblings(0, 4, 6, 6)));

        SumNode[] memory short = new SumNode[](2);
        short[0] = siblings[0];
        short[1] = siblings[1];
        assertFalse(ledger.verifyInclusion(0, commitments[2], amounts[2], 2, short));

        SumNode[] memory summed = _siblings(0, 4, 6, 2);
        summed[0].sum = 1;
        assertFalse(ledger.verifyInclusion(0, commitments[2], amounts[2], 2, summed));

        SumNode[] memory aliased = _siblings(0, 4, 6, 2);
        aliased[1].hash = bytes32(uint256(aliased[1].hash) + Poseidon2.PRIME);
        assertFalse(ledger.verifyInclusion(0, commitments[2], amounts[2], 2, aliased));
    }

    function test_RevertWhen_AnInclusionNamesAnUnknownEpoch() public {
        SumNode[] memory siblings = _siblings(0, 4, 6, 2);

        vm.expectRevert(abi.encodeWithSelector(ILiabilityLedger.UnknownEpoch.selector, 0));
        ledger.verifyInclusion(0, bytes32(0), 0, 0, siblings);
    }

    function test_RevertWhen_TheStatementOfAnUnknownEpochIsAsked() public {
        vm.expectRevert(abi.encodeWithSelector(ILiabilityLedger.UnknownEpoch.selector, 0));
        ledger.treeRootOf(0);

        vm.expectRevert(abi.encodeWithSelector(ILiabilityLedger.UnknownEpoch.selector, 0));
        ledger.totalCommitmentOf(0);
    }

    function test_TheAdminSealsAtTheHead() public {
        bytes32 head = _proveAndAnchor(_statement(0, 1, 5), bytes32(0));

        vm.expectEmit(address(ledger));
        emit ILiabilityLedger.LedgerSealed(head, 1, ledger.lastAsOf());
        vm.prank(admin);
        ledger.seal(head);

        assertTrue(ledger.isSealed());
    }

    function test_RevertWhen_TheSealNamesAnotherHead() public {
        bytes32 head = _proveAndAnchor(_statement(0, 1, 5), bytes32(0));

        vm.expectRevert(abi.encodeWithSelector(ILiabilityLedger.HeadMismatch.selector, head, bytes32(uint256(1))));
        vm.prank(admin);
        ledger.seal(bytes32(uint256(1)));
    }

    function test_RevertWhen_AStrangerSeals() public {
        vm.expectRevert(
            abi.encodeWithSelector(IAccessControl.AccessControlUnauthorizedAccount.selector, stranger, bytes32(0))
        );
        vm.prank(stranger);
        ledger.seal(bytes32(0));
    }

    function test_ASealedLedgerNeitherProvesNorAnchors() public {
        ConfidentialStatement memory statement = _statement(0, 1, 5);
        _prove(statement);
        vm.prank(admin);
        ledger.seal(bytes32(0));

        vm.expectRevert(ILiabilityLedger.LedgerIsSealed.selector);
        _anchor(_header(statement, bytes32(0), _now() - 60));

        ConfidentialStatement memory other = _statement(0, 2, 5);
        vm.expectRevert(ILiabilityLedger.LedgerIsSealed.selector);
        ledger.proveEpoch(0, other.treeRoot, other.totalCommitment, 5, other.proof);

        vm.expectRevert(ILiabilityLedger.LedgerIsSealed.selector);
        vm.prank(admin);
        ledger.seal(bytes32(0));
    }

    function test_ASuccessorStartsAtThePredecessorsHead() public {
        bytes32 genesis = keccak256("sealed head");
        LiabilityLedgerV2 successor = new LiabilityLedgerV2(
            admin, ADMIN_DELAY, anchorer, 40, genesis, _now() - 600, PlonkVerifier(address(verifier)), DEPTH
        );
        ConfidentialStatement memory statement = _statement(40, 1, 5);
        verifier.accept(statement.proof, _inputs(statement));
        successor.proveEpoch(40, statement.treeRoot, statement.totalCommitment, 5, statement.proof);

        vm.expectRevert(abi.encodeWithSelector(ILiabilityLedger.UnknownEpoch.selector, 39));
        successor.epochAt(39);

        EpochHeader memory early = _header(_statement(39, 1, 5), genesis, _now() - 60);
        vm.expectRevert(abi.encodeWithSelector(ILiabilityLedger.EpochBeforeGenesis.selector, 40, 39));
        vm.prank(anchorer);
        successor.anchor(early);

        vm.prank(anchorer);
        bytes32 digest = successor.anchor(_header(statement, genesis, _now() - 60));
        assertEq(successor.head(), digest);
        assertEq(successor.genesisEpoch(), 40);
        assertEq(successor.genesisDigest(), genesis);
    }

    function anchorAs(address who, EpochHeader memory header) external returns (bytes32) {
        vm.prank(who);
        return ledger.anchor(header);
    }
}
