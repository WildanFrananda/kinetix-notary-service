// SPDX-License-Identifier: MIT
pragma solidity 0.8.37;

import {Test} from "forge-std/src/Test.sol";

import {BatchChain} from "../../src/evidence/BatchChain.sol";
import {BatchHeader} from "../../src/evidence/BatchHeader.sol";
import {EvidenceRegistry} from "../../src/evidence/EvidenceRegistry.sol";
import {EvidenceTree} from "../../src/evidence/EvidenceTree.sol";
import {StreamGenesis} from "../../src/evidence/StreamGenesis.sol";
import {EvidenceTreeBuilder} from "../support/EvidenceTreeBuilder.sol";

contract EvidenceDifferentialTest is Test {
    uint256 private constant START = 1_760_000_000;

    string private fixture;
    address private admin;
    address private anchorer;

    function setUp() public {
        fixture = vm.readFile(string.concat(vm.projectRoot(), "/test/fixtures/evidence-v1.json"));
        admin = makeAddr("admin");
        anchorer = makeAddr("evidence-anchorer");
        vm.warp(START + 30 days);
    }

    function test_TheSchemeVersionAndThePaddingAgree() public view {
        assertEq(vm.parseJsonUint(fixture, ".scheme"), EvidenceTree.SCHEME);
        assertEq(vm.parseJsonBytes32(fixture, ".empty"), EvidenceTree.empty());
    }

    function test_EveryGoLeafAndTreeIsTheSameInSolidity() public view {
        uint256 count = vm.parseJsonUint(fixture, ".treeCount");
        assertGt(count, 0);
        for (uint256 t = 0; t < count; t++) {
            string memory key = _tree(t);
            uint8 stream = _streamOf(key);
            bytes32[] memory salts = vm.parseJsonBytes32Array(fixture, string.concat(key, ".salts"));
            bytes32[] memory records = vm.parseJsonBytes32Array(fixture, string.concat(key, ".recordDigests"));
            bytes32[] memory leaves = vm.parseJsonBytes32Array(fixture, string.concat(key, ".leaves"));

            for (uint256 i = 0; i < leaves.length; i++) {
                assertEq(EvidenceTree.leaf(stream, salts[i], records[i]), leaves[i], key);
            }

            bytes32[][] memory levels = EvidenceTreeBuilder.build(stream, salts, records);
            assertEq(EvidenceTreeBuilder.root(levels), vm.parseJsonBytes32(fixture, string.concat(key, ".root")), key);
            assertEq(
                EvidenceTree.depthOf(uint64(leaves.length)),
                vm.parseJsonUint(fixture, string.concat(key, ".depth")),
                key
            );
        }
    }

    function test_EveryProofGoIssuesVerifiesOnTheRegistry() public {
        uint256 count = vm.parseJsonUint(fixture, ".treeCount");
        uint256 verified = 0;
        for (uint256 t = 0; t < count; t++) {
            string memory key = _tree(t);
            uint8 stream = _streamOf(key);
            EvidenceRegistry registry = _registryHolding(key);
            bytes32[] memory salts = vm.parseJsonBytes32Array(fixture, string.concat(key, ".salts"));
            bytes32[] memory records = vm.parseJsonBytes32Array(fixture, string.concat(key, ".recordDigests"));

            for (uint256 i = 0; i < salts.length; i++) {
                string memory proofKey = string.concat(key, ".proofs[", vm.toString(i), "]");
                (uint64 index, bytes32[] memory siblings) = _proof(proofKey);

                assertTrue(registry.verifyInclusion(stream, 0, salts[i], records[i], index, siblings), proofKey);
                verified++;
            }
        }

        assertGt(verified, 70);
    }

    function test_EveryProofGoRejectsTheRegistryRejects() public {
        uint256 count = vm.parseJsonUint(fixture, ".rejectionCount");
        assertGt(count, 0);
        for (uint256 r = 0; r < count; r++) {
            string memory key = string.concat(".rejections[", vm.toString(r), "]");
            EvidenceRegistry registry = _registryHolding(_tree(vm.parseJsonUint(fixture, string.concat(key, ".tree"))));
            (uint64 index, bytes32[] memory siblings) = _proof(string.concat(key, ".proof"));

            assertFalse(
                registry.verifyInclusion(
                    uint8(vm.parseJsonUint(fixture, string.concat(key, ".stream"))),
                    0,
                    vm.parseJsonBytes32(fixture, string.concat(key, ".salt")),
                    vm.parseJsonBytes32(fixture, string.concat(key, ".recordDigest")),
                    index,
                    siblings
                ),
                vm.parseJsonString(fixture, string.concat(key, ".reason"))
            );
        }
    }

    function test_EveryGoBatchChainIsTheChainTheRegistryKeeps() public {
        EvidenceRegistry registry = _fresh();
        uint256 length = vm.parseJsonUint(fixture, ".chainLength");
        assertGt(length, 0);
        for (uint256 i = 0; i < length; i++) {
            string memory key = string.concat(".chain[", vm.toString(i), "]");
            BatchHeader memory header = BatchHeader({
                stream: uint8(vm.parseJsonUint(fixture, string.concat(key, ".stream"))),
                number: uint64(vm.parseJsonUint(fixture, string.concat(key, ".number"))),
                previousDigest: vm.parseJsonBytes32(fixture, string.concat(key, ".previousDigest")),
                root: vm.parseJsonBytes32(fixture, string.concat(key, ".root")),
                leafCount: uint64(vm.parseJsonUint(fixture, string.concat(key, ".leafCount"))),
                throughTransaction: uint64(vm.parseJsonUint(fixture, string.concat(key, ".throughTransaction"))),
                throughSequence: uint64(vm.parseJsonUint(fixture, string.concat(key, ".throughSequence"))),
                closedAt: uint64(vm.parseJsonUint(fixture, string.concat(key, ".closedAt")))
            });
            bytes32 expected = vm.parseJsonBytes32(fixture, string.concat(key, ".digest"));

            assertEq(BatchChain.digest(header), expected, key);
            vm.prank(anchorer);
            assertEq(registry.anchor(header), expected, key);
            assertEq(registry.stateOf(header.stream).head, expected, key);
        }
    }

    function _registryHolding(string memory key) private returns (EvidenceRegistry) {
        EvidenceRegistry registry = _fresh();
        bytes32 root = vm.parseJsonBytes32(fixture, string.concat(key, ".root"));
        uint64 leafCount = uint64(vm.parseJsonUint(fixture, string.concat(key, ".leafCount")));
        for (uint8 stream = 1; stream <= 3; stream++) {
            BatchHeader memory header = BatchHeader({
                stream: stream,
                number: 0,
                previousDigest: bytes32(0),
                root: root,
                leafCount: leafCount,
                throughTransaction: 1,
                throughSequence: 1,
                closedAt: uint64(vm.getBlockTimestamp())
            });
            vm.prank(anchorer);
            registry.anchor(header);
        }

        return registry;
    }

    function _fresh() private returns (EvidenceRegistry) {
        StreamGenesis[] memory genesis = new StreamGenesis[](3);
        for (uint8 i = 0; i < 3; i++) {
            genesis[i] = StreamGenesis({
                stream: i + 1,
                batch: 0,
                digest: bytes32(0),
                closedAt: 0,
                throughTransaction: 0,
                throughSequence: 0
            });
        }

        return new EvidenceRegistry(admin, 3 days, anchorer, genesis);
    }

    function _proof(string memory key) private view returns (uint64, bytes32[] memory) {
        return (
            uint64(vm.parseJsonUint(fixture, string.concat(key, ".index"))),
            vm.parseJsonBytes32Array(fixture, string.concat(key, ".siblings"))
        );
    }

    function _tree(uint256 index) private pure returns (string memory) {
        return string.concat(".trees[", vm.toString(index), "]");
    }

    function _streamOf(string memory key) private view returns (uint8) {
        return uint8(vm.parseJsonUint(fixture, string.concat(key, ".stream")));
    }
}
