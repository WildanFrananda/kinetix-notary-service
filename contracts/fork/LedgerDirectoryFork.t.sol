// SPDX-License-Identifier: MIT
pragma solidity 0.8.37;

import {Test} from "forge-std/src/Test.sol";

import {EpochHeader} from "../src/EpochHeader.sol";
import {LiabilityLedger} from "../src/LiabilityLedger.sol";
import {SumNode} from "../src/SumNode.sol";
import {StatementCourt} from "../src/court/StatementCourt.sol";
import {Generation} from "../src/directory/Generation.sol";
import {ILedgerDirectory} from "../src/directory/ILedgerDirectory.sol";
import {LedgerDirectory} from "../src/directory/LedgerDirectory.sol";
import {WithholdingCourt} from "../src/withholding/WithholdingCourt.sol";
import {LeafSet} from "../test/support/LeafSet.sol";
import {TreeBuilder} from "../test/support/TreeBuilder.sol";

contract LedgerDirectoryForkTest is Test {
    LiabilityLedger private ledger;
    StatementCourt private court;
    WithholdingCourt private withholding;
    LedgerDirectory private directory;
    address private admin;

    function setUp() public {
        vm.createSelectFork(vm.envString("KINETIX_NOTARY_FORK_RPC_URL"), vm.envUint("KINETIX_NOTARY_FORK_BLOCK"));
        ledger = LiabilityLedger(vm.envAddress("KINETIX_NOTARY_FORK_LEDGER"));
        court = StatementCourt(vm.envAddress("KINETIX_NOTARY_FORK_COURT"));
        withholding = new WithholdingCourt(court, 7 days, 24 hours);
        admin = ledger.defaultAdmin();
        directory = new LedgerDirectory(admin, 3 days, ledger, court, withholding, 7 days);
    }

    function test_GenerationZeroIsTheLiveLedgerAndItsCourt() public view {
        Generation memory zero = directory.current();

        assertEq(block.chainid, 84_532);
        assertEq(zero.chainId, 84_532);
        assertEq(zero.ledger, address(ledger));
        assertEq(zero.court, address(court));
        assertEq(zero.genesisEpoch, ledger.genesisEpoch());
        assertEq(zero.genesisDigest, ledger.genesisDigest());
        assertEq(zero.scheme, ledger.scheme());
    }

    function test_TheLiveLedgerCannotBeSucceededWhileItAnchors() public {
        (LiabilityLedger next, StatementCourt nextCourt, WithholdingCourt nextWithholding) = _successor();

        vm.expectRevert(abi.encodeWithSelector(ILedgerDirectory.PredecessorNotSealed.selector, address(ledger)));
        vm.prank(admin);
        directory.registerSuccessor(next, nextCourt, nextWithholding);
    }

    function test_ASuccessorOfTheSealedLiveLedgerContinuesItsChainExactly() public {
        uint64 last = ledger.nextEpoch() - 1;
        bytes32 head = ledger.head();
        vm.prank(admin);
        ledger.seal(head);
        (LiabilityLedger next, StatementCourt nextCourt, WithholdingCourt nextWithholding) = _successor();

        vm.prank(admin);
        uint256 index = directory.registerSuccessor(next, nextCourt, nextWithholding);

        (uint256 before,) = directory.generationOf(last);
        (uint256 after_, Generation memory found) = directory.generationOf(last + 1);
        assertEq(index, 1);
        assertEq(before, 0);
        assertEq(after_, 1);
        assertEq(found.ledger, address(next));
        assertEq(found.genesisDigest, head);

        EpochHeader memory header = _nextHeader(next, 3);
        assertEq(header.previousDigest, head);
        vm.prank(makeAddr("fork-successor-anchorer"));
        bytes32 digest = next.anchor(header);
        assertEq(next.nextEpoch(), last + 2);
        assertEq(next.epochAt(last + 1).digest, digest);
    }

    function _successor() private returns (LiabilityLedger, StatementCourt, WithholdingCourt) {
        LiabilityLedger next = new LiabilityLedger(
            admin, 3 days, makeAddr("fork-successor-anchorer"), ledger.nextEpoch(), ledger.head(), ledger.lastAsOf()
        );
        StatementCourt nextCourt = new StatementCourt(admin, 3 days, next, makeAddr("fork-successor-signer"), 7 days);
        return (next, nextCourt, new WithholdingCourt(nextCourt, 7 days, 24 hours));
    }

    function _nextHeader(LiabilityLedger target, uint256 count) private view returns (EpochHeader memory) {
        uint64 number = target.nextEpoch();
        SumNode[][] memory levels =
            TreeBuilder.build(number, LeafSet.commitments(number, count), LeafSet.amounts(number, count));
        SumNode memory root = TreeBuilder.root(levels);
        return EpochHeader({
            number: number,
            previousDigest: target.head(),
            root: root.hash,
            total: root.sum,
            leafCount: uint64(count),
            asOf: target.lastAsOf() + 1
        });
    }
}
