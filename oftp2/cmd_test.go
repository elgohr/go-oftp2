package oftp2_test

import (
	"bifroest/oftp2"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestNewStreamTransmissionHeader(t *testing.T) {
	for _, scenario := range []struct {
		cmd            oftp2.Command
		expectedHeader []byte
	}{
		{
			cmd:            oftp2.NewStartSessionReadyMessage(),
			expectedHeader: []byte{0x10, 0x00, 0x00, 0x17},
		},
	} {
		t.Run(string(scenario.cmd), func(t *testing.T) {
			require.Equal(t, scenario.expectedHeader, scenario.cmd.StreamTransmissionBuffer()[:4])
		})
	}
}

func TestCommandCmd(t *testing.T) {
	for _, scenario := range []struct {
		cmd         func(t *testing.T) oftp2.Command
		expectedCmd oftp2.Cmd
	}{
		{
			cmd: func(t *testing.T) oftp2.Command {
				return oftp2.NewStartSessionReadyMessage()
			},
			expectedCmd: oftp2.StartSessionReadyMessage,
		},
		{
			cmd: func(t *testing.T) oftp2.Command {
				return validSessionStart(t)
			},
			expectedCmd: oftp2.StartSessionMessage,
		},
		{
			cmd: func(t *testing.T) oftp2.Command {
				return []byte("-")
			},
			expectedCmd: oftp2.Unknown,
		},
	} {
		t.Run(string(scenario.expectedCmd), func(t *testing.T) {
			cmd := scenario.cmd(t).Cmd()
			require.Equal(t, scenario.expectedCmd, cmd, string(cmd))
		})
	}
}
