package oftp2

import (
	"fmt"
	"strconv"
)

// o-------------------------------------------------------------------o
// |       SFNA        Start File Negative Answer                      |
// |                                                                   |
// |       Start File Phase           Speaker <---- Listener           |
// |-------------------------------------------------------------------|
// | Pos | Field     | Description                           | Format  |
// |-----+-----------+---------------------------------------+---------|
// |   0 | SFNACMD   | SFNA Command, '3'                     | F X(1)  |
// |   1 | SFNAREAS  | Answer Reason                         | F 9(2)  |
// |   3 | SFNARRTR  | Retry Indicator, (Y/N)                | F X(1)  |
// |   4 | SFNAREASL | Answer Reason Text Length             | V 9(3)  |
// |   7 | SFNAREAST | Answer Reason Text                    | V T(n)  |
// o-------------------------------------------------------------------o
//
// https://datatracker.ietf.org/doc/html/rfc5024#section-5.3.5

type StartFileNegativeAnswerCmd []byte

func (c StartFileNegativeAnswerCmd) Valid() error {
	fixLength := 8 // prefix + CR
	variableLength, err := strconv.Atoi(string(c[4:7]))
	if err != nil {
		return err
	}
	totalLength := fixLength + variableLength
	if length := len(c); length != totalLength {
		return NewInvalidLengthError(totalLength, length)
	} else if StartFileNegativeMessage.Byte() != c[0] {
		return NewInvalidPrefixError(StartFileNegativeMessage.String(), string(c[0]))
	} else if 0 == c.ReasonCode() {
		return fmt.Errorf("invalid reason code")
	} else if c[3] != 'Y' && c[3] != 'N' {
		return fmt.Errorf("invalid retry")
	} else if cmd := string(c[totalLength-1]); CarriageReturn != cmd {
		return NewNoCrSuffixError(cmd)
	}
	return nil
}

func (c StartFileNegativeAnswerCmd) ReasonCode() AnswerReason {
	i, _ := strconv.Atoi(string(c[1:3]))
	return AnswerReason(i)
}

func (c StartFileNegativeAnswerCmd) Retry() bool {
	return c[3] == 'Y'
}

func (c StartFileNegativeAnswerCmd) ReasonText() string {
	return string(c[7 : len(c)-1])
}

func NewStartFileNegativeAnswer(input NegativeFileInput) (Command, error) {
	if _, exists := KnownReasonCodes[input.Reason]; !exists {
		return nil, fmt.Errorf("unknown answer reason: %d", input.Reason)
	}
	length := len(input.ReasonText)
	if length > 999 {
		return nil, fmt.Errorf("reason text is too long: %d", length)
	}
	r, _ := fillUpInt(int(input.Reason), 2)
	l, _ := fillUpInt(length, 3)

	return Command(
		string(StartFileNegativeMessage) +
			r +
			boolToString(input.Retry) +
			l +
			input.ReasonText +
			CarriageReturn,
	), nil
}

type NegativeFileInput struct {
	Reason     AnswerReason
	Retry      bool
	ReasonText string
}

type AnswerReason int

var KnownReasonCodes = map[AnswerReason]struct{}{
	AnswerInvalidFilename:                 {},
	AnswerInvalidDestination:              {},
	AnswerInvalidOrigin:                   {},
	AnswerStorageRecordFormatNotSupported: {},
	AnswerMaximumRecordLengthNotSupported: {},
	AnswerFilesizeTooBig:                  {},
	AnswerInvalidRecordCount:              {},
	AnswerInvalidByteCount:                {},
	AnswerAccessMethodFailure:             {},
	AnswerDuplicateFile:                   {},
	AnswerFileDirectionRefused:            {},
	AnswerCipherSuiteNotSupported:         {},
	AnswerEncryptedFileNotAllowed:         {},
	AnswerUnencryptedFileNotAllowed:       {},
	AnswerCompressionNotAllowed:           {},
	AnswerSignedFileNotAllowed:            {},
	AnswerUnsignedFileNotAllowed:          {},
	AnswerUnspecified:                     {},
}

const (
	AnswerInvalidFilename                 AnswerReason = 01
	AnswerInvalidDestination              AnswerReason = 02
	AnswerInvalidOrigin                   AnswerReason = 03
	AnswerStorageRecordFormatNotSupported AnswerReason = 04
	AnswerMaximumRecordLengthNotSupported AnswerReason = 05
	AnswerFilesizeTooBig                  AnswerReason = 06
	AnswerInvalidRecordCount              AnswerReason = 10
	AnswerInvalidByteCount                AnswerReason = 11
	AnswerAccessMethodFailure             AnswerReason = 12
	AnswerDuplicateFile                   AnswerReason = 13
	AnswerFileDirectionRefused            AnswerReason = 14
	AnswerCipherSuiteNotSupported         AnswerReason = 15
	AnswerEncryptedFileNotAllowed         AnswerReason = 16
	AnswerUnencryptedFileNotAllowed       AnswerReason = 17
	AnswerCompressionNotAllowed           AnswerReason = 18
	AnswerSignedFileNotAllowed            AnswerReason = 19
	AnswerUnsignedFileNotAllowed          AnswerReason = 20
	AnswerUnspecified                     AnswerReason = 99
)
