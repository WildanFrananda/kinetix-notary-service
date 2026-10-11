// SPDX-License-Identifier: UNLICENSED
pragma solidity 0.8.37;

import {
    AccessControlDefaultAdminRules
} from "@openzeppelin-contracts/access/extensions/AccessControlDefaultAdminRules.sol";
import {SafeCast} from "@openzeppelin-contracts/utils/math/SafeCast.sol";

import {Epoch} from "./Epoch.sol";
import {EpochChain} from "./EpochChain.sol";
import {EpochHeader} from "./EpochHeader.sol";
import {ILiabilityLedger} from "./ILiabilityLedger.sol";
import {SumNode} from "./SumNode.sol";
import {SumTree} from "./SumTree.sol";

contract LiabilityLedger is ILiabilityLedger, AccessControlDefaultAdminRules {
    bytes32 public constant ANCHORER_ROLE = keccak256("kinetix.notary.ANCHORER");

    uint64 public immutable genesisEpoch;
    bytes32 public immutable genesisDigest;

    mapping(uint64 number => Epoch) private _epochs;
    bytes32 private _head;
    uint64 private _nextEpoch;
    uint64 private _lastAsOf;
    bool private _sealed;

    constructor(
        address admin,
        uint48 adminTransferDelay,
        address anchorer,
        uint64 genesisEpoch_,
        bytes32 genesisDigest_,
        uint64 genesisAsOf
    ) AccessControlDefaultAdminRules(adminTransferDelay, admin) {
        if ((genesisEpoch_ == 0) != (genesisDigest_ == bytes32(0))) {
            revert InvalidGenesis(genesisEpoch_, genesisDigest_);
        }
        if (anchorer == address(0)) {
            revert ZeroAnchorer();
        }
        _grantRole(ANCHORER_ROLE, anchorer);
        genesisEpoch = genesisEpoch_;
        genesisDigest = genesisDigest_;
        _head = genesisDigest_;
        _nextEpoch = genesisEpoch_;
        _lastAsOf = genesisAsOf;
    }

    function anchor(EpochHeader calldata header) external onlyRole(ANCHORER_ROLE) returns (bytes32) {
        if (_sealed) {
            revert LedgerIsSealed();
        }
        bytes32 digest = EpochChain.digest(header);
        if (header.number < genesisEpoch) {
            revert EpochBeforeGenesis(genesisEpoch, header.number);
        }
        if (header.number < _nextEpoch) {
            bytes32 anchored = _epochs[header.number].digest;
            if (anchored == digest) {
                revert EpochAlreadyAnchored(header.number, digest);
            }
            revert EpochConflict(header.number, anchored, digest);
        }
        if (header.number > _nextEpoch) {
            revert EpochOutOfSequence(_nextEpoch, header.number);
        }
        if (header.previousDigest != _head) {
            revert PreviousDigestMismatch(_head, header.previousDigest);
        }
        _requireWellFormed(header);

        _epochs[header.number] = Epoch({
            root: header.root,
            digest: digest,
            total: header.total,
            leafCount: header.leafCount,
            asOf: header.asOf,
            anchoredAt: SafeCast.toUint64(block.timestamp)
        });
        _head = digest;
        _nextEpoch = header.number + 1;
        _lastAsOf = header.asOf;

        emit EpochAnchored(
            header.number, digest, header.previousDigest, header.root, header.total, header.leafCount, header.asOf
        );
        return digest;
    }

    function seal(bytes32 expectedHead) external onlyRole(DEFAULT_ADMIN_ROLE) {
        if (_sealed) {
            revert LedgerIsSealed();
        }
        if (expectedHead != _head) {
            revert HeadMismatch(_head, expectedHead);
        }
        _sealed = true;
        emit LedgerSealed(_head, _nextEpoch, _lastAsOf);
    }

    function verifyInclusion(
        uint64 number,
        bytes32 commitment,
        uint64 amount,
        uint64 index,
        SumNode[] calldata siblings
    ) external view returns (bool) {
        Epoch storage anchored = _anchored(number);
        if (index >= anchored.leafCount) {
            return false;
        }
        if (siblings.length != SumTree.depthOf(anchored.leafCount)) {
            return false;
        }
        if (SumTree.pathSum(amount, siblings) != anchored.total) {
            return false;
        }
        return SumTree.computeRoot(SumTree.leaf(number, commitment, amount), index, siblings).hash == anchored.root;
    }

    function epochAt(uint64 number) external view returns (Epoch memory) {
        return _anchored(number);
    }

    function head() external view returns (bytes32) {
        return _head;
    }

    function nextEpoch() external view returns (uint64) {
        return _nextEpoch;
    }

    function lastAsOf() external view returns (uint64) {
        return _lastAsOf;
    }

    function isSealed() external view returns (bool) {
        return _sealed;
    }

    function scheme() external pure returns (uint16) {
        return SumTree.SCHEME;
    }

    function supportsInterface(bytes4 interfaceId) public view override returns (bool) {
        return interfaceId == type(ILiabilityLedger).interfaceId || super.supportsInterface(interfaceId);
    }

    function _anchored(uint64 number) private view returns (Epoch storage) {
        if (number < genesisEpoch || number >= _nextEpoch) {
            revert UnknownEpoch(number);
        }
        return _epochs[number];
    }

    function _requireWellFormed(EpochHeader calldata header) private view {
        if (header.leafCount > SumTree.MAX_LEAVES) {
            revert TooManyLeaves(header.leafCount, SumTree.MAX_LEAVES);
        }
        if (header.leafCount == 0 && (header.root != SumTree.empty().hash || header.total != 0)) {
            revert EmptyEpochMismatch(header.root, header.total);
        }
        if (header.asOf <= _lastAsOf) {
            revert SnapshotNotAfterPrevious(header.asOf, _lastAsOf);
        }
        if (header.asOf > block.timestamp) {
            revert SnapshotInFuture(header.asOf, SafeCast.toUint64(block.timestamp));
        }
    }
}
