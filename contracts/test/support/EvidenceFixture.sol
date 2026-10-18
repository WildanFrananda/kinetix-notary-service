// SPDX-License-Identifier: MIT
pragma solidity 0.8.37;

import {Test} from "forge-std/src/Test.sol";

import {BatchHeader} from "../../src/evidence/BatchHeader.sol";
import {EvidenceRegistry} from "../../src/evidence/EvidenceRegistry.sol";
import {StreamGenesis} from "../../src/evidence/StreamGenesis.sol";
import {EvidenceTreeBuilder} from "./EvidenceTreeBuilder.sol";

abstract contract EvidenceFixture is Test {
    uint48 internal constant ADMIN_DELAY = 3 days;
    uint256 internal constant START = 1_760_000_000;
    uint8 internal constant LIFECYCLE = 1;
    uint8 internal constant DELIVERY = 2;
    uint8 internal constant REVIEW = 3;

    address internal admin;
    address internal anchorer;
    address internal stranger;
    EvidenceRegistry internal registry;

    function setUp() public virtual {
        admin = makeAddr("admin");
        anchorer = makeAddr("evidence-anchorer");
        stranger = makeAddr("stranger");
        vm.warp(START);
        registry = new EvidenceRegistry(admin, ADMIN_DELAY, anchorer, _freshGenesis());
    }

    function _freshGenesis() internal pure returns (StreamGenesis[] memory) {
        StreamGenesis[] memory genesis = new StreamGenesis[](3);
        genesis[0] = _start(LIFECYCLE);
        genesis[1] = _start(DELIVERY);
        genesis[2] = _start(REVIEW);
        return genesis;
    }

    function _start(uint8 stream) internal pure returns (StreamGenesis memory) {
        return StreamGenesis({
            stream: stream,
            batch: 0,
            digest: bytes32(0),
            closedAt: 0,
            throughTransaction: 0,
            throughSequence: 0
        });
    }

    function _salts(uint256 seed, uint256 count) internal pure returns (bytes32[] memory) {
        bytes32[] memory salts = new bytes32[](count);
        for (uint256 i = 0; i < count; i++) {
            salts[i] = keccak256(abi.encode(seed, "salt", i));
        }

        return salts;
    }

    function _records(uint256 seed, uint256 count) internal pure returns (bytes32[] memory) {
        bytes32[] memory records = new bytes32[](count);
        for (uint256 i = 0; i < count; i++) {
            records[i] = keccak256(abi.encode(seed, "record", i));
        }

        return records;
    }

    function _tree(uint8 stream, uint256 seed, uint256 count) internal pure returns (bytes32[][] memory) {
        return EvidenceTreeBuilder.build(stream, _salts(seed, count), _records(seed, count));
    }

    function _next(uint8 stream, bytes32[][] memory levels, uint256 count) internal view returns (BatchHeader memory) {
        (,, uint64 nextBatch, uint64 throughTransaction, uint64 throughSequence) = _position(stream);
        return BatchHeader({
            stream: stream,
            number: nextBatch,
            previousDigest: registry.stateOf(stream).head,
            root: EvidenceTreeBuilder.root(levels),
            leafCount: uint64(count),
            throughTransaction: throughTransaction + 1,
            throughSequence: throughSequence + 1,
            closedAt: _now()
        });
    }

    function _position(uint8 stream) internal view returns (bytes32, uint64, uint64, uint64, uint64) {
        return (
            registry.stateOf(stream).head,
            registry.stateOf(stream).lastClosedAt,
            registry.stateOf(stream).nextBatch,
            registry.stateOf(stream).throughTransaction,
            registry.stateOf(stream).throughSequence
        );
    }

    function _anchor(BatchHeader memory header) internal returns (bytes32) {
        vm.prank(anchorer);
        return registry.anchor(header);
    }

    function _anchorNext(uint8 stream, uint256 seed, uint256 count) internal returns (bytes32[][] memory) {
        vm.warp(vm.getBlockTimestamp() + 1 days);
        bytes32[][] memory levels = _tree(stream, seed, count);
        _anchor(_next(stream, levels, count));
        return levels;
    }

    function _now() internal view returns (uint64) {
        return uint64(vm.getBlockTimestamp());
    }
}
