// SPDX-License-Identifier: MIT
pragma solidity 0.8.37;

import {AnchorCouncil} from "../../src/council/AnchorCouncil.sol";
import {LedgerFixture} from "../support/LedgerFixture.sol";
import {CouncilHandler} from "./CouncilHandler.sol";

contract AnchorCouncilInvariantTest is LedgerFixture {
    uint256 private constant POOL = 5;

    AnchorCouncil private council;
    CouncilHandler private handler;

    function setUp() public override {
        super.setUp();
        uint256[] memory pool = new uint256[](POOL);
        for (uint256 i = 0; i < POOL; i++) {
            (, pool[i]) = makeAddrAndKey(string.concat("pool", vm.toString(i)));
        }

        address[] memory founders = new address[](2);
        (founders[0], founders[1]) = vm.addr(pool[0]) < vm.addr(pool[1])
            ? (vm.addr(pool[0]), vm.addr(pool[1]))
            : (vm.addr(pool[1]), vm.addr(pool[0]));
        council = new AnchorCouncil(admin, ADMIN_DELAY, ledger, founders, 2, 7 days);
        bytes32 role = ledger.ANCHORER_ROLE();
        vm.prank(admin);
        ledger.grantRole(role, address(council));

        handler = new CouncilHandler(council, admin, pool);
        targetContract(address(handler));
    }

    function invariant_NoActionHadAnOutcomeTheRulesForbid() public view {
        assertEq(handler.unexpected(), 0);
    }

    function invariant_TheLedgerHoldsExactlyWhatTheCouncilAnchored() public view {
        assertEq(ledger.nextEpoch(), handler.anchored());
    }

    function invariant_EveryEpochHadAQuorumOfItsTime() public view {
        for (uint64 epoch = 0; epoch < ledger.nextEpoch(); epoch++) {
            assertGe(handler.thresholdAt(epoch), council.MIN_THRESHOLD());
            assertGe(handler.approversOf(epoch), handler.thresholdAt(epoch));
        }
    }

    function invariant_TheMembershipIsAlwaysWellFormed() public view {
        address[] memory members = council.members();
        uint8 threshold = council.threshold();

        assertGe(threshold, council.MIN_THRESHOLD());
        assertGe(members.length, threshold);
        assertLe(members.length, council.MAX_MEMBERS());
        for (uint256 i = 0; i < members.length; i++) {
            assertTrue(council.isMember(members[i]));
            if (i > 0) {
                assertTrue(members[i - 1] < members[i]);
            }
        }
    }

    function invariant_NobodyOutsideTheListIsAMember() public view {
        address[] memory members = council.members();
        uint256 counted = 0;
        for (uint256 i = 0; i < handler.poolSize(); i++) {
            if (council.isMember(handler.poolMember(i))) {
                counted++;
            }
        }

        assertEq(counted, members.length);
        assertFalse(council.isMember(address(0)));
        assertFalse(council.isMember(admin));
    }
}
