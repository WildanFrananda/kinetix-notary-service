// SPDX-License-Identifier: MIT
pragma solidity 0.8.37;

import {SafeCast} from "@openzeppelin-contracts/utils/math/SafeCast.sol";
import {Script} from "forge-std/src/Script.sol";

import {ILiabilityLedger} from "../src/ILiabilityLedger.sol";
import {LiabilityLedgerV2} from "../src/confidential/LiabilityLedgerV2.sol";
import {PlonkVerifier} from "../src/confidential/PlonkVerifier.sol";
import {Create2Factory} from "./Create2Factory.sol";
import {LedgerV2Deployment} from "./LedgerV2Deployment.sol";

contract DeployLiabilityLedgerV2 is Script {
    error PredecessorNotSealed(address predecessor);

    function run() external returns (LiabilityLedgerV2) {
        return deploy(fromEnvironment());
    }

    function fromEnvironment() public view returns (LedgerV2Deployment memory) {
        return LedgerV2Deployment({
            admin: vm.envAddress("KINETIX_NOTARY_ADMIN_ADDRESS"),
            adminTransferDelay: SafeCast.toUint48(vm.envUint("KINETIX_NOTARY_ADMIN_TRANSFER_DELAY_SECONDS")),
            anchorer: vm.envAddress("KINETIX_NOTARY_ANCHORER_ADDRESS"),
            predecessor: ILiabilityLedger(vm.envAddress("KINETIX_NOTARY_PREDECESSOR_LEDGER_ADDRESS")),
            depth: SafeCast.toUint8(vm.envUint("KINETIX_NOTARY_LEDGER_DEPTH")),
            verifierSalt: vm.envBytes32("KINETIX_NOTARY_PLONK_VERIFIER_DEPLOY_SALT"),
            salt: vm.envBytes32("KINETIX_NOTARY_DEPLOY_SALT")
        });
    }

    function deploy(LedgerV2Deployment memory deployment) public returns (LiabilityLedgerV2) {
        if (!deployment.predecessor.isSealed()) {
            revert PredecessorNotSealed(address(deployment.predecessor));
        }

        bytes memory ledgerCode = initCode(deployment);

        vm.startBroadcast();
        Create2Factory.deploy(deployment.verifierSalt, type(PlonkVerifier).creationCode);
        address ledger = Create2Factory.deploy(deployment.salt, ledgerCode);
        vm.stopBroadcast();
        return LiabilityLedgerV2(ledger);
    }

    function verifierAddressOf(LedgerV2Deployment memory deployment) public pure returns (address) {
        return Create2Factory.addressOf(deployment.verifierSalt, type(PlonkVerifier).creationCode);
    }

    function addressOf(LedgerV2Deployment memory deployment) public view returns (address) {
        return Create2Factory.addressOf(deployment.salt, initCode(deployment));
    }

    function initCode(LedgerV2Deployment memory deployment) public view returns (bytes memory) {
        return abi.encodePacked(
            type(LiabilityLedgerV2).creationCode,
            abi.encode(
                deployment.admin,
                deployment.adminTransferDelay,
                deployment.anchorer,
                deployment.predecessor.nextEpoch(),
                deployment.predecessor.head(),
                deployment.predecessor.lastAsOf(),
                verifierAddressOf(deployment),
                deployment.depth
            )
        );
    }
}
