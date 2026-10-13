// SPDX-License-Identifier: UNLICENSED
pragma solidity 0.8.37;

import {SafeCast} from "@openzeppelin-contracts/utils/math/SafeCast.sol";
import {Script} from "forge-std/src/Script.sol";

import {LiabilityLedger} from "../src/LiabilityLedger.sol";
import {Create2Factory} from "./Create2Factory.sol";
import {LedgerDeployment} from "./LedgerDeployment.sol";

contract DeployLiabilityLedger is Script {
    function run() external returns (LiabilityLedger) {
        return deploy(fromEnvironment());
    }

    function fromEnvironment() public view returns (LedgerDeployment memory) {
        return LedgerDeployment({
            admin: vm.envAddress("KINETIX_NOTARY_ADMIN_ADDRESS"),
            adminTransferDelay: SafeCast.toUint48(vm.envUint("KINETIX_NOTARY_ADMIN_TRANSFER_DELAY_SECONDS")),
            anchorer: vm.envAddress("KINETIX_NOTARY_ANCHORER_ADDRESS"),
            genesisEpoch: SafeCast.toUint64(vm.envUint("KINETIX_NOTARY_GENESIS_EPOCH")),
            genesisDigest: vm.envBytes32("KINETIX_NOTARY_GENESIS_DIGEST"),
            genesisAsOf: SafeCast.toUint64(vm.envUint("KINETIX_NOTARY_GENESIS_AS_OF")),
            salt: vm.envBytes32("KINETIX_NOTARY_DEPLOY_SALT")
        });
    }

    function deploy(LedgerDeployment memory deployment) public returns (LiabilityLedger) {
        bytes memory code = initCode(deployment);

        vm.startBroadcast();
        address ledger = Create2Factory.deploy(deployment.salt, code);
        vm.stopBroadcast();
        return LiabilityLedger(ledger);
    }

    function addressOf(LedgerDeployment memory deployment) public pure returns (address) {
        return Create2Factory.addressOf(deployment.salt, initCode(deployment));
    }

    function initCode(LedgerDeployment memory deployment) public pure returns (bytes memory) {
        return abi.encodePacked(
            type(LiabilityLedger).creationCode,
            abi.encode(
                deployment.admin,
                deployment.adminTransferDelay,
                deployment.anchorer,
                deployment.genesisEpoch,
                deployment.genesisDigest,
                deployment.genesisAsOf
            )
        );
    }
}
