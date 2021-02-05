package oftp2_test

import (
	cmd "bifroest/oftp2"
	oftp2 "bifroest/oftp2/cmd"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestStartSessionErrors(t *testing.T) {
	for _, scenario := range []struct {
		name   string
		input  func() oftp2.StartSessionInput
		expect func(t *testing.T, cmd cmd.Command, err error)
	}{
		{
			name: "with a standard input",
			input: func() oftp2.StartSessionInput {
				return oftp2.StartSessionInput{
					IdentificationCode: validSsidIdCode,
					DataExchangeBufferSize: 99999,
					Capabilities: oftp2.CapabilityReceive,
				}
			},
			expect: func(t *testing.T, cmd cmd.Command, err error) {
				require.NoError(t, err)
				require.NotNil(t, cmd)
			},
		},
		{
			name: "with invalid ssid id code",
			input: func() oftp2.StartSessionInput {
				return oftp2.StartSessionInput{
					IdentificationCode: invalidSsidIdCode,
				}
			},
			expect: func(t *testing.T, cmd cmd.Command, err error) {
				require.Error(t, err)
				require.Nil(t, cmd)
			},
		},
		{
			name: "with invalid ssid id code",
			input: func() oftp2.StartSessionInput {
				return oftp2.StartSessionInput{
					IdentificationCode: nil,
				}
			},
			expect: func(t *testing.T, cmd cmd.Command, err error) {
				require.EqualError(t, err, "missing identification code")
				require.Nil(t, cmd)
			},
		},
		{
			name: "with invalid password",
			input: func() oftp2.StartSessionInput {
				return oftp2.StartSessionInput{
					IdentificationCode: validSsidIdCode,
					Password: "WAY_TOO_LONG_RIGHT?",
				}
			},
			expect: func(t *testing.T, cmd cmd.Command, err error) {
				require.EqualError(t, err, "exceeded capacity: WAY_TOO_LONG_RIGHT? (8)")
				require.Nil(t, cmd)
			},
		},
		{
			name: "with invalid data exchange buffer",
			input: func() oftp2.StartSessionInput {
				return oftp2.StartSessionInput{
					IdentificationCode: validSsidIdCode,
					Password: "PASSWORD",
					DataExchangeBufferSize: 100000,
				}
			},
			expect: func(t *testing.T, cmd cmd.Command, err error) {
				require.EqualError(t, err, "exceeded capacity: 100000 (5)")
				require.Nil(t, cmd)
			},
		},
		{
			name: "with invalid credit",
			input: func() oftp2.StartSessionInput {
				return oftp2.StartSessionInput{
					IdentificationCode: validSsidIdCode,
					Password: "PASSWORD",
					Capabilities: oftp2.CapabilityBoth,
					DataExchangeBufferSize: 99999,
					Credit: 9999,
				}
			},
			expect: func(t *testing.T, cmd cmd.Command, err error) {
				require.EqualError(t, err, "exceeded capacity: 9999 (3)")
				require.Nil(t, cmd)
			},
		},
		{
			name: "with invalid user data",
			input: func() oftp2.StartSessionInput {
				return oftp2.StartSessionInput{
					IdentificationCode: validSsidIdCode,
					Password: "PASSWORD",
					Capabilities: oftp2.CapabilityReceive,
					DataExchangeBufferSize: 99999,
					Credit: 999,
					UserData: "12345678910",
				}
			},
			expect: func(t *testing.T, cmd cmd.Command, err error) {
				require.EqualError(t, err, "exceeded capacity: 12345678910 (8)")
				require.Nil(t, cmd)
			},
		},
		{
			name: "with invalid capabilities",
			input: func() oftp2.StartSessionInput {
				return oftp2.StartSessionInput{
					IdentificationCode: validSsidIdCode,
					Password: "PASSWORD",
					DataExchangeBufferSize: 99999,
					Capabilities: "T",
				}
			},
			expect: func(t *testing.T, cmd cmd.Command, err error) {
				require.EqualError(t, err, "unknown capability: T")
				require.Nil(t, cmd)
			},
		},
	} {
		t.Run(scenario.name, func(t *testing.T) {
			session, err := oftp2.StartSession(scenario.input())
			scenario.expect(t, session, err)
		})
	}
}

func TestStartSessionCmd_Valid(t *testing.T) {
	for _, scenario := range []struct {
		name   string
		input  func(t *testing.T) []byte
		expect func(t *testing.T, ssid oftp2.StartSessionCmd)
	}{
		{
			name: "with a standard message",
			input: func(t *testing.T) []byte {
				return validSession(t)
			},
			expect: func(t *testing.T, ssid oftp2.StartSessionCmd) {
				require.Equal(t, "X", string(ssid.Command()))
				require.Equal(t, "5", string(ssid.ProtocolLevel()))
				c := ssid.IdentificationCode()
				require.NoError(t, c.Valid())
				require.Equal(t, "password", string(ssid.Password()))
				require.Equal(t, 99999, ssid.DataExchangeBufferSize())
				require.Equal(t, "B", string(ssid.Capabilities()))
				require.True(t, ssid.BufferCompression())
				require.True(t, ssid.Restart())
				require.True(t, ssid.SpecialLogic())
				require.Equal(t, 999, ssid.Credit())
				require.True(t, ssid.Authentication())
				require.Equal(t, "        ", string(ssid.User()))
			},
		},
		{
			name: "with a wrong command",
			input: func(t *testing.T) []byte {
				session := validSession(t)
				session[0] = 'Y'
				return session
			},
			expect: func(t *testing.T, ssid oftp2.StartSessionCmd) {
				require.Error(t, ssid.Valid())
				require.Equal(t, "Y", string(ssid.Command()))
				require.Equal(t, "5", string(ssid.ProtocolLevel()))
			},
		},
		{
			name: "with different protocol level",
			input: func(t *testing.T) []byte {
				session := validSession(t)
				session[1] = '3'
				return session
			},
			expect: func(t *testing.T, ssid oftp2.StartSessionCmd) {
				require.Error(t, ssid.Valid())
				require.Equal(t, "3", string(ssid.ProtocolLevel()))
			},
		},
		{
			name: "with invalid ssid code",
			input: func(t *testing.T) []byte {
				session := validSession(t)
				session[4] = ')'
				return session
			},
			expect: func(t *testing.T, ssid oftp2.StartSessionCmd) {
				require.Error(t, ssid.Valid()) // error from ssid code
				require.Equal(t, "X", string(ssid.Command()))
				require.Equal(t, "5", string(ssid.ProtocolLevel()))
			},
		},
		{
			name: "with invalid data exchange buffer size",
			input: func(t *testing.T) []byte {
				session := validSession(t)
				session[36] = ' '
				return session
			},
			expect: func(t *testing.T, ssid oftp2.StartSessionCmd) {
				require.Error(t, ssid.Valid())
				require.Equal(t, 0, ssid.DataExchangeBufferSize())
			},
		},
		{
			name: "with data exchange buffer size at lower end",
			input: func(t *testing.T) []byte {
				input := validSsidInput
				input.DataExchangeBufferSize = 128
				session, err := oftp2.StartSession(input)
				require.NoError(t, err)
				return session
			},
			expect: func(t *testing.T, ssid oftp2.StartSessionCmd) {
				require.NoError(t, ssid.Valid())
				require.Equal(t, 128, ssid.DataExchangeBufferSize())
			},
		},
		{
			name: "with data exchange buffer size below lower end",
			input: func(t *testing.T) []byte {
				input := validSsidInput
				input.DataExchangeBufferSize = 127
				session, err := oftp2.StartSession(input)
				require.NoError(t, err)
				return session
			},
			expect: func(t *testing.T, ssid oftp2.StartSessionCmd) {
				require.EqualError(t, ssid.Valid(), "invalid DataExchangeBufferSize: 127")
			},
		},
		{
			name: "with data exchange buffer size at upper end",
			input: func(t *testing.T) []byte {
				input := validSsidInput
				input.DataExchangeBufferSize = 99999
				session, err := oftp2.StartSession(input)
				require.NoError(t, err)
				return session
			},
			expect: func(t *testing.T, ssid oftp2.StartSessionCmd) {
				require.NoError(t, ssid.Valid())
				require.Equal(t, 99999, ssid.DataExchangeBufferSize())
			},
		},
		{
			name: "with unknown capability",
			input: func(t *testing.T) []byte {
				session := validSession(t)
				session[40] = 'U'
				return session
			},
			expect: func(t *testing.T, ssid oftp2.StartSessionCmd) {
				require.EqualError(t, ssid.Valid(), "unknown capability: U")
			},
		},
		{
			name: "with unknown compression indicator",
			input: func(t *testing.T) []byte {
				session := validSession(t)
				session[41] = 'U'
				return session
			},
			expect: func(t *testing.T, ssid oftp2.StartSessionCmd) {
				require.EqualError(t, ssid.Valid(), "unknown BufferCompressionIndicator: U")
			},
		},
		{
			name: "with unknown restart indicator",
			input: func(t *testing.T) []byte {
				session := validSession(t)
				session[42] = 'U'
				return session
			},
			expect: func(t *testing.T, ssid oftp2.StartSessionCmd) {
				require.EqualError(t, ssid.Valid(), "unknown RestartIndicator: U")
			},
		},
		{
			name: "with unknown special logic indicator",
			input: func(t *testing.T) []byte {
				session := validSession(t)
				session[43] = 'U'
				return session
			},
			expect: func(t *testing.T, ssid oftp2.StartSessionCmd) {
				require.EqualError(t, ssid.Valid(), "unknown SpecialLogicIndicator: U")
			},
		},
		{
			name: "with invalid credit",
			input: func(t *testing.T) []byte {
				session := validSession(t)
				session[45] = 'U'
				return session
			},
			expect: func(t *testing.T, ssid oftp2.StartSessionCmd) {
				require.EqualError(t, ssid.Valid(), `invalid Credit: strconv.Atoi: parsing "9U9": invalid syntax`)
			},
		},
		{
			name: "with negative credit",
			input: func(t *testing.T) []byte {
				session := validSession(t)
				session[44] = '-'
				return session
			},
			expect: func(t *testing.T, ssid oftp2.StartSessionCmd) {
				require.EqualError(t, ssid.Valid(), `invalid Credit: -99`)
			},
		},
		{
			name: "with unknown authentication",
			input: func(t *testing.T) []byte {
				session := validSession(t)
				session[47] = 'U'
				return session
			},
			expect: func(t *testing.T, ssid oftp2.StartSessionCmd) {
				require.EqualError(t, ssid.Valid(), `unknown Authentication: U`)
			},
		},
	} {
		t.Run(scenario.name, func(t *testing.T) {
			scenario.expect(t, scenario.input(t))
		})
	}
}

var (
	validSsidIdCode   = oftp2.SsidIdentificationCode("X", 1234, "ORG           ", "abcdef")
	invalidSsidIdCode = oftp2.SsidIdentificationCode("XA", 1234, "ORG           ", "abcdef")

	validSsidInput = oftp2.StartSessionInput{
		IdentificationCode:     validSsidIdCode,
		Password:               "password",
		DataExchangeBufferSize: 99999,
		Capabilities:           oftp2.CapabilityBoth,
		BufferCompression:      true,
		Restart:                true,
		SpecialLogic:           true,
		Credit:                 999,
		SecureAuthentication:   true,
		UserData:               "        ",
	}
)

func validSession(t *testing.T) cmd.Command {
	session, err := oftp2.StartSession(validSsidInput)
	require.NoError(t, err)
	return session
}
