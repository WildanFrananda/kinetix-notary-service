// SPDX-License-Identifier: MIT
pragma solidity 0.8.37;

import {SumNode} from "../../src/SumNode.sol";
import {AccountCommitments} from "../../src/withholding/AccountCommitments.sol";
import {Registration} from "../../src/withholding/Registration.sol";
import {WithholdingCourt} from "../../src/withholding/WithholdingCourt.sol";
import {CourtFixture} from "./CourtFixture.sol";
import {LeafSet} from "./LeafSet.sol";
import {TreeBuilder} from "./TreeBuilder.sol";

abstract contract WithholdingFixture is CourtFixture {
    uint64 internal constant WINDOW = 7 days;
    uint64 internal constant PERIOD = 24 hours;

    address internal holder;
    uint256 internal holderKey;
    bytes32 internal account;
    WithholdingCourt internal withholding;

    function setUp() public virtual override {
        super.setUp();
        (holder, holderKey) = makeAddrAndKey("holder");
        account = keccak256("account");
        withholding = new WithholdingCourt(court, WINDOW, PERIOD);
    }

    function _registration(uint64 fromEpoch) internal view returns (Registration memory) {
        return Registration({holder: holder, account: account, fromEpoch: fromEpoch});
    }

    function _signRegistration(uint256 key, Registration memory registration) internal view returns (bytes memory) {
        (uint8 v, bytes32 r, bytes32 s) = vm.sign(key, withholding.registrationDigest(registration));
        return abi.encodePacked(r, s, v);
    }

    function _salt(uint64 epoch) internal pure returns (bytes32) {
        return keccak256(abi.encode("salt", epoch));
    }

    function _accountTree(
        uint64 number,
        uint256 seed,
        uint256 count,
        uint256 index
    ) internal view returns (SumNode[][] memory, uint64) {
        bytes32[] memory commitments = LeafSet.commitments(seed, count);
        uint64[] memory amounts = LeafSet.amounts(seed, count);
        commitments[index] = AccountCommitments.commitment(account, _salt(number));
        return (TreeBuilder.build(number, commitments, amounts), amounts[index]);
    }

    function _anchorWithAccount(uint256 seed, uint256 count, uint256 index) internal returns (uint64, uint64) {
        uint64 number = ledger.nextEpoch();
        (SumNode[][] memory levels, uint64 amount) = _accountTree(number, seed, count, index);
        _anchor(_header(number, ledger.head(), levels, count, _now() - 1));
        vm.warp(vm.getBlockTimestamp() + 1 hours);
        return (number, amount);
    }

    function _anchorWithoutAccount(uint256 seed, uint256 count) internal returns (uint64) {
        uint64 number = ledger.nextEpoch();
        SumNode[][] memory levels =
            TreeBuilder.build(number, LeafSet.commitments(seed, count), LeafSet.amounts(seed, count));
        _anchor(_header(number, ledger.head(), levels, count, _now() - 1));
        vm.warp(vm.getBlockTimestamp() + 1 hours);
        return number;
    }

    function _accountSiblings(
        uint64 number,
        uint256 seed,
        uint256 count,
        uint256 index
    ) internal view returns (SumNode[] memory) {
        (SumNode[][] memory levels,) = _accountTree(number, seed, count, index);
        return TreeBuilder.proof(levels, index);
    }

    function _request(Registration memory registration, uint64 epoch) internal returns (bytes32) {
        bytes memory signature = _signRegistration(signerKey, registration);
        vm.prank(registration.holder);
        return withholding.requestStatement(registration, signature, epoch);
    }
}
