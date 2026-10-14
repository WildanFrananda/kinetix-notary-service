// SPDX-License-Identifier: MIT
pragma solidity 0.8.37;

import {CouncilDeployment} from "../../script/CouncilDeployment.sol";
import {Create2Factory} from "../../script/Create2Factory.sol";
import {DeployAnchorCouncil} from "../../script/DeployAnchorCouncil.s.sol";
import {ILiabilityLedger} from "../../src/ILiabilityLedger.sol";
import {AnchorCouncil} from "../../src/council/AnchorCouncil.sol";
import {LedgerFixture} from "../support/LedgerFixture.sol";

contract DeployAnchorCouncilTest is LedgerFixture {
    address private notary;
    address private watch;
    DeployAnchorCouncil private script;

    function setUp() public override {
        super.setUp();
        notary = makeAddr("council-notary");
        watch = makeAddr("council-watch");
        script = new DeployAnchorCouncil();
    }

    function test_DeploysACouncilBesideTheLedger() public {
        AnchorCouncil council = script.deploy(_fresh(keccak256("salt")));

        assertEq(council.defaultAdmin(), admin);
        assertEq(council.defaultAdminDelay(), 3 days);
        assertEq(address(council.ledger()), address(ledger));
        assertEq(council.members(), _members());
        assertEq(council.threshold(), 2);
        assertEq(council.membershipChangeDelay(), 7 days);
    }

    function test_DeploysThroughTheFactoryWhereTheSaltSays() public {
        CouncilDeployment memory deployment = _fresh(keccak256("salt"));
        address expected = vm.computeCreate2Address(keccak256("salt"), keccak256(script.initCode(deployment)));

        assertEq(script.addressOf(deployment), expected);
        assertEq(address(script.deploy(deployment)), expected);
    }

    function test_TheSaltDecidesTheAddress() public {
        address first = address(script.deploy(_fresh(keccak256("first"))));
        address second = address(script.deploy(_fresh(keccak256("second"))));

        assertNotEq(first, second);
    }

    function test_RefusesASecondCouncilAtTheSameAddress() public {
        CouncilDeployment memory deployment = _fresh(keccak256("salt"));
        script.deploy(deployment);
        address taken = script.addressOf(deployment);

        vm.expectRevert(abi.encodeWithSelector(Create2Factory.DeploymentFailed.selector, taken));
        script.deploy(deployment);
    }

    function test_RefusesWithoutTheFactory() public {
        CouncilDeployment memory deployment = _fresh(keccak256("salt"));
        vm.etch(Create2Factory.FACTORY, "");

        vm.expectRevert(abi.encodeWithSelector(Create2Factory.FactoryMissing.selector, Create2Factory.FACTORY));
        script.deploy(deployment);
    }

    function test_RefusesToStandBesideSomethingThatIsNotALedger() public {
        CouncilDeployment memory notALedger = _fresh(keccak256("salt"));
        notALedger.ledger = ILiabilityLedger(stranger);
        address expected = script.addressOf(notALedger);

        vm.expectRevert(abi.encodeWithSelector(Create2Factory.DeploymentFailed.selector, expected));
        script.deploy(notALedger);
    }

    function test_RefusesACouncilOfOne() public {
        CouncilDeployment memory alone = _fresh(keccak256("salt"));
        alone.threshold = 1;
        address expected = script.addressOf(alone);

        vm.expectRevert(abi.encodeWithSelector(Create2Factory.DeploymentFailed.selector, expected));
        script.deploy(alone);
    }

    function test_SortsTheMembersItIsGiven() public view {
        address[] memory given = new address[](3);
        given[0] = address(0x30);
        given[1] = address(0x10);
        given[2] = address(0x20);

        address[] memory sorted = script.ascending(given);

        assertEq(sorted[0], address(0x10));
        assertEq(sorted[1], address(0x20));
        assertEq(sorted[2], address(0x30));
    }

    function test_ReadsEveryValueFromTheEnvironmentAndRefusesWhatDoesNotFit() public {
        address[] memory expected = _members();
        vm.setEnv("KINETIX_NOTARY_COUNCIL_ADMIN_ADDRESS", vm.toString(admin));
        vm.setEnv("KINETIX_NOTARY_COUNCIL_ADMIN_TRANSFER_DELAY_SECONDS", "259200");
        vm.setEnv("KINETIX_NOTARY_LEDGER_ADDRESS", vm.toString(address(ledger)));
        vm.setEnv(
            "KINETIX_NOTARY_COUNCIL_MEMBERS", string.concat(vm.toString(expected[1]), ",", vm.toString(expected[0]))
        );
        vm.setEnv("KINETIX_NOTARY_COUNCIL_THRESHOLD", "2");
        vm.setEnv("KINETIX_NOTARY_MEMBERSHIP_CHANGE_DELAY_SECONDS", "604800");
        vm.setEnv("KINETIX_NOTARY_COUNCIL_DEPLOY_SALT", vm.toString(keccak256("salt")));

        CouncilDeployment memory read = script.fromEnvironment();

        assertEq(read.admin, admin);
        assertEq(read.adminTransferDelay, 259_200);
        assertEq(address(read.ledger), address(ledger));
        assertEq(read.members, expected);
        assertEq(read.threshold, 2);
        assertEq(read.membershipChangeDelay, 604_800);
        assertEq(read.salt, keccak256("salt"));

        vm.setEnv("KINETIX_NOTARY_COUNCIL_THRESHOLD", "256");
        vm.expectRevert();
        script.fromEnvironment();
    }

    function _members() private view returns (address[] memory) {
        address[] memory members = new address[](2);
        (members[0], members[1]) = notary < watch ? (notary, watch) : (watch, notary);
        return members;
    }

    function _fresh(bytes32 salt) private view returns (CouncilDeployment memory) {
        return CouncilDeployment({
            admin: admin,
            adminTransferDelay: 3 days,
            ledger: ledger,
            members: _members(),
            threshold: 2,
            membershipChangeDelay: 7 days,
            salt: salt
        });
    }
}
