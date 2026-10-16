// SPDX-License-Identifier: MIT
pragma solidity 0.8.37;

import {EntryPoint} from "account-abstraction/contracts/core/EntryPoint.sol";
import {IEntryPoint} from "account-abstraction/contracts/interfaces/IEntryPoint.sol";
import {PackedUserOperation} from "account-abstraction/contracts/interfaces/PackedUserOperation.sol";
import {Test} from "forge-std/src/Test.sol";

contract DifferentialOperationTest is Test {
    string private fixture;
    IEntryPoint private entryPoint;

    function setUp() public {
        fixture = vm.readFile(string.concat(vm.projectRoot(), "/test/fixtures/operation-v1.json"));
        vm.chainId(vm.parseJsonUint(fixture, ".chainId"));
        address where = vm.parseJsonAddress(fixture, ".entryPoint");
        vm.etch(where, address(new EntryPoint()).code);
        entryPoint = IEntryPoint(where);
    }

    function test_TheDomainAgrees() public view {
        assertEq(vm.parseJsonString(fixture, ".name"), "ERC4337");
        assertEq(vm.parseJsonString(fixture, ".version"), "1");
        assertEq(block.chainid, 84_532);
    }

    function test_EveryOperationGoHashesHashesTheSameInTheEntryPoint() public {
        uint256 count = vm.parseJsonUint(fixture, ".operationCount");
        assertGt(count, 0);

        for (uint256 i = 0; i < count; i++) {
            string memory key = string.concat(".operations[", vm.toString(i), "]");
            PackedUserOperation memory op = _operation(key);

            assertEq(entryPoint.getUserOpHash(op), vm.parseJsonBytes32(fixture, string.concat(key, ".hash")), key);
        }
    }

    function _operation(string memory key) private returns (PackedUserOperation memory op) {
        op.sender = vm.parseJsonAddress(fixture, string.concat(key, ".sender"));
        op.nonce = vm.parseJsonUint(fixture, string.concat(key, ".nonce"));
        address delegate = vm.parseJsonAddress(fixture, string.concat(key, ".delegate"));
        if (delegate != address(0)) {
            vm.etch(op.sender, abi.encodePacked(hex"ef0100", delegate));
            op.initCode = abi.encodePacked(bytes20(bytes2(0x7702)));
        }

        op.callData = vm.parseJsonBytes(fixture, string.concat(key, ".callData"));
        op.accountGasLimits = _pack(_uint(key, ".verificationGasLimit"), _uint(key, ".callGasLimit"));
        op.preVerificationGas = _uint(key, ".preVerificationGas");
        op.gasFees = _pack(_uint(key, ".maxPriorityFeePerGas"), _uint(key, ".maxFeePerGas"));
        address paymaster = vm.parseJsonAddress(fixture, string.concat(key, ".paymaster"));
        if (paymaster != address(0)) {
            op.paymasterAndData = abi.encodePacked(
                paymaster,
                uint128(_uint(key, ".paymasterVerificationGasLimit")),
                uint128(_uint(key, ".paymasterPostOpGasLimit")),
                vm.parseJsonBytes(fixture, string.concat(key, ".paymasterData"))
            );
        }

        return op;
    }

    function _uint(string memory key, string memory field) private view returns (uint256) {
        return vm.parseJsonUint(fixture, string.concat(key, field));
    }

    function _pack(uint256 high, uint256 low) private pure returns (bytes32) {
        return bytes32(high << 128 | low);
    }
}
