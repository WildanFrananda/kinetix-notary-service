// SPDX-License-Identifier: MIT
pragma solidity 0.8.37;

import {IAccessControl} from "@openzeppelin-contracts/access/IAccessControl.sol";
import {IERC165} from "@openzeppelin-contracts/utils/introspection/IERC165.sol";

import {ILiabilityLedger} from "../../src/ILiabilityLedger.sol";
import {SumNode} from "../../src/SumNode.sol";
import {EpochStatement} from "../../src/court/EpochStatement.sol";
import {Fault} from "../../src/court/Fault.sol";
import {FaultKind} from "../../src/court/FaultKind.sol";
import {IStatementCourt} from "../../src/court/IStatementCourt.sol";
import {PendingSigner} from "../../src/court/PendingSigner.sol";
import {SignerTerm} from "../../src/court/SignerTerm.sol";
import {StatementCourt} from "../../src/court/StatementCourt.sol";
import {StatementHashing} from "../../src/court/StatementHashing.sol";
import {CourtFixture} from "../support/CourtFixture.sol";

contract StatementCourtTest is CourtFixture {
    function test_StartsWithTheInitialSignerFromGenesis() public view {
        SignerTerm[] memory terms = court.signerTerms();

        assertEq(terms.length, 1);
        assertEq(terms[0].signer, signer);
        assertEq(terms[0].fromEpoch, 0);
        assertEq(court.signerOf(0), signer);
        assertEq(court.signerOf(type(uint64).max), signer);
        assertEq(address(court.ledger()), address(ledger));
        assertEq(court.signerChangeDelay(), SIGNER_DELAY);
        assertEq(court.defaultAdmin(), admin);
    }

    function test_RefusesAnythingThatIsNotALedger() public {
        vm.expectRevert(abi.encodeWithSelector(IStatementCourt.NotALedger.selector, stranger));
        new StatementCourt(admin, ADMIN_DELAY, _ledgerAt(stranger), signer, SIGNER_DELAY);

        vm.expectRevert(abi.encodeWithSelector(IStatementCourt.NotALedger.selector, address(court)));
        new StatementCourt(admin, ADMIN_DELAY, _ledgerAt(address(court)), signer, SIGNER_DELAY);
    }

    function test_RefusesAZeroSignerOrAZeroDelay() public {
        vm.expectRevert(IStatementCourt.ZeroSigner.selector);
        new StatementCourt(admin, ADMIN_DELAY, ledger, address(0), SIGNER_DELAY);

        vm.expectRevert(IStatementCourt.ZeroDelay.selector);
        new StatementCourt(admin, ADMIN_DELAY, ledger, signer, 0);
    }

    function test_AnnouncesItsInterfaces() public view {
        assertTrue(court.supportsInterface(type(IStatementCourt).interfaceId));
        assertTrue(court.supportsInterface(type(IAccessControl).interfaceId));
        assertTrue(court.supportsInterface(type(IERC165).interfaceId));
        assertFalse(court.supportsInterface(0xffffffff));
    }

    function test_TheDomainNamesThisCourtOnThisChain() public view {
        (, string memory name, string memory version, uint256 chainId, address verifyingContract,,) =
            court.eip712Domain();

        assertEq(name, "Kinetix Statement Court");
        assertEq(version, "1");
        assertEq(chainId, block.chainid);
        assertEq(verifyingContract, address(court));
    }

    function test_TheDigestIsEIP712OverTheStatement() public {
        _anchorNext(1, 5);
        EpochStatement memory statement = _true(0, 1, 5, 2);
        bytes32 domain = keccak256(
            abi.encode(
                keccak256("EIP712Domain(string name,string version,uint256 chainId,address verifyingContract)"),
                keccak256("Kinetix Statement Court"),
                keccak256("1"),
                block.chainid,
                address(court)
            )
        );
        bytes32 structHash = keccak256(
            abi.encode(
                keccak256(
                    "EpochStatement(address ledger,uint64 epoch,bytes32 commitment,uint64 amount,uint64 index,bytes32 siblingsHash)"
                ),
                statement.ledger,
                statement.epoch,
                statement.commitment,
                statement.amount,
                statement.index,
                statement.siblingsHash
            )
        );

        assertEq(court.statementDigest(statement), keccak256(abi.encodePacked("\x19\x01", domain, structHash)));
    }

    function test_TheSiblingsHashIsEachHashThenEachSumPacked() public view {
        SumNode[] memory siblings = _siblings(0, 1, 5, 3);
        bytes memory packed;
        for (uint256 i = 0; i < siblings.length; i++) {
            packed = abi.encodePacked(packed, siblings[i].hash, siblings[i].sum);
        }

        assertEq(siblings.length, 3);
        assertEq(court.siblingsHash(siblings), keccak256(packed));
        assertEq(court.siblingsHash(new SumNode[](0)), keccak256(""));
    }

    function test_RecoversTheSignerAndNothingFromABrokenSignature() public {
        _anchorNext(1, 5);
        EpochStatement memory statement = _true(0, 1, 5, 2);
        bytes memory signature = _sign(signerKey, statement);

        assertEq(court.recoverSigner(statement, signature), signer);
        assertEq(court.recoverSigner(statement, abi.encodePacked(bytes32(0), bytes32(0))), address(0));
        assertEq(court.recoverSigner(statement, _highS(signature)), address(0));
    }

    function test_ATrueStatementIsNoFault() public {
        _anchorNext(1, 5);
        EpochStatement memory statement = _true(0, 1, 5, 2);
        bytes memory signature = _sign(signerKey, statement);

        vm.expectRevert(abi.encodeWithSelector(IStatementCourt.ProofHolds.selector, uint64(0)));
        court.proveBadProof(statement, signature, _siblings(0, 1, 5, 2));
    }

    function test_ASignedStatementWithTheWrongAmountIsAFault() public {
        _anchorNext(1, 5);
        EpochStatement memory statement = _true(0, 1, 5, 2);
        statement.amount += 1;
        bytes memory signature = _sign(signerKey, statement);
        bytes32 digest = court.statementDigest(statement);

        vm.expectEmit(address(court));
        emit IStatementCourt.FaultProven(0, FaultKind.BadProof, digest, stranger);
        vm.prank(stranger);
        bytes32 evidence = court.proveBadProof(statement, signature, _siblings(0, 1, 5, 2));

        assertEq(evidence, digest);
        assertTrue(court.isProven(digest));
        assertEq(court.faultCount(0), 1);
        assertEq(court.totalFaults(), 1);
        Fault[] memory faults = court.faultsOf(0);
        assertEq(uint8(faults[0].kind), uint8(FaultKind.BadProof));
        assertEq(faults[0].evidence, digest);
        assertEq(faults[0].prover, stranger);
        assertEq(faults[0].provenAt, _now());
    }

    function test_EveryWayAStatementCanLieIsAFault() public {
        _anchorNext(1, 5);
        SumNode[] memory siblings = _siblings(0, 1, 5, 2);

        EpochStatement memory wrongIndex = _true(0, 1, 5, 2);
        wrongIndex.index = 3;
        court.proveBadProof(wrongIndex, _sign(signerKey, wrongIndex), siblings);

        EpochStatement memory wrongCommitment = _true(0, 1, 5, 2);
        wrongCommitment.commitment = keccak256("an account that is not there");
        court.proveBadProof(wrongCommitment, _sign(signerKey, wrongCommitment), siblings);

        SumNode[] memory shortened = new SumNode[](2);
        shortened[0] = siblings[0];
        shortened[1] = siblings[1];
        EpochStatement memory shortPath = _true(0, 1, 5, 2);
        shortPath.siblingsHash = StatementHashing.siblingsHash(shortened);
        court.proveBadProof(shortPath, _sign(signerKey, shortPath), shortened);

        EpochStatement memory outside = _true(0, 1, 5, 2);
        outside.index = 9;
        court.proveBadProof(outside, _sign(signerKey, outside), siblings);

        assertEq(court.faultCount(0), 4);
        assertEq(court.totalFaults(), 4);
    }

    function test_AFaultIsRecordedOnce() public {
        _anchorNext(1, 5);
        EpochStatement memory statement = _true(0, 1, 5, 2);
        statement.amount += 1;
        bytes memory signature = _sign(signerKey, statement);
        court.proveBadProof(statement, signature, _siblings(0, 1, 5, 2));

        vm.expectRevert(
            abi.encodeWithSelector(IStatementCourt.AlreadyProven.selector, court.statementDigest(statement))
        );
        court.proveBadProof(statement, signature, _siblings(0, 1, 5, 2));
    }

    function test_TheSiblingsMustBeTheOnesThatWereSigned() public {
        _anchorNext(1, 5);
        EpochStatement memory statement = _true(0, 1, 5, 2);
        statement.amount += 1;
        SumNode[] memory other = _siblings(0, 1, 5, 3);
        bytes memory signature = _sign(signerKey, statement);

        vm.expectRevert(
            abi.encodeWithSelector(
                IStatementCourt.SiblingsMismatch.selector, statement.siblingsHash, StatementHashing.siblingsHash(other)
            )
        );
        court.proveBadProof(statement, signature, other);
    }

    function test_OnlyTheEpochsSignerCanBeConvicted() public {
        _anchorNext(1, 5);
        EpochStatement memory statement = _true(0, 1, 5, 2);
        statement.amount += 1;
        (address forger, uint256 forgerKey) = makeAddrAndKey("forger");
        bytes memory forged = _sign(forgerKey, statement);

        vm.expectRevert(abi.encodeWithSelector(IStatementCourt.NotSignedBySigner.selector, 0, signer, forger));
        court.proveBadProof(statement, forged, _siblings(0, 1, 5, 2));

        vm.expectRevert(abi.encodeWithSelector(IStatementCourt.NotSignedBySigner.selector, 0, signer, address(0)));
        court.proveBadProof(statement, hex"1234", _siblings(0, 1, 5, 2));
    }

    function test_AStatementAboutAnotherLedgerProvesNothingHere() public {
        _anchorNext(1, 5);
        EpochStatement memory statement = _true(0, 1, 5, 2);
        statement.ledger = stranger;
        bytes memory signature = _sign(signerKey, statement);

        vm.expectRevert(abi.encodeWithSelector(IStatementCourt.WrongLedger.selector, address(ledger), stranger));
        court.proveBadProof(statement, signature, _siblings(0, 1, 5, 2));
    }

    function test_AStatementAboutAnEpochNotYetAnchoredProvesNothing() public {
        _anchorNext(1, 5);
        EpochStatement memory statement = _true(1, 1, 5, 2);
        bytes memory signature = _sign(signerKey, statement);

        vm.expectRevert(abi.encodeWithSelector(IStatementCourt.EpochNotAnchored.selector, uint64(1)));
        court.proveBadProof(statement, signature, _siblings(1, 1, 5, 2));
    }

    function test_OneLeafWithTwoAmountsIsEquivocation() public {
        _anchorNext(1, 5);
        EpochStatement memory first = _true(0, 1, 5, 2);
        EpochStatement memory second = _true(0, 1, 5, 2);
        second.amount += 7;

        bytes32 evidence = court.proveEquivocation(first, _sign(signerKey, first), second, _sign(signerKey, second));

        assertEq(evidence, _pair(court.statementDigest(first), court.statementDigest(second)));
        Fault[] memory faults = court.faultsOf(0);
        assertEq(faults.length, 1);
        assertEq(uint8(faults[0].kind), uint8(FaultKind.Equivocation));
    }

    function test_EveryContradictionBetweenTwoStatementsIsEquivocation() public {
        _anchorNext(1, 5);
        EpochStatement memory honest = _true(0, 1, 5, 2);
        bytes memory honestSignature = _sign(signerKey, honest);

        EpochStatement memory twoIndexes = _true(0, 1, 5, 2);
        twoIndexes.index = 4;
        court.proveEquivocation(honest, honestSignature, twoIndexes, _sign(signerKey, twoIndexes));

        EpochStatement memory twoLeaves = _true(0, 1, 5, 2);
        twoLeaves.commitment = keccak256("another account");
        court.proveEquivocation(honest, honestSignature, twoLeaves, _sign(signerKey, twoLeaves));

        EpochStatement memory twoPaths = _true(0, 1, 5, 2);
        twoPaths.siblingsHash = keccak256("another path");
        court.proveEquivocation(honest, honestSignature, twoPaths, _sign(signerKey, twoPaths));

        assertEq(court.faultCount(0), 3);
    }

    function test_TwoDifferentLeavesDoNotConflict() public {
        _anchorNext(1, 5);
        EpochStatement memory first = _true(0, 1, 5, 1);
        EpochStatement memory second = _true(0, 1, 5, 3);
        bytes memory firstSignature = _sign(signerKey, first);
        bytes memory secondSignature = _sign(signerKey, second);

        vm.expectRevert(IStatementCourt.NoConflict.selector);
        court.proveEquivocation(first, firstSignature, second, secondSignature);

        vm.expectRevert(IStatementCourt.NoConflict.selector);
        court.proveEquivocation(first, firstSignature, first, firstSignature);
    }

    function test_EquivocationCountsOnceWhicheverWayRound() public {
        _anchorNext(1, 5);
        EpochStatement memory first = _true(0, 1, 5, 2);
        EpochStatement memory second = _true(0, 1, 5, 2);
        second.amount += 1;
        bytes memory firstSignature = _sign(signerKey, first);
        bytes memory secondSignature = _sign(signerKey, second);
        bytes32 evidence = court.proveEquivocation(first, firstSignature, second, secondSignature);

        vm.expectRevert(abi.encodeWithSelector(IStatementCourt.AlreadyProven.selector, evidence));
        court.proveEquivocation(second, secondSignature, first, firstSignature);
    }

    function test_EquivocationNeedsOneEpoch() public {
        _anchorNext(1, 5);
        _anchorNext(2, 5);
        EpochStatement memory first = _true(0, 1, 5, 2);
        EpochStatement memory second = _true(1, 2, 5, 2);
        bytes memory firstSignature = _sign(signerKey, first);
        bytes memory secondSignature = _sign(signerKey, second);

        vm.expectRevert(abi.encodeWithSelector(IStatementCourt.DifferentEpochs.selector, uint64(0), uint64(1)));
        court.proveEquivocation(first, firstSignature, second, secondSignature);
    }

    function test_OnlyTheAdminSchedulesASigner() public {
        vm.expectRevert(
            abi.encodeWithSelector(
                IAccessControl.AccessControlUnauthorizedAccount.selector, stranger, court.DEFAULT_ADMIN_ROLE()
            )
        );
        vm.prank(stranger);
        court.scheduleSigner(stranger, 5);
    }

    function test_ASignerTakesOverAfterTheDelayFromItsEpoch() public {
        _anchorNext(1, 5);
        (address next,) = makeAddrAndKey("next-signer");

        vm.expectEmit(address(court));
        emit IStatementCourt.SignerScheduled(next, 3, _now() + SIGNER_DELAY);
        vm.prank(admin);
        court.scheduleSigner(next, 3);

        PendingSigner memory pending = court.pendingSigner();
        assertEq(pending.signer, next);
        assertEq(pending.fromEpoch, 3);
        assertEq(court.signerOf(3), signer);

        vm.expectRevert(
            abi.encodeWithSelector(IStatementCourt.SignerChangeTooEarly.selector, pending.activatesAt, _now())
        );
        court.activateSigner();

        vm.warp(pending.activatesAt);
        vm.expectEmit(address(court));
        emit IStatementCourt.SignerActivated(next, 3);
        vm.prank(stranger);
        court.activateSigner();

        assertEq(court.signerOf(0), signer);
        assertEq(court.signerOf(2), signer);
        assertEq(court.signerOf(3), next);
        assertEq(court.signerOf(type(uint64).max), next);
        assertEq(court.pendingSigner().signer, address(0));
        assertEq(court.signerTerms().length, 2);
    }

    function test_ARotationCannotDisownAStatementAlreadySigned() public {
        _anchorNext(1, 5);
        EpochStatement memory statement = _true(0, 1, 5, 2);
        statement.amount += 1;
        bytes memory signature = _sign(signerKey, statement);
        (address next,) = makeAddrAndKey("next-signer");

        vm.prank(admin);
        court.scheduleSigner(next, 1);
        vm.warp(court.pendingSigner().activatesAt);
        court.activateSigner();

        court.proveBadProof(statement, signature, _siblings(0, 1, 5, 2));
        assertEq(court.faultCount(0), 1);
    }

    function test_ASignerCannotStartAtAnEpochThatAlreadyExists() public {
        _anchorNext(1, 5);
        _anchorNext(2, 5);
        (address next,) = makeAddrAndKey("next-signer");

        vm.expectRevert(abi.encodeWithSelector(IStatementCourt.FromEpochTooLow.selector, uint64(1), uint64(2)));
        vm.prank(admin);
        court.scheduleSigner(next, 1);
    }

    function test_ASignerScheduledTooCloseGoesStale() public {
        (address next,) = makeAddrAndKey("next-signer");
        vm.prank(admin);
        court.scheduleSigner(next, 1);
        _anchorNext(1, 5);
        _anchorNext(2, 5);
        vm.warp(court.pendingSigner().activatesAt);

        vm.expectRevert(abi.encodeWithSelector(IStatementCourt.SignerChangeStale.selector, uint64(1), uint64(2)));
        court.activateSigner();
    }

    function test_OneChangeAtATimeAndTheAdminMayCancelIt() public {
        (address next,) = makeAddrAndKey("next-signer");
        (address other,) = makeAddrAndKey("other-signer");
        vm.startPrank(admin);
        court.scheduleSigner(next, 4);

        vm.expectRevert(abi.encodeWithSelector(IStatementCourt.SignerChangePending.selector, next, uint64(4)));
        court.scheduleSigner(other, 5);

        vm.expectEmit(address(court));
        emit IStatementCourt.SignerCancelled(next, 4);
        court.cancelSigner();

        vm.expectRevert(IStatementCourt.NoSignerChangePending.selector);
        court.cancelSigner();
        vm.stopPrank();

        vm.expectRevert(IStatementCourt.NoSignerChangePending.selector);
        court.activateSigner();
    }

    function test_TheScheduleRefusesNonsense() public {
        vm.startPrank(admin);

        vm.expectRevert(IStatementCourt.ZeroSigner.selector);
        court.scheduleSigner(address(0), 4);

        vm.expectRevert(abi.encodeWithSelector(IStatementCourt.SameSigner.selector, signer));
        court.scheduleSigner(signer, 4);
        vm.stopPrank();
    }

    function test_ALaterTermMustStartAfterTheLastOne() public {
        (address next,) = makeAddrAndKey("next-signer");
        (address other,) = makeAddrAndKey("other-signer");
        vm.prank(admin);
        court.scheduleSigner(next, 10);
        vm.warp(court.pendingSigner().activatesAt);
        court.activateSigner();

        vm.expectRevert(abi.encodeWithSelector(IStatementCourt.FromEpochTooLow.selector, uint64(10), uint64(11)));
        vm.prank(admin);
        court.scheduleSigner(other, 10);
    }

    function test_NoTermCanFollowOneThatStartsAtTheLastEpoch() public {
        (address next,) = makeAddrAndKey("next-signer");
        (address other,) = makeAddrAndKey("other-signer");
        vm.prank(admin);
        court.scheduleSigner(next, type(uint64).max);
        vm.warp(court.pendingSigner().activatesAt);
        court.activateSigner();

        vm.expectRevert(
            abi.encodeWithSelector(IStatementCourt.FromEpochTooLow.selector, type(uint64).max, type(uint64).max)
        );
        vm.prank(admin);
        court.scheduleSigner(other, type(uint64).max);
    }

    function _ledgerAt(address where) private pure returns (ILiabilityLedger) {
        return ILiabilityLedger(where);
    }

    function _pair(bytes32 a, bytes32 b) private pure returns (bytes32) {
        if (a < b) {
            return keccak256(abi.encodePacked(a, b));
        }

        return keccak256(abi.encodePacked(b, a));
    }

    function _highS(bytes memory signature) private pure returns (bytes memory) {
        bytes32 r;
        bytes32 s;
        uint8 v;
        assembly ("memory-safe") {
            r := mload(add(signature, 0x20))
            s := mload(add(signature, 0x40))
            v := byte(0, mload(add(signature, 0x60)))
        }
        uint256 order = 0xFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFEBAAEDCE6AF48A03BBFD25E8CD0364141;
        return abi.encodePacked(r, bytes32(order - uint256(s)), v == 27 ? uint8(28) : uint8(27));
    }
}
