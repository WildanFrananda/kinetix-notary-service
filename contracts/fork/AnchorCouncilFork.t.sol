// SPDX-License-Identifier: MIT
pragma solidity 0.8.37;

import {IAccessControl} from "@openzeppelin-contracts/access/IAccessControl.sol";
import {Test} from "forge-std/src/Test.sol";

import {EpochHeader} from "../src/EpochHeader.sol";
import {ILiabilityLedger} from "../src/ILiabilityLedger.sol";
import {LiabilityLedger} from "../src/LiabilityLedger.sol";
import {SumNode} from "../src/SumNode.sol";
import {AnchorCouncil} from "../src/council/AnchorCouncil.sol";
import {IAnchorCouncil} from "../src/council/IAnchorCouncil.sol";
import {LeafSet} from "../test/support/LeafSet.sol";
import {TreeBuilder} from "../test/support/TreeBuilder.sol";

contract AnchorCouncilForkTest is Test {
    LiabilityLedger private ledger;
    address private liveAnchorer;
    address private notary;
    uint256 private notaryKey;
    address private watch;
    uint256 private watchKey;
    AnchorCouncil private council;

    function setUp() public {
        vm.createSelectFork(vm.envString("KINETIX_NOTARY_FORK_RPC_URL"), vm.envUint("KINETIX_NOTARY_FORK_BLOCK"));
        ledger = LiabilityLedger(vm.envAddress("KINETIX_NOTARY_FORK_LEDGER"));
        liveAnchorer = vm.envAddress("KINETIX_NOTARY_FORK_ANCHORER");
        (notary, notaryKey) = makeAddrAndKey("fork-council-notary");
        (watch, watchKey) = makeAddrAndKey("fork-council-watch");
        address[] memory members = new address[](2);
        (members[0], members[1]) = notary < watch ? (notary, watch) : (watch, notary);
        council = new AnchorCouncil(ledger.defaultAdmin(), 3 days, ledger, members, 2, 7 days);
    }

    function test_StandsBesideTheLiveLedger() public view {
        assertEq(block.chainid, 84_532);
        assertEq(address(council.ledger()), address(ledger));
        assertTrue(ledger.hasRole(ledger.ANCHORER_ROLE(), liveAnchorer));
        assertFalse(ledger.hasRole(ledger.ANCHORER_ROLE(), address(council)));
    }

    function test_TheHandOverMovesAnchoringFromTheHotKeyToTheCouncil() public {
        _handOver();
        EpochHeader memory header = _nextHeader(7);
        bytes32 expected = council.epochDigest(header);
        bytes[] memory approvals = _approvals(header);
        bytes32 role = ledger.ANCHORER_ROLE();

        vm.expectRevert(
            abi.encodeWithSelector(IAccessControl.AccessControlUnauthorizedAccount.selector, liveAnchorer, role)
        );
        vm.prank(liveAnchorer);
        ledger.anchor(header);

        assertEq(council.anchor(header, approvals), expected);
        assertEq(ledger.head(), expected);
        assertEq(ledger.epochAt(header.number).digest, expected);
    }

    function test_OneApprovalAnchorsNothingOnTheLiveLedger() public {
        _handOver();
        EpochHeader memory header = _nextHeader(3);
        uint64 next = ledger.nextEpoch();
        bytes[] memory one = new bytes[](1);
        one[0] = _approve(notaryKey, header);

        vm.expectRevert(abi.encodeWithSelector(IAnchorCouncil.NotEnoughApprovals.selector, 1, 2));
        council.anchor(header, one);
        assertEq(ledger.nextEpoch(), next);
    }

    function test_TheLiveLedgersOwnRefusalsComeBackThroughTheCouncil() public {
        _handOver();
        EpochHeader memory header = _nextHeader(2);
        header.previousDigest = keccak256("not the head");
        bytes[] memory approvals = _approvals(header);
        bytes32 head = ledger.head();

        vm.expectRevert(
            abi.encodeWithSelector(ILiabilityLedger.PreviousDigestMismatch.selector, head, header.previousDigest)
        );
        council.anchor(header, approvals);
    }

    function _handOver() private {
        bytes32 role = ledger.ANCHORER_ROLE();
        address admin = ledger.defaultAdmin();
        vm.startPrank(admin);
        ledger.grantRole(role, address(council));
        ledger.revokeRole(role, liveAnchorer);
        vm.stopPrank();
    }

    function _nextHeader(uint256 count) private view returns (EpochHeader memory) {
        uint64 number = ledger.nextEpoch();
        SumNode[][] memory levels =
            TreeBuilder.build(number, LeafSet.commitments(number, count), LeafSet.amounts(number, count));
        SumNode memory root = TreeBuilder.root(levels);
        return EpochHeader({
            number: number,
            previousDigest: ledger.head(),
            root: root.hash,
            total: root.sum,
            leafCount: uint64(count),
            asOf: ledger.lastAsOf() + 1
        });
    }

    function _approve(uint256 key, EpochHeader memory header) private view returns (bytes memory) {
        (uint8 v, bytes32 r, bytes32 s) =
            vm.sign(key, council.approvalDigest(header.number, council.epochDigest(header)));
        return abi.encodePacked(r, s, v);
    }

    function _approvals(EpochHeader memory header) private view returns (bytes[] memory) {
        bytes[] memory approvals = new bytes[](2);
        bytes memory byNotary = _approve(notaryKey, header);
        bytes memory byWatch = _approve(watchKey, header);
        (approvals[0], approvals[1]) = notary < watch ? (byNotary, byWatch) : (byWatch, byNotary);
        return approvals;
    }
}
