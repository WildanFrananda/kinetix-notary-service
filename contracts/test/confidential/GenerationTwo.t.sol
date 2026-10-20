// SPDX-License-Identifier: MIT
pragma solidity 0.8.37;

import {Simple7702Account} from "account-abstraction/contracts/accounts/Simple7702Account.sol";
import {BaseAccount} from "account-abstraction/contracts/core/BaseAccount.sol";
import {EntryPoint} from "account-abstraction/contracts/core/EntryPoint.sol";
import {IEntryPoint} from "account-abstraction/contracts/interfaces/IEntryPoint.sol";

import {EpochHeader} from "../../src/EpochHeader.sol";
import {SumNode} from "../../src/SumNode.sol";
import {ConfidentialTree} from "../../src/confidential/ConfidentialTree.sol";
import {LiabilityLedgerV2} from "../../src/confidential/LiabilityLedgerV2.sol";
import {PlonkVerifier} from "../../src/confidential/PlonkVerifier.sol";
import {EpochStatement} from "../../src/court/EpochStatement.sol";
import {IStatementCourt} from "../../src/court/IStatementCourt.sol";
import {StatementCourt} from "../../src/court/StatementCourt.sol";
import {StatementHashing} from "../../src/court/StatementHashing.sol";
import {Generation} from "../../src/directory/Generation.sol";
import {INotaryPaymaster} from "../../src/paymaster/INotaryPaymaster.sol";
import {NotaryPaymaster} from "../../src/paymaster/NotaryPaymaster.sol";
import {AccountCommitments} from "../../src/withholding/AccountCommitments.sol";
import {IWithholdingCourt} from "../../src/withholding/IWithholdingCourt.sol";
import {Registration} from "../../src/withholding/Registration.sol";
import {RequestState} from "../../src/withholding/RequestState.sol";
import {WithholdingCourt} from "../../src/withholding/WithholdingCourt.sol";
import {ConfidentialTreeBuilder} from "../support/ConfidentialTreeBuilder.sol";
import {DirectoryFixture} from "../support/DirectoryFixture.sol";
import {LeafSet} from "../support/LeafSet.sol";
import {ScriptedVerifier} from "../support/ScriptedVerifier.sol";

