// SPDX-License-Identifier: UNLICENSED
pragma solidity 0.8.37;

import {SumNode} from "./SumNode.sol";

library SumTree {
    uint16 internal constant SCHEME = 1;
    uint8 internal constant MAX_DEPTH = 32;
    uint64 internal constant MAX_LEAVES = uint64(1) << MAX_DEPTH;

    bytes1 private constant LEAF_TAG = 0x00;
    bytes1 private constant NODE_TAG = 0x01;
    bytes1 private constant EMPTY_TAG = 0x02;

    function leaf(uint64 epochNumber, bytes32 commitment, uint64 amount) internal pure returns (SumNode memory) {
        return SumNode({hash: keccak256(abi.encodePacked(LEAF_TAG, epochNumber, commitment, amount)), sum: amount});
    }

    function empty() internal pure returns (SumNode memory) {
        return SumNode({hash: keccak256(abi.encodePacked(EMPTY_TAG)), sum: 0});
    }

    function parent(SumNode memory left, SumNode memory right) internal pure returns (SumNode memory) {
        return SumNode({
            hash: keccak256(abi.encodePacked(NODE_TAG, left.hash, left.sum, right.hash, right.sum)),
            sum: left.sum + right.sum
        });
    }

    function depthOf(uint64 leafCount) internal pure returns (uint8) {
        if (leafCount <= 1) {
            return 0;
        }
        uint64 remaining = leafCount - 1;
        uint8 depth = 0;
        while (remaining != 0) {
            remaining >>= 1;
            depth++;
        }
        return depth;
    }

    function pathSum(uint64 amount, SumNode[] calldata siblings) internal pure returns (uint256) {
        uint256 total = amount;
        for (uint256 level = 0; level < siblings.length; level++) {
            total += siblings[level].sum;
        }
        return total;
    }

    function computeRoot(
        SumNode memory start,
        uint64 index,
        SumNode[] calldata siblings
    ) internal pure returns (SumNode memory) {
        SumNode memory node = start;
        uint64 position = index;
        for (uint256 level = 0; level < siblings.length; level++) {
            if (position & 1 == 0) {
                node = parent(node, siblings[level]);
            } else {
                node = parent(siblings[level], node);
            }
            position >>= 1;
        }
        return node;
    }
}
