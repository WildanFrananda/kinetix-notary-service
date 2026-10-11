// SPDX-License-Identifier: UNLICENSED
pragma solidity 0.8.37;

import {SafeCast} from "@openzeppelin-contracts/utils/math/SafeCast.sol";
import {Script} from "forge-std/src/Script.sol";

import {LiabilityLedger} from "../src/LiabilityLedger.sol";
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
        vm.startBroadcast();
        LiabilityLedger ledger = new LiabilityLedger{salt: deployment.salt}(
            deployment.admin,
            deployment.adminTransferDelay,
            deployment.anchorer,
            deployment.genesisEpoch,
            deployment.genesisDigest,
            deployment.genesisAsOf
        );
        vm.stopBroadcast();
        return ledger;
    }
}
