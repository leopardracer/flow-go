package epochs

import (
	"testing"

	"github.com/stretchr/testify/require"
	"pgregory.net/rapid"

	"github.com/onflow/flow-go/model/flow"
	"github.com/onflow/flow-go/state/protocol"
	"github.com/onflow/flow-go/utils/unittest"
)

// TestEjectorRapid performs a rapid check on the ejector structure, ensuring that it correctly tracks and ejects nodes.
// This test covers only happy-path scenario.
func TestEjectorRapid(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		ej := ejector{}
		trackedIdentities := rapid.Map(rapid.SliceOfN(rapid.IntRange(3, 10), 1, 3), func(n []int) []flow.DynamicIdentityEntryList {
			var result []flow.DynamicIdentityEntryList
			for _, count := range n {
				result = append(result, unittest.DynamicIdentityEntryListFixture(count))
			}
			return result
		}).Draw(t, "tracked-identities")

		for _, list := range trackedIdentities {
			err := ej.TrackDynamicIdentityList(list)
			require.NoError(t, err)
		}

		var ejectedIdentities []flow.Identifier
		for _, list := range trackedIdentities {
			nodeID := rapid.SampledFrom(list).Draw(t, "ejected-identity").NodeID
			require.True(t, ej.Eject(nodeID))
			ejectedIdentities = append(ejectedIdentities, nodeID)
		}

		for _, ejectedNodeID := range ejectedIdentities {
			ejected := false
			for _, list := range trackedIdentities {
				if entry, found := list.ByNodeID(ejectedNodeID); found && entry.Ejected {
					ejected = true
				}
			}
			require.True(t, ejected, "identity should be ejected in tracked list")
		}
	})
}

// TestEjector_ReadmitEjectedIdentity ensures that a node that was ejected cannot be readmitted with subsequent track requests.
func TestEjector_ReadmitEjectedIdentity(t *testing.T) {
	list := unittest.DynamicIdentityEntryListFixture(3)
	ej := ejector{}
	ejectedNodeID := list[0].NodeID
	require.NoError(t, ej.TrackDynamicIdentityList(list))
	require.True(t, ej.Eject(ejectedNodeID))
	readmit := append(unittest.DynamicIdentityEntryListFixture(3), &flow.DynamicIdentityEntry{
		NodeID:  ejectedNodeID,
		Ejected: false,
	})
	err := ej.TrackDynamicIdentityList(readmit)
	require.Error(t, err)
	require.True(t, protocol.IsInvalidServiceEventError(err))
}
