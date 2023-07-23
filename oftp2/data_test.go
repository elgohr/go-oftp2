package oftp2_test

import (
	"github.com/elgohr/go-oftp2/oftp2"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestDataExchangeBuffer(t *testing.T) {
	inputPayload := []byte("INPUT")
	d := oftp2.NewDataExchangeBuffer(inputPayload)
	exchangeBuffer := oftp2.DataExchangeBuffer(d)
	require.NoError(t, exchangeBuffer.Valid())
	require.Equal(t, inputPayload, exchangeBuffer.Payload())
}

func TestDataExchangeBuffer_Invalid(t *testing.T) {
	d := oftp2.NewDataExchangeBuffer([]byte("INPUT"))
	d[0] = 'F'
	exchangeBuffer := oftp2.DataExchangeBuffer(d)
	require.EqualError(t, exchangeBuffer.Valid(), "does not start with D, but with F")
}
