// SPDX-License-Identifier: MIT
pragma solidity 0.8.37;

import {Test} from "forge-std/src/Test.sol";

import {Poseidon2} from "../../src/confidential/Poseidon2.sol";
import {TestPlonkVerifier} from "../support/TestPlonkVerifier.sol";

contract PlonkVerifierTest is Test {
    uint256 private constant VERIFY_GAS = 1_000_000;

    string private fixture;
    TestPlonkVerifier private verifier;
    bytes private proof;
    uint256[] private inputs;

    function setUp() public {
        fixture = vm.readFile(string.concat(vm.projectRoot(), "/test/fixtures/plonk-v2.json"));
        verifier = new TestPlonkVerifier();
        proof = vm.parseJsonBytes(fixture, ".proof");
        inputs = vm.parseJsonUintArray(fixture, ".publicInputs");
    }

    function test_ThePublicInputsAreTheStatementInTheLedgersOrder() public view {
        assertEq(inputs.length, 4);
        assertEq(inputs[0], vm.parseJsonUint(fixture, ".epoch"));
        assertEq(inputs[1], vm.parseJsonUint(fixture, ".treeRoot"));
        assertEq(inputs[2], vm.parseJsonUint(fixture, ".totalCommitment"));
        assertEq(inputs[3], vm.parseJsonUint(fixture, ".leafCount"));
    }

    function test_AProofMadeInGoVerifiesOnChain() public view {
        assertTrue(verifier.Verify{gas: VERIFY_GAS}(proof, inputs));
    }

    function test_AProofIsRefusedForAnyOtherStatement() public view {
        for (uint256 i = 0; i < inputs.length; i++) {
            uint256[] memory other = _copy(inputs);
            other[i] += 1;
            assertFalse(_accepts(proof, other), vm.toString(i));
        }
    }

    function test_AnInputAliasedByTheModulusIsRefused() public view {
        for (uint256 i = 0; i < inputs.length; i++) {
            uint256[] memory aliased = _copy(inputs);
            aliased[i] += Poseidon2.PRIME;
            assertFalse(_accepts(proof, aliased), vm.toString(i));
        }
    }

    function test_AMissingOrAnExtraInputIsRefused() public view {
        uint256[] memory fewer = new uint256[](inputs.length - 1);
        uint256[] memory more = new uint256[](inputs.length + 1);
        for (uint256 i = 0; i < fewer.length; i++) {
            fewer[i] = inputs[i];
        }
        for (uint256 i = 0; i < inputs.length; i++) {
            more[i] = inputs[i];
        }

        assertFalse(_accepts(proof, fewer));
        assertFalse(_accepts(proof, more));
    }

    function test_EveryWordOfTheProofIsBinding() public view {
        assertEq(proof.length % 32, 0);
        for (uint256 offset = 0; offset < proof.length; offset += 32) {
            bytes memory tampered = bytes.concat(proof);
            tampered[offset + 31] ^= 0x01;
            assertFalse(_accepts(tampered, inputs), vm.toString(offset));
        }
    }

    function test_AShortOrALongProofIsRefused() public view {
        bytes memory short = new bytes(proof.length - 32);
        for (uint256 i = 0; i < short.length; i++) {
            short[i] = proof[i];
        }

        assertFalse(_accepts(short, inputs));
        assertFalse(_accepts(bytes.concat(proof, bytes32(0)), inputs));
    }

    function _accepts(bytes memory candidate, uint256[] memory statement) private view returns (bool) {
        try verifier.Verify{gas: VERIFY_GAS}(candidate, statement) returns (bool verified) {
            return verified;
        } catch {
            return false;
        }
    }

    function _copy(uint256[] memory values) private pure returns (uint256[] memory) {
        uint256[] memory out = new uint256[](values.length);
        for (uint256 i = 0; i < values.length; i++) {
            out[i] = values[i];
        }
        return out;
    }
}
