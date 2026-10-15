// SPDX-License-Identifier: MIT
pragma solidity 0.8.37;

import {Generation} from "./Generation.sol";

struct PendingGeneration {
    Generation generation;
    uint64 activatesAt;
}
