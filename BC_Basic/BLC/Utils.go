package BLC

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"log"
)

// int 转 []byte
func IntToHex(data int64) []byte {
	// 创建缓冲
	buffer := new(bytes.Buffer)
	err := binary.Write(buffer, binary.BigEndian, data)
	if err != nil {
		log.Panicf("int failed to []byte. %v\n", err)
	}
	return buffer.Bytes()
}

// Json 与 Array 之间转换
// windows下的转义问题
// main.exe send -from "[\"Amy\",\"Bob\"]" -to "[\"Bob\",\"James\"]" -amount "[\"10\",\"5\"]"

func JsonToArray(jsonString string) []string {
	var strArr []string

	if err := json.Unmarshal([]byte(jsonString), &strArr); err != nil {
		log.Panicf("Unmarshal to Array failed. %v\n ", err)
	}
	return strArr
}
