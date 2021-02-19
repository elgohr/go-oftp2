package oftp2_test

import (
	"bifroest/oftp2"
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
				return oftp2.NewStartSessionReadyMessage()
			},
			expect: func(t *testing.T, ssrm oftp2.StartSessionReadyMessageCmd) {
				require.NoError(t, ssrm.Valid())
				require.Equal(t, "ODETTE FTP READY ", string(ssrm.Message()))
			},
		},
		{
			name: "with a wrong cmd id",
			input: func() []byte {
				m := oftp2.NewStartSessionReadyMessage()
				m[0] = 'X'
				return m
			},
			expect: func(t *testing.T, ssrm oftp2.StartSessionReadyMessageCmd) {
				require.EqualError(t, ssrm.Valid(), "does not start with I, but with X")
				require.Equal(t, "ODETTE FTP READY ", string(ssrm.Message()))
			},
		},
		{
			name: "with missing CR",
			input: func() []byte {
				m := oftp2.NewStartSessionReadyMessage()
				m[18] = ' '
				return m
			},
			expect: func(t *testing.T, ssrm oftp2.StartSessionReadyMessageCmd) {
				require.EqualError(t, ssrm.Valid(), "does not end on carriage return, but on 32")
				require.Equal(t, "ODETTE FTP READY ", string(ssrm.Message()))
			},
		},
		{
			name: "with an exceeding message",
			input: func() []byte {
				m := oftp2.NewStartSessionReadyMessage()
				return append(m, ' ')
			},
			expect: func(t *testing.T, ssrm oftp2.StartSessionReadyMessageCmd) {
				require.EqualError(t, ssrm.Valid(), "expected the length of 19, but got 20")
				require.Equal(t, "ODETTE FTP READY ", string(ssrm.Message()))
			},
		},
	} {
		t.Run(scenario.name, func(t *testing.T) {
			scenario.expect(t, scenario.input())
		})
	}
}
