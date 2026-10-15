// SPDX-License-Identifier: MIT
pragma solidity 0.8.37;

import {SafeCast} from "@openzeppelin-contracts/utils/math/SafeCast.sol";
import {Script} from "forge-std/src/Script.sol";

import {ILiabilityLedger} from "../src/ILiabilityLedger.sol";
import {IStatementCourt} from "../src/court/IStatementCourt.sol";
import {LedgerDirectory} from "../src/directory/LedgerDirectory.sol";
import {IWithholdingCourt} from "../src/withholding/IWithholdingCourt.sol";
import {Create2Factory} from "./Create2Factory.sol";
import {DirectoryDeployment} from "./DirectoryDeployment.sol";

contract DeployLedgerDirectory is Script {
    function run() external returns (LedgerDirectory) {
        return deploy(fromEnvironment());
    }

    function fromEnvironment() public view returns (DirectoryDeployment memory) {
        return DirectoryDeployment({
            admin: vm.envAddress("KINETIX_NOTARY_DIRECTORY_ADMIN_ADDRESS"),
            adminTransferDelay: SafeCast.toUint48(vm.envUint("KINETIX_NOTARY_DIRECTORY_ADMIN_TRANSFER_DELAY_SECONDS")),
            ledger: ILiabilityLedger(vm.envAddress("KINETIX_NOTARY_LEDGER_ADDRESS")),
            court: IStatementCourt(vm.envAddress("KINETIX_NOTARY_COURT_ADDRESS")),
            withholding: IWithholdingCourt(vm.envAddress("KINETIX_NOTARY_WITHHOLDING_COURT_ADDRESS")),
            foreignSuccessorDelay: SafeCast.toUint64(vm.envUint("KINETIX_NOTARY_FOREIGN_SUCCESSOR_DELAY_SECONDS")),
            salt: vm.envBytes32("KINETIX_NOTARY_DIRECTORY_DEPLOY_SALT")
        });
    }

    function deploy(DirectoryDeployment memory deployment) public returns (LedgerDirectory) {
        bytes memory code = initCode(deployment);

        vm.startBroadcast();
        address directory = Create2Factory.deploy(deployment.salt, code);
        vm.stopBroadcast();
        return LedgerDirectory(directory);
    }

    function addressOf(DirectoryDeployment memory deployment) public pure returns (address) {
        return Create2Factory.addressOf(deployment.salt, initCode(deployment));
    }

    function initCode(DirectoryDeployment memory deployment) public pure returns (bytes memory) {
        return abi.encodePacked(
            type(LedgerDirectory).creationCode,
            abi.encode(
                deployment.admin,
                deployment.adminTransferDelay,
                deployment.ledger,
                deployment.court,
                deployment.withholding,
                deployment.foreignSuccessorDelay
            )
        );
    }
}
