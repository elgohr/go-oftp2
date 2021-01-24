package oftp2

import (
	"bifroest/oftp2"
	"fmt"
	"strconv"
)

// o-------------------------------------------------------------------o
// | Pos | Field     | Description                           | Format  |
// |-----+-----------+---------------------------------------+---------|
// |   0 | SIOOID    | ODETTE Identifier                     | F X(1)  |
// |   1 | SIOICD    | International Code Designator         | V 9(4)  |
// |   5 | SIOORG    | Organisation Code                     | V X(14) |
// |  19 | SIOCSA    | Computer Subaddress                   | V X(6)  |
// o-------------------------------------------------------------------o
//
// https://tools.ietf.org/html/rfc5024#section-5.4

type IdentificationCode []byte

func (c IdentificationCode) Valid() error {
	if len(c) != 25 {
		return fmt.Errorf(oftp2.InvalidLengthErrorFormat, 25, len(c))
	} else if _, err := strconv.ParseInt(string(c[1:5]), 10, 32); err != nil {
		return fmt.Errorf("international code designator is not a number, but %v", string(c[1:3]))
	}
	return nil
}

func (c IdentificationCode) Oid() byte {
	return c[0]
}

func (c IdentificationCode) Icd() int {
	i , _ := strconv.ParseInt(string(c[1:5]), 10, 32)
	return int(i)
}

func (c IdentificationCode) Org() []byte {
	return c[5:19]
}

func (c IdentificationCode) Csa() []byte {
	return c[19:25]
}

func SsidIdentificationCode(siooid string, sioicd int, sioorg string, siocsa string) IdentificationCode {
	return IdentificationCode(siooid + strconv.Itoa(sioicd) + sioorg + siocsa)
}
