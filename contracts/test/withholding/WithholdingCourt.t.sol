// SPDX-License-Identifier: MIT
pragma solidity 0.8.37;

import {IERC165} from "@openzeppelin-contracts/utils/introspection/IERC165.sol";

import {SumNode} from "../../src/SumNode.sol";
import {IStatementCourt} from "../../src/court/IStatementCourt.sol";
import {IWithholdingCourt} from "../../src/withholding/IWithholdingCourt.sol";
import {OpenRequest} from "../../src/withholding/OpenRequest.sol";
import {Registration} from "../../src/withholding/Registration.sol";
import {RequestState} from "../../src/withholding/RequestState.sol";
import {Silence} from "../../src/withholding/Silence.sol";
import {StatementRequest} from "../../src/withholding/StatementRequest.sol";
import {WithholdingCourt} from "../../src/withholding/WithholdingCourt.sol";
import {TreeBuilder} from "../support/TreeBuilder.sol";
import {WithholdingFixture} from "../support/WithholdingFixture.sol";

contract WithholdingCourtTest is WithholdingFixture {
    function test_StartsWithItsCourtItsLedgerAndItsDurations() public view {
        assertEq(address(withholding.court()), address(court));
        assertEq(address(withholding.ledger()), address(ledger));
        assertEq(withholding.requestWindow(), WINDOW);
        assertEq(withholding.answerPeriod(), PERIOD);
        assertEq(withholding.totalSilences(), 0);
    }

    function test_RefusesAnythingThatIsNotACourt() public {
        vm.expectRevert(abi.encodeWithSelector(IWithholdingCourt.NotACourt.selector, stranger));
        new WithholdingCourt(IStatementCourt(stranger), WINDOW, PERIOD);

        vm.expectRevert(abi.encodeWithSelector(IWithholdingCourt.NotACourt.selector, address(ledger)));
        new WithholdingCourt(IStatementCourt(address(ledger)), WINDOW, PERIOD);
    }

    function test_RefusesAZeroOrAnOverlongDuration() public {
        vm.expectRevert(abi.encodeWithSelector(IWithholdingCourt.DurationOutOfRange.selector, uint64(0), 365 days));
        new WithholdingCourt(court, 0, PERIOD);

        vm.expectRevert(abi.encodeWithSelector(IWithholdingCourt.DurationOutOfRange.selector, uint64(0), 365 days));
        new WithholdingCourt(court, WINDOW, 0);

        vm.expectRevert(
            abi.encodeWithSelector(IWithholdingCourt.DurationOutOfRange.selector, uint64(365 days + 1), 365 days)
        );
        new WithholdingCourt(court, 365 days + 1, PERIOD);

        WithholdingCourt longest = new WithholdingCourt(court, 365 days, 365 days);
        assertEq(longest.requestWindow(), 365 days);
    }

    function test_AnnouncesItsInterfaces() public view {
        assertTrue(withholding.supportsInterface(type(IWithholdingCourt).interfaceId));
        assertTrue(withholding.supportsInterface(type(IERC165).interfaceId));
        assertFalse(withholding.supportsInterface(type(IStatementCourt).interfaceId));
        assertFalse(withholding.supportsInterface(0xffffffff));
    }

    function test_TheDomainNamesThisCourtOnThisChain() public view {
        (, string memory name, string memory version, uint256 chainId, address verifyingContract,,) =
            withholding.eip712Domain();

        assertEq(name, "Kinetix Withholding Court");
        assertEq(version, "1");
        assertEq(chainId, block.chainid);
        assertEq(verifyingContract, address(withholding));
    }

    function test_TheRegistrationDigestIsEIP712OverTheRegistration() public view {
        Registration memory registration = _registration(3);
        bytes32 domain = keccak256(
            abi.encode(
                keccak256("EIP712Domain(string name,string version,uint256 chainId,address verifyingContract)"),
                keccak256("Kinetix Withholding Court"),
                keccak256("1"),
                block.chainid,
                address(withholding)
            )
        );
        bytes32 structHash = keccak256(
            abi.encode(
                keccak256("Registration(address holder,bytes32 account,uint64 fromEpoch)"),
                registration.holder,
                registration.account,
                registration.fromEpoch
            )
        );

        assertEq(
            withholding.registrationDigest(registration), keccak256(abi.encodePacked("\x19\x01", domain, structHash))
        );
    }

    function test_TheCommitmentIsTheTagThenTheAccountThenTheSalt() public view {
        bytes32 salt = _salt(4);

        assertEq(withholding.commitmentOf(account, salt), keccak256(abi.encodePacked(bytes1(0x06), account, salt)));
    }

    function test_RecoversTheSignerAndNothingFromABrokenSignature() public view {
        Registration memory registration = _registration(0);
        bytes memory signature = _signRegistration(signerKey, registration);

        assertEq(withholding.recoverSigner(registration, signature), signer);
        assertEq(withholding.recoverSigner(registration, abi.encodePacked(bytes32(0), bytes32(0))), address(0));
        assertEq(withholding.recoverSigner(registration, hex"00"), address(0));
    }

    function test_AHolderRequestsAStatement() public {
        (uint64 epoch,) = _anchorWithAccount(1, 5, 2);
        Registration memory registration = _registration(0);
        bytes memory signature = _signRegistration(signerKey, registration);
        bytes32 digest = withholding.registrationDigest(registration);
        uint64 dueAt = _now() + PERIOD;

        vm.expectEmit(address(withholding));
        emit IWithholdingCourt.StatementRequested(digest, account, epoch, holder, dueAt);
        vm.prank(holder);
        bytes32 returned = withholding.requestStatement(registration, signature, epoch);

        StatementRequest memory request = withholding.requestOf(digest, epoch);
        OpenRequest memory open = withholding.openRequestOf(digest);
        assertEq(returned, digest);
        assertEq(uint8(request.state), uint8(RequestState.Open));
        assertEq(request.dueAt, dueAt);
        assertTrue(open.active);
        assertEq(open.epoch, epoch);
    }

    function test_OnlyTheHolderMayRequest() public {
        (uint64 epoch,) = _anchorWithAccount(1, 5, 2);
        Registration memory registration = _registration(0);
        bytes memory signature = _signRegistration(signerKey, registration);

        vm.expectRevert(abi.encodeWithSelector(IWithholdingCourt.NotTheHolder.selector, holder, stranger));
        vm.prank(stranger);
        withholding.requestStatement(registration, signature, epoch);
    }

    function test_RefusesARegistrationTheSignerDidNotSign() public {
        (uint64 epoch,) = _anchorWithAccount(1, 5, 2);
        Registration memory registration = _registration(0);
        (address forger, uint256 forgerKey) = makeAddrAndKey("forger");
        bytes memory forged = _signRegistration(forgerKey, registration);
        bytes memory byTheHolder = _signRegistration(holderKey, registration);

        vm.expectRevert(abi.encodeWithSelector(IWithholdingCourt.NotSignedBySigner.selector, uint64(0), signer, forger));
        vm.prank(holder);
        withholding.requestStatement(registration, forged, epoch);

        vm.expectRevert(abi.encodeWithSelector(IWithholdingCourt.NotSignedBySigner.selector, uint64(0), signer, holder));
        vm.prank(holder);
        withholding.requestStatement(registration, byTheHolder, epoch);

        vm.expectRevert(
            abi.encodeWithSelector(IWithholdingCourt.NotSignedBySigner.selector, uint64(0), signer, address(0))
        );
        vm.prank(holder);
        withholding.requestStatement(registration, hex"00", epoch);
    }

    function test_ARegistrationWithAnyFieldChangedIsNoLongerSigned() public {
        (uint64 epoch,) = _anchorWithAccount(1, 5, 2);
        Registration memory registration = _registration(0);
        bytes memory signature = _signRegistration(signerKey, registration);
        registration.account = keccak256("another account");
        address recovered = withholding.recoverSigner(registration, signature);

        assertTrue(recovered != signer);
        vm.expectRevert(
            abi.encodeWithSelector(IWithholdingCourt.NotSignedBySigner.selector, uint64(0), signer, recovered)
        );
        vm.prank(holder);
        withholding.requestStatement(registration, signature, epoch);
    }

    function test_ARegistrationIsSignedByTheSignerOfItsFirstEpoch() public {
        _anchorWithAccount(1, 5, 2);
        (address next, uint256 nextKey) = makeAddrAndKey("next-signer");
        vm.prank(admin);
        court.scheduleSigner(next, 2);
        vm.warp(vm.getBlockTimestamp() + SIGNER_DELAY);
        court.activateSigner();
        _anchorWithAccount(2, 5, 1);
        (uint64 epoch,) = _anchorWithAccount(3, 5, 4);

        Registration memory early = _registration(0);
        bytes memory earlySignature = _signRegistration(signerKey, early);
        vm.prank(holder);
        withholding.requestStatement(early, earlySignature, epoch);

        Registration memory late = Registration({holder: stranger, account: account, fromEpoch: 2});
        bytes memory byTheOldSigner = _signRegistration(signerKey, late);
        bytes memory byTheNewSigner = _signRegistration(nextKey, late);
        vm.expectRevert(abi.encodeWithSelector(IWithholdingCourt.NotSignedBySigner.selector, uint64(2), next, signer));
        vm.prank(stranger);
        withholding.requestStatement(late, byTheOldSigner, epoch);

        vm.prank(stranger);
        withholding.requestStatement(late, byTheNewSigner, epoch);
    }

    function test_RefusesAnEpochBeforeTheRegistration() public {
        (uint64 epoch,) = _anchorWithAccount(1, 5, 2);
        Registration memory registration = _registration(epoch + 1);
        bytes memory signature = _signRegistration(signerKey, registration);

        vm.expectRevert(abi.encodeWithSelector(IWithholdingCourt.BeforeRegistration.selector, epoch, epoch + 1));
        vm.prank(holder);
        withholding.requestStatement(registration, signature, epoch);
    }

    function test_RefusesAnEpochThatIsNotAnchored() public {
        _anchorWithAccount(1, 5, 2);
        Registration memory registration = _registration(0);
        bytes memory signature = _signRegistration(signerKey, registration);

        vm.expectRevert(abi.encodeWithSelector(IWithholdingCourt.EpochNotAnchored.selector, uint64(1)));
        vm.prank(holder);
        withholding.requestStatement(registration, signature, 1);
    }

    function test_TheWindowClosesSevenDaysAfterAnchoring() public {
        uint64 anchoredAt = _now();
        (uint64 epoch,) = _anchorWithAccount(1, 5, 2);
        Registration memory registration = _registration(0);
        bytes memory signature = _signRegistration(signerKey, registration);
        uint256 closesAt = uint256(anchoredAt) + WINDOW;

        vm.warp(closesAt + 1);
        vm.expectRevert(abi.encodeWithSelector(IWithholdingCourt.RequestWindowClosed.selector, epoch, closesAt));
        vm.prank(holder);
        withholding.requestStatement(registration, signature, epoch);

        vm.warp(closesAt);
        vm.prank(holder);
        withholding.requestStatement(registration, signature, epoch);
    }

    function test_ARegistrationHoldsOneOpenRequestAtATime() public {
        (uint64 first,) = _anchorWithAccount(1, 5, 2);
        (uint64 second,) = _anchorWithAccount(2, 5, 3);
        Registration memory registration = _registration(0);
        bytes32 digest = _request(registration, first);
        bytes memory signature = _signRegistration(signerKey, registration);

        vm.expectRevert(abi.encodeWithSelector(IWithholdingCourt.RequestStillOpen.selector, digest, first));
        vm.prank(holder);
        withholding.requestStatement(registration, signature, second);
    }

    function test_AnAnswerFromAnyoneClosesTheRequest() public {
        (uint64 epoch, uint64 amount) = _anchorWithAccount(1, 5, 2);
        Registration memory registration = _registration(0);
        bytes32 digest = _request(registration, epoch);
        SumNode[] memory siblings = _accountSiblings(epoch, 1, 5, 2);
        bytes32 commitment = withholding.commitmentOf(account, _salt(epoch));

        vm.expectEmit(address(withholding));
        emit IWithholdingCourt.StatementAnswered(digest, account, epoch, commitment, amount, stranger);
        vm.prank(stranger);
        bytes32 returned = withholding.answer(registration, epoch, _salt(epoch), amount, 2, siblings);

        assertEq(returned, commitment);
        assertEq(uint8(withholding.requestOf(digest, epoch).state), uint8(RequestState.Answered));
        assertFalse(withholding.openRequestOf(digest).active);
    }

    function test_AnAnsweredEpochCannotBeRequestedAgainButTheNextCan() public {
        (uint64 first, uint64 amount) = _anchorWithAccount(1, 5, 2);
        (uint64 second,) = _anchorWithAccount(2, 5, 3);
        Registration memory registration = _registration(0);
        bytes32 digest = _request(registration, first);
        withholding.answer(registration, first, _salt(first), amount, 2, _accountSiblings(first, 1, 5, 2));
        bytes memory signature = _signRegistration(signerKey, registration);

        vm.expectRevert(abi.encodeWithSelector(IWithholdingCourt.AlreadyRequested.selector, digest, first));
        vm.prank(holder);
        withholding.requestStatement(registration, signature, first);

        vm.prank(holder);
        withholding.requestStatement(registration, signature, second);
        assertEq(withholding.openRequestOf(digest).epoch, second);
    }

    function test_AnAnswerThatDoesNotHoldIsRefused() public {
        (uint64 epoch, uint64 amount) = _anchorWithAccount(1, 5, 2);
        Registration memory registration = _registration(0);
        _request(registration, epoch);
        SumNode[] memory siblings = _accountSiblings(epoch, 1, 5, 2);

        vm.expectRevert(abi.encodeWithSelector(IWithholdingCourt.AnswerDoesNotHold.selector, epoch));
        withholding.answer(registration, epoch, _salt(epoch + 1), amount, 2, siblings);

        vm.expectRevert(abi.encodeWithSelector(IWithholdingCourt.AnswerDoesNotHold.selector, epoch));
        withholding.answer(registration, epoch, _salt(epoch), amount + 1, 2, siblings);

        vm.expectRevert(abi.encodeWithSelector(IWithholdingCourt.AnswerDoesNotHold.selector, epoch));
        withholding.answer(registration, epoch, _salt(epoch), amount, 3, siblings);
    }

    function test_AnotherLeafOfTheTreeIsNotAnAnswer() public {
        (uint64 epoch,) = _anchorWithAccount(1, 5, 2);
        Registration memory registration = _registration(0);
        _request(registration, epoch);
        (SumNode[][] memory levels,) = _accountTree(epoch, 1, 5, 2);
        SumNode[] memory siblings = TreeBuilder.proof(levels, 4);
        uint64 amount = _amounts(1, 5)[4];

        assertTrue(ledger.verifyInclusion(epoch, _commitments(1, 5)[4], amount, 4, siblings));
        vm.expectRevert(abi.encodeWithSelector(IWithholdingCourt.AnswerDoesNotHold.selector, epoch));
        withholding.answer(registration, epoch, _salt(epoch), amount, 4, siblings);
    }

    function test_AnAnswerNamingAnotherRegistrationFindsNoRequest() public {
        (uint64 epoch, uint64 amount) = _anchorWithAccount(1, 5, 2);
        Registration memory registration = _registration(0);
        _request(registration, epoch);
        Registration memory other = Registration({holder: stranger, account: account, fromEpoch: 0});
        bytes32 otherDigest = withholding.registrationDigest(other);
        SumNode[] memory siblings = _accountSiblings(epoch, 1, 5, 2);

        vm.expectRevert(abi.encodeWithSelector(IWithholdingCourt.NoOpenRequest.selector, otherDigest, epoch));
        withholding.answer(other, epoch, _salt(epoch), amount, 2, siblings);
    }

    function test_TheAnswerIsDueTwentyFourHoursAfterTheRequest() public {
        (uint64 epoch, uint64 amount) = _anchorWithAccount(1, 5, 2);
        Registration memory registration = _registration(0);
        bytes32 digest = _request(registration, epoch);
        uint64 dueAt = withholding.requestOf(digest, epoch).dueAt;
        SumNode[] memory siblings = _accountSiblings(epoch, 1, 5, 2);
        uint256 snapshot = vm.snapshotState();

        vm.warp(dueAt + 1);
        vm.expectRevert(abi.encodeWithSelector(IWithholdingCourt.AnswerTooLate.selector, dueAt, dueAt + 1));
        withholding.answer(registration, epoch, _salt(epoch), amount, 2, siblings);

        vm.revertToState(snapshot);
        vm.warp(dueAt);
        withholding.answer(registration, epoch, _salt(epoch), amount, 2, siblings);
        assertEq(uint8(withholding.requestOf(digest, epoch).state), uint8(RequestState.Answered));
    }

    function test_NoAnswerWithoutARequest() public {
        (uint64 epoch, uint64 amount) = _anchorWithAccount(1, 5, 2);
        Registration memory registration = _registration(0);
        bytes32 digest = withholding.registrationDigest(registration);
        SumNode[] memory siblings = _accountSiblings(epoch, 1, 5, 2);

        vm.expectRevert(abi.encodeWithSelector(IWithholdingCourt.NoOpenRequest.selector, digest, epoch));
        withholding.answer(registration, epoch, _salt(epoch), amount, 2, siblings);
    }

    function test_AnUnansweredRequestBecomesASilence() public {
        (uint64 epoch,) = _anchorWithAccount(1, 5, 2);
        Registration memory registration = _registration(0);
        bytes32 digest = _request(registration, epoch);
        uint64 dueAt = withholding.requestOf(digest, epoch).dueAt;
        vm.warp(dueAt + 1);

        vm.expectEmit(address(withholding));
        emit IWithholdingCourt.SilenceRecorded(digest, account, epoch, stranger);
        vm.prank(stranger);
        bytes32 returned = withholding.recordUnanswered(registration, epoch);

        Silence[] memory silences = withholding.silencesOf(epoch);
        assertEq(returned, digest);
        assertEq(silences.length, 1);
        assertEq(silences[0].registration, digest);
        assertEq(silences[0].account, account);
        assertEq(silences[0].holder, holder);
        assertEq(silences[0].dueAt, dueAt);
        assertEq(silences[0].recorder, stranger);
        assertEq(silences[0].recordedAt, dueAt + 1);
        assertEq(withholding.silenceCount(epoch), 1);
        assertEq(withholding.totalSilences(), 1);
        assertEq(uint8(withholding.requestOf(digest, epoch).state), uint8(RequestState.Unanswered));
        assertFalse(withholding.openRequestOf(digest).active);
    }

    function test_ASilenceWaitsForTheDeadline() public {
        (uint64 epoch,) = _anchorWithAccount(1, 5, 2);
        Registration memory registration = _registration(0);
        bytes32 digest = _request(registration, epoch);
        uint64 dueAt = withholding.requestOf(digest, epoch).dueAt;
        vm.warp(dueAt);

        vm.expectRevert(abi.encodeWithSelector(IWithholdingCourt.NotOverdue.selector, dueAt, dueAt));
        withholding.recordUnanswered(registration, epoch);
    }

    function test_ASilenceIsRecordedOnceAndNeverOverAnAnswer() public {
        (uint64 first, uint64 amount) = _anchorWithAccount(1, 5, 2);
        (uint64 second,) = _anchorWithAccount(2, 5, 3);
        Registration memory registration = _registration(0);
        bytes32 digest = _request(registration, first);
        withholding.answer(registration, first, _salt(first), amount, 2, _accountSiblings(first, 1, 5, 2));
        vm.warp(vm.getBlockTimestamp() + PERIOD + 1);

        vm.expectRevert(abi.encodeWithSelector(IWithholdingCourt.NoOpenRequest.selector, digest, first));
        withholding.recordUnanswered(registration, first);

        _request(registration, second);
        vm.warp(vm.getBlockTimestamp() + PERIOD + 1);
        withholding.recordUnanswered(registration, second);

        vm.expectRevert(abi.encodeWithSelector(IWithholdingCourt.NoOpenRequest.selector, digest, second));
        withholding.recordUnanswered(registration, second);
        assertEq(withholding.totalSilences(), 1);
        assertEq(uint8(withholding.requestOf(digest, first).state), uint8(RequestState.Answered));
        assertEq(uint8(withholding.requestOf(digest, second).state), uint8(RequestState.Unanswered));
    }

    function test_AnAccountLeftOutOfTheTreeCanOnlyEndInSilence() public {
        uint64 epoch = _anchorWithoutAccount(1, 5);
        Registration memory registration = _registration(0);
        bytes32 digest = _request(registration, epoch);

        for (uint64 index = 0; index < 5; index++) {
            SumNode[] memory siblings = _siblings(epoch, 1, 5, index);
            uint64 amount = _amounts(1, 5)[index];
            vm.expectRevert(abi.encodeWithSelector(IWithholdingCourt.AnswerDoesNotHold.selector, epoch));
            withholding.answer(registration, epoch, _salt(epoch), amount, index, siblings);
        }

        vm.warp(vm.getBlockTimestamp() + PERIOD + 1);
        withholding.recordUnanswered(registration, epoch);
        assertEq(uint8(withholding.requestOf(digest, epoch).state), uint8(RequestState.Unanswered));
        assertEq(withholding.totalSilences(), 1);
    }

    function test_TheNextRequestMayFollowASilence() public {
        uint64 first = _anchorWithoutAccount(1, 5);
        (uint64 second,) = _anchorWithAccount(2, 5, 3);
        Registration memory registration = _registration(0);
        bytes32 digest = _request(registration, first);
        vm.warp(vm.getBlockTimestamp() + PERIOD + 1);
        withholding.recordUnanswered(registration, first);

        _request(registration, second);

        assertEq(withholding.openRequestOf(digest).epoch, second);
    }
}
