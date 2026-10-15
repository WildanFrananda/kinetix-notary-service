// SPDX-License-Identifier: MIT
pragma solidity 0.8.37;

import {Create2Factory} from "../../script/Create2Factory.sol";
import {DeployLedgerDirectory} from "../../script/DeployLedgerDirectory.s.sol";
import {DirectoryDeployment} from "../../script/DirectoryDeployment.sol";
import {IStatementCourt} from "../../src/court/IStatementCourt.sol";
import {Generation} from "../../src/directory/Generation.sol";
import {LedgerDirectory} from "../../src/directory/LedgerDirectory.sol";
import {WithholdingFixture} from "../support/WithholdingFixture.sol";

contract DeployLedgerDirectoryTest is WithholdingFixture {
    DeployLedgerDirectory private script;

    function setUp() public override {
        super.setUp();
        script = new DeployLedgerDirectory();
    }

    function test_DeploysADirectoryWhoseFirstGenerationIsTheLedgerAndItsCourts() public {
        LedgerDirectory directory = script.deploy(_fresh(keccak256("salt")));
        Generation memory zero = directory.generation(0);

        assertEq(directory.defaultAdmin(), admin);
        assertEq(directory.defaultAdminDelay(), 3 days);
        assertEq(directory.foreignSuccessorDelay(), 7 days);
        assertEq(directory.generationCount(), 1);
        assertEq(zero.ledger, address(ledger));
        assertEq(zero.court, address(court));
        assertEq(zero.withholding, address(withholding));
    }

    function test_DeploysThroughTheFactoryWhereTheSaltSays() public {
        DirectoryDeployment memory deployment = _fresh(keccak256("salt"));
        address expected = vm.computeCreate2Address(keccak256("salt"), keccak256(script.initCode(deployment)));

        assertEq(script.addressOf(deployment), expected);
        assertEq(address(script.deploy(deployment)), expected);
    }

    function test_TheSaltDecidesTheAddress() public {
        address first = address(script.deploy(_fresh(keccak256("first"))));
        address second = address(script.deploy(_fresh(keccak256("second"))));

        assertNotEq(first, second);
    }

    function test_RefusesASecondDirectoryAtTheSameAddress() public {
        DirectoryDeployment memory deployment = _fresh(keccak256("salt"));
        script.deploy(deployment);
        address taken = script.addressOf(deployment);

        vm.expectRevert(abi.encodeWithSelector(Create2Factory.DeploymentFailed.selector, taken));
        script.deploy(deployment);
    }

    function test_RefusesWithoutTheFactory() public {
        DirectoryDeployment memory deployment = _fresh(keccak256("salt"));
        vm.etch(Create2Factory.FACTORY, "");

        vm.expectRevert(abi.encodeWithSelector(Create2Factory.FactoryMissing.selector, Create2Factory.FACTORY));
        script.deploy(deployment);
    }

    function test_RefusesACourtThatJudgesAnotherLedger() public {
        DirectoryDeployment memory mismatched = _fresh(keccak256("salt"));
        mismatched.court = IStatementCourt(stranger);
        address expected = script.addressOf(mismatched);

        vm.expectRevert(abi.encodeWithSelector(Create2Factory.DeploymentFailed.selector, expected));
        script.deploy(mismatched);
    }

    function test_ReadsEveryValueFromTheEnvironmentAndRefusesWhatDoesNotFit() public {
        vm.setEnv("KINETIX_NOTARY_DIRECTORY_ADMIN_ADDRESS", vm.toString(admin));
        vm.setEnv("KINETIX_NOTARY_DIRECTORY_ADMIN_TRANSFER_DELAY_SECONDS", "259200");
        vm.setEnv("KINETIX_NOTARY_LEDGER_ADDRESS", vm.toString(address(ledger)));
        vm.setEnv("KINETIX_NOTARY_COURT_ADDRESS", vm.toString(address(court)));
        vm.setEnv("KINETIX_NOTARY_WITHHOLDING_COURT_ADDRESS", vm.toString(address(withholding)));
        vm.setEnv("KINETIX_NOTARY_FOREIGN_SUCCESSOR_DELAY_SECONDS", "604800");
        vm.setEnv("KINETIX_NOTARY_DIRECTORY_DEPLOY_SALT", vm.toString(keccak256("salt")));

        DirectoryDeployment memory read = script.fromEnvironment();

        assertEq(read.admin, admin);
        assertEq(read.adminTransferDelay, 259_200);
        assertEq(address(read.ledger), address(ledger));
        assertEq(address(read.court), address(court));
        assertEq(address(read.withholding), address(withholding));
        assertEq(read.foreignSuccessorDelay, 604_800);
        assertEq(read.salt, keccak256("salt"));

        vm.setEnv("KINETIX_NOTARY_DIRECTORY_ADMIN_TRANSFER_DELAY_SECONDS", "281474976710656");
        vm.expectRevert();
        script.fromEnvironment();
    }

    function _fresh(bytes32 salt) private view returns (DirectoryDeployment memory) {
        return DirectoryDeployment({
            admin: admin,
            adminTransferDelay: 3 days,
            ledger: ledger,
            court: court,
            withholding: withholding,
            foreignSuccessorDelay: 7 days,
            salt: salt
        });
    }
}
