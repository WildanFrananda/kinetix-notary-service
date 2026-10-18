// SPDX-License-Identifier: MIT
pragma solidity 0.8.37;

import {IEvidenceRegistry} from "../src/evidence/IEvidenceRegistry.sol";

struct EvidenceDeployment {
    address admin;
    uint48 adminTransferDelay;
    address anchorer;
    IEvidenceRegistry predecessor;
    bytes32 salt;
}
