package oftp2

import (
	"errors"
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
// https://tools.ietf.org/html/rfc5024#section-5.3.5

type StartFileNegativeAnswerCmd []byte

func (c StartFileNegativeAnswerCmd) Valid() error {
	fixLength := 8 // prefix + CR
	variableLength, err := strconv.Atoi(string(c[4:7]))
	if err != nil {
		return err
	}
	totalLength := fixLength + variableLength
	if l := len(c); l != totalLength {
		return fmt.Errorf(InvalidLengthErrorFormat, totalLength, l)
	} else if Cmd(c[0]) != StartFileNegativeMessage {
		return fmt.Errorf(InvalidPrefixErrorFormat, string(StartFileNegativeMessage), string(c[0]))
	} else if c.ReasonCode() == 0 {
		return fmt.Errorf("invalid reason code")
	}else if c[3] != 'Y' && c[3] != 'N' {
		return fmt.Errorf("invalid retry")
	}else if string(c[totalLength-1]) != CarriageReturn {
		return fmt.Errorf(InvalidSuffixErrorFormat, string(c[totalLength-1]))
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

func (c StartFileNegativeAnswerCmd) ReasonText() []byte {
	return c[7:len(c)-1]
}

func NewStartFileNegativeAnswer(reason AnswerReason, retry bool, reasonText string) (Command, error) {
	if _, exists := KnownReasonCodes[reason]; !exists {
		return nil, fmt.Errorf("unknown answer reason: %d", reason)
	}
	i := len(reasonText)
	if i > 999 {
		return nil, errors.New("reason text is too long")
	}
	r, err := fillUpInt(int(reason), 2)
	if err != nil {
		return nil, err
	}

	l, err := fillUpInt(i, 3)
	if err != nil {
		return nil, err
	}
	return Command(
		string(StartFileNegativeMessage) +
			r +
			boolToString(retry) +
			l +
			reasonText +
			CarriageReturn,
	), nil
}

type AnswerReason int

var KnownReasonCodes = map[AnswerReason]struct{}{
	AnswerInvalidFilename: {},
	AnswerInvalidDestination: {},
	AnswerInvalidOrigin: {},
	AnswerStorageRecordFormatNotSupported: {},
	AnswerMaximumRecordLengthNotSupported: {},
	AnswerFilesizeTooBig: {},
	AnswerInvalidRecordCount: {},
	AnswerInvalidByteCount: {},
	AnswerAccessMethodFailure: {},
	AnswerDuplicateFile: {},
	AnswerFileDirectionRefused: {},
	AnswerCipherSuiteNotSupported: {},
	AnswerEncryptedFileNotAllowed: {},
	AnswerUnencryptedFileNotAllowed: {},
	AnswerCompressionNotAllowed: {},
	AnswerSignedFileNotAllowed: {},
	AnswerUnsignedFileNotAllowed: {},
	AnswerUnspecified: {},
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
