// SPDX-License-Identifier: MIT
pragma solidity 0.8.37;

struct Generation {
    uint64 chainId;
    address ledger;
    address court;
    address withholding;
    uint64 genesisEpoch;
    bytes32 genesisDigest;
    uint16 scheme;
    uint64 registeredAt;
}
