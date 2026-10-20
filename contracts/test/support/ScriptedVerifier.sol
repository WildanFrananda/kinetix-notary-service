// SPDX-License-Identifier: MIT
pragma solidity 0.8.37;

contract ScriptedVerifier {
    error Scripted();

    mapping(bytes32 call => bool) private _accepted;
    bool private _reverts;

    function accept(bytes calldata proof, uint256[] calldata publicInputs) external {
        _accepted[keccak256(abi.encode(proof, publicInputs))] = true;
    }

    function revertEverything() external {
        _reverts = true;
    }

    function Verify(bytes calldata proof, uint256[] calldata publicInputs) external view returns (bool) {
        if (_reverts) {
            revert Scripted();
        }
        return _accepted[keccak256(abi.encode(proof, publicInputs))];
    }
}
