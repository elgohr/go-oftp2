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

func intToHexBytes(i int32) []byte {
	b := new(bytes.Buffer)
	if err := binary.Write(b, binary.BigEndian, i); err != nil {
		log.Println(err)
		return nil
	}
	return b.Bytes()
}
