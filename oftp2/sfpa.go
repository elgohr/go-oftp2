package oftp2

import (
	"errors"
	"fmt"
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
		return fmt.Errorf(InvalidLengthErrorFormat, 19, l)
	} else if Cmd(c[0]) != StartFilePositiveMessage {
		return fmt.Errorf(InvalidPrefixErrorFormat, string(StartFilePositiveMessage), string(c[0]))
	} else if string(c[18]) != CarriageReturn {
		return fmt.Errorf(InvalidSuffixErrorFormat, string(c[18]))
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
