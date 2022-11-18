package oftp2

import (
	"errors"
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
	if size := len(c); size != 61 {
		return fmt.Errorf("invalid size: %d", size)
	} else if StartSessionMessage.Byte() != c[0] {
		return NewInvalidPrefixError(StartSessionMessage.String(), string(c[0]))
	} else if cmd := string(c[60]); CarriageReturn != cmd {
		return NewNoCrSuffixError(cmd)
	} else if err := c.IdentificationCode().Valid(); err != nil {
		return err
	} else if level := c.ProtocolLevel(); level != '5' {
		return fmt.Errorf("invalid protocol level: %d", level)
	} else if de, err := strconv.Atoi(string(c[35:40])); err != nil {
		return fmt.Errorf("invalid DataExchangeBufferSize: %w", err)
	} else if de < 128 || de > 99999 {
		return fmt.Errorf("invalid DataExchangeBufferSize: %d", de)
	} else if ca := c.Capabilities(); !isCapability(ca) {
		return fmt.Errorf("unknown capability: %s", ca)
	} else if bc := string(c[41]); !isBool(bc) {
		return fmt.Errorf("unknown BufferCompressionIndicator: %s", bc)
	} else if ri := string(c[42]); !isBool(ri) {
		return fmt.Errorf("unknown RestartIndicator: %s", ri)
	} else if sli := string(c[43]); !isBool(sli) {
		return fmt.Errorf("unknown SpecialLogicIndicator: %s", sli)
	} else if cred, err := strconv.Atoi(string(c[44:47])); err != nil {
		return fmt.Errorf("invalid Credit: %w", err)
	} else if cred < 0 || cred > 999 {
		return fmt.Errorf("invalid Credit: %d", cred)
	} else if auth := string(c[47]); !isBool(auth) {
		return fmt.Errorf("unknown Authentication: %s", auth)
	}

	return nil
}

func (c StartSessionCmd) ProtocolLevel() byte {
	return c[1]
}

func (c StartSessionCmd) IdentificationCode() IdentificationCode {
	return IdentificationCode(c[2:27])
}

func (c StartSessionCmd) Password() []byte {
	return c[27:35]
}

func (c StartSessionCmd) DataExchangeBufferSize() int {
	i, _ := strconv.Atoi(string(c[35:40]))
	return i
}

func (c StartSessionCmd) Capabilities() SsidCapability {
	return SsidCapability(c[40])
}

func (c StartSessionCmd) BufferCompression() bool {
	return c[41] == 'Y'
}

func (c StartSessionCmd) Restart() bool {
	return c[42] == 'Y'
}

func (c StartSessionCmd) SpecialLogic() bool {
	return c[43] == 'Y'
}

func (c StartSessionCmd) Credit() int {
	i, _ := strconv.Atoi(string(c[44:47]))
	return i
}

func (c StartSessionCmd) Authentication() bool {
	return c[47] == 'Y'
}

func (c StartSessionCmd) User() []byte {
	return c[48:56]
}

type SsidCapability string

const (
	CapabilitySend    SsidCapability = "S"
	CapabilityReceive SsidCapability = "R"
	CapabilityBoth    SsidCapability = "B"
)

type StartSessionInput struct {
	IdentificationCode     IdentificationCode
	Password               string
	DataExchangeBufferSize int
	Capabilities           SsidCapability
	BufferCompression      bool
	Restart                bool
	SpecialLogic           bool
	Credit                 int
	SecureAuthentication   bool
	UserData               string
}

func NewStartSession(input StartSessionInput) (Command, error) {
	if input.IdentificationCode == nil {
		return nil, errors.New("missing identification code")
	}

	if err := input.IdentificationCode.Valid(); err != nil {
		return nil, err
	}

	password, err := fillUpString(input.Password, 8)
	if err != nil {
		return nil, err
	}

	bufferSize, err := fillUpInt(input.DataExchangeBufferSize, 5)
	if err != nil {
		return nil, err
	}

	if !isCapability(input.Capabilities) {
		return nil, fmt.Errorf("unknown capability: %s", input.Capabilities)
	}

	credit, err := fillUpInt(input.Credit, 3)
	if err != nil {
		return nil, err
	}

	userData, err := fillUpString(input.UserData, 8)
	if err != nil {
		return nil, err
	}

	return Command(string(StartSessionMessage) +
		"5" + // OFTP-2
		string(input.IdentificationCode) +
		password +
		bufferSize +
		string(input.Capabilities) +
		boolToString(input.BufferCompression) +
		boolToString(input.Restart) +
		boolToString(input.SpecialLogic) +
		credit +
		boolToString(input.SecureAuthentication) +
		reserved(4) +
		userData +
		CarriageReturn), nil
}

func isCapability(input SsidCapability) bool {
	return input == CapabilitySend ||
		input == CapabilityReceive ||
		input == CapabilityBoth
}

func reserved(c int) string {
	res := ""
	for i := 0; i < c; i++ {
		res += " "
	}
	return res
}
