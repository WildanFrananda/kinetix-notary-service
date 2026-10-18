// SPDX-License-Identifier: MIT
pragma solidity 0.8.37;

struct BatchHeader {
    uint8 stream;
    uint64 number;
    bytes32 previousDigest;
    bytes32 root;
    uint64 leafCount;
    uint64 throughTransaction;
    uint64 throughSequence;
    uint64 closedAt;
}
