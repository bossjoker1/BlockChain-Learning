package BLC

import (
	"BlockChain-Learning/BC_Basic/Utils"
	"bytes"
	"math/big"
)

// 实现base58编码

// base58字符表
var b58Alphabet = []byte("123456789" +
	"ABCDEFGHJKLMNPQRSTUVWXYZ" +
	"abcdefghijkmnopqrstuvwxyz")

// 编码函数
func Base58Encode(input []byte) []byte {
	var result []byte
	x := big.NewInt(0).SetBytes(input) // bytes 转换为bigint
	//fmt.Printf("x : %v\n", x)
	base := big.NewInt(int64(len(b58Alphabet))) // 设置一个base58求模的基数
	zero := big.NewInt(0)
	mod := &big.Int{} // 余数
	for x.Cmp(zero) != 0 {
		x.DivMod(x, base, mod) // 求余
		// 以余数为下标，取值
		result = append(result, b58Alphabet[mod.Int64()])
	}
	// 反转切片
	Utils.Reverse(result)

	for b := range input { // b 代表切片下标
		if b == 0x00 {
			result = append([]byte{b58Alphabet[0]}, result...)
		} else {
			break
		}
	}
	// fmt.Printf("result : %s\n", result)
	return result
}

// 解码函数
func Base58Decode(input []byte) []byte {
	result := big.NewInt(0)
	zeroBytes := 0
	for b := range input {
		if b == 0x00 {
			zeroBytes++
		}
	}
	data := input[zeroBytes:]
	for _, b := range data {
		// 获取bytes数组中指定数字第一次出现的索引
		charIndex := bytes.IndexByte(b58Alphabet, b)
		result.Mul(result, big.NewInt(58))
		result.Add(result, big.NewInt(int64(charIndex)))
	}
	decoded := result.Bytes()
	decoded = append(bytes.Repeat([]byte{byte(0x00)}, zeroBytes), decoded...)
	return decoded
}
