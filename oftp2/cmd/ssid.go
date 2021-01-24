package oftp2

import "bifroest/oftp2"

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

func (c StartSessionCmd) Cmd() byte {
	return c[0]
}

func (c StartSessionCmd) Lev() byte {
	return c[1]
}

func StartSession(identification IdentificationCode) oftp2.Command {
	return oftp2.Command(ssidCmd + ssidLev + identification + oftp2.CarriageReturn)
}

const (
	ssidCmd = 'X'
	ssidLev = '5'
)
