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
			name:  "with a standard input",
			input: negativeInput{
				reasonCode: oftp2.AnswerInvalidFilename,
				retry: true,
				reasonText: "BECAUSE",
			},
			expect: func(t *testing.T, cmd oftp2.Command, err error) {
				require.NoError(t, err)
				require.NotNil(t, cmd)
			},
		},
		{
			name:  "with an unknown reasonCode",
			input: negativeInput{
				reasonCode: 98,
				retry: true,
				reasonText: "",
			},
			expect: func(t *testing.T, cmd oftp2.Command, err error) {
				require.EqualError(t, err, "unknown answer reason: 98")
				require.Nil(t, cmd)
			},
		},
		{
			name:  "with an reason text that is too long",
			input: negativeInput{
				reasonCode: oftp2.AnswerInvalidFilename,
				retry: false,
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

type negativeInput struct {
	reasonCode oftp2.AnswerReason
	retry bool
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
