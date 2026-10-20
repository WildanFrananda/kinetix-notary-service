// SPDX-License-Identifier: MIT
pragma solidity 0.8.37;

struct ConfidentialStatement {
    uint64 number;
    bytes32 treeRoot;
    bytes32 totalCommitment;
    uint64 leafCount;
    bytes proof;
}
