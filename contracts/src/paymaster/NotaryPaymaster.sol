// SPDX-License-Identifier: MIT
pragma solidity 0.8.37;

import {
    AccessControlDefaultAdminRules
} from "@openzeppelin-contracts/access/extensions/AccessControlDefaultAdminRules.sol";
import {ECDSA} from "@openzeppelin-contracts/utils/cryptography/ECDSA.sol";
import {MessageHashUtils} from "@openzeppelin-contracts/utils/cryptography/MessageHashUtils.sol";
import {ERC165Checker} from "@openzeppelin-contracts/utils/introspection/ERC165Checker.sol";
import {SafeCast} from "@openzeppelin-contracts/utils/math/SafeCast.sol";
import {BaseAccount} from "account-abstraction/contracts/core/BaseAccount.sol";
import {_packValidationData} from "account-abstraction/contracts/core/Helpers.sol";
import {UserOperationLib} from "account-abstraction/contracts/core/UserOperationLib.sol";
import {IEntryPoint} from "account-abstraction/contracts/interfaces/IEntryPoint.sol";
import {IStakeManager} from "account-abstraction/contracts/interfaces/IStakeManager.sol";
import {PackedUserOperation} from "account-abstraction/contracts/interfaces/PackedUserOperation.sol";

import {ILiabilityLedger} from "../ILiabilityLedger.sol";
import {SumNode} from "../SumNode.sol";
import {EpochStatement} from "../court/EpochStatement.sol";
import {IStatementCourt} from "../court/IStatementCourt.sol";
import {StatementHashing} from "../court/StatementHashing.sol";
import {IWithholdingCourt} from "../withholding/IWithholdingCourt.sol";
import {OpenRequest} from "../withholding/OpenRequest.sol";
import {Registration} from "../withholding/Registration.sol";
import {RegistrationHashing} from "../withholding/RegistrationHashing.sol";
import {RequestState} from "../withholding/RequestState.sol";
import {StatementRequest} from "../withholding/StatementRequest.sol";
import {INotaryPaymaster} from "./INotaryPaymaster.sol";

