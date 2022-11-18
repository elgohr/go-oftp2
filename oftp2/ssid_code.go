package oftp2

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
	if size := len(c); size != 25 {
		return NewInvalidLengthError(25, size)
	}
	return nil
}

func (c IdentificationCode) OdetteIdentifier() byte {
	return c[0]
}

func (c IdentificationCode) InternationalCodeDesignator() []byte {
	return c[1:5]
}

func (c IdentificationCode) OrganisationCode() []byte {
	return c[5:19]
}

func (c IdentificationCode) ComputerSubaddress() []byte {
	return c[19:25]
}

type SsidIdentificationCodeInput struct {
	OdetteIdentifier            string
	InternationalCodeDesignator string
	OrganisationCode            string
	ComputerSubaddress          string
}

func SsidIdentificationCode(input SsidIdentificationCodeInput) (IdentificationCode, error) {
	oid, err := fillUpString(input.OdetteIdentifier, 1)
	if err != nil {
		return nil, err
	}
	intCode, err := fillUpString(input.InternationalCodeDesignator, 4)
	if err != nil {
		return nil, err
	}
	orgCode, err := fillUpString(input.OrganisationCode, 14)
	if err != nil {
		return nil, err
	}
	subAddr, err := fillUpString(input.ComputerSubaddress, 6)
	if err != nil {
		return nil, err
	}
	return IdentificationCode(
		oid +
			intCode +
			orgCode +
			subAddr,
	), nil
}
