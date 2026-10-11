// SPDX-License-Identifier: UNLICENSED
pragma solidity 0.8.37;

import {Test} from "forge-std/src/Test.sol";

import {EpochChain} from "../src/EpochChain.sol";
import {EpochHeader} from "../src/EpochHeader.sol";
import {LiabilityLedger} from "../src/LiabilityLedger.sol";
import {SumNode} from "../src/SumNode.sol";
import {SumTree} from "../src/SumTree.sol";
import {TreeBuilder} from "./support/TreeBuilder.sol";

contract DifferentialTest is Test {
    string private fixture;
    address private admin;
    address private anchorer;

    function setUp() public {
        fixture = vm.readFile(string.concat(vm.projectRoot(), "/test/fixtures/scheme-v1.json"));
        admin = makeAddr("admin");
        anchorer = makeAddr("anchorer");
        vm.warp(1_760_000_000);
    }

    function test_TheSchemeVersionAgrees() public view {
        assertEq(vm.parseJsonUint(fixture, ".scheme"), SumTree.SCHEME);
    }

    function test_EveryGoTreeHasTheSameRootInSolidity() public view {
        uint256 count = vm.parseJsonUint(fixture, ".treeCount");
        assertGt(count, 0);
        for (uint256 t = 0; t < count; t++) {
            string memory key = _tree(t);
            (bytes32[] memory commitments, uint64[] memory amounts) = _leaves(key);
            SumNode memory root = TreeBuilder.root(TreeBuilder.build(_epochOf(key), commitments, amounts));

            assertEq(root.hash, vm.parseJsonBytes32(fixture, string.concat(key, ".root")), key);
            assertEq(root.sum, vm.parseJsonUint(fixture, string.concat(key, ".total")), key);
            assertEq(
                SumTree.depthOf(uint64(commitments.length)),
                vm.parseJsonUint(fixture, string.concat(key, ".depth")),
                key
            );
        }
    }

    function test_EveryProofGoIssuesVerifiesOnTheLedger() public {
        uint256 count = vm.parseJsonUint(fixture, ".treeCount");
        uint256 verified = 0;
        for (uint256 t = 0; t < count; t++) {
            string memory key = _tree(t);
            (bytes32[] memory commitments, uint64[] memory amounts) = _leaves(key);
            if (commitments.length == 0) {
                continue;
            }
            LiabilityLedger ledger = _anchoredLedger(key);
            for (uint256 j = 0; j < commitments.length; j++) {
                string memory proofKey = string.concat(key, ".proofs[", vm.toString(j), "]");
                (uint64 index, SumNode[] memory siblings) = _proof(proofKey);

                assertTrue(ledger.verifyInclusion(_epochOf(key), commitments[j], amounts[j], index, siblings), proofKey);
                verified++;
            }
        }
        assertGt(verified, 100);
    }

    function test_EveryProofGoRejectsTheLedgerRejects() public {
        uint256 count = vm.parseJsonUint(fixture, ".rejectionCount");
        assertGt(count, 0);
        for (uint256 r = 0; r < count; r++) {
            string memory key = string.concat(".rejections[", vm.toString(r), "]");
            string memory treeKey = _tree(vm.parseJsonUint(fixture, string.concat(key, ".tree")));
            LiabilityLedger ledger = _anchoredLedger(treeKey);
            (uint64 index, SumNode[] memory siblings) = _proof(string.concat(key, ".proof"));

            assertFalse(
                ledger.verifyInclusion(
                    _epochOf(treeKey),
                    vm.parseJsonBytes32(fixture, string.concat(key, ".commitment")),
                    uint64(vm.parseJsonUint(fixture, string.concat(key, ".amount"))),
                    index,
                    siblings
                ),
                vm.parseJsonString(fixture, string.concat(key, ".reason"))
            );
        }
    }

    function test_TheEpochChainDigestsAgree() public view {
        uint256 length = vm.parseJsonUint(fixture, ".chainLength");
        bytes32 previous = bytes32(0);
        for (uint256 i = 0; i < length; i++) {
            (EpochHeader memory header, bytes32 digest) = _link(i);

            assertEq(header.previousDigest, previous);
            assertEq(EpochChain.digest(header), digest);
            previous = digest;
        }
    }

    function test_TheGoChainAnchorsOnTheLedgerWithTheSameDigests() public {
        LiabilityLedger ledger = new LiabilityLedger(admin, 3 days, anchorer, 0, bytes32(0), 0);
        uint256 length = vm.parseJsonUint(fixture, ".chainLength");
        for (uint256 i = 0; i < length; i++) {
            (EpochHeader memory header, bytes32 digest) = _link(i);
            vm.warp(header.asOf + 1);

            vm.prank(anchorer);
            assertEq(ledger.anchor(header), digest);
        }
        assertEq(ledger.nextEpoch(), length);
    }

    function _tree(uint256 index) private pure returns (string memory) {
        return string.concat(".trees[", vm.toString(index), "]");
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
        assertEq(amounts.length, leafCount);
        return (commitments, amounts);
    }

    function _proof(string memory key) private view returns (uint64, SumNode[] memory) {
        uint64 index = uint64(vm.parseJsonUint(fixture, string.concat(key, ".index")));
        string memory hashesKey = string.concat(key, ".siblingHashes");
        if (!_hasElements(hashesKey)) {
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

    function _hasElements(string memory arrayKey) private view returns (bool) {
        return vm.keyExistsJson(fixture, string.concat(arrayKey, "[0]"));
    }

    function _link(uint256 i) private view returns (EpochHeader memory, bytes32) {
        string memory key = string.concat(".chain[", vm.toString(i), "]");
        EpochHeader memory header = EpochHeader({
            number: uint64(vm.parseJsonUint(fixture, string.concat(key, ".number"))),
            previousDigest: vm.parseJsonBytes32(fixture, string.concat(key, ".previousDigest")),
            root: vm.parseJsonBytes32(fixture, string.concat(key, ".root")),
            total: uint64(vm.parseJsonUint(fixture, string.concat(key, ".total"))),
            leafCount: uint64(vm.parseJsonUint(fixture, string.concat(key, ".leafCount"))),
            asOf: uint64(vm.parseJsonUint(fixture, string.concat(key, ".asOf")))
        });
        return (header, vm.parseJsonBytes32(fixture, string.concat(key, ".digest")));
    }

    function _anchoredLedger(string memory key) private returns (LiabilityLedger) {
        uint64 epoch = _epochOf(key);
        bytes32 genesis = epoch == 0 ? bytes32(0) : keccak256("differential genesis");
        LiabilityLedger ledger = new LiabilityLedger(admin, 3 days, anchorer, epoch, genesis, 0);
        EpochHeader memory header = EpochHeader({
            number: epoch,
            previousDigest: genesis,
            root: vm.parseJsonBytes32(fixture, string.concat(key, ".root")),
            total: uint64(vm.parseJsonUint(fixture, string.concat(key, ".total"))),
            leafCount: uint64(vm.parseJsonUint(fixture, string.concat(key, ".leafCount"))),
            asOf: uint64(vm.getBlockTimestamp() - 1)
        });
        vm.prank(anchorer);
        ledger.anchor(header);
        return ledger;
    }
}
