// SPDX-License-Identifier: MIT
pragma solidity 0.8.37;

import {SumNode} from "../../src/SumNode.sol";
import {EpochStatement} from "../../src/court/EpochStatement.sol";
import {StatementCourt} from "../../src/court/StatementCourt.sol";
import {StatementHashing} from "../../src/court/StatementHashing.sol";
import {LedgerFixture} from "./LedgerFixture.sol";
import {TreeBuilder} from "./TreeBuilder.sol";

abstract contract CourtFixture is LedgerFixture {
    uint64 internal constant SIGNER_DELAY = 7 days;

    address internal signer;
    uint256 internal signerKey;
    StatementCourt internal court;

    function setUp() public virtual override {
        super.setUp();
        (signer, signerKey) = makeAddrAndKey("statement-signer");
        court = new StatementCourt(admin, ADMIN_DELAY, ledger, signer, SIGNER_DELAY);
    }

    function _anchorNext(uint256 seed, uint256 count) internal returns (SumNode[][] memory) {
        uint64 number = ledger.nextEpoch();
        SumNode[][] memory levels = _tree(number, seed, count);
        _anchor(_header(number, ledger.head(), levels, count, _now() - 1));
        vm.warp(vm.getBlockTimestamp() + 1 hours);
        return levels;
    }

    function _true(
        uint64 number,
        uint256 seed,
        uint256 count,
        uint256 index
    ) internal view returns (EpochStatement memory) {
        SumNode[][] memory levels = _tree(number, seed, count);
        return EpochStatement({
            ledger: address(ledger),
            epoch: number,
            commitment: _commitments(seed, count)[index],
            amount: _amounts(seed, count)[index],
            index: uint64(index),
            siblingsHash: StatementHashing.siblingsHash(TreeBuilder.proof(levels, index))
        });
    }

    function _siblings(
        uint64 number,
        uint256 seed,
        uint256 count,
        uint256 index
    ) internal pure returns (SumNode[] memory) {
        return TreeBuilder.proof(_tree(number, seed, count), index);
    }

    function _sign(uint256 key, EpochStatement memory statement) internal view returns (bytes memory) {
        (uint8 v, bytes32 r, bytes32 s) = vm.sign(key, court.statementDigest(statement));
        return abi.encodePacked(r, s, v);
    }
}
