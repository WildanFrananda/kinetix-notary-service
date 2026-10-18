// SPDX-License-Identifier: MIT
pragma solidity 0.8.37;

import {IAccessControl} from "@openzeppelin-contracts/access/IAccessControl.sol";

import {Batch} from "../../src/evidence/Batch.sol";
import {BatchChain} from "../../src/evidence/BatchChain.sol";
import {BatchHeader} from "../../src/evidence/BatchHeader.sol";
import {EvidenceRegistry} from "../../src/evidence/EvidenceRegistry.sol";
import {EvidenceTree} from "../../src/evidence/EvidenceTree.sol";
import {IEvidenceRegistry} from "../../src/evidence/IEvidenceRegistry.sol";
import {StreamGenesis} from "../../src/evidence/StreamGenesis.sol";
import {StreamState} from "../../src/evidence/StreamState.sol";
import {EvidenceFixture} from "../support/EvidenceFixture.sol";
import {EvidenceTreeBuilder} from "../support/EvidenceTreeBuilder.sol";

contract EvidenceRegistryTest is EvidenceFixture {
    function test_StartsWithEveryStreamItWasGivenAndNothingAnchored() public view {
        uint8[] memory ids = registry.streams();

        assertEq(ids.length, 3);
        assertEq(ids[0], LIFECYCLE);
        assertEq(ids[1], DELIVERY);
        assertEq(ids[2], REVIEW);
        assertEq(registry.scheme(), 1);
        assertFalse(registry.isSealed());
        assertTrue(registry.supportsInterface(type(IEvidenceRegistry).interfaceId));
        assertTrue(registry.hasRole(registry.ANCHORER_ROLE(), anchorer));
        assertEq(registry.defaultAdmin(), admin);
        assertEq(registry.defaultAdminDelay(), ADMIN_DELAY);

        StreamState memory fresh = registry.stateOf(DELIVERY);
        assertTrue(fresh.known);
        assertEq(fresh.head, bytes32(0));
        assertEq(fresh.nextBatch, 0);
    }

    function test_RefusesAZeroAnchor() public {
        vm.expectRevert(IEvidenceRegistry.ZeroAnchorer.selector);
        new EvidenceRegistry(admin, ADMIN_DELAY, address(0), _freshGenesis());
    }

    function test_RefusesARegistryOfNoStreams() public {
        vm.expectRevert(IEvidenceRegistry.NoStreams.selector);
        new EvidenceRegistry(admin, ADMIN_DELAY, anchorer, new StreamGenesis[](0));
    }

    function test_RefusesMoreStreamsThanTheRegistryCanList() public {
        StreamGenesis[] memory many = new StreamGenesis[](33);
        for (uint256 i = 0; i < many.length; i++) {
            many[i] = _start(uint8(i + 1));
        }

        vm.expectRevert(abi.encodeWithSelector(IEvidenceRegistry.TooManyStreams.selector, 33, 32));
        new EvidenceRegistry(admin, ADMIN_DELAY, anchorer, many);

        StreamGenesis[] memory most = new StreamGenesis[](32);
        for (uint256 i = 0; i < most.length; i++) {
            most[i] = _start(uint8(i + 1));
        }

        uint8[] memory listed = new EvidenceRegistry(admin, ADMIN_DELAY, anchorer, most).streams();
        assertEq(listed.length, 32);
        assertEq(listed[0], 1);
        assertEq(listed[31], 32);
    }

    function test_RefusesStreamsThatAreNotAscendingOrAreZero() public {
        StreamGenesis[] memory repeated = _freshGenesis();
        repeated[2].stream = DELIVERY;
        vm.expectRevert(abi.encodeWithSelector(IEvidenceRegistry.StreamsNotAscending.selector, DELIVERY, DELIVERY));
        new EvidenceRegistry(admin, ADMIN_DELAY, anchorer, repeated);

        StreamGenesis[] memory zero = _freshGenesis();
        zero[0].stream = 0;
        vm.expectRevert(abi.encodeWithSelector(IEvidenceRegistry.StreamsNotAscending.selector, 0, 0));
        new EvidenceRegistry(admin, ADMIN_DELAY, anchorer, zero);
    }

    function test_RefusesAGenesisWhoseBatchAndDigestDisagree() public {
        StreamGenesis[] memory digestOnly = _freshGenesis();
        digestOnly[1].digest = keccak256("head");
        vm.expectRevert(
            abi.encodeWithSelector(IEvidenceRegistry.InvalidGenesis.selector, DELIVERY, 0, keccak256("head"))
        );
        new EvidenceRegistry(admin, ADMIN_DELAY, anchorer, digestOnly);

        StreamGenesis[] memory batchOnly = _freshGenesis();
        batchOnly[2].batch = 4;
        vm.expectRevert(abi.encodeWithSelector(IEvidenceRegistry.InvalidGenesis.selector, REVIEW, 4, bytes32(0)));
        new EvidenceRegistry(admin, ADMIN_DELAY, anchorer, batchOnly);
    }

    function test_AnchorsAStreamsFirstBatchAndRecordsIt() public {
        vm.warp(START + 1 days);
        bytes32[][] memory levels = _tree(DELIVERY, 1, 5);
        BatchHeader memory header = _next(DELIVERY, levels, 5);

        vm.expectEmit(true, true, false, true, address(registry));
        emit IEvidenceRegistry.BatchAnchored(DELIVERY, 0, BatchChain.digest(header), header);
        bytes32 digest = _anchor(header);

        Batch memory anchored = registry.batchAt(DELIVERY, 0);
        assertEq(anchored.root, header.root);
        assertEq(anchored.digest, digest);
        assertEq(anchored.leafCount, 5);
        assertEq(anchored.throughTransaction, 1);
        assertEq(anchored.throughSequence, 1);
        assertEq(anchored.closedAt, START + 1 days);
        assertEq(anchored.anchoredAt, START + 1 days);

        StreamState memory state = registry.stateOf(DELIVERY);
        assertEq(state.head, digest);
        assertEq(state.nextBatch, 1);
        assertEq(state.lastClosedAt, START + 1 days);
    }

    function test_EachStreamKeepsItsOwnChain() public {
        _anchorNext(LIFECYCLE, 1, 3);
        _anchorNext(LIFECYCLE, 2, 3);

        assertEq(registry.stateOf(LIFECYCLE).nextBatch, 2);
        assertEq(registry.stateOf(DELIVERY).nextBatch, 0);
        assertEq(registry.stateOf(DELIVERY).head, bytes32(0));

        _anchorNext(DELIVERY, 3, 1);
        assertEq(registry.batchAt(DELIVERY, 0).leafCount, 1);
    }

    function test_EveryRecordOfABatchProvesItsInclusion() public {
        bytes32[][] memory levels = _anchorNext(REVIEW, 7, 5);
        bytes32[] memory salts = _salts(7, 5);
        bytes32[] memory records = _records(7, 5);

        for (uint256 i = 0; i < 5; i++) {
            bytes32[] memory siblings = EvidenceTreeBuilder.proof(levels, i);
            assertTrue(registry.verifyInclusion(REVIEW, 0, salts[i], records[i], uint64(i), siblings));
        }
    }

    function test_AProofThatDoesNotProveInclusionIsFalseNotARevert() public {
        bytes32[][] memory levels = _anchorNext(REVIEW, 7, 5);
        bytes32[] memory salts = _salts(7, 5);
        bytes32[] memory records = _records(7, 5);
        bytes32[] memory siblings = EvidenceTreeBuilder.proof(levels, 2);

        assertFalse(registry.verifyInclusion(REVIEW, 0, keccak256("other salt"), records[2], 2, siblings));
        assertFalse(registry.verifyInclusion(REVIEW, 0, salts[2], keccak256("other record"), 2, siblings));
        assertFalse(registry.verifyInclusion(REVIEW, 0, salts[2], records[2], 3, siblings));
        assertFalse(registry.verifyInclusion(REVIEW, 0, salts[2], records[2], 5, EvidenceTreeBuilder.proof(levels, 5)));
        assertFalse(registry.verifyInclusion(REVIEW, 0, salts[2], records[2], 2, new bytes32[](2)));
    }

    function test_AProofIsBoundToItsStream() public {
        bytes32[][] memory review = _anchorNext(REVIEW, 9, 4);
        _anchorNext(DELIVERY, 9, 4);
        bytes32[] memory salts = _salts(9, 4);
        bytes32[] memory records = _records(9, 4);

        bytes32[] memory siblings = EvidenceTreeBuilder.proof(review, 1);
        assertTrue(registry.verifyInclusion(REVIEW, 0, salts[1], records[1], 1, siblings));
        assertFalse(registry.verifyInclusion(DELIVERY, 0, salts[1], records[1], 1, siblings));
    }

    function test_ABatchThatDoesNotExistIsUnknownRatherThanFalse() public {
        bytes32[] memory none = new bytes32[](0);

        vm.expectRevert(abi.encodeWithSelector(IEvidenceRegistry.UnknownBatch.selector, LIFECYCLE, 0));
        registry.verifyInclusion(LIFECYCLE, 0, bytes32(0), bytes32(0), 0, none);

        vm.expectRevert(abi.encodeWithSelector(IEvidenceRegistry.UnknownStream.selector, 9));
        registry.verifyInclusion(9, 0, bytes32(0), bytes32(0), 0, none);

        vm.expectRevert(abi.encodeWithSelector(IEvidenceRegistry.UnknownBatch.selector, DELIVERY, 0));
        registry.batchAt(DELIVERY, 0);
    }

    function test_AnUnknownStreamTakesNoBatch() public {
        vm.warp(START + 1 days);
        BatchHeader memory header = _next(LIFECYCLE, _tree(LIFECYCLE, 1, 2), 2);
        header.stream = 4;

        vm.expectRevert(abi.encodeWithSelector(IEvidenceRegistry.UnknownStream.selector, 4));
        _anchor(header);
    }

    function test_OnlyTheAnchorerAnchors() public {
        vm.warp(START + 1 days);
        BatchHeader memory header = _next(LIFECYCLE, _tree(LIFECYCLE, 1, 2), 2);
        bytes32 role = registry.ANCHORER_ROLE();

        vm.expectRevert(
            abi.encodeWithSelector(IAccessControl.AccessControlUnauthorizedAccount.selector, stranger, role)
        );
        vm.prank(stranger);
        registry.anchor(header);
    }

    function test_AResentBatchIsAlreadyAnchoredAndADifferentOneIsAConflict() public {
        vm.warp(START + 1 days);
        BatchHeader memory header = _next(LIFECYCLE, _tree(LIFECYCLE, 1, 2), 2);
        bytes32 digest = _anchor(header);

        vm.expectRevert(abi.encodeWithSelector(IEvidenceRegistry.BatchAlreadyAnchored.selector, LIFECYCLE, 0, digest));
        _anchor(header);

        BatchHeader memory other = header;
        other.root = keccak256("another root");
        vm.expectRevert(
            abi.encodeWithSelector(
                IEvidenceRegistry.BatchConflict.selector, LIFECYCLE, 0, digest, BatchChain.digest(other)
            )
        );
        _anchor(other);
    }

    function test_RefusesABatchOutOfSequenceOrNotOnTheHead() public {
        vm.warp(START + 1 days);
        BatchHeader memory skipping = _next(LIFECYCLE, _tree(LIFECYCLE, 1, 2), 2);
        skipping.number = 1;
        vm.expectRevert(abi.encodeWithSelector(IEvidenceRegistry.BatchOutOfSequence.selector, LIFECYCLE, 0, 1));
        _anchor(skipping);

        BatchHeader memory forked = _next(LIFECYCLE, _tree(LIFECYCLE, 1, 2), 2);
        forked.previousDigest = keccak256("somewhere else");
        vm.expectRevert(
            abi.encodeWithSelector(
                IEvidenceRegistry.PreviousDigestMismatch.selector, LIFECYCLE, bytes32(0), keccak256("somewhere else")
            )
        );
        _anchor(forked);
    }

    function test_RefusesAnEmptyBatchAndOneTooLargeForTheScheme() public {
        vm.warp(START + 1 days);
        BatchHeader memory empty = _next(LIFECYCLE, _tree(LIFECYCLE, 1, 1), 1);
        empty.leafCount = 0;
        vm.expectRevert(abi.encodeWithSelector(IEvidenceRegistry.EmptyBatch.selector, LIFECYCLE));
        _anchor(empty);

        BatchHeader memory huge = _next(LIFECYCLE, _tree(LIFECYCLE, 1, 1), 1);
        huge.leafCount = EvidenceTree.MAX_LEAVES + 1;
        vm.expectRevert(
            abi.encodeWithSelector(
                IEvidenceRegistry.TooManyLeaves.selector, EvidenceTree.MAX_LEAVES + 1, EvidenceTree.MAX_LEAVES
            )
        );
        _anchor(huge);
    }

    function test_RefusesACloseTimeNotAfterThePreviousOrInTheFuture() public {
        _anchorNext(LIFECYCLE, 1, 2);
        uint64 previous = registry.stateOf(LIFECYCLE).lastClosedAt;

        BatchHeader memory stale = _next(LIFECYCLE, _tree(LIFECYCLE, 2, 2), 2);
        stale.closedAt = previous;
        vm.expectRevert(
            abi.encodeWithSelector(IEvidenceRegistry.CloseNotAfterPrevious.selector, LIFECYCLE, previous, previous)
        );
        _anchor(stale);

        BatchHeader memory early = _next(LIFECYCLE, _tree(LIFECYCLE, 2, 2), 2);
        early.closedAt = _now() + 1;
        vm.expectRevert(abi.encodeWithSelector(IEvidenceRegistry.CloseInFuture.selector, LIFECYCLE, _now() + 1, _now()));
        _anchor(early);
    }

    function test_RefusesAFeedPositionNotAfterThePrevious() public {
        _anchorNext(LIFECYCLE, 1, 2);
        vm.warp(vm.getBlockTimestamp() + 1 days);

        BatchHeader memory sameSpot = _next(LIFECYCLE, _tree(LIFECYCLE, 2, 2), 2);
        sameSpot.throughTransaction = 1;
        sameSpot.throughSequence = 1;
        vm.expectRevert(
            abi.encodeWithSelector(IEvidenceRegistry.PositionNotAfterPrevious.selector, LIFECYCLE, 1, 1, 1, 1)
        );
        _anchor(sameSpot);

        BatchHeader memory earlierTransaction = _next(LIFECYCLE, _tree(LIFECYCLE, 2, 2), 2);
        earlierTransaction.throughTransaction = 0;
        earlierTransaction.throughSequence = 99;
        vm.expectRevert(
            abi.encodeWithSelector(IEvidenceRegistry.PositionNotAfterPrevious.selector, LIFECYCLE, 0, 99, 1, 1)
        );
        _anchor(earlierTransaction);

        BatchHeader memory laterSequence = _next(LIFECYCLE, _tree(LIFECYCLE, 2, 2), 2);
        laterSequence.throughTransaction = 1;
        laterSequence.throughSequence = 2;
        _anchor(laterSequence);
        assertEq(registry.stateOf(LIFECYCLE).throughSequence, 2);
    }

    function test_SealingStopsEveryStreamAndHandsOverEachHead() public {
        _anchorNext(LIFECYCLE, 1, 2);
        _anchorNext(REVIEW, 2, 3);
        bytes32[] memory heads = _heads();

        vm.expectEmit(true, false, false, true, address(registry));
        emit IEvidenceRegistry.StreamSealed(LIFECYCLE, heads[0], 1, registry.stateOf(LIFECYCLE).lastClosedAt, 1, 1);
        vm.prank(admin);
        registry.seal(heads);

        assertTrue(registry.isSealed());
        vm.warp(vm.getBlockTimestamp() + 1 days);
        BatchHeader memory late = _next(DELIVERY, _tree(DELIVERY, 3, 1), 1);
        vm.expectRevert(IEvidenceRegistry.RegistryIsSealed.selector);
        _anchor(late);

        vm.expectRevert(IEvidenceRegistry.RegistryIsSealed.selector);
        vm.prank(admin);
        registry.seal(heads);
    }

    function test_SealingNeedsTheAdminAndEveryCurrentHead() public {
        _anchorNext(LIFECYCLE, 1, 2);
        bytes32[] memory heads = _heads();

        bytes32 adminRole = registry.DEFAULT_ADMIN_ROLE();
        vm.expectRevert(
            abi.encodeWithSelector(IAccessControl.AccessControlUnauthorizedAccount.selector, anchorer, adminRole)
        );
        vm.prank(anchorer);
        registry.seal(heads);

        bytes32[] memory short = new bytes32[](2);
        vm.expectRevert(abi.encodeWithSelector(IEvidenceRegistry.HeadCountMismatch.selector, 3, 2));
        vm.prank(admin);
        registry.seal(short);

        bytes32[] memory stale = _heads();
        stale[0] = bytes32(0);
        vm.expectRevert(
            abi.encodeWithSelector(IEvidenceRegistry.HeadMismatch.selector, LIFECYCLE, heads[0], bytes32(0))
        );
        vm.prank(admin);
        registry.seal(stale);
    }

    function test_ASealedRegistryStillProves() public {
        bytes32[][] memory levels = _anchorNext(DELIVERY, 4, 3);
        bytes32[] memory heads = _heads();
        vm.prank(admin);
        registry.seal(heads);

        bytes32[] memory siblings = EvidenceTreeBuilder.proof(levels, 2);
        assertTrue(registry.verifyInclusion(DELIVERY, 0, _salts(4, 3)[2], _records(4, 3)[2], 2, siblings));
    }

    function test_ASuccessorContinuesEachStreamExactlyWhereTheSealedOneStopped() public {
        _anchorNext(LIFECYCLE, 1, 2);
        _anchorNext(LIFECYCLE, 2, 2);
        bytes32[] memory heads = _heads();
        vm.prank(admin);
        registry.seal(heads);

        StreamGenesis[] memory genesis = new StreamGenesis[](3);
        uint8[] memory ids = registry.streams();
        for (uint256 i = 0; i < ids.length; i++) {
            StreamState memory last = registry.stateOf(ids[i]);
            genesis[i] = StreamGenesis({
                stream: ids[i],
                batch: last.nextBatch,
                digest: last.head,
                closedAt: last.lastClosedAt,
                throughTransaction: last.throughTransaction,
                throughSequence: last.throughSequence
            });
        }

        EvidenceRegistry predecessor = registry;
        registry = new EvidenceRegistry(admin, ADMIN_DELAY, anchorer, genesis);

        vm.expectRevert(abi.encodeWithSelector(IEvidenceRegistry.UnknownBatch.selector, LIFECYCLE, 1));
        registry.batchAt(LIFECYCLE, 1);

        vm.warp(vm.getBlockTimestamp() + 1 days);
        BatchHeader memory replayed = _next(LIFECYCLE, _tree(LIFECYCLE, 3, 2), 2);
        replayed.number = 1;
        vm.expectRevert(abi.encodeWithSelector(IEvidenceRegistry.BatchBeforeGenesis.selector, LIFECYCLE, 2, 1));
        _anchor(replayed);

        bytes32 digest = _anchor(_next(LIFECYCLE, _tree(LIFECYCLE, 3, 2), 2));
        assertEq(registry.batchAt(LIFECYCLE, 2).digest, digest);
        assertEq(registry.stateOf(DELIVERY).head, predecessor.stateOf(DELIVERY).head);
    }

    function _heads() private view returns (bytes32[] memory) {
        uint8[] memory ids = registry.streams();
        bytes32[] memory heads = new bytes32[](ids.length);
        for (uint256 i = 0; i < ids.length; i++) {
            heads[i] = registry.stateOf(ids[i]).head;
        }

        return heads;
    }
}
