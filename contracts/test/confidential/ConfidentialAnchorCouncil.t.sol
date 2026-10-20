// SPDX-License-Identifier: MIT
pragma solidity 0.8.37;

import {EpochChain} from "../../src/EpochChain.sol";
import {EpochHeader} from "../../src/EpochHeader.sol";
import {LiabilityLedger} from "../../src/LiabilityLedger.sol";
import {ConfidentialAnchorCouncil} from "../../src/confidential/ConfidentialAnchorCouncil.sol";
import {ConfidentialEpochChain} from "../../src/confidential/ConfidentialEpochChain.sol";
import {IConfidentialLedger} from "../../src/confidential/IConfidentialLedger.sol";
import {IAnchorCouncil} from "../../src/council/IAnchorCouncil.sol";
import {ConfidentialLedgerFixture} from "../support/ConfidentialLedgerFixture.sol";
import {ConfidentialStatement} from "../support/ConfidentialStatement.sol";

contract ConfidentialAnchorCouncilTest is ConfidentialLedgerFixture {
    uint64 internal constant MEMBERSHIP_DELAY = 7 days;

    address private notary;
    uint256 private notaryKey;
    address private watch;
    uint256 private watchKey;
    ConfidentialAnchorCouncil private council;

    function setUp() public override {
        super.setUp();
        (notary, notaryKey) = makeAddrAndKey("council-notary");
        (watch, watchKey) = makeAddrAndKey("council-watch");
        council = new ConfidentialAnchorCouncil(admin, ADMIN_DELAY, ledger, _pair(), 2, MEMBERSHIP_DELAY);
        bytes32 role = ledger.ANCHORER_ROLE();
        vm.prank(admin);
        ledger.grantRole(role, address(council));
    }

    function test_RevertWhen_TheLedgerSpeaksSchemeOne() public {
        LiabilityLedger v1 = new LiabilityLedger(admin, ADMIN_DELAY, anchorer, 0, bytes32(0), 0);

        vm.expectRevert(abi.encodeWithSelector(IAnchorCouncil.NotALedger.selector, address(v1)));
        new ConfidentialAnchorCouncil(admin, ADMIN_DELAY, v1, _pair(), 2, MEMBERSHIP_DELAY);
    }

    function test_AnnouncesTheCouncilInterface() public view {
        assertTrue(council.supportsInterface(type(IAnchorCouncil).interfaceId));
    }

    function test_ItsEpochDigestIsTheLedgersScheme() public view {
        EpochHeader memory header = _header(_statement(0, 1, 5), bytes32(0), _now() - 60);

        assertEq(council.epochDigest(header), ConfidentialEpochChain.digest(header));
        assertNotEq(council.epochDigest(header), EpochChain.digest(header));
    }

    function test_AnchorsAProvenEpochOnTwoApprovals() public {
        ConfidentialStatement memory statement = _statement(0, 1, 5);
        _prove(statement);
        EpochHeader memory header = _header(statement, bytes32(0), _now() - 60);
        bytes32 digest = ConfidentialEpochChain.digest(header);
        bytes[] memory approvals = _approvals(header, digest);

        vm.expectEmit(address(council));
        emit IAnchorCouncil.EpochApproved(0, digest, _pair());
        bytes32 anchored = council.anchor(header, approvals);

        assertEq(anchored, digest);
        assertEq(ledger.head(), digest);
        assertEq(ledger.epochAt(0).digest, digest);
    }

    function test_ApprovalsOfTheSchemeOneDigestAreRefused() public {
        ConfidentialStatement memory statement = _statement(0, 1, 5);
        _prove(statement);
        EpochHeader memory header = _header(statement, bytes32(0), _now() - 60);
        bytes[] memory approvals = _approvals(header, EpochChain.digest(header));

        vm.expectRevert();
        council.anchor(header, approvals);

        assertEq(ledger.nextEpoch(), 0);
    }

    function test_TheCouncilCannotAnchorWhatNobodyProved() public {
        EpochHeader memory header = _header(_statement(0, 1, 5), bytes32(0), _now() - 60);
        bytes[] memory approvals = _approvals(header, ConfidentialEpochChain.digest(header));

        vm.expectRevert(
            abi.encodeWithSelector(IConfidentialLedger.EpochNotProven.selector, 0, header.root, header.leafCount)
        );
        council.anchor(header, approvals);
    }

    function _pair() private view returns (address[] memory) {
        address[] memory pair = new address[](2);
        (pair[0], pair[1]) = notary < watch ? (notary, watch) : (watch, notary);
        return pair;
    }

    function _approvals(EpochHeader memory header, bytes32 digest) private view returns (bytes[] memory) {
        bytes32 hash = council.approvalDigest(header.number, digest);
        (uint8 v1, bytes32 r1, bytes32 s1) = vm.sign(notaryKey, hash);
        (uint8 v2, bytes32 r2, bytes32 s2) = vm.sign(watchKey, hash);
        bytes memory byNotary = abi.encodePacked(r1, s1, v1);
        bytes memory byWatch = abi.encodePacked(r2, s2, v2);
        bytes[] memory approvals = new bytes[](2);
        (approvals[0], approvals[1]) = notary < watch ? (byNotary, byWatch) : (byWatch, byNotary);
        return approvals;
    }
}
