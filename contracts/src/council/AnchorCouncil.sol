// SPDX-License-Identifier: MIT
pragma solidity 0.8.37;

import {
    AccessControlDefaultAdminRules
} from "@openzeppelin-contracts/access/extensions/AccessControlDefaultAdminRules.sol";
import {ECDSA} from "@openzeppelin-contracts/utils/cryptography/ECDSA.sol";
import {EIP712} from "@openzeppelin-contracts/utils/cryptography/EIP712.sol";
import {ERC165Checker} from "@openzeppelin-contracts/utils/introspection/ERC165Checker.sol";
import {SafeCast} from "@openzeppelin-contracts/utils/math/SafeCast.sol";

import {EpochChain} from "../EpochChain.sol";
import {EpochHeader} from "../EpochHeader.sol";
import {ILiabilityLedger} from "../ILiabilityLedger.sol";
import {ApprovalHashing} from "./ApprovalHashing.sol";
import {IAnchorCouncil} from "./IAnchorCouncil.sol";
import {PendingMembership} from "./PendingMembership.sol";

contract AnchorCouncil is IAnchorCouncil, AccessControlDefaultAdminRules, EIP712 {
    uint256 public constant MAX_MEMBERS = 16;
    uint8 public constant MIN_THRESHOLD = 2;

    ILiabilityLedger public immutable ledger;
    uint64 public immutable membershipChangeDelay;

    address[] private _members;
    uint8 private _threshold;
    PendingMembership private _pending;

    constructor(
        address admin,
        uint48 adminTransferDelay,
        ILiabilityLedger ledger_,
        address[] memory members_,
        uint8 threshold_,
        uint64 membershipChangeDelay_
    ) AccessControlDefaultAdminRules(adminTransferDelay, admin) EIP712(ApprovalHashing.NAME, ApprovalHashing.VERSION) {
        if (!_isLedger(ledger_)) {
            revert NotALedger(address(ledger_));
        }

        if (membershipChangeDelay_ == 0) {
            revert ZeroDelay();
        }

        _requireValid(members_, threshold_);
        _members = members_;
        _threshold = threshold_;
        membershipChangeDelay = membershipChangeDelay_;
        ledger = ledger_;
    }

    function anchor(EpochHeader calldata header, bytes[] calldata signatures) external returns (bytes32) {
        bytes32 digest = EpochChain.digest(header);
        address[] memory approvers = _approvers(header.number, digest, signatures);

        emit EpochApproved(header.number, digest, approvers);
        return ledger.anchor(header);
    }

    function scheduleMembership(
        address[] calldata nextMembers,
        uint8 nextThreshold
    ) external onlyRole(DEFAULT_ADMIN_ROLE) {
        if (_hasPending()) {
            revert MembershipChangePending(_pending.activatesAt);
        }

        _requireValid(nextMembers, nextThreshold);
        if (nextThreshold == _threshold && _same(nextMembers)) {
            revert SameMembership();
        }

        uint64 activatesAt = SafeCast.toUint64(block.timestamp) + membershipChangeDelay;
        _pending.members = nextMembers;
        _pending.threshold = nextThreshold;
        _pending.activatesAt = activatesAt;
        emit MembershipScheduled(nextMembers, nextThreshold, activatesAt);
    }

    function cancelMembership() external onlyRole(DEFAULT_ADMIN_ROLE) {
        if (!_hasPending()) {
            revert NoMembershipChangePending();
        }

        PendingMembership memory pending = _pending;
        delete _pending;
        emit MembershipCancelled(pending.members, pending.threshold);
    }

    function activateMembership() external {
        if (!_hasPending()) {
            revert NoMembershipChangePending();
        }

        PendingMembership memory pending = _pending;
        uint64 blockTime = SafeCast.toUint64(block.timestamp);
        if (blockTime < pending.activatesAt) {
            revert MembershipChangeTooEarly(pending.activatesAt, blockTime);
        }

        delete _pending;
        _members = pending.members;
        _threshold = pending.threshold;
        emit MembershipActivated(pending.members, pending.threshold);
    }

    function members() external view returns (address[] memory) {
        return _members;
    }

    function threshold() external view returns (uint8) {
        return _threshold;
    }

    function isMember(address account) external view returns (bool) {
        address[] memory current = _members;
        return _contains(current, account);
    }

    function pendingMembership() external view returns (PendingMembership memory) {
        return _pending;
    }

    function epochDigest(EpochHeader calldata header) external pure returns (bytes32) {
        return EpochChain.digest(header);
    }

    function approvalDigest(uint64 number, bytes32 digest) public view returns (bytes32) {
        return _hashTypedDataV4(ApprovalHashing.structHash(address(ledger), number, digest));
    }

    function recoverApprover(uint64 number, bytes32 digest, bytes calldata signature) external view returns (address) {
        return _recover(approvalDigest(number, digest), signature);
    }

    function supportsInterface(bytes4 interfaceId) public view override returns (bool) {
        return interfaceId == type(IAnchorCouncil).interfaceId || super.supportsInterface(interfaceId);
    }

    function _approvers(
        uint64 number,
        bytes32 digest,
        bytes[] calldata signatures
    ) private view returns (address[] memory) {
        uint256 given = signatures.length;
        if (given < _threshold) {
            revert NotEnoughApprovals(given, _threshold);
        }

        bytes32 hash = approvalDigest(number, digest);
        address[] memory approvers = new address[](given);
        for (uint256 i = 0; i < given; i++) {
            approvers[i] = _recover(hash, signatures[i]);
        }

        address[] memory current = _members;
        uint256 refused = _firstRefused(current, approvers);
        if (refused == given) {
            return approvers;
        }

        if (!_contains(current, approvers[refused])) {
            revert NotAMember(approvers[refused]);
        }

        revert ApprovalsOutOfOrder(approvers[refused - 1], approvers[refused]);
    }

    function _firstRefused(address[] memory current, address[] memory approvers) private pure returns (uint256) {
        address previous = address(0);
        for (uint256 i = 0; i < approvers.length; i++) {
            if (!_contains(current, approvers[i]) || approvers[i] <= previous) {
                return i;
            }

            previous = approvers[i];
        }

        return approvers.length;
    }

    function _hasPending() private view returns (bool) {
        return _pending.members.length != 0;
    }

    function _same(address[] calldata members_) private view returns (bool) {
        if (members_.length != _members.length) {
            return false;
        }

        for (uint256 i = 0; i < members_.length; i++) {
            if (members_[i] != _members[i]) {
                return false;
            }
        }

        return true;
    }

    function _contains(address[] memory list, address account) private pure returns (bool) {
        for (uint256 i = 0; i < list.length; i++) {
            if (list[i] == account) {
                return true;
            }
        }

        return false;
    }

    function _isLedger(ILiabilityLedger candidate) private view returns (bool) {
        return ERC165Checker.supportsInterface(address(candidate), type(ILiabilityLedger).interfaceId);
    }

    function _requireValid(address[] memory members_, uint8 threshold_) private pure {
        if (threshold_ < MIN_THRESHOLD) {
            revert ThresholdTooLow(threshold_, MIN_THRESHOLD);
        }

        if (members_.length < threshold_) {
            revert TooFewMembers(members_.length, threshold_);
        }

        if (members_.length > MAX_MEMBERS) {
            revert TooManyMembers(members_.length, MAX_MEMBERS);
        }

        if (members_[0] == address(0)) {
            revert ZeroMember();
        }

        uint256 unordered = _firstUnordered(members_);
        if (unordered < members_.length) {
            revert MembersNotAscending(unordered);
        }
    }

    function _firstUnordered(address[] memory members_) private pure returns (uint256) {
        for (uint256 i = 1; i < members_.length; i++) {
            if (members_[i] <= members_[i - 1]) {
                return i;
            }
        }

        return members_.length;
    }

    function _recover(bytes32 digest, bytes calldata signature) private pure returns (address) {
        (address recovered, ECDSA.RecoverError failure, bytes32 detail) = ECDSA.tryRecoverCalldata(digest, signature);
        if (failure != ECDSA.RecoverError.NoError || detail != bytes32(0)) {
            return address(0);
        }

        return recovered;
    }
}
