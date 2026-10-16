// SPDX-License-Identifier: MIT
pragma solidity 0.8.37;

import {IEntryPoint} from "account-abstraction/contracts/interfaces/IEntryPoint.sol";

import {IStatementCourt} from "../src/court/IStatementCourt.sol";
import {IWithholdingCourt} from "../src/withholding/IWithholdingCourt.sol";

struct PaymasterDeployment {
    address admin;
    uint48 adminTransferDelay;
    IEntryPoint entryPoint;
    address account;
    IStatementCourt court;
    IWithholdingCourt withholding;
    uint256 maxCostPerOperation;
    uint32 unstakeDelay;
    bytes32 salt;
}
