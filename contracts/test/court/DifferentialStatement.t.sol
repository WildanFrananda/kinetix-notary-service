// SPDX-License-Identifier: MIT
pragma solidity 0.8.37;

import {SumNode} from "../../src/SumNode.sol";
import {EpochStatement} from "../../src/court/EpochStatement.sol";
import {StatementCourt} from "../../src/court/StatementCourt.sol";
import {LedgerFixture} from "../support/LedgerFixture.sol";

contract DifferentialStatementTest is LedgerFixture {
    string private fixture;
    StatementCourt private court;

    function setUp() public override {
        super.setUp();
        fixture = vm.readFile(string.concat(vm.projectRoot(), "/test/fixtures/statement-v1.json"));
        vm.chainId(vm.parseJsonUint(fixture, ".chainId"));
        address where = vm.parseJsonAddress(fixture, ".court");
        deployCodeTo(
            "StatementCourt.sol:StatementCourt",
            abi.encode(admin, ADMIN_DELAY, ledger, vm.parseJsonAddress(fixture, ".signer"), uint64(7 days)),
            where
        );
        court = StatementCourt(where);
    }

    function test_TheDomainAgrees() public view {
        (, string memory name, string memory version, uint256 chainId, address verifyingContract,,) =
            court.eip712Domain();

        assertEq(name, vm.parseJsonString(fixture, ".name"));
        assertEq(version, vm.parseJsonString(fixture, ".version"));
        assertEq(chainId, vm.parseJsonUint(fixture, ".chainId"));
        assertEq(verifyingContract, vm.parseJsonAddress(fixture, ".court"));
    }

    function test_EveryStatementGoSignsHasTheSameDigestAndSignerHere() public view {
        uint256 count = vm.parseJsonUint(fixture, ".statementCount");
        address signer = vm.parseJsonAddress(fixture, ".signer");
        assertGt(count, 0);

        for (uint256 i = 0; i < count; i++) {
            string memory key = string.concat(".statements[", vm.toString(i), "]");
            SumNode[] memory siblings = _siblings(key);
            EpochStatement memory statement = EpochStatement({
                ledger: vm.parseJsonAddress(fixture, string.concat(key, ".ledger")),
                epoch: uint64(vm.parseJsonUint(fixture, string.concat(key, ".epoch"))),
                commitment: vm.parseJsonBytes32(fixture, string.concat(key, ".commitment")),
                amount: uint64(vm.parseJsonUint(fixture, string.concat(key, ".amount"))),
                index: uint64(vm.parseJsonUint(fixture, string.concat(key, ".index"))),
                siblingsHash: vm.parseJsonBytes32(fixture, string.concat(key, ".siblingsHash"))
            });
            bytes memory signature = vm.parseJsonBytes(fixture, string.concat(key, ".signature"));

            assertEq(court.siblingsHash(siblings), statement.siblingsHash, key);
            assertEq(court.statementDigest(statement), vm.parseJsonBytes32(fixture, string.concat(key, ".digest")), key);
            assertEq(court.recoverSigner(statement, signature), signer, key);
        }
    }

    function _siblings(string memory key) private view returns (SumNode[] memory) {
        string memory hashesKey = string.concat(key, ".siblingHashes");
        if (!vm.keyExistsJson(fixture, string.concat(hashesKey, "[0]"))) {
            return new SumNode[](0);
        }

        bytes32[] memory hashes = vm.parseJsonBytes32Array(fixture, hashesKey);
        uint256[] memory sums = vm.parseJsonUintArray(fixture, string.concat(key, ".siblingSums"));
        SumNode[] memory siblings = new SumNode[](hashes.length);
        for (uint256 i = 0; i < hashes.length; i++) {
            siblings[i] = SumNode({hash: hashes[i], sum: uint64(sums[i])});
        }

        return siblings;
    }
}
