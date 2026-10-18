// SPDX-License-Identifier: MIT
pragma solidity 0.8.37;

import {EvidenceFixture} from "../support/EvidenceFixture.sol";
import {EvidenceTreeBuilder} from "../support/EvidenceTreeBuilder.sol";

contract EvidenceTreeFuzzTest is EvidenceFixture {
    function testFuzz_EveryLeafOfAnyBatchProvesAndNoPaddingDoes(
        uint256 seed,
        uint256 rawCount,
        uint256 rawPick
    ) public {
        uint256 count = bound(rawCount, 1, 64);
        bytes32[][] memory levels = _anchorNext(DELIVERY, seed, count);
        bytes32[] memory salts = _salts(seed, count);
        bytes32[] memory records = _records(seed, count);
        uint256 pick = bound(rawPick, 0, count - 1);

        bytes32[] memory siblings = EvidenceTreeBuilder.proof(levels, pick);
        assertTrue(registry.verifyInclusion(DELIVERY, 0, salts[pick], records[pick], uint64(pick), siblings));

        uint256 width = levels[0].length;
        if (count < width) {
            bytes32[] memory paddingProof = EvidenceTreeBuilder.proof(levels, count);
            assertFalse(registry.verifyInclusion(DELIVERY, 0, bytes32(0), bytes32(0), uint64(count), paddingProof));
        }
    }

    function testFuzz_AnyChangedSiblingBreaksTheProof(uint256 seed, uint256 rawCount, uint256 rawLevel) public {
        uint256 count = bound(rawCount, 2, 64);
        bytes32[][] memory levels = _anchorNext(REVIEW, seed, count);
        bytes32[] memory siblings = EvidenceTreeBuilder.proof(levels, 0);
        uint256 level = bound(rawLevel, 0, siblings.length - 1);
        siblings[level] = keccak256(abi.encode(siblings[level]));

        assertFalse(registry.verifyInclusion(REVIEW, 0, _salts(seed, count)[0], _records(seed, count)[0], 0, siblings));
    }
}
