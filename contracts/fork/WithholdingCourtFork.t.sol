// SPDX-License-Identifier: MIT
pragma solidity 0.8.37;

import {Test} from "forge-std/src/Test.sol";

import {Epoch} from "../src/Epoch.sol";
import {ILiabilityLedger} from "../src/ILiabilityLedger.sol";
import {SumNode} from "../src/SumNode.sol";
import {SumTree} from "../src/SumTree.sol";
import {IStatementCourt} from "../src/court/IStatementCourt.sol";
import {StatementCourt} from "../src/court/StatementCourt.sol";
import {IWithholdingCourt} from "../src/withholding/IWithholdingCourt.sol";
import {Registration} from "../src/withholding/Registration.sol";
import {RequestState} from "../src/withholding/RequestState.sol";
import {WithholdingCourt} from "../src/withholding/WithholdingCourt.sol";

contract WithholdingCourtForkTest is Test {
    ILiabilityLedger private ledger;
    IStatementCourt private liveCourt;
    address private holder;

    function setUp() public {
        vm.createSelectFork(vm.envString("KINETIX_NOTARY_FORK_RPC_URL"), vm.envUint("KINETIX_NOTARY_FORK_BLOCK"));
        ledger = ILiabilityLedger(vm.envAddress("KINETIX_NOTARY_FORK_LEDGER"));
        liveCourt = IStatementCourt(vm.envAddress("KINETIX_NOTARY_FORK_COURT"));
        holder = makeAddr("fork-holder");
    }

    function test_StandsBesideTheDeployedCourtAndItsLedger() public {
        WithholdingCourt withholding = new WithholdingCourt(liveCourt, 7 days, 24 hours);

        assertEq(block.chainid, 84_532);
        assertEq(address(withholding.court()), address(liveCourt));
        assertEq(address(withholding.ledger()), address(ledger));
        assertEq(address(liveCourt.ledger()), address(ledger));
    }

    function test_TheDeployedCourtsSignerIsTheOnlyOneThatRegisters() public {
        WithholdingCourt withholding = new WithholdingCourt(liveCourt, 7 days, 24 hours);
        uint64 epoch = ledger.nextEpoch() - 1;
        Registration memory registration =
            Registration({holder: holder, account: keccak256("fork account"), fromEpoch: epoch});
        (address impostor, uint256 impostorKey) = makeAddrAndKey("fork-impostor");
        bytes memory signature = _sign(withholding, impostorKey, registration);
        address liveSigner = liveCourt.signerOf(epoch);

        assertTrue(liveSigner != address(0));
        vm.expectRevert(
            abi.encodeWithSelector(IWithholdingCourt.NotSignedBySigner.selector, epoch, liveSigner, impostor)
        );
        vm.prank(holder);
        withholding.requestStatement(registration, signature, epoch);
    }

    function test_AnAccountTheRealTreeLacksEndsInSilence() public {
        (address signer, uint256 signerKey) = makeAddrAndKey("fork-statement-signer");
        StatementCourt court = new StatementCourt(makeAddr("fork-admin"), 3 days, ledger, signer, 7 days);
        WithholdingCourt withholding = new WithholdingCourt(court, 7 days, 24 hours);
        uint64 epoch = ledger.nextEpoch() - 1;
        Epoch memory anchored = ledger.epochAt(epoch);
        Registration memory registration =
            Registration({holder: holder, account: keccak256("fork account"), fromEpoch: 0});
        bytes memory signature = _sign(withholding, signerKey, registration);
        SumNode[] memory siblings = new SumNode[](SumTree.depthOf(anchored.leafCount));

        vm.prank(holder);
        bytes32 digest = withholding.requestStatement(registration, signature, epoch);

        vm.expectRevert(abi.encodeWithSelector(IWithholdingCourt.AnswerDoesNotHold.selector, epoch));
        withholding.answer(registration, epoch, keccak256("any salt"), anchored.total, 0, siblings);

        vm.warp(vm.getBlockTimestamp() + 24 hours + 1);
        withholding.recordUnanswered(registration, epoch);

        assertEq(uint8(withholding.requestOf(digest, epoch).state), uint8(RequestState.Unanswered));
        assertEq(withholding.silenceCount(epoch), 1);
    }

    function test_TheWindowOverTheFirstLiveEpochClosesAWeekAfterItsAnchoring() public {
        (address signer, uint256 signerKey) = makeAddrAndKey("fork-statement-signer");
        StatementCourt court = new StatementCourt(makeAddr("fork-admin"), 3 days, ledger, signer, 7 days);
        WithholdingCourt withholding = new WithholdingCourt(court, 7 days, 24 hours);
        uint256 closesAt = uint256(ledger.epochAt(0).anchoredAt) + 7 days;
        Registration memory registration =
            Registration({holder: holder, account: keccak256("fork account"), fromEpoch: 0});
        bytes memory signature = _sign(withholding, signerKey, registration);
        vm.warp(closesAt + 1);

        vm.expectRevert(abi.encodeWithSelector(IWithholdingCourt.RequestWindowClosed.selector, uint64(0), closesAt));
        vm.prank(holder);
        withholding.requestStatement(registration, signature, 0);
    }

    function _sign(
        WithholdingCourt withholding,
        uint256 key,
        Registration memory registration
    ) private view returns (bytes memory) {
        (uint8 v, bytes32 r, bytes32 s) = vm.sign(key, withholding.registrationDigest(registration));
        return abi.encodePacked(r, s, v);
    }
}
