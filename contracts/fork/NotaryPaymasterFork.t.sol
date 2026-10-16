// SPDX-License-Identifier: MIT
pragma solidity 0.8.37;

import {BaseAccount} from "account-abstraction/contracts/core/BaseAccount.sol";
import {IEntryPoint} from "account-abstraction/contracts/interfaces/IEntryPoint.sol";
import {IStakeManager} from "account-abstraction/contracts/interfaces/IStakeManager.sol";
import {PackedUserOperation} from "account-abstraction/contracts/interfaces/PackedUserOperation.sol";
import {Test} from "forge-std/src/Test.sol";

import {ILiabilityLedger} from "../src/ILiabilityLedger.sol";
import {StatementCourt} from "../src/court/StatementCourt.sol";
import {NotaryPaymaster} from "../src/paymaster/NotaryPaymaster.sol";
import {Registration} from "../src/withholding/Registration.sol";
import {RequestState} from "../src/withholding/RequestState.sol";
import {WithholdingCourt} from "../src/withholding/WithholdingCourt.sol";

contract NotaryPaymasterForkTest is Test {
    IEntryPoint private constant ENTRY_POINT = IEntryPoint(0x4337084D9E255Ff0702461CF8895CE9E3b5Ff108);
    address private constant ACCOUNT = 0x4Cd241E8d1510e30b2076397afc7508Ae59C66c9;
    bytes32 private constant ENTRY_POINT_CODEHASH = 0xa06a3208a3f27999ffae96d26acf86a40a16e582c5467513cbf2c7e96e68029c;
    bytes32 private constant ACCOUNT_CODEHASH = 0x82c1e6c0f83d22eef579344e8eff26baf24db4dabe5408d681b00d0512bc3ec4;

    ILiabilityLedger private ledger;
    WithholdingCourt private withholding;
    NotaryPaymaster private paymaster;
    address private signer;
    uint256 private signerKey;
    address private holder;
    uint256 private holderKey;
    address private admin;

    function setUp() public {
        vm.createSelectFork(vm.envString("KINETIX_NOTARY_FORK_RPC_URL"), vm.envUint("KINETIX_NOTARY_FORK_BLOCK"));
        ledger = ILiabilityLedger(vm.envAddress("KINETIX_NOTARY_FORK_LEDGER"));
        (signer, signerKey) = makeAddrAndKey("fork-paymaster-signer");
        (holder, holderKey) = makeAddrAndKey("fork-paymaster-holder");
        admin = makeAddr("fork-paymaster-admin");
        StatementCourt court = new StatementCourt(admin, 3 days, ledger, signer, 7 days);
        withholding = new WithholdingCourt(court, 7 days, 24 hours);
        paymaster = new NotaryPaymaster(admin, 3 days, ENTRY_POINT, ACCOUNT, court, withholding, 0.0001 ether, 1 days);
        vm.deal(address(this), 0.01 ether);
        ENTRY_POINT.depositTo{value: 0.01 ether}(address(paymaster));
    }

    function test_TheCanonicalEntryPointAndAccountAreTheOnesPinned() public view {
        assertEq(block.chainid, 84_532);
        assertEq(address(ENTRY_POINT).codehash, ENTRY_POINT_CODEHASH);
        assertEq(ACCOUNT.codehash, ACCOUNT_CODEHASH);
    }

    function test_TheAdminStakesAtTheCanonicalEntryPoint() public {
        vm.deal(admin, 1 ether);
        vm.prank(admin);
        paymaster.addStake{value: 1 ether}();

        IStakeManager.DepositInfo memory info = ENTRY_POINT.getDepositInfo(address(paymaster));
        assertTrue(info.staked);
        assertEq(info.stake, 1 ether);
        assertEq(info.unstakeDelaySec, 1 days);
        assertEq(paymaster.getDepositInfo().stake, 1 ether);
    }

    function test_AHolderWithNoEtherRequestsAStatementOnItsFirstDelegatedOperation() public {
        vm.signAndAttachDelegation(ACCOUNT, holderKey);
        uint64 epoch = ledger.nextEpoch() - 1;
        Registration memory registration =
            Registration({holder: holder, account: keccak256("fork-paymaster-account"), fromEpoch: 0});
        (uint8 v, bytes32 r, bytes32 s) = vm.sign(signerKey, withholding.registrationDigest(registration));
        bytes memory call =
            abi.encodeCall(withholding.requestStatement, (registration, abi.encodePacked(r, s, v), epoch));

        PackedUserOperation memory op;
        op.sender = holder;
        op.nonce = ENTRY_POINT.getNonce(holder, 0);
        op.initCode = abi.encodePacked(bytes20(bytes2(0x7702)));
        op.callData = abi.encodeCall(BaseAccount.execute, (address(withholding), 0, call));
        op.accountGasLimits = bytes32(uint256(150_000) << 128 | 300_000);
        op.preVerificationGas = 100_000;
        op.gasFees = bytes32(uint256(0.01 gwei) << 128 | 0.01 gwei);
        op.paymasterAndData = abi.encodePacked(address(paymaster), uint128(100_000), uint128(0));
        (v, r, s) = vm.sign(holderKey, ENTRY_POINT.getUserOpHash(op));
        op.signature = abi.encodePacked(r, s, v);
        uint256 before = paymaster.getDeposit();

        PackedUserOperation[] memory ops = new PackedUserOperation[](1);
        ops[0] = op;
        address bundler = makeAddr("fork-bundler");
        vm.prank(bundler, bundler);
        ENTRY_POINT.handleOps(ops, payable(bundler));

        bytes32 digest = withholding.registrationDigest(registration);
        assertEq(uint8(withholding.requestOf(digest, epoch).state), uint8(RequestState.Open));
        assertEq(holder.balance, 0);
        assertLt(paymaster.getDeposit(), before);
    }
}
