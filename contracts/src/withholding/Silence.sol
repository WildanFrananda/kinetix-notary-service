// SPDX-License-Identifier: MIT
pragma solidity 0.8.37;

struct Silence {
    bytes32 registration;
    bytes32 account;
    address holder;
    uint64 dueAt;
    address recorder;
    uint64 recordedAt;
}
