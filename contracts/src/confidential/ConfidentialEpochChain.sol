// SPDX-License-Identifier: MIT
pragma solidity 0.8.37;

import {EpochHeader} from "../EpochHeader.sol";
import {ConfidentialTree} from "./ConfidentialTree.sol";

library ConfidentialEpochChain {
    bytes1 private constant EPOCH_TAG = 0x03;

    function digest(EpochHeader memory header) internal pure returns (bytes32) {
        return keccak256(
            abi.encodePacked(
                EPOCH_TAG,
                ConfidentialTree.SCHEME,
                header.number,
                header.previousDigest,
                header.root,
                header.total,
                header.leafCount,
                header.asOf
            )
        );
    }
}
