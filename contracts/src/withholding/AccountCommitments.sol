// SPDX-License-Identifier: MIT
pragma solidity 0.8.37;

library AccountCommitments {
    bytes1 internal constant COMMITMENT_TAG = 0x06;

    function commitment(bytes32 account, bytes32 salt) internal pure returns (bytes32) {
        return keccak256(abi.encodePacked(COMMITMENT_TAG, account, salt));
    }
}
