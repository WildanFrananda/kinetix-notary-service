// SPDX-License-Identifier: MIT
pragma solidity 0.8.37;

import {EpochHeader} from "../../src/EpochHeader.sol";
import {AnchorCouncil} from "../../src/council/AnchorCouncil.sol";
import {LedgerFixture} from "../support/LedgerFixture.sol";

contract DifferentialApprovalTest is LedgerFixture {
    string private fixture;
    AnchorCouncil private council;
    address private member;

    function setUp() public override {
        super.setUp();
        fixture = vm.readFile(string.concat(vm.projectRoot(), "/test/fixtures/approval-v1.json"));
        vm.chainId(vm.parseJsonUint(fixture, ".chainId"));
        address ledgerAt = vm.parseJsonAddress(fixture, ".ledger");
        deployCodeTo(
            "LiabilityLedger.sol:LiabilityLedger", abi.encode(admin, ADMIN_DELAY, anchorer, 0, bytes32(0), 0), ledgerAt
        );

        member = vm.parseJsonAddress(fixture, ".member");
        address other = makeAddr("other-member");
        address[] memory members = new address[](2);
        (members[0], members[1]) = member < other ? (member, other) : (other, member);
        address where = vm.parseJsonAddress(fixture, ".council");
        deployCodeTo(
            "AnchorCouncil.sol:AnchorCouncil",
            abi.encode(admin, ADMIN_DELAY, ledgerAt, members, uint8(2), uint64(7 days)),
            where
        );
        council = AnchorCouncil(where);
    }

    function test_TheDomainAgrees() public view {
        (, string memory name, string memory version, uint256 chainId, address verifyingContract,,) =
            council.eip712Domain();

        assertEq(name, vm.parseJsonString(fixture, ".name"));
        assertEq(version, vm.parseJsonString(fixture, ".version"));
        assertEq(chainId, vm.parseJsonUint(fixture, ".chainId"));
        assertEq(verifyingContract, vm.parseJsonAddress(fixture, ".council"));
        assertEq(address(council.ledger()), vm.parseJsonAddress(fixture, ".ledger"));
    }

    function test_EveryApprovalGoSignsHasTheSameDigestsAndMemberHere() public view {
        uint256 count = vm.parseJsonUint(fixture, ".approvalCount");
        assertGt(count, 0);

        for (uint256 i = 0; i < count; i++) {
            string memory key = string.concat(".approvals[", vm.toString(i), "]");
            EpochHeader memory header = EpochHeader({
                number: uint64(vm.parseJsonUint(fixture, string.concat(key, ".number"))),
                previousDigest: vm.parseJsonBytes32(fixture, string.concat(key, ".previousDigest")),
                root: vm.parseJsonBytes32(fixture, string.concat(key, ".root")),
                total: uint64(vm.parseJsonUint(fixture, string.concat(key, ".total"))),
                leafCount: uint64(vm.parseJsonUint(fixture, string.concat(key, ".leafCount"))),
                asOf: uint64(vm.parseJsonUint(fixture, string.concat(key, ".asOf")))
            });
            bytes32 epoch = council.epochDigest(header);
            bytes memory signature = vm.parseJsonBytes(fixture, string.concat(key, ".signature"));

            assertEq(epoch, vm.parseJsonBytes32(fixture, string.concat(key, ".epochDigest")), key);
            assertEq(
                council.approvalDigest(header.number, epoch),
                vm.parseJsonBytes32(fixture, string.concat(key, ".digest")),
                key
            );
            assertEq(council.recoverApprover(header.number, epoch, signature), member, key);
            assertTrue(council.isMember(member), key);
        }
    }
}
