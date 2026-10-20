// SPDX-License-Identifier: MIT
pragma solidity 0.8.37;

import {Test} from "forge-std/src/Test.sol";

import {CouncilDeployment} from "../script/CouncilDeployment.sol";
import {DeployConfidentialAnchorCouncil} from "../script/DeployConfidentialAnchorCouncil.s.sol";
import {DeployLiabilityLedgerV2} from "../script/DeployLiabilityLedgerV2.s.sol";
import {LedgerV2Deployment} from "../script/LedgerV2Deployment.sol";
import {EpochHeader} from "../src/EpochHeader.sol";
import {LiabilityLedger} from "../src/LiabilityLedger.sol";
import {SumNode} from "../src/SumNode.sol";
import {ConfidentialAnchorCouncil} from "../src/confidential/ConfidentialAnchorCouncil.sol";
import {ConfidentialEpochChain} from "../src/confidential/ConfidentialEpochChain.sol";
import {LiabilityLedgerV2} from "../src/confidential/LiabilityLedgerV2.sol";
import {StatementCourt} from "../src/court/StatementCourt.sol";
import {Generation} from "../src/directory/Generation.sol";
import {LedgerDirectory} from "../src/directory/LedgerDirectory.sol";
import {WithholdingCourt} from "../src/withholding/WithholdingCourt.sol";

contract GenerationTwoForkTest is Test {
    string private fixture;
    LiabilityLedger private ledger;
    LedgerDirectory private directory;
    address private admin;
    address private notary;
    uint256 private notaryKey;
    address private watch;
    uint256 private watchKey;

    function setUp() public {
        vm.createSelectFork(
            vm.envString("KINETIX_NOTARY_FORK_RPC_URL"), vm.envUint("KINETIX_NOTARY_FORK_GENERATION_BLOCK")
        );
        fixture = vm.readFile(string.concat(vm.projectRoot(), "/test/fixtures/plonk-production.json"));
        ledger = LiabilityLedger(vm.envAddress("KINETIX_NOTARY_FORK_LEDGER"));
        directory = LedgerDirectory(vm.envAddress("KINETIX_NOTARY_FORK_DIRECTORY"));
        admin = ledger.defaultAdmin();
        (notary, notaryKey) = makeAddrAndKey("fork-generation-notary");
        (watch, watchKey) = makeAddrAndKey("fork-generation-watch");
    }

    function test_TheProofFixtureIsForTheEpochAfterTheLiveHead() public view {
        assertEq(vm.parseJsonUint(fixture, ".epoch"), ledger.nextEpoch());
        assertEq(vm.parseJsonUint(fixture, ".depth"), 8);
        assertFalse(ledger.isSealed());
        assertEq(directory.generationCount(), 1);
        assertEq(directory.current().ledger, address(ledger));
    }

    function test_GenerationTwoSucceedsTheLiveLedgerThroughTheLiveDirectory() public {
        bytes32 head = ledger.head();
        uint64 number = ledger.nextEpoch();
        uint64 lastAsOf = ledger.lastAsOf();
        vm.prank(admin);
        ledger.seal(head);

        LiabilityLedgerV2 next = new DeployLiabilityLedgerV2()
            .deploy(
                LedgerV2Deployment({
                    admin: admin,
                    adminTransferDelay: 3 days,
                    anchorer: notary,
                    predecessor: ledger,
                    depth: 8,
                    verifierSalt: keccak256("kinetix.notary.fork.verifier"),
                    salt: keccak256("kinetix.notary.fork.ledger-v2")
                })
            );
        StatementCourt nextCourt = new StatementCourt(admin, 3 days, next, makeAddr("fork-generation-signer"), 7 days);
        WithholdingCourt nextWithholding = new WithholdingCourt(nextCourt, 7 days, 24 hours);
        ConfidentialAnchorCouncil council = _council(next);

        bytes32 role = next.ANCHORER_ROLE();
        vm.startPrank(admin);
        next.grantRole(role, address(council));
        next.revokeRole(role, notary);
        uint256 index = directory.registerSuccessor(next, nextCourt, nextWithholding);
        vm.stopPrank();

        Generation memory current = directory.current();
        assertEq(index, 1);
        assertEq(current.ledger, address(next));
        assertEq(current.scheme, 2);
        assertEq(current.genesisEpoch, number);
        assertEq(current.genesisDigest, head);

        bytes32 root = next.proveEpoch(
            number,
            vm.parseJsonBytes32(fixture, ".treeRoot"),
            vm.parseJsonBytes32(fixture, ".totalCommitment"),
            uint64(vm.parseJsonUint(fixture, ".leafCount")),
            vm.parseJsonBytes(fixture, ".proof")
        );
        EpochHeader memory header = EpochHeader({
            number: number,
            previousDigest: head,
            root: root,
            total: 0,
            leafCount: uint64(vm.parseJsonUint(fixture, ".leafCount")),
            asOf: lastAsOf + 1
        });
        bytes32 digest = council.anchor(header, _approvals(council, header));
        assertEq(digest, ConfidentialEpochChain.digest(header));
        assertEq(next.head(), digest);

        (uint256 generation, Generation memory found) = directory.generationOf(number);
        assertEq(generation, 1);
        _everyAccountIsIncluded(LiabilityLedgerV2(found.ledger), number);
    }

    function _council(LiabilityLedgerV2 next) private returns (ConfidentialAnchorCouncil) {
        DeployConfidentialAnchorCouncil script = new DeployConfidentialAnchorCouncil();
        address[] memory members = new address[](2);
        members[0] = notary;
        members[1] = watch;
        return script.deploy(
            CouncilDeployment({
                admin: admin,
                adminTransferDelay: 3 days,
                ledger: next,
                members: script.ascending(members),
                threshold: 2,
                membershipChangeDelay: 7 days,
                salt: keccak256("kinetix.notary.fork.council-v2")
            })
        );
    }

    function _approvals(
        ConfidentialAnchorCouncil council,
        EpochHeader memory header
    ) private view returns (bytes[] memory) {
        bytes32 hash = council.approvalDigest(header.number, council.epochDigest(header));
        (uint8 v1, bytes32 r1, bytes32 s1) = vm.sign(notaryKey, hash);
        (uint8 v2, bytes32 r2, bytes32 s2) = vm.sign(watchKey, hash);
        bytes[] memory approvals = new bytes[](2);
        (approvals[0], approvals[1]) = notary < watch
            ? (abi.encodePacked(r1, s1, v1), abi.encodePacked(r2, s2, v2))
            : (abi.encodePacked(r2, s2, v2), abi.encodePacked(r1, s1, v1));
        return approvals;
    }

    function _everyAccountIsIncluded(LiabilityLedgerV2 next, uint64 number) private view {
        bytes32[] memory commitments = vm.parseJsonBytes32Array(fixture, ".commitments");
        uint256[] memory amounts = vm.parseJsonUintArray(fixture, ".amounts");
        assertEq(commitments.length, 5);
        for (uint256 i = 0; i < commitments.length; i++) {
            string memory key = string.concat(".inclusions[", vm.toString(i), "]");
            bytes32[] memory hashes = vm.parseJsonBytes32Array(fixture, string.concat(key, ".siblingHashes"));
            SumNode[] memory siblings = new SumNode[](hashes.length);
            for (uint256 level = 0; level < hashes.length; level++) {
                siblings[level] = SumNode({hash: hashes[level], sum: 0});
            }

            uint64 index = uint64(vm.parseJsonUint(fixture, string.concat(key, ".index")));
            assertTrue(next.verifyInclusion(number, commitments[i], uint64(amounts[i]), index, siblings), key);
        }
    }
}
