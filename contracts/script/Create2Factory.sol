// SPDX-License-Identifier: MIT
pragma solidity 0.8.37;

library Create2Factory {
    address internal constant FACTORY = 0x4e59b44847b379578588920cA78FbF26c0B4956C;

    error FactoryMissing(address factory);
    error DeploymentFailed(address expected);

    function deploy(bytes32 salt, bytes memory initCode) internal returns (address deployed) {
        if (FACTORY.code.length == 0) {
            revert FactoryMissing(FACTORY);
        }

        deployed = addressOf(salt, initCode);
        (bool succeeded, bytes memory returned) = FACTORY.call(abi.encodePacked(salt, initCode));

        if (!succeeded || returned.length != 20 || address(bytes20(returned)) != deployed) {
            revert DeploymentFailed(deployed);
        }

        return deployed;
    }

    function addressOf(bytes32 salt, bytes memory initCode) internal pure returns (address) {
        return address(uint160(uint256(keccak256(abi.encodePacked(bytes1(0xff), FACTORY, salt, keccak256(initCode))))));
    }
}
