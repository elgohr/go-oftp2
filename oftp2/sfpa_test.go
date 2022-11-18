package oftp2_test

import (
	"bifroest/oftp2"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestStartFilePositiveAnswer(t *testing.T) {
	for _, scenario := range []struct {
		with  string
		input int
		expect func(t *testing.T, cmd oftp2.Command, err error)
	}{
		{
			with:  "a standard input",
			input: 1,
			expect: func(t *testing.T, cmd oftp2.Command, err error) {
				require.NoError(t, err)
				require.NotNil(t, cmd)
			},
		},
		{
			with:  "a negative input",
			input: -1,
			expect: func(t *testing.T, cmd oftp2.Command, err error) {
				require.Error(t, err)
				require.Nil(t, cmd)
			},
		},
		{
			with:  "a exceeding input",
			input: 100000000000000000,
			expect: func(t *testing.T, cmd oftp2.Command, err error) {
				require.Error(t, err)
				require.Nil(t, cmd)
			},
		},
	} {
		t.Run(scenario.with, func(t *testing.T) {
			session, err := oftp2.NewStartFilePositiveAnswer(scenario.input)
			scenario.expect(t, session, err)
		})
	}
}

func TestStartFilePositiveAnswer_Valid(t *testing.T) {
	for _, scenario := range []struct {
		with  string
		input func(t *testing.T) []byte
		expect func(t *testing.T, sfpa oftp2.StartFilePositiveAnswerCmd)
	}{
		{
			with: "a standard message",
			input: func(t *testing.T) []byte {
				return validStartFilePositive(t)
			},
			expect: func(t *testing.T, sfpa oftp2.StartFilePositiveAnswerCmd) {
				require.NoError(t, sfpa.Valid())
				require.Equal(t, 1, sfpa.AnswerCount())
			},
		},
		{
			with: "a wrong cmd type",
			input: func(t *testing.T) []byte {
				p := validStartFilePositive(t)
				p[0] = '^'
				return p
			},
			expect: func(t *testing.T, sfpa oftp2.StartFilePositiveAnswerCmd) {
				require.EqualError(t, sfpa.Valid(), "does not start with 2, but with ^")
				require.Equal(t, 1, sfpa.AnswerCount())
			},
		},
		{
			with: "a wrong length",
			input: func(t *testing.T) []byte {
				p := validStartFilePositive(t)
				return append(p, ' ')
			},
			expect: func(t *testing.T, sfpa oftp2.StartFilePositiveAnswerCmd) {
				require.EqualError(t, sfpa.Valid(), "expected the length of 19, but got 20")
				require.Equal(t, 1, sfpa.AnswerCount())
			},
		},
		{
			with: "missing carriage return",
			input: func(t *testing.T) []byte {
				p := validStartFilePositive(t)
				p[len(p)-1] = 'd'
				return p
			},
			expect: func(t *testing.T, sfpa oftp2.StartFilePositiveAnswerCmd) {
				require.EqualError(t, sfpa.Valid(), "does not end on carriage return, but on d")
				require.Equal(t, 1, sfpa.AnswerCount())
			},
		},
		{
			with: "corrupted answer count",
			input: func(t *testing.T) []byte {
				p := validStartFilePositive(t)
				p[3] = 'd'
				return p
			},
			expect: func(t *testing.T, sfpa oftp2.StartFilePositiveAnswerCmd) {
				require.EqualError(t, sfpa.Valid(), `strconv.Atoi: parsing "00d00000000000001": invalid syntax`)
				require.Equal(t, 0, sfpa.AnswerCount())
			},
		},
		{
			with: "negative answer count",
			input: func(t *testing.T) []byte {
				p := validStartFilePositive(t)
				p[1] = '-'
				return p
			},
			expect: func(t *testing.T, sfpa oftp2.StartFilePositiveAnswerCmd) {
				require.EqualError(t, sfpa.Valid(), `answer count can't be negative`)
				require.Equal(t, -1, sfpa.AnswerCount())
			},
		},
	} {
		t.Run(scenario.with, func(t *testing.T) {
			scenario.expect(t, scenario.input(t))
		})
	}
}

func validStartFilePositive(t *testing.T) oftp2.Command {
	file, err := oftp2.NewStartFilePositiveAnswer(1)
	require.NoError(t, err)
	return file
}
