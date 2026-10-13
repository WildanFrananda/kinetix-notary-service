// SPDX-License-Identifier: MIT
pragma solidity 0.8.37;

import {Fault} from "../../src/court/Fault.sol";
import {SignerTerm} from "../../src/court/SignerTerm.sol";
import {CourtFixture} from "../support/CourtFixture.sol";
import {CourtHandler} from "./CourtHandler.sol";

contract StatementCourtInvariantTest is CourtFixture {
    CourtHandler private handler;

    function setUp() public override {
        super.setUp();
        handler = new CourtHandler(court, admin, anchorer, signer, signerKey);
        targetContract(address(handler));
    }

    function invariant_NoActionHadAnOutcomeTheRulesForbid() public view {
        assertEq(handler.unexpected(), 0);
    }

    function invariant_EveryProvenFaultIsCountedOnce() public view {
        assertEq(court.totalFaults(), handler.faults());

        uint256 counted = 0;
        for (uint256 i = 0; i < handler.anchoredCount(); i++) {
            counted += court.faultCount(handler.anchoredAt(i));
        }

        assertEq(counted, court.totalFaults());
    }

    function invariant_AnAnchoredEpochNeverChangesSigner() public view {
        for (uint256 i = 0; i < handler.anchoredCount(); i++) {
            uint64 epoch = handler.anchoredAt(i);
            assertEq(court.signerOf(epoch), handler.observedSigner(epoch));
        }
    }

    function invariant_TermsStartStrictlyLater() public view {
        SignerTerm[] memory terms = court.signerTerms();
        assertEq(terms[0].fromEpoch, ledger.genesisEpoch());
        for (uint256 i = 1; i < terms.length; i++) {
            assertGt(terms[i].fromEpoch, terms[i - 1].fromEpoch);
        }
    }

    function invariant_EveryRecordedFaultIsEvidenced() public view {
        for (uint256 i = 0; i < handler.anchoredCount(); i++) {
            Fault[] memory faults = court.faultsOf(handler.anchoredAt(i));
            for (uint256 j = 0; j < faults.length; j++) {
                assertTrue(court.isProven(faults[j].evidence));
                assertLt(handler.anchoredAt(i), ledger.nextEpoch());
            }
        }
    }
}
