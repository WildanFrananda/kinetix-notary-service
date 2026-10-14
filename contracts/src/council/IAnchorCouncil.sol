// SPDX-License-Identifier: MIT
pragma solidity 0.8.37;

import {EpochHeader} from "../EpochHeader.sol";
import {ILiabilityLedger} from "../ILiabilityLedger.sol";
import {PendingMembership} from "./PendingMembership.sol";

interface IAnchorCouncil {
    event EpochApproved(uint64 indexed number, bytes32 indexed digest, address[] approvers);
    event MembershipScheduled(address[] members, uint8 threshold, uint64 activatesAt);
    event MembershipCancelled(address[] members, uint8 threshold);
    event MembershipActivated(address[] members, uint8 threshold);

    error NotALedger(address ledger);
    error ZeroDelay();
    error ThresholdTooLow(uint8 threshold, uint8 minimum);
    error TooFewMembers(uint256 count, uint8 threshold);
    error TooManyMembers(uint256 count, uint256 maximum);
    error ZeroMember();
    error MembersNotAscending(uint256 index);
    error SameMembership();
    error MembershipChangePending(uint64 activatesAt);
    error NoMembershipChangePending();
    error MembershipChangeTooEarly(uint64 activatesAt, uint64 blockTime);
    error NotEnoughApprovals(uint256 given, uint8 threshold);
    error NotAMember(address recovered);
    error ApprovalsOutOfOrder(address previous, address next);

    function anchor(EpochHeader calldata header, bytes[] calldata signatures) external returns (bytes32);

    function scheduleMembership(address[] calldata nextMembers, uint8 nextThreshold) external;

    function cancelMembership() external;

    function activateMembership() external;

    function ledger() external view returns (ILiabilityLedger);

    function membershipChangeDelay() external view returns (uint64);

    function members() external view returns (address[] memory);

    function threshold() external view returns (uint8);

    function isMember(address account) external view returns (bool);

    function pendingMembership() external view returns (PendingMembership memory);

    function epochDigest(EpochHeader calldata header) external pure returns (bytes32);

    function approvalDigest(uint64 number, bytes32 digest) external view returns (bytes32);

    function recoverApprover(uint64 number, bytes32 digest, bytes calldata signature) external view returns (address);
}
