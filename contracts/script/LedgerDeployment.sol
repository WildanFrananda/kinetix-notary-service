// SPDX-License-Identifier: UNLICENSED
pragma solidity 0.8.37;

struct LedgerDeployment {
    address admin;
    uint48 adminTransferDelay;
    address anchorer;
    uint64 genesisEpoch;
    bytes32 genesisDigest;
    uint64 genesisAsOf;
    bytes32 salt;
}
