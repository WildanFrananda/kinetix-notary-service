// SPDX-License-Identifier: MIT
pragma solidity 0.8.37;

interface IConfidentialLedger {
    event EpochProven(
        uint64 indexed number,
        bytes32 indexed root,
        bytes32 treeRoot,
        bytes32 totalCommitment,
        uint64 leafCount
    );

    error AnchorerNotGranted(address anchorer);
    error NotAVerifier(address verifier);
    error DepthOutOfRange(uint8 depth, uint8 maximum);
    error EpochClosed(uint64 number, uint64 nextEpoch);
    error NotAFieldElement(bytes32 value);
    error StatementAlreadyProven(uint64 number, bytes32 root, uint64 leafCount);
    error ProofRejected(uint64 number);
    error EpochNotProven(uint64 number, bytes32 root, uint64 leafCount);
    error TotalNotHidden(uint64 total);

    function proveEpoch(
        uint64 number,
        bytes32 treeRoot,
        bytes32 totalCommitment,
        uint64 leafCount,
        bytes calldata proof
    ) external returns (bytes32);

    function isProven(uint64 number, bytes32 root, uint64 leafCount) external view returns (bool);

    function treeRootOf(uint64 number) external view returns (bytes32);

    function totalCommitmentOf(uint64 number) external view returns (bytes32);

    function verifier() external view returns (address);

    function depth() external view returns (uint8);

    function capacity() external view returns (uint64);
}
