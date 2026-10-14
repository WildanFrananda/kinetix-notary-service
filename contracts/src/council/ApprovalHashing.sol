// SPDX-License-Identifier: MIT
pragma solidity 0.8.37;

library ApprovalHashing {
    string internal constant NAME = "Kinetix Anchor Council";
    string internal constant VERSION = "1";

    bytes32 internal constant APPROVAL_TYPEHASH =
        keccak256("EpochApproval(address ledger,uint64 number,bytes32 digest)");

    function structHash(address ledger, uint64 number, bytes32 digest) internal pure returns (bytes32) {
        return keccak256(abi.encode(APPROVAL_TYPEHASH, ledger, number, digest));
    }
}
