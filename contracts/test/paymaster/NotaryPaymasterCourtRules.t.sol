// SPDX-License-Identifier: MIT
pragma solidity 0.8.37;

import {PackedUserOperation} from "account-abstraction/contracts/interfaces/PackedUserOperation.sol";

import {SumNode} from "../../src/SumNode.sol";
import {EpochStatement} from "../../src/court/EpochStatement.sol";
import {INotaryPaymaster} from "../../src/paymaster/INotaryPaymaster.sol";
import {Registration} from "../../src/withholding/Registration.sol";
import {RequestState} from "../../src/withholding/RequestState.sol";
import {PaymasterFixture} from "../support/PaymasterFixture.sol";

contract NotaryPaymasterCourtRulesTest is PaymasterFixture {
    function test_ARequestIsSponsoredUntilTheCourtClosesItsWindow() public {
        (uint64 epoch,) = _anchorWithAccount(1, 5, 2);
        uint64 closesAt = ledger.epochAt(epoch).anchoredAt + WINDOW;
        PackedUserOperation memory op = _signedBy(holderKey, _op(holder, address(withholding), _requestCall(epoch), ""));

        (uint48 validAfter, uint48 validUntil) = _validate(op);
        assertEq(validAfter, 0);
        assertEq(validUntil, closesAt);

        vm.warp(closesAt + 1);
        _expectOutOfTime();
        _handle(op);

        vm.warp(closesAt);
        _handle(op);
        bytes32 digest = withholding.registrationDigest(_registration(0));
        assertEq(uint8(withholding.requestOf(digest, epoch).state), uint8(RequestState.Open));
    }

    function test_AnUnansweredRecordIsSponsoredOnlyOnceTheAnswerIsOverdue() public {
        (uint64 epoch,) = _anchorWithAccount(1, 5, 2);
        bytes32 digest = _request(_registration(0), epoch);
        uint64 dueAt = withholding.requestOf(digest, epoch).dueAt;
        PackedUserOperation memory op =
            _signedBy(holderKey, _op(holder, address(withholding), _unansweredCall(epoch), _registrationSignature()));

        (uint48 validAfter, uint48 validUntil) = _validate(op);
        assertEq(validAfter, dueAt);
        assertEq(validUntil, 0);

        vm.warp(dueAt);
        _expectOutOfTime();
        _handle(op);

        vm.warp(dueAt + 1);
        _handle(op);
        assertEq(withholding.silenceCount(epoch), 1);
    }

    function test_ASecondRequestWhileTheFirstIsOpenIsNotSponsored() public {
        (uint64 first,) = _anchorWithAccount(1, 5, 2);
        (uint64 second,) = _anchorWithAccount(2, 5, 2);
        bytes32 digest = _request(_registration(0), first);

        _refusedAsTheCourtWould(
            address(withholding),
            _requestCall(second),
            "",
            abi.encodeWithSelector(INotaryPaymaster.RequestStillOpen.selector, digest, first)
        );
    }

    function test_ARequestAlreadyMadeIsNotSponsoredAgain() public {
        (uint64 epoch,) = _anchorWithAccount(1, 5, 2);
        bytes32 digest = _request(_registration(0), epoch);
        vm.warp(withholding.requestOf(digest, epoch).dueAt + 1);
        withholding.recordUnanswered(_registration(0), epoch);

        _refusedAsTheCourtWould(
            address(withholding),
            _requestCall(epoch),
            "",
            abi.encodeWithSelector(INotaryPaymaster.AlreadyRequested.selector, digest, epoch)
        );
    }

    function test_ARequestBeforeTheRegistrationOrForAnEpochNotAnchoredIsNotSponsored() public {
        (uint64 epoch,) = _anchorWithAccount(1, 5, 2);
        Registration memory later = _registration(epoch + 1);
        bytes memory call =
            abi.encodeCall(withholding.requestStatement, (later, _signRegistration(signerKey, later), epoch));
        _refusedAsTheCourtWould(
            address(withholding),
            call,
            "",
            abi.encodeWithSelector(INotaryPaymaster.BeforeRegistration.selector, epoch, epoch + 1)
        );

        _refusedAsTheCourtWould(
            address(withholding),
            _requestCall(epoch + 1),
            "",
            abi.encodeWithSelector(INotaryPaymaster.EpochNotAnchored.selector, epoch + 1)
        );
    }

    function test_AnUnansweredRecordWithoutAnOpenRequestIsNotSponsored() public {
        (uint64 epoch, uint64 amount) = _anchorWithAccount(1, 5, 2);
        bytes32 digest = withholding.registrationDigest(_registration(0));
        _refusedAsTheCourtWould(
            address(withholding),
            _unansweredCall(epoch),
            _registrationSignature(),
            abi.encodeWithSelector(INotaryPaymaster.NoOpenRequest.selector, digest, epoch)
        );

        _request(_registration(0), epoch);
        withholding.answer(_registration(0), epoch, _salt(epoch), amount, 2, _accountSiblings(epoch, 1, 5, 2));
        vm.warp(withholding.requestOf(digest, epoch).dueAt + 1);
        _refusedAsTheCourtWould(
            address(withholding),
            _unansweredCall(epoch),
            _registrationSignature(),
            abi.encodeWithSelector(INotaryPaymaster.NoOpenRequest.selector, digest, epoch)
        );
    }

    function test_ASoundProofIsNotSponsored() public {
        _anchorNext(1, 5);
        EpochStatement memory statement = _true(0, 1, 5, 2);

        _refusedAsTheCourtWould(
            address(court),
            _badProofCall(statement, _siblings(0, 1, 5, 2)),
            "",
            abi.encodeWithSelector(INotaryPaymaster.ProofHolds.selector, uint64(0))
        );
    }

    function test_AFaultAlreadyProvenIsNotSponsoredAgain() public {
        _anchorNext(1, 5);
        EpochStatement memory statement = _true(0, 1, 5, 2);
        statement.amount += 1;
        court.proveBadProof(statement, _sign(signerKey, statement), _siblings(0, 1, 5, 2));
        _refusedAsTheCourtWould(
            address(court),
            _badProofCall(statement, _siblings(0, 1, 5, 2)),
            "",
            abi.encodeWithSelector(INotaryPaymaster.AlreadyProven.selector, court.statementDigest(statement))
        );

        EpochStatement memory honest = _true(0, 1, 5, 2);
        bytes32 evidence =
            court.proveEquivocation(honest, _sign(signerKey, honest), statement, _sign(signerKey, statement));
        _refusedAsTheCourtWould(
            address(court),
            _equivocationCall(statement, honest),
            "",
            abi.encodeWithSelector(INotaryPaymaster.AlreadyProven.selector, evidence)
        );
    }

    function test_AProofAgainstAnotherLedgerOtherSiblingsOrAnEpochNotAnchoredIsNotSponsored() public {
        _anchorNext(1, 5);
        EpochStatement memory foreign = _true(0, 1, 5, 2);
        foreign.ledger = stranger;
        _refusedAsTheCourtWould(
            address(court),
            _badProofCall(foreign, _siblings(0, 1, 5, 2)),
            "",
            abi.encodeWithSelector(INotaryPaymaster.WrongLedger.selector, address(ledger), stranger)
        );

        EpochStatement memory statement = _true(0, 1, 5, 2);
        statement.amount += 1;
        SumNode[] memory others = _siblings(0, 1, 5, 3);
        _refusedAsTheCourtWould(
            address(court),
            _badProofCall(statement, others),
            "",
            abi.encodeWithSelector(
                INotaryPaymaster.SiblingsMismatch.selector, statement.siblingsHash, court.siblingsHash(others)
            )
        );

        EpochStatement memory unanchored = _true(1, 1, 5, 2);
        _refusedAsTheCourtWould(
            address(court),
            _badProofCall(unanchored, _siblings(1, 1, 5, 2)),
            "",
            abi.encodeWithSelector(INotaryPaymaster.EpochNotAnchored.selector, uint64(1))
        );
    }

    function test_StatementsThatDoNotContradictEachOtherAreNotSponsored() public {
        _anchorNext(1, 5);
        _anchorNext(2, 5);
        EpochStatement memory first = _true(0, 1, 5, 2);
        EpochStatement memory second = _true(1, 2, 5, 2);
        _refusedAsTheCourtWould(
            address(court),
            _equivocationCall(first, second),
            "",
            abi.encodeWithSelector(INotaryPaymaster.DifferentEpochs.selector, uint64(0), uint64(1))
        );

        _refusedAsTheCourtWould(
            address(court),
            _equivocationCall(first, first),
            "",
            abi.encodeWithSelector(INotaryPaymaster.NoConflict.selector)
        );

        EpochStatement memory neighbour = _true(0, 1, 5, 3);
        _refusedAsTheCourtWould(
            address(court),
            _equivocationCall(first, neighbour),
            "",
            abi.encodeWithSelector(INotaryPaymaster.NoConflict.selector)
        );
    }

    function test_AnEpochAnotherSignerAnswersForIsNotSponsored() public {
        _anchorNext(1, 5);
        (address successor,) = makeAddrAndKey("successor-signer");
        uint64 from = ledger.nextEpoch();
        vm.prank(admin);
        court.scheduleSigner(successor, from);
        vm.warp(vm.getBlockTimestamp() + SIGNER_DELAY);
        court.activateSigner();
        _anchorNext(2, 5);

        EpochStatement memory statement = _true(from, 2, 5, 2);
        statement.amount += 1;
        _refusedAsTheCourtWould(
            address(court),
            _badProofCall(statement, _siblings(from, 2, 5, 2)),
            "",
            abi.encodeWithSelector(INotaryPaymaster.AnotherSignersEpoch.selector, from, successor)
        );

        Registration memory registration = _registration(from);
        bytes memory call = abi.encodeCall(
            withholding.requestStatement, (registration, _signRegistration(signerKey, registration), from)
        );
        _refusedAsTheCourtWould(
            address(withholding),
            call,
            "",
            abi.encodeWithSelector(INotaryPaymaster.AnotherSignersEpoch.selector, from, successor)
        );
    }

    function test_TheSponsorshipViewAnswersAsTheValidationDoes() public {
        (uint64 epoch,) = _anchorWithAccount(1, 5, 2);
        PackedUserOperation memory op = _op(holder, address(withholding), _requestCall(epoch), "");

        (uint48 validAfter, uint48 validUntil) = paymaster.sponsorship(holder, op.callData, "");
        (uint48 validatedAfter, uint48 validatedUntil) = _validate(op);

        assertEq(validAfter, validatedAfter);
        assertEq(validUntil, validatedUntil);

        vm.expectRevert(abi.encodeWithSelector(INotaryPaymaster.NotTheHolder.selector, holder, stranger));
        paymaster.sponsorship(stranger, op.callData, "");
    }

    function _refusedAsTheCourtWould(
        address target,
        bytes memory call,
        bytes memory paymasterData,
        bytes memory reason
    ) private {
        uint256 snapshot = vm.snapshotState();
        vm.prank(holder);
        (bool taken,) = target.call(call);
        vm.revertToState(snapshot);
        assertFalse(taken, "the court would have taken the call");

        PackedUserOperation memory op = _signedBy(holderKey, _op(holder, target, call, paymasterData));
        uint256 before = paymaster.getDeposit();
        _expectPaymasterRefusal(reason);
        _handle(op);
        assertEq(paymaster.getDeposit(), before);
    }

    function _badProofCall(
        EpochStatement memory statement,
        SumNode[] memory siblings
    ) private view returns (bytes memory) {
        return abi.encodeCall(court.proveBadProof, (statement, _sign(signerKey, statement), siblings));
    }

    function _equivocationCall(
        EpochStatement memory first,
        EpochStatement memory second
    ) private view returns (bytes memory) {
        return abi.encodeCall(
            court.proveEquivocation, (first, _sign(signerKey, first), second, _sign(signerKey, second))
        );
    }
}
