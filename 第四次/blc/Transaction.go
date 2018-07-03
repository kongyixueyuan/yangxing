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
	//创世区块没有消费，没有交易，没有输入
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
	//需要组装一个最新的transaction
	//首先需要找到from用户的可以满足value的可以花费的OUTput
	//然后根据组装transaction的方式来组装新的transaction
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
