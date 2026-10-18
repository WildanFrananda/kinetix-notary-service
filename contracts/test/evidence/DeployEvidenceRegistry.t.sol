// SPDX-License-Identifier: MIT
pragma solidity 0.8.37;

import {Create2Factory} from "../../script/Create2Factory.sol";
import {DeployEvidenceRegistry} from "../../script/DeployEvidenceRegistry.s.sol";
import {EvidenceDeployment} from "../../script/EvidenceDeployment.sol";
import {EvidenceRegistry} from "../../src/evidence/EvidenceRegistry.sol";
import {IEvidenceRegistry} from "../../src/evidence/IEvidenceRegistry.sol";
import {EvidenceFixture} from "../support/EvidenceFixture.sol";

contract DeployEvidenceRegistryTest is EvidenceFixture {
    DeployEvidenceRegistry private script;

    function setUp() public override {
        super.setUp();
        script = new DeployEvidenceRegistry();
    }

    function test_DeploysAFreshRegistryOfTheThreeStreams() public {
        EvidenceRegistry fresh = script.deploy(_deployment(keccak256("salt"), IEvidenceRegistry(address(0))));
        uint8[] memory ids = fresh.streams();

        assertEq(ids.length, 3);
        assertEq(ids[0], 1);
        assertEq(ids[1], 2);
        assertEq(ids[2], 3);
        assertEq(fresh.defaultAdmin(), admin);
        assertEq(fresh.defaultAdminDelay(), 3 days);
        assertTrue(fresh.hasRole(fresh.ANCHORER_ROLE(), anchorer));
        assertEq(fresh.stateOf(2).nextBatch, 0);
    }

    function test_DeploysThroughTheFactoryWhereTheSaltSays() public {
        EvidenceDeployment memory deployment = _deployment(keccak256("salt"), IEvidenceRegistry(address(0)));
        address expected = vm.computeCreate2Address(keccak256("salt"), keccak256(script.initCode(deployment)));

        assertEq(script.addressOf(deployment), expected);
        assertEq(address(script.deploy(deployment)), expected);
    }

    function test_RefusesASecondRegistryAtTheSameAddress() public {
        EvidenceDeployment memory deployment = _deployment(keccak256("salt"), IEvidenceRegistry(address(0)));
        script.deploy(deployment);
        address taken = script.addressOf(deployment);

        vm.expectRevert(abi.encodeWithSelector(Create2Factory.DeploymentFailed.selector, taken));
        script.deploy(deployment);
    }

    function test_RefusesAZeroAnchorerThroughTheFactory() public {
        EvidenceDeployment memory deployment = _deployment(keccak256("salt"), IEvidenceRegistry(address(0)));
        deployment.anchorer = address(0);
        address expected = script.addressOf(deployment);

        vm.expectRevert(abi.encodeWithSelector(Create2Factory.DeploymentFailed.selector, expected));
        script.deploy(deployment);
    }

    function test_ASuccessorTakesEachStreamFromTheSealedPredecessor() public {
        _anchorNext(DELIVERY, 1, 3);
        _anchorNext(DELIVERY, 2, 3);
        _anchorNext(REVIEW, 3, 1);
        bytes32[] memory heads = new bytes32[](3);
        heads[0] = registry.stateOf(LIFECYCLE).head;
        heads[1] = registry.stateOf(DELIVERY).head;
        heads[2] = registry.stateOf(REVIEW).head;
        vm.prank(admin);
        registry.seal(heads);

        EvidenceRegistry successor = script.deploy(_deployment(keccak256("next"), registry));

        assertEq(successor.stateOf(DELIVERY).nextBatch, 2);
        assertEq(successor.stateOf(DELIVERY).head, heads[1]);
        assertEq(successor.stateOf(DELIVERY).throughSequence, registry.stateOf(DELIVERY).throughSequence);
        assertEq(successor.stateOf(REVIEW).nextBatch, 1);
        assertEq(successor.stateOf(LIFECYCLE).nextBatch, 0);
    }

    function test_RefusesToContinueARegistryThatIsNotSealed() public {
        EvidenceDeployment memory deployment = _deployment(keccak256("next"), registry);

        vm.expectRevert(abi.encodeWithSelector(DeployEvidenceRegistry.PredecessorNotSealed.selector, address(registry)));
        script.genesisOf(deployment);
    }

    function test_ReadsEveryValueFromTheEnvironment() public {
        vm.setEnv("KINETIX_NOTARY_EVIDENCE_ADMIN_ADDRESS", vm.toString(admin));
        vm.setEnv("KINETIX_NOTARY_EVIDENCE_ADMIN_TRANSFER_DELAY_SECONDS", "259200");
        vm.setEnv("KINETIX_NOTARY_EVIDENCE_ANCHORER_ADDRESS", vm.toString(anchorer));
        vm.setEnv("KINETIX_NOTARY_EVIDENCE_DEPLOY_SALT", vm.toString(keccak256("salt")));

        EvidenceDeployment memory read = script.fromEnvironment();

        assertEq(read.admin, admin);
        assertEq(read.adminTransferDelay, 259_200);
        assertEq(read.anchorer, anchorer);
        assertEq(address(read.predecessor), address(0));
        assertEq(read.salt, keccak256("salt"));

        vm.setEnv("KINETIX_NOTARY_EVIDENCE_ADMIN_TRANSFER_DELAY_SECONDS", "281474976710656");
        vm.expectRevert();
        script.fromEnvironment();
    }

    function _deployment(bytes32 salt, IEvidenceRegistry predecessor) private view returns (EvidenceDeployment memory) {
        return EvidenceDeployment({
            admin: admin,
            adminTransferDelay: 3 days,
            anchorer: anchorer,
            predecessor: predecessor,
            salt: salt
        });
    }
}
