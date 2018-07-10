package blc

import "bytes"

type TXInput struct {
	Txhash []byte
	Vout      int
	ScriptSig []byte
	PublicKey []byte
}

func (txinput *TXInput) UnLockRipemd160Hash(ripemd160Hash []byte) bool {
	publicKey := Ripemd160Hash(txinput.PublicKey)

	return bytes.Compare(publicKey,ripemd160Hash) == 0
}

