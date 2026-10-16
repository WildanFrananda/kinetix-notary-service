// SPDX-License-Identifier: MIT
pragma solidity 0.8.37;

import {IStatementCourt} from "../../src/court/IStatementCourt.sol";
import {IWithholdingCourt} from "../../src/withholding/IWithholdingCourt.sol";
import {Registration} from "../../src/withholding/Registration.sol";

contract LyingWithholding {
    IStatementCourt private immutable _court;

    constructor(IStatementCourt court_) {
        _court = court_;
    }

    function supportsInterface(bytes4 interfaceId) external pure returns (bool) {
        return interfaceId == type(IWithholdingCourt).interfaceId || interfaceId == 0x01ffc9a7;
    }

    function court() external view returns (IStatementCourt) {
        return _court;
    }

    function registrationDigest(Registration calldata) external pure returns (bytes32) {
        return keccak256("not the withholding court's domain");
    }
}
