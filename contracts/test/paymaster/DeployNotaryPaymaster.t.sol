// SPDX-License-Identifier: MIT
pragma solidity 0.8.37;

import {Create2Factory} from "../../script/Create2Factory.sol";
import {DeployNotaryPaymaster} from "../../script/DeployNotaryPaymaster.s.sol";
import {PaymasterDeployment} from "../../script/PaymasterDeployment.sol";
import {IStatementCourt} from "../../src/court/IStatementCourt.sol";
import {NotaryPaymaster} from "../../src/paymaster/NotaryPaymaster.sol";
import {PaymasterFixture} from "../support/PaymasterFixture.sol";

contract DeployNotaryPaymasterTest is PaymasterFixture {
    DeployNotaryPaymaster private script;

    function setUp() public override {
        super.setUp();
        script = new DeployNotaryPaymaster();
    }

    function test_DeploysAPaymasterBesideTheCourts() public {
        NotaryPaymaster deployed = script.deploy(_fresh(keccak256("salt")));

        assertEq(deployed.defaultAdmin(), admin);
        assertEq(deployed.defaultAdminDelay(), 3 days);
        assertEq(address(deployed.entryPoint()), CANONICAL_ENTRY_POINT);
        assertEq(deployed.account(), address(delegate));
        assertEq(address(deployed.court()), address(court));
        assertEq(address(deployed.withholding()), address(withholding));
        assertEq(deployed.statementSigner(), signer);
        assertEq(deployed.maxCostPerOperation(), CEILING);
        assertEq(deployed.unstakeDelay(), UNSTAKE_DELAY);
    }

    function test_DeploysThroughTheFactoryWhereTheSaltSays() public {
        PaymasterDeployment memory deployment = _fresh(keccak256("salt"));
        address expected = vm.computeCreate2Address(keccak256("salt"), keccak256(script.initCode(deployment)));

        assertEq(script.addressOf(deployment), expected);
        assertEq(address(script.deploy(deployment)), expected);
    }

    function test_RefusesASecondPaymasterAtTheSameAddress() public {
        PaymasterDeployment memory deployment = _fresh(keccak256("salt"));
        script.deploy(deployment);
        address taken = script.addressOf(deployment);

        vm.expectRevert(abi.encodeWithSelector(Create2Factory.DeploymentFailed.selector, taken));
        script.deploy(deployment);
    }

    function test_RefusesToStandBesideSomethingThatIsNotACourt() public {
        PaymasterDeployment memory deployment = _fresh(keccak256("salt"));
        deployment.court = IStatementCourt(stranger);
        address expected = script.addressOf(deployment);

        vm.expectRevert(abi.encodeWithSelector(Create2Factory.DeploymentFailed.selector, expected));
        script.deploy(deployment);
    }

    function test_ReadsEveryValueFromTheEnvironmentAndRefusesWhatDoesNotFit() public {
        vm.setEnv("KINETIX_NOTARY_PAYMASTER_ADMIN_ADDRESS", vm.toString(admin));
        vm.setEnv("KINETIX_NOTARY_PAYMASTER_ADMIN_TRANSFER_DELAY_SECONDS", "259200");
        vm.setEnv("KINETIX_NOTARY_ENTRY_POINT_ADDRESS", vm.toString(CANONICAL_ENTRY_POINT));
        vm.setEnv("KINETIX_NOTARY_7702_ACCOUNT_ADDRESS", vm.toString(address(delegate)));
        vm.setEnv("KINETIX_NOTARY_COURT_ADDRESS", vm.toString(address(court)));
        vm.setEnv("KINETIX_NOTARY_WITHHOLDING_COURT_ADDRESS", vm.toString(address(withholding)));
        vm.setEnv("KINETIX_NOTARY_PAYMASTER_MAX_COST_WEI", "100000000000000");
        vm.setEnv("KINETIX_NOTARY_PAYMASTER_UNSTAKE_DELAY_SECONDS", "86400");
        vm.setEnv("KINETIX_NOTARY_PAYMASTER_DEPLOY_SALT", vm.toString(keccak256("salt")));

        PaymasterDeployment memory read = script.fromEnvironment();

        assertEq(read.admin, admin);
        assertEq(read.adminTransferDelay, 259_200);
        assertEq(address(read.entryPoint), CANONICAL_ENTRY_POINT);
        assertEq(read.account, address(delegate));
        assertEq(address(read.court), address(court));
        assertEq(address(read.withholding), address(withholding));
        assertEq(read.maxCostPerOperation, 100_000_000_000_000);
        assertEq(read.unstakeDelay, 86_400);
        assertEq(read.salt, keccak256("salt"));

        vm.setEnv("KINETIX_NOTARY_PAYMASTER_UNSTAKE_DELAY_SECONDS", "4294967296");
        vm.expectRevert();
        script.fromEnvironment();

        vm.setEnv("KINETIX_NOTARY_PAYMASTER_UNSTAKE_DELAY_SECONDS", "86400");
        vm.setEnv("KINETIX_NOTARY_PAYMASTER_ADMIN_TRANSFER_DELAY_SECONDS", "281474976710656");
        vm.expectRevert();
        script.fromEnvironment();
    }

    function _fresh(bytes32 salt) private view returns (PaymasterDeployment memory) {
        return PaymasterDeployment({
            admin: admin,
            adminTransferDelay: 3 days,
            entryPoint: entryPoint,
            account: address(delegate),
            court: court,
            withholding: withholding,
            maxCostPerOperation: CEILING,
            unstakeDelay: UNSTAKE_DELAY,
            salt: salt
        });
    }
}
