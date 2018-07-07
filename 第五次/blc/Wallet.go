package blc

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"fmt"
	"crypto/rand"
	"log"
	"crypto/sha256"
	"golang.org/x/crypto/ripemd160"
)

const version = byte(0x00)

type Wallet struct {
	PrivateKey ecdsa.PrivateKey
	PublicKey []byte
}

func (w *Wallet)getAddress(publicKey []byte) []byte {

	//get ripemd160Hash
	ripemd160Hash := Ripemd160Hash(publicKey)

	//add version
	ripemd160HashWithVersion := append([]byte{version},ripemd160Hash...)
	//add checksum
	ripemd160WithChecksum := append(ripemd160HashWithVersion,CheckSumHash(ripemd160HashWithVersion)...)

	return Base58Encode(ripemd160WithChecksum)
}

func CheckSumHash(input []byte) []byte {
	checkSumHash := sha256.New()
	checkSumHash.Write(input)
	checkSumHashbyte := checkSumHash.Sum(nil)
	checkSumHash.Write(checkSumHashbyte)
	return checkSumHash.Sum(nil)
}

func Ripemd160Hash(publicKey [] byte) []byte {

	hash := sha256.New()
	hash.Write(publicKey)
	bytehash := hash.Sum(nil)

	hash160 := ripemd160.New()
	hash160.Write(bytehash)
	return hash160.Sum(nil)
}

func NewWallet() *Wallet {
	publicKey, privateKey := newPairKey()
	return &Wallet{publicKey,privateKey}
}

func newPairKey() (ecdsa.PrivateKey,[]byte) {
	curve := elliptic.P256()
	fmt.Println(curve)
	privatekey, err := ecdsa.GenerateKey(curve,rand.Reader)
	if err != nil {
		log.Panic(err)
	}
	publickey := append(privatekey.X.Bytes(),privatekey.X.Bytes()...)
	fmt.Println(publickey)
	return *privatekey, publickey
}