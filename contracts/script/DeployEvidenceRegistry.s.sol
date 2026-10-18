// SPDX-License-Identifier: MIT
pragma solidity 0.8.37;

import {SafeCast} from "@openzeppelin-contracts/utils/math/SafeCast.sol";
import {Script} from "forge-std/src/Script.sol";

import {EvidenceRegistry} from "../src/evidence/EvidenceRegistry.sol";
import {IEvidenceRegistry} from "../src/evidence/IEvidenceRegistry.sol";
import {StreamGenesis} from "../src/evidence/StreamGenesis.sol";
import {StreamState} from "../src/evidence/StreamState.sol";
import {Create2Factory} from "./Create2Factory.sol";
import {EvidenceDeployment} from "./EvidenceDeployment.sol";

contract DeployEvidenceRegistry is Script {
    uint8 public constant LIFECYCLE = 1;
    uint8 public constant DELIVERY = 2;
    uint8 public constant REVIEW = 3;

    error PredecessorNotSealed(address predecessor);

    function run() external returns (EvidenceRegistry) {
        return deploy(fromEnvironment());
    }

    function fromEnvironment() public view returns (EvidenceDeployment memory) {
        return EvidenceDeployment({
            admin: vm.envAddress("KINETIX_NOTARY_EVIDENCE_ADMIN_ADDRESS"),
            adminTransferDelay: SafeCast.toUint48(vm.envUint("KINETIX_NOTARY_EVIDENCE_ADMIN_TRANSFER_DELAY_SECONDS")),
            anchorer: vm.envAddress("KINETIX_NOTARY_EVIDENCE_ANCHORER_ADDRESS"),
            predecessor: IEvidenceRegistry(vm.envOr("KINETIX_NOTARY_EVIDENCE_PREDECESSOR_ADDRESS", address(0))),
            salt: vm.envBytes32("KINETIX_NOTARY_EVIDENCE_DEPLOY_SALT")
        });
    }

    function deploy(EvidenceDeployment memory deployment) public returns (EvidenceRegistry) {
        bytes memory code = initCode(deployment);

        vm.startBroadcast();
        address registry = Create2Factory.deploy(deployment.salt, code);
        vm.stopBroadcast();
        return EvidenceRegistry(registry);
    }

    function addressOf(EvidenceDeployment memory deployment) public view returns (address) {
        return Create2Factory.addressOf(deployment.salt, initCode(deployment));
    }

    function initCode(EvidenceDeployment memory deployment) public view returns (bytes memory) {
        return abi.encodePacked(
            type(EvidenceRegistry).creationCode,
            abi.encode(deployment.admin, deployment.adminTransferDelay, deployment.anchorer, genesisOf(deployment))
        );
    }

    function genesisOf(EvidenceDeployment memory deployment) public view returns (StreamGenesis[] memory) {
        if (address(deployment.predecessor) == address(0)) {
            StreamGenesis[] memory fresh = new StreamGenesis[](3);
            fresh[0] = _fresh(LIFECYCLE);
            fresh[1] = _fresh(DELIVERY);
            fresh[2] = _fresh(REVIEW);
            return fresh;
        }

        if (!deployment.predecessor.isSealed()) {
            revert PredecessorNotSealed(address(deployment.predecessor));
        }

        uint8[] memory ids = deployment.predecessor.streams();
        StreamGenesis[] memory carried = new StreamGenesis[](ids.length);
        for (uint256 i = 0; i < ids.length; i++) {
            StreamState memory last = deployment.predecessor.stateOf(ids[i]);
            carried[i] = StreamGenesis({
                stream: ids[i],
                batch: last.nextBatch,
                digest: last.head,
                closedAt: last.lastClosedAt,
                throughTransaction: last.throughTransaction,
                throughSequence: last.throughSequence
            });
        }

        return carried;
    }

    function _fresh(uint8 stream) private pure returns (StreamGenesis memory) {
        return StreamGenesis({
            stream: stream,
            batch: 0,
            digest: bytes32(0),
            closedAt: 0,
            throughTransaction: 0,
            throughSequence: 0
        });
    }
}
