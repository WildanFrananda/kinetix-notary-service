// SPDX-License-Identifier: MIT
pragma solidity 0.8.37;

import {CourtDeployment} from "../../script/CourtDeployment.sol";
import {Create2Factory} from "../../script/Create2Factory.sol";
import {DeployStatementCourt} from "../../script/DeployStatementCourt.s.sol";
import {ILiabilityLedger} from "../../src/ILiabilityLedger.sol";
import {StatementCourt} from "../../src/court/StatementCourt.sol";
import {LedgerFixture} from "../support/LedgerFixture.sol";

contract DeployStatementCourtTest is LedgerFixture {
    address private signer;
    DeployStatementCourt private script;

    function setUp() public override {
        super.setUp();
        signer = makeAddr("statement-signer");
        script = new DeployStatementCourt();
    }

    function test_DeploysACourtBesideTheLedger() public {
        StatementCourt court = script.deploy(_fresh(keccak256("salt")));

        assertEq(court.defaultAdmin(), admin);
        assertEq(court.defaultAdminDelay(), 3 days);
        assertEq(address(court.ledger()), address(ledger));
        assertEq(court.signerOf(0), signer);
        assertEq(court.signerChangeDelay(), 7 days);
        assertEq(court.totalFaults(), 0);
    }

    function test_DeploysThroughTheFactoryWhereTheSaltSays() public {
        CourtDeployment memory deployment = _fresh(keccak256("salt"));
        address expected = vm.computeCreate2Address(keccak256("salt"), keccak256(script.initCode(deployment)));

        assertEq(script.addressOf(deployment), expected);
        assertEq(address(script.deploy(deployment)), expected);
    }

    function test_TheSaltDecidesTheAddress() public {
        address first = address(script.deploy(_fresh(keccak256("first"))));
        address second = address(script.deploy(_fresh(keccak256("second"))));

        assertNotEq(first, second);
    }

    function test_RefusesASecondCourtAtTheSameAddress() public {
        CourtDeployment memory deployment = _fresh(keccak256("salt"));
        script.deploy(deployment);
        address taken = script.addressOf(deployment);

        vm.expectRevert(abi.encodeWithSelector(Create2Factory.DeploymentFailed.selector, taken));
        script.deploy(deployment);
    }

    function test_RefusesWithoutTheFactory() public {
        CourtDeployment memory deployment = _fresh(keccak256("salt"));
        vm.etch(Create2Factory.FACTORY, "");

        vm.expectRevert(abi.encodeWithSelector(Create2Factory.FactoryMissing.selector, Create2Factory.FACTORY));
        script.deploy(deployment);
    }

    function test_RefusesToJudgeSomethingThatIsNotALedger() public {
        CourtDeployment memory deployment = _fresh(keccak256("salt"));
        deployment.ledger = ILiabilityLedger(stranger);
        address expected = script.addressOf(deployment);

        vm.expectRevert(abi.encodeWithSelector(Create2Factory.DeploymentFailed.selector, expected));
        script.deploy(deployment);
    }

    function test_ReadsEveryValueFromTheEnvironmentAndRefusesWhatDoesNotFit() public {
        vm.setEnv("KINETIX_NOTARY_COURT_ADMIN_ADDRESS", vm.toString(admin));
        vm.setEnv("KINETIX_NOTARY_COURT_ADMIN_TRANSFER_DELAY_SECONDS", "259200");
        vm.setEnv("KINETIX_NOTARY_LEDGER_ADDRESS", vm.toString(address(ledger)));
        vm.setEnv("KINETIX_NOTARY_STATEMENT_SIGNER_ADDRESS", vm.toString(signer));
        vm.setEnv("KINETIX_NOTARY_SIGNER_CHANGE_DELAY_SECONDS", "604800");
        vm.setEnv("KINETIX_NOTARY_COURT_DEPLOY_SALT", vm.toString(keccak256("salt")));

        CourtDeployment memory read = script.fromEnvironment();

        assertEq(read.admin, admin);
        assertEq(read.adminTransferDelay, 259_200);
        assertEq(address(read.ledger), address(ledger));
        assertEq(read.signer, signer);
        assertEq(read.signerChangeDelay, 604_800);
        assertEq(read.salt, keccak256("salt"));

        vm.setEnv("KINETIX_NOTARY_SIGNER_CHANGE_DELAY_SECONDS", vm.toString(uint256(type(uint64).max) + 1));
        vm.expectRevert();
        script.fromEnvironment();
    }

    function _fresh(bytes32 salt) private view returns (CourtDeployment memory) {
        return CourtDeployment({
            admin: admin,
            adminTransferDelay: 3 days,
            ledger: ledger,
            signer: signer,
            signerChangeDelay: 7 days,
            salt: salt
        });
    }
}
