// SPDX-License-Identifier: MIT
pragma solidity 0.8.37;

import {EpochHeader} from "../../src/EpochHeader.sol";
import {AnchorCouncil} from "../../src/council/AnchorCouncil.sol";
import {IAnchorCouncil} from "../../src/council/IAnchorCouncil.sol";
import {CouncilFixture} from "../support/CouncilFixture.sol";

contract AnchorCouncilFuzzTest is CouncilFixture {
    uint256 private constant POOL = 8;

    function testFuzz_AnchorsExactlyWhenEnoughMembersApproveInOrder(
        uint256 rawSize,
        uint256 rawThreshold,
        uint256 memberMask,
        uint256 outsiderMask
    ) public {
        uint256 size = bound(rawSize, 2, POOL);
        uint8 threshold = uint8(bound(rawThreshold, 2, size));
        uint256[] memory keys = _keys("member", size);
        AnchorCouncil chosen = _council(keys, threshold);
        EpochHeader memory header = _next(rawSize, 4);

        uint256[] memory signing = _subset(keys, memberMask, _keys("outsider", 3), outsiderMask);
        bytes[] memory approvals = _ordered(chosen, header, signing);
        bool outsiders = outsiderMask % 8 != 0;

        if (approvals.length < threshold) {
            vm.expectRevert(
                abi.encodeWithSelector(IAnchorCouncil.NotEnoughApprovals.selector, approvals.length, threshold)
            );
        } else if (outsiders) {
            vm.expectPartialRevert(IAnchorCouncil.NotAMember.selector);
        }

        chosen.anchor(header, approvals);

        bool anchored = approvals.length >= threshold && !outsiders;
        assertEq(ledger.nextEpoch(), anchored ? 1 : 0);
    }

    function testFuzz_AnyDisorderIsRefused(uint256 rawSize, uint256 rawFirst, uint256 rawSecond) public {
        uint256 size = bound(rawSize, 2, POOL);
        uint256[] memory keys = _keys("member", size);
        AnchorCouncil chosen = _council(keys, 2);
        EpochHeader memory header = _next(rawSize, 3);
        bytes[] memory approvals = _ordered(chosen, header, keys);
        uint256 first = bound(rawFirst, 0, size - 1);
        uint256 second = bound(rawSecond, 0, size - 1);
        vm.assume(first != second);
        (approvals[first], approvals[second]) = (approvals[second], approvals[first]);

        vm.expectPartialRevert(IAnchorCouncil.ApprovalsOutOfOrder.selector);
        chosen.anchor(header, approvals);
    }

    function testFuzz_ARepeatedApproverIsRefused(uint256 rawSize, uint256 rawPosition) public {
        uint256 size = bound(rawSize, 2, POOL);
        uint256[] memory keys = _keys("member", size);
        AnchorCouncil chosen = _council(keys, 2);
        EpochHeader memory header = _next(rawSize, 3);
        bytes[] memory approvals = _ordered(chosen, header, keys);
        uint256 position = bound(rawPosition, 1, size - 1);
        approvals[position] = approvals[position - 1];

        vm.expectPartialRevert(IAnchorCouncil.ApprovalsOutOfOrder.selector);
        chosen.anchor(header, approvals);
        assertEq(ledger.nextEpoch(), 0);
    }

    function testFuzz_AnApprovalBindsEveryFieldOfTheHeader(uint256 field, uint256 noise) public {
        EpochHeader memory header = _next(noise, 6);
        bytes[] memory approvals = _approvals(header);
        EpochHeader memory changed = header;
        uint256 which = field % 6;
        uint64 delta = uint64(bound(noise, 1, 1000));
        if (which == 0) {
            changed.number += delta;
        } else if (which == 1) {
            changed.previousDigest = keccak256(abi.encode(noise));
        } else if (which == 2) {
            changed.root = keccak256(abi.encode(noise, "root"));
        } else if (which == 3) {
            changed.total += delta;
        } else if (which == 4) {
            changed.leafCount += delta;
        } else {
            changed.asOf -= delta;
        }

        vm.expectPartialRevert(IAnchorCouncil.NotAMember.selector);
        council.anchor(changed, approvals);
    }

    function testFuzz_AMembershipIsAcceptedExactlyWhenItIsWellFormed(
        address[] memory candidates,
        uint8 threshold
    ) public {
        vm.assume(candidates.length > 0);
        bool wellFormed = threshold >= 2 && candidates.length >= threshold && candidates.length <= 16
            && candidates[0] != address(0) && _ascending(candidates);

        if (!wellFormed) {
            vm.expectRevert();
        }

        AnchorCouncil built = new AnchorCouncil(admin, ADMIN_DELAY, ledger, candidates, threshold, MEMBERSHIP_DELAY);

        if (wellFormed) {
            assertEq(built.members(), candidates);
            assertEq(built.threshold(), threshold);
        }
    }

    function _keys(string memory label, uint256 count) private returns (uint256[] memory) {
        uint256[] memory keys = new uint256[](count);
        for (uint256 i = 0; i < count; i++) {
            (, keys[i]) = makeAddrAndKey(string.concat(label, vm.toString(i)));
        }

        return keys;
    }

    function _council(uint256[] memory keys, uint8 threshold) private returns (AnchorCouncil) {
        address[] memory addresses = new address[](keys.length);
        for (uint256 i = 0; i < keys.length; i++) {
            addresses[i] = vm.addr(keys[i]);
        }

        AnchorCouncil built =
            new AnchorCouncil(admin, ADMIN_DELAY, ledger, _sort(addresses), threshold, MEMBERSHIP_DELAY);
        _handOver(address(built));
        return built;
    }

    function _subset(
        uint256[] memory members,
        uint256 memberMask,
        uint256[] memory outsiders,
        uint256 outsiderMask
    ) private pure returns (uint256[] memory) {
        uint256 count = _countOf(members, memberMask) + _countOf(outsiders, outsiderMask);
        uint256[] memory chosen = new uint256[](count);
        uint256 next = 0;
        for (uint256 i = 0; i < members.length; i++) {
            if ((memberMask >> i) & 1 == 1) {
                chosen[next++] = members[i];
            }
        }

        for (uint256 i = 0; i < outsiders.length; i++) {
            if ((outsiderMask >> i) & 1 == 1) {
                chosen[next++] = outsiders[i];
            }
        }

        return chosen;
    }

    function _countOf(uint256[] memory keys, uint256 mask) private pure returns (uint256) {
        uint256 count = 0;
        for (uint256 i = 0; i < keys.length; i++) {
            if ((mask >> i) & 1 == 1) {
                count++;
            }
        }

        return count;
    }

    function _ordered(
        AnchorCouncil chosen,
        EpochHeader memory header,
        uint256[] memory keys
    ) private view returns (bytes[] memory) {
        uint256[] memory sorted = new uint256[](keys.length);
        for (uint256 i = 0; i < keys.length; i++) {
            sorted[i] = keys[i];
        }

        for (uint256 i = 0; i < sorted.length; i++) {
            for (uint256 j = i + 1; j < sorted.length; j++) {
                if (vm.addr(sorted[j]) < vm.addr(sorted[i])) {
                    (sorted[i], sorted[j]) = (sorted[j], sorted[i]);
                }
            }
        }

        bytes32 digest = chosen.approvalDigest(header.number, chosen.epochDigest(header));
        bytes[] memory approvals = new bytes[](sorted.length);
        for (uint256 i = 0; i < sorted.length; i++) {
            (uint8 v, bytes32 r, bytes32 s) = vm.sign(sorted[i], digest);
            approvals[i] = abi.encodePacked(r, s, v);
        }

        return approvals;
    }

    function _sort(address[] memory addresses) private pure returns (address[] memory) {
        for (uint256 i = 0; i < addresses.length; i++) {
            for (uint256 j = i + 1; j < addresses.length; j++) {
                if (addresses[j] < addresses[i]) {
                    (addresses[i], addresses[j]) = (addresses[j], addresses[i]);
                }
            }
        }

        return addresses;
    }

    function _ascending(address[] memory addresses) private pure returns (bool) {
        for (uint256 i = 1; i < addresses.length; i++) {
            if (addresses[i] <= addresses[i - 1]) {
                return false;
            }
        }

        return true;
    }
}
