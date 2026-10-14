// SPDX-License-Identifier: MIT
pragma solidity 0.8.37;

import {Test} from "forge-std/src/Test.sol";

import {EpochHeader} from "../../src/EpochHeader.sol";
import {LiabilityLedger} from "../../src/LiabilityLedger.sol";
import {SumNode} from "../../src/SumNode.sol";
import {AccountCommitments} from "../../src/withholding/AccountCommitments.sol";
import {OpenRequest} from "../../src/withholding/OpenRequest.sol";
import {Registration} from "../../src/withholding/Registration.sol";
import {RequestState} from "../../src/withholding/RequestState.sol";
import {StatementRequest} from "../../src/withholding/StatementRequest.sol";
import {WithholdingCourt} from "../../src/withholding/WithholdingCourt.sol";
import {LeafSet} from "../support/LeafSet.sol";
import {TreeBuilder} from "../support/TreeBuilder.sol";

contract WithholdingHandler is Test {
    uint256 public constant HOLDERS = 3;

    WithholdingCourt public immutable withholding;
    LiabilityLedger public immutable ledger;
    address public immutable anchorer;

    Registration[HOLDERS] private _registrations;
    bytes[HOLDERS] private _signatures;
    mapping(uint64 epoch => uint256 seed) private _seeds;
    mapping(uint64 epoch => uint256 count) private _counts;
    mapping(uint64 epoch => uint256 mask) private _masks;
    uint64[] private _anchored;

    uint256 public silences;
    uint256 public answers;
    uint256 public unexpected;

    constructor(WithholdingCourt withholding_, address anchorer_, uint256 signerKey) {
        withholding = withholding_;
        ledger = LiabilityLedger(address(withholding_.ledger()));
        anchorer = anchorer_;

        for (uint256 i = 0; i < HOLDERS; i++) {
            address holder = makeAddr(string(abi.encode("withholding-holder", i)));
            _registrations[i] =
                Registration({holder: holder, account: keccak256(abi.encode("account", i)), fromEpoch: 0});
            (uint8 v, bytes32 r, bytes32 s) = vm.sign(signerKey, withholding_.registrationDigest(_registrations[i]));
            _signatures[i] = abi.encodePacked(r, s, v);
        }
    }

    function anchorNext(uint256 seed, uint256 rawCount, uint256 rawMask, uint256 rawAdvance) external {
        vm.warp(vm.getBlockTimestamp() + bound(rawAdvance, 1 minutes, 1 days));
        uint64 number = ledger.nextEpoch();
        uint256 count = bound(rawCount, HOLDERS, 12);
        uint256 mask = bound(rawMask, 0, (uint256(1) << HOLDERS) - 1);
        SumNode memory root = TreeBuilder.root(_tree(number, seed, count, mask));
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
        _masks[number] = mask;
        _anchored.push(number);
    }

    function request(uint256 rawHolder, uint256 rawEpoch) external {
        if (_anchored.length == 0) {
            return;
        }

        uint256 i = bound(rawHolder, 0, HOLDERS - 1);
        uint64 epoch = _anchored[bound(rawEpoch, 0, _anchored.length - 1)];
        Registration memory registration = _registrations[i];
        bytes32 digest = withholding.registrationDigest(registration);
        bool allowed = vm.getBlockTimestamp() <= uint256(ledger.epochAt(epoch).anchoredAt) + withholding.requestWindow()
            && !withholding.openRequestOf(digest).active
            && withholding.requestOf(digest, epoch).state == RequestState.None;
        bytes memory signature = _signatures[i];

        vm.prank(registration.holder);
        try withholding.requestStatement(registration, signature, epoch) {
            if (!allowed) {
                unexpected++;
            }
        } catch {
            if (allowed) {
                unexpected++;
            }
        }
    }

    function answerTruthfully(uint256 rawHolder) external {
        uint256 i = bound(rawHolder, 0, HOLDERS - 1);
        (bool open, uint64 epoch, StatementRequest memory pending) = _openOf(i);
        if (!open || _masks[epoch] & (uint256(1) << i) == 0) {
            return;
        }

        bool allowed = vm.getBlockTimestamp() <= pending.dueAt;
        SumNode[] memory siblings = TreeBuilder.proof(_tree(epoch, _seeds[epoch], _counts[epoch], _masks[epoch]), i);
        uint64 amount = LeafSet.amounts(_seeds[epoch], _counts[epoch])[i];
        try withholding.answer(_registrations[i], epoch, _salt(epoch, i), amount, uint64(i), siblings) {
            if (!allowed) {
                unexpected++;
            }

            answers++;
        } catch {
            if (allowed) {
                unexpected++;
            }
        }
    }

    function answerFalsely(uint256 rawHolder, bytes32 salt, uint256 rawIndex) external {
        uint256 i = bound(rawHolder, 0, HOLDERS - 1);
        (bool open, uint64 epoch,) = _openOf(i);
        if (!open) {
            return;
        }

        uint256 index = bound(rawIndex, 0, _counts[epoch] - 1);
        bool included = _masks[epoch] & (uint256(1) << i) != 0;
        if (included && index == i && salt == _salt(epoch, i)) {
            return;
        }

        SumNode[] memory siblings = TreeBuilder.proof(_tree(epoch, _seeds[epoch], _counts[epoch], _masks[epoch]), index);
        uint64 amount = LeafSet.amounts(_seeds[epoch], _counts[epoch])[index];
        try withholding.answer(_registrations[i], epoch, salt, amount, uint64(index), siblings) {
            unexpected++;
        } catch {}
    }

    function recordSilence(uint256 rawHolder) external {
        uint256 i = bound(rawHolder, 0, HOLDERS - 1);
        (bool open, uint64 epoch, StatementRequest memory pending) = _openOf(i);
        bool allowed = open && vm.getBlockTimestamp() > pending.dueAt;

        try withholding.recordUnanswered(_registrations[i], epoch) {
            if (!allowed) {
                unexpected++;
            }

            silences++;
        } catch {
            if (allowed) {
                unexpected++;
            }
        }
    }

    function wait(uint256 rawAdvance) external {
        vm.warp(vm.getBlockTimestamp() + bound(rawAdvance, 0, 2 days));
    }

    function anchoredCount() external view returns (uint256) {
        return _anchored.length;
    }

    function anchoredAt(uint256 position) external view returns (uint64) {
        return _anchored[position];
    }

    function registrationOf(uint256 i) external view returns (Registration memory) {
        return _registrations[i];
    }

    function _openOf(uint256 i) private view returns (bool, uint64, StatementRequest memory) {
        bytes32 digest = withholding.registrationDigest(_registrations[i]);
        OpenRequest memory open = withholding.openRequestOf(digest);
        return (open.active, open.epoch, withholding.requestOf(digest, open.epoch));
    }

    function _tree(uint64 number, uint256 seed, uint256 count, uint256 mask) private view returns (SumNode[][] memory) {
        bytes32[] memory commitments = LeafSet.commitments(seed, count);
        for (uint256 i = 0; i < HOLDERS; i++) {
            if (mask & (uint256(1) << i) != 0) {
                commitments[i] = AccountCommitments.commitment(_registrations[i].account, _salt(number, i));
            }
        }

        return TreeBuilder.build(number, commitments, LeafSet.amounts(seed, count));
    }

    function _salt(uint64 epoch, uint256 i) private pure returns (bytes32) {
        return keccak256(abi.encode("salt", epoch, i));
    }
}
