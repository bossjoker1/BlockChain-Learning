package BLC

import (
	"bytes"
	"encoding/binary"
	"encoding/gob"
	"encoding/json"
	"fmt"
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

// 反转切片
func Reverse(data []byte) {
	for i, j := 0, len(data)-1; i < j; i, j = i+1, j-1 {
		data[i], data[j] = data[j], data[i]
	}
}

// 将结构体序列化为字节数组
func GobEncode(data interface{}) []byte {
	var buff bytes.Buffer
	enc := gob.NewEncoder(&buff)
	err := enc.Encode(data)
	if err != nil {
		log.Panicf("encode the data failed. %v\n", err)
	}
	return buff.Bytes()
}

// 命令格式
// xxxxx(指令)xxxxx(数据...)

// 将命令转化为字节数组
// 指令长度最长为12位

func CommandToBytes(cmd string) []byte {
	var cmdbytes [12]byte
	for i, c := range cmd {
		cmdbytes[i] = byte(c)
	}

	return cmdbytes[:]
}

// 将字节数组转化成命令

func bytesToCommand(cmdbytes []byte) string {
	var cmd []byte // 接受命令

	for _, b := range cmdbytes {
		if b != 0x0 {
			cmd = append(cmd, b)
		}
	}
	return fmt.Sprintf("%s", cmd)
}
