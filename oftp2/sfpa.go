package oftp2

import (
	"errors"
	"strconv"
)

// o-------------------------------------------------------------------o
// |       SFPA        Start File Positive Answer                      |
// |                                                                   |
// |       Start File Phase           Speaker <---- Listener           |
// |-------------------------------------------------------------------|
// | Pos | Field     | Description                           | Format  |
// |-----+-----------+---------------------------------------+---------|
// |   0 | SFPACMD   | SFPA Command, '2'                     | F X(1)  |
// |   1 | SFPAACNT  | Answer Count                          | V 9(17) |
// o-------------------------------------------------------------------o
//
// https://tools.ietf.org/html/rfc5024#section-5.3.4

type StartFilePositiveAnswerCmd []byte

func (c StartFilePositiveAnswerCmd) Valid() error {
	if l := len(c); l != 19 {
		return NewInvalidLengthError(19, l)
	} else if StartFilePositiveMessage.Byte() != c[0] {
		return NewInvalidPrefixError(StartFilePositiveMessage.String(), string(c[0]))
	} else if cmd := string(c[18]); CarriageReturn != cmd {
		return NewNoCrSuffixError(cmd)
	} else if val, err := strconv.Atoi(string(c[1:18])); err != nil {
		return err
	} else if val < 0 {
		return errors.New("answer count can't be negative")
	}
	return nil
}

func (c StartFilePositiveAnswerCmd) AnswerCount() int {
	i, _ := strconv.Atoi(string(c[1:18]))
	return i
}

func NewStartFilePositiveAnswer(count int) (Command, error) {
	if count < 0 {
		return nil, errors.New("answer count can't be negative")
	}
	c, err := fillUpInt(count, 17)
	if err != nil {
		return nil, err
	}
	return Command(
		string(StartFilePositiveMessage) +
			c +
			CarriageReturn), nil
}
