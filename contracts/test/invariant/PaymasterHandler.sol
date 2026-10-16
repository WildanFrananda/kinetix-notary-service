// SPDX-License-Identifier: MIT
pragma solidity 0.8.37;

import {BaseAccount} from "account-abstraction/contracts/core/BaseAccount.sol";
import {IEntryPoint} from "account-abstraction/contracts/interfaces/IEntryPoint.sol";
import {PackedUserOperation} from "account-abstraction/contracts/interfaces/PackedUserOperation.sol";
import {Test} from "forge-std/src/Test.sol";

import {NotaryPaymaster} from "../../src/paymaster/NotaryPaymaster.sol";
import {Registration} from "../../src/withholding/Registration.sol";
import {WithholdingCourt} from "../../src/withholding/WithholdingCourt.sol";

contract PaymasterHandler is Test {
    NotaryPaymaster public immutable paymaster;
    IEntryPoint public immutable entryPoint;
    WithholdingCourt public immutable withholding;
    address public immutable holder;
    address public immutable admin;
    address public immutable bundler;

    uint256 private immutable _holderKey;
    uint256 private immutable _signerKey;
    Registration private _registration;
    uint64 private immutable _epoch;

    uint256 public deposited;
    uint256 public withdrawn;
    uint256 public spent;
    uint256 public staked;
    uint256 public unexpected;

    constructor(
        NotaryPaymaster paymaster_,
        WithholdingCourt withholding_,
        uint256 holderKey,
        uint256 signerKey,
        Registration memory registration,
        uint64 epoch,
        address admin_
    ) {
        paymaster = paymaster_;
        entryPoint = paymaster_.entryPoint();
        withholding = withholding_;
        holder = vm.addr(holderKey);
        admin = admin_;
        bundler = makeAddr("handler-bundler");
        _holderKey = holderKey;
        _signerKey = signerKey;
        _registration = registration;
        _epoch = epoch;
        deposited = paymaster_.getDeposit();
    }

    function request(uint256 rawFee) external {
        bytes memory call =
            abi.encodeCall(withholding.requestStatement, (_registration, _signature(_signerKey), _epoch));
        _judged(call, "", bound(rawFee, 1, 1 gwei));
    }

    function unanswered(uint256 rawWait, uint256 rawFee) external {
        vm.warp(vm.getBlockTimestamp() + bound(rawWait, 0, 2 days));
        bytes memory call = abi.encodeCall(withholding.recordUnanswered, (_registration, _epoch));
        _judged(call, _signature(_signerKey), bound(rawFee, 1, 1 gwei));
    }

    function forged(uint256 rawKey, uint256 rawFee) external {
        uint256 key = bound(rawKey, 1, 1_000_000);
        if (key == _signerKey) {
            return;
        }

        bytes memory call = abi.encodeCall(withholding.requestStatement, (_registration, _signature(key), _epoch));
        _refused(_op(address(withholding), call, "", bound(rawFee, 1, 1 gwei)));
    }

    function elsewhere(address target, bytes calldata call, uint256 rawFee) external {
        if (target == address(withholding) || target == address(paymaster.court())) {
            return;
        }

        _refused(_op(target, call, "", bound(rawFee, 1, 1 gwei)));
    }

    function deposit(uint256 rawAmount) external {
        uint256 amount = bound(rawAmount, 1, 1 ether);
        vm.deal(address(this), amount);
        entryPoint.depositTo{value: amount}(address(paymaster));
        deposited += amount;
    }

    function withdraw(uint256 rawAmount) external {
        uint256 available = paymaster.getDeposit();
        if (available == 0) {
            return;
        }

        uint256 amount = bound(rawAmount, 1, available);
        vm.prank(admin);
        paymaster.withdrawTo(payable(makeAddr("handler-treasury")), amount);
        withdrawn += amount;
    }

    function stake(uint256 rawAmount) external {
        uint256 amount = bound(rawAmount, 1, 1 ether);
        vm.deal(admin, amount);
        vm.prank(admin);
        paymaster.addStake{value: amount}();
        staked += amount;
    }

    function _judged(bytes memory call, bytes memory paymasterData, uint256 fee) private {
        uint256 snapshot = vm.snapshotState();
        vm.prank(holder);
        (bool taken,) = address(withholding).call(call);
        vm.revertToState(snapshot);

        PackedUserOperation memory op = _op(address(withholding), call, paymasterData, fee);
        uint256 before = paymaster.getDeposit();
        if (before < 600_000 * fee) {
            _refused(op);
            return;
        }

        bool sponsored = _handle(op);
        if (sponsored != taken) {
            unexpected++;
            return;
        }

        spent += before - paymaster.getDeposit();
    }

    function _refused(PackedUserOperation memory op) private {
        uint256 before = paymaster.getDeposit();
        if (_handle(op) || paymaster.getDeposit() != before) {
            unexpected++;
        }
    }

    function _handle(PackedUserOperation memory op) private returns (bool) {
        PackedUserOperation[] memory ops = new PackedUserOperation[](1);
        ops[0] = op;
        vm.prank(bundler, bundler);
        try entryPoint.handleOps(ops, payable(bundler)) {
            return true;
        } catch {
            return false;
        }
    }

    function _op(
        address target,
        bytes memory call,
        bytes memory paymasterData,
        uint256 fee
    ) private view returns (PackedUserOperation memory op) {
        op.sender = holder;
        op.nonce = entryPoint.getNonce(holder, 0);
        op.callData = abi.encodeCall(BaseAccount.execute, (target, 0, call));
        op.accountGasLimits = bytes32(uint256(150_000) << 128 | 300_000);
        op.preVerificationGas = 50_000;
        op.gasFees = bytes32(fee << 128 | fee);
        op.paymasterAndData = abi.encodePacked(address(paymaster), uint128(100_000), uint128(0), paymasterData);
        (uint8 v, bytes32 r, bytes32 s) = vm.sign(_holderKey, entryPoint.getUserOpHash(op));
        op.signature = abi.encodePacked(r, s, v);
        return op;
    }

    function _signature(uint256 key) private view returns (bytes memory) {
        (uint8 v, bytes32 r, bytes32 s) = vm.sign(key, withholding.registrationDigest(_registration));
        return abi.encodePacked(r, s, v);
    }
}
