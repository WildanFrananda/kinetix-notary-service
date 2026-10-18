// SPDX-License-Identifier: MIT
pragma solidity 0.8.37;

library EvidenceTree {
    uint16 internal constant SCHEME = 1;
    uint8 internal constant MAX_DEPTH = 32;
    uint64 internal constant MAX_LEAVES = uint64(1) << MAX_DEPTH;

    bytes1 private constant LEAF_TAG = 0x10;
    bytes1 private constant NODE_TAG = 0x11;
    bytes1 private constant EMPTY_TAG = 0x12;

    function leaf(uint8 stream, bytes32 salt, bytes32 recordDigest) internal pure returns (bytes32) {
        return keccak256(abi.encodePacked(LEAF_TAG, stream, salt, recordDigest));
    }

    function empty() internal pure returns (bytes32) {
        return keccak256(abi.encodePacked(EMPTY_TAG));
    }

    function parent(bytes32 left, bytes32 right) internal pure returns (bytes32) {
        return keccak256(abi.encodePacked(NODE_TAG, left, right));
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

    function computeRoot(bytes32 start, uint64 index, bytes32[] calldata siblings) internal pure returns (bytes32) {
        bytes32 node = start;
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
