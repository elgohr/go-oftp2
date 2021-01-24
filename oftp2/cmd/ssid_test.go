package oftp2_test

import (
	oftp2 "bifroest/oftp2/cmd"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestStartSession(t *testing.T) {
	validSsidIdCode := oftp2.SsidIdentificationCode("X", 1234, "ORG           ", "abcdef")
	invalidSsidIdCode := oftp2.SsidIdentificationCode("XA", 1234, "ORG           ", "abcdef")

	for _, scenario := range []struct {
		name   string
		input  func() []byte
		expect func(t *testing.T, ssid oftp2.StartSessionCmd)
	}{
		{
			name: "with a standard message",
			input: func() []byte {
				return oftp2.StartSession(validSsidIdCode)
			},
			expect: func(t *testing.T, ssid oftp2.StartSessionCmd) {
				require.Equal(t, "X", string(ssid.Cmd()))
				require.Equal(t, "5", string(ssid.Lev()))
			},
		},
		{
			name: "with a wrong cmd id",
			input: func() []byte {
				m := oftp2.StartSession(nil)
				m[0] = 'Y'
				return m
			},
			expect: func(t *testing.T, ssid oftp2.StartSessionCmd) {
				require.Equal(t, "Y", string(ssid.Cmd()))
				require.Equal(t, "5", string(ssid.Lev()))
			},
		},
		{
			name: "with missing CR",
			input: func() []byte {
				m := oftp2.StartSession(nil)
				//m[18] = ' '
				return m
			},
			expect: func(t *testing.T, ssid oftp2.StartSessionCmd) {
				//require.Error(t, ssrm.Valid())
				require.Equal(t, "X", string(ssid.Cmd()))
				require.Equal(t, "5", string(ssid.Lev()))
			},
		},
	} {
		t.Run(scenario.name, func(t *testing.T) {
			scenario.expect(t, scenario.input())
		})
	}
}
