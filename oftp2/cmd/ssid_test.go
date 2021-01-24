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
				return oftp2.StartSession(validSsidIdCode, "password")
			},
			expect: func(t *testing.T, ssid oftp2.StartSessionCmd) {
				require.Equal(t, "X", string(ssid.Cmd()))
				require.Equal(t, "5", string(ssid.Lev()))
				c := ssid.Code()
				require.NoError(t, c.Valid())
				require.Equal(t, "password", string(ssid.Pswd()))
				require.Equal(t, 99999, ssid.Sdeb())
			},
		},
		{
			name: "with a wrong cmd id",
			input: func() []byte {
				m := oftp2.StartSession(validSsidIdCode, "password")
				m[0] = 'Y'
				return m
			},
			expect: func(t *testing.T, ssid oftp2.StartSessionCmd) {
				require.Error(t, ssid.Valid())
				require.Equal(t, "Y", string(ssid.Cmd()))
				require.Equal(t, "5", string(ssid.Lev()))
			},
		},
		//{
		//	name: "with missing CR",
		//	input: func() []byte {
		//		m := oftp2.StartSession(validSsidIdCode, "password")
		//		m[len(m)-1] = ' '
		//		return m
		//	},
		//	expect: func(t *testing.T, ssid oftp2.StartSessionCmd) {
		//		require.Error(t, ssid.Valid())
		//		require.Equal(t, "X", string(ssid.Cmd()))
		//		require.Equal(t, "5", string(ssid.Lev()))
		//	},
		//},
		{
			name: "with invalid ssid code",
			input: func() []byte {
				m := oftp2.StartSession(invalidSsidIdCode, "password")
				//m[18] = ' '
				return m
			},
			expect: func(t *testing.T, ssid oftp2.StartSessionCmd) {
				require.Error(t, ssid.Valid()) // error from ssid code
				require.Equal(t, "X", string(ssid.Cmd()))
				require.Equal(t, "5", string(ssid.Lev()))
			},
		},
		{
			name: "with invalid data exchange buffer size",
			input: func() []byte {
				m := oftp2.StartSession(invalidSsidIdCode, "password")
				m[36] = ' '
				return m
			},
			expect: func(t *testing.T, ssid oftp2.StartSessionCmd) {
				require.Error(t, ssid.Valid())
				require.Equal(t, 0, ssid.Sdeb())
			},
		},
	} {
		t.Run(scenario.name, func(t *testing.T) {
			scenario.expect(t, scenario.input())
		})
	}
}
