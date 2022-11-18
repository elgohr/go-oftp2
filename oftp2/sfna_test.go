package oftp2_test

import (
	"bifroest/oftp2"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestStartFileNegativeAnswer(t *testing.T) {
	for _, scenario := range []struct {
		with  string
		input oftp2.NegativeFileInput
		expect func(t *testing.T, cmd oftp2.Command, err error)
	}{
		{
			with: "a standard input",
			input: oftp2.NegativeFileInput{
				Reason:     oftp2.AnswerInvalidFilename,
				Retry:      true,
				ReasonText: "BECAUSE",
			},
			expect: func(t *testing.T, cmd oftp2.Command, err error) {
				require.NoError(t, err)
				require.NotNil(t, cmd)
			},
		},
		{
			with: "an unknown reasonCode",
			input: oftp2.NegativeFileInput{
				Reason:     98,
				Retry:      true,
				ReasonText: "",
			},
			expect: func(t *testing.T, cmd oftp2.Command, err error) {
				require.EqualError(t, err, "unknown answer reason: 98")
				require.Nil(t, cmd)
			},
		},
		{
			with: "a reason text that is too long",
			input: oftp2.NegativeFileInput{
				Reason:     oftp2.AnswerInvalidFilename,
				Retry:      false,
				ReasonText: generateLongString(1000),
			},
			expect: func(t *testing.T, cmd oftp2.Command, err error) {
				require.EqualError(t, err, "reason text is too long: 1000")
				require.Nil(t, cmd)
			},
		},
	} {
		t.Run(scenario.with, func(t *testing.T) {
			s, err := oftp2.NewStartFileNegativeAnswer(scenario.input)
			scenario.expect(t, s, err)
		})
	}
}

func TestStartFileNegativeAnswer_Valid(t *testing.T) {
	for _, scenario := range []struct {
		with  string
		input func(t *testing.T) []byte
		expect func(t *testing.T, sfna oftp2.StartFileNegativeAnswerCmd)
	}{
		{
			with: "a standard message",
			input: func(t *testing.T) []byte {
				file, err := oftp2.NewStartFileNegativeAnswer(oftp2.NegativeFileInput{
					Reason:     oftp2.AnswerInvalidFilename,
					Retry:      false,
					ReasonText: "MY_TEXT",
				})
				require.NoError(t, err)
				return file
			},
			expect: func(t *testing.T, sfna oftp2.StartFileNegativeAnswerCmd) {
				require.NoError(t, sfna.Valid())
				require.Equal(t, oftp2.AnswerInvalidFilename, sfna.ReasonCode())
				require.Equal(t, false, sfna.Retry())
				require.Equal(t, "MY_TEXT", sfna.ReasonText())
			},
		},
		{
			with: "a wrong cmd type",
			input: func(t *testing.T) []byte {
				p := validStartFileNegative(t)
				p[0] = '^'
				return p
			},
			expect: func(t *testing.T, sfna oftp2.StartFileNegativeAnswerCmd) {
				require.EqualError(t, sfna.Valid(), "does not start with 3, but with ^")
				require.Equal(t, oftp2.AnswerInvalidFilename, sfna.ReasonCode())
				require.Equal(t, true, sfna.Retry())
				require.Equal(t, "", sfna.ReasonText())
			},
		},
		{
			with: "a wrong length",
			input: func(t *testing.T) []byte {
				p := validStartFileNegative(t)
				return append(p, ' ')
			},
			expect: func(t *testing.T, sfna oftp2.StartFileNegativeAnswerCmd) {
				require.EqualError(t, sfna.Valid(), "expected the length of 8, but got 9")
				require.Equal(t, oftp2.AnswerInvalidFilename, sfna.ReasonCode())
				require.Equal(t, true, sfna.Retry())
				require.Equal(t, "\r", sfna.ReasonText())
			},
		},
		{
			with: "missing carriage return",
			input: func(t *testing.T) []byte {
				p := validStartFileNegative(t)
				p[len(p)-1] = 'd'
				return p
			},
			expect: func(t *testing.T, sfna oftp2.StartFileNegativeAnswerCmd) {
				require.EqualError(t, sfna.Valid(), "does not end on carriage return, but on d")
				require.Equal(t, oftp2.AnswerInvalidFilename, sfna.ReasonCode())
				require.Equal(t, true, sfna.Retry())
				require.Equal(t, "", sfna.ReasonText())
			},
		},
		{
			with: "corrupted reason code",
			input: func(t *testing.T) []byte {
				p := validStartFileNegative(t)
				p[2] = 'd'
				return p
			},
			expect: func(t *testing.T, sfna oftp2.StartFileNegativeAnswerCmd) {
				require.EqualError(t, sfna.Valid(), `invalid reason code`)
				require.Equal(t, oftp2.AnswerReason(0), sfna.ReasonCode())
				require.Equal(t, true, sfna.Retry())
				require.Equal(t, "", sfna.ReasonText())
			},
		},
		{
			with: "corrupted retry",
			input: func(t *testing.T) []byte {
				p := validStartFileNegative(t)
				p[3] = 'd'
				return p
			},
			expect: func(t *testing.T, sfna oftp2.StartFileNegativeAnswerCmd) {
				require.EqualError(t, sfna.Valid(), `invalid retry`)
				require.Equal(t, oftp2.AnswerInvalidFilename, sfna.ReasonCode())
				require.Equal(t, false, sfna.Retry())
				require.Equal(t, "", sfna.ReasonText())
			},
		},
		{
			with: "corrupted reason length",
			input: func(t *testing.T) []byte {
				file, err := oftp2.NewStartFileNegativeAnswer(oftp2.NegativeFileInput{
					Reason:     oftp2.AnswerInvalidFilename,
					ReasonText: "MY_TEXT",
				})
				require.NoError(t, err)
				file[5] = 'd'
				return file
			},
			expect: func(t *testing.T, sfna oftp2.StartFileNegativeAnswerCmd) {
				require.EqualError(t, sfna.Valid(), `strconv.Atoi: parsing "0d7": invalid syntax`)
				require.Equal(t, oftp2.AnswerInvalidFilename, sfna.ReasonCode())
				require.Equal(t, false, sfna.Retry())
				require.Equal(t, "MY_TEXT", sfna.ReasonText())
			},
		},
	} {
		t.Run(scenario.with, func(t *testing.T) {
			scenario.expect(t, scenario.input(t))
		})
	}
}

func validStartFileNegative(t *testing.T) oftp2.Command {
	file, err := oftp2.NewStartFileNegativeAnswer(oftp2.NegativeFileInput{
		Reason: oftp2.AnswerInvalidFilename,
		Retry: true,
	})
	require.NoError(t, err)
	return file
}

func generateLongString(size int) string {
	str := ""
	for i := 0; i < size; i++ {
		str += "a"
	}
	return str
}
