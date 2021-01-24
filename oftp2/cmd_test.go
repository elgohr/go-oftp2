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
			cmd:            oftp2.Command("IODETTE FTP READY " + oftp2.CarriageReturn),
			expectedHeader: []byte{0x10, 0x00, 0x00, 0x17},
		},
		//{
		//	bufferLength:   41,
		//	expectedHeader: []byte{0x10, 0x00, 0x00, 0x41},
		//},
		//{
		//	bufferLength:  4,
		//	expectedError: errors.New("wrong buffer length (4)"),
		//},
		//{
		//	bufferLength:   1003,
		//	expectedHeader: []byte{0x10, 0x00, 0x10, 0x03},
		//},
		//{
		//	bufferLength:   100003,
		//	expectedHeader: []byte{0x10, 0x10, 0x00, 0x03},
		//},
		//{
		//	bufferLength:  100004,
		//	expectedError: errors.New("wrong buffer length (100004)"),
		//},
	} {
		t.Run(string(scenario.cmd), func(t *testing.T) {
			require.Equal(t, scenario.expectedHeader, scenario.cmd.StreamTransmissionBuffer()[:4])
		})
	}

}
