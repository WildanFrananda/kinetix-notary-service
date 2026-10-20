// SPDX-License-Identifier: MIT
pragma solidity 0.8.37;

import {Test} from "forge-std/src/Test.sol";

import {EpochHeader} from "../../src/EpochHeader.sol";
import {SumNode} from "../../src/SumNode.sol";
import {ConfidentialTree} from "../../src/confidential/ConfidentialTree.sol";
import {IConfidentialLedger} from "../../src/confidential/IConfidentialLedger.sol";
import {LiabilityLedgerV2} from "../../src/confidential/LiabilityLedgerV2.sol";
import {PlonkVerifier} from "../../src/confidential/PlonkVerifier.sol";
import {TestPlonkVerifier} from "../support/TestPlonkVerifier.sol";

contract LiabilityLedgerV2PlonkTest is Test {
    string private fixture;
    LiabilityLedgerV2 private ledger;
    address private anchorer;
    uint64 private number;
    bytes32 private treeRoot;
    bytes32 private totalCommitment;
    uint64 private leafCount;
    bytes private proof;

    function setUp() public {
        fixture = vm.readFile(string.concat(vm.projectRoot(), "/test/fixtures/plonk-v2.json"));
        anchorer = makeAddr("anchorer");
        number = uint64(vm.parseJsonUint(fixture, ".epoch"));
        treeRoot = vm.parseJsonBytes32(fixture, ".treeRoot");
        totalCommitment = vm.parseJsonBytes32(fixture, ".totalCommitment");
        leafCount = uint64(vm.parseJsonUint(fixture, ".leafCount"));
        proof = vm.parseJsonBytes(fixture, ".proof");

        uint64 asOf = uint64(vm.parseJsonUint(fixture, ".asOf"));
        vm.warp(asOf + 60);
        ledger = new LiabilityLedgerV2(
            makeAddr("admin"),
            3 days,
            anchorer,
            number,
            vm.parseJsonBytes32(fixture, ".previousDigest"),
            asOf - 3600,
            PlonkVerifier(address(new TestPlonkVerifier())),
            uint8(vm.parseJsonUint(fixture, ".depth"))
        );
    }

    function test_AGoProofOfAnEpochIsAcceptedAndAnchored() public {
        bytes32 root = ledger.proveEpoch(number, treeRoot, totalCommitment, leafCount, proof);
        assertEq(root, vm.parseJsonBytes32(fixture, ".root"));

        vm.prank(anchorer);
        bytes32 digest = ledger.anchor(_header());

        assertEq(digest, vm.parseJsonBytes32(fixture, ".digest"));
        assertEq(ledger.treeRootOf(number), treeRoot);
        assertEq(ledger.totalCommitmentOf(number), totalCommitment);
        assertEq(ledger.epochAt(number).total, 0);
    }

    function test_EveryAccountOfTheProvenEpochVerifiesItsInclusion() public {
        ledger.proveEpoch(number, treeRoot, totalCommitment, leafCount, proof);
        vm.prank(anchorer);
        ledger.anchor(_header());

        bytes32[] memory commitments = vm.parseJsonBytes32Array(fixture, ".commitments");
        uint256[] memory amounts = vm.parseJsonUintArray(fixture, ".amounts");
        assertEq(commitments.length, leafCount);
        for (uint256 i = 0; i < commitments.length; i++) {
            string memory key = string.concat(".inclusions[", vm.toString(i), "]");
            bytes32[] memory hashes = vm.parseJsonBytes32Array(fixture, string.concat(key, ".siblingHashes"));
            uint256[] memory sums = vm.parseJsonUintArray(fixture, string.concat(key, ".siblingSums"));
            SumNode[] memory siblings = new SumNode[](hashes.length);
            for (uint256 level = 0; level < hashes.length; level++) {
                assertEq(sums[level], 0);
                siblings[level] = SumNode({hash: hashes[level], sum: 0});
            }

            uint64 index = uint64(vm.parseJsonUint(fixture, string.concat(key, ".index")));
            assertTrue(ledger.verifyInclusion(number, commitments[i], uint64(amounts[i]), index, siblings), key);
            assertFalse(ledger.verifyInclusion(number, commitments[i], uint64(amounts[i]) + 1, index, siblings), key);
        }
    }

    function test_TheRealVerifierRefusesTheProofForAnyOtherStatement() public {
        bytes32 otherRoot = bytes32(uint256(treeRoot) + 1);
        bytes32 otherCommitment = bytes32(uint256(totalCommitment) + 1);

        vm.expectRevert(abi.encodeWithSelector(IConfidentialLedger.ProofRejected.selector, number));
        ledger.proveEpoch(number, otherRoot, totalCommitment, leafCount, proof);

        vm.expectRevert(abi.encodeWithSelector(IConfidentialLedger.ProofRejected.selector, number));
        ledger.proveEpoch(number, treeRoot, otherCommitment, leafCount, proof);

        vm.expectRevert(abi.encodeWithSelector(IConfidentialLedger.ProofRejected.selector, number));
        ledger.proveEpoch(number, treeRoot, totalCommitment, leafCount - 1, proof);

        vm.expectRevert(abi.encodeWithSelector(IConfidentialLedger.ProofRejected.selector, number + 1));
        ledger.proveEpoch(number + 1, treeRoot, totalCommitment, leafCount, proof);
    }

    function test_AProofWithABrokenPointIsRefusedNotPanicked() public {
        bytes memory broken = bytes.concat(proof);
        broken[31] ^= 0x01;

        vm.expectRevert(abi.encodeWithSelector(IConfidentialLedger.ProofRejected.selector, number));
        ledger.proveEpoch{gas: 3_000_000}(number, treeRoot, totalCommitment, leafCount, broken);
    }

    function test_AnUnprovenHeaderIsNeverAnchoredEvenWithAValidProofOnFile() public {
        ledger.proveEpoch(number, treeRoot, totalCommitment, leafCount, proof);
        EpochHeader memory header = _header();
        header.root = ConfidentialTree.bindingRoot(treeRoot, bytes32(uint256(totalCommitment) + 1));

        vm.expectRevert(
            abi.encodeWithSelector(IConfidentialLedger.EpochNotProven.selector, number, header.root, leafCount)
        );
        vm.prank(anchorer);
        ledger.anchor(header);
    }

    function _header() private view returns (EpochHeader memory) {
        return EpochHeader({
            number: number,
            previousDigest: vm.parseJsonBytes32(fixture, ".previousDigest"),
            root: vm.parseJsonBytes32(fixture, ".root"),
            total: 0,
            leafCount: leafCount,
            asOf: uint64(vm.parseJsonUint(fixture, ".asOf"))
        });
    }
}
