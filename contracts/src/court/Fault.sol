// SPDX-License-Identifier: MIT
pragma solidity 0.8.37;

import {FaultKind} from "./FaultKind.sol";

struct Fault {
    FaultKind kind;
    bytes32 evidence;
    address prover;
    uint64 provenAt;
}
