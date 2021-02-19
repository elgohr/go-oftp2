package oftp2

import (
	"fmt"
)

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

const (
	ssrmCmd = "I"
	ssrmMsg = "ODETTE FTP READY "
)

type StartSessionReadyMessageCmd []byte

func (c StartSessionReadyMessageCmd) Valid() error {
	if size := len(c); size != 19 {
		return fmt.Errorf(InvalidLengthErrorFormat, 19, size)
	} else if cmd := string(c[0]); cmd != "I" {
		return fmt.Errorf(InvalidPrefixErrorFormat, "I", cmd)
	} else if string(c[18]) != CarriageReturn {
		return fmt.Errorf(InvalidSuffixErrorFormat, c[18])
	}
	return nil
}

func (c StartSessionReadyMessageCmd) Message() []byte {
	return c[1:18]
}

func NewStartSessionReadyMessage() Command {
	return Command(ssrmCmd + ssrmMsg + CarriageReturn)
}
