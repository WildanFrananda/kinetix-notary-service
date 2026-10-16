// SPDX-License-Identifier: MIT
pragma solidity 0.8.37;

import {IAccessControl} from "@openzeppelin-contracts/access/IAccessControl.sol";
import {Simple7702Account} from "account-abstraction/contracts/accounts/Simple7702Account.sol";
import {BaseAccount} from "account-abstraction/contracts/core/BaseAccount.sol";
import {IEntryPoint} from "account-abstraction/contracts/interfaces/IEntryPoint.sol";
import {IPaymaster} from "account-abstraction/contracts/interfaces/IPaymaster.sol";
import {PackedUserOperation} from "account-abstraction/contracts/interfaces/PackedUserOperation.sol";

import {EpochStatement} from "../../src/court/EpochStatement.sol";
import {IStatementCourt} from "../../src/court/IStatementCourt.sol";
import {StatementCourt} from "../../src/court/StatementCourt.sol";
import {INotaryPaymaster} from "../../src/paymaster/INotaryPaymaster.sol";
import {NotaryPaymaster} from "../../src/paymaster/NotaryPaymaster.sol";
import {IWithholdingCourt} from "../../src/withholding/IWithholdingCourt.sol";
import {Registration} from "../../src/withholding/Registration.sol";
import {RequestState} from "../../src/withholding/RequestState.sol";
import {WithholdingCourt} from "../../src/withholding/WithholdingCourt.sol";
import {LyingCourt} from "../support/LyingCourt.sol";
import {LyingWithholding} from "../support/LyingWithholding.sol";
import {PaymasterFixture} from "../support/PaymasterFixture.sol";

