// SPDX-License-Identifier: MIT
pragma solidity 0.8.37;

import {Test} from "forge-std/src/Test.sol";

import {EpochHeader} from "../../src/EpochHeader.sol";
import {ConfidentialTree} from "../../src/confidential/ConfidentialTree.sol";
import {IConfidentialLedger} from "../../src/confidential/IConfidentialLedger.sol";
import {LiabilityLedgerV2} from "../../src/confidential/LiabilityLedgerV2.sol";
import {PlonkVerifier} from "../../src/confidential/PlonkVerifier.sol";
import {TestPlonkVerifier} from "../support/TestPlonkVerifier.sol";

contract ProductionPlonkVerifierTest is Test {
    uint256 private constant VERIFY_GAS = 1_000_000;

    string private fixture;
    PlonkVerifier private verifier;
    bytes private proof;
    uint256[] private inputs;

    function setUp() public {
        fixture = vm.readFile(string.concat(vm.projectRoot(), "/test/fixtures/plonk-production.json"));
        verifier = new PlonkVerifier();
        proof = vm.parseJsonBytes(fixture, ".proof");
        inputs = vm.parseJsonUintArray(fixture, ".publicInputs");
    }

    function test_TheFixtureIsForGenerationTwosDepth() public view {
        assertEq(vm.parseJsonUint(fixture, ".depth"), 8);
        assertEq(inputs.length, 4);
        assertEq(inputs[0], vm.parseJsonUint(fixture, ".epoch"));
        assertEq(inputs[1], vm.parseJsonUint(fixture, ".treeRoot"));
        assertEq(inputs[2], vm.parseJsonUint(fixture, ".totalCommitment"));
        assertEq(inputs[3], vm.parseJsonUint(fixture, ".leafCount"));
    }

    function test_AProofMadeWithTheIgnitionSRSVerifiesOnChain() public view {
        assertTrue(verifier.Verify{gas: VERIFY_GAS}(proof, inputs));
    }

    function test_TheProofIsRefusedForAnyOtherStatement() public view {
        for (uint256 i = 0; i < inputs.length; i++) {
            uint256[] memory other = new uint256[](inputs.length);
            for (uint256 j = 0; j < inputs.length; j++) {
                other[j] = inputs[j];
            }
            other[i] += 1;
            assertFalse(_accepts(proof, other), vm.toString(i));
        }
    }

    function test_TheTestVerifierIsNotTheProductionOne() public pure {
        assertNotEq(keccak256(type(PlonkVerifier).runtimeCode), keccak256(type(TestPlonkVerifier).runtimeCode));
    }

    function test_AGenerationTwoLedgerAnchorsTheProvenEpoch() public {
        uint64 number = uint64(inputs[0]);
        bytes32 treeRoot = bytes32(inputs[1]);
        bytes32 totalCommitment = bytes32(inputs[2]);
        uint64 leafCount = uint64(inputs[3]);
        address anchorer = makeAddr("anchorer");
        vm.warp(1_790_000_060);
        LiabilityLedgerV2 ledger = new LiabilityLedgerV2(
            makeAddr("admin"),
            3 days,
            anchorer,
            number,
            bytes32(uint256(0x07) << 248),
            1_789_996_400,
            PlonkVerifier(address(verifier)),
            8
        );

        bytes32 root = ledger.proveEpoch(number, treeRoot, totalCommitment, leafCount, proof);
        assertEq(root, ConfidentialTree.bindingRoot(treeRoot, totalCommitment));
        assertEq(ledger.capacity(), 256);

        vm.prank(anchorer);
        ledger.anchor(
            EpochHeader({
                number: number,
                previousDigest: bytes32(uint256(0x07) << 248),
                root: root,
                total: 0,
                leafCount: leafCount,
                asOf: 1_790_000_000
            })
        );
        assertEq(ledger.treeRootOf(number), treeRoot);

        vm.expectRevert(abi.encodeWithSelector(IConfidentialLedger.EpochClosed.selector, number, number + 1));
        ledger.proveEpoch(number, treeRoot, totalCommitment, leafCount, proof);
    }

    function _accepts(bytes memory candidate, uint256[] memory statement) private view returns (bool) {
        try verifier.Verify{gas: VERIFY_GAS}(candidate, statement) returns (bool verified) {
            return verified;
        } catch {
            return false;
        }
    }
}
