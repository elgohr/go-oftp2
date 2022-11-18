package oftp2_test

import (
	"bifroest/oftp2"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestSid(t *testing.T) {
	for _, scenario := range []struct {
		with   string
		input  func() oftp2.SidInput
		expect func(t *testing.T, sid oftp2.Sid, err error)
	}{
		{
			with:  "a standard input",
			input: validSidInput,
			expect: func(t *testing.T, sid oftp2.Sid, err error) {
				require.NoError(t, err)
				require.NotNil(t, sid)
			},
		},
		{
			with: "an exceeding code designator",
			input: func() oftp2.SidInput {
				i := validSidInput()
				i.CodeDesignator = "10000"
				return i
			},
			expect: func(t *testing.T, sid oftp2.Sid, err error) {
				require.EqualError(t, err, "code designator is too long: 10000")
				require.Nil(t, sid)
			},
		},
		{
			with: "an invalid organisation code",
			input: func() oftp2.SidInput {
				i := validSidInput()
				i.OrganisationCode = "!123456789101"
				return i
			},
			expect: func(t *testing.T, sid oftp2.Sid, err error) {
				require.EqualError(t, err, "organisation code is may contain ^[a-zA-Z0-9- ]+$")
				require.Nil(t, sid)
			},
		},
		{
			with: "an exceeding organisation code",
			input: func() oftp2.SidInput {
				i := validSidInput()
				i.OrganisationCode = "123456789101112"
				return i
			},
			expect: func(t *testing.T, sid oftp2.Sid, err error) {
				require.EqualError(t, err, "organisation code is too long: 123456789101112")
				require.Nil(t, sid)
			},
		},
		{
			with: "an exceeding subaddress",
			input: func() oftp2.SidInput {
				i := validSidInput()
				i.SubAddress = "1234567"
				return i
			},
			expect: func(t *testing.T, sid oftp2.Sid, err error) {
				require.EqualError(t, err, "subaddress is too long: 1234567")
				require.Nil(t, sid)
			},
		},
	} {
		t.Run(scenario.with, func(t *testing.T) {
			s, err := oftp2.NewSid(scenario.input())
			scenario.expect(t, s, err)
		})
	}
}

func TestSid_Valid(t *testing.T) {
	for _, scenario := range []struct {
		with   string
		input  func(t *testing.T) []byte
		expect func(t *testing.T, sid oftp2.Sid)
	}{
		{
			with: "a standard message",
			input: func(t *testing.T) []byte {
				return validSid(t)
			},
			expect: func(t *testing.T, sid oftp2.Sid) {
				require.NoError(t, sid.Valid())
				require.Equal(t, "Test", sid.CodeDesignator())
				require.Equal(t, "Org", sid.OrganisationCode())
				require.Equal(t, "Addres", sid.SubAddress())
			},
		},
		{
			with: "a wrong cmd type",
			input: func(t *testing.T) []byte {
				p := validSid(t)
				p[0] = '^'
				return p
			},
			expect: func(t *testing.T, sid oftp2.Sid) {
				require.EqualError(t, sid.Valid(), "does not start with O, but with ^")
				require.Equal(t, "Test", sid.CodeDesignator())
				require.Equal(t, "Org", sid.OrganisationCode())
				require.Equal(t, "Addres", sid.SubAddress())
			},
		},
		{
			with: "a wrong length",
			input: func(t *testing.T) []byte {
				p := validSid(t)
				return append(p, ' ')
			},
			expect: func(t *testing.T, sid oftp2.Sid) {
				require.EqualError(t, sid.Valid(), "expected the length of 25, but got 26")
				require.Equal(t, "Test", sid.CodeDesignator())
				require.Equal(t, "Org", sid.OrganisationCode())
				require.Equal(t, "Addres", sid.SubAddress())
			},
		},
		{
			with: "corrupted organisation code",
			input: func(t *testing.T) []byte {
				p := validSid(t)
				p[6] = '!'
				return p
			},
			expect: func(t *testing.T, sid oftp2.Sid) {
				require.EqualError(t, sid.Valid(), `organisation code is may contain ^[a-zA-Z0-9- ]+$`)
				require.Equal(t, "Test", sid.CodeDesignator())
				require.Equal(t, "!         Org", sid.OrganisationCode())
				require.Equal(t, "Addres", sid.SubAddress())
			},
		},
	} {
		t.Run(scenario.with, func(t *testing.T) {
			scenario.expect(t, scenario.input(t))
		})
	}
}

func validSidInput() oftp2.SidInput {
	return oftp2.SidInput{
		CodeDesignator:   "9999",
		OrganisationCode: "Org",
		SubAddress:       "Addres",
	}
}

func validSid(t *testing.T) oftp2.Sid {
	file, err := oftp2.NewSid(oftp2.SidInput{
		CodeDesignator:   "Test",
		OrganisationCode: "Org",
		SubAddress:       "Addres",
	})
	require.NoError(t, err)
	return file
}
