// SPDX-License-Identifier: MIT
pragma solidity 0.8.37;

import {RequestState} from "./RequestState.sol";

struct StatementRequest {
    RequestState state;
    uint64 dueAt;
}
