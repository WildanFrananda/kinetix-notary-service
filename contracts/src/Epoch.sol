// SPDX-License-Identifier: UNLICENSED
pragma solidity 0.8.37;

struct Epoch {
    bytes32 root;
    bytes32 digest;
    uint64 total;
    uint64 leafCount;
    uint64 asOf;
    uint64 anchoredAt;
}