contract NotaryPaymaster is INotaryPaymaster, AccessControlDefaultAdminRules {
    bytes32 private constant DOMAIN_TYPEHASH =
        keccak256("EIP712Domain(string name,string version,uint256 chainId,address verifyingContract)");
    bytes32 private constant SAMPLE = keccak256("Kinetix Notary Paymaster domain check");
    uint256 private constant SELECTOR_SIZE = 4;
    uint256 private constant WORD_SIZE = 32;
    uint32 private constant LEAST_UNSTAKE_DELAY = 1 days;

    IEntryPoint public immutable entryPoint;
    address public immutable account;
    IStatementCourt public immutable court;
    IWithholdingCourt public immutable withholding;
    address public immutable statementSigner;
    uint256 public immutable maxCostPerOperation;
    uint32 public immutable unstakeDelay;

    bytes32 private immutable _delegation;
    bytes32 private immutable _courtDomain;
    bytes32 private immutable _withholdingDomain;
    ILiabilityLedger private immutable _ledger;
    uint64 private immutable _requestWindow;

    constructor(
        address admin,
        uint48 adminTransferDelay,
        IEntryPoint entryPoint_,
        address account_,
        IStatementCourt court_,
        IWithholdingCourt withholding_,
        uint256 maxCostPerOperation_,
        uint32 unstakeDelay_
    ) AccessControlDefaultAdminRules(adminTransferDelay, admin) {
        if (account_ == address(0)) {
            revert NoContract(account_);
        }

        _requireCode(address(entryPoint_));
        _requireCode(account_);
        if (!ERC165Checker.supportsInterface(address(court_), type(IStatementCourt).interfaceId)) {
            revert NotAStatementCourt(address(court_));
        }

        if (!ERC165Checker.supportsInterface(address(withholding_), type(IWithholdingCourt).interfaceId)) {
            revert NotAWithholdingCourt(address(withholding_));
        }

        if (address(_courtOf(withholding_)) != address(court_)) {
            revert WithholdingOfAnotherCourt(address(withholding_), address(court_));
        }

        if (maxCostPerOperation_ == 0) {
            revert ZeroCostCeiling();
        }

        if (unstakeDelay_ < LEAST_UNSTAKE_DELAY) {
            revert UnstakeDelayTooShort(unstakeDelay_, LEAST_UNSTAKE_DELAY);
        }

        address signer = _currentSignerOf(court_);
        if (signer == address(0)) {
            revert NoStatementSigner(address(court_));
        }

        bytes32 courtDomain = _domain(StatementHashing.NAME, StatementHashing.VERSION, address(court_));
        bytes32 withholdingDomain =
            _domain(RegistrationHashing.NAME, RegistrationHashing.VERSION, address(withholding_));
        _requireSameDomains(court_, withholding_, courtDomain, withholdingDomain);

        entryPoint = entryPoint_;
        account = account_;
        court = court_;
        withholding = withholding_;
        statementSigner = signer;
        maxCostPerOperation = maxCostPerOperation_;
        unstakeDelay = unstakeDelay_;
        _delegation = keccak256(abi.encodePacked(bytes3(0xef0100), account_));
        _courtDomain = courtDomain;
        _withholdingDomain = withholdingDomain;
        _ledger = _ledgerOf(court_);
        _requestWindow = _requestWindowOf(withholding_);
    }

    function validatePaymasterUserOp(
        PackedUserOperation calldata userOp,
        bytes32,
        uint256 maxCost
    ) external view returns (bytes memory context, uint256 validationData) {
        if (msg.sender != address(entryPoint)) {
            revert NotTheEntryPoint(msg.sender);
        }

        if (maxCost > maxCostPerOperation) {
            revert CostTooHigh(maxCost, maxCostPerOperation);
        }

        if (userOp.sender.codehash != _delegation) {
            revert NotDelegated(userOp.sender);
        }

        (uint48 validAfter, uint48 validUntil) = _sponsorship(
            userOp.sender, userOp.callData, userOp.paymasterAndData[UserOperationLib.PAYMASTER_DATA_OFFSET:]
        );
        return ("", _packValidationData(false, validUntil, validAfter));
    }

    function postOp(PostOpMode, bytes calldata, uint256, uint256) external pure {
        revert NoPostOp();
    }

    function sponsorship(
        address sender,
        bytes calldata callData,
        bytes calldata paymasterData
    ) external view returns (uint48 validAfter, uint48 validUntil) {
        return _sponsorship(sender, callData, paymasterData);
    }

    function withdrawTo(address payable to, uint256 amount) external onlyRole(DEFAULT_ADMIN_ROLE) {
        emit Withdrawn(to, amount);
        entryPoint.withdrawTo(to, amount);
    }

    function addStake() external payable onlyRole(DEFAULT_ADMIN_ROLE) {
        entryPoint.addStake{value: msg.value}(unstakeDelay);
    }

    function withdrawBalance() external {
        address admin = defaultAdmin();
        if (msg.sender != admin) {
            revert AccessControlUnauthorizedAccount(msg.sender, DEFAULT_ADMIN_ROLE);
        }

        uint256 balance = address(this).balance;
        if (!payable(admin).send(balance)) {
            revert BalanceNotSent(admin, balance);
        }
    }

    function unlockStake() external onlyRole(DEFAULT_ADMIN_ROLE) {
        entryPoint.unlockStake();
    }

    function withdrawStake(address payable to) external onlyRole(DEFAULT_ADMIN_ROLE) {
        entryPoint.withdrawStake(to);
    }

    function getDeposit() external view returns (uint256) {
        return entryPoint.balanceOf(address(this));
    }

    function getDepositInfo() external view returns (IStakeManager.DepositInfo memory) {
        return entryPoint.getDepositInfo(address(this));
    }

    function supportsInterface(bytes4 interfaceId) public view override returns (bool) {
        return interfaceId == type(INotaryPaymaster).interfaceId || super.supportsInterface(interfaceId);
    }

    function _sponsorship(
        address sender,
        bytes calldata callData,
        bytes calldata paymasterData
    ) private view returns (uint48 validAfter, uint48 validUntil) {
        (address target, bytes calldata call) = _execution(callData);
        bytes4 selector = _selectorOf(call);
        if (target == address(withholding)) {
            return _withholdingSponsorship(sender, selector, call, paymasterData);
        }

        if (target == address(court)) {
            _checkCourtCall(selector, call);
            return (0, 0);
        }

        revert TargetNotSponsored(target);
    }

    function _withholdingSponsorship(
        address sender,
        bytes4 selector,
        bytes calldata call,
        bytes calldata paymasterData
    ) private view returns (uint48 validAfter, uint48 validUntil) {
        if (selector == IWithholdingCourt.requestStatement.selector) {
            (Registration memory registration, bytes memory signature, uint64 epoch) =
                abi.decode(call[SELECTOR_SIZE:], (Registration, bytes, uint64));
            bytes32 digest = _requireRegistration(sender, registration, signature);
            return (0, _requestClosesAt(digest, registration, epoch));
        }

        if (selector == IWithholdingCourt.recordUnanswered.selector) {
            (Registration memory registration, uint64 epoch) = abi.decode(call[SELECTOR_SIZE:], (Registration, uint64));
            bytes32 digest = _requireRegistration(sender, registration, paymasterData);
            return (_overdueAfter(digest, epoch), 0);
        }

        revert CallNotSponsored(address(withholding), selector);
    }

    function _requestClosesAt(
        bytes32 digest,
        Registration memory registration,
        uint64 epoch
    ) private view returns (uint48) {
        if (epoch < registration.fromEpoch) {
            revert BeforeRegistration(epoch, registration.fromEpoch);
        }

        uint64 anchoredAt = _anchoredAt(epoch);
        OpenRequest memory open = withholding.openRequestOf(digest);
        if (open.active) {
            revert RequestStillOpen(digest, open.epoch);
        }

        if (withholding.requestOf(digest, epoch).state != RequestState.None) {
            revert AlreadyRequested(digest, epoch);
        }

        return SafeCast.toUint48(uint256(anchoredAt) + _requestWindow);
    }

    function _overdueAfter(bytes32 digest, uint64 epoch) private view returns (uint48) {
        StatementRequest memory request = withholding.requestOf(digest, epoch);
        if (request.state != RequestState.Open) {
            revert NoOpenRequest(digest, epoch);
        }

        return SafeCast.toUint48(request.dueAt);
    }

    function _checkCourtCall(bytes4 selector, bytes calldata call) private view {
        if (selector == IStatementCourt.proveBadProof.selector) {
            (EpochStatement memory statement, bytes memory signature, SumNode[] memory siblings) =
                abi.decode(call[SELECTOR_SIZE:], (EpochStatement, bytes, SumNode[]));
            _checkBadProof(statement, signature, siblings);
            return;
        }

        if (selector == IStatementCourt.proveEquivocation.selector) {
            (
                EpochStatement memory first,
                bytes memory firstSignature,
                EpochStatement memory second,
                bytes memory secondSignature
            ) = abi.decode(call[SELECTOR_SIZE:], (EpochStatement, bytes, EpochStatement, bytes));
            _checkEquivocation(first, firstSignature, second, secondSignature);
            return;
        }

        revert CallNotSponsored(address(court), selector);
    }

    function _checkBadProof(
        EpochStatement memory statement,
        bytes memory signature,
        SumNode[] memory siblings
    ) private view {
        bytes32 digest = _requireStatement(statement, signature);
        bytes32 given = StatementHashing.siblingsHash(siblings);
        if (given != statement.siblingsHash) {
            revert SiblingsMismatch(statement.siblingsHash, given);
        }

        if (_ledger.verifyInclusion(statement.epoch, statement.commitment, statement.amount, statement.index, siblings))
        {
            revert ProofHolds(statement.epoch);
        }

        _requireUnproven(digest);
    }

    function _checkEquivocation(
        EpochStatement memory first,
        bytes memory firstSignature,
        EpochStatement memory second,
        bytes memory secondSignature
    ) private view {
        if (first.epoch != second.epoch) {
            revert DifferentEpochs(first.epoch, second.epoch);
        }

        bytes32 firstDigest = _requireStatement(first, firstSignature);
        bytes32 secondDigest = _requireStatement(second, secondSignature);
        if (firstDigest == secondDigest || !StatementHashing.conflict(first, second)) {
            revert NoConflict();
        }

        _requireUnproven(
            firstDigest < secondDigest
                ? keccak256(abi.encodePacked(firstDigest, secondDigest))
                : keccak256(abi.encodePacked(secondDigest, firstDigest))
        );
    }

    function _requireRegistration(
        address sender,
        Registration memory registration,
        bytes memory signature
    ) private view returns (bytes32) {
        if (registration.holder != sender) {
            revert NotTheHolder(registration.holder, sender);
        }

        bytes32 digest = _requireSigned(_withholdingDomain, RegistrationHashing.structHash(registration), signature);
        _requireSignersEpoch(registration.fromEpoch);
        return digest;
    }

    function _requireStatement(EpochStatement memory statement, bytes memory signature) private view returns (bytes32) {
        if (statement.ledger != address(_ledger)) {
            revert WrongLedger(address(_ledger), statement.ledger);
        }

        _requireAnchored(statement.epoch);
        bytes32 digest = _requireSigned(_courtDomain, StatementHashing.structHash(statement), signature);
        _requireSignersEpoch(statement.epoch);
        return digest;
    }

    function _requireSigned(bytes32 domain, bytes32 structHash, bytes memory signature) private view returns (bytes32) {
        bytes32 digest = MessageHashUtils.toTypedDataHash(domain, structHash);
        (address recovered, ECDSA.RecoverError failure, bytes32 detail) = ECDSA.tryRecover(digest, signature);
        if (failure != ECDSA.RecoverError.NoError || detail != bytes32(0) || recovered != statementSigner) {
            revert NotSignedByTheNotary(recovered);
        }

        return digest;
    }

    function _requireSignersEpoch(uint64 epoch) private view {
        address signer = court.signerOf(epoch);
        if (signer != statementSigner) {
            revert AnotherSignersEpoch(epoch, signer);
        }
    }

    function _requireUnproven(bytes32 evidence) private view {
        if (court.isProven(evidence)) {
            revert AlreadyProven(evidence);
        }
    }

    function _requireAnchored(uint64 epoch) private view {
        if (epoch < _ledger.genesisEpoch() || epoch >= _ledger.nextEpoch()) {
            revert EpochNotAnchored(epoch);
        }
    }

    function _anchoredAt(uint64 epoch) private view returns (uint64) {
        _requireAnchored(epoch);
        return _ledger.epochAt(epoch).anchoredAt;
    }

    function _execution(bytes calldata callData) private pure returns (address target, bytes calldata call) {
        bytes4 selector = _selectorOf(callData);
        if (selector != BaseAccount.execute.selector) {
            revert NotAnExecute(selector);
        }

        bytes calldata arguments = callData[SELECTOR_SIZE:];
        uint256 value;
        uint256 offset;
        (target, value, offset) = abi.decode(arguments, (address, uint256, uint256));
        if (value != 0) {
            revert ValueNotZero(value);
        }

        uint256 start = offset + WORD_SIZE;
        uint256 length = uint256(bytes32(arguments[offset:start]));
        call = arguments[start:start + length];
        return (target, call);
    }

    function _selectorOf(bytes calldata data) private pure returns (bytes4) {
        if (data.length < SELECTOR_SIZE) {
            return bytes4(0);
        }

        return bytes4(data[:SELECTOR_SIZE]);
    }

    function _requireSameDomains(
        IStatementCourt court_,
        IWithholdingCourt withholding_,
        bytes32 courtDomain,
        bytes32 withholdingDomain
    ) private view {
        EpochStatement memory statement = EpochStatement({
            ledger: address(court_),
            epoch: type(uint64).max,
            commitment: SAMPLE,
            amount: type(uint64).max,
            index: type(uint64).max,
            siblingsHash: SAMPLE
        });
        Registration memory registration =
            Registration({holder: address(withholding_), account: SAMPLE, fromEpoch: type(uint64).max});
        if (
            _statementDigestOf(court_, statement)
                != MessageHashUtils.toTypedDataHash(courtDomain, StatementHashing.structHash(statement))
        ) {
            revert DomainMismatch(address(court_));
        }

        if (
            _registrationDigestOf(withholding_, registration)
                != MessageHashUtils.toTypedDataHash(withholdingDomain, RegistrationHashing.structHash(registration))
        ) {
            revert DomainMismatch(address(withholding_));
        }
    }

    function _domain(
        string memory name,
        string memory version,
        address verifyingContract
    ) private view returns (bytes32) {
        return keccak256(
            abi.encode(
                DOMAIN_TYPEHASH, keccak256(bytes(name)), keccak256(bytes(version)), block.chainid, verifyingContract
            )
        );
    }

    function _requireCode(address target) private view {
        if (target.code.length == 0) {
            revert NoContract(target);
        }
    }

    function _courtOf(IWithholdingCourt withholding_) private view returns (IStatementCourt) {
        return withholding_.court();
    }

    function _ledgerOf(IStatementCourt court_) private view returns (ILiabilityLedger) {
        return court_.ledger();
    }

    function _requestWindowOf(IWithholdingCourt withholding_) private view returns (uint64) {
        return withholding_.requestWindow();
    }

    function _currentSignerOf(IStatementCourt court_) private view returns (address) {
        return court_.signerOf(type(uint64).max);
    }

    function _statementDigestOf(
        IStatementCourt court_,
        EpochStatement memory statement
    ) private view returns (bytes32) {
        return court_.statementDigest(statement);
    }

    function _registrationDigestOf(
        IWithholdingCourt withholding_,
        Registration memory registration
    ) private view returns (bytes32) {
        return withholding_.registrationDigest(registration);
    }
}
