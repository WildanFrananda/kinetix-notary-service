// SPDX-License-Identifier: MIT
pragma solidity 0.8.37;

import {SumNode} from "../../src/SumNode.sol";
import {ConfidentialTree} from "../../src/confidential/ConfidentialTree.sol";

library ConfidentialTreeBuilder {
    function root(
        uint64 epochNumber,
        uint256 depth,
        bytes32[] memory commitments,
        uint64[] memory amounts
    ) internal pure returns (uint256) {
        uint256[][] memory levels = build(epochNumber, depth, commitments, amounts);
        return levels[depth][0];
    }

    function build(
        uint64 epochNumber,
        uint256 depth,
        bytes32[] memory commitments,
        uint64[] memory amounts
    ) internal pure returns (uint256[][] memory) {
        uint256[][] memory levels = new uint256[][](depth + 1);
        uint256 width = uint256(1) << depth;
        levels[0] = new uint256[](width);
        for (uint256 i = 0; i < width; i++) {
            if (i < commitments.length) {
                levels[0][i] = ConfidentialTree.leaf(epochNumber, commitments[i], amounts[i]);
            } else {
                levels[0][i] = ConfidentialTree.empty();
            }
        }
        for (uint256 level = 1; level <= depth; level++) {
            uint256 size = width >> level;
            levels[level] = new uint256[](size);
            for (uint256 i = 0; i < size; i++) {
                levels[level][i] = ConfidentialTree.node(levels[level - 1][2 * i], levels[level - 1][2 * i + 1]);
            }
        }
        return levels;
    }

    function siblings(uint256[][] memory levels, uint256 index) internal pure returns (SumNode[] memory) {
        SumNode[] memory path = new SumNode[](levels.length - 1);
        uint256 position = index;
        for (uint256 level = 0; level < path.length; level++) {
            path[level] = SumNode({hash: bytes32(levels[level][position ^ 1]), sum: 0});
            position >>= 1;
        }
        return path;
    }
}
