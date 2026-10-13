// SPDX-License-Identifier: MIT
pragma solidity 0.8.37;

import {Test} from "forge-std/src/Test.sol";

import {Epoch} from "../src/Epoch.sol";
import {ILiabilityLedger} from "../src/ILiabilityLedger.sol";
import {SumNode} from "../src/SumNode.sol";
import {SumTree} from "../src/SumTree.sol";
import {EpochStatement} from "../src/court/EpochStatement.sol";
import {FaultKind} from "../src/court/FaultKind.sol";
import {IStatementCourt} from "../src/court/IStatementCourt.sol";
import {StatementCourt} from "../src/court/StatementCourt.sol";
import {StatementHashing} from "../src/court/StatementHashing.sol";

contract StatementCourtForkTest is Test {
    ILiabilityLedger private ledger;
    StatementCourt private court;
    address private signer;
    uint256 private signerKey;

    function setUp() public {
        vm.createSelectFork(vm.envString("KINETIX_NOTARY_FORK_RPC_URL"), vm.envUint("KINETIX_NOTARY_FORK_BLOCK"));
        ledger = ILiabilityLedger(vm.envAddress("KINETIX_NOTARY_FORK_LEDGER"));
        (signer, signerKey) = makeAddrAndKey("fork-statement-signer");
        court = new StatementCourt(makeAddr("fork-admin"), 3 days, ledger, signer, 7 days);
    }

    function test_TheCourtStandsBesideTheDeployedLedger() public view {
        assertEq(block.chainid, 84_532);
        assertEq(ledger.scheme(), SumTree.SCHEME);
        assertGt(ledger.nextEpoch(), 0);
        assertEq(court.signerOf(0), signer);
        assertEq(address(court.ledger()), address(ledger));
    }

    function test_AForgedStatementAboutEpochZeroIsConvictedByTheRealRoot() public {
        (EpochStatement memory statement, SumNode[] memory siblings) = _forged(keccak256("forged account"));
        bytes memory signature = _sign(statement);

        bytes32 evidence = court.proveBadProof(statement, signature, siblings);

        assertEq(evidence, court.statementDigest(statement));
        assertEq(court.faultCount(0), 1);
        assertEq(uint8(court.faultsOf(0)[0].kind), uint8(FaultKind.BadProof));
    }

    function test_TwoStatementsClaimingOneIndexOfEpochZeroAreEquivocation() public {
        (EpochStatement memory first,) = _forged(keccak256("first account"));
        (EpochStatement memory second,) = _forged(keccak256("second account"));
        bytes memory firstSignature = _sign(first);
        bytes memory secondSignature = _sign(second);

        court.proveEquivocation(first, firstSignature, second, secondSignature);

        assertEq(court.totalFaults(), 1);
    }

    function test_AnEpochTheLedgerHasNotReachedProvesNothing() public {
        (EpochStatement memory statement, SumNode[] memory siblings) = _forged(keccak256("future"));
        statement.epoch = ledger.nextEpoch();
        bytes memory signature = _sign(statement);

        vm.expectRevert(abi.encodeWithSelector(IStatementCourt.EpochNotAnchored.selector, statement.epoch));
        court.proveBadProof(statement, signature, siblings);
    }

    function _forged(bytes32 commitment) private view returns (EpochStatement memory, SumNode[] memory) {
        Epoch memory epoch = ledger.epochAt(0);
        SumNode[] memory siblings = new SumNode[](SumTree.depthOf(epoch.leafCount));
        EpochStatement memory statement = EpochStatement({
            ledger: address(ledger),
            epoch: 0,
            commitment: commitment,
            amount: epoch.total,
            index: 0,
            siblingsHash: StatementHashing.siblingsHash(siblings)
        });
        return (statement, siblings);
    }

    function _sign(EpochStatement memory statement) private view returns (bytes memory) {
        (uint8 v, bytes32 r, bytes32 s) = vm.sign(signerKey, court.statementDigest(statement));
        return abi.encodePacked(r, s, v);
    }
}
