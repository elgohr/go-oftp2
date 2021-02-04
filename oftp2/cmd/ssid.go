package oftp2

import (
	"bifroest/oftp2"
	"fmt"
	"strconv"
)

// o-------------------------------------------------------------------o
// |       SSID        Start Session                                   |
// |                                                                   |
// |       Start Session Phase     Initiator <---> Responder           |
// |-------------------------------------------------------------------|
// | Pos | Field     | Description                           | Format  |
// |-----+-----------+---------------------------------------+---------|
// |   0 | SSIDCMD   | SSID Command 'X'                      | F X(1)  |
// |   1 | SSIDLEV   | Protocol Release Level                | F 9(1)  |
// |   2 | SSIDCODE  | Initiator's Identification Code       | V X(25) |
// |  27 | SSIDPSWD  | Initiator's Password                  | V X(8)  |
// |  35 | SSIDSDEB  | Data Exchange Buffer Size             | V 9(5)  |
// |  40 | SSIDSR    | Send / Receive Capabilities (S/R/B)   | F X(1)  |
// |  41 | SSIDCMPR  | Buffer Compression Indicator (Y/N)    | F X(1)  |
// |  42 | SSIDREST  | Restart Indicator (Y/N)               | F X(1)  |
// |  43 | SSIDSPEC  | Special Logic Indicator (Y/N)         | F X(1)  |
// |  44 | SSIDCRED  | Credit                                | V 9(3)  |
// |  47 | SSIDAUTH  | Secure Authentication (Y/N)           | F X(1)  |
// |  48 | SSIDRSV1  | Reserved                              | F X(4)  |
// |  52 | SSIDUSER  | User Data                             | V X(8)  |
// |  60 | SSIDCR    | Carriage Return                       | F X(1)  |
// o-------------------------------------------------------------------o
//
// https://tools.ietf.org/html/rfc5024#section-5.3.2

type StartSessionCmd []byte

func (c StartSessionCmd) Valid() error {
	if string(c[0]) != "X" {
		return fmt.Errorf(oftp2.InvalidPrefixErrorFormat, "X", c[0])
	} else if err := c.Code().Valid(); err != nil {
		return err
	}
	return nil
}

func (c StartSessionCmd) Cmd() byte {
	return c[0]
}

func (c StartSessionCmd) Lev() byte {
	return c[1]
}

func (c StartSessionCmd) Code() IdentificationCode {
	return IdentificationCode(c[2:27])
}

func (c StartSessionCmd) Pswd() []byte {
	return c[27:35]
}

func (c StartSessionCmd) Sdeb() int {
	i, _ := strconv.Atoi(string(c[35:40]))
	return i
}

func (c StartSessionCmd) Dsr() byte {
	return c[40]
}

func (c StartSessionCmd) Cmpr() bool {
	return c[41] == 'Y'
}

func (c StartSessionCmd) Rest() bool {
	return c[42] == 'Y'
}

func (c StartSessionCmd) Spec() bool {
	return c[43] == 'Y'
}

func (c StartSessionCmd) Cred() int {
	i, _ := strconv.Atoi(string(c[44:47]))
	return i
}

func (c StartSessionCmd) Auth() bool {
	return c[47] == 'Y'
}

func (c StartSessionCmd) User() []byte {
	return c[48:56]
}

func StartSession(identification IdentificationCode, password string) oftp2.Command {
	id := string(identification)
	return oftp2.Command(ssidCmd +
		ssidLev +
		id +
		password +
		ssidDeb +
		ssidDsr +
		ssidCmpr +
		ssidRest +
		ssidSpec +
		ssidCred +
		ssidAuth +
		ssidUser +
		oftp2.CarriageReturn)
}

const (
	ssidCmd  = "X"
	ssidLev  = "5"
	ssidDeb  = "99999"
	ssidDsr  = "B"
	ssidCmpr = "Y"
	ssidRest = "Y"
	ssidSpec = "Y"
	ssidCred = "999"
	ssidAuth = "Y"
	ssidUser = "        "
)
