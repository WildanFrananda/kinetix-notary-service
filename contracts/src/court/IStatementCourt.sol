// SPDX-License-Identifier: MIT
pragma solidity 0.8.37;

import {ILiabilityLedger} from "../ILiabilityLedger.sol";
import {SumNode} from "../SumNode.sol";
import {EpochStatement} from "./EpochStatement.sol";
import {Fault} from "./Fault.sol";
import {FaultKind} from "./FaultKind.sol";
import {PendingSigner} from "./PendingSigner.sol";
import {SignerTerm} from "./SignerTerm.sol";

interface IStatementCourt {
    event SignerScheduled(address indexed signer, uint64 fromEpoch, uint64 activatesAt);
    event SignerCancelled(address indexed signer, uint64 fromEpoch);
    event SignerActivated(address indexed signer, uint64 fromEpoch);
    event FaultProven(uint64 indexed epoch, FaultKind indexed kind, bytes32 indexed evidence, address prover);

    error NotALedger(address ledger);
    error ZeroSigner();
    error ZeroDelay();
    error WrongLedger(address expected, address given);
    error EpochNotAnchored(uint64 epoch);
    error NotSignedBySigner(uint64 epoch, address signer, address recovered);
    error SiblingsMismatch(bytes32 signed, bytes32 given);
    error ProofHolds(uint64 epoch);
    error DifferentEpochs(uint64 first, uint64 second);
    error NoConflict();
    error AlreadyProven(bytes32 evidence);
    error SameSigner(address signer);
    error FromEpochTooLow(uint64 fromEpoch, uint64 least);
    error SignerChangePending(address signer, uint64 fromEpoch);
    error NoSignerChangePending();
    error SignerChangeTooEarly(uint64 activatesAt, uint64 blockTime);
    error SignerChangeStale(uint64 fromEpoch, uint64 nextEpoch);

    function proveBadProof(
        EpochStatement calldata statement,
        bytes calldata signature,
        SumNode[] calldata siblings
    ) external returns (bytes32);

    function proveEquivocation(
        EpochStatement calldata first,
        bytes calldata firstSignature,
        EpochStatement calldata second,
        bytes calldata secondSignature
    ) external returns (bytes32);

    function scheduleSigner(address signer, uint64 fromEpoch) external;

    function cancelSigner() external;

    function activateSigner() external;

    function ledger() external view returns (ILiabilityLedger);

    function signerChangeDelay() external view returns (uint64);

    function signerOf(uint64 epoch) external view returns (address);

    function signerTerms() external view returns (SignerTerm[] memory);

    function pendingSigner() external view returns (PendingSigner memory);

    function statementDigest(EpochStatement calldata statement) external view returns (bytes32);

    function siblingsHash(SumNode[] calldata siblings) external pure returns (bytes32);

    function recoverSigner(EpochStatement calldata statement, bytes calldata signature) external view returns (address);

    function faultCount(uint64 epoch) external view returns (uint256);

    function faultsOf(uint64 epoch) external view returns (Fault[] memory);

    function totalFaults() external view returns (uint256);

    function isProven(bytes32 evidence) external view returns (bool);
}
