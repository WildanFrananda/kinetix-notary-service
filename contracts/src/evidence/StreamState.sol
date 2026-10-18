// SPDX-License-Identifier: MIT
pragma solidity 0.8.37;

struct StreamState {
    bool known;
    bytes32 head;
    uint64 genesisBatch;
    uint64 nextBatch;
    uint64 lastClosedAt;
    uint64 throughTransaction;
    uint64 throughSequence;
}
