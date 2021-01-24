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
	if len(c) != 19 {
		return fmt.Errorf(oftp2.InvalidLengthErrorFormat, 19, len(c))
	} else if string(c[0]) != "I" {
		return fmt.Errorf(oftp2.InvalidPrefixErrorFormat, "I", c[0])
	} else if string(c[18]) != oftp2.CarriageReturn {
		return fmt.Errorf(oftp2.InvalidSuffixErrorFormat, c[18])
	}
	return nil
}

func (c StartSessionReadyMessageCmd) Cmd() byte {
	return c[0]
}

func (c StartSessionReadyMessageCmd) Msg() []byte {
	return c[1:18]
}

func StartSessionReadyMessage() oftp2.Command {
	return oftp2.Command(ssrmCmd + ssrmMsg + oftp2.CarriageReturn)
}
