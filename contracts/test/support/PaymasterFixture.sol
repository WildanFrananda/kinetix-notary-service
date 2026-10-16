// SPDX-License-Identifier: MIT
pragma solidity 0.8.37;

import {Simple7702Account} from "account-abstraction/contracts/accounts/Simple7702Account.sol";
import {BaseAccount} from "account-abstraction/contracts/core/BaseAccount.sol";
import {EntryPoint} from "account-abstraction/contracts/core/EntryPoint.sol";
import {IEntryPoint} from "account-abstraction/contracts/interfaces/IEntryPoint.sol";
import {PackedUserOperation} from "account-abstraction/contracts/interfaces/PackedUserOperation.sol";

import {NotaryPaymaster} from "../../src/paymaster/NotaryPaymaster.sol";
import {WithholdingFixture} from "./WithholdingFixture.sol";

abstract contract PaymasterFixture is WithholdingFixture {
    address internal constant CANONICAL_ENTRY_POINT = 0x4337084D9E255Ff0702461CF8895CE9E3b5Ff108;
    uint256 internal constant CEILING = 0.001 ether;
    uint32 internal constant UNSTAKE_DELAY = 1 days;
    uint256 internal constant VERIFICATION_GAS = 150_000;
    uint256 internal constant CALL_GAS = 300_000;
    uint256 internal constant PAYMASTER_GAS = 100_000;
    uint256 internal constant PRE_VERIFICATION_GAS = 50_000;
    uint256 internal constant FEE = 1 gwei;

    IEntryPoint internal entryPoint;
    Simple7702Account internal delegate;
    NotaryPaymaster internal paymaster;
    address internal bundler;

    function setUp() public virtual override {
        super.setUp();
        bundler = makeAddr("bundler");
        vm.etch(CANONICAL_ENTRY_POINT, address(new EntryPoint()).code);
        entryPoint = IEntryPoint(CANONICAL_ENTRY_POINT);
        delegate = new Simple7702Account();
        paymaster = new NotaryPaymaster(
            admin, ADMIN_DELAY, entryPoint, address(delegate), court, withholding, CEILING, UNSTAKE_DELAY
        );
        vm.deal(address(this), 1 ether);
        entryPoint.depositTo{value: 1 ether}(address(paymaster));
        vm.signAndAttachDelegation(address(delegate), holderKey);
    }

    function _op(
        address sender,
        address target,
        bytes memory call,
        bytes memory paymasterData
    ) internal view returns (PackedUserOperation memory op) {
        op.sender = sender;
        op.nonce = entryPoint.getNonce(sender, 0);
        op.callData = abi.encodeCall(BaseAccount.execute, (target, 0, call));
        op.accountGasLimits = _pack(VERIFICATION_GAS, CALL_GAS);
        op.preVerificationGas = PRE_VERIFICATION_GAS;
        op.gasFees = _pack(FEE, FEE);
        op.paymasterAndData = abi.encodePacked(address(paymaster), uint128(PAYMASTER_GAS), uint128(0), paymasterData);
        return op;
    }

    function _signedBy(uint256 key, PackedUserOperation memory op) internal view returns (PackedUserOperation memory) {
        bytes32 hash = entryPoint.getUserOpHash(op);
        (uint8 v, bytes32 r, bytes32 s) = vm.sign(key, hash);
        op.signature = abi.encodePacked(r, s, v);
        return op;
    }

    function _handle(PackedUserOperation memory op) internal {
        PackedUserOperation[] memory ops = new PackedUserOperation[](1);
        ops[0] = op;
        vm.prank(bundler, bundler);
        entryPoint.handleOps(ops, payable(bundler));
    }

    function _validate(PackedUserOperation memory op) internal returns (uint48 validAfter, uint48 validUntil) {
        vm.prank(CANONICAL_ENTRY_POINT);
        (bytes memory context, uint256 validationData) = paymaster.validatePaymasterUserOp(op, bytes32(0), CEILING);
        assertEq(context.length, 0);
        assertEq(uint160(validationData), 0);
        return (uint48(validationData >> 208), uint48(validationData >> 160));
    }

    function _expectPaymasterRefusal(bytes memory reason) internal {
        vm.expectRevert(abi.encodeWithSelector(IEntryPoint.FailedOpWithRevert.selector, 0, "AA33 reverted", reason));
    }

    function _expectOutOfTime() internal {
        vm.expectRevert(abi.encodeWithSelector(IEntryPoint.FailedOp.selector, 0, "AA32 paymaster expired or not due"));
    }

    function _requestCall(uint64 epoch) internal view returns (bytes memory) {
        return abi.encodeCall(
            withholding.requestStatement, (_registration(0), _signRegistration(signerKey, _registration(0)), epoch)
        );
    }

    function _unansweredCall(uint64 epoch) internal view returns (bytes memory) {
        return abi.encodeCall(withholding.recordUnanswered, (_registration(0), epoch));
    }

    function _registrationSignature() internal view returns (bytes memory) {
        return _signRegistration(signerKey, _registration(0));
    }

    function _pack(uint256 high, uint256 low) internal pure returns (bytes32) {
        return bytes32(high << 128 | low);
    }
}
