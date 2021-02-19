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

type Cmd byte

const (
	StartSessionReadyMessage Cmd = 'I'
	StartSessionMessage      Cmd = 'X'
	StartFilePositiveMessage Cmd = '2'
	StartFileNegativeMessage Cmd = '3'
	Unknown                  Cmd = '0'
)

var KnownCommands = map[Cmd]struct{}{
	StartSessionReadyMessage: {},
	StartSessionMessage:      {},
	StartFilePositiveMessage: {},
	StartFileNegativeMessage: {},
}

func (c Command) Cmd() Cmd {
	if len(c) == 0 {
		return Unknown
	}
	if _, exists := KnownCommands[Cmd(c[0])]; !exists {
		return Unknown
	}
	return Cmd(c[0])
}

func intToHexBytes(i int32) []byte {
	b := new(bytes.Buffer)
	if err := binary.Write(b, binary.BigEndian, i); err != nil {
		log.Println(err)
		return nil
	}
	return b.Bytes()
}
