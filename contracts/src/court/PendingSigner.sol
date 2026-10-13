// SPDX-License-Identifier: MIT
pragma solidity 0.8.37;

struct PendingSigner {
    address signer;
    uint64 fromEpoch;
    uint64 activatesAt;
}