contract GenerationTwoTest is DirectoryFixture {
    uint8 private constant DEPTH = 3;
    uint256 private constant COUNT = 5;
    uint256 private constant INDEX = 2;

    ScriptedVerifier private verifier;
    LiabilityLedgerV2 private next;
    StatementCourt private nextCourt;
    WithholdingCourt private nextWithholding;

    function setUp() public override {
        super.setUp();
        _anchorInto(ledger, 1);
        _anchorInto(ledger, 2);
        _seal(ledger);

        verifier = new ScriptedVerifier();
        next = new LiabilityLedgerV2(
            admin,
            ADMIN_DELAY,
            anchorer,
            ledger.nextEpoch(),
            ledger.head(),
            ledger.lastAsOf(),
            PlonkVerifier(address(verifier)),
            DEPTH
        );
        nextCourt = new StatementCourt(admin, ADMIN_DELAY, next, signer, SIGNER_DELAY);
        nextWithholding = new WithholdingCourt(nextCourt, WINDOW, PERIOD);
    }

    function test_TheDirectoryRegistersGenerationTwoAtTheSealedHead() public {
        vm.prank(admin);
        uint256 index = directory.registerSuccessor(next, nextCourt, nextWithholding);

        Generation memory current = directory.current();
        assertEq(index, 1);
        assertEq(current.ledger, address(next));
        assertEq(current.court, address(nextCourt));
        assertEq(current.withholding, address(nextWithholding));
        assertEq(current.scheme, 2);
        assertEq(current.genesisEpoch, 2);
        assertEq(current.genesisDigest, ledger.head());

        _anchorWithAccount(10);
        (uint256 old,) = directory.generationOf(1);
        (uint256 later,) = directory.generationOf(2);
        assertEq(old, 0);
        assertEq(later, 1);
    }

    function test_TheUnchangedCourtRecordsAFalseStatementAboutAGenerationTwoEpoch() public {
        (uint64 number, bytes32 commitment, uint64 amount, SumNode[] memory siblings) = _anchorWithAccount(10);
        EpochStatement memory lie = _statement(number, commitment, amount + 1, siblings);
        bytes memory signature = _signFor(lie);

        nextCourt.proveBadProof(lie, signature, siblings);

        assertEq(nextCourt.faultCount(number), 1);
    }

    function test_TheUnchangedCourtRefusesToFaultATrueStatement() public {
        (uint64 number, bytes32 commitment, uint64 amount, SumNode[] memory siblings) = _anchorWithAccount(10);
        EpochStatement memory truth = _statement(number, commitment, amount, siblings);
        bytes memory signature = _signFor(truth);

        vm.expectRevert(abi.encodeWithSelector(IStatementCourt.ProofHolds.selector, number));
        nextCourt.proveBadProof(truth, signature, siblings);
    }

    function test_AStatementCarryingASumIsAFaultInSchemeTwo() public {
        (uint64 number, bytes32 commitment, uint64 amount, SumNode[] memory siblings) = _anchorWithAccount(10);
        siblings[0].sum = 1;
        EpochStatement memory summed = _statement(number, commitment, amount, siblings);
        bytes memory signature = _signFor(summed);

        nextCourt.proveBadProof(summed, signature, siblings);

        assertEq(nextCourt.faultCount(number), 1);
    }

    function test_TheUnchangedWithholdingCourtTakesAGenerationTwoAnswer() public {
        Registration memory registration = _registration(2);
        (uint64 number, bytes32 commitment, uint64 amount, SumNode[] memory siblings) = _anchorWithAccount(20);
        bytes32 digest = _requestFrom(registration, number);

        vm.expectRevert(abi.encodeWithSelector(IWithholdingCourt.AnswerDoesNotHold.selector, number));
        nextWithholding.answer(registration, number, _salt(number), amount + 1, uint64(INDEX), siblings);

        bytes32 answered = nextWithholding.answer(registration, number, _salt(number), amount, uint64(INDEX), siblings);

        assertEq(answered, commitment);
        assertEq(uint8(nextWithholding.requestOf(digest, number).state), uint8(RequestState.Answered));
    }

    function test_AnUnansweredGenerationTwoRequestBecomesASilence() public {
        Registration memory registration = _registration(2);
        (uint64 number,,,) = _anchorWithAccount(20);
        _requestFrom(registration, number);
        vm.warp(vm.getBlockTimestamp() + PERIOD + 1);

        nextWithholding.recordUnanswered(registration, number);

        assertEq(nextWithholding.silenceCount(number), 1);
    }

    function test_AGenerationTwoPaymasterSponsorsOnlyWhatItsCourtWouldTake() public {
        address canonical = 0x4337084D9E255Ff0702461CF8895CE9E3b5Ff108;
        vm.etch(canonical, address(new EntryPoint()).code);
        NotaryPaymaster paymaster = new NotaryPaymaster(
            admin,
            ADMIN_DELAY,
            IEntryPoint(canonical),
            address(new Simple7702Account()),
            nextCourt,
            nextWithholding,
            0.001 ether,
            1 days
        );
        (uint64 number, bytes32 commitment, uint64 amount, SumNode[] memory siblings) = _anchorWithAccount(10);
        EpochStatement memory lie = _statement(number, commitment, amount + 1, siblings);
        EpochStatement memory truth = _statement(number, commitment, amount, siblings);
        bytes memory lieCall = abi.encodeCall(nextCourt.proveBadProof, (lie, _signFor(lie), siblings));
        bytes memory truthCall = abi.encodeCall(nextCourt.proveBadProof, (truth, _signFor(truth), siblings));

        (uint48 validAfter, uint48 validUntil) =
            paymaster.sponsorship(holder, abi.encodeCall(BaseAccount.execute, (address(nextCourt), 0, lieCall)), "");
        assertEq(validAfter, 0);
        assertEq(validUntil, 0);

        vm.expectRevert(abi.encodeWithSelector(IStatementCourt.ProofHolds.selector, number));
        paymaster.sponsorship(holder, abi.encodeCall(BaseAccount.execute, (address(nextCourt), 0, truthCall)), "");

        vm.expectRevert(abi.encodeWithSelector(INotaryPaymaster.TargetNotSponsored.selector, address(court)));
        paymaster.sponsorship(holder, abi.encodeCall(BaseAccount.execute, (address(court), 0, lieCall)), "");
    }

    function _requestFrom(Registration memory registration, uint64 number) private returns (bytes32) {
        (uint8 v, bytes32 r, bytes32 s) = vm.sign(signerKey, nextWithholding.registrationDigest(registration));
        vm.prank(registration.holder);
        return nextWithholding.requestStatement(registration, abi.encodePacked(r, s, v), number);
    }

    function _anchorWithAccount(uint256 seed) private returns (uint64, bytes32, uint64, SumNode[] memory) {
        uint64 number = next.nextEpoch();
        bytes32[] memory commitments = LeafSet.commitments(seed, COUNT);
        uint64[] memory amounts = LeafSet.amounts(seed, COUNT);
        commitments[INDEX] = AccountCommitments.commitment(account, _salt(number));
        uint256[][] memory levels = ConfidentialTreeBuilder.build(number, DEPTH, commitments, amounts);

        bytes32 treeRoot = bytes32(levels[DEPTH][0]);
        bytes32 totalCommitment = keccak256(abi.encode(seed, "total")) >> 8;
        bytes memory proof = abi.encode(seed, "proof");
        uint256[] memory inputs = new uint256[](4);
        inputs[0] = number;
        inputs[1] = uint256(treeRoot);
        inputs[2] = uint256(totalCommitment);
        inputs[3] = COUNT;
        verifier.accept(proof, inputs);
        bytes32 root = next.proveEpoch(number, treeRoot, totalCommitment, uint64(COUNT), proof);

        EpochHeader memory header = EpochHeader({
            number: number,
            previousDigest: next.head(),
            root: root,
            total: 0,
            leafCount: uint64(COUNT),
            asOf: _now() - 1
        });
        vm.prank(anchorer);
        next.anchor(header);
        vm.warp(vm.getBlockTimestamp() + 1 hours);

        assertEq(root, ConfidentialTree.bindingRoot(treeRoot, totalCommitment));
        return (number, commitments[INDEX], amounts[INDEX], ConfidentialTreeBuilder.siblings(levels, INDEX));
    }

    function _statement(
        uint64 number,
        bytes32 commitment,
        uint64 amount,
        SumNode[] memory siblings
    ) private view returns (EpochStatement memory) {
        return EpochStatement({
            ledger: address(next),
            epoch: number,
            commitment: commitment,
            amount: amount,
            index: uint64(INDEX),
            siblingsHash: StatementHashing.siblingsHash(siblings)
        });
    }

    function _signFor(EpochStatement memory statement) private view returns (bytes memory) {
        (uint8 v, bytes32 r, bytes32 s) = vm.sign(signerKey, nextCourt.statementDigest(statement));
        return abi.encodePacked(r, s, v);
    }
}
