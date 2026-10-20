// SPDX-License-Identifier: MIT
pragma solidity 0.8.37;

import {Test} from "forge-std/src/Test.sol";

import {EpochHeader} from "../../src/EpochHeader.sol";
import {ConfidentialTree} from "../../src/confidential/ConfidentialTree.sol";
import {LiabilityLedgerV2} from "../../src/confidential/LiabilityLedgerV2.sol";
import {Poseidon2} from "../../src/confidential/Poseidon2.sol";
import {ScriptedVerifier} from "../support/ScriptedVerifier.sol";

contract ConfidentialLedgerHandler is Test {
    LiabilityLedgerV2 public immutable ledger;
    ScriptedVerifier public immutable verifier;
    address public immutable admin;
    address public immutable anchorer;

    EpochHeader[] private _headers;
    bytes32[] private _digests;
    bool public sealedByAdmin;
    uint256 public unexpected;

    constructor(LiabilityLedgerV2 ledger_, ScriptedVerifier verifier_, address admin_, address anchorer_) {
        ledger = ledger_;
        verifier = verifier_;
        admin = admin_;
        anchorer = anchorer_;
    }

    function proveAndAnchorNext(uint256 seed, uint64 rawLeafCount, uint256 rawAdvance, uint256 rawLag) external {
        vm.warp(vm.getBlockTimestamp() + bound(rawAdvance, 1, 7 days));
        uint64 number = ledger.nextEpoch();
        uint64 leafCount = uint64(bound(rawLeafCount, 0, ledger.capacity()));
        (bytes32 treeRoot, bytes32 totalCommitment) = _statement(seed);
        bool proven = _prove(number, treeRoot, totalCommitment, leafCount, seed, true);
        if (!proven) {
            return;
        }

        uint64 last = ledger.lastAsOf();
        uint64 asOf = uint64(vm.getBlockTimestamp() - bound(rawLag, 0, vm.getBlockTimestamp() - last - 1));
        EpochHeader memory header = EpochHeader({
            number: number,
            previousDigest: ledger.head(),
            root: ConfidentialTree.bindingRoot(treeRoot, totalCommitment),
            total: 0,
            leafCount: leafCount,
            asOf: asOf
        });

        vm.prank(anchorer);
        try ledger.anchor(header) returns (bytes32 digest) {
            _headers.push(header);
            _digests.push(digest);
        } catch {
            unexpected++;
        }
    }

    function proveAStrangersStatement(uint256 seed, uint64 rawAhead, address stranger) external {
        uint64 number = ledger.nextEpoch() + uint64(bound(rawAhead, 0, 3));
        (bytes32 treeRoot, bytes32 totalCommitment) = _statement(seed);
        bytes32 head = ledger.head();
        vm.startPrank(stranger);
        _prove(number, treeRoot, totalCommitment, uint64(bound(seed, 0, ledger.capacity())), seed, true);
        vm.stopPrank();
        if (ledger.head() != head) {
            unexpected++;
        }
    }

    function proveWithARefusedProof(uint256 seed) external {
        (bytes32 treeRoot, bytes32 totalCommitment) = _statement(seed);
        uint64 number = ledger.nextEpoch();
        try ledger.proveEpoch(number, treeRoot, totalCommitment, 1, abi.encode(seed, "refused")) {
            unexpected++;
        } catch {}
    }

    function anchorAnUnprovenEpoch(uint256 seed) external {
        (bytes32 treeRoot, bytes32 totalCommitment) = _statement(seed);
        EpochHeader memory header = EpochHeader({
            number: ledger.nextEpoch(),
            previousDigest: ledger.head(),
            root: ConfidentialTree.bindingRoot(treeRoot, totalCommitment),
            total: 0,
            leafCount: 1,
            asOf: uint64(vm.getBlockTimestamp())
        });
        if (ledger.isProven(header.number, header.root, header.leafCount)) {
            return;
        }

        vm.prank(anchorer);
        try ledger.anchor(header) {
            unexpected++;
        } catch {}
    }

    function anchorWithATotal(uint256 seed, uint64 total) external {
        if (total == 0) {
            return;
        }

        uint64 number = ledger.nextEpoch();
        (bytes32 treeRoot, bytes32 totalCommitment) = _statement(seed);
        _prove(number, treeRoot, totalCommitment, 1, seed, false);
        EpochHeader memory header = EpochHeader({
            number: number,
            previousDigest: ledger.head(),
            root: ConfidentialTree.bindingRoot(treeRoot, totalCommitment),
            total: total,
            leafCount: 1,
            asOf: uint64(vm.getBlockTimestamp())
        });

        vm.prank(anchorer);
        try ledger.anchor(header) {
            unexpected++;
        } catch {}
    }

    function replayAnAnchoredEpoch(uint256 rawWhich) external {
        if (_headers.length == 0) {
            return;
        }

        EpochHeader memory header = _headers[bound(rawWhich, 0, _headers.length - 1)];
        vm.prank(anchorer);
        try ledger.anchor(header) {
            unexpected++;
        } catch {}
    }

    function seal(uint256 chance) external {
        if (sealedByAdmin || chance % 64 != 0) {
            return;
        }

        bytes32 head = ledger.head();
        vm.prank(admin);
        try ledger.seal(head) {
            sealedByAdmin = true;
        } catch {
            unexpected++;
        }
    }

    function anchoredCount() external view returns (uint256) {
        return _headers.length;
    }

    function headerAt(uint256 index) external view returns (EpochHeader memory) {
        return _headers[index];
    }

    function digestAt(uint256 index) external view returns (bytes32) {
        return _digests[index];
    }

    function _prove(
        uint64 number,
        bytes32 treeRoot,
        bytes32 totalCommitment,
        uint64 leafCount,
        uint256 seed,
        bool counted
    ) private returns (bool) {
        bytes memory proof = abi.encode(seed, "proof");
        uint256[] memory inputs = new uint256[](4);
        inputs[0] = number;
        inputs[1] = uint256(treeRoot);
        inputs[2] = uint256(totalCommitment);
        inputs[3] = leafCount;
        verifier.accept(proof, inputs);

        bytes32 root = ConfidentialTree.bindingRoot(treeRoot, totalCommitment);
        bool already = ledger.isProven(number, root, leafCount);
        try ledger.proveEpoch(number, treeRoot, totalCommitment, leafCount, proof) {
            if (counted && (sealedByAdmin || already)) {
                unexpected++;
            }
            return true;
        } catch {
            if (counted && !sealedByAdmin && !already) {
                unexpected++;
            }
            return false;
        }
    }

    function _statement(uint256 seed) private pure returns (bytes32, bytes32) {
        return (
            bytes32(uint256(keccak256(abi.encode(seed, "tree"))) % Poseidon2.PRIME),
            bytes32(uint256(keccak256(abi.encode(seed, "total"))) % Poseidon2.PRIME)
        );
    }
}
