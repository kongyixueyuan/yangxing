package blc

import(
	"bytes"
	"encoding/binary"
	"log"
	"encoding/json"
)

func Int2Byte(number int64) []byte {
	buff := new(bytes.Buffer)
	err := binary.Write(buff,binary.BigEndian,number)
	if err != nil {
		log.Panic(err)
	}
	return buff.Bytes()
}

func JSONToArray(jsonString string) []string {

	//json 到 []string
	var sArr []string
	if err := json.Unmarshal([]byte(jsonString), &sArr); err != nil {
		log.Panic(err)
	}
	return sArr
}