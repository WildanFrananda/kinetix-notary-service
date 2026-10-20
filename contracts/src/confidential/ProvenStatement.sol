// SPDX-License-Identifier: MIT
pragma solidity 0.8.37;

struct ProvenStatement {
    bytes32 treeRoot;
    bytes32 totalCommitment;
    bool proven;
}
