// SPDX-License-Identifier: MIT
pragma solidity 0.8.37;

import {IStatementCourt} from "../src/court/IStatementCourt.sol";

struct WithholdingDeployment {
    IStatementCourt court;
    uint64 requestWindow;
    uint64 answerPeriod;
    bytes32 salt;
}
