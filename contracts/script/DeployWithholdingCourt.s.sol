// SPDX-License-Identifier: MIT
pragma solidity 0.8.37;

import {SafeCast} from "@openzeppelin-contracts/utils/math/SafeCast.sol";
import {Script} from "forge-std/src/Script.sol";

import {IStatementCourt} from "../src/court/IStatementCourt.sol";
import {WithholdingCourt} from "../src/withholding/WithholdingCourt.sol";
import {Create2Factory} from "./Create2Factory.sol";
import {WithholdingDeployment} from "./WithholdingDeployment.sol";

contract DeployWithholdingCourt is Script {
    function run() external returns (WithholdingCourt) {
        return deploy(fromEnvironment());
    }

    function fromEnvironment() public view returns (WithholdingDeployment memory) {
        return WithholdingDeployment({
            court: IStatementCourt(vm.envAddress("KINETIX_NOTARY_COURT_ADDRESS")),
            requestWindow: SafeCast.toUint64(vm.envUint("KINETIX_NOTARY_WITHHOLDING_REQUEST_WINDOW_SECONDS")),
            answerPeriod: SafeCast.toUint64(vm.envUint("KINETIX_NOTARY_WITHHOLDING_ANSWER_PERIOD_SECONDS")),
            salt: vm.envBytes32("KINETIX_NOTARY_WITHHOLDING_DEPLOY_SALT")
        });
    }

    function deploy(WithholdingDeployment memory deployment) public returns (WithholdingCourt) {
        bytes memory code = initCode(deployment);

        vm.startBroadcast();
        address withholding = Create2Factory.deploy(deployment.salt, code);
        vm.stopBroadcast();
        return WithholdingCourt(withholding);
    }

    function addressOf(WithholdingDeployment memory deployment) public pure returns (address) {
        return Create2Factory.addressOf(deployment.salt, initCode(deployment));
    }

    function initCode(WithholdingDeployment memory deployment) public pure returns (bytes memory) {
        return abi.encodePacked(
            type(WithholdingCourt).creationCode,
            abi.encode(deployment.court, deployment.requestWindow, deployment.answerPeriod)
        );
    }
}
