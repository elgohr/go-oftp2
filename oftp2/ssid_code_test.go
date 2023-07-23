package oftp2_test

import (
	"github.com/elgohr/go-oftp2/oftp2"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestIdentificationCode(t *testing.T) {
	for _, scenario := range []struct {
		with   string
		input  func() oftp2.SsidIdentificationCodeInput
		expect func(t *testing.T, code oftp2.IdentificationCode, err error)
	}{
		{
			with: "a standard message",
			input: func() oftp2.SsidIdentificationCodeInput {
				return validSsidIdCodeInput
			},
			expect: func(t *testing.T, code oftp2.IdentificationCode, err error) {
				require.NotNil(t, code)
				require.NoError(t, err)
			},
		},
		{
			with: "expanding odette Id",
			input: func() oftp2.SsidIdentificationCodeInput {
				in := validSsidIdCodeInput
				in.OdetteIdentifier = "LONG"
				return in
			},
			expect: func(t *testing.T, code oftp2.IdentificationCode, err error) {
				require.EqualError(t, err, "exceeded capacity: LONG (1)")
				require.Nil(t, code)
			},
		},
		{
			with: "expanding InternationalCodeDesignator",
			input: func() oftp2.SsidIdentificationCodeInput {
				in := validSsidIdCodeInput
				in.InternationalCodeDesignator = "12345"
				return in
			},
			expect: func(t *testing.T, code oftp2.IdentificationCode, err error) {
				require.EqualError(t, err, "exceeded capacity: 12345 (4)")
				require.Nil(t, code)
			},
		},
		{
			with: "expanding OrganisationCode",
			input: func() oftp2.SsidIdentificationCodeInput {
				in := validSsidIdCodeInput
				in.OrganisationCode = "123456712345678"
				return in
			},
			expect: func(t *testing.T, code oftp2.IdentificationCode, err error) {
				require.EqualError(t, err, "exceeded capacity: 123456712345678 (14)")
				require.Nil(t, code)
			},
		},
		{
			with: "expanding ComputerSubaddress",
			input: func() oftp2.SsidIdentificationCodeInput {
				in := validSsidIdCodeInput
				in.ComputerSubaddress = "1234567"
				return in
			},
			expect: func(t *testing.T, code oftp2.IdentificationCode, err error) {
				require.EqualError(t, err, "exceeded capacity: 1234567 (6)")
				require.Nil(t, code)
			},
		},
	} {
		t.Run(scenario.with, func(t *testing.T) {
			code, err := oftp2.SsidIdentificationCode(scenario.input())
			scenario.expect(t, code, err)
		})
	}
}

func TestIdentificationCode_Valid(t *testing.T) {
	for _, scenario := range []struct {
		with   string
		input  func() oftp2.IdentificationCode
		expect func(t *testing.T, code oftp2.IdentificationCode)
	}{
		{
			with: "a standard message",
			input: func() oftp2.IdentificationCode {
				return validSsidCode(t)
			},
			expect: func(t *testing.T, code oftp2.IdentificationCode) {
				require.NoError(t, code.Valid())
				require.Equal(t, "X", string(code.OdetteIdentifier()))
				require.Equal(t, "1234", string(code.InternationalCodeDesignator()))
				require.Equal(t, "           ORG", string(code.OrganisationCode()))
				require.Equal(t, "abcdef", string(code.ComputerSubaddress()))
			},
		},
		{
			with: "expanding message",
			input: func() oftp2.IdentificationCode {
				return []byte("THIS_MESSAGE_IS_WAY_TOO_LONG")
			},
			expect: func(t *testing.T, code oftp2.IdentificationCode) {
				require.EqualError(t, code.Valid(), "expected the length of 25, but got 28")
			},
		},
	} {
		t.Run(scenario.with, func(t *testing.T) {
			scenario.expect(t, scenario.input())
		})
	}
}

var (
	validSsidIdCodeInput = oftp2.SsidIdentificationCodeInput{
		OdetteIdentifier:            "X",
		InternationalCodeDesignator: "1234",
		OrganisationCode:            "ORG",
		ComputerSubaddress:          "abcdef",
	}
	invalidSsidIdCode = oftp2.IdentificationCode("THIS_MESSAGE_IS_WAY_TOO_LONG")
)

func validSsidCode(t *testing.T) oftp2.IdentificationCode {
	session, err := oftp2.SsidIdentificationCode(validSsidIdCodeInput)
	require.NoError(t, err)
	return session
}
