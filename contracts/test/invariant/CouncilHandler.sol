// SPDX-License-Identifier: MIT
pragma solidity 0.8.37;

import {Test} from "forge-std/src/Test.sol";

import {EpochHeader} from "../../src/EpochHeader.sol";
import {LiabilityLedger} from "../../src/LiabilityLedger.sol";
import {SumNode} from "../../src/SumNode.sol";
import {AnchorCouncil} from "../../src/council/AnchorCouncil.sol";
import {PendingMembership} from "../../src/council/PendingMembership.sol";
import {LeafSet} from "../support/LeafSet.sol";
import {TreeBuilder} from "../support/TreeBuilder.sol";

contract CouncilHandler is Test {
    LiabilityLedger public immutable ledger;
    AnchorCouncil public immutable council;
    address public immutable admin;

    uint256[] private _pool;
    mapping(uint64 epoch => uint256 approvers) private _approvers;
    mapping(uint64 epoch => uint8 threshold) private _thresholdAt;

    uint256 public anchored;
    uint256 public refused;
    uint256 public unexpected;

    constructor(AnchorCouncil council_, address admin_, uint256[] memory pool) {
        council = council_;
        ledger = LiabilityLedger(address(council_.ledger()));
        admin = admin_;
        _pool = pool;
    }

    function poolSize() external view returns (uint256) {
        return _pool.length;
    }

    function poolMember(uint256 index) external view returns (address) {
        return vm.addr(_pool[index]);
    }

    function approversOf(uint64 epoch) external view returns (uint256) {
        return _approvers[epoch];
    }

    function thresholdAt(uint64 epoch) external view returns (uint8) {
        return _thresholdAt[epoch];
    }

    function anchorNext(uint256 seed, uint256 rawCount, uint256 mask, uint256 rawSwap, uint256 rawAdvance) external {
        vm.warp(vm.getBlockTimestamp() + bound(rawAdvance, 1 minutes, 1 days));
        uint256 count = bound(rawCount, 0, 8);
        uint64 number = ledger.nextEpoch();
        SumNode[][] memory levels =
            TreeBuilder.build(number, LeafSet.commitments(seed, count), LeafSet.amounts(seed, count));
        SumNode memory root = TreeBuilder.root(levels);
        EpochHeader memory header = EpochHeader({
            number: number,
            previousDigest: ledger.head(),
            root: root.hash,
            total: root.sum,
            leafCount: uint64(count),
            asOf: uint64(vm.getBlockTimestamp() - 1)
        });

        uint256[] memory keys = _chosen(mask);
        bytes[] memory approvals = _signed(header, keys);
        bool tampered = keys.length > 1 && rawSwap % 4 < 2;
        if (tampered && rawSwap % 4 == 0) {
            (approvals[0], approvals[1]) = (approvals[1], approvals[0]);
        }

        if (tampered && rawSwap % 4 == 1) {
            approvals[1] = approvals[0];
        }

        bool shouldAnchor = !tampered && keys.length >= council.threshold() && _allMembers(keys);
        uint8 threshold = council.threshold();
        try council.anchor(header, approvals) {
            if (!shouldAnchor) {
                unexpected++;
            }

            _approvers[number] = keys.length;
            _thresholdAt[number] = threshold;
            anchored++;
        } catch {
            if (shouldAnchor) {
                unexpected++;
            }

            refused++;
        }
    }

    function schedule(uint256 mask, uint256 rawThreshold) external {
        address[] memory members = _addresses(_chosen(mask));
        uint8 threshold = uint8(bound(rawThreshold, 0, 4));
        bool pending = council.pendingMembership().activatesAt != 0;
        bool valid = threshold >= 2 && members.length >= threshold;
        bool same = threshold == council.threshold() && _equal(members, council.members());

        vm.prank(admin);
        try council.scheduleMembership(members, threshold) {
            if (pending || !valid || same) {
                unexpected++;
            }
        } catch {
            if (!pending && valid && !same) {
                unexpected++;
            }
        }
    }

    function activate(uint256 rawAdvance) external {
        vm.warp(vm.getBlockTimestamp() + bound(rawAdvance, 0, 8 days));
        PendingMembership memory pending = council.pendingMembership();
        bool due = pending.activatesAt != 0 && vm.getBlockTimestamp() >= pending.activatesAt;

        try council.activateMembership() {
            if (!due) {
                unexpected++;
            }

            if (!_equal(council.members(), pending.members) || council.threshold() != pending.threshold) {
                unexpected++;
            }
        } catch {
            if (due) {
                unexpected++;
            }
        }
    }

    function cancel() external {
        bool pending = council.pendingMembership().activatesAt != 0;

        vm.prank(admin);
        try council.cancelMembership() {
            if (!pending) {
                unexpected++;
            }
        } catch {
            if (pending) {
                unexpected++;
            }
        }
    }

    function _chosen(uint256 mask) private view returns (uint256[] memory) {
        uint256 count = 0;
        for (uint256 i = 0; i < _pool.length; i++) {
            if ((mask >> i) & 1 == 1) {
                count++;
            }
        }

        uint256[] memory keys = new uint256[](count);
        uint256 next = 0;
        for (uint256 i = 0; i < _pool.length; i++) {
            if ((mask >> i) & 1 == 1) {
                keys[next++] = _pool[i];
            }
        }

        for (uint256 i = 0; i < keys.length; i++) {
            for (uint256 j = i + 1; j < keys.length; j++) {
                if (vm.addr(keys[j]) < vm.addr(keys[i])) {
                    (keys[i], keys[j]) = (keys[j], keys[i]);
                }
            }
        }

        return keys;
    }

    function _signed(EpochHeader memory header, uint256[] memory keys) private view returns (bytes[] memory) {
        bytes32 digest = council.approvalDigest(header.number, council.epochDigest(header));
        bytes[] memory approvals = new bytes[](keys.length);
        for (uint256 i = 0; i < keys.length; i++) {
            (uint8 v, bytes32 r, bytes32 s) = vm.sign(keys[i], digest);
            approvals[i] = abi.encodePacked(r, s, v);
        }

        return approvals;
    }

    function _allMembers(uint256[] memory keys) private view returns (bool) {
        for (uint256 i = 0; i < keys.length; i++) {
            if (!council.isMember(vm.addr(keys[i]))) {
                return false;
            }
        }

        return true;
    }

    function _addresses(uint256[] memory keys) private pure returns (address[] memory) {
        address[] memory addresses = new address[](keys.length);
        for (uint256 i = 0; i < keys.length; i++) {
            addresses[i] = vm.addr(keys[i]);
        }

        return addresses;
    }

    function _equal(address[] memory a, address[] memory b) private pure returns (bool) {
        if (a.length != b.length) {
            return false;
        }

        for (uint256 i = 0; i < a.length; i++) {
            if (a[i] != b[i]) {
                return false;
            }
        }

        return true;
    }
}
