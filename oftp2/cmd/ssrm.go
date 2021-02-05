package oftp2

import (
	"bifroest/oftp2"
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
		return fmt.Errorf(oftp2.InvalidLengthErrorFormat, 19, size)
	} else if cmd := string(c[0]); cmd != "I" {
		return fmt.Errorf(oftp2.InvalidPrefixErrorFormat, "I", cmd)
	} else if string(c[18]) != oftp2.CarriageReturn {
		return fmt.Errorf(oftp2.InvalidSuffixErrorFormat, c[18])
	}
	return nil
}

func (c StartSessionReadyMessageCmd) Command() byte {
	return c[0]
}

func (c StartSessionReadyMessageCmd) Message() []byte {
	return c[1:18]
}

func StartSessionReadyMessage() oftp2.Command {
	return oftp2.Command(ssrmCmd + ssrmMsg + oftp2.CarriageReturn)
}
