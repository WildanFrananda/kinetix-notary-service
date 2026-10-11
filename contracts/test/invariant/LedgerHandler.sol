// SPDX-License-Identifier: UNLICENSED
pragma solidity 0.8.37;

import {Test} from "forge-std/src/Test.sol";

import {EpochHeader} from "../../src/EpochHeader.sol";
import {LiabilityLedger} from "../../src/LiabilityLedger.sol";
import {SumTree} from "../../src/SumTree.sol";

contract LedgerHandler is Test {
    LiabilityLedger public immutable ledger;
    address public immutable admin;
    address public immutable anchorer;

    EpochHeader[] private _headers;
    bytes32[] private _digests;
    bool public sealedByAdmin;
    uint256 public unexpected;

    constructor(LiabilityLedger ledger_, address admin_, address anchorer_) {
        ledger = ledger_;
        admin = admin_;
        anchorer = anchorer_;
    }

    function anchorNext(uint256 rootSeed, uint64 rawLeafCount, uint64 total, uint256 rawAdvance, uint256 rawLag)
        external
    {
        vm.warp(vm.getBlockTimestamp() + bound(rawAdvance, 1, 7 days));
        uint64 leafCount = uint64(bound(rawLeafCount, 1, SumTree.MAX_LEAVES));
        uint64 last = ledger.lastAsOf();
        uint64 asOf = uint64(vm.getBlockTimestamp() - bound(rawLag, 0, vm.getBlockTimestamp() - last - 1));
        EpochHeader memory header = EpochHeader({
            number: ledger.nextEpoch(),
            previousDigest: ledger.head(),
            root: keccak256(abi.encode(rootSeed)),
            total: total,
            leafCount: leafCount,
            asOf: asOf
        });

        vm.prank(anchorer);
        try ledger.anchor(header) returns (bytes32 digest) {
            if (sealedByAdmin) {
                unexpected++;
            }
            _headers.push(header);
            _digests.push(digest);
        } catch {
            if (!sealedByAdmin) {
                unexpected++;
            }
        }
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

    function anchorOutOfSequence(uint64 rawSkip, uint256 rootSeed) external {
        uint64 skip = uint64(bound(rawSkip, 1, 1000));
        EpochHeader memory header = EpochHeader({
            number: ledger.nextEpoch() + skip,
            previousDigest: ledger.head(),
            root: keccak256(abi.encode(rootSeed)),
            total: 1,
            leafCount: 1,
            asOf: uint64(vm.getBlockTimestamp())
        });
        vm.prank(anchorer);
        try ledger.anchor(header) {
            unexpected++;
        } catch {}
    }

    function anchorOnAForgedHead(bytes32 forged, uint256 rootSeed) external {
        if (forged == ledger.head()) {
            return;
        }
        EpochHeader memory header = EpochHeader({
            number: ledger.nextEpoch(),
            previousDigest: forged,
            root: keccak256(abi.encode(rootSeed)),
            total: 1,
            leafCount: 1,
            asOf: uint64(vm.getBlockTimestamp())
        });
        vm.prank(anchorer);
        try ledger.anchor(header) {
            unexpected++;
        } catch {}
    }

    function anchorAsSomeoneElse(address caller, uint256 rootSeed) external {
        if (caller == anchorer) {
            return;
        }
        EpochHeader memory header = EpochHeader({
            number: ledger.nextEpoch(),
            previousDigest: ledger.head(),
            root: keccak256(abi.encode(rootSeed)),
            total: 1,
            leafCount: 1,
            asOf: uint64(vm.getBlockTimestamp())
        });
        vm.prank(caller);
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
}
