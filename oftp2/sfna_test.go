package oftp2_test

import (
	"bifroest/oftp2"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestStartFileNegativeAnswer(t *testing.T) {
	for _, scenario := range []struct {
		name   string
		input  negativeInput
		expect func(t *testing.T, cmd oftp2.Command, err error)
	}{
		{
			name: "with a standard input",
			input: negativeInput{
				reasonCode: oftp2.AnswerInvalidFilename,
				retry:      true,
				reasonText: "BECAUSE",
			},
			expect: func(t *testing.T, cmd oftp2.Command, err error) {
				require.NoError(t, err)
				require.NotNil(t, cmd)
			},
		},
		{
			name: "with an unknown reasonCode",
			input: negativeInput{
				reasonCode: 98,
				retry:      true,
				reasonText: "",
			},
			expect: func(t *testing.T, cmd oftp2.Command, err error) {
				require.EqualError(t, err, "unknown answer reason: 98")
				require.Nil(t, cmd)
			},
		},
		{
			name: "with an reason text that is too long",
			input: negativeInput{
				reasonCode: oftp2.AnswerInvalidFilename,
				retry:      false,
				reasonText: generateLongString(1000),
			},
			expect: func(t *testing.T, cmd oftp2.Command, err error) {
				require.EqualError(t, err, "reason text is too long")
				require.Nil(t, cmd)
			},
		},
	} {
		t.Run(scenario.name, func(t *testing.T) {
			in := scenario.input
			s, err := oftp2.NewStartFileNegativeAnswer(in.reasonCode, in.retry, in.reasonText)
			scenario.expect(t, s, err)
		})
	}
}

func TestStartFileNegativeAnswer_Valid(t *testing.T) {
	for _, scenario := range []struct {
		name   string
		input  func(t *testing.T) []byte
		expect func(t *testing.T, sfna oftp2.StartFileNegativeAnswerCmd)
	}{
		{
			name: "with a standard message",
			input: func(t *testing.T) []byte {
				file, err := oftp2.NewStartFileNegativeAnswer(oftp2.AnswerInvalidFilename, false, "MY_TEXT")
				require.NoError(t, err)
				return file
			},
			expect: func(t *testing.T, sfna oftp2.StartFileNegativeAnswerCmd) {
				require.NoError(t, sfna.Valid())
				require.Equal(t, oftp2.AnswerInvalidFilename, sfna.ReasonCode())
				require.Equal(t, false, sfna.Retry())
				require.Equal(t, "MY_TEXT", string(sfna.ReasonText()))
			},
		},
		{
			name: "with a wrong cmd type",
			input: func(t *testing.T) []byte {
				p := validStartFileNegative(t)
				p[0] = '^'
				return p
			},
			expect: func(t *testing.T, sfna oftp2.StartFileNegativeAnswerCmd) {
				require.EqualError(t, sfna.Valid(), "does not start with 3, but with ^")
				require.Equal(t, oftp2.AnswerInvalidFilename, sfna.ReasonCode())
				require.Equal(t, true, sfna.Retry())
				require.Equal(t, "", string(sfna.ReasonText()))
			},
		},
		{
			name: "with a wrong length",
			input: func(t *testing.T) []byte {
				p := validStartFileNegative(t)
				return append(p, ' ')
			},
			expect: func(t *testing.T, sfna oftp2.StartFileNegativeAnswerCmd) {
				require.EqualError(t, sfna.Valid(), "expected the length of 8, but got 9")
				require.Equal(t, oftp2.AnswerInvalidFilename, sfna.ReasonCode())
				require.Equal(t, true, sfna.Retry())
				require.Equal(t, "\r", string(sfna.ReasonText()))
			},
		},
		{
			name: "missing carriage return",
			input: func(t *testing.T) []byte {
				p := validStartFileNegative(t)
				p[len(p)-1] = 'd'
				return p
			},
			expect: func(t *testing.T, sfna oftp2.StartFileNegativeAnswerCmd) {
				require.EqualError(t, sfna.Valid(), "does not end on carriage return, but on d")
				require.Equal(t, oftp2.AnswerInvalidFilename, sfna.ReasonCode())
				require.Equal(t, true, sfna.Retry())
				require.Equal(t, "", string(sfna.ReasonText()))
			},
		},
		{
			name: "with corrupted reason code",
			input: func(t *testing.T) []byte {
				p := validStartFileNegative(t)
				p[2] = 'd'
				return p
			},
			expect: func(t *testing.T, sfna oftp2.StartFileNegativeAnswerCmd) {
				require.EqualError(t, sfna.Valid(), `invalid reason code`)
				require.Equal(t, oftp2.AnswerReason(0), sfna.ReasonCode())
				require.Equal(t, true, sfna.Retry())
				require.Equal(t, "", string(sfna.ReasonText()))
			},
		},
		{
			name: "with corrupted retry",
			input: func(t *testing.T) []byte {
				p := validStartFileNegative(t)
				p[3] = 'd'
				return p
			},
			expect: func(t *testing.T, sfna oftp2.StartFileNegativeAnswerCmd) {
				require.EqualError(t, sfna.Valid(), `invalid retry`)
				require.Equal(t, oftp2.AnswerInvalidFilename, sfna.ReasonCode())
				require.Equal(t, false, sfna.Retry())
				require.Equal(t, "", string(sfna.ReasonText()))
			},
		},
		{
			name: "with corrupted reason length",
			input: func(t *testing.T) []byte {
				file, err := oftp2.NewStartFileNegativeAnswer(oftp2.AnswerInvalidFilename, false, "MY_TEXT")
				require.NoError(t, err)
				file[5] = 'd'
				return file
			},
			expect: func(t *testing.T, sfna oftp2.StartFileNegativeAnswerCmd) {
				require.EqualError(t, sfna.Valid(), `strconv.Atoi: parsing "0d7": invalid syntax`)
				require.Equal(t, oftp2.AnswerInvalidFilename, sfna.ReasonCode())
				require.Equal(t, false, sfna.Retry())
				require.Equal(t, "MY_TEXT", string(sfna.ReasonText()))
			},
		},
	} {
		t.Run(scenario.name, func(t *testing.T) {
			scenario.expect(t, scenario.input(t))
		})
	}
}

type negativeInput struct {
	reasonCode oftp2.AnswerReason
	retry      bool
	reasonText string
}

func validStartFileNegative(t *testing.T) oftp2.Command {
	file, err := oftp2.NewStartFileNegativeAnswer(oftp2.AnswerInvalidFilename, true, "")
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
