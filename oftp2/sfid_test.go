package oftp2_test

import (
	"bifroest/oftp2"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestStartFile(t *testing.T) {
	for _, scenario := range []struct {
		with   string
		input  func(t *testing.T) oftp2.StartFileInput
		expect func(t *testing.T, cmd oftp2.Command, err error)
	}{
		{
			with:  "a standard input",
			input: validStartFileInput,
			expect: func(t *testing.T, cmd oftp2.Command, err error) {
				require.NoError(t, err)
				require.NotNil(t, cmd)
			},
		},
		{
			with: "an exceeding filename",
			input: func(t *testing.T) oftp2.StartFileInput {
				i := validStartFileInput(t)
				i.Name = "123456789101112131415161718"
				return i
			},
			expect: func(t *testing.T, cmd oftp2.Command, err error) {
				require.EqualError(t, err, "name is too long: 123456789101112131415161718")
				require.Nil(t, cmd)
			},
		},
		{
			with: "an exceeding user data",
			input: func(t *testing.T) oftp2.StartFileInput {
				i := validStartFileInput(t)
				i.UserData = []byte("123456789")
				return i
			},
			expect: func(t *testing.T, cmd oftp2.Command, err error) {
				require.EqualError(t, err, "user data is too long: 123456789")
				require.Nil(t, cmd)
			},
		},
		{
			with: "an invalid destination",
			input: func(t *testing.T) oftp2.StartFileInput {
				i := validStartFileInput(t)
				i.Destination = []byte("!")
				return i
			},
			expect: func(t *testing.T, cmd oftp2.Command, err error) {
				require.EqualError(t, err, "expected the length of 25, but got 1")
				require.Nil(t, cmd)
			},
		},
		{
			with: "an invalid origin",
			input: func(t *testing.T) oftp2.StartFileInput {
				i := validStartFileInput(t)
				i.Origin = []byte("!")
				return i
			},
			expect: func(t *testing.T, cmd oftp2.Command, err error) {
				require.EqualError(t, err, "expected the length of 25, but got 1")
				require.Nil(t, cmd)
			},
		},
		{
			with: "an unknown file format",
			input: func(t *testing.T) oftp2.StartFileInput {
				i := validStartFileInput(t)
				i.Format = '?'
				return i
			},
			expect: func(t *testing.T, cmd oftp2.Command, err error) {
				require.EqualError(t, err, "unknown file format: ?")
				require.Nil(t, cmd)
			},
		},
		{
			with: "negative max record size",
			input: func(t *testing.T) oftp2.StartFileInput {
				i := validStartFileInput(t)
				i.MaxRecordSize = -1
				return i
			},
			expect: func(t *testing.T, cmd oftp2.Command, err error) {
				require.EqualError(t, err, "invalid max record size: -1")
				require.Nil(t, cmd)
			},
		},
		{
			with: "exceeding max record size",
			input: func(t *testing.T) oftp2.StartFileInput {
				i := validStartFileInput(t)
				i.MaxRecordSize = 100000
				return i
			},
			expect: func(t *testing.T, cmd oftp2.Command, err error) {
				require.EqualError(t, err, "invalid max record size: 100000")
				require.Nil(t, cmd)
			},
		},
		{
			with: "negative transmitted size",
			input: func(t *testing.T) oftp2.StartFileInput {
				i := validStartFileInput(t)
				i.TransmittedSize = -1
				return i
			},
			expect: func(t *testing.T, cmd oftp2.Command, err error) {
				require.EqualError(t, err, "invalid transmitted size: -1")
				require.Nil(t, cmd)
			},
		},
		{
			with: "exceeding transmitted size",
			input: func(t *testing.T) oftp2.StartFileInput {
				i := validStartFileInput(t)
				i.TransmittedSize = 10000000000000
				return i
			},
			expect: func(t *testing.T, cmd oftp2.Command, err error) {
				require.EqualError(t, err, "invalid transmitted size: 10000000000000")
				require.Nil(t, cmd)
			},
		},
		{
			with: "negative original size",
			input: func(t *testing.T) oftp2.StartFileInput {
				i := validStartFileInput(t)
				i.OriginalSize = -1
				return i
			},
			expect: func(t *testing.T, cmd oftp2.Command, err error) {
				require.EqualError(t, err, "invalid original size: -1")
				require.Nil(t, cmd)
			},
		},
		{
			with: "exceeding original size",
			input: func(t *testing.T) oftp2.StartFileInput {
				i := validStartFileInput(t)
				i.OriginalSize = 10000000000000
				return i
			},
			expect: func(t *testing.T, cmd oftp2.Command, err error) {
				require.EqualError(t, err, "invalid original size: 10000000000000")
				require.Nil(t, cmd)
			},
		},
		{
			with: "non matching size without compression and encryption",
			input: func(t *testing.T) oftp2.StartFileInput {
				i := validStartFileInput(t)
				i.TransmittedSize = 999
				i.OriginalSize = 1000
				i.Compression = oftp2.NoCompression
				i.Security = oftp2.SecurityNoServices
				return i
			},
			expect: func(t *testing.T, cmd oftp2.Command, err error) {
				require.EqualError(t, err, "transmitted size (999) does not match original size (1000)")
				require.Nil(t, cmd)
			},
		},
		{
			with: "unknown security level",
			input: func(t *testing.T) oftp2.StartFileInput {
				i := validStartFileInput(t)
				i.Security = -1
				return i
			},
			expect: func(t *testing.T, cmd oftp2.Command, err error) {
				require.EqualError(t, err, "unknown security level: -1")
				require.Nil(t, cmd)
			},
		},
		{
			with: "unknown cipher",
			input: func(t *testing.T) oftp2.StartFileInput {
				i := validStartFileInput(t)
				i.Cipher = -1
				return i
			},
			expect: func(t *testing.T, cmd oftp2.Command, err error) {
				require.EqualError(t, err, "unknown cipher: -1")
				require.Nil(t, cmd)
			},
		},
		{
			with: "unknown compression",
			input: func(t *testing.T) oftp2.StartFileInput {
				i := validStartFileInput(t)
				i.Compression = -1
				return i
			},
			expect: func(t *testing.T, cmd oftp2.Command, err error) {
				require.EqualError(t, err, "unknown compression: -1")
				require.Nil(t, cmd)
			},
		},
		{
			with: "exceeding description",
			input: func(t *testing.T) oftp2.StartFileInput {
				i := validStartFileInput(t)
				i.Description = generateLongString(1000)
				return i
			},
			expect: func(t *testing.T, cmd oftp2.Command, err error) {
				require.EqualError(t, err, "description is too long: 1000")
				require.Nil(t, cmd)
			},
		},
	} {
		t.Run(scenario.with, func(t *testing.T) {
			s, err := oftp2.NewStartFile(scenario.input(t))
			scenario.expect(t, s, err)
		})
	}
}

func TestStartFile_Valid(t *testing.T) {
	for _, scenario := range []struct {
		with   string
		input  func(t *testing.T) []byte
		expect func(t *testing.T, sfid oftp2.StartFileCmd)
	}{
		{
			with: "a standard message",
			input: func(t *testing.T) []byte {
				return validStartFile(t)
			},
			expect: func(t *testing.T, sfid oftp2.StartFileCmd) {
				require.NoError(t, sfid.Valid())
				//require.Equal(t, oftp2.AnswerInvalidFilename, sfna.ReasonCode())
				//require.Equal(t, false, sfna.Retry())
				//require.Equal(t, "MY_TEXT", sfna.ReasonText())
			},
		},
		{
			with: "a wrong cmd type",
			input: func(t *testing.T) []byte {
				p := validStartFile(t)
				p[0] = '^'
				return p
			},
			expect: func(t *testing.T, sfid oftp2.StartFileCmd) {
				require.EqualError(t, sfid.Valid(), "wrong command id: ^")
				require.Equal(t, "MY_FILE", sfid.Name())
				stamp, err := oftp2.NewTimeStamp([]byte("20200102030405060708"))
				require.NoError(t, err)
				require.Equal(t, stamp, sfid.Date())
				require.Equal(t, []byte("        "), sfid.UserData())
				destinationSid, err := oftp2.NewSid(oftp2.SidInput{
					CodeDesignator:   "Test",
					OrganisationCode: "Org",
					SubAddress:       "Sender",
				})
				require.NoError(t, err)
				require.Equal(t, destinationSid, sfid.Destination())
				originSid, err := oftp2.NewSid(oftp2.SidInput{
					CodeDesignator:   "Test",
					OrganisationCode: "Org",
					SubAddress:       "Origin",
				})
				require.NoError(t, err)
				require.Equal(t, originSid, sfid.Origin())
			},
		},
		//{
		//	with: "a wrong length",
		//	input: func(t *testing.T) []byte {
		//		p := validStartFileNegative(t)
		//		return append(p, ' ')
		//	},
		//	expect: func(t *testing.T, sfna oftp2.StartFileNegativeAnswerCmd) {
		//		require.EqualError(t, sfna.Valid(), "expected the length of 8, but got 9")
		//		require.Equal(t, oftp2.AnswerInvalidFilename, sfna.ReasonCode())
		//		require.Equal(t, true, sfna.Retry())
		//		require.Equal(t, "\r", sfna.ReasonText())
		//	},
		//},
		//{
		//	with: "missing carriage return",
		//	input: func(t *testing.T) []byte {
		//		p := validStartFileNegative(t)
		//		p[len(p)-1] = 'd'
		//		return p
		//	},
		//	expect: func(t *testing.T, sfna oftp2.StartFileNegativeAnswerCmd) {
		//		require.EqualError(t, sfna.Valid(), "does not end on carriage return, but on d")
		//		require.Equal(t, oftp2.AnswerInvalidFilename, sfna.ReasonCode())
		//		require.Equal(t, true, sfna.Retry())
		//		require.Equal(t, "", sfna.ReasonText())
		//	},
		//},
		//{
		//	with: "corrupted reason code",
		//	input: func(t *testing.T) []byte {
		//		p := validStartFileNegative(t)
		//		p[2] = 'd'
		//		return p
		//	},
		//	expect: func(t *testing.T, sfna oftp2.StartFileNegativeAnswerCmd) {
		//		require.EqualError(t, sfna.Valid(), `invalid reason code`)
		//		require.Equal(t, oftp2.AnswerReason(0), sfna.ReasonCode())
		//		require.Equal(t, true, sfna.Retry())
		//		require.Equal(t, "", sfna.ReasonText())
		//	},
		//},
		//{
		//	with: "corrupted retry",
		//	input: func(t *testing.T) []byte {
		//		p := validStartFileNegative(t)
		//		p[3] = 'd'
		//		return p
		//	},
		//	expect: func(t *testing.T, sfna oftp2.StartFileNegativeAnswerCmd) {
		//		require.EqualError(t, sfna.Valid(), `invalid retry`)
		//		require.Equal(t, oftp2.AnswerInvalidFilename, sfna.ReasonCode())
		//		require.Equal(t, false, sfna.Retry())
		//		require.Equal(t, "", sfna.ReasonText())
		//	},
		//},
		//{
		//	with: "corrupted reason length",
		//	input: func(t *testing.T) []byte {
		//		file, err := oftp2.NewStartFileNegativeAnswer(oftp2.NegativeFileInput{
		//			Reason:     oftp2.AnswerInvalidFilename,
		//			ReasonText: "MY_TEXT",
		//		})
		//		require.NoError(t, err)
		//		file[5] = 'd'
		//		return file
		//	},
		//	expect: func(t *testing.T, sfna oftp2.StartFileNegativeAnswerCmd) {
		//		require.EqualError(t, sfna.Valid(), `strconv.Atoi: parsing "0d7": invalid syntax`)
		//		require.Equal(t, oftp2.AnswerInvalidFilename, sfna.ReasonCode())
		//		require.Equal(t, false, sfna.Retry())
		//		require.Equal(t, "MY_TEXT", sfna.ReasonText())
		//	},
		//},
	} {
		t.Run(scenario.with, func(t *testing.T) {
			scenario.expect(t, scenario.input(t))
		})
	}
}

func validStartFile(t *testing.T) oftp2.Command {
	file, err := oftp2.NewStartFile(validStartFileInput(t))
	require.NoError(t, err)
	return file
}

func validStartFileInput(t *testing.T) oftp2.StartFileInput {
	stamp, err := oftp2.NewTimeStamp([]byte("20200102030405060708"))
	require.NoError(t, err)
	destinationSid, err := oftp2.NewSid(oftp2.SidInput{
		CodeDesignator:   "Test",
		OrganisationCode: "Org",
		SubAddress:       "Sender",
	})
	require.NoError(t, err)
	originSid, err := oftp2.NewSid(oftp2.SidInput{
		CodeDesignator:   "Test",
		OrganisationCode: "Org",
		SubAddress:       "Origin",
	})
	require.NoError(t, err)
	return oftp2.StartFileInput{
		Name:            "MY_FILE",
		Date:            stamp,
		UserData:        []byte("        "),
		Destination:     destinationSid,
		Origin:          originSid,
		Format:          oftp2.FileFormatFixed,
		MaxRecordSize:   10,
		TransmittedSize: 10,
		OriginalSize:    20,
		RestartPosition: 0,
		Security:        oftp2.SecurityEncrypted,
		Cipher:          oftp2.CipherAes256Cbc,
		Compression:     oftp2.NoCompression,
		SignedReceipt:   false,
		Description:     "Description",
	}
}
