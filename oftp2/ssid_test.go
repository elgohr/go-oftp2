package oftp2_test

import (
	"github.com/elgohr/go-oftp2/oftp2"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestStartSessionErrors(t *testing.T) {
	for _, scenario := range []struct {
		with   string
		input  func(t *testing.T) oftp2.StartSessionInput
		expect func(t *testing.T, cmd oftp2.Command, err error)
	}{
		{
			with: "a standard input",
			input: func(t *testing.T) oftp2.StartSessionInput {
				return oftp2.StartSessionInput{
					IdentificationCode:     validSsidCode(t),
					DataExchangeBufferSize: 99999,
					Capabilities:           oftp2.CapabilityReceive,
				}
			},
			expect: func(t *testing.T, cmd oftp2.Command, err error) {
				require.NoError(t, err)
				require.NotNil(t, cmd)
			},
		},
		{
			with: "invalid ssid id code",
			input: func(t *testing.T) oftp2.StartSessionInput {
				return oftp2.StartSessionInput{
					IdentificationCode: invalidSsidIdCode,
				}
			},
			expect: func(t *testing.T, cmd oftp2.Command, err error) {
				require.Error(t, err)
				require.Nil(t, cmd)
			},
		},
		{
			with: "invalid ssid id code",
			input: func(t *testing.T) oftp2.StartSessionInput {
				return oftp2.StartSessionInput{
					IdentificationCode: nil,
				}
			},
			expect: func(t *testing.T, cmd oftp2.Command, err error) {
				require.EqualError(t, err, "missing identification code")
				require.Nil(t, cmd)
			},
		},
		{
			with: "invalid password",
			input: func(t *testing.T) oftp2.StartSessionInput {
				return oftp2.StartSessionInput{
					IdentificationCode: validSsidCode(t),
					Password:           "WAY_TOO_LONG_RIGHT?",
				}
			},
			expect: func(t *testing.T, cmd oftp2.Command, err error) {
				require.EqualError(t, err, "exceeded capacity: WAY_TOO_LONG_RIGHT? (8)")
				require.Nil(t, cmd)
			},
		},
		{
			with: "invalid data exchange buffer",
			input: func(t *testing.T) oftp2.StartSessionInput {
				return oftp2.StartSessionInput{
					IdentificationCode:     validSsidCode(t),
					Password:               "PASSWORD",
					DataExchangeBufferSize: 100000,
				}
			},
			expect: func(t *testing.T, cmd oftp2.Command, err error) {
				require.EqualError(t, err, "exceeded capacity: 100000 (5)")
				require.Nil(t, cmd)
			},
		},
		{
			with: "invalid credit",
			input: func(t *testing.T) oftp2.StartSessionInput {
				return oftp2.StartSessionInput{
					IdentificationCode:     validSsidCode(t),
					Password:               "PASSWORD",
					Capabilities:           oftp2.CapabilityBoth,
					DataExchangeBufferSize: 99999,
					Credit:                 9999,
				}
			},
			expect: func(t *testing.T, cmd oftp2.Command, err error) {
				require.EqualError(t, err, "exceeded capacity: 9999 (3)")
				require.Nil(t, cmd)
			},
		},
		{
			with: "invalid user data",
			input: func(t *testing.T) oftp2.StartSessionInput {
				return oftp2.StartSessionInput{
					IdentificationCode:     validSsidCode(t),
					Password:               "PASSWORD",
					Capabilities:           oftp2.CapabilityReceive,
					DataExchangeBufferSize: 99999,
					Credit:                 999,
					UserData:               "12345678910",
				}
			},
			expect: func(t *testing.T, cmd oftp2.Command, err error) {
				require.EqualError(t, err, "exceeded capacity: 12345678910 (8)")
				require.Nil(t, cmd)
			},
		},
		{
			with: "invalid capabilities",
			input: func(t *testing.T) oftp2.StartSessionInput {
				return oftp2.StartSessionInput{
					IdentificationCode:     validSsidCode(t),
					Password:               "PASSWORD",
					DataExchangeBufferSize: 99999,
					Capabilities:           "T",
				}
			},
			expect: func(t *testing.T, cmd oftp2.Command, err error) {
				require.EqualError(t, err, "unknown capability: T")
				require.Nil(t, cmd)
			},
		},
	} {
		t.Run(scenario.with, func(t *testing.T) {
			session, err := oftp2.NewStartSession(scenario.input(t))
			scenario.expect(t, session, err)
		})
	}
}

func TestStartSessionCmd_Valid(t *testing.T) {
	for _, scenario := range []struct {
		with   string
		input  func(t *testing.T) []byte
		expect func(t *testing.T, ssid oftp2.StartSessionCmd)
	}{
		{
			with: "a standard message",
			input: func(t *testing.T) []byte {
				return validSessionStart(t)
			},
			expect: func(t *testing.T, ssid oftp2.StartSessionCmd) {
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
			with: "a wrong command",
			input: func(t *testing.T) []byte {
				session := validSessionStart(t)
				session[0] = 'Y'
				return session
			},
			expect: func(t *testing.T, ssid oftp2.StartSessionCmd) {
				require.Error(t, ssid.Valid())
				require.Equal(t, "5", string(ssid.ProtocolLevel()))
			},
		},
		{
			with: "different protocol level",
			input: func(t *testing.T) []byte {
				session := validSessionStart(t)
				session[1] = '3'
				return session
			},
			expect: func(t *testing.T, ssid oftp2.StartSessionCmd) {
				require.Error(t, ssid.Valid())
				require.Equal(t, "3", string(ssid.ProtocolLevel()))
			},
		},
		{
			with: "invalid data exchange buffer size",
			input: func(t *testing.T) []byte {
				session := validSessionStart(t)
				session[36] = ' '
				return session
			},
			expect: func(t *testing.T, ssid oftp2.StartSessionCmd) {
				require.Error(t, ssid.Valid())
				require.Equal(t, 0, ssid.DataExchangeBufferSize())
			},
		},
		{
			with: "data exchange buffer size at lower end",
			input: func(t *testing.T) []byte {
				input := oftp2.StartSessionInput{
					IdentificationCode:     validSsidCode(t),
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
				input.DataExchangeBufferSize = 128
				session, err := oftp2.NewStartSession(input)
				require.NoError(t, err)
				return session
			},
			expect: func(t *testing.T, ssid oftp2.StartSessionCmd) {
				require.NoError(t, ssid.Valid())
				require.Equal(t, 128, ssid.DataExchangeBufferSize())
			},
		},
		{
			with: "data exchange buffer size below lower end",
			input: func(t *testing.T) []byte {
				input := oftp2.StartSessionInput{
					IdentificationCode:     validSsidCode(t),
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
				input.DataExchangeBufferSize = 127
				session, err := oftp2.NewStartSession(input)
				require.NoError(t, err)
				return session
			},
			expect: func(t *testing.T, ssid oftp2.StartSessionCmd) {
				require.EqualError(t, ssid.Valid(), "invalid DataExchangeBufferSize: 127")
			},
		},
		{
			with: "data exchange buffer size at upper end",
			input: func(t *testing.T) []byte {
				input := oftp2.StartSessionInput{
					IdentificationCode:     validSsidCode(t),
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
				input.DataExchangeBufferSize = 99999
				session, err := oftp2.NewStartSession(input)
				require.NoError(t, err)
				return session
			},
			expect: func(t *testing.T, ssid oftp2.StartSessionCmd) {
				require.NoError(t, ssid.Valid())
				require.Equal(t, 99999, ssid.DataExchangeBufferSize())
			},
		},
		{
			with: "unknown capability",
			input: func(t *testing.T) []byte {
				session := validSessionStart(t)
				session[40] = 'U'
				return session
			},
			expect: func(t *testing.T, ssid oftp2.StartSessionCmd) {
				require.EqualError(t, ssid.Valid(), "unknown capability: U")
			},
		},
		{
			with: "unknown compression indicator",
			input: func(t *testing.T) []byte {
				session := validSessionStart(t)
				session[41] = 'U'
				return session
			},
			expect: func(t *testing.T, ssid oftp2.StartSessionCmd) {
				require.EqualError(t, ssid.Valid(), "unknown BufferCompressionIndicator: U")
			},
		},
		{
			with: "unknown restart indicator",
			input: func(t *testing.T) []byte {
				session := validSessionStart(t)
				session[42] = 'U'
				return session
			},
			expect: func(t *testing.T, ssid oftp2.StartSessionCmd) {
				require.EqualError(t, ssid.Valid(), "unknown RestartIndicator: U")
			},
		},
		{
			with: "unknown special logic indicator",
			input: func(t *testing.T) []byte {
				session := validSessionStart(t)
				session[43] = 'U'
				return session
			},
			expect: func(t *testing.T, ssid oftp2.StartSessionCmd) {
				require.EqualError(t, ssid.Valid(), "unknown SpecialLogicIndicator: U")
			},
		},
		{
			with: "invalid credit",
			input: func(t *testing.T) []byte {
				session := validSessionStart(t)
				session[45] = 'U'
				return session
			},
			expect: func(t *testing.T, ssid oftp2.StartSessionCmd) {
				require.EqualError(t, ssid.Valid(), `invalid Credit: strconv.Atoi: parsing "9U9": invalid syntax`)
			},
		},
		{
			with: "negative credit",
			input: func(t *testing.T) []byte {
				session := validSessionStart(t)
				session[44] = '-'
				return session
			},
			expect: func(t *testing.T, ssid oftp2.StartSessionCmd) {
				require.EqualError(t, ssid.Valid(), `invalid Credit: -99`)
			},
		},
		{
			with: "unknown authentication",
			input: func(t *testing.T) []byte {
				session := validSessionStart(t)
				session[47] = 'U'
				return session
			},
			expect: func(t *testing.T, ssid oftp2.StartSessionCmd) {
				require.EqualError(t, ssid.Valid(), `unknown Authentication: U`)
			},
		},
		{
			with: "missing CarriageReturn",
			input: func(t *testing.T) []byte {
				session := validSessionStart(t)
				session[60] = 'y'
				return session
			},
			expect: func(t *testing.T, ssid oftp2.StartSessionCmd) {
				require.EqualError(t, ssid.Valid(), `does not end on carriage return, but on y`)
			},
		},
		{
			with: "exceeding message",
			input: func(t *testing.T) []byte {
				session := validSessionStart(t)
				return append(session, ' ')
			},
			expect: func(t *testing.T, ssid oftp2.StartSessionCmd) {
				require.EqualError(t, ssid.Valid(), `invalid size: 62`)
			},
		},
	} {
		t.Run(scenario.with, func(t *testing.T) {
			scenario.expect(t, scenario.input(t))
		})
	}
}

func validSessionStart(t *testing.T) oftp2.Command {
	session, err := oftp2.NewStartSession(oftp2.StartSessionInput{
		IdentificationCode:     validSsidCode(t),
		Password:               "password",
		DataExchangeBufferSize: 99999,
		Capabilities:           oftp2.CapabilityBoth,
		BufferCompression:      true,
		Restart:                true,
		SpecialLogic:           true,
		Credit:                 999,
		SecureAuthentication:   true,
		UserData:               "        ",
	})
	require.NoError(t, err)
	return session
}
