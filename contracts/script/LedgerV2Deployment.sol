// SPDX-License-Identifier: MIT
pragma solidity 0.8.37;

import {ILiabilityLedger} from "../src/ILiabilityLedger.sol";

struct LedgerV2Deployment {
    address admin;
    uint48 adminTransferDelay;
    address anchorer;
    ILiabilityLedger predecessor;
    uint8 depth;
    bytes32 verifierSalt;
    bytes32 salt;
}
