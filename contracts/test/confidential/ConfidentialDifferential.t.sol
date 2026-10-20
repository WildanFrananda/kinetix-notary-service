// SPDX-License-Identifier: MIT
pragma solidity 0.8.37;

import {Test} from "forge-std/src/Test.sol";

import {EpochHeader} from "../../src/EpochHeader.sol";
import {SumNode} from "../../src/SumNode.sol";
import {ConfidentialEpochChain} from "../../src/confidential/ConfidentialEpochChain.sol";
import {ConfidentialTree} from "../../src/confidential/ConfidentialTree.sol";
import {Poseidon2} from "../../src/confidential/Poseidon2.sol";
import {ConfidentialTreeBuilder} from "../support/ConfidentialTreeBuilder.sol";
import {ConfidentialTreeHarness} from "../support/ConfidentialTreeHarness.sol";

contract ConfidentialDifferentialTest is Test {
    string private fixture;
    ConfidentialTreeHarness private harness;

    function setUp() public {
        fixture = vm.readFile(string.concat(vm.projectRoot(), "/test/fixtures/scheme-v2.json"));
        harness = new ConfidentialTreeHarness();
    }

    function test_TheSchemeVersionAgrees() public view {
        assertEq(vm.parseJsonUint(fixture, ".scheme"), ConfidentialTree.SCHEME);
    }

    function test_EveryCompressionMatchesGnarkCrypto() public view {
        uint256 count = vm.parseJsonUint(fixture, ".compressionCount");
        assertGt(count, 20);
        for (uint256 i = 0; i < count; i++) {
            string memory key = string.concat(".compressions[", vm.toString(i), "]");
            uint256 left = vm.parseJsonUint(fixture, string.concat(key, ".left"));
            uint256 right = vm.parseJsonUint(fixture, string.concat(key, ".right"));
            assertEq(Poseidon2.compress(left, right), vm.parseJsonUint(fixture, string.concat(key, ".out")), key);
        }
    }

    function test_TheEmptyLeafMatchesGo() public view {
        assertEq(ConfidentialTree.empty(), vm.parseJsonUint(fixture, ".empty"));
    }

    function test_EveryLeafMatchesGo() public view {
        uint256 count = vm.parseJsonUint(fixture, ".leafCount");
        assertGt(count, 0);
        for (uint256 i = 0; i < count; i++) {
            string memory key = string.concat(".leaves[", vm.toString(i), "]");
            uint256 hash = ConfidentialTree.leaf(
                uint64(vm.parseJsonUint(fixture, string.concat(key, ".epoch"))),
                vm.parseJsonBytes32(fixture, string.concat(key, ".commitment")),
                uint64(vm.parseJsonUint(fixture, string.concat(key, ".amount")))
            );
            assertEq(hash, vm.parseJsonUint(fixture, string.concat(key, ".hash")), key);
        }
    }

    function test_EveryBindingRootMatchesGo() public view {
        uint256 count = vm.parseJsonUint(fixture, ".bindingCount");
        assertGt(count, 0);
        for (uint256 i = 0; i < count; i++) {
            string memory key = string.concat(".bindings[", vm.toString(i), "]");
            bytes32 root = ConfidentialTree.bindingRoot(
                vm.parseJsonBytes32(fixture, string.concat(key, ".treeRoot")),
                vm.parseJsonBytes32(fixture, string.concat(key, ".totalCommitment"))
            );
            assertEq(root, vm.parseJsonBytes32(fixture, string.concat(key, ".root")), key);
        }
    }

    function test_EveryGoTreeHasTheSameRootInSolidity() public view {
        uint256 count = vm.parseJsonUint(fixture, ".treeCount");
        assertGt(count, 0);
        for (uint256 t = 0; t < count; t++) {
            string memory key = _tree(t);
            (bytes32[] memory commitments, uint64[] memory amounts) = _leaves(key);
            uint256 root = ConfidentialTreeBuilder.root(
                _epochOf(key), vm.parseJsonUint(fixture, string.concat(key, ".depth")), commitments, amounts
            );
            assertEq(root, vm.parseJsonUint(fixture, string.concat(key, ".treeRoot")), key);
        }
    }

    function test_EveryProofGoIssuesVerifiesInSolidity() public view {
        uint256 count = vm.parseJsonUint(fixture, ".treeCount");
        uint256 verified = 0;
        for (uint256 t = 0; t < count; t++) {
            string memory key = _tree(t);
            (bytes32[] memory commitments, uint64[] memory amounts) = _leaves(key);
            for (uint256 j = 0; j < commitments.length; j++) {
                string memory proofKey = string.concat(key, ".proofs[", vm.toString(j), "]");
                (uint64 index, SumNode[] memory siblings) = _proof(proofKey);
                assertTrue(_verifies(key, commitments[j], amounts[j], index, siblings), proofKey);
                verified++;
            }
        }
        assertGt(verified, 50);
    }

    function test_EveryProofGoRejectsSolidityRejects() public view {
        uint256 count = vm.parseJsonUint(fixture, ".rejectionCount");
        assertGt(count, 8);
        for (uint256 r = 0; r < count; r++) {
            string memory key = string.concat(".rejections[", vm.toString(r), "]");
            string memory treeKey = _tree(vm.parseJsonUint(fixture, string.concat(key, ".tree")));
            (uint64 index, SumNode[] memory siblings) = _proof(string.concat(key, ".proof"));
            bool verifies = _verifies(
                treeKey,
                vm.parseJsonBytes32(fixture, string.concat(key, ".commitment")),
                uint64(vm.parseJsonUint(fixture, string.concat(key, ".amount"))),
                index,
                siblings
            );
            assertFalse(verifies, vm.parseJsonString(fixture, string.concat(key, ".reason")));
        }
    }

    function test_TheDigestChainMatchesGo() public view {
        uint256 count = vm.parseJsonUint(fixture, ".chainLength");
        assertGt(count, 0);
        bytes32 previous = bytes32(0);
        for (uint256 i = 0; i < count; i++) {
            string memory key = string.concat(".chain[", vm.toString(i), "]");
            EpochHeader memory header = EpochHeader({
                number: uint64(vm.parseJsonUint(fixture, string.concat(key, ".number"))),
                previousDigest: vm.parseJsonBytes32(fixture, string.concat(key, ".previousDigest")),
                root: vm.parseJsonBytes32(fixture, string.concat(key, ".root")),
                total: uint64(vm.parseJsonUint(fixture, string.concat(key, ".total"))),
                leafCount: uint64(vm.parseJsonUint(fixture, string.concat(key, ".leafCount"))),
                asOf: uint64(vm.parseJsonUint(fixture, string.concat(key, ".asOf")))
            });
            assertEq(header.previousDigest, previous, key);
            assertEq(header.total, 0, key);
            bytes32 digest = ConfidentialEpochChain.digest(header);
            assertEq(digest, vm.parseJsonBytes32(fixture, string.concat(key, ".digest")), key);
            previous = digest;
        }
    }

    function _verifies(
        string memory treeKey,
        bytes32 commitment,
        uint64 amount,
        uint64 index,
        SumNode[] memory siblings
    ) private view returns (bool) {
        return harness.verifies(
            _epochOf(treeKey),
            vm.parseJsonUint(fixture, string.concat(treeKey, ".treeRoot")),
            uint64(vm.parseJsonUint(fixture, string.concat(treeKey, ".leafCount"))),
            vm.parseJsonUint(fixture, string.concat(treeKey, ".depth")),
            commitment,
            amount,
            index,
            siblings
        );
    }

    function _tree(uint256 t) private pure returns (string memory) {
        return string.concat(".trees[", vm.toString(t), "]");
    }

    function _epochOf(string memory key) private view returns (uint64) {
        return uint64(vm.parseJsonUint(fixture, string.concat(key, ".epoch")));
    }

    function _leaves(string memory key) private view returns (bytes32[] memory, uint64[] memory) {
        uint256 leafCount = vm.parseJsonUint(fixture, string.concat(key, ".leafCount"));
        if (leafCount == 0) {
            return (new bytes32[](0), new uint64[](0));
        }
        bytes32[] memory commitments = vm.parseJsonBytes32Array(fixture, string.concat(key, ".commitments"));
        uint256[] memory raw = vm.parseJsonUintArray(fixture, string.concat(key, ".amounts"));
        uint64[] memory amounts = new uint64[](raw.length);
        for (uint256 i = 0; i < raw.length; i++) {
            amounts[i] = uint64(raw[i]);
        }
        assertEq(commitments.length, leafCount);
        return (commitments, amounts);
    }

    function _proof(string memory key) private view returns (uint64, SumNode[] memory) {
        uint64 index = uint64(vm.parseJsonUint(fixture, string.concat(key, ".index")));
        string memory hashesKey = string.concat(key, ".siblingHashes");
        if (!vm.keyExistsJson(fixture, string.concat(hashesKey, "[0]"))) {
            return (index, new SumNode[](0));
        }
        bytes32[] memory hashes = vm.parseJsonBytes32Array(fixture, hashesKey);
        uint256[] memory sums = vm.parseJsonUintArray(fixture, string.concat(key, ".siblingSums"));
        SumNode[] memory siblings = new SumNode[](hashes.length);
        for (uint256 i = 0; i < hashes.length; i++) {
            siblings[i] = SumNode({hash: hashes[i], sum: uint64(sums[i])});
        }
        return (index, siblings);
    }
}
