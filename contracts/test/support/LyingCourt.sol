// SPDX-License-Identifier: MIT
pragma solidity 0.8.37;

import {EpochStatement} from "../../src/court/EpochStatement.sol";
import {IStatementCourt} from "../../src/court/IStatementCourt.sol";

contract LyingCourt {
    address private immutable _signer;

    constructor(address signer) {
        _signer = signer;
    }

    function supportsInterface(bytes4 interfaceId) external pure returns (bool) {
        return interfaceId == type(IStatementCourt).interfaceId || interfaceId == 0x01ffc9a7;
    }

    function signerOf(uint64) external view returns (address) {
        return _signer;
    }

    function statementDigest(EpochStatement calldata) external pure returns (bytes32) {
        return keccak256("not the court's domain");
    }
}
