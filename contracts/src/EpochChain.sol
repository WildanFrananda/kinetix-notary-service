// SPDX-License-Identifier: UNLICENSED
pragma solidity 0.8.37;

import {EpochHeader} from "./EpochHeader.sol";
import {SumTree} from "./SumTree.sol";

library EpochChain {
    bytes1 private constant EPOCH_TAG = 0x03;

    function digest(EpochHeader memory header) internal pure returns (bytes32) {
        return keccak256(
            abi.encodePacked(
                EPOCH_TAG,
                SumTree.SCHEME,
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
