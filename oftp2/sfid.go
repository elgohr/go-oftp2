package oftp2

import (
	"fmt"
	"log"
	"strings"
)

// o-------------------------------------------------------------------o
// |       SFID        Start File                                      |
// |                                                                   |
// |       Start File Phase           Speaker ----> Listener           |
// |-------------------------------------------------------------------|
// | Pos | Field     | Description                           | Format  |
// |-----+-----------+---------------------------------------+---------|
// |   0 | SFIDCMD   | SFID Command, 'H'                     | F X(1)  |
// |   1 | SFIDDSN   | Virtual File Dataset Name             | V X(26) |
// |  27 | SFIDRSV1  | Reserved                              | F X(3)  |
// |  30 | SFIDDATE  | Virtual File Date stamp, (CCYYMMDD)   | V 9(8)  |
// |  38 | SFIDTIME  | Virtual File Time stamp, (HHMMSScccc) | V 9(10) |
// |  48 | SFIDUSER  | User Data                             | V X(8)  |
// |  56 | SFIDDEST  | Destination                           | V X(25) |
// |  81 | SFIDORIG  | Originator                            | V X(25) |
// | 106 | SFIDFMT   | File Format (F/V/U/T)                 | F X(1)  |
// | 107 | SFIDLRECL | Maximum Record Size                   | V 9(5)  |
// | 112 | SFIDFSIZ  | File Size, 1K blocks                  | V 9(13) |
// | 125 | SFIDOSIZ  | Original File Size, 1K blocks         | V 9(13) |
// | 138 | SFIDREST  | Restart Position                      | V 9(17) |
// | 155 | SFIDSEC   | Security Level                        | F 9(2)  |
// | 157 | SFIDCIPH  | Cipher suite selection                | F 9(2)  |
// | 159 | SFIDCOMP  | File compression algorithm            | F 9(1)  |
// | 160 | SFIDENV   | File enveloping format                | F 9(1)  |
// | 161 | SFIDSIGN  | Signed EERP request                   | F X(1)  |
// | 162 | SFIDDESCL | Virtual File Description length       | V 9(3)  |
// | 165 | SFIDDESC  | Virtual File Description              | V T(n)  |
// o-------------------------------------------------------------------o
//
// https://datatracker.ietf.org/doc/html/rfc5024#section-5.3.3

type StartFileCmd []byte

func (c StartFileCmd) Valid() error {
	if Id(c[0]) != StartFile {
		return fmt.Errorf("wrong command id: %v", string(c[0]))
	}
	return nil
}

func (c StartFileCmd) Name() string {
	return strings.TrimSpace(string(c[1:27]))
}

func (c StartFileCmd) Date() Timestamp {
	t, err := NewTimeStamp(c[30:48])
	if err != nil {
		log.Println(err)
	}
	return t
}

func (c StartFileCmd) UserData() []byte {
	return c[48:56]
}

func (c StartFileCmd) Destination() Sid {
	return Sid(c[56:81])
}

func (c StartFileCmd) Origin() Sid {
	return Sid(c[81:106])
}

