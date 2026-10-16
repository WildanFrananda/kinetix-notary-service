// SPDX-License-Identifier: MIT
pragma solidity 0.8.37;

import {IEntryPoint} from "account-abstraction/contracts/interfaces/IEntryPoint.sol";
import {PackedUserOperation} from "account-abstraction/contracts/interfaces/PackedUserOperation.sol";
import {Vm} from "forge-std/src/Vm.sol";

import {SumNode} from "../../src/SumNode.sol";
import {EpochStatement} from "../../src/court/EpochStatement.sol";
import {PaymasterFixture} from "../support/PaymasterFixture.sol";

contract NotaryPaymasterEquivalenceTest is PaymasterFixture {
    uint256 private constant STEPS = 8;

    uint64[2] private _epochs;
    uint64[2] private _accountAmounts;

    function testFuzz_AHoldersWithholdingCallIsSponsoredExactlyWhenTheCourtWouldTakeIt(
        uint8[STEPS] calldata actions,
        uint32[STEPS] calldata waits
    ) public {
        (_epochs[0], _accountAmounts[0]) = _anchorWithAccount(1, 5, 2);
        (_epochs[1], _accountAmounts[1]) = _anchorWithAccount(2, 5, 2);

        for (uint256 i = 0; i < STEPS; i++) {
            vm.warp(vm.getBlockTimestamp() + bound(waits[i], 0, 4 days));
            uint256 which = actions[i] & 1;
            uint64 epoch = _epochs[which];
            uint8 kind = (actions[i] >> 1) % 3;
            if (kind == 0) {
                _sponsoredExactlyWhenTaken(address(withholding), _requestCall(epoch), "");
            } else if (kind == 1) {
                _sponsoredExactlyWhenTaken(address(withholding), _unansweredCall(epoch), _registrationSignature());
            } else {
                _answer(which);
            }
        }
    }

    function testFuzz_AProofIsSponsoredExactlyWhenTheCourtWouldConvict(
        uint8 field,
        uint64 rawShift,
        bool provenBefore
    ) public {
        _anchorNext(1, 5);
        EpochStatement memory statement = _true(0, 1, 5, 2);
        uint64 shift = uint64(bound(rawShift, 0, 1000));
        if (field % 3 == 0) {
            statement.amount += shift;
        } else if (field % 3 == 1) {
            statement.index = shift % 8;
        } else {
            statement.commitment = bytes32(uint256(statement.commitment) ^ shift);
        }

        bytes memory signature = _sign(signerKey, statement);
        if (provenBefore) {
            try court.proveBadProof(statement, signature, _siblings(0, 1, 5, 2)) {} catch {}
        }

        bytes memory call = abi.encodeCall(court.proveBadProof, (statement, signature, _siblings(0, 1, 5, 2)));

        _sponsoredExactlyWhenTaken(address(court), call, "");
    }

    function testFuzz_AnEquivocationIsSponsoredExactlyWhenTheCourtWouldConvict(
        uint8 rawIndex,
        uint64 rawShift,
        bool sameCommitment,
        bool provenBefore
    ) public {
        _anchorNext(1, 5);
        EpochStatement memory first = _true(0, 1, 5, 2);
        EpochStatement memory second = _true(0, 1, 5, uint256(rawIndex) % 5);
        second.amount += uint64(bound(rawShift, 0, 1000));
        if (sameCommitment) {
            second.commitment = first.commitment;
        }

        bytes memory firstSignature = _sign(signerKey, first);
        bytes memory secondSignature = _sign(signerKey, second);
        if (provenBefore) {
            try court.proveEquivocation(first, firstSignature, second, secondSignature) {} catch {}
        }

        bytes memory call = abi.encodeCall(court.proveEquivocation, (first, firstSignature, second, secondSignature));

        _sponsoredExactlyWhenTaken(address(court), call, "");
    }

    function _sponsoredExactlyWhenTaken(address target, bytes memory call, bytes memory paymasterData) private {
        uint256 snapshot = vm.snapshotState();
        vm.prank(holder);
        (bool taken,) = target.call(call);
        vm.revertToState(snapshot);

        PackedUserOperation[] memory ops = new PackedUserOperation[](1);
        ops[0] = _signedBy(holderKey, _op(holder, target, call, paymasterData));
        uint256 before = paymaster.getDeposit();
        vm.recordLogs();
        vm.prank(bundler, bundler);
        (bool sponsored,) = address(entryPoint).call(abi.encodeCall(IEntryPoint.handleOps, (ops, payable(bundler))));

        assertEq(sponsored, taken, "the paymaster and the court disagree");
        if (!sponsored) {
            assertEq(paymaster.getDeposit(), before);
            return;
        }

        assertTrue(_succeeded(vm.getRecordedLogs()), "a sponsored call failed");
    }

    function _succeeded(Vm.Log[] memory logs) private pure returns (bool) {
        for (uint256 i = 0; i < logs.length; i++) {
            if (logs[i].topics[0] == IEntryPoint.UserOperationEvent.selector) {
                (, bool success,,) = abi.decode(logs[i].data, (uint256, bool, uint256, uint256));
                return success;
            }
        }

        return false;
    }

    function _answer(uint256 which) private {
        uint64 epoch = _epochs[which];
        SumNode[] memory siblings = _accountSiblings(epoch, which + 1, 5, 2);
        try withholding.answer(_registration(0), epoch, _salt(epoch), _accountAmounts[which], 2, siblings) {} catch {}
    }
}
