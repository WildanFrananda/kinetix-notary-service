// SPDX-License-Identifier: MIT
pragma solidity 0.8.37;

import {SafeCast} from "@openzeppelin-contracts/utils/math/SafeCast.sol";
import {IEntryPoint} from "account-abstraction/contracts/interfaces/IEntryPoint.sol";
import {Script} from "forge-std/src/Script.sol";

import {IStatementCourt} from "../src/court/IStatementCourt.sol";
import {NotaryPaymaster} from "../src/paymaster/NotaryPaymaster.sol";
import {IWithholdingCourt} from "../src/withholding/IWithholdingCourt.sol";
import {Create2Factory} from "./Create2Factory.sol";
import {PaymasterDeployment} from "./PaymasterDeployment.sol";

contract DeployNotaryPaymaster is Script {
    function run() external returns (NotaryPaymaster) {
        return deploy(fromEnvironment());
    }

    function fromEnvironment() public view returns (PaymasterDeployment memory) {
        return PaymasterDeployment({
            admin: vm.envAddress("KINETIX_NOTARY_PAYMASTER_ADMIN_ADDRESS"),
            adminTransferDelay: SafeCast.toUint48(vm.envUint("KINETIX_NOTARY_PAYMASTER_ADMIN_TRANSFER_DELAY_SECONDS")),
            entryPoint: IEntryPoint(vm.envAddress("KINETIX_NOTARY_ENTRY_POINT_ADDRESS")),
            account: vm.envAddress("KINETIX_NOTARY_7702_ACCOUNT_ADDRESS"),
            court: IStatementCourt(vm.envAddress("KINETIX_NOTARY_COURT_ADDRESS")),
            withholding: IWithholdingCourt(vm.envAddress("KINETIX_NOTARY_WITHHOLDING_COURT_ADDRESS")),
            maxCostPerOperation: vm.envUint("KINETIX_NOTARY_PAYMASTER_MAX_COST_WEI"),
            unstakeDelay: SafeCast.toUint32(vm.envUint("KINETIX_NOTARY_PAYMASTER_UNSTAKE_DELAY_SECONDS")),
            salt: vm.envBytes32("KINETIX_NOTARY_PAYMASTER_DEPLOY_SALT")
        });
    }

    function deploy(PaymasterDeployment memory deployment) public returns (NotaryPaymaster) {
        bytes memory code = initCode(deployment);

        vm.startBroadcast();
        address paymaster = Create2Factory.deploy(deployment.salt, code);
        vm.stopBroadcast();
        return NotaryPaymaster(paymaster);
    }

    function addressOf(PaymasterDeployment memory deployment) public pure returns (address) {
        return Create2Factory.addressOf(deployment.salt, initCode(deployment));
    }

    function initCode(PaymasterDeployment memory deployment) public pure returns (bytes memory) {
        return abi.encodePacked(
            type(NotaryPaymaster).creationCode,
            abi.encode(
                deployment.admin,
                deployment.adminTransferDelay,
                deployment.entryPoint,
                deployment.account,
                deployment.court,
                deployment.withholding,
                deployment.maxCostPerOperation,
                deployment.unstakeDelay
            )
        );
    }
}
