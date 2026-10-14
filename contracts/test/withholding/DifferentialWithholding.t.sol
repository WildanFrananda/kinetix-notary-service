// SPDX-License-Identifier: MIT
pragma solidity 0.8.37;

import {IStatementCourt} from "../../src/court/IStatementCourt.sol";
import {Registration} from "../../src/withholding/Registration.sol";
import {WithholdingCourt} from "../../src/withholding/WithholdingCourt.sol";
import {LedgerFixture} from "../support/LedgerFixture.sol";

contract DifferentialWithholdingTest is LedgerFixture {
    string private fixture;
    WithholdingCourt private withholding;

    function setUp() public override {
        super.setUp();
        fixture = vm.readFile(string.concat(vm.projectRoot(), "/test/fixtures/withholding-v1.json"));
        vm.chainId(vm.parseJsonUint(fixture, ".chainId"));

        address court = vm.parseJsonAddress(fixture, ".court");
        deployCodeTo(
            "StatementCourt.sol:StatementCourt",
            abi.encode(admin, ADMIN_DELAY, ledger, vm.parseJsonAddress(fixture, ".signer"), uint64(7 days)),
            court
        );

        address where = vm.parseJsonAddress(fixture, ".withholdingCourt");
        deployCodeTo(
            "WithholdingCourt.sol:WithholdingCourt",
            abi.encode(IStatementCourt(court), uint64(7 days), uint64(24 hours)),
            where
        );
        withholding = WithholdingCourt(where);
    }

    function test_TheDomainAgrees() public view {
        (, string memory name, string memory version, uint256 chainId, address verifyingContract,,) =
            withholding.eip712Domain();

        assertEq(name, vm.parseJsonString(fixture, ".name"));
        assertEq(version, vm.parseJsonString(fixture, ".version"));
        assertEq(chainId, vm.parseJsonUint(fixture, ".chainId"));
        assertEq(verifyingContract, vm.parseJsonAddress(fixture, ".withholdingCourt"));
    }

    function test_EveryRegistrationGoSignsHasTheSameDigestAndSignerHere() public view {
        uint256 count = vm.parseJsonUint(fixture, ".registrationCount");
        address signer = vm.parseJsonAddress(fixture, ".signer");
        assertGt(count, 0);

        for (uint256 i = 0; i < count; i++) {
            string memory key = string.concat(".registrations[", vm.toString(i), "]");
            Registration memory registration = Registration({
                holder: vm.parseJsonAddress(fixture, string.concat(key, ".holder")),
                account: vm.parseJsonBytes32(fixture, string.concat(key, ".account")),
                fromEpoch: uint64(vm.parseJsonUint(fixture, string.concat(key, ".fromEpoch")))
            });
            bytes memory signature = vm.parseJsonBytes(fixture, string.concat(key, ".signature"));

            assertEq(
                withholding.registrationDigest(registration),
                vm.parseJsonBytes32(fixture, string.concat(key, ".digest")),
                key
            );
            assertEq(withholding.recoverSigner(registration, signature), signer, key);
        }
    }

    function test_EveryAccountGoDerivesCommitsTheSameHere() public view {
        uint256 count = vm.parseJsonUint(fixture, ".accountCount");
        assertGt(count, 0);

        for (uint256 i = 0; i < count; i++) {
            string memory key = string.concat(".accounts[", vm.toString(i), "]");
            bytes32 account = vm.parseJsonBytes32(fixture, string.concat(key, ".account"));
            bytes32 preimage = keccak256(
                abi.encodePacked(
                    bytes1(0x05),
                    uint8(vm.parseJsonUint(fixture, string.concat(key, ".kind"))),
                    vm.parseJsonBytes32(fixture, string.concat(key, ".secret")),
                    bytes(vm.parseJsonString(fixture, string.concat(key, ".holder")))
                )
            );

            assertEq(account, preimage, key);
            assertEq(
                withholding.commitmentOf(account, vm.parseJsonBytes32(fixture, string.concat(key, ".salt"))),
                vm.parseJsonBytes32(fixture, string.concat(key, ".commitment")),
                key
            );
        }
    }
}
