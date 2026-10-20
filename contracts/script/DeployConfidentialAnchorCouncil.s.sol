// SPDX-License-Identifier: MIT
pragma solidity 0.8.37;

import {SafeCast} from "@openzeppelin-contracts/utils/math/SafeCast.sol";
import {Script} from "forge-std/src/Script.sol";

import {ILiabilityLedger} from "../src/ILiabilityLedger.sol";
import {ConfidentialAnchorCouncil} from "../src/confidential/ConfidentialAnchorCouncil.sol";
import {CouncilDeployment} from "./CouncilDeployment.sol";
import {Create2Factory} from "./Create2Factory.sol";

contract DeployConfidentialAnchorCouncil is Script {
    function run() external returns (ConfidentialAnchorCouncil) {
        return deploy(fromEnvironment());
    }

    function fromEnvironment() public view returns (CouncilDeployment memory) {
        return CouncilDeployment({
            admin: vm.envAddress("KINETIX_NOTARY_COUNCIL_ADMIN_ADDRESS"),
            adminTransferDelay: SafeCast.toUint48(vm.envUint("KINETIX_NOTARY_COUNCIL_ADMIN_TRANSFER_DELAY_SECONDS")),
            ledger: ILiabilityLedger(vm.envAddress("KINETIX_NOTARY_LEDGER_ADDRESS")),
            members: ascending(vm.envAddress("KINETIX_NOTARY_COUNCIL_MEMBERS", ",")),
            threshold: SafeCast.toUint8(vm.envUint("KINETIX_NOTARY_COUNCIL_THRESHOLD")),
            membershipChangeDelay: SafeCast.toUint64(vm.envUint("KINETIX_NOTARY_MEMBERSHIP_CHANGE_DELAY_SECONDS")),
            salt: vm.envBytes32("KINETIX_NOTARY_COUNCIL_DEPLOY_SALT")
        });
    }

    function deploy(CouncilDeployment memory deployment) public returns (ConfidentialAnchorCouncil) {
        bytes memory code = initCode(deployment);

        vm.startBroadcast();
        address council = Create2Factory.deploy(deployment.salt, code);
        vm.stopBroadcast();
        return ConfidentialAnchorCouncil(council);
    }

    function addressOf(CouncilDeployment memory deployment) public pure returns (address) {
        return Create2Factory.addressOf(deployment.salt, initCode(deployment));
    }

    function initCode(CouncilDeployment memory deployment) public pure returns (bytes memory) {
        return abi.encodePacked(
            type(ConfidentialAnchorCouncil).creationCode,
            abi.encode(
                deployment.admin,
                deployment.adminTransferDelay,
                deployment.ledger,
                deployment.members,
                deployment.threshold,
                deployment.membershipChangeDelay
            )
        );
    }

    function ascending(address[] memory members) public pure returns (address[] memory) {
        address[] memory sorted = new address[](members.length);
        for (uint256 i = 0; i < members.length; i++) {
            address member = members[i];
            uint256 j = i;
            while (j > 0 && sorted[j - 1] > member) {
                sorted[j] = sorted[j - 1];
                j--;
            }

            sorted[j] = member;
        }

        return sorted;
    }
}
