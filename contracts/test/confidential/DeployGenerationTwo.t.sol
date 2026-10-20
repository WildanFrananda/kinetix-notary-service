// SPDX-License-Identifier: MIT
pragma solidity 0.8.37;

import {CouncilDeployment} from "../../script/CouncilDeployment.sol";
import {Create2Factory} from "../../script/Create2Factory.sol";
import {DeployConfidentialAnchorCouncil} from "../../script/DeployConfidentialAnchorCouncil.s.sol";
import {DeployLiabilityLedgerV2} from "../../script/DeployLiabilityLedgerV2.s.sol";
import {LedgerV2Deployment} from "../../script/LedgerV2Deployment.sol";
import {ILiabilityLedger} from "../../src/ILiabilityLedger.sol";
import {ConfidentialAnchorCouncil} from "../../src/confidential/ConfidentialAnchorCouncil.sol";
import {LiabilityLedgerV2} from "../../src/confidential/LiabilityLedgerV2.sol";
import {PlonkVerifier} from "../../src/confidential/PlonkVerifier.sol";
import {LedgerFixture} from "../support/LedgerFixture.sol";

contract DeployGenerationTwoTest is LedgerFixture {
    DeployLiabilityLedgerV2 private ledgerScript;
    DeployConfidentialAnchorCouncil private councilScript;

    function setUp() public override {
        super.setUp();
        ledgerScript = new DeployLiabilityLedgerV2();
        councilScript = new DeployConfidentialAnchorCouncil();
        _anchor(_header(0, bytes32(0), _tree(0, 1, 3), 3, _now() - 60));
        vm.warp(vm.getBlockTimestamp() + 1 hours);
    }

    function test_RefusesWhileThePredecessorStillAnchors() public {
        LedgerV2Deployment memory deployment = _deployment(keccak256("salt"));

        vm.expectRevert(abi.encodeWithSelector(DeployLiabilityLedgerV2.PredecessorNotSealed.selector, address(ledger)));
        ledgerScript.deploy(deployment);
    }

    function test_DeploysTheVerifierAndALedgerThatStartsAtTheSealedHead() public {
        _sealLedger();
        LedgerV2Deployment memory deployment = _deployment(keccak256("salt"));
        address expectedLedger = ledgerScript.addressOf(deployment);
        address expectedVerifier = ledgerScript.verifierAddressOf(deployment);

        LiabilityLedgerV2 next = ledgerScript.deploy(deployment);

        assertEq(address(next), expectedLedger);
        assertEq(next.verifier(), expectedVerifier);
        assertEq(keccak256(expectedVerifier.code), keccak256(type(PlonkVerifier).runtimeCode));
        assertEq(next.genesisEpoch(), ledger.nextEpoch());
        assertEq(next.genesisDigest(), ledger.head());
        assertEq(next.lastAsOf(), ledger.lastAsOf());
        assertEq(next.depth(), 8);
        assertEq(next.defaultAdmin(), admin);
        assertTrue(next.hasRole(next.ANCHORER_ROLE(), anchorer));
        assertEq(next.scheme(), 2);
    }

    function test_RefusesASecondGenerationAtTheSameAddresses() public {
        _sealLedger();
        LedgerV2Deployment memory deployment = _deployment(keccak256("salt"));
        ledgerScript.deploy(deployment);
        address taken = ledgerScript.verifierAddressOf(deployment);

        vm.expectRevert(abi.encodeWithSelector(Create2Factory.DeploymentFailed.selector, taken));
        ledgerScript.deploy(deployment);
    }

    function test_RefusesWithoutTheFactory() public {
        _sealLedger();
        LedgerV2Deployment memory deployment = _deployment(keccak256("salt"));
        vm.etch(Create2Factory.FACTORY, "");

        vm.expectRevert(abi.encodeWithSelector(Create2Factory.FactoryMissing.selector, Create2Factory.FACTORY));
        ledgerScript.deploy(deployment);
    }

    function test_DeploysTheSchemeTwoCouncilBesideTheNewLedger() public {
        _sealLedger();
        LiabilityLedgerV2 next = ledgerScript.deploy(_deployment(keccak256("salt")));
        CouncilDeployment memory deployment = _council(next, keccak256("council"));

        ConfidentialAnchorCouncil council = councilScript.deploy(deployment);

        assertEq(address(council), councilScript.addressOf(deployment));
        assertEq(address(council.ledger()), address(next));
        assertEq(council.threshold(), 2);
    }

    function test_RefusesASchemeTwoCouncilBesideASchemeOneLedger() public {
        CouncilDeployment memory deployment = _council(ledger, keccak256("council"));
        address expected = councilScript.addressOf(deployment);

        vm.expectRevert(abi.encodeWithSelector(Create2Factory.DeploymentFailed.selector, expected));
        councilScript.deploy(deployment);
    }

    function _sealLedger() private {
        bytes32 head = ledger.head();
        vm.prank(admin);
        ledger.seal(head);
    }

    function _deployment(bytes32 salt) private view returns (LedgerV2Deployment memory) {
        return LedgerV2Deployment({
            admin: admin,
            adminTransferDelay: ADMIN_DELAY,
            anchorer: anchorer,
            predecessor: ledger,
            depth: 8,
            verifierSalt: keccak256(abi.encode(salt, "verifier")),
            salt: salt
        });
    }

    function _council(ILiabilityLedger target, bytes32 salt) private returns (CouncilDeployment memory) {
        return CouncilDeployment({
            admin: admin,
            adminTransferDelay: ADMIN_DELAY,
            ledger: target,
            members: councilScript.ascending(_members()),
            threshold: 2,
            membershipChangeDelay: 7 days,
            salt: salt
        });
    }

    function _members() private returns (address[] memory) {
        address[] memory members = new address[](2);
        members[0] = makeAddr("council-notary");
        members[1] = makeAddr("council-watch");
        return members;
    }
}
