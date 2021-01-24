package oftp2_test

import (
	oftp2 "bifroest/oftp2/cmd"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestStartSessionReadyMessage(t *testing.T) {
	for _, scenario := range []struct {
		name   string
		input  func() []byte
		expect func(t *testing.T, ssrm oftp2.StartSessionReadyMessageCmd)
	}{
		{
			name: "with a standard message",
			input: func() []byte {
				return oftp2.StartSessionReadyMessage()
			},
			expect: func(t *testing.T, ssrm oftp2.StartSessionReadyMessageCmd) {
				require.NoError(t, ssrm.Valid())
				require.Equal(t, "I", string(ssrm.Cmd()))
				require.Equal(t, "ODETTE FTP READY ", string(ssrm.Msg()))
			},
		},
		{
			name: "with a wrong cmd id",
			input: func() []byte {
				m := oftp2.StartSessionReadyMessage()
				m[0] = 'X'
				return m
			},
			expect: func(t *testing.T, ssrm oftp2.StartSessionReadyMessageCmd) {
				require.Error(t, ssrm.Valid())
				require.Equal(t, "X", string(ssrm.Cmd()))
				require.Equal(t, "ODETTE FTP READY ", string(ssrm.Msg()))
			},
		},
		{
			name: "with missing CR",
			input: func() []byte {
				m := oftp2.StartSessionReadyMessage()
				m[18] = ' '
				return m
			},
			expect: func(t *testing.T, ssrm oftp2.StartSessionReadyMessageCmd) {
				require.Error(t, ssrm.Valid())
				require.Equal(t, "I", string(ssrm.Cmd()))
				require.Equal(t, "ODETTE FTP READY ", string(ssrm.Msg()))
			},
		},
	} {
		t.Run(scenario.name, func(t *testing.T) {
			scenario.expect(t, scenario.input())
		})
	}
}
