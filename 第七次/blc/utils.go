package blc

import(
	"bytes"
	"encoding/binary"
	"log"
	"encoding/json"
	"fmt"
	"encoding/gob"
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

func ReverseByte(input []byte) []byte {
	for i,j :=0,len(input)-1;i<j;i,j = i+1,j-1 {
		input[i],input[j] = input[j],input[i]
	}
	return  input
}



//version 转字节数组
func commandToBytes(command string) []byte {
	var bytes [COMMANDLENGTH]byte

	for i, c := range command {
		bytes[i] = byte(c)
	}

	return bytes[:]
}


//字节数组转version
func bytesToCommand(bytes []byte) string {
	var command []byte

	for _, b := range bytes {
		if b != 0x0 {
			command = append(command, b)
		}
	}

	return fmt.Sprintf("%s", command)
}


// 将结构体序列化成字节数组
func gobEncode(data interface{}) []byte {
	var buff bytes.Buffer

	enc := gob.NewEncoder(&buff)
	err := enc.Encode(data)
	if err != nil {
		log.Panic(err)
	}

	return buff.Bytes()
}