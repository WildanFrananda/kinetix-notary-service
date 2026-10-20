// SPDX-License-Identifier: MIT
pragma solidity 0.8.37;

import {
    AccessControlDefaultAdminRules
} from "@openzeppelin-contracts/access/extensions/AccessControlDefaultAdminRules.sol";
import {SafeCast} from "@openzeppelin-contracts/utils/math/SafeCast.sol";

import {Epoch} from "../Epoch.sol";
import {EpochHeader} from "../EpochHeader.sol";
import {ILiabilityLedger} from "../ILiabilityLedger.sol";
import {SumNode} from "../SumNode.sol";
import {ConfidentialEpochChain} from "./ConfidentialEpochChain.sol";
import {ConfidentialTree} from "./ConfidentialTree.sol";
import {IConfidentialLedger} from "./IConfidentialLedger.sol";
import {PlonkVerifier} from "./PlonkVerifier.sol";
import {Poseidon2} from "./Poseidon2.sol";
import {ProvenStatement} from "./ProvenStatement.sol";

contract LiabilityLedgerV2 is ILiabilityLedger, IConfidentialLedger, AccessControlDefaultAdminRules {
    bytes32 public constant ANCHORER_ROLE = keccak256("kinetix.notary.ANCHORER");

    uint64 public immutable genesisEpoch;
    bytes32 public immutable genesisDigest;
    PlonkVerifier private immutable _verifier;
    uint8 public immutable depth;

    mapping(uint64 number => Epoch) private _epochs;
    mapping(uint64 number => ProvenStatement) private _anchoredStatements;
    mapping(bytes32 key => ProvenStatement) private _proven;
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
        uint64 genesisAsOf,
        PlonkVerifier verifier_,
        uint8 depth_
    ) AccessControlDefaultAdminRules(adminTransferDelay, admin) {
        if ((genesisEpoch_ == 0) != (genesisDigest_ == bytes32(0))) {
            revert InvalidGenesis(genesisEpoch_, genesisDigest_);
        }

        if (anchorer == address(0)) {
            revert ZeroAnchorer();
        }

        if (address(verifier_).code.length == 0) {
            revert NotAVerifier(address(verifier_));
        }

        if (depth_ > ConfidentialTree.MAX_DEPTH) {
            revert DepthOutOfRange(depth_, ConfidentialTree.MAX_DEPTH);
        }

        if (!_grantRole(ANCHORER_ROLE, anchorer)) {
            revert AnchorerNotGranted(anchorer);
        }

        genesisEpoch = genesisEpoch_;
        genesisDigest = genesisDigest_;
        _verifier = verifier_;
        depth = depth_;
        _head = genesisDigest_;
        _nextEpoch = genesisEpoch_;
        _lastAsOf = genesisAsOf;
    }

    function proveEpoch(
        uint64 number,
        bytes32 treeRoot,
        bytes32 totalCommitment,
        uint64 leafCount,
        bytes calldata proof
    ) external returns (bytes32) {
        if (_sealed) {
            revert LedgerIsSealed();
        }

        if (number < _nextEpoch) {
            revert EpochClosed(number, _nextEpoch);
        }

        if (leafCount > capacity()) {
            revert TooManyLeaves(leafCount, capacity());
        }

        if (uint256(treeRoot) >= Poseidon2.PRIME) {
            revert NotAFieldElement(treeRoot);
        }

        if (uint256(totalCommitment) >= Poseidon2.PRIME) {
            revert NotAFieldElement(totalCommitment);
        }

        bytes32 root = ConfidentialTree.bindingRoot(treeRoot, totalCommitment);
        bytes32 key = _statementKey(number, root, leafCount);
        if (_proven[key].proven) {
            revert StatementAlreadyProven(number, root, leafCount);
        }

        if (!_verifies(proof, number, treeRoot, totalCommitment, leafCount)) {
            revert ProofRejected(number);
        }

        _proven[key] = ProvenStatement({treeRoot: treeRoot, totalCommitment: totalCommitment, proven: true});
        emit EpochProven(number, root, treeRoot, totalCommitment, leafCount);
        return root;
    }

    function anchor(EpochHeader calldata header) external onlyRole(ANCHORER_ROLE) returns (bytes32) {
        if (_sealed) {
            revert LedgerIsSealed();
        }

        bytes32 digest = ConfidentialEpochChain.digest(header);
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

        ProvenStatement memory statement = _proven[_statementKey(header.number, header.root, header.leafCount)];
        if (!statement.proven) {
            revert EpochNotProven(header.number, header.root, header.leafCount);
        }

        _epochs[header.number] = Epoch({
            root: header.root,
            digest: digest,
            total: 0,
            leafCount: header.leafCount,
            asOf: header.asOf,
            anchoredAt: SafeCast.toUint64(block.timestamp)
        });
        _anchoredStatements[header.number] = statement;
        _head = digest;
        _nextEpoch = header.number + 1;
        _lastAsOf = header.asOf;

        emit EpochAnchored(header.number, digest, header.previousDigest, header.root, 0, header.leafCount, header.asOf);
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

        if (siblings.length != depth) {
            return false;
        }

        if (!ConfidentialTree.wellFormed(siblings)) {
            return false;
        }

        uint256 reached =
            ConfidentialTree.computeRoot(ConfidentialTree.leaf(number, commitment, amount), index, siblings);
        return bytes32(reached) == _anchoredStatements[number].treeRoot;
    }

    function epochAt(uint64 number) external view returns (Epoch memory) {
        return _anchored(number);
    }

    function isProven(uint64 number, bytes32 root, uint64 leafCount) external view returns (bool) {
        return _proven[_statementKey(number, root, leafCount)].proven;
    }

    function treeRootOf(uint64 number) external view returns (bytes32) {
        _requireAnchored(number);
        return _anchoredStatements[number].treeRoot;
    }

    function totalCommitmentOf(uint64 number) external view returns (bytes32) {
        _requireAnchored(number);
        return _anchoredStatements[number].totalCommitment;
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
        return ConfidentialTree.SCHEME;
    }

    function verifier() external view returns (address) {
        return address(_verifier);
    }

    function capacity() public view returns (uint64) {
        return uint64(1) << depth;
    }

    function supportsInterface(bytes4 interfaceId) public view override returns (bool) {
        return interfaceId == type(ILiabilityLedger).interfaceId || interfaceId == type(IConfidentialLedger).interfaceId
            || super.supportsInterface(interfaceId);
    }

    function _verifies(
        bytes calldata proof,
        uint64 number,
        bytes32 treeRoot,
        bytes32 totalCommitment,
        uint64 leafCount
    ) private view returns (bool) {
        uint256[] memory inputs = new uint256[](4);
        inputs[0] = number;
        inputs[1] = uint256(treeRoot);
        inputs[2] = uint256(totalCommitment);
        inputs[3] = leafCount;
        try _verifier.Verify(proof, inputs) returns (bool verified) {
            return verified;
        } catch {
            return false;
        }
    }

    function _anchored(uint64 number) private view returns (Epoch storage) {
        _requireAnchored(number);
        return _epochs[number];
    }

    function _requireAnchored(uint64 number) private view {
        if (number < genesisEpoch || number >= _nextEpoch) {
            revert UnknownEpoch(number);
        }
    }

    function _requireWellFormed(EpochHeader calldata header) private view {
        if (header.total != 0) {
            revert TotalNotHidden(header.total);
        }

        if (header.leafCount > capacity()) {
            revert TooManyLeaves(header.leafCount, capacity());
        }

        if (header.asOf <= _lastAsOf) {
            revert SnapshotNotAfterPrevious(header.asOf, _lastAsOf);
        }

        if (header.asOf > block.timestamp) {
            revert SnapshotInFuture(header.asOf, SafeCast.toUint64(block.timestamp));
        }
    }

    function _statementKey(uint64 number, bytes32 root, uint64 leafCount) private pure returns (bytes32) {
        return keccak256(abi.encode(number, root, leafCount));
    }
}
