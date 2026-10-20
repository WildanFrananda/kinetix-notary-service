package sources_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

var (
	councilPath             = filepath.Join("..", "..", "contracts", "src", "council", "AnchorCouncil.sol")
	confidentialCouncilPath = filepath.Join(
		"..",
		"..",
		"contracts",
		"src",
		"confidential",
		"ConfidentialAnchorCouncil.sol",
	)
)

var (
	councilScriptPath             = filepath.Join("..", "..", "contracts", "script", "DeployAnchorCouncil.s.sol")
	confidentialCouncilScriptPath = filepath.Join(
		"..",
		"..",
		"contracts",
		"script",
		"DeployConfidentialAnchorCouncil.s.sol",
	)
)

type substitution struct {
	old   string
	new   string
	count int
}

func TestTheSchemeTwoCouncilIsTheCouncilWithOnlyItsDigestAndLedgerChanged(t *testing.T) {
	original, err := os.ReadFile(councilPath)
	require.NoError(t, err)

	derived, err := os.ReadFile(confidentialCouncilPath)
	require.NoError(t, err)

	want := string(original)
	for _, s := range councilSubstitutions() {
		require.Equal(t, s.count, strings.Count(want, s.old), "AnchorCouncil.sol changed under %q", s.old)

		want = strings.ReplaceAll(want, s.old, s.new)
	}

	require.Equal(
		t,
		want,
		string(derived),
		"ConfidentialAnchorCouncil.sol must be AnchorCouncil.sol with these substitutions and nothing else",
	)
}

func TestTheSchemeTwoCouncilsDeployScriptIsTheCouncilsWithOnlyItsTypeChanged(t *testing.T) {
	original, err := os.ReadFile(councilScriptPath)
	require.NoError(t, err)

	derived, err := os.ReadFile(confidentialCouncilScriptPath)
	require.NoError(t, err)

	want := strings.ReplaceAll(string(original), "AnchorCouncil", "ConfidentialAnchorCouncil")
	want = strings.Replace(
		want,
		"\"../src/council/ConfidentialAnchorCouncil.sol\"",
		"\"../src/confidential/ConfidentialAnchorCouncil.sol\"",
		1,
	)

	require.Equal(t, 7, strings.Count(string(original), "AnchorCouncil"), "DeployAnchorCouncil.s.sol changed")
	require.Equal(
		t,
		want,
		string(derived),
		"DeployConfidentialAnchorCouncil.s.sol must be DeployAnchorCouncil.s.sol with the council's type changed",
	)
}

func councilSubstitutions() []substitution {
	return []substitution{
		{
			old: "import {EpochChain} from \"../EpochChain.sol\";\n" +
				"import {EpochHeader} from \"../EpochHeader.sol\";\n" +
				"import {ILiabilityLedger} from \"../ILiabilityLedger.sol\";\n" +
				"import {ApprovalHashing} from \"./ApprovalHashing.sol\";\n" +
				"import {IAnchorCouncil} from \"./IAnchorCouncil.sol\";\n" +
				"import {PendingMembership} from \"./PendingMembership.sol\";\n",
			new: "import {EpochHeader} from \"../EpochHeader.sol\";\n" +
				"import {ILiabilityLedger} from \"../ILiabilityLedger.sol\";\n" +
				"import {ApprovalHashing} from \"../council/ApprovalHashing.sol\";\n" +
				"import {IAnchorCouncil} from \"../council/IAnchorCouncil.sol\";\n" +
				"import {PendingMembership} from \"../council/PendingMembership.sol\";\n" +
				"import {ConfidentialEpochChain} from \"./ConfidentialEpochChain.sol\";\n" +
				"import {IConfidentialLedger} from \"./IConfidentialLedger.sol\";\n",
			count: 1,
		},
		{
			old:   "contract AnchorCouncil is",
			new:   "contract ConfidentialAnchorCouncil is",
			count: 1,
		},
		{
			old:   "EpochChain.digest(header)",
			new:   "ConfidentialEpochChain.digest(header)",
			count: 2,
		},
		{
			old: "        return ERC165Checker.supportsInterface(address(candidate), type(ILiabilityLedger).interfaceId);\n",
			new: "        return ERC165Checker.supportsInterface(address(candidate), type(ILiabilityLedger).interfaceId)\n" +
				"            && ERC165Checker.supportsInterface(address(candidate), type(IConfidentialLedger).interfaceId);\n",
			count: 1,
		},
	}
}
