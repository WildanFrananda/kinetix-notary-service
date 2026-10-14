// SPDX-License-Identifier: MIT
pragma solidity 0.8.37;

import {IAccessControl} from "@openzeppelin-contracts/access/IAccessControl.sol";
import {IERC165} from "@openzeppelin-contracts/utils/introspection/IERC165.sol";

import {EpochHeader} from "../../src/EpochHeader.sol";
import {ILiabilityLedger} from "../../src/ILiabilityLedger.sol";
import {LiabilityLedger} from "../../src/LiabilityLedger.sol";
import {AnchorCouncil} from "../../src/council/AnchorCouncil.sol";
import {IAnchorCouncil} from "../../src/council/IAnchorCouncil.sol";
import {PendingMembership} from "../../src/council/PendingMembership.sol";
import {CouncilFixture} from "../support/CouncilFixture.sol";

contract AnchorCouncilTest is CouncilFixture {
    function test_StartsWithItsMembersAndThreshold() public view {
        address[] memory current = council.members();

        assertEq(current.length, 2);
        assertTrue(current[0] < current[1]);
        assertTrue(council.isMember(notary));
        assertTrue(council.isMember(watch));
        assertFalse(council.isMember(stranger));
        assertFalse(council.isMember(address(0)));
        assertEq(council.threshold(), 2);
        assertEq(address(council.ledger()), address(ledger));
        assertEq(council.membershipChangeDelay(), MEMBERSHIP_DELAY);
        assertEq(council.defaultAdmin(), admin);
        assertEq(council.defaultAdminDelay(), ADMIN_DELAY);
        assertEq(council.pendingMembership().activatesAt, 0);
    }

    function test_RefusesAnythingThatIsNotALedger() public {
        address[] memory members = _pair(notary, watch);

        vm.expectRevert(abi.encodeWithSelector(IAnchorCouncil.NotALedger.selector, stranger));
        new AnchorCouncil(admin, ADMIN_DELAY, ILiabilityLedger(stranger), members, 2, MEMBERSHIP_DELAY);

        vm.expectRevert(abi.encodeWithSelector(IAnchorCouncil.NotALedger.selector, address(council)));
        new AnchorCouncil(admin, ADMIN_DELAY, ILiabilityLedger(address(council)), members, 2, MEMBERSHIP_DELAY);
    }

    function test_RefusesAZeroDelay() public {
        address[] memory members = _pair(notary, watch);

        vm.expectRevert(IAnchorCouncil.ZeroDelay.selector);
        new AnchorCouncil(admin, ADMIN_DELAY, ledger, members, 2, 0);
    }

    function test_RefusesAThresholdBelowTwo() public {
        address[] memory members = _pair(notary, watch);

        vm.expectRevert(abi.encodeWithSelector(IAnchorCouncil.ThresholdTooLow.selector, 1, 2));
        new AnchorCouncil(admin, ADMIN_DELAY, ledger, members, 1, MEMBERSHIP_DELAY);

        vm.expectRevert(abi.encodeWithSelector(IAnchorCouncil.ThresholdTooLow.selector, 0, 2));
        new AnchorCouncil(admin, ADMIN_DELAY, ledger, members, 0, MEMBERSHIP_DELAY);
    }

    function test_RefusesFewerMembersThanTheThreshold() public {
        address[] memory members = _pair(notary, watch);

        vm.expectRevert(abi.encodeWithSelector(IAnchorCouncil.TooFewMembers.selector, 2, 3));
        new AnchorCouncil(admin, ADMIN_DELAY, ledger, members, 3, MEMBERSHIP_DELAY);
    }

    function test_RefusesMoreThanSixteenMembers() public {
        address[] memory members = _ascending(17);

        vm.expectRevert(abi.encodeWithSelector(IAnchorCouncil.TooManyMembers.selector, 17, 16));
        new AnchorCouncil(admin, ADMIN_DELAY, ledger, members, 2, MEMBERSHIP_DELAY);

        AnchorCouncil largest = new AnchorCouncil(admin, ADMIN_DELAY, ledger, _ascending(16), 16, MEMBERSHIP_DELAY);
        assertEq(largest.members().length, 16);
    }

    function test_RefusesAZeroMember() public {
        address[] memory members = _pair(address(0), watch);

        vm.expectRevert(IAnchorCouncil.ZeroMember.selector);
        new AnchorCouncil(admin, ADMIN_DELAY, ledger, members, 2, MEMBERSHIP_DELAY);
    }

    function test_RefusesMembersOutOfOrderOrRepeated() public {
        address[] memory descending = _pair(notary, watch);
        (descending[0], descending[1]) = (descending[1], descending[0]);

        vm.expectRevert(abi.encodeWithSelector(IAnchorCouncil.MembersNotAscending.selector, 1));
        new AnchorCouncil(admin, ADMIN_DELAY, ledger, descending, 2, MEMBERSHIP_DELAY);

        address[] memory repeated = new address[](2);
        repeated[0] = watch;
        repeated[1] = watch;

        vm.expectRevert(abi.encodeWithSelector(IAnchorCouncil.MembersNotAscending.selector, 1));
        new AnchorCouncil(admin, ADMIN_DELAY, ledger, repeated, 2, MEMBERSHIP_DELAY);
    }

    function test_AnnouncesItsInterfaces() public view {
        assertTrue(council.supportsInterface(type(IAnchorCouncil).interfaceId));
        assertTrue(council.supportsInterface(type(IAccessControl).interfaceId));
        assertTrue(council.supportsInterface(type(IERC165).interfaceId));
        assertFalse(council.supportsInterface(type(ILiabilityLedger).interfaceId));
        assertFalse(council.supportsInterface(0xffffffff));
    }

    function test_TheDomainNamesThisCouncilOnThisChain() public view {
        (, string memory name, string memory version, uint256 chainId, address verifyingContract,,) =
            council.eip712Domain();

        assertEq(name, "Kinetix Anchor Council");
        assertEq(version, "1");
        assertEq(chainId, block.chainid);
        assertEq(verifyingContract, address(council));
    }

    function test_TheApprovalIsEIP712OverTheLedgerTheNumberAndTheEpochDigest() public view {
        bytes32 digest = keccak256("epoch digest");
        bytes32 domain = keccak256(
            abi.encode(
                keccak256("EIP712Domain(string name,string version,uint256 chainId,address verifyingContract)"),
                keccak256("Kinetix Anchor Council"),
                keccak256("1"),
                block.chainid,
                address(council)
            )
        );
        bytes32 structHash = keccak256(
            abi.encode(
                keccak256("EpochApproval(address ledger,uint64 number,bytes32 digest)"), address(ledger), 7, digest
            )
        );

        assertEq(council.approvalDigest(7, digest), keccak256(abi.encodePacked("\x19\x01", domain, structHash)));
    }

    function test_TheEpochDigestIsTheOneTheLedgerAnchors() public {
        EpochHeader memory header = _next(1, 5);
        bytes32 expected = council.epochDigest(header);
        bytes[] memory approvals = _approvals(header);

        assertEq(council.anchor(header, approvals), expected);
        assertEq(ledger.head(), expected);
    }

    function test_AnchorsTheNextEpochWithEveryMembersApproval() public {
        EpochHeader memory header = _next(1, 5);
        bytes32 digest = council.epochDigest(header);
        bytes[] memory approvals = _approvals(header);

        vm.expectEmit(address(council));
        emit IAnchorCouncil.EpochApproved(0, digest, council.members());
        vm.prank(stranger);
        council.anchor(header, approvals);

        assertEq(ledger.nextEpoch(), 1);
        assertEq(ledger.epochAt(0).digest, digest);
    }

    function test_AnchorsAnEmptyEpoch() public {
        EpochHeader memory header = _next(1, 0);
        bytes[] memory approvals = _approvals(header);

        council.anchor(header, approvals);

        assertEq(ledger.nextEpoch(), 1);
    }

    function test_AnchorsEpochAfterEpoch() public {
        for (uint256 i = 0; i < 4; i++) {
            _anchorNext(i, i + 1);
        }

        assertEq(ledger.nextEpoch(), 4);
    }

    function test_RefusesFewerApprovalsThanTheThreshold() public {
        EpochHeader memory header = _next(1, 5);
        bytes[] memory one = new bytes[](1);
        one[0] = _approve(notaryKey, header);

        vm.expectRevert(abi.encodeWithSelector(IAnchorCouncil.NotEnoughApprovals.selector, 1, 2));
        council.anchor(header, one);

        bytes[] memory none = new bytes[](0);

        vm.expectRevert(abi.encodeWithSelector(IAnchorCouncil.NotEnoughApprovals.selector, 0, 2));
        council.anchor(header, none);
        assertEq(ledger.nextEpoch(), 0);
    }

    function test_RefusesAnApprovalFromSomeoneWhoIsNotAMember() public {
        EpochHeader memory header = _next(1, 5);
        (address outsider, uint256 outsiderKey) = makeAddrAndKey("outsider");
        bytes[] memory approvals = new bytes[](2);
        bytes memory byOutsider = _approve(outsiderKey, header);
        bytes memory byWatch = _approve(watchKey, header);
        (approvals[0], approvals[1]) = outsider < watch ? (byOutsider, byWatch) : (byWatch, byOutsider);

        vm.expectRevert(abi.encodeWithSelector(IAnchorCouncil.NotAMember.selector, outsider));
        council.anchor(header, approvals);
    }

    function test_RefusesTheSameMemberTwice() public {
        EpochHeader memory header = _next(1, 5);
        bytes[] memory approvals = new bytes[](2);
        approvals[0] = _approve(watchKey, header);
        approvals[1] = approvals[0];

        vm.expectRevert(abi.encodeWithSelector(IAnchorCouncil.ApprovalsOutOfOrder.selector, watch, watch));
        council.anchor(header, approvals);
    }

    function test_RefusesApprovalsInDescendingOrder() public {
        EpochHeader memory header = _next(1, 5);
        bytes[] memory approvals = _approvals(header);
        (approvals[0], approvals[1]) = (approvals[1], approvals[0]);
        (address low, address high) = notary < watch ? (notary, watch) : (watch, notary);

        vm.expectRevert(abi.encodeWithSelector(IAnchorCouncil.ApprovalsOutOfOrder.selector, high, low));
        council.anchor(header, approvals);
    }

    function test_RefusesAMalformedSignature() public {
        EpochHeader memory header = _next(1, 5);
        bytes[] memory approvals = _approvals(header);
        approvals[0] = hex"deadbeef";

        vm.expectRevert(abi.encodeWithSelector(IAnchorCouncil.NotAMember.selector, address(0)));
        council.anchor(header, approvals);
    }

    function test_RefusesAHighSSignature() public {
        EpochHeader memory header = _next(1, 5);
        bytes[] memory approvals = _approvals(header);
        bytes memory signature = approvals[0];
        bytes32 r;
        bytes32 s;
        uint8 v;
        assembly {
            r := mload(add(signature, 32))
            s := mload(add(signature, 64))
            v := byte(0, mload(add(signature, 96)))
        }
        uint256 order = 0xFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFEBAAEDCE6AF48A03BBFD25E8CD0364141;
        approvals[0] = abi.encodePacked(r, bytes32(order - uint256(s)), v == 27 ? uint8(28) : uint8(27));

        vm.expectRevert(abi.encodeWithSelector(IAnchorCouncil.NotAMember.selector, address(0)));
        council.anchor(header, approvals);
    }

    function test_AnApprovalCoversExactlyOneHeader() public {
        EpochHeader memory header = _next(1, 5);
        bytes[] memory approvals = _approvals(header);
        EpochHeader memory other = _next(2, 5);
        address recovered = council.recoverApprover(other.number, council.epochDigest(other), approvals[0]);

        assertFalse(council.isMember(recovered));
        vm.expectRevert(abi.encodeWithSelector(IAnchorCouncil.NotAMember.selector, recovered));
        council.anchor(other, approvals);
    }

    function test_AnApprovalForAnotherCouncilCountsForNothing() public {
        AnchorCouncil twin = new AnchorCouncil(admin, ADMIN_DELAY, ledger, _pair(notary, watch), 2, MEMBERSHIP_DELAY);
        _handOver(address(twin));
        EpochHeader memory header = _next(1, 5);
        bytes[] memory approvals = _approvals(header);
        address recovered = twin.recoverApprover(header.number, twin.epochDigest(header), approvals[0]);

        assertFalse(twin.isMember(recovered));
        vm.expectRevert(abi.encodeWithSelector(IAnchorCouncil.NotAMember.selector, recovered));
        twin.anchor(header, approvals);
    }

    function test_RecoversEachMembersApproval() public view {
        EpochHeader memory header = _next(1, 5);
        bytes32 digest = council.epochDigest(header);

        assertEq(council.recoverApprover(0, digest, _approve(notaryKey, header)), notary);
        assertEq(council.recoverApprover(0, digest, _approve(watchKey, header)), watch);
        assertEq(council.recoverApprover(0, digest, hex"00"), address(0));
    }

    function test_TheLedgersRefusalReachesTheCallerUnchanged() public {
        EpochHeader memory header = _anchorNext(1, 5);
        bytes32 digest = ledger.epochAt(0).digest;
        bytes[] memory again = _approvals(header);

        vm.expectRevert(abi.encodeWithSelector(ILiabilityLedger.EpochAlreadyAnchored.selector, 0, digest));
        council.anchor(header, again);

        EpochHeader memory skipped = _next(2, 5);
        skipped.number = 5;
        bytes[] memory approvals = _approvals(skipped);

        vm.expectRevert(abi.encodeWithSelector(ILiabilityLedger.EpochOutOfSequence.selector, 1, 5));
        council.anchor(skipped, approvals);

        EpochHeader memory future = _next(3, 5);
        future.asOf = _now() + 1 hours;
        bytes[] memory early = _approvals(future);

        vm.expectRevert(abi.encodeWithSelector(ILiabilityLedger.SnapshotInFuture.selector, future.asOf, _now()));
        council.anchor(future, early);
    }

    function test_AnchorsNothingWithoutTheAnchorerRole() public {
        LiabilityLedger other = new LiabilityLedger(admin, ADMIN_DELAY, anchorer, 0, bytes32(0), 0);
        AnchorCouncil powerless =
            new AnchorCouncil(admin, ADMIN_DELAY, other, _pair(notary, watch), 2, MEMBERSHIP_DELAY);
        EpochHeader memory header = _header(0, bytes32(0), _tree(0, 1, 5), 5, _now() - 1);
        bytes32 digest = powerless.epochDigest(header);
        bytes[] memory approvals = new bytes[](2);
        bytes memory byNotary = _signed(notaryKey, powerless.approvalDigest(0, digest));
        bytes memory byWatch = _signed(watchKey, powerless.approvalDigest(0, digest));
        (approvals[0], approvals[1]) = notary < watch ? (byNotary, byWatch) : (byWatch, byNotary);
        bytes32 role = other.ANCHORER_ROLE();

        vm.expectRevert(
            abi.encodeWithSelector(IAccessControl.AccessControlUnauthorizedAccount.selector, address(powerless), role)
        );
        powerless.anchor(header, approvals);
    }

    function test_OnlyTheAdminSchedulesAMembership() public {
        address[] memory next = _pair(notary, stranger);
        bytes32 role = council.DEFAULT_ADMIN_ROLE();

        vm.expectRevert(
            abi.encodeWithSelector(IAccessControl.AccessControlUnauthorizedAccount.selector, stranger, role)
        );
        vm.prank(stranger);
        council.scheduleMembership(next, 2);
    }

    function test_AScheduledMembershipWaitsForItsDelay() public {
        (address newcomer,) = makeAddrAndKey("newcomer");
        address[] memory next = _pair(notary, newcomer);
        uint64 activatesAt = _now() + MEMBERSHIP_DELAY;

        vm.expectEmit(address(council));
        emit IAnchorCouncil.MembershipScheduled(next, 2, activatesAt);
        vm.prank(admin);
        council.scheduleMembership(next, 2);

        PendingMembership memory pending = council.pendingMembership();
        assertEq(pending.members, next);
        assertEq(pending.threshold, 2);
        assertEq(pending.activatesAt, activatesAt);
        assertTrue(council.isMember(watch));

        vm.warp(activatesAt - 1);
        vm.expectRevert(
            abi.encodeWithSelector(IAnchorCouncil.MembershipChangeTooEarly.selector, activatesAt, activatesAt - 1)
        );
        council.activateMembership();
    }

    function test_AnyoneActivatesADueMembershipAndTheOldMemberStopsCounting() public {
        (address newcomer, uint256 newcomerKey) = makeAddrAndKey("newcomer");
        address[] memory next = _pair(notary, newcomer);
        vm.prank(admin);
        council.scheduleMembership(next, 2);
        EpochHeader memory header = _next(1, 5);
        bytes[] memory before = _approvals(header);

        vm.warp(_now() + MEMBERSHIP_DELAY);
        vm.expectEmit(address(council));
        emit IAnchorCouncil.MembershipActivated(next, 2);
        vm.prank(stranger);
        council.activateMembership();

        assertEq(council.members(), next);
        assertFalse(council.isMember(watch));
        assertTrue(council.isMember(newcomer));
        assertEq(council.pendingMembership().activatesAt, 0);
        assertEq(council.pendingMembership().members.length, 0);

        vm.expectRevert(abi.encodeWithSelector(IAnchorCouncil.NotAMember.selector, watch));
        council.anchor(header, before);

        bytes[] memory approvals = new bytes[](2);
        bytes memory byNotary = _approve(notaryKey, header);
        bytes memory byNewcomer = _approve(newcomerKey, header);
        (approvals[0], approvals[1]) = notary < newcomer ? (byNotary, byNewcomer) : (byNewcomer, byNotary);
        council.anchor(header, approvals);
        assertEq(ledger.nextEpoch(), 1);
    }

    function test_AMembershipCanRaiseTheThreshold() public {
        (address third, uint256 thirdKey) = makeAddrAndKey("third");
        address[] memory three = _ascendingOf(notary, watch, third);
        vm.prank(admin);
        council.scheduleMembership(three, 3);
        vm.warp(_now() + MEMBERSHIP_DELAY);
        council.activateMembership();
        EpochHeader memory header = _next(1, 5);
        bytes[] memory two = _approvals(header);

        vm.expectRevert(abi.encodeWithSelector(IAnchorCouncil.NotEnoughApprovals.selector, 2, 3));
        council.anchor(header, two);

        uint256[] memory keys = new uint256[](3);
        keys[0] = notaryKey;
        keys[1] = watchKey;
        keys[2] = thirdKey;
        council.anchor(header, _sortedApprovals(header, keys));
        assertEq(ledger.nextEpoch(), 1);
    }

    function test_TheAdminCancelsAPendingMembership() public {
        address[] memory next = _pair(notary, stranger);
        vm.startPrank(admin);
        council.scheduleMembership(next, 2);

        vm.expectEmit(address(council));
        emit IAnchorCouncil.MembershipCancelled(next, 2);
        council.cancelMembership();
        vm.stopPrank();

        assertEq(council.pendingMembership().activatesAt, 0);
        vm.warp(_now() + MEMBERSHIP_DELAY);
        vm.expectRevert(IAnchorCouncil.NoMembershipChangePending.selector);
        council.activateMembership();
        assertTrue(council.isMember(watch));
    }

    function test_OnlyTheAdminCancels() public {
        address[] memory next = _pair(notary, stranger);
        vm.prank(admin);
        council.scheduleMembership(next, 2);
        bytes32 role = council.DEFAULT_ADMIN_ROLE();

        vm.expectRevert(
            abi.encodeWithSelector(IAccessControl.AccessControlUnauthorizedAccount.selector, stranger, role)
        );
        vm.prank(stranger);
        council.cancelMembership();
    }

    function test_RefusesASecondScheduleWhileOneIsPending() public {
        address[] memory next = _pair(notary, stranger);
        vm.startPrank(admin);
        council.scheduleMembership(next, 2);
        uint64 activatesAt = council.pendingMembership().activatesAt;

        vm.expectRevert(abi.encodeWithSelector(IAnchorCouncil.MembershipChangePending.selector, activatesAt));
        council.scheduleMembership(next, 2);
        vm.stopPrank();
    }

    function test_RefusesToScheduleTheCurrentMembership() public {
        address[] memory same = _pair(notary, watch);

        vm.expectRevert(IAnchorCouncil.SameMembership.selector);
        vm.prank(admin);
        council.scheduleMembership(same, 2);
    }

    function test_RefusesToScheduleAMembershipTheConstructorWouldRefuse() public {
        address[] memory pair = _pair(notary, stranger);
        address[] memory descending = _pair(notary, stranger);
        (descending[0], descending[1]) = (descending[1], descending[0]);
        vm.startPrank(admin);

        vm.expectRevert(abi.encodeWithSelector(IAnchorCouncil.ThresholdTooLow.selector, 1, 2));
        council.scheduleMembership(pair, 1);

        vm.expectRevert(abi.encodeWithSelector(IAnchorCouncil.TooFewMembers.selector, 2, 3));
        council.scheduleMembership(pair, 3);

        vm.expectRevert(abi.encodeWithSelector(IAnchorCouncil.MembersNotAscending.selector, 1));
        council.scheduleMembership(descending, 2);
        vm.stopPrank();
    }

    function test_RefusesToActivateOrCancelNothing() public {
        vm.expectRevert(IAnchorCouncil.NoMembershipChangePending.selector);
        council.activateMembership();

        vm.expectRevert(IAnchorCouncil.NoMembershipChangePending.selector);
        vm.prank(admin);
        council.cancelMembership();
    }

    function _ascending(uint256 count) private pure returns (address[] memory) {
        address[] memory members = new address[](count);
        for (uint256 i = 0; i < count; i++) {
            members[i] = address(uint160(i + 1));
        }

        return members;
    }

    function _ascendingOf(address a, address b, address c) private pure returns (address[] memory) {
        address[] memory three = new address[](3);
        three[0] = a;
        three[1] = b;
        three[2] = c;
        for (uint256 i = 0; i < 3; i++) {
            for (uint256 j = i + 1; j < 3; j++) {
                if (three[j] < three[i]) {
                    (three[i], three[j]) = (three[j], three[i]);
                }
            }
        }

        return three;
    }

    function _sortedApprovals(EpochHeader memory header, uint256[] memory keys) private view returns (bytes[] memory) {
        for (uint256 i = 0; i < keys.length; i++) {
            for (uint256 j = i + 1; j < keys.length; j++) {
                if (vm.addr(keys[j]) < vm.addr(keys[i])) {
                    (keys[i], keys[j]) = (keys[j], keys[i]);
                }
            }
        }

        bytes[] memory approvals = new bytes[](keys.length);
        for (uint256 i = 0; i < keys.length; i++) {
            approvals[i] = _approve(keys[i], header);
        }

        return approvals;
    }

    function _signed(uint256 key, bytes32 digest) private pure returns (bytes memory) {
        (uint8 v, bytes32 r, bytes32 s) = vm.sign(key, digest);
        return abi.encodePacked(r, s, v);
    }
}
