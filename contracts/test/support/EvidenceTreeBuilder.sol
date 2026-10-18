// SPDX-License-Identifier: MIT
pragma solidity 0.8.37;

import {EvidenceTree} from "../../src/evidence/EvidenceTree.sol";

library EvidenceTreeBuilder {
    function build(
        uint8 stream,
        bytes32[] memory salts,
        bytes32[] memory recordDigests
    ) internal pure returns (bytes32[][] memory) {
        uint256 count = salts.length;
        uint8 depth = EvidenceTree.depthOf(uint64(count));
        uint256 width = uint256(1) << depth;
        bytes32[][] memory levels = new bytes32[][](uint256(depth) + 1);
        levels[0] = new bytes32[](width);
        for (uint256 i = 0; i < width; i++) {
            if (i < count) {
                levels[0][i] = EvidenceTree.leaf(stream, salts[i], recordDigests[i]);
            } else {
                levels[0][i] = EvidenceTree.empty();
            }
        }

        for (uint256 level = 1; level <= depth; level++) {
            uint256 size = width >> level;
            levels[level] = new bytes32[](size);
            for (uint256 i = 0; i < size; i++) {
                levels[level][i] = EvidenceTree.parent(levels[level - 1][2 * i], levels[level - 1][2 * i + 1]);
            }
        }

        return levels;
    }

    function root(bytes32[][] memory levels) internal pure returns (bytes32) {
        return levels[levels.length - 1][0];
    }

    function proof(bytes32[][] memory levels, uint256 index) internal pure returns (bytes32[] memory) {
        uint256 depth = levels.length - 1;
        bytes32[] memory siblings = new bytes32[](depth);
        uint256 position = index;
        for (uint256 level = 0; level < depth; level++) {
            siblings[level] = levels[level][position ^ 1];
            position >>= 1;
        }

        return siblings;
    }
}
