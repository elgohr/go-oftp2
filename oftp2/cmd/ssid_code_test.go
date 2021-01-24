package oftp2_test

import (
	oftp2 "bifroest/oftp2/cmd"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestIdentificationCode(t *testing.T) {
	for _, scenario := range []struct {
		name   string
		input  func() []byte
		expect func(t *testing.T, code oftp2.IdentificationCode)
	}{
		{
			name: "with a standard message",
			input: func() []byte {
				return oftp2.SsidIdentificationCode("X", 1234, "ORG           ", "abcdef")
			},
			expect: func(t *testing.T, code oftp2.IdentificationCode) {
				require.NoError(t, code.Valid())
				require.Equal(t, "X", string(code.Oid()))
				require.Equal(t, 1234, code.Icd())
				require.Equal(t, "ORG           ", string(code.Org()))
				require.Equal(t, "abcdef", string(code.Csa()))
			},
		},
		{
			name: "with an expanding message",
			input: func() []byte {
				return oftp2.SsidIdentificationCode("XA", 1234, "ORG          ", "abcdef")
			},
			expect: func(t *testing.T, code oftp2.IdentificationCode) {
				require.Error(t, code.Valid())
				require.Equal(t, 0, code.Icd())
			},
		},
		{
			name: "with a corrupt Icd",
			input: func() []byte {
				c := oftp2.SsidIdentificationCode("X", 1234, "ORG           ", "abcdef")
				c[2] = 'A'
				return c
			},
			expect: func(t *testing.T, code oftp2.IdentificationCode) {
				require.EqualError(t, code.Valid(), "international code designator is not a number, but 1A")
			},
		},
	} {
		t.Run(scenario.name, func(t *testing.T) {
			scenario.expect(t, scenario.input())
		})
	}
}
