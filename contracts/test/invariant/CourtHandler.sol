// SPDX-License-Identifier: MIT
pragma solidity 0.8.37;

import {Test} from "forge-std/src/Test.sol";

import {EpochHeader} from "../../src/EpochHeader.sol";
import {LiabilityLedger} from "../../src/LiabilityLedger.sol";
import {SumNode} from "../../src/SumNode.sol";
import {EpochStatement} from "../../src/court/EpochStatement.sol";
import {PendingSigner} from "../../src/court/PendingSigner.sol";
import {StatementCourt} from "../../src/court/StatementCourt.sol";
import {StatementHashing} from "../../src/court/StatementHashing.sol";
import {LeafSet} from "../support/LeafSet.sol";
import {TreeBuilder} from "../support/TreeBuilder.sol";

contract CourtHandler is Test {
    LiabilityLedger public immutable ledger;
    StatementCourt public immutable court;
    address public immutable admin;
    address public immutable anchorer;

    mapping(address signer => uint256 key) private _keys;
    mapping(uint64 epoch => uint256 seed) private _seeds;
    mapping(uint64 epoch => uint256 count) private _counts;
    mapping(uint64 epoch => address signer) private _observed;
    mapping(bytes32 evidence => bool proven) private _proven;
    uint64[] private _anchored;
    uint256 private _signerCount;

    uint256 public faults;
    uint256 public unexpected;

    constructor(StatementCourt court_, address admin_, address anchorer_, address signer, uint256 signerKey) {
        court = court_;
        ledger = LiabilityLedger(address(court_.ledger()));
        admin = admin_;
        anchorer = anchorer_;
        _keys[signer] = signerKey;
    }

    function anchorNext(uint256 seed, uint256 rawCount, uint256 rawAdvance) external {
        vm.warp(vm.getBlockTimestamp() + bound(rawAdvance, 1 minutes, 2 days));
        uint256 count = bound(rawCount, 1, 12);
        uint64 number = ledger.nextEpoch();
        SumNode[][] memory levels =
            TreeBuilder.build(number, LeafSet.commitments(seed, count), LeafSet.amounts(seed, count));
        SumNode memory root = TreeBuilder.root(levels);
        EpochHeader memory header = EpochHeader({
            number: number,
            previousDigest: ledger.head(),
            root: root.hash,
            total: root.sum,
            leafCount: uint64(count),
            asOf: uint64(vm.getBlockTimestamp() - 1)
        });

        vm.prank(anchorer);
        ledger.anchor(header);
        _seeds[number] = seed;
        _counts[number] = count;
        _observed[number] = court.signerOf(number);
        _anchored.push(number);
    }

    function proveATrueStatement(uint256 rawEpoch, uint256 rawIndex) external {
        if (_anchored.length == 0) {
            return;
        }

        (EpochStatement memory statement, SumNode[] memory siblings) = _statement(rawEpoch, rawIndex);
        bytes memory signature = _sign(_keys[court.signerOf(statement.epoch)], statement);
        try court.proveBadProof(statement, signature, siblings) {
            unexpected++;
        } catch {}
    }

    function proveAFalseStatement(uint256 rawEpoch, uint256 rawIndex, uint64 rawDelta) external {
        if (_anchored.length == 0) {
            return;
        }

        (EpochStatement memory statement, SumNode[] memory siblings) = _statement(rawEpoch, rawIndex);
        unchecked {
            statement.amount += uint64(bound(rawDelta, 1, type(uint64).max));
        }

        bytes memory signature = _sign(_keys[court.signerOf(statement.epoch)], statement);
        bytes32 digest = court.statementDigest(statement);
        try court.proveBadProof(statement, signature, siblings) returns (bytes32 evidence) {
            if (_proven[digest] || evidence != digest) {
                unexpected++;
            }

            _proven[evidence] = true;
            faults++;
        } catch {
            if (!_proven[digest]) {
                unexpected++;
            }
        }
    }

    function proveAForgery(uint256 rawEpoch, uint256 rawIndex, uint256 forgerSeed) external {
        if (_anchored.length == 0) {
            return;
        }

        (EpochStatement memory statement, SumNode[] memory siblings) = _statement(rawEpoch, rawIndex);
        statement.amount ^= 1;
        (, uint256 forgerKey) = makeAddrAndKey(string(abi.encode("forger", forgerSeed)));
        bytes memory signature = _sign(forgerKey, statement);
        try court.proveBadProof(statement, signature, siblings) {
            unexpected++;
        } catch {}
    }

    function equivocate(uint256 rawEpoch, uint256 rawIndex, uint64 rawDelta) external {
        if (_anchored.length == 0) {
            return;
        }

        (EpochStatement memory first,) = _statement(rawEpoch, rawIndex);
        EpochStatement memory second = EpochStatement({
            ledger: first.ledger,
            epoch: first.epoch,
            commitment: first.commitment,
            amount: first.amount,
            index: first.index,
            siblingsHash: first.siblingsHash
        });
        unchecked {
            second.amount += uint64(bound(rawDelta, 1, type(uint64).max));
        }

        uint256 key = _keys[court.signerOf(first.epoch)];
        bytes memory firstSignature = _sign(key, first);
        bytes memory secondSignature = _sign(key, second);
        bytes32 a = court.statementDigest(first);
        bytes32 b = court.statementDigest(second);
        bytes32 expected = a < b ? keccak256(abi.encodePacked(a, b)) : keccak256(abi.encodePacked(b, a));
        try court.proveEquivocation(first, firstSignature, second, secondSignature) returns (bytes32 evidence) {
            if (_proven[expected] || evidence != expected) {
                unexpected++;
            }

            _proven[evidence] = true;
            faults++;
        } catch {
            if (!_proven[expected]) {
                unexpected++;
            }
        }
    }

    function scheduleSigner(uint256 rawAhead) external {
        _signerCount++;
        (address next, uint256 key) = makeAddrAndKey(string(abi.encode("signer", _signerCount)));
        _keys[next] = key;
        uint64 fromEpoch = ledger.nextEpoch() + uint64(bound(rawAhead, 0, 5));
        PendingSigner memory before = court.pendingSigner();

        vm.prank(admin);
        try court.scheduleSigner(next, fromEpoch) {
            if (before.signer != address(0)) {
                unexpected++;
            }
        } catch {}
    }

    function activateSigner(uint256 rawWait) external {
        vm.warp(vm.getBlockTimestamp() + bound(rawWait, 0, 9 days));
        PendingSigner memory pending = court.pendingSigner();
        bool allowed = pending.signer != address(0) && vm.getBlockTimestamp() >= pending.activatesAt
            && pending.fromEpoch >= ledger.nextEpoch();
        try court.activateSigner() {
            if (!allowed) {
                unexpected++;
            }
        } catch {
            if (allowed) {
                unexpected++;
            }
        }
    }

    function cancelSigner() external {
        bool pending = court.pendingSigner().signer != address(0);
        vm.prank(admin);
        try court.cancelSigner() {
            if (!pending) {
                unexpected++;
            }
        } catch {
            if (pending) {
                unexpected++;
            }
        }
    }

    function anchoredCount() external view returns (uint256) {
        return _anchored.length;
    }

    function anchoredAt(uint256 i) external view returns (uint64) {
        return _anchored[i];
    }

    function observedSigner(uint64 epoch) external view returns (address) {
        return _observed[epoch];
    }

    function _statement(
        uint256 rawEpoch,
        uint256 rawIndex
    ) private view returns (EpochStatement memory, SumNode[] memory) {
        uint64 epoch = _anchored[bound(rawEpoch, 0, _anchored.length - 1)];
        uint256 seed = _seeds[epoch];
        uint256 count = _counts[epoch];
        uint256 index = bound(rawIndex, 0, count - 1);
        SumNode[][] memory levels =
            TreeBuilder.build(epoch, LeafSet.commitments(seed, count), LeafSet.amounts(seed, count));
        SumNode[] memory siblings = TreeBuilder.proof(levels, index);
        EpochStatement memory statement = EpochStatement({
            ledger: address(ledger),
            epoch: epoch,
            commitment: LeafSet.commitments(seed, count)[index],
            amount: LeafSet.amounts(seed, count)[index],
            index: uint64(index),
            siblingsHash: StatementHashing.siblingsHash(siblings)
        });
        return (statement, siblings);
    }

    function _sign(uint256 key, EpochStatement memory statement) private view returns (bytes memory) {
        (uint8 v, bytes32 r, bytes32 s) = vm.sign(key, court.statementDigest(statement));
        return abi.encodePacked(r, s, v);
    }
}
