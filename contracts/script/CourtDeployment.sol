// SPDX-License-Identifier: MIT
pragma solidity 0.8.37;

import {ILiabilityLedger} from "../src/ILiabilityLedger.sol";

struct CourtDeployment {
    address admin;
    uint48 adminTransferDelay;
    ILiabilityLedger ledger;
    address signer;
    uint64 signerChangeDelay;
    bytes32 salt;
}
