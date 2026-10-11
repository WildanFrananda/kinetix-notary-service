// SPDX-License-Identifier: UNLICENSED
pragma solidity 0.8.37;

import {Test} from "forge-std/src/Test.sol";

import {DeployLiabilityLedger} from "../script/DeployLiabilityLedger.s.sol";
import {LedgerDeployment} from "../script/LedgerDeployment.sol";
import {LiabilityLedger} from "../src/LiabilityLedger.sol";

contract DeployLiabilityLedgerTest is Test {
    address private admin;
    address private anchorer;
    DeployLiabilityLedger private script;

    function setUp() public {
        admin = makeAddr("admin");
        anchorer = makeAddr("anchorer");
        script = new DeployLiabilityLedger();
    }

    function test_DeploysAFreshLedger() public {
        LiabilityLedger ledger = script.deploy(_fresh(keccak256("salt")));

        assertEq(ledger.defaultAdmin(), admin);
        assertEq(ledger.defaultAdminDelay(), 3 days);
        assertTrue(ledger.hasRole(ledger.ANCHORER_ROLE(), anchorer));
        assertEq(ledger.genesisEpoch(), 0);
        assertEq(ledger.genesisDigest(), bytes32(0));
        assertEq(ledger.nextEpoch(), 0);
        assertEq(ledger.lastAsOf(), 0);
        assertFalse(ledger.isSealed());
    }

    function test_DeploysASuccessorThatContinuesAChain() public {
        bytes32 head = keccak256("sealed head");
        LedgerDeployment memory deployment = _fresh(keccak256("salt"));
        deployment.genesisEpoch = 42;
        deployment.genesisDigest = head;
        deployment.genesisAsOf = 1_760_000_000;

        LiabilityLedger ledger = script.deploy(deployment);

        assertEq(ledger.genesisEpoch(), 42);
        assertEq(ledger.head(), head);
        assertEq(ledger.nextEpoch(), 42);
        assertEq(ledger.lastAsOf(), 1_760_000_000);
    }

    function test_TheSaltDecidesTheAddress() public {
        address first = address(script.deploy(_fresh(keccak256("first"))));
        address second = address(script.deploy(_fresh(keccak256("second"))));

        assertNotEq(first, second);
    }

    function test_ReadsEveryValueFromTheEnvironmentAndRefusesWhatDoesNotFit() public {
        vm.setEnv("KINETIX_NOTARY_ADMIN_ADDRESS", vm.toString(admin));
        vm.setEnv("KINETIX_NOTARY_ADMIN_TRANSFER_DELAY_SECONDS", "259200");
        vm.setEnv("KINETIX_NOTARY_ANCHORER_ADDRESS", vm.toString(anchorer));
        vm.setEnv("KINETIX_NOTARY_GENESIS_EPOCH", "7");
        vm.setEnv("KINETIX_NOTARY_GENESIS_DIGEST", vm.toString(keccak256("head")));
        vm.setEnv("KINETIX_NOTARY_GENESIS_AS_OF", "1760000000");
        vm.setEnv("KINETIX_NOTARY_DEPLOY_SALT", vm.toString(keccak256("salt")));

        LedgerDeployment memory read = script.fromEnvironment();

        assertEq(read.admin, admin);
        assertEq(read.adminTransferDelay, 259_200);
        assertEq(read.anchorer, anchorer);
        assertEq(read.genesisEpoch, 7);
        assertEq(read.genesisDigest, keccak256("head"));
        assertEq(read.genesisAsOf, 1_760_000_000);
        assertEq(read.salt, keccak256("salt"));

        vm.setEnv("KINETIX_NOTARY_ADMIN_TRANSFER_DELAY_SECONDS", vm.toString(uint256(type(uint48).max) + 1));
        vm.expectRevert();
        script.fromEnvironment();
    }

    function _fresh(bytes32 salt) private view returns (LedgerDeployment memory) {
        return LedgerDeployment({
            admin: admin,
            adminTransferDelay: 3 days,
            anchorer: anchorer,
            genesisEpoch: 0,
            genesisDigest: bytes32(0),
            genesisAsOf: 0,
            salt: salt
        });
    }
}
