// SPDX-License-Identifier: MIT
pragma solidity 0.8.37;

import {Test} from "forge-std/src/Test.sol";

import {EpochHeader} from "../../src/EpochHeader.sol";
import {SumNode} from "../../src/SumNode.sol";
import {ConfidentialTree} from "../../src/confidential/ConfidentialTree.sol";
import {LiabilityLedgerV2} from "../../src/confidential/LiabilityLedgerV2.sol";
import {PlonkVerifier} from "../../src/confidential/PlonkVerifier.sol";
import {ConfidentialStatement} from "./ConfidentialStatement.sol";
import {ConfidentialTreeBuilder} from "./ConfidentialTreeBuilder.sol";
import {ScriptedVerifier} from "./ScriptedVerifier.sol";

abstract contract ConfidentialLedgerFixture is Test {
    uint48 internal constant ADMIN_DELAY = 3 days;
    uint256 internal constant START = 1_760_000_000;
    uint8 internal constant DEPTH = 3;

    address internal admin;
    address internal anchorer;
    address internal stranger;
    ScriptedVerifier internal verifier;
    LiabilityLedgerV2 internal ledger;

    function setUp() public virtual {
        admin = makeAddr("admin");
        anchorer = makeAddr("anchorer");
        stranger = makeAddr("stranger");
        vm.warp(START);
        verifier = new ScriptedVerifier();
        ledger = new LiabilityLedgerV2(
            admin, ADMIN_DELAY, anchorer, 0, bytes32(0), 0, PlonkVerifier(address(verifier)), DEPTH
        );
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
        for (uint256 i = 0; i < count; i++) {
            amounts[i] = uint64(uint256(keccak256(abi.encode(seed, "amount", i))) % 1_000_000_000);
        }
        return amounts;
    }

    function _levels(uint64 number, uint256 seed, uint256 count) internal pure returns (uint256[][] memory) {
        return ConfidentialTreeBuilder.build(number, DEPTH, _commitments(seed, count), _amounts(seed, count));
    }

    function _statement(
        uint64 number,
        uint256 seed,
        uint256 count
    ) internal pure returns (ConfidentialStatement memory) {
        uint256[][] memory levels = _levels(number, seed, count);
        return ConfidentialStatement({
            number: number,
            treeRoot: bytes32(levels[DEPTH][0]),
            totalCommitment: keccak256(abi.encode(seed, "total")) >> 8,
            leafCount: uint64(count),
            proof: abi.encodePacked(keccak256(abi.encode(seed, "proof")))
        });
    }

    function _inputs(ConfidentialStatement memory statement) internal pure returns (uint256[] memory) {
        uint256[] memory inputs = new uint256[](4);
        inputs[0] = statement.number;
        inputs[1] = uint256(statement.treeRoot);
        inputs[2] = uint256(statement.totalCommitment);
        inputs[3] = statement.leafCount;
        return inputs;
    }

    function _prove(ConfidentialStatement memory statement) internal returns (bytes32) {
        verifier.accept(statement.proof, _inputs(statement));
        return ledger.proveEpoch(
            statement.number, statement.treeRoot, statement.totalCommitment, statement.leafCount, statement.proof
        );
    }

    function _header(
        ConfidentialStatement memory statement,
        bytes32 previousDigest,
        uint64 asOf
    ) internal pure returns (EpochHeader memory) {
        return EpochHeader({
            number: statement.number,
            previousDigest: previousDigest,
            root: ConfidentialTree.bindingRoot(statement.treeRoot, statement.totalCommitment),
            total: 0,
            leafCount: statement.leafCount,
            asOf: asOf
        });
    }

    function _proveAndAnchor(
        ConfidentialStatement memory statement,
        bytes32 previousDigest
    ) internal returns (bytes32) {
        _prove(statement);
        return _anchor(_header(statement, previousDigest, _now() - 60));
    }

    function _anchor(EpochHeader memory header) internal returns (bytes32) {
        vm.prank(anchorer);
        return ledger.anchor(header);
    }

    function _siblings(
        uint64 number,
        uint256 seed,
        uint256 count,
        uint256 index
    ) internal pure returns (SumNode[] memory) {
        return ConfidentialTreeBuilder.siblings(_levels(number, seed, count), index);
    }

    function _now() internal view returns (uint64) {
        return uint64(vm.getBlockTimestamp());
    }
}
