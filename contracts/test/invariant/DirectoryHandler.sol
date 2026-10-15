// SPDX-License-Identifier: MIT
pragma solidity 0.8.37;

import {Test} from "forge-std/src/Test.sol";

import {EpochHeader} from "../../src/EpochHeader.sol";
import {LiabilityLedger} from "../../src/LiabilityLedger.sol";
import {SumNode} from "../../src/SumNode.sol";
import {StatementCourt} from "../../src/court/StatementCourt.sol";
import {Generation} from "../../src/directory/Generation.sol";
import {LedgerDirectory} from "../../src/directory/LedgerDirectory.sol";
import {WithholdingCourt} from "../../src/withholding/WithholdingCourt.sol";
import {LeafSet} from "../support/LeafSet.sol";
import {TreeBuilder} from "../support/TreeBuilder.sol";

contract DirectoryHandler is Test {
    uint48 private constant ADMIN_DELAY = 3 days;
    uint64 private constant SIGNER_DELAY = 7 days;
    uint64 private constant WINDOW = 7 days;
    uint64 private constant PERIOD = 24 hours;

    LedgerDirectory public immutable directory;
    address public immutable admin;
    address public immutable anchorer;
    address public immutable signer;
    address public immutable stranger;

    bytes32[] private _recorded;

    uint256 public unexpected;

    constructor(LedgerDirectory directory_, address admin_, address anchorer_, address signer_, address stranger_) {
        directory = directory_;
        admin = admin_;
        anchorer = anchorer_;
        signer = signer_;
        stranger = stranger_;
        _record(0);
    }

    function recordedCount() external view returns (uint256) {
        return _recorded.length;
    }

    function recorded(uint256 index) external view returns (bytes32) {
        return _recorded[index];
    }

    function anchor(uint256 seed, uint256 rawAdvance) external {
        vm.warp(vm.getBlockTimestamp() + bound(rawAdvance, 1 minutes, 1 days));
        LiabilityLedger current = _currentLocal();
        if (address(current) == address(0) || current.isSealed()) {
            return;
        }

        uint64 number = current.nextEpoch();
        uint256 count = 1 + seed % 4;
        SumNode[][] memory levels =
            TreeBuilder.build(number, LeafSet.commitments(seed, count), LeafSet.amounts(seed, count));
        SumNode memory root = TreeBuilder.root(levels);
        EpochHeader memory header = EpochHeader({
            number: number,
            previousDigest: current.head(),
            root: root.hash,
            total: root.sum,
            leafCount: uint64(count),
            asOf: uint64(vm.getBlockTimestamp()) - 1
        });
        vm.prank(anchorer);
        current.anchor(header);
    }

    function succeed() external {
        LiabilityLedger current = _currentLocal();
        if (address(current) == address(0)) {
            _expectRefusal(_registration(LiabilityLedger(directory.generation(0).ledger)));
            return;
        }

        _sealIfOpen(current);
        bool pending = directory.pendingSuccessor().generation.chainId != 0;
        bytes memory call = _registration(current);
        vm.prank(admin);
        (bool succeeded,) = address(directory).call(call);
        if (succeeded == pending) {
            unexpected++;
        }

        if (succeeded) {
            _record(directory.generationCount() - 1);
        }
    }

    function forge(uint256 rawShift) external {
        LiabilityLedger current = _currentLocal();
        if (address(current) == address(0)) {
            return;
        }

        _sealIfOpen(current);
        uint64 shift = uint64(bound(rawShift, 1, 3));
        LiabilityLedger skipping = new LiabilityLedger(
            admin,
            ADMIN_DELAY,
            anchorer,
            current.nextEpoch() + shift,
            keccak256(abi.encode("forged", rawShift)),
            current.lastAsOf()
        );
        _expectRefusal(_registrationOf(skipping));
    }

    function intrude() external {
        LiabilityLedger current = _currentLocal();
        if (address(current) == address(0)) {
            return;
        }

        _sealIfOpen(current);
        bytes memory call = _registration(current);
        vm.prank(stranger);
        (bool succeeded,) = address(directory).call(call);
        if (succeeded) {
            unexpected++;
            _record(directory.generationCount() - 1);
        }
    }

    function propose(uint64 rawChain) external {
        LiabilityLedger current = _currentLocal();
        if (address(current) == address(0)) {
            return;
        }

        _sealIfOpen(current);
        uint64 chainId = uint64(bound(rawChain, 1, type(uint64).max));
        bool foreign = chainId != block.chainid;
        bool pending = directory.pendingSuccessor().generation.chainId != 0;
        vm.prank(admin);
        try directory.proposeForeignSuccessor(chainId, address(0xF1), address(0xF2), address(0xF3), 1) {
            if (pending || !foreign) {
                unexpected++;
            }
        } catch {
            if (!pending && foreign) {
                unexpected++;
            }
        }
    }

    function cancel() external {
        bool pending = directory.pendingSuccessor().generation.chainId != 0;
        vm.prank(admin);
        try directory.cancelForeignSuccessor() {
            if (!pending) {
                unexpected++;
            }
        } catch {
            if (pending) {
                unexpected++;
            }
        }
    }

    function activate(uint256 rawWait) external {
        vm.warp(vm.getBlockTimestamp() + bound(rawWait, 0, 8 days));
        uint64 activatesAt = directory.pendingSuccessor().activatesAt;
        bool due = activatesAt != 0 && vm.getBlockTimestamp() >= activatesAt;
        try directory.activateForeignSuccessor() returns (uint256 index) {
            if (!due) {
                unexpected++;
            }

            _record(index);
        } catch {
            if (due) {
                unexpected++;
            }
        }
    }

    function _currentLocal() private view returns (LiabilityLedger) {
        Generation memory current = directory.current();
        if (current.chainId != block.chainid) {
            return LiabilityLedger(address(0));
        }

        return LiabilityLedger(current.ledger);
    }

    function _sealIfOpen(LiabilityLedger current) private {
        if (current.isSealed()) {
            return;
        }

        bytes32 head = current.head();
        vm.prank(admin);
        current.seal(head);
    }

    function _registration(LiabilityLedger predecessor) private returns (bytes memory) {
        LiabilityLedger next = new LiabilityLedger(
            admin, ADMIN_DELAY, anchorer, predecessor.nextEpoch(), predecessor.head(), predecessor.lastAsOf()
        );
        return _registrationOf(next);
    }

    function _registrationOf(LiabilityLedger next) private returns (bytes memory) {
        StatementCourt nextCourt = new StatementCourt(admin, ADMIN_DELAY, next, signer, SIGNER_DELAY);
        WithholdingCourt nextWithholding = new WithholdingCourt(nextCourt, WINDOW, PERIOD);
        return abi.encodeCall(LedgerDirectory.registerSuccessor, (next, nextCourt, nextWithholding));
    }

    function _expectRefusal(bytes memory call) private {
        vm.prank(admin);
        (bool succeeded,) = address(directory).call(call);
        if (succeeded) {
            unexpected++;
            _record(directory.generationCount() - 1);
        }
    }

    function _record(uint256 index) private {
        _recorded.push(keccak256(abi.encode(directory.generation(index))));
    }
}
