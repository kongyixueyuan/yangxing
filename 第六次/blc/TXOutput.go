package blc

import "bytes"

//money
//addr
type TXOutput struct {
	Value int64
	PubKeyHash  []byte//160Hash
}

func (tx *TXOutput)UnLockScriptPubKeyWithAddress(inputaddr string) bool {
	//inputaddr is public addr
	PubHash := Base58Decode([]byte(inputaddr))
	Get160Hash := PubHash[1:len(PubHash)-4]
	return bytes.Compare(tx.PubKeyHash,Get160Hash) == 0
}

func (tx *TXOutput)Set160Hash(addr string) []byte {
	PubHash := Base58Decode([]byte(addr))
	Get160Hash := PubHash[1:len(PubHash)-4]

	return Get160Hash
}

func NewTxOutput(value int64, addr string) *TXOutput {
	txoutput := &TXOutput{value,nil}
	txoutput.PubKeyHash = txoutput.Set160Hash(addr)

	return txoutput
}