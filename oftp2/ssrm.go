package oftp2

// o-------------------------------------------------------------------o
// |       SSRM        Start Session Ready Message                     |
// |                                                                   |
// |       Start Session Phase     Initiator <---- Responder           |
// |-------------------------------------------------------------------|
// | Pos | Field     | Description                           | Format  |
// |-----+-----------+---------------------------------------+---------|
// |   0 | SSRMCMD   | SSRM Command, 'I'                     | F X(1)  |
// |   1 | SSRMMSG   | Ready Message, 'ODETTE FTP READY '    | F X(17) |
// |  18 | SSRMCR    | Carriage Return                       | F X(1)  |
// o-------------------------------------------------------------------o
//
// https://tools.ietf.org/html/rfc5024#section-5.3.1

type StartSessionReadyMessageCmd []byte

func (c StartSessionReadyMessageCmd) Valid() error {
	if size := len(c); size != 19 {
		return NewInvalidLengthError(19, size)
	} else if StartSessionReadyMessage.Byte() != c[0] {
		return NewInvalidPrefixError(StartSessionReadyMessage.String(), string(c[0]))
	} else if cmd := string(c[18]); cmd != CarriageReturn {
		return NewNoCrSuffixError(cmd)
	}
	return nil
}

func (c StartSessionReadyMessageCmd) Message() []byte {
	return c[1:18]
}

func NewStartSessionReadyMessage() Command {
	return Command(string(StartSessionReadyMessage) + "ODETTE FTP READY " + CarriageReturn)
}
