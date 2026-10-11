// SPDX-License-Identifier: UNLICENSED
pragma solidity 0.8.37;

import {Test} from "forge-std/src/Test.sol";

import {EpochHeader} from "../../src/EpochHeader.sol";
import {LiabilityLedger} from "../../src/LiabilityLedger.sol";
import {SumNode} from "../../src/SumNode.sol";
import {TreeBuilder} from "./TreeBuilder.sol";

abstract contract LedgerFixture is Test {
    uint48 internal constant ADMIN_DELAY = 3 days;
    uint256 internal constant START = 1_760_000_000;

    address internal admin;
    address internal anchorer;
    address internal stranger;
    LiabilityLedger internal ledger;

    function setUp() public virtual {
        admin = makeAddr("admin");
        anchorer = makeAddr("anchorer");
        stranger = makeAddr("stranger");
        vm.warp(START);
        ledger = new LiabilityLedger(admin, ADMIN_DELAY, anchorer, 0, bytes32(0), 0);
    }

    function _commitments(uint256 seed, uint256 count) internal pure returns (bytes32[] memory) {
        bytes32[] memory commitments = new bytes32[](count);
        for (uint256 i = 0; i < count; i++) {
            commitments[i] = keccak256(abi.encode(seed, "commitment", i));
        }
        return commitments;
    }

    function _amounts(uint256 seed, uint256 count) internal pure returns (uint64[] memory) {
        uint64[] memory amounts = new uint64[](count);
        uint256 cap = count == 0 ? 0 : uint256(type(uint64).max) / count;
        for (uint256 i = 0; i < count; i++) {
            amounts[i] = uint64(uint256(keccak256(abi.encode(seed, "amount", i))) % (cap + 1));
        }
        return amounts;
    }

    function _tree(uint64 number, uint256 seed, uint256 count) internal pure returns (SumNode[][] memory) {
        return TreeBuilder.build(number, _commitments(seed, count), _amounts(seed, count));
    }

    function _header(
        uint64 number,
        bytes32 previousDigest,
        SumNode[][] memory levels,
        uint256 count,
        uint64 asOf
    ) internal pure returns (EpochHeader memory) {
        SumNode memory root = TreeBuilder.root(levels);
        return EpochHeader({
            number: number,
            previousDigest: previousDigest,
            root: root.hash,
            total: root.sum,
            leafCount: uint64(count),
            asOf: asOf
        });
    }

    function _anchor(EpochHeader memory header) internal returns (bytes32) {
        vm.prank(anchorer);
        return ledger.anchor(header);
    }

    function _now() internal view returns (uint64) {
        return uint64(vm.getBlockTimestamp());
    }
}
