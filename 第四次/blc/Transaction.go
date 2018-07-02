package blc

import (
	"bytes"
	"encoding/gob"
	"log"
	"crypto/sha256"
	"encoding/hex"
)

type Transaction struct {
	TxHash []byte
	Vins []*TXInput
	Vouts []*TXOutput
}

func NewCoinbaseTransaction(address string) *Transaction {
	txInput := &TXInput{[]byte{},-1,"Genesis Data"}
	txOutput := &TXOutput{100,address}//coinbase get 100
	txCoinbase := &Transaction{[]byte{},[]*TXInput{txInput},[]*TXOutput{txOutput}}
	//设置hash值
	txCoinbase.HashTransaction()
	return txCoinbase
}

func (tx *Transaction) HashTransaction()  {
	var result bytes.Buffer
	encoder := gob.NewEncoder(&result)
	err := encoder.Encode(tx)
	if err != nil {
		log.Panic(err)
	}
	hash := sha256.Sum256(result.Bytes())
	tx.TxHash = hash[:]
}

func NewSimpleTransaction(from string,to string,amount int,blockchain *Blockchain,txs []*Transaction) *Transaction {
	//$ ./bc send -from '["juncheng"]' -to '["zhangqiang"]' -amount '["2"]'
	//	[juncheng]
	//	[zhangqiang]
	//	[2]
	// 通过一个函数，返回
	money,spendableUTXODic := blockchain.FindSpendableUTXOS(from,amount,txs)
	//
	//	{hash1:[0],hash2:[2,3]}

	var txIntputs []*TXInput
	var txOutputs []*TXOutput

	for txHash,indexArray := range spendableUTXODic  {

		txHashBytes,_ := hex.DecodeString(txHash)
		for _,index := range indexArray  {
			txInput := &TXInput{txHashBytes,index,from}
			txIntputs = append(txIntputs,txInput)
		}

	}

	// 转账
	txOutput := &TXOutput{int64(amount),to}
	txOutputs = append(txOutputs,txOutput)

	// 找零
	txOutput = &TXOutput{int64(money) - int64(amount),from}
	txOutputs = append(txOutputs,txOutput)

	tx := &Transaction{[]byte{},txIntputs,txOutputs}

	//设置hash值
	tx.HashTransaction()

	return tx

}

func (tx *Transaction) IsCoinbaseTransaction() bool {

	return len(tx.Vins[0].Txhash) == 0 && tx.Vins[0].Vout == -1
}
