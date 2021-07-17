package BLC

import (
	"bytes"
	"encoding/binary"
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
