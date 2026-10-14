// SPDX-License-Identifier: MIT
pragma solidity 0.8.37;

import {SumNode} from "../../src/SumNode.sol";
import {IWithholdingCourt} from "../../src/withholding/IWithholdingCourt.sol";
import {Registration} from "../../src/withholding/Registration.sol";
import {RequestState} from "../../src/withholding/RequestState.sol";
import {WithholdingFixture} from "../support/WithholdingFixture.sol";

contract WithholdingCourtFuzzTest is WithholdingFixture {
    uint256 private constant CURVE_ORDER = 0xFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFEBAAEDCE6AF48A03BBFD25E8CD0364141;

    function testFuzz_AnIncludedAccountCanAlwaysBeAnswered(uint256 seed, uint256 rawCount, uint256 rawIndex) public {
        uint256 count = bound(rawCount, 1, 40);
        uint256 index = bound(rawIndex, 0, count - 1);
        (uint64 epoch, uint64 amount) = _anchorWithAccount(seed, count, index);
        Registration memory registration = _registration(0);
        bytes32 digest = _request(registration, epoch);
        SumNode[] memory siblings = _accountSiblings(epoch, seed, count, index);

        withholding.answer(registration, epoch, _salt(epoch), amount, uint64(index), siblings);

        assertEq(uint8(withholding.requestOf(digest, epoch).state), uint8(RequestState.Answered));
    }

    function testFuzz_AWrongSaltOrAmountNeverAnswers(
        uint256 seed,
        uint256 rawCount,
        uint256 rawIndex,
        bytes32 salt,
        uint64 rawDelta
    ) public {
        uint256 count = bound(rawCount, 1, 40);
        uint256 index = bound(rawIndex, 0, count - 1);
        (uint64 epoch, uint64 amount) = _anchorWithAccount(seed, count, index);
        vm.assume(salt != _salt(epoch));
        Registration memory registration = _registration(0);
        _request(registration, epoch);
        SumNode[] memory siblings = _accountSiblings(epoch, seed, count, index);
        uint64 wrong;
        unchecked {
            wrong = amount + uint64(bound(rawDelta, 1, type(uint64).max));
        }

        vm.expectRevert(abi.encodeWithSelector(IWithholdingCourt.AnswerDoesNotHold.selector, epoch));
        withholding.answer(registration, epoch, salt, amount, uint64(index), siblings);

        vm.expectRevert(abi.encodeWithSelector(IWithholdingCourt.AnswerDoesNotHold.selector, epoch));
        withholding.answer(registration, epoch, _salt(epoch), wrong, uint64(index), siblings);
    }

    function testFuzz_AnAccountLeftOutIsNeverAnswered(
        uint256 seed,
        uint256 rawCount,
        uint256 rawIndex,
        bytes32 salt
    ) public {
        uint256 count = bound(rawCount, 1, 40);
        uint256 index = bound(rawIndex, 0, count - 1);
        uint64 epoch = _anchorWithoutAccount(seed, count);
        Registration memory registration = _registration(0);
        _request(registration, epoch);
        SumNode[] memory siblings = _siblings(epoch, seed, count, index);
        uint64 amount = _amounts(seed, count)[index];

        vm.expectRevert(abi.encodeWithSelector(IWithholdingCourt.AnswerDoesNotHold.selector, epoch));
        withholding.answer(registration, epoch, salt, amount, uint64(index), siblings);
    }

    function testFuzz_NobodyButTheHolderMayRequest(address caller) public {
        vm.assume(caller != holder);
        (uint64 epoch,) = _anchorWithAccount(1, 5, 2);
        Registration memory registration = _registration(0);
        bytes memory signature = _signRegistration(signerKey, registration);

        vm.expectRevert(abi.encodeWithSelector(IWithholdingCourt.NotTheHolder.selector, holder, caller));
        vm.prank(caller);
        withholding.requestStatement(registration, signature, epoch);
    }

    function testFuzz_ARegistrationSignedByAnyOtherKeyIsRefused(uint256 rawKey) public {
        uint256 key = bound(rawKey, 1, CURVE_ORDER - 1);
        vm.assume(key != signerKey);
        (uint64 epoch,) = _anchorWithAccount(1, 5, 2);
        Registration memory registration = _registration(0);
        bytes memory signature = _signRegistration(key, registration);

        vm.expectRevert(
            abi.encodeWithSelector(IWithholdingCourt.NotSignedBySigner.selector, uint64(0), signer, vm.addr(key))
        );
        vm.prank(holder);
        withholding.requestStatement(registration, signature, epoch);
    }

    function testFuzz_TheDeadlineSeparatesAnAnswerFromASilence(uint256 rawWait) public {
        uint256 wait = bound(rawWait, 0, 3 * uint256(PERIOD));
        (uint64 epoch, uint64 amount) = _anchorWithAccount(1, 5, 2);
        Registration memory registration = _registration(0);
        bytes32 digest = _request(registration, epoch);
        SumNode[] memory siblings = _accountSiblings(epoch, 1, 5, 2);
        vm.warp(vm.getBlockTimestamp() + wait);

        if (wait <= PERIOD) {
            withholding.answer(registration, epoch, _salt(epoch), amount, 2, siblings);
            assertEq(uint8(withholding.requestOf(digest, epoch).state), uint8(RequestState.Answered));
        } else {
            uint64 dueAt = withholding.requestOf(digest, epoch).dueAt;
            uint64 blockTime = uint64(vm.getBlockTimestamp());
            vm.expectRevert(abi.encodeWithSelector(IWithholdingCourt.AnswerTooLate.selector, dueAt, blockTime));
            withholding.answer(registration, epoch, _salt(epoch), amount, 2, siblings);

            withholding.recordUnanswered(registration, epoch);
            assertEq(uint8(withholding.requestOf(digest, epoch).state), uint8(RequestState.Unanswered));
        }
    }
}
