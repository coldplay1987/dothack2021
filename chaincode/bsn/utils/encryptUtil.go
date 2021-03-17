/*
@Time : 2020/5/15 16:02
@Author : leixianting
@File : encryptUtil
@Software: GoLand
*/
package utils

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"encoding/hex"
	"golang.org/x/crypto/sha3"
)

type Address [20]byte

func GenerateAddress() (string, error) {

	prikey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)

	if err != nil {
		return "", err
	}

	pubKey := prikey.PublicKey

	pubBytes := elliptic.Marshal(elliptic.P256(), pubKey.X, pubKey.Y)

	return PubkeyBytesToAddress(pubBytes).Hex(), nil

}

// bytes 转 Address 类型
func PubkeyBytesToAddress(b []byte) Address {

	sha3Bytes := keccak256(b[1:])[12:]

	var a Address
	a.setBytes(sha3Bytes)
	return a
}

func (a *Address) setBytes(b []byte) {
	if len(b) > len(a) {
		b = b[len(b)-20:]
	}
	copy(a[20-len(b):], b)
}

// byte转16进制字符串显示
func (a Address) Hex() string {
	unchecksummed := hex.EncodeToString(a[:])
	sha := sha3.NewLegacyKeccak256()
	sha.Write([]byte(unchecksummed))
	hash := sha.Sum(nil)

	result := []byte(unchecksummed)
	for i := 0; i < len(result); i++ {
		hashByte := hash[i/2]
		if i%2 == 0 {
			hashByte = hashByte >> 4
		} else {
			hashByte &= 0xf
		}
		if result[i] > '9' && hashByte > 7 {
			result[i] -= 32
		}
	}
	return string(result)
}

// 通过sha3生成32字节的byte
func keccak256(data ...[]byte) []byte {
	d := sha3.NewLegacyKeccak256()
	for _, b := range data {
		d.Write(b)
	}
	return d.Sum(nil)
}
