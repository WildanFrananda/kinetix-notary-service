// SPDX-License-Identifier: MIT
pragma solidity 0.8.37;

import {EpochHeader} from "../../src/EpochHeader.sol";
import {AnchorCouncil} from "../../src/council/AnchorCouncil.sol";
import {LedgerFixture} from "./LedgerFixture.sol";

abstract contract CouncilFixture is LedgerFixture {
    uint64 internal constant MEMBERSHIP_DELAY = 7 days;

    address internal notary;
    uint256 internal notaryKey;
    address internal watch;
    uint256 internal watchKey;
    AnchorCouncil internal council;

    function setUp() public virtual override {
        super.setUp();
        (notary, notaryKey) = makeAddrAndKey("council-notary");
        (watch, watchKey) = makeAddrAndKey("council-watch");
        council = new AnchorCouncil(admin, ADMIN_DELAY, ledger, _pair(notary, watch), 2, MEMBERSHIP_DELAY);
        _handOver(address(council));
    }

    function _handOver(address to) internal {
        bytes32 role = ledger.ANCHORER_ROLE();
        vm.prank(admin);
        ledger.grantRole(role, to);
    }

    function _pair(address a, address b) internal pure returns (address[] memory) {
        address[] memory pair = new address[](2);
        (pair[0], pair[1]) = a < b ? (a, b) : (b, a);
        return pair;
    }

    function _next(uint256 seed, uint256 count) internal view returns (EpochHeader memory) {
        uint64 number = ledger.nextEpoch();
        return _header(number, ledger.head(), _tree(number, seed, count), count, _now() - 1);
    }

    function _approve(uint256 key, EpochHeader memory header) internal view returns (bytes memory) {
        bytes32 digest = council.approvalDigest(header.number, council.epochDigest(header));
        (uint8 v, bytes32 r, bytes32 s) = vm.sign(key, digest);
        return abi.encodePacked(r, s, v);
    }

    function _approvals(EpochHeader memory header) internal view returns (bytes[] memory) {
        bytes[] memory approvals = new bytes[](2);
        bytes memory byNotary = _approve(notaryKey, header);
        bytes memory byWatch = _approve(watchKey, header);
        (approvals[0], approvals[1]) = notary < watch ? (byNotary, byWatch) : (byWatch, byNotary);
        return approvals;
    }

    function _anchorNext(uint256 seed, uint256 count) internal returns (EpochHeader memory) {
        EpochHeader memory header = _next(seed, count);
        bytes[] memory approvals = _approvals(header);
        council.anchor(header, approvals);
        vm.warp(vm.getBlockTimestamp() + 1 hours);
        return header;
    }
}
