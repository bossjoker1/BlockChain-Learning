package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"fmt"
	"log"
	"math/big"
)

func main() {
	curve := elliptic.P256()
	// 基于椭圆加密
	priv, err := ecdsa.GenerateKey(curve, rand.Reader)
	if err != nil {
		log.Panicf("ecdsa generate key failed. %v\n", err)
	}

	pubKey := append(priv.PublicKey.X.Bytes(), priv.PublicKey.Y.Bytes()...)

	r, s, err := ecdsa.Sign(rand.Reader, priv, []byte("test ecdsa"))

	x := big.Int{}
	y := big.Int{}
	pubKeyLen := len(pubKey)
	x.SetBytes(pubKey[:(pubKeyLen / 2)])
	y.SetBytes(pubKey[(pubKeyLen / 2):])

	rawPubkey := ecdsa.PublicKey{Curve: curve, X: &x, Y: &y}

	if !ecdsa.Verify(&rawPubkey, []byte("test ecdsa"), r, s) {
		fmt.Println("false!")
	} else {
		fmt.Println("true")
	}

}
