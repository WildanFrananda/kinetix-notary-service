// SPDX-License-Identifier: MIT
pragma solidity 0.8.37;

import {Registration} from "./Registration.sol";

library RegistrationHashing {
    string internal constant NAME = "Kinetix Withholding Court";
    string internal constant VERSION = "1";

    bytes32 internal constant REGISTRATION_TYPEHASH =
        keccak256("Registration(address holder,bytes32 account,uint64 fromEpoch)");

    function structHash(Registration memory registration) internal pure returns (bytes32) {
        return
            keccak256(
                abi.encode(REGISTRATION_TYPEHASH, registration.holder, registration.account, registration.fromEpoch)
            );
    }
}
