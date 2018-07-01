package blc

import(
	"bytes"
	"encoding/binary"
	"log"
)

func Int2Byte(number int64) []byte {
	buff := new(bytes.Buffer)
	err := binary.Write(buff,binary.BigEndian,number)
	if err != nil {
		log.Panic(err)
	}
	return buff.Bytes()
}
