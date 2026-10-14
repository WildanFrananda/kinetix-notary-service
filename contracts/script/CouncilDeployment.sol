// SPDX-License-Identifier: MIT
pragma solidity 0.8.37;

import {ILiabilityLedger} from "../src/ILiabilityLedger.sol";

struct CouncilDeployment {
    address admin;
    uint48 adminTransferDelay;
    ILiabilityLedger ledger;
    address[] members;
    uint8 threshold;
    uint64 membershipChangeDelay;
    bytes32 salt;
}
