// SPDX-License-Identifier: MIT
pragma solidity 0.8.37;

import {OpenRequest} from "../../src/withholding/OpenRequest.sol";
import {RequestState} from "../../src/withholding/RequestState.sol";
import {Silence} from "../../src/withholding/Silence.sol";
import {StatementRequest} from "../../src/withholding/StatementRequest.sol";
import {WithholdingFixture} from "../support/WithholdingFixture.sol";
import {WithholdingHandler} from "./WithholdingHandler.sol";

contract WithholdingCourtInvariantTest is WithholdingFixture {
    WithholdingHandler private handler;

    function setUp() public override {
        super.setUp();
        handler = new WithholdingHandler(withholding, anchorer, signerKey);
        targetContract(address(handler));
    }

    function invariant_NoActionHadAnOutcomeTheRulesForbid() public view {
        assertEq(handler.unexpected(), 0);
    }

    function invariant_EverySilenceIsCountedOnce() public view {
        assertEq(withholding.totalSilences(), handler.silences());

        uint256 counted = 0;
        for (uint256 i = 0; i < handler.anchoredCount(); i++) {
            counted += withholding.silenceCount(handler.anchoredAt(i));
        }

        assertEq(counted, withholding.totalSilences());
    }

    function invariant_EverySilenceCameAfterItsDeadlineAndEndedItsRequest() public view {
        for (uint256 i = 0; i < handler.anchoredCount(); i++) {
            uint64 epoch = handler.anchoredAt(i);
            Silence[] memory silences = withholding.silencesOf(epoch);
            for (uint256 j = 0; j < silences.length; j++) {
                StatementRequest memory request = withholding.requestOf(silences[j].registration, epoch);
                assertGt(silences[j].recordedAt, silences[j].dueAt);
                assertEq(request.dueAt, silences[j].dueAt);
                assertEq(uint8(request.state), uint8(RequestState.Unanswered));
            }
        }
    }

    function invariant_ARegistrationIsOpenExactlyWhenItsRequestIs() public view {
        for (uint256 h = 0; h < handler.HOLDERS(); h++) {
            bytes32 digest = withholding.registrationDigest(handler.registrationOf(h));
            OpenRequest memory open = withholding.openRequestOf(digest);
            uint256 openCount = 0;
            for (uint256 i = 0; i < handler.anchoredCount(); i++) {
                uint64 epoch = handler.anchoredAt(i);
                if (withholding.requestOf(digest, epoch).state == RequestState.Open) {
                    openCount++;
                    assertTrue(open.active);
                    assertEq(open.epoch, epoch);
                }
            }

            assertEq(openCount, open.active ? 1 : 0);
        }
    }
}
