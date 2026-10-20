// SPDX-License-Identifier: MIT
pragma solidity 0.8.37;

import {SumNode} from "../../src/SumNode.sol";
import {ConfidentialTree} from "../../src/confidential/ConfidentialTree.sol";

contract ConfidentialTreeHarness {
    function verifies(
        uint64 epochNumber,
        uint256 treeRoot,
        uint64 leafCount,
        uint256 depth,
        bytes32 commitment,
        uint64 amount,
        uint64 index,
        SumNode[] calldata siblings
    ) external pure returns (bool) {
        if (index >= leafCount) {
            return false;
        }
        if (siblings.length != depth) {
            return false;
        }
        if (!ConfidentialTree.wellFormed(siblings)) {
            return false;
        }
        return ConfidentialTree.computeRoot(ConfidentialTree.leaf(epochNumber, commitment, amount), index, siblings)
            == treeRoot;
    }
}
