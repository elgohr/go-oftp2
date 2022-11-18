package oftp2

// o-------------------------------------------------------------------o
// |       DATA        Data Exchange Buffer                            |
// |                                                                   |
// |       Data Transfer Phase        Speaker ----> Listener           |
// |-------------------------------------------------------------------|
// | Pos | Field     | Description                           | Format  |
// |-----+-----------+---------------------------------------+---------|
// |   0 | DATACMD   | DATA Command, 'D'                     | F X(1)  |
// |   1 | DATABUF   | Data Exchange Buffer payload          | V U(n)  |
// o-------------------------------------------------------------------o
//
// https://www.rfc-editor.org/rfc/rfc5024#section-5.3.6

type DataExchangeBuffer []byte

func NewDataExchangeBuffer(data []byte) Command {
	return append([]byte{DataExchangeBufferMessage.Byte()}, data...)
}

func (c DataExchangeBuffer) Valid() error {
	if DataExchangeBufferMessage.Byte() != c[0] {
		return NewInvalidPrefixError(DataExchangeBufferMessage.String(), string(c[0]))
	}
	return nil
}

func (c DataExchangeBuffer) Payload() []byte {
	return c[1:]
}
