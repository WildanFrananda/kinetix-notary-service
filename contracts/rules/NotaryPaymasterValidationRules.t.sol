// SPDX-License-Identifier: MIT
pragma solidity 0.8.37;

import {PackedUserOperation} from "account-abstraction/contracts/interfaces/PackedUserOperation.sol";
import {Vm} from "forge-std/src/Vm.sol";

import {EpochStatement} from "../src/court/EpochStatement.sol";
import {Registration} from "../src/withholding/Registration.sol";
import {PaymasterFixture} from "../test/support/PaymasterFixture.sol";

contract NotaryPaymasterValidationRulesTest is PaymasterFixture {
    uint8 private constant GAS = 0x5a;

    function test_ValidationUsesNoOpcodeTheBundlersForbid() public {
        PackedUserOperation[] memory ops = _theFourCalls();
        for (uint256 i = 0; i < ops.length; i++) {
            vm.startDebugTraceRecording();
            _validate(ops[i]);
            Vm.DebugStep[] memory steps = vm.stopAndReturnDebugTraceRecording();

            for (uint256 j = 0; j < steps.length; j++) {
                if (steps[j].contractAddr == address(this) || steps[j].contractAddr == address(vm)) {
                    continue;
                }

                assertFalse(_forbidden(steps[j].opcode), "a forbidden opcode ran during validation");
                if (steps[j].opcode == GAS) {
                    assertTrue(j + 1 < steps.length && _isCall(steps[j + 1].opcode), "GAS not followed by a call");
                }
            }
        }
    }

    function test_ValidationReadsOnlyTheCourtsAndTheLedgerAndWritesNothing() public {
        PackedUserOperation[] memory ops = _theFourCalls();
        address[4] memory written = [address(paymaster), address(court), address(withholding), address(ledger)];
        for (uint256 i = 0; i < ops.length; i++) {
            vm.record();
            _validate(ops[i]);

            (bytes32[] memory ownReads,) = vm.accesses(address(paymaster));
            assertEq(ownReads.length, 0, "the paymaster read its own storage");
            (bytes32[] memory senderReads,) = vm.accesses(holder);
            assertEq(senderReads.length, 0, "the paymaster read the sender's storage");
            for (uint256 j = 0; j < written.length; j++) {
                (, bytes32[] memory writes) = vm.accesses(written[j]);
                assertEq(writes.length, 0, "validation wrote storage");
            }
        }
    }

    function _theFourCalls() private returns (PackedUserOperation[] memory ops) {
        (uint64 epoch,) = _anchorWithAccount(1, 5, 2);
        Registration memory other = Registration({holder: holder, account: keccak256("other"), fromEpoch: 0});
        _request(other, epoch);
        EpochStatement memory honest = _true(epoch, 1, 5, 2);
        EpochStatement memory lie = _true(epoch, 1, 5, 2);
        lie.amount += 1;

        ops = new PackedUserOperation[](4);
        ops[0] = _op(holder, address(withholding), _requestCall(epoch), "");
        ops[1] = _op(
            holder,
            address(withholding),
            abi.encodeCall(withholding.recordUnanswered, (other, epoch)),
            _signRegistration(signerKey, other)
        );
        ops[2] = _op(
            holder,
            address(court),
            abi.encodeCall(court.proveBadProof, (lie, _sign(signerKey, lie), _siblings(epoch, 1, 5, 2))),
            ""
        );
        ops[3] = _op(
            holder,
            address(court),
            abi.encodeCall(court.proveEquivocation, (honest, _sign(signerKey, honest), lie, _sign(signerKey, lie))),
            ""
        );
        return ops;
    }

    function _forbidden(uint8 opcode) private pure returns (bool) {
        return opcode == 0x31 || opcode == 0x32 || opcode == 0x3a || (opcode >= 0x40 && opcode <= 0x45)
            || opcode == 0x47 || opcode == 0x48 || opcode == 0x49 || opcode == 0x4a || opcode == 0x5c || opcode == 0x5d
            || opcode == 0xf0 || opcode == 0xf5 || opcode == 0xff;
    }

    function _isCall(uint8 opcode) private pure returns (bool) {
        return opcode == 0xf1 || opcode == 0xf2 || opcode == 0xf4 || opcode == 0xfa;
    }
}
