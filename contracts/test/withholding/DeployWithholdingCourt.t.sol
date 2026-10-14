// SPDX-License-Identifier: MIT
pragma solidity 0.8.37;

import {Create2Factory} from "../../script/Create2Factory.sol";
import {DeployWithholdingCourt} from "../../script/DeployWithholdingCourt.s.sol";
import {WithholdingDeployment} from "../../script/WithholdingDeployment.sol";
import {IStatementCourt} from "../../src/court/IStatementCourt.sol";
import {WithholdingCourt} from "../../src/withholding/WithholdingCourt.sol";
import {CourtFixture} from "../support/CourtFixture.sol";

contract DeployWithholdingCourtTest is CourtFixture {
    DeployWithholdingCourt private script;

    function setUp() public override {
        super.setUp();
        script = new DeployWithholdingCourt();
    }

    function test_DeploysAWithholdingCourtBesideTheCourt() public {
        WithholdingCourt withholding = script.deploy(_fresh(keccak256("salt")));

        assertEq(address(withholding.court()), address(court));
        assertEq(address(withholding.ledger()), address(ledger));
        assertEq(withholding.requestWindow(), 7 days);
        assertEq(withholding.answerPeriod(), 24 hours);
        assertEq(withholding.totalSilences(), 0);
    }

    function test_DeploysThroughTheFactoryWhereTheSaltSays() public {
        WithholdingDeployment memory deployment = _fresh(keccak256("salt"));
        address expected = vm.computeCreate2Address(keccak256("salt"), keccak256(script.initCode(deployment)));

        assertEq(script.addressOf(deployment), expected);
        assertEq(address(script.deploy(deployment)), expected);
    }

    function test_TheSaltDecidesTheAddress() public {
        address first = address(script.deploy(_fresh(keccak256("first"))));
        address second = address(script.deploy(_fresh(keccak256("second"))));

        assertNotEq(first, second);
    }

    function test_RefusesASecondWithholdingCourtAtTheSameAddress() public {
        WithholdingDeployment memory deployment = _fresh(keccak256("salt"));
        script.deploy(deployment);
        address taken = script.addressOf(deployment);

        vm.expectRevert(abi.encodeWithSelector(Create2Factory.DeploymentFailed.selector, taken));
        script.deploy(deployment);
    }

    function test_RefusesWithoutTheFactory() public {
        WithholdingDeployment memory deployment = _fresh(keccak256("salt"));
        vm.etch(Create2Factory.FACTORY, "");

        vm.expectRevert(abi.encodeWithSelector(Create2Factory.FactoryMissing.selector, Create2Factory.FACTORY));
        script.deploy(deployment);
    }

    function test_RefusesToStandBesideSomethingThatIsNotACourt() public {
        WithholdingDeployment memory deployment = _fresh(keccak256("salt"));
        deployment.court = IStatementCourt(address(ledger));
        address expected = script.addressOf(deployment);

        vm.expectRevert(abi.encodeWithSelector(Create2Factory.DeploymentFailed.selector, expected));
        script.deploy(deployment);
    }

    function test_ReadsEveryValueFromTheEnvironmentAndRefusesWhatDoesNotFit() public {
        vm.setEnv("KINETIX_NOTARY_COURT_ADDRESS", vm.toString(address(court)));
        vm.setEnv("KINETIX_NOTARY_WITHHOLDING_REQUEST_WINDOW_SECONDS", "604800");
        vm.setEnv("KINETIX_NOTARY_WITHHOLDING_ANSWER_PERIOD_SECONDS", "86400");
        vm.setEnv("KINETIX_NOTARY_WITHHOLDING_DEPLOY_SALT", vm.toString(keccak256("salt")));

        WithholdingDeployment memory read = script.fromEnvironment();

        assertEq(address(read.court), address(court));
        assertEq(read.requestWindow, 604_800);
        assertEq(read.answerPeriod, 86_400);
        assertEq(read.salt, keccak256("salt"));

        vm.setEnv("KINETIX_NOTARY_WITHHOLDING_ANSWER_PERIOD_SECONDS", vm.toString(uint256(type(uint64).max) + 1));
        vm.expectRevert();
        script.fromEnvironment();
    }

    function _fresh(bytes32 salt) private view returns (WithholdingDeployment memory) {
        return WithholdingDeployment({court: court, requestWindow: 7 days, answerPeriod: 24 hours, salt: salt});
    }
}
