// SPDX-License-Identifier: UNLICENSED
pragma solidity 0.8.37;

import {SumNode} from "../../src/SumNode.sol";
import {SumTree} from "../../src/SumTree.sol";

library TreeBuilder {
    function build(uint64 epochNumber, bytes32[] memory commitments, uint64[] memory amounts)
        internal
        pure
        returns (SumNode[][] memory)
    {
        uint256 count = commitments.length;
        uint8 depth = SumTree.depthOf(uint64(count));
        uint256 width = uint256(1) << depth;
        SumNode[][] memory levels = new SumNode[][](uint256(depth) + 1);
        levels[0] = new SumNode[](width);
        for (uint256 i = 0; i < width; i++) {
            if (i < count) {
                levels[0][i] = SumTree.leaf(epochNumber, commitments[i], amounts[i]);
            } else {
                levels[0][i] = SumTree.empty();
            }
        }
        for (uint256 level = 1; level <= depth; level++) {
            uint256 size = width >> level;
            levels[level] = new SumNode[](size);
            for (uint256 i = 0; i < size; i++) {
                levels[level][i] = SumTree.parent(levels[level - 1][2 * i], levels[level - 1][2 * i + 1]);
            }
        }
        return levels;
    }

    function root(SumNode[][] memory levels) internal pure returns (SumNode memory) {
        return levels[levels.length - 1][0];
    }

    function proof(SumNode[][] memory levels, uint256 index) internal pure returns (SumNode[] memory) {
        uint256 depth = levels.length - 1;
        SumNode[] memory siblings = new SumNode[](depth);
        uint256 position = index;
        for (uint256 level = 0; level < depth; level++) {
            siblings[level] = levels[level][position ^ 1];
            position >>= 1;
        }
        return siblings;
    }
}
