// SPDX-License-Identifier: MIT
pragma solidity 0.8.37;

struct PendingMembership {
    address[] members;
    uint8 threshold;
    uint64 activatesAt;
}
