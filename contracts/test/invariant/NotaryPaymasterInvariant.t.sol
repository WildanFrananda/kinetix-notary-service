// SPDX-License-Identifier: MIT
pragma solidity 0.8.37;

import {PaymasterFixture} from "../support/PaymasterFixture.sol";
import {PaymasterHandler} from "./PaymasterHandler.sol";

contract NotaryPaymasterInvariantTest is PaymasterFixture {
    PaymasterHandler private handler;

    function setUp() public override {
        super.setUp();
        (uint64 epoch,) = _anchorWithAccount(1, 5, 2);
        handler = new PaymasterHandler(paymaster, withholding, holderKey, signerKey, _registration(0), epoch, admin);
        targetContract(address(handler));
    }

    function invariant_NoActionHadAnOutcomeTheRulesForbid() public view {
        assertEq(handler.unexpected(), 0);
    }

    function invariant_TheDepositMovesOnlyForSponsoredOperationsAndTheAdmin() public view {
        assertEq(paymaster.getDeposit(), handler.deposited() - handler.withdrawn() - handler.spent());
    }

    function invariant_TheStakeIsWhatTheAdminLocked() public view {
        assertEq(paymaster.getDepositInfo().stake, handler.staked());
    }

    function invariant_TheHolderNeverPays() public view {
        assertEq(holder.balance, 0);
    }
}
