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

func (w *Wallet)getAddress() []byte {

	//get ripemd160Hash
	ripemd160Hash := Ripemd160Hash(w.PublicKey)
	fmt.Printf("ripemd160Hash = 0x%x and len= %d\n",ripemd160Hash,len(ripemd160Hash))
	//add version
	ripemd160HashWithVersion := append([]byte{version},ripemd160Hash...)
	fmt.Printf("ripemd160HashWithVersion = 0x%x and len= %d\n",ripemd160HashWithVersion,len(ripemd160HashWithVersion))
	//add checksum
	ripemd160WithChecksum := append(ripemd160HashWithVersion,CheckSumHash(ripemd160HashWithVersion)...)
	fmt.Printf("ripemd160WithChecksum = 0x%x and len= %d\n",ripemd160WithChecksum,len(ripemd160WithChecksum))
	data := Base58Encode(ripemd160WithChecksum)
	fmt.Printf("Base58Encode = 0x%x and len= %d\n",data,len(data))
	return data
}

func CheckSumHash(input []byte) []byte {
	//checkSumHash := sha256.New()
	//checkSumHash.Write(input)
	//checkSumHashbyte := checkSumHash.Sum(nil)
	//checkSumHash.Write(checkSumHashbyte)
	//return checkSumHash.Sum(nil)

	hash1 := sha256.Sum256(input)

	hash2 := sha256.Sum256(hash1[:])

	return hash2[:4]
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