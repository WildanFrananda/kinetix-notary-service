// SPDX-License-Identifier: MIT
pragma solidity 0.8.37;

import {EpochStatement} from "../../src/court/EpochStatement.sol";
import {IStatementCourt} from "../../src/court/IStatementCourt.sol";
import {StatementHashing} from "../../src/court/StatementHashing.sol";
import {CourtFixture} from "../support/CourtFixture.sol";

contract StatementCourtFuzzTest is CourtFixture {
    uint256 private constant CURVE_ORDER = 0xFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFEBAAEDCE6AF48A03BBFD25E8CD0364141;

    function testFuzz_ATrueStatementIsNeverConvicted(uint256 seed, uint256 rawCount, uint256 rawIndex) public {
        uint256 count = bound(rawCount, 1, 40);
        uint256 index = bound(rawIndex, 0, count - 1);
        _anchorNext(seed, count);
        EpochStatement memory statement = _true(0, seed, count, index);
        bytes memory signature = _sign(signerKey, statement);

        vm.expectRevert(abi.encodeWithSelector(IStatementCourt.ProofHolds.selector, uint64(0)));
        court.proveBadProof(statement, signature, _siblings(0, seed, count, index));
    }

    function testFuzz_AnyWrongAmountSignedIsConvicted(
        uint256 seed,
        uint256 rawCount,
        uint256 rawIndex,
        uint64 rawDelta
    ) public {
        uint256 count = bound(rawCount, 1, 40);
        uint256 index = bound(rawIndex, 0, count - 1);
        _anchorNext(seed, count);
        EpochStatement memory statement = _true(0, seed, count, index);
        unchecked {
            statement.amount += uint64(bound(rawDelta, 1, type(uint64).max));
        }

        court.proveBadProof(statement, _sign(signerKey, statement), _siblings(0, seed, count, index));

        assertEq(court.faultCount(0), 1);
    }

    function testFuzz_AnyWrongIndexSignedIsConvicted(
        uint256 seed,
        uint256 rawCount,
        uint256 rawIndex,
        uint64 rawOther
    ) public {
        uint256 count = bound(rawCount, 2, 40);
        uint256 index = bound(rawIndex, 0, count - 1);
        uint64 other = uint64(bound(rawOther, 0, type(uint64).max - 1));
        if (other >= index) {
            other += 1;
        }

        _anchorNext(seed, count);
        EpochStatement memory statement = _true(0, seed, count, index);
        statement.index = other;

        court.proveBadProof(statement, _sign(signerKey, statement), _siblings(0, seed, count, index));

        assertEq(court.faultCount(0), 1);
    }

    function testFuzz_TwoStatementsConflictExactlyWhenTheyShareALeafOrAnIndex(
        bytes32 firstCommitment,
        bytes32 secondCommitment,
        uint64 firstIndex,
        uint64 secondIndex
    ) public pure {
        EpochStatement memory first = EpochStatement({
            ledger: address(1),
            epoch: 0,
            commitment: firstCommitment,
            amount: 1,
            index: firstIndex,
            siblingsHash: bytes32(0)
        });
        EpochStatement memory second = EpochStatement({
            ledger: address(1),
            epoch: 0,
            commitment: secondCommitment,
            amount: 2,
            index: secondIndex,
            siblingsHash: bytes32(0)
        });

        bool expected = firstCommitment == secondCommitment || firstIndex == secondIndex;
        assertEq(StatementHashing.conflict(first, second), expected);
    }

    function testFuzz_TheCourtRecoversWhoeverSigned(uint256 rawKey, bytes32 commitment, uint64 amount) public view {
        uint256 key = bound(rawKey, 1, CURVE_ORDER - 1);
        EpochStatement memory statement = EpochStatement({
            ledger: address(ledger),
            epoch: 0,
            commitment: commitment,
            amount: amount,
            index: 0,
            siblingsHash: bytes32(0)
        });

        assertEq(court.recoverSigner(statement, _sign(key, statement)), vm.addr(key));
    }
}
