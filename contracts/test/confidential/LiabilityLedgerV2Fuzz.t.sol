// SPDX-License-Identifier: MIT
pragma solidity 0.8.37;

import {EpochHeader} from "../../src/EpochHeader.sol";
import {SumNode} from "../../src/SumNode.sol";
import {Poseidon2} from "../../src/confidential/Poseidon2.sol";
import {ConfidentialLedgerFixture} from "../support/ConfidentialLedgerFixture.sol";

contract LiabilityLedgerV2FuzzTest is ConfidentialLedgerFixture {
    uint256 private constant SEED = 4;
    uint256 private constant COUNT = 6;

    function setUp() public override {
        super.setUp();
        _proveAndAnchor(_statement(0, SEED, COUNT), bytes32(0));
    }

    function testFuzz_AnInclusionNeverHoldsForAnotherAmount(uint8 rawIndex, uint64 other) public view {
        uint256 index = bound(rawIndex, 0, COUNT - 1);
        bytes32[] memory commitments = _commitments(SEED, COUNT);
        uint64[] memory amounts = _amounts(SEED, COUNT);
        vm.assume(other != amounts[index]);

        assertFalse(
            ledger.verifyInclusion(0, commitments[index], other, uint64(index), _siblings(0, SEED, COUNT, index))
        );
    }

    function testFuzz_AnInclusionNeverHoldsForAnotherCommitment(uint8 rawIndex, bytes32 other) public view {
        uint256 index = bound(rawIndex, 0, COUNT - 1);
        bytes32[] memory commitments = _commitments(SEED, COUNT);
        vm.assume(other != commitments[index]);

        assertFalse(
            ledger.verifyInclusion(
                0, other, _amounts(SEED, COUNT)[index], uint64(index), _siblings(0, SEED, COUNT, index)
            )
        );
    }

    function testFuzz_ASiblingOutsideTheFieldIsNeverAccepted(uint8 rawIndex, uint8 rawLevel, uint256 raw) public view {
        uint256 index = bound(rawIndex, 0, COUNT - 1);
        SumNode[] memory siblings = _siblings(0, SEED, COUNT, index);
        siblings[bound(rawLevel, 0, siblings.length - 1)].hash = bytes32(bound(raw, Poseidon2.PRIME, type(uint256).max));

        assertFalse(
            ledger.verifyInclusion(
                0, _commitments(SEED, COUNT)[index], _amounts(SEED, COUNT)[index], uint64(index), siblings
            )
        );
    }

    function testFuzz_AHeaderNobodyProvedIsNeverAnchored(bytes32 root, uint64 rawLeafCount) public {
        uint64 leafCount = uint64(bound(rawLeafCount, 0, ledger.capacity()));
        vm.assume(!ledger.isProven(1, root, leafCount));
        vm.warp(vm.getBlockTimestamp() + 1 hours);
        EpochHeader memory header = EpochHeader({
            number: 1,
            previousDigest: ledger.head(),
            root: root,
            total: 0,
            leafCount: leafCount,
            asOf: _now() - 1
        });

        vm.prank(anchorer);
        try ledger.anchor(header) {
            fail();
        } catch {}
    }
}
