// SPDX-License-Identifier: UNLICENSED
pragma solidity 0.8.37;

import {EpochHeader} from "../src/EpochHeader.sol";
import {SumNode} from "../src/SumNode.sol";
import {SumTree} from "../src/SumTree.sol";
import {LedgerFixture} from "./support/LedgerFixture.sol";
import {TreeBuilder} from "./support/TreeBuilder.sol";

contract SumTreeTest is LedgerFixture {
    uint256 private constant MAX_FUZZ_LEAVES = 70;

    function testFuzz_EveryAccountProvesItsInclusion(uint256 seed, uint256 rawCount) public {
        uint256 count = bound(rawCount, 1, MAX_FUZZ_LEAVES);
        SumNode[][] memory levels = _tree(0, seed, count);
        _anchor(_header(0, bytes32(0), levels, count, _now() - 1));
        bytes32[] memory commitments = _commitments(seed, count);
        uint64[] memory amounts = _amounts(seed, count);

        for (uint256 i = 0; i < count; i++) {
            assertTrue(ledger.verifyInclusion(0, commitments[i], amounts[i], uint64(i), TreeBuilder.proof(levels, i)));
        }
    }

    function testFuzz_TheRootSumIsTheTotalOfEveryAccount(uint256 seed, uint256 rawCount) public pure {
        uint256 count = bound(rawCount, 0, MAX_FUZZ_LEAVES);
        uint64[] memory amounts = _amounts(seed, count);
        uint256 expected = 0;
        for (uint256 i = 0; i < count; i++) {
            expected += amounts[i];
        }

        assertEq(TreeBuilder.root(_tree(0, seed, count)).sum, expected);
    }

    function testFuzz_AnyChangeToTheClaimIsRejected(
        uint256 seed,
        uint256 rawCount,
        uint256 rawIndex,
        uint256 rawLevel,
        uint8 tamper,
        uint64 delta
    ) public {
        uint256 count = bound(rawCount, 2, MAX_FUZZ_LEAVES);
        uint256 index = bound(rawIndex, 0, count - 1);
        delta = uint64(bound(delta, 1, type(uint64).max));
        SumNode[][] memory levels = _tree(0, seed, count);
        _anchor(_header(0, bytes32(0), levels, count, _now() - 1));
        bytes32 commitment = _commitments(seed, count)[index];
        uint64 amount = _amounts(seed, count)[index];
        SumNode[] memory siblings = TreeBuilder.proof(levels, index);
        uint256 level = bound(rawLevel, 0, siblings.length - 1);
        uint64 claimedIndex = uint64(index);

        uint8 kind = tamper % 5;
        if (kind == 0) {
            commitment = keccak256(abi.encode(commitment, delta));
        } else if (kind == 1) {
            unchecked {
                amount += delta;
            }
        } else if (kind == 2) {
            siblings[level].hash = keccak256(abi.encode(siblings[level].hash, delta));
        } else if (kind == 3) {
            unchecked {
                siblings[level].sum += delta;
            }
        } else {
            claimedIndex = uint64(index ^ (uint256(1) << level));
        }

        assertFalse(ledger.verifyInclusion(0, commitment, amount, claimedIndex, siblings));
    }

    function testFuzz_ValueMovedBetweenAccountAndSiblingIsRejected(
        uint256 seed,
        uint256 rawCount,
        uint256 rawIndex,
        uint256 rawLevel,
        uint64 rawMoved
    ) public {
        uint256 count = bound(rawCount, 2, MAX_FUZZ_LEAVES);
        uint256 index = bound(rawIndex, 0, count - 1);
        SumNode[] memory siblings = _anchoredProof(seed, count, index);
        rawLevel = bound(rawLevel, 0, siblings.length - 1);
        vm.assume(siblings[rawLevel].sum > 0);
        uint64 moved = uint64(bound(rawMoved, 1, siblings[rawLevel].sum));
        siblings[rawLevel].sum -= moved;
        uint64 claimed = _amounts(seed, count)[index] + moved;

        assertFalse(ledger.verifyInclusion(0, _commitments(seed, count)[index], claimed, uint64(index), siblings));
    }

    function testFuzz_DepthIsTheSmallestThatHoldsEveryLeaf(uint64 leafCount) public pure {
        uint8 depth = SumTree.depthOf(leafCount);

        if (leafCount <= 1) {
            assertEq(depth, 0);
            return;
        }
        assertGe(uint256(1) << depth, leafCount);
        assertLt(uint256(1) << (depth - 1), leafCount);
    }

    function testFuzz_TheSameAccountsInAnotherEpochCommitToADifferentRoot(uint256 seed, uint256 rawCount, uint64 other)
        public
        pure
    {
        uint256 count = bound(rawCount, 1, MAX_FUZZ_LEAVES);
        vm.assume(other != 0);

        assertNotEq(TreeBuilder.root(_tree(0, seed, count)).hash, TreeBuilder.root(_tree(other, seed, count)).hash);
    }

    function testFuzz_AProofFromAnotherEpochIsRejected(uint256 seed, uint256 rawCount, uint256 rawIndex) public {
        uint256 count = bound(rawCount, 2, MAX_FUZZ_LEAVES);
        uint256 index = bound(rawIndex, 0, count - 1);
        bytes32 head = _anchor(_header(0, bytes32(0), _tree(0, seed, count), count, _now() - 2));
        SumNode[][] memory later = _tree(1, seed, count);
        _anchor(_header(1, head, later, count, _now() - 1));

        assertFalse(
            ledger.verifyInclusion(
                0,
                _commitments(seed, count)[index],
                _amounts(seed, count)[index],
                uint64(index),
                TreeBuilder.proof(later, index)
            )
        );
    }

    function _anchoredProof(uint256 seed, uint256 count, uint256 index) private returns (SumNode[] memory) {
        SumNode[][] memory levels = _tree(0, seed, count);
        _anchor(_header(0, bytes32(0), levels, count, _now() - 1));
        return TreeBuilder.proof(levels, index);
    }
}
