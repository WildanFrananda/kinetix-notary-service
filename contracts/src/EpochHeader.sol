// SPDX-License-Identifier: UNLICENSED
pragma solidity 0.8.37;

struct EpochHeader {
    uint64 number;
    bytes32 previousDigest;
    bytes32 root;
    uint64 total;
    uint64 leafCount;
    uint64 asOf;
}
