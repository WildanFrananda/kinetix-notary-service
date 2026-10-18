// SPDX-License-Identifier: MIT
pragma solidity 0.8.37;

struct Batch {
    bytes32 root;
    bytes32 digest;
    uint64 leafCount;
    uint64 throughTransaction;
    uint64 throughSequence;
    uint64 closedAt;
    uint64 anchoredAt;
}
