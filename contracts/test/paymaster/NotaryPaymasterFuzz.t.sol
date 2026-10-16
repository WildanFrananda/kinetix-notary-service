// SPDX-License-Identifier: MIT
pragma solidity 0.8.37;

import {BaseAccount} from "account-abstraction/contracts/core/BaseAccount.sol";
import {PackedUserOperation} from "account-abstraction/contracts/interfaces/PackedUserOperation.sol";

import {INotaryPaymaster} from "../../src/paymaster/INotaryPaymaster.sol";
import {PaymasterFixture} from "../support/PaymasterFixture.sol";

contract NotaryPaymasterFuzzTest is PaymasterFixture {
    function testFuzz_NoCallButTheFourHolderCallsIsEverSponsored(address target, bytes calldata call) public {
        PackedUserOperation memory op = _op(holder, target, call, "");

        vm.expectRevert();
        vm.prank(CANONICAL_ENTRY_POINT);
        paymaster.validatePaymasterUserOp(op, bytes32(0), 0);
    }

    function testFuzz_OnlyTheNotarysKeySignsAnythingSponsored(uint256 rawKey) public {
        uint256 key = bound(
            rawKey,
            1,
            115_792_089_237_316_195_423_570_985_008_687_907_852_837_564_279_074_904_382_605_163_141_518_161_494_336
        );
        vm.assume(key != signerKey);
        (uint64 epoch,) = _anchorWithAccount(1, 5, 2);
        bytes memory call = abi.encodeCall(
            withholding.requestStatement, (_registration(0), _signRegistration(key, _registration(0)), epoch)
        );
        PackedUserOperation memory op = _op(holder, address(withholding), call, "");

        vm.expectRevert(abi.encodeWithSelector(INotaryPaymaster.NotSignedByTheNotary.selector, vm.addr(key)));
        vm.prank(CANONICAL_ENTRY_POINT);
        paymaster.validatePaymasterUserOp(op, bytes32(0), 0);
    }

    function testFuzz_TheCeilingIsExact(uint256 maxCost) public {
        (uint64 epoch,) = _anchorWithAccount(1, 5, 2);
        PackedUserOperation memory op = _op(holder, address(withholding), _requestCall(epoch), "");

        if (maxCost > CEILING) {
            vm.expectRevert(abi.encodeWithSelector(INotaryPaymaster.CostTooHigh.selector, maxCost, CEILING));
            vm.prank(CANONICAL_ENTRY_POINT);
            paymaster.validatePaymasterUserOp(op, bytes32(0), maxCost);
            return;
        }

        vm.prank(CANONICAL_ENTRY_POINT);
        (bytes memory context, uint256 validationData) = paymaster.validatePaymasterUserOp(op, bytes32(0), maxCost);

        assertEq(context.length, 0);
        assertEq(validationData, uint256(ledger.epochAt(epoch).anchoredAt + WINDOW) << 160);
    }

    function testFuzz_AnyValueSentAlongIsRefused(uint256 value) public {
        vm.assume(value != 0);
        (uint64 epoch,) = _anchorWithAccount(1, 5, 2);
        PackedUserOperation memory op = _op(holder, address(withholding), _requestCall(epoch), "");
        op.callData = abi.encodeCall(BaseAccount.execute, (address(withholding), value, _requestCall(epoch)));

        vm.expectRevert(abi.encodeWithSelector(INotaryPaymaster.ValueNotZero.selector, value));
        vm.prank(CANONICAL_ENTRY_POINT);
        paymaster.validatePaymasterUserOp(op, bytes32(0), 0);
    }
}
