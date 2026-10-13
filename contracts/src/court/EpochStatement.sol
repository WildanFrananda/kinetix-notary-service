// SPDX-License-Identifier: MIT
pragma solidity 0.8.37;

struct EpochStatement {
    address ledger;
    uint64 epoch;
    bytes32 commitment;
    uint64 amount;
    uint64 index;
    bytes32 siblingsHash;
}
