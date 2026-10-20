// SPDX-License-Identifier: MIT
pragma solidity 0.8.37;

import {SumNode} from "../SumNode.sol";
import {Poseidon2} from "./Poseidon2.sol";

library ConfidentialTree {
    uint16 internal constant SCHEME = 2;
    uint8 internal constant MAX_DEPTH = 16;

    uint256 private constant TAG_SHIFT = 128;
    uint256 private constant LEAF_TAG = 1;
    uint256 private constant EMPTY_TAG = 2;
    bytes1 private constant BINDING_TAG = 0x04;

    function leaf(uint64 epochNumber, bytes32 commitment, uint64 amount) internal pure returns (uint256) {
        uint256 head = (LEAF_TAG << TAG_SHIFT) | (uint256(epochNumber) << 64) | uint256(amount);
        uint256 high = uint256(commitment) >> 128;
        uint256 low = uint256(commitment) & type(uint128).max;
        return Poseidon2.compress(Poseidon2.compress(Poseidon2.compress(0, head), high), low);
    }

    function empty() internal pure returns (uint256) {
        return Poseidon2.compress(0, EMPTY_TAG << TAG_SHIFT);
    }

    function node(uint256 left, uint256 right) internal pure returns (uint256) {
        return Poseidon2.compress(left, right);
    }

    function wellFormed(SumNode[] calldata siblings) internal pure returns (bool) {
        for (uint256 level = 0; level < siblings.length; level++) {
            if (siblings[level].sum != 0 || uint256(siblings[level].hash) >= Poseidon2.PRIME) {
                return false;
            }
        }
        return true;
    }

    function computeRoot(uint256 start, uint64 index, SumNode[] calldata siblings) internal pure returns (uint256) {
        uint256 current = start;
        uint64 position = index;
        for (uint256 level = 0; level < siblings.length; level++) {
            uint256 sibling = uint256(siblings[level].hash);
            if (position & 1 == 0) {
                current = node(current, sibling);
            } else {
                current = node(sibling, current);
            }
            position >>= 1;
        }
        return current;
    }

    function bindingRoot(bytes32 treeRoot, bytes32 totalCommitment) internal pure returns (bytes32) {
        return keccak256(abi.encodePacked(BINDING_TAG, treeRoot, totalCommitment));
    }
}
