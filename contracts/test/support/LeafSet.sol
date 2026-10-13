// SPDX-License-Identifier: MIT
pragma solidity 0.8.37;

library LeafSet {
    function commitments(uint256 seed, uint256 count) internal pure returns (bytes32[] memory) {
        bytes32[] memory out = new bytes32[](count);
        for (uint256 i = 0; i < count; i++) {
            out[i] = keccak256(abi.encode(seed, "commitment", i));
        }

        return out;
    }

    function amounts(uint256 seed, uint256 count) internal pure returns (uint64[] memory) {
        uint64[] memory out = new uint64[](count);
        uint256 cap = count == 0 ? 0 : uint256(type(uint64).max) / count;
        for (uint256 i = 0; i < count; i++) {
            out[i] = uint64(uint256(keccak256(abi.encode(seed, "amount", i))) % (cap + 1));
        }

        return out;
    }
}
