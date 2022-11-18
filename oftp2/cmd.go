package oftp2

import (
	"bytes"
	"encoding/binary"
	"log"
)

type Command []byte

// Header Length
const StreamTransmissionHeaderLength = 4

func (c Command) StreamTransmissionBuffer() []byte {
	sth := intToHexBytes(int32(len(c) + StreamTransmissionHeaderLength))
	sth[0] = 0x10
	return append(sth, c...)
}

type Id byte

const (
	StartSessionReadyMessage  Id = 'I'
	StartSessionMessage       Id = 'X'
	SidId                     Id = 'O'
	StartFile                 Id = 'H'
	StartFilePositiveMessage  Id = '2'
	StartFileNegativeMessage  Id = '3'
	DataExchangeBufferMessage Id = 'D'
	Unknown                   Id = '0'
)

func (i Id) Byte() byte {
	return byte(i)
}

func (i Id) String() string {
	return string(i)
}

var KnownIds = map[Id]struct{}{
	StartSessionReadyMessage: {},
	StartSessionMessage:      {},
	StartFilePositiveMessage: {},
	StartFileNegativeMessage: {},
}

func (c Command) Cmd() Id {
	if len(c) == 0 {
		return Unknown
	}
	if _, exists := KnownIds[Id(c[0])]; !exists {
		return Unknown
	}
	return Id(c[0])
}

func intToHexBytes(i int32) []byte {
	b := &bytes.Buffer{}
	if err := binary.Write(b, binary.BigEndian, i); err != nil {
		log.Println(err)
		return nil
	}
	return b.Bytes()
}
