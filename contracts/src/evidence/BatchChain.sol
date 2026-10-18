// SPDX-License-Identifier: MIT
pragma solidity 0.8.37;

import {BatchHeader} from "./BatchHeader.sol";
import {EvidenceTree} from "./EvidenceTree.sol";

library BatchChain {
    bytes1 private constant BATCH_TAG = 0x13;

    function digest(BatchHeader memory header) internal pure returns (bytes32) {
        return keccak256(
            abi.encodePacked(
                BATCH_TAG,
                EvidenceTree.SCHEME,
                header.stream,
                header.number,
                header.previousDigest,
                header.root,
                header.leafCount,
                header.throughTransaction,
                header.throughSequence,
                header.closedAt
            )
        );
    }
}
