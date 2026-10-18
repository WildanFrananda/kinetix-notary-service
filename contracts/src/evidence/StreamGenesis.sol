// SPDX-License-Identifier: MIT
pragma solidity 0.8.37;

struct StreamGenesis {
    uint8 stream;
    uint64 batch;
    bytes32 digest;
    uint64 closedAt;
    uint64 throughTransaction;
    uint64 throughSequence;
}
