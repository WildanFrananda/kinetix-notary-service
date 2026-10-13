// SPDX-License-Identifier: MIT
pragma solidity 0.8.37;

import {SumNode} from "../SumNode.sol";
import {EpochStatement} from "./EpochStatement.sol";

library StatementHashing {
    string internal constant NAME = "Kinetix Statement Court";
    string internal constant VERSION = "1";

    bytes32 internal constant STATEMENT_TYPEHASH = keccak256(
        "EpochStatement(address ledger,uint64 epoch,bytes32 commitment,uint64 amount,uint64 index,bytes32 siblingsHash)"
    );

    function structHash(EpochStatement memory statement) internal pure returns (bytes32) {
        return keccak256(
            abi.encode(
                STATEMENT_TYPEHASH,
                statement.ledger,
                statement.epoch,
                statement.commitment,
                statement.amount,
                statement.index,
                statement.siblingsHash
            )
        );
    }

    function siblingsHash(SumNode[] memory siblings) internal pure returns (bytes32) {
        bytes memory packed = new bytes(0);
        for (uint256 i = 0; i < siblings.length; i++) {
            packed = bytes.concat(packed, siblings[i].hash, bytes8(siblings[i].sum));
        }

        return keccak256(packed);
    }

    function conflict(EpochStatement memory a, EpochStatement memory b) internal pure returns (bool) {
        return a.commitment == b.commitment || a.index == b.index;
    }
}
