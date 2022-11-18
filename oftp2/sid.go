package oftp2

import (
	"fmt"
	"regexp"
	"strings"
)

type Sid []byte

func (c Sid) Valid() error {
	if length := len(c); length != 25 {
		return fmt.Errorf("expected the length of %d, but got %d", 25, length)
	} else if Id(c[0]) != SidId {
		return fmt.Errorf("does not start with %v, but with %v", string(SidId), string(c[0]))
	} else if err := isValidOrganisationCode(string(c[5:19])); err != nil {
		return err
	}
	return nil
}

func (c Sid) CodeDesignator() string {
	return string(c[1:5])
}

func (c Sid) OrganisationCode() string {
	return strings.TrimSpace(string(c[5:19]))
}

func (c Sid) SubAddress() string {
	return strings.TrimSpace(string(c[19:25]))
}

func NewSid(input SidInput) (Sid, error) {
	if len(input.CodeDesignator) > 4 {
		return nil, fmt.Errorf("code designator is too long: %v", input.CodeDesignator)
	} else if len(input.OrganisationCode) > 14 {
		return nil, fmt.Errorf("organisation code is too long: %v", input.OrganisationCode)
	} else if len(input.SubAddress) > 6 {
		return nil, fmt.Errorf("subaddress is too long: %v", input.SubAddress)
	} else if err := isValidOrganisationCode(input.OrganisationCode); err != nil {
		return nil, err
	}

	designator, err := fillUpString(input.CodeDesignator, 4)
	if err != nil {
		return nil, err
	}
	orgCode, err := fillUpString(input.OrganisationCode, 14)
	if err != nil {
		return nil, err
	}
	subAddr, err := fillUpString(input.SubAddress, 6)
	if err != nil {
		return nil, err
	}
	return Sid(
		"O" +
			designator +
			orgCode +
			subAddr,
	), nil
}

const organisationCodePattern = "^[a-zA-Z0-9- ]+$"

func isValidOrganisationCode(rawOrgCode string) error {
	if matched, err := regexp.Match(organisationCodePattern, []byte(rawOrgCode)); err != nil {
		return err
	} else if !matched {
		return fmt.Errorf("organisation code is may contain %v", organisationCodePattern)
	}
	return nil
}

type SidInput struct {
	CodeDesignator   string
	OrganisationCode string
	SubAddress       string
}