contract NotaryPaymasterTest is PaymasterFixture {
    function test_StandsBesideTheCourtsAndTakesTheirSignerAndDomains() public view {
        assertEq(address(paymaster.entryPoint()), CANONICAL_ENTRY_POINT);
        assertEq(paymaster.account(), address(delegate));
        assertEq(address(paymaster.court()), address(court));
        assertEq(address(paymaster.withholding()), address(withholding));
        assertEq(paymaster.statementSigner(), signer);
        assertEq(paymaster.maxCostPerOperation(), CEILING);
        assertEq(paymaster.unstakeDelay(), UNSTAKE_DELAY);
        assertEq(paymaster.getDeposit(), 1 ether);
        assertEq(paymaster.getDepositInfo().deposit, 1 ether);
        assertFalse(paymaster.getDepositInfo().staked);
        assertEq(paymaster.defaultAdmin(), admin);
        assertTrue(paymaster.supportsInterface(type(INotaryPaymaster).interfaceId));
    }

    function test_AHolderWithNoEtherRequestsAStatement() public {
        (uint64 epoch,) = _anchorWithAccount(1, 5, 2);
        uint256 before = paymaster.getDeposit();

        _handle(_signedBy(holderKey, _op(holder, address(withholding), _requestCall(epoch), "")));

        bytes32 digest = withholding.registrationDigest(_registration(0));
        assertEq(uint8(withholding.requestOf(digest, epoch).state), uint8(RequestState.Open));
        assertEq(holder.balance, 0);
        assertLt(paymaster.getDeposit(), before);
        assertGt(bundler.balance, 0);
    }

    function test_AHolderRecordsAnUnansweredRequestWithTheRegistrationSignatureBesideIt() public {
        (uint64 epoch,) = _anchorWithAccount(1, 5, 2);
        _handle(_signedBy(holderKey, _op(holder, address(withholding), _requestCall(epoch), "")));
        bytes32 digest = withholding.registrationDigest(_registration(0));
        vm.warp(withholding.requestOf(digest, epoch).dueAt + 1);
        bytes memory call = abi.encodeCall(withholding.recordUnanswered, (_registration(0), epoch));
        bytes memory registrationSignature = _signRegistration(signerKey, _registration(0));

        _handle(_signedBy(holderKey, _op(holder, address(withholding), call, registrationSignature)));

        assertEq(withholding.silenceCount(epoch), 1);
        assertEq(withholding.silencesOf(epoch)[0].recorder, holder);
    }

    function test_AHolderProvesABadProof() public {
        _anchorNext(1, 5);
        EpochStatement memory statement = _true(0, 1, 5, 2);
        statement.amount += 1;
        bytes memory call =
            abi.encodeCall(court.proveBadProof, (statement, _sign(signerKey, statement), _siblings(0, 1, 5, 2)));

        _handle(_signedBy(holderKey, _op(holder, address(court), call, "")));

        assertEq(court.faultCount(0), 1);
        assertEq(court.faultsOf(0)[0].prover, holder);
    }

    function test_AHolderProvesAnEquivocation() public {
        _anchorNext(1, 5);
        EpochStatement memory first = _true(0, 1, 5, 2);
        EpochStatement memory second = _true(0, 1, 5, 2);
        second.amount += 7;
        bytes memory call =
            abi.encodeCall(court.proveEquivocation, (first, _sign(signerKey, first), second, _sign(signerKey, second)));

        _handle(_signedBy(holderKey, _op(holder, address(court), call, "")));

        assertEq(court.faultCount(0), 1);
    }

    function test_ASenderDelegatedToAnyOtherAccountIsNotSponsored() public {
        (uint64 epoch,) = _anchorWithAccount(1, 5, 2);
        Simple7702Account lookalike = new Simple7702Account();
        vm.signAndAttachDelegation(address(lookalike), holderKey);
        PackedUserOperation memory op = _signedBy(holderKey, _op(holder, address(withholding), _requestCall(epoch), ""));

        _expectPaymasterRefusal(abi.encodeWithSelector(INotaryPaymaster.NotDelegated.selector, holder));
        _handle(op);
    }

    function test_OnlyTheTwoCourtsAreSponsored() public {
        bytes memory call = abi.encodeCall(ledger.head, ());
        PackedUserOperation memory op = _signedBy(holderKey, _op(holder, address(ledger), call, ""));

        _expectPaymasterRefusal(abi.encodeWithSelector(INotaryPaymaster.TargetNotSponsored.selector, address(ledger)));
        _handle(op);
    }

    function test_OnlyTheHoldersCallsOnTheCourtsAreSponsored() public {
        bytes memory answer = abi.encodeCall(
            withholding.answer, (_registration(0), uint64(0), bytes32(0), uint64(0), uint64(0), _siblings(0, 1, 5, 2))
        );
        PackedUserOperation memory op = _signedBy(holderKey, _op(holder, address(withholding), answer, ""));
        _expectPaymasterRefusal(
            abi.encodeWithSelector(
                INotaryPaymaster.CallNotSponsored.selector, address(withholding), IWithholdingCourt.answer.selector
            )
        );
        _handle(op);

        bytes memory schedule = abi.encodeCall(court.scheduleSigner, (stranger, uint64(9)));
        op = _signedBy(holderKey, _op(holder, address(court), schedule, ""));
        _expectPaymasterRefusal(
            abi.encodeWithSelector(
                INotaryPaymaster.CallNotSponsored.selector, address(court), IStatementCourt.scheduleSigner.selector
            )
        );
        _handle(op);
    }

    function test_OnlyASingleExecuteWithoutValueIsSponsored() public {
        (uint64 epoch,) = _anchorWithAccount(1, 5, 2);
        PackedUserOperation memory op = _op(holder, address(withholding), _requestCall(epoch), "");
        op.callData = abi.encodeCall(BaseAccount.execute, (address(withholding), 1, _requestCall(epoch)));
        op = _signedBy(holderKey, op);
        _expectPaymasterRefusal(abi.encodeWithSelector(INotaryPaymaster.ValueNotZero.selector, 1));
        _handle(op);

        op = _op(holder, address(withholding), _requestCall(epoch), "");
        op.callData = abi.encodeWithSelector(bytes4(keccak256("executeBatch((address,uint256,bytes)[])")), "");
        op = _signedBy(holderKey, op);
        _expectPaymasterRefusal(
            abi.encodeWithSelector(
                INotaryPaymaster.NotAnExecute.selector, bytes4(keccak256("executeBatch((address,uint256,bytes)[])"))
            )
        );
        _handle(op);
    }

    function test_ARegistrationForAnotherHolderIsNotSponsored() public {
        (uint64 epoch,) = _anchorWithAccount(1, 5, 2);
        Registration memory theirs = _registration(0);
        theirs.holder = stranger;
        bytes memory call =
            abi.encodeCall(withholding.requestStatement, (theirs, _signRegistration(signerKey, theirs), epoch));
        PackedUserOperation memory op = _signedBy(holderKey, _op(holder, address(withholding), call, ""));

        _expectPaymasterRefusal(abi.encodeWithSelector(INotaryPaymaster.NotTheHolder.selector, stranger, holder));
        _handle(op);
    }

    function test_ARegistrationTheNotaryNeverSignedIsNotSponsored() public {
        (uint64 epoch,) = _anchorWithAccount(1, 5, 2);
        (address forger, uint256 forgerKey) = makeAddrAndKey("forger");
        bytes memory call = abi.encodeCall(
            withholding.requestStatement, (_registration(0), _signRegistration(forgerKey, _registration(0)), epoch)
        );
        PackedUserOperation memory op = _signedBy(holderKey, _op(holder, address(withholding), call, ""));

        _expectPaymasterRefusal(abi.encodeWithSelector(INotaryPaymaster.NotSignedByTheNotary.selector, forger));
        _handle(op);
    }

    function test_AnUnansweredRecordWithoutTheNotarysSignatureIsNotSponsored() public {
        bytes memory call = abi.encodeCall(withholding.recordUnanswered, (_registration(0), uint64(0)));
        PackedUserOperation memory op = _signedBy(holderKey, _op(holder, address(withholding), call, ""));

        _expectPaymasterRefusal(abi.encodeWithSelector(INotaryPaymaster.NotSignedByTheNotary.selector, address(0)));
        _handle(op);
    }

    function test_AStatementTheNotaryNeverSignedIsNotSponsored() public {
        _anchorNext(1, 5);
        (address forger, uint256 forgerKey) = makeAddrAndKey("forger");
        EpochStatement memory statement = _true(0, 1, 5, 2);
        statement.amount += 1;
        bytes memory call =
            abi.encodeCall(court.proveBadProof, (statement, _sign(forgerKey, statement), _siblings(0, 1, 5, 2)));
        PackedUserOperation memory op = _signedBy(holderKey, _op(holder, address(court), call, ""));
        _expectPaymasterRefusal(abi.encodeWithSelector(INotaryPaymaster.NotSignedByTheNotary.selector, forger));
        _handle(op);

        EpochStatement memory second = _true(0, 1, 5, 2);
        second.amount += 7;
        call = abi.encodeCall(
            court.proveEquivocation, (second, _sign(signerKey, second), statement, _sign(forgerKey, statement))
        );
        op = _signedBy(holderKey, _op(holder, address(court), call, ""));
        _expectPaymasterRefusal(abi.encodeWithSelector(INotaryPaymaster.NotSignedByTheNotary.selector, forger));
        _handle(op);
    }

    function test_AnOperationThatMayCostMoreThanTheCeilingIsNotSponsored() public {
        (uint64 epoch,) = _anchorWithAccount(1, 5, 2);
        PackedUserOperation memory op = _op(holder, address(withholding), _requestCall(epoch), "");
        op.gasFees = _pack(FEE, 10 gwei);
        op = _signedBy(holderKey, op);
        uint256 maxCost = (VERIFICATION_GAS + CALL_GAS + PAYMASTER_GAS + PRE_VERIFICATION_GAS) * 10 gwei;

        _expectPaymasterRefusal(abi.encodeWithSelector(INotaryPaymaster.CostTooHigh.selector, maxCost, CEILING));
        _handle(op);
    }

    function test_TheCeilingIsTheMostAnOperationMayCost() public {
        (uint64 epoch,) = _anchorWithAccount(1, 5, 2);
        PackedUserOperation memory op = _op(holder, address(withholding), _requestCall(epoch), "");

        vm.prank(CANONICAL_ENTRY_POINT);
        paymaster.validatePaymasterUserOp(op, bytes32(0), CEILING);

        vm.expectRevert(abi.encodeWithSelector(INotaryPaymaster.CostTooHigh.selector, CEILING + 1, CEILING));
        vm.prank(CANONICAL_ENTRY_POINT);
        paymaster.validatePaymasterUserOp(op, bytes32(0), CEILING + 1);
    }

    function test_OnlyTheEntryPointAsksAndNothingIsLeftForAPostOp() public {
        PackedUserOperation memory op = _op(holder, address(withholding), "", "");

        vm.expectRevert(abi.encodeWithSelector(INotaryPaymaster.NotTheEntryPoint.selector, stranger));
        vm.prank(stranger);
        paymaster.validatePaymasterUserOp(op, bytes32(0), 0);

        vm.expectRevert(INotaryPaymaster.NoPostOp.selector);
        paymaster.postOp(IPaymaster.PostOpMode.opSucceeded, "", 0, 0);
    }

    function test_AnyoneDepositsAndOnlyTheAdminWithdraws() public {
        vm.deal(stranger, 1 ether);
        vm.prank(stranger);
        entryPoint.depositTo{value: 0.5 ether}(address(paymaster));
        assertEq(paymaster.getDeposit(), 1.5 ether);

        vm.expectRevert(
            abi.encodeWithSelector(
                IAccessControl.AccessControlUnauthorizedAccount.selector, stranger, paymaster.DEFAULT_ADMIN_ROLE()
            )
        );
        vm.prank(stranger);
        paymaster.withdrawTo(payable(stranger), 1 ether);

        address payable treasury = payable(makeAddr("treasury"));
        vm.prank(admin);
        paymaster.withdrawTo(treasury, 1.5 ether);
        assertEq(paymaster.getDeposit(), 0);
        assertEq(treasury.balance, 1.5 ether);
    }

    function test_RefusesToStandBesideAnythingButTheCourtsAndTheirSigner() public {
        vm.expectRevert(abi.encodeWithSelector(INotaryPaymaster.NoContract.selector, stranger));
        new NotaryPaymaster(
            admin, ADMIN_DELAY, IEntryPoint(stranger), address(delegate), court, withholding, CEILING, UNSTAKE_DELAY
        );

        vm.expectRevert(abi.encodeWithSelector(INotaryPaymaster.NoContract.selector, address(0)));
        new NotaryPaymaster(admin, ADMIN_DELAY, entryPoint, address(0), court, withholding, CEILING, UNSTAKE_DELAY);

        vm.expectRevert(abi.encodeWithSelector(INotaryPaymaster.NotAStatementCourt.selector, address(ledger)));
        new NotaryPaymaster(
            admin,
            ADMIN_DELAY,
            entryPoint,
            address(delegate),
            IStatementCourt(address(ledger)),
            withholding,
            CEILING,
            UNSTAKE_DELAY
        );

        vm.expectRevert(abi.encodeWithSelector(INotaryPaymaster.NotAWithholdingCourt.selector, address(court)));
        new NotaryPaymaster(
            admin,
            ADMIN_DELAY,
            entryPoint,
            address(delegate),
            court,
            IWithholdingCourt(address(court)),
            CEILING,
            UNSTAKE_DELAY
        );

        StatementCourt otherCourt = new StatementCourt(admin, ADMIN_DELAY, ledger, signer, SIGNER_DELAY);
        vm.expectRevert(
            abi.encodeWithSelector(
                INotaryPaymaster.WithholdingOfAnotherCourt.selector, address(withholding), address(otherCourt)
            )
        );
        new NotaryPaymaster(
            admin, ADMIN_DELAY, entryPoint, address(delegate), otherCourt, withholding, CEILING, UNSTAKE_DELAY
        );

        vm.expectRevert(INotaryPaymaster.ZeroCostCeiling.selector);
        new NotaryPaymaster(admin, ADMIN_DELAY, entryPoint, address(delegate), court, withholding, 0, UNSTAKE_DELAY);
    }

    function test_RefusesCourtsWithoutASignerOrWithAnotherDomain() public {
        LyingCourt silent = new LyingCourt(address(0));
        LyingWithholding besideSilent = new LyingWithholding(IStatementCourt(address(silent)));
        vm.expectRevert(abi.encodeWithSelector(INotaryPaymaster.NoStatementSigner.selector, address(silent)));
        new NotaryPaymaster(
            admin,
            ADMIN_DELAY,
            entryPoint,
            address(delegate),
            IStatementCourt(address(silent)),
            IWithholdingCourt(address(besideSilent)),
            CEILING,
            UNSTAKE_DELAY
        );

        LyingCourt foreign = new LyingCourt(signer);
        LyingWithholding besideForeign = new LyingWithholding(IStatementCourt(address(foreign)));
        vm.expectRevert(abi.encodeWithSelector(INotaryPaymaster.DomainMismatch.selector, address(foreign)));
        new NotaryPaymaster(
            admin,
            ADMIN_DELAY,
            entryPoint,
            address(delegate),
            IStatementCourt(address(foreign)),
            IWithholdingCourt(address(besideForeign)),
            CEILING,
            UNSTAKE_DELAY
        );

        LyingWithholding besideReal = new LyingWithholding(court);
        vm.expectRevert(abi.encodeWithSelector(INotaryPaymaster.DomainMismatch.selector, address(besideReal)));
        new NotaryPaymaster(
            admin,
            ADMIN_DELAY,
            entryPoint,
            address(delegate),
            court,
            IWithholdingCourt(address(besideReal)),
            CEILING,
            UNSTAKE_DELAY
        );
    }

    function test_TheFirstOperationCarriesTheDelegationMarkerAndStillHashesAsTheEntryPointDoes() public {
        (uint64 epoch,) = _anchorWithAccount(1, 5, 2);
        PackedUserOperation memory op = _op(holder, address(withholding), _requestCall(epoch), "");
        op.initCode = abi.encodePacked(bytes20(bytes2(0x7702)));

        _handle(_signedBy(holderKey, op));

        bytes32 digest = withholding.registrationDigest(_registration(0));
        assertEq(uint8(withholding.requestOf(digest, epoch).state), uint8(RequestState.Open));
    }
}
