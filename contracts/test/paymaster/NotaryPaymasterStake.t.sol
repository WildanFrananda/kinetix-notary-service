// SPDX-License-Identifier: MIT
pragma solidity 0.8.37;

import {IAccessControl} from "@openzeppelin-contracts/access/IAccessControl.sol";
import {IStakeManager} from "account-abstraction/contracts/interfaces/IStakeManager.sol";

import {INotaryPaymaster} from "../../src/paymaster/INotaryPaymaster.sol";
import {NotaryPaymaster} from "../../src/paymaster/NotaryPaymaster.sol";
import {PaymasterFixture} from "../support/PaymasterFixture.sol";

contract NotaryPaymasterStakeTest is PaymasterFixture {
    function test_TheAdminStakesWithThePinnedDelayAndMayAddToIt() public {
        vm.deal(admin, 2 ether);

        vm.prank(admin);
        paymaster.addStake{value: 1 ether}();

        IStakeManager.DepositInfo memory info = paymaster.getDepositInfo();
        assertTrue(info.staked);
        assertEq(info.stake, 1 ether);
        assertEq(info.unstakeDelaySec, UNSTAKE_DELAY);
        assertEq(info.withdrawTime, 0);
        assertEq(info.deposit, 1 ether);

        vm.prank(admin);
        paymaster.addStake{value: 0.5 ether}();
        assertEq(paymaster.getDepositInfo().stake, 1.5 ether);
    }

    function test_OnlyTheAdminTouchesTheStake() public {
        vm.deal(stranger, 1 ether);
        bytes memory refusal = abi.encodeWithSelector(
            IAccessControl.AccessControlUnauthorizedAccount.selector, stranger, paymaster.DEFAULT_ADMIN_ROLE()
        );

        vm.expectRevert(refusal);
        vm.prank(stranger);
        paymaster.addStake{value: 1 ether}();

        vm.expectRevert(refusal);
        vm.prank(stranger);
        paymaster.unlockStake();

        vm.expectRevert(refusal);
        vm.prank(stranger);
        paymaster.withdrawStake(payable(stranger));
    }

    function test_TheStakeComesBackOnlyAfterTheDelay() public {
        vm.deal(admin, 1 ether);
        vm.prank(admin);
        paymaster.addStake{value: 1 ether}();

        vm.prank(admin);
        paymaster.unlockStake();
        IStakeManager.DepositInfo memory info = paymaster.getDepositInfo();
        assertFalse(info.staked);
        assertEq(info.withdrawTime, vm.getBlockTimestamp() + UNSTAKE_DELAY);

        address payable treasury = payable(makeAddr("treasury"));
        vm.expectRevert("Stake withdrawal is not due");
        vm.prank(admin);
        paymaster.withdrawStake(treasury);

        vm.warp(info.withdrawTime);
        vm.prank(admin);
        paymaster.withdrawStake(treasury);
        assertEq(treasury.balance, 1 ether);
        assertEq(paymaster.getDepositInfo().stake, 0);
        assertEq(paymaster.getDeposit(), 1 ether);
    }

    function test_OnlyTheAdminTakesBackEtherForcedIntoThePaymaster() public {
        vm.deal(address(paymaster), 1 ether);

        vm.expectRevert(
            abi.encodeWithSelector(
                IAccessControl.AccessControlUnauthorizedAccount.selector, stranger, paymaster.DEFAULT_ADMIN_ROLE()
            )
        );
        vm.prank(stranger);
        paymaster.withdrawBalance();

        uint256 before = admin.balance;
        vm.prank(admin);
        paymaster.withdrawBalance();
        assertEq(address(paymaster).balance, 0);
        assertEq(admin.balance, before + 1 ether);
        assertEq(paymaster.getDeposit(), 1 ether);
    }

    function test_EtherTheAdminCannotTakeStaysWhereItIs() public {
        vm.deal(address(paymaster), 1 ether);
        vm.etch(admin, hex"fe");

        vm.expectRevert(abi.encodeWithSelector(INotaryPaymaster.BalanceNotSent.selector, admin, 1 ether));
        vm.prank(admin);
        paymaster.withdrawBalance();
        assertEq(address(paymaster).balance, 1 ether);
    }

    function test_RefusesAnUnstakeDelayShorterThanADay() public {
        vm.expectRevert(
            abi.encodeWithSelector(INotaryPaymaster.UnstakeDelayTooShort.selector, uint32(1 days - 1), uint32(1 days))
        );
        new NotaryPaymaster(admin, ADMIN_DELAY, entryPoint, address(delegate), court, withholding, CEILING, 1 days - 1);
    }
}