func NewStartFile(input StartFileInput) (Command, error) {
	if len(input.Name) > 26 {
		return nil, fmt.Errorf("name is too long: %v", input.Name)
	} else if len(input.UserData) > 8 {
		return nil, fmt.Errorf("user data is too long: %v", string(input.UserData))
	} else if err := input.Destination.Valid(); err != nil {
		return nil, err
	} else if err := input.Origin.Valid(); err != nil {
		return nil, err
	} else if _, exists := KnownFileFormats[input.Format]; !exists {
		return nil, fmt.Errorf("unknown file format: %v", string(input.Format))
	} else if input.MaxRecordSize < 0 || input.MaxRecordSize > 99999 {
		return nil, fmt.Errorf("invalid max record size: %d", input.MaxRecordSize)
	} else if input.TransmittedSize < 0 || input.TransmittedSize > 9999999999999 {
		return nil, fmt.Errorf("invalid transmitted size: %d", input.TransmittedSize)
	} else if input.OriginalSize < 0 || input.OriginalSize > 9999999999999 {
		return nil, fmt.Errorf("invalid original size: %d", input.OriginalSize)
	} else if _, exists := KnownSecurityLevels[input.Security]; !exists {
		return nil, fmt.Errorf("unknown security level: %d", input.Security)
	} else if input.Security == SecurityNoServices && input.Compression == NoCompression && input.TransmittedSize != input.OriginalSize {
		return nil, fmt.Errorf("transmitted size (%d) does not match original size (%d)", input.TransmittedSize, input.OriginalSize)
	} else if _, exists := KnownCiphers[input.Cipher]; !exists {
		return nil, fmt.Errorf("unknown cipher: %d", input.Cipher)
	} else if _, exists := KnownCompressions[input.Compression]; !exists {
		return nil, fmt.Errorf("unknown compression: %d", input.Compression)
	} else if length := len(input.Description); length > 999 {
		return nil, fmt.Errorf("description is too long: %d", length)
	}

	name, err := fillUpString(input.Name, 26)
	if err != nil {
		return nil, err
	}
	userData, err := fillUpString(string(input.UserData), 8)
	if err != nil {
		return nil, err
	}

	return Command(
		string(StartFile) +
			name +
			reserved(3) +
			input.Date.ToString() +
			userData +
			string(input.Destination) +
			string(input.Origin) +
			CarriageReturn), nil
}

type StartFileInput struct {
	Name            string
	Date            Timestamp
	UserData        []byte
	Destination     Sid
	Origin          Sid
	Format          FileFormat
	MaxRecordSize   int
	TransmittedSize int64
	OriginalSize    int64
	RestartPosition int64
	Security        SecurityLevel
	Cipher          Cipher
	Compression     Compression
	SignedReceipt   bool
	Description     string
}

type FileFormat byte

var KnownFileFormats = map[FileFormat]struct{}{
	FileFormatFixed:        {},
	FileFormatVariable:     {},
	FileFormatUnstructured: {},
	FileFormatText:         {},
}

const (
	FileFormatFixed        FileFormat = 'F'
	FileFormatVariable     FileFormat = 'V'
	FileFormatUnstructured FileFormat = 'U'
	FileFormatText         FileFormat = 'T'
)

type SecurityLevel int

var KnownSecurityLevels = map[SecurityLevel]struct{}{
	SecurityNoServices:         {},
	SecurityEncrypted:          {},
	SecuritySigned:             {},
	SecurityEncryptedAndSigned: {},
}

const (
	SecurityNoServices         SecurityLevel = 00
	SecurityEncrypted          SecurityLevel = 01
	SecuritySigned             SecurityLevel = 02
	SecurityEncryptedAndSigned SecurityLevel = 03
)

type Cipher int

var KnownCiphers = map[Cipher]CipherMapping{
	NoCipher: {},
	Cipher3DesEdeCbc3Key: {
		Symmetric:  "AES_256_CBC",
		Asymmetric: "RSA_PKCS1_15",
		Hashing:    "SHA-1",
	},
	CipherAes256Cbc: {
		Symmetric:  "3DES_EDE_CBC_3KEY",
		Asymmetric: "RSA_PKCS1_15",
		Hashing:    "SHA-1",
	},
}

const (
	NoCipher             Cipher = 00
	Cipher3DesEdeCbc3Key Cipher = 01
	CipherAes256Cbc      Cipher = 02
)

type CipherMapping struct {
	Symmetric  string
	Asymmetric string
	Hashing    string
}

type Compression int

var KnownCompressions = map[Compression]struct{}{
	NoCompression:   {},
	CompressionZlib: {},
}

const (
	NoCompression   Compression = 0
	CompressionZlib Compression = 1
)
