// SPDX-License-Identifier: MIT
pragma solidity 0.8.37;

import {ILiabilityLedger} from "../src/ILiabilityLedger.sol";
import {IStatementCourt} from "../src/court/IStatementCourt.sol";
import {IWithholdingCourt} from "../src/withholding/IWithholdingCourt.sol";

struct DirectoryDeployment {
    address admin;
    uint48 adminTransferDelay;
    ILiabilityLedger ledger;
    IStatementCourt court;
    IWithholdingCourt withholding;
    uint64 foreignSuccessorDelay;
    bytes32 salt;
}
