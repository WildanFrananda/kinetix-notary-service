// SPDX-License-Identifier: MIT
pragma solidity 0.8.37;

import {EpochHeader} from "../../src/EpochHeader.sol";
import {LiabilityLedger} from "../../src/LiabilityLedger.sol";
import {SumNode} from "../../src/SumNode.sol";
import {StatementCourt} from "../../src/court/StatementCourt.sol";
import {LedgerDirectory} from "../../src/directory/LedgerDirectory.sol";
import {WithholdingCourt} from "../../src/withholding/WithholdingCourt.sol";
import {LeafSet} from "./LeafSet.sol";
import {TreeBuilder} from "./TreeBuilder.sol";
import {WithholdingFixture} from "./WithholdingFixture.sol";

abstract contract DirectoryFixture is WithholdingFixture {
    uint64 internal constant FOREIGN_DELAY = 7 days;
    uint64 internal constant FOREIGN_CHAIN = 11_155_111;

    LedgerDirectory internal directory;

    function setUp() public virtual override {
        super.setUp();
        directory = new LedgerDirectory(admin, ADMIN_DELAY, ledger, court, withholding, FOREIGN_DELAY);
    }

    function _anchorInto(LiabilityLedger target, uint256 seed) internal returns (uint64) {
        uint64 number = target.nextEpoch();
        uint256 count = 1 + seed % 3;
        SumNode[][] memory levels =
            TreeBuilder.build(number, LeafSet.commitments(seed, count), LeafSet.amounts(seed, count));
        EpochHeader memory header = _header(number, target.head(), levels, count, _now() - 1);
        vm.prank(anchorer);
        target.anchor(header);
        vm.warp(vm.getBlockTimestamp() + 1 hours);
        return number;
    }

    function _seal(LiabilityLedger target) internal {
        bytes32 head = target.head();
        vm.prank(admin);
        target.seal(head);
    }

    function _successorOf(LiabilityLedger predecessor)
        internal
        returns (LiabilityLedger, StatementCourt, WithholdingCourt)
    {
        LiabilityLedger next = new LiabilityLedger(
            admin, ADMIN_DELAY, anchorer, predecessor.nextEpoch(), predecessor.head(), predecessor.lastAsOf()
        );
        StatementCourt nextCourt = new StatementCourt(admin, ADMIN_DELAY, next, signer, SIGNER_DELAY);
        return (next, nextCourt, new WithholdingCourt(nextCourt, WINDOW, PERIOD));
    }

    function _succeed(LiabilityLedger predecessor) internal returns (LiabilityLedger) {
        _seal(predecessor);
        (LiabilityLedger next, StatementCourt nextCourt, WithholdingCourt nextWithholding) = _successorOf(predecessor);
        vm.prank(admin);
        directory.registerSuccessor(next, nextCourt, nextWithholding);
        return next;
    }

    function _propose(uint64 chainId) internal {
        vm.prank(admin);
        directory.proposeForeignSuccessor(chainId, address(0xF1), address(0xF2), address(0xF3), 1);
    }
}
