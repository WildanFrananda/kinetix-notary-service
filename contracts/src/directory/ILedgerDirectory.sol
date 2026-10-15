// SPDX-License-Identifier: MIT
pragma solidity 0.8.37;

import {ILiabilityLedger} from "../ILiabilityLedger.sol";
import {IStatementCourt} from "../court/IStatementCourt.sol";
import {IWithholdingCourt} from "../withholding/IWithholdingCourt.sol";
import {Generation} from "./Generation.sol";
import {PendingGeneration} from "./PendingGeneration.sol";

interface ILedgerDirectory {
    event GenerationRegistered(
        uint256 indexed index,
        uint64 indexed chainId,
        address indexed ledger,
        address court,
        address withholding,
        uint64 genesisEpoch,
        bytes32 genesisDigest,
        uint16 scheme
    );
    event ForeignSuccessorProposed(
        uint64 indexed chainId,
        address indexed ledger,
        address court,
        address withholding,
        uint16 scheme,
        uint64 activatesAt
    );
    event ForeignSuccessorCancelled(uint64 indexed chainId, address indexed ledger);

    error ZeroDelay();
    error NotALedger(address ledger);
    error NotAStatementCourt(address court);
    error NotAWithholdingCourt(address withholding);
    error CourtOfAnotherLedger(address court, address ledger);
    error WithholdingOfAnotherCourt(address withholding, address court);
    error SuccessorPending(uint64 activatesAt);
    error DirectoryMoved(uint64 chainId);
    error PredecessorNotSealed(address ledger);
    error SuccessorSealed(address ledger);
    error GenesisMismatch(uint64 expectedEpoch, bytes32 expectedDigest, uint64 givenEpoch, bytes32 givenDigest);
    error NotForeign(uint64 chainId);
    error ZeroAddress();
    error ZeroScheme();
    error NoSuccessorPending();
    error SuccessorTooEarly(uint64 activatesAt, uint64 blockTime);
    error UnknownGeneration(uint256 index, uint256 count);
    error EpochBeforeDirectory(uint64 epoch, uint64 genesisEpoch);

    function registerSuccessor(
        ILiabilityLedger ledger,
        IStatementCourt court,
        IWithholdingCourt withholding
    ) external returns (uint256);

    function proposeForeignSuccessor(
        uint64 chainId,
        address ledger,
        address court,
        address withholding,
        uint16 scheme
    ) external;

    function cancelForeignSuccessor() external;

    function activateForeignSuccessor() external returns (uint256);

    function foreignSuccessorDelay() external view returns (uint64);

    function generationCount() external view returns (uint256);

    function generation(uint256 index) external view returns (Generation memory);

    function current() external view returns (Generation memory);

    function generationOf(uint64 epoch) external view returns (uint256, Generation memory);

    function pendingSuccessor() external view returns (PendingGeneration memory);
}
