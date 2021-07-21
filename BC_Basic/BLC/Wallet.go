package BLC

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"golang.org/x/crypto/ripemd160"
	"log"
)

// 钱包结构

type Wallet struct {
	// 私钥
	PrivateKey ecdsa.PrivateKey

	// 公钥
	PublicKey []byte
}

// 创建钱包
func NewWallet() *Wallet {
	privKey, pubKey := NewPair()
	return &Wallet{privKey, pubKey}
}

// 生成公-私钥对
func NewPair() (ecdsa.PrivateKey, []byte) {
	curve := elliptic.P256()
	// 基于椭圆加密
	priv, err := ecdsa.GenerateKey(curve, rand.Reader)
	if err != nil {
		log.Panicf("ecdsa generate key failed. %v\n", err)
	}

	pubKey := append(priv.PublicKey.X.Bytes(), priv.PublicKey.Y.Bytes()...)
	return *priv, pubKey
}

// 对公钥进行双hash
func Ripemd160_SHA256(pubKey []byte) []byte {
	// 1. SHA256
	hash256 := sha256.New()
	hash256.Write(pubKey)
	hash := hash256.Sum(nil)

	// REPEMD160
	rmd160 := ripemd160.New()
	rmd160.Write(hash)

	return rmd160.Sum(nil)

}

// 通过钱包获取地址
func (w *Wallet) GetAddr() []byte {
	// 获取 pubKey hash
	pubHash := Ripemd160_SHA256(w.PublicKey)
	// 获取version, 加到前缀
	verison_pubHash := append([]byte{VERSION}, pubHash...)
	// 生成checkSum
	checkSum := CheckSum(verison_pubHash)
	bytes := append(verison_pubHash, checkSum...)

	// Base58编码
	return Base58Encode(bytes)
}

// 生成校验和
// 两次hash,
func CheckSum(payload []byte) []byte {
	first_hash := sha256.Sum256(payload)
	second_hash := sha256.Sum256(first_hash[:])

	return second_hash[:CHECKSUMLEN] // 取固定长度

}

// 判断地址有效性
func IsValidforAddr(addr []byte) bool {
	// base58解码
	decodeAddr := Base58Decode(addr)
	// 拆分，checkSum发挥作用
	checkSum := decodeAddr[len(decodeAddr)-CHECKSUMLEN:]
	version_pubKeyHash := decodeAddr[:len(decodeAddr)-CHECKSUMLEN]
	checkBytes := CheckSum(version_pubKeyHash)
	if bytes.Compare(checkSum, checkBytes) == 0 {
		return true
	}

	return false
}
