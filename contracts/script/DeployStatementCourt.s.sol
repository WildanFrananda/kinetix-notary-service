// SPDX-License-Identifier: MIT
pragma solidity 0.8.37;

import {SafeCast} from "@openzeppelin-contracts/utils/math/SafeCast.sol";
import {Script} from "forge-std/src/Script.sol";

import {ILiabilityLedger} from "../src/ILiabilityLedger.sol";
import {StatementCourt} from "../src/court/StatementCourt.sol";
import {CourtDeployment} from "./CourtDeployment.sol";
import {Create2Factory} from "./Create2Factory.sol";

contract DeployStatementCourt is Script {
    function run() external returns (StatementCourt) {
        return deploy(fromEnvironment());
    }

    function fromEnvironment() public view returns (CourtDeployment memory) {
        return CourtDeployment({
            admin: vm.envAddress("KINETIX_NOTARY_COURT_ADMIN_ADDRESS"),
            adminTransferDelay: SafeCast.toUint48(vm.envUint("KINETIX_NOTARY_COURT_ADMIN_TRANSFER_DELAY_SECONDS")),
            ledger: ILiabilityLedger(vm.envAddress("KINETIX_NOTARY_LEDGER_ADDRESS")),
            signer: vm.envAddress("KINETIX_NOTARY_STATEMENT_SIGNER_ADDRESS"),
            signerChangeDelay: SafeCast.toUint64(vm.envUint("KINETIX_NOTARY_SIGNER_CHANGE_DELAY_SECONDS")),
            salt: vm.envBytes32("KINETIX_NOTARY_COURT_DEPLOY_SALT")
        });
    }

    function deploy(CourtDeployment memory deployment) public returns (StatementCourt) {
        bytes memory code = initCode(deployment);

        vm.startBroadcast();
        address court = Create2Factory.deploy(deployment.salt, code);
        vm.stopBroadcast();
        return StatementCourt(court);
    }

    function addressOf(CourtDeployment memory deployment) public pure returns (address) {
        return Create2Factory.addressOf(deployment.salt, initCode(deployment));
    }

    function initCode(CourtDeployment memory deployment) public pure returns (bytes memory) {
        return abi.encodePacked(
            type(StatementCourt).creationCode,
            abi.encode(
                deployment.admin,
                deployment.adminTransferDelay,
                deployment.ledger,
                deployment.signer,
                deployment.signerChangeDelay
            )
        );
    }
}
