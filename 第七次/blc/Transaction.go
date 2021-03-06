package blc

import (
	"bytes"
	"encoding/gob"
	"log"
	"crypto/sha256"
	"encoding/hex"
	"crypto/elliptic"
	"math/big"
	"crypto/ecdsa"
	"crypto/rand"
	"time"
)

type Transaction struct {
	TxHash []byte
	Vins []*TXInput
	Vouts []*TXOutput
}

func NewCoinbaseTransaction(address string) *Transaction {
	//创世区块没有消费，没有交易，没有输入
	txInput := &TXInput{[]byte{},-1,nil,[]byte{}}

	txOutput := NewTxOutput(100,address)

	//txOutput := &TXOutput{100,address}//coinbase get 100
	txCoinbase := &Transaction{[]byte{},[]*TXInput{txInput},[]*TXOutput{txOutput}}
	//设置hash值
	txCoinbase.HashTransaction()
	log.Println()
	log.Println("创世区块hash = ",txCoinbase.TxHash)
	return txCoinbase
}

func (tx *Transaction) HashTransaction()  {
	var result bytes.Buffer

	encoder := gob.NewEncoder(&result)

	err := encoder.Encode(tx)
	if err != nil {
		log.Panic(err)
	}

	resultBytes := bytes.Join([][]byte{Int2Byte(time.Now().Unix()),result.Bytes()},[]byte{})

	hash := sha256.Sum256(resultBytes)

	tx.TxHash = hash[:]
}


func NewSimpleTransaction(from string,to string,amount int64,utxoSet *UTXOSet,txs []*Transaction,nodeID string) *Transaction {
	//需要组装一个最新的transaction
	//首先需要找到from用户的可以满足value的可以花费的OUTput
	//然后根据组装transaction的方式来组装新的transaction


	wallets,_ := YX_NewWallets(nodeID)
	wallet := wallets.Walletmap[from]
	money,spendableUTXODic := utxoSet.FindSpendableUTXOS(from,amount,txs)
	//
	//	{hash1:[0],hash2:[2,3]}

	var txIntputs []*TXInput
	var txOutputs []*TXOutput

	for txHash,indexArray := range spendableUTXODic  {

		txHashBytes,_ := hex.DecodeString(txHash)
		for _,index := range indexArray  {

			txInput := &TXInput{txHashBytes,index,nil,wallet.PublicKey}
			txIntputs = append(txIntputs,txInput)
		}

	}

	// 转账
	//txOutput := &TXOutput{int64(amount),to}
	txOutput := NewTxOutput(int64(amount),to)
	txOutputs = append(txOutputs,txOutput)

	// 找零
	//txOutput = &TXOutput{int64(money) - int64(amount),from}
	txOutput = NewTxOutput(int64(money) - int64(amount),from)
	txOutputs = append(txOutputs,txOutput)

	tx := &Transaction{[]byte{},txIntputs,txOutputs}

	//设置hash值
	tx.HashTransaction()
	//utxoSet.blockchain.SignTransaction(tx, wallet.PrivateKey)
	return tx

}

func (tx *Transaction) IsCoinbaseTransaction() bool {

	return len(tx.Vins[0].Txhash) == 0 && tx.Vins[0].Vout == -1
}


func (tx *Transaction) Sign(privKey ecdsa.PrivateKey, prevTXs map[string]Transaction) {

	if tx.IsCoinbaseTransaction() {
		return
	}


	for _, vin := range tx.Vins {
		if prevTXs[hex.EncodeToString(vin.Txhash)].TxHash == nil {
			log.Panic("ERROR: Previous transaction is not correct")
		}
	}


	txCopy := tx.TrimmedCopy()

	for inID, vin := range txCopy.Vins {
		prevTx := prevTXs[hex.EncodeToString(vin.Txhash)]
		txCopy.Vins[inID].ScriptSig = nil
		txCopy.Vins[inID].PublicKey = prevTx.Vouts[vin.Vout].PubKeyHash
		txCopy.TxHash = txCopy.Hash()
		txCopy.Vins[inID].PublicKey = nil

		// 签名代码
		r, s, err := ecdsa.Sign(rand.Reader, &privKey, txCopy.TxHash)
		if err != nil {
			log.Panic(err)
		}
		signature := append(r.Bytes(), s.Bytes()...)

		tx.Vins[inID].ScriptSig = signature
	}
}


// 拷贝一份新的Transaction用于签名                                    T
func (tx *Transaction) TrimmedCopy() Transaction {
	var inputs []*TXInput
	var outputs []*TXOutput

	for _, vin := range tx.Vins {
		inputs = append(inputs, &TXInput{vin.Txhash, vin.Vout, nil, nil})
	}

	for _, vout := range tx.Vouts {
		outputs = append(outputs, &TXOutput{vout.Value, vout.PubKeyHash})
	}

	txCopy := Transaction{tx.TxHash, inputs, outputs}

	return txCopy
}


// 数字签名验证

func (tx *Transaction) Verify(prevTXs map[string]Transaction) bool {
	if tx.IsCoinbaseTransaction() {
		return true
	}

	for _, vin := range tx.Vins {
		if prevTXs[hex.EncodeToString(vin.Txhash)].TxHash == nil {
			log.Panic("ERROR: Previous transaction is not correct")
		}
	}

	txCopy := tx.TrimmedCopy()

	curve := elliptic.P256()

	for inID, vin := range tx.Vins {
		prevTx := prevTXs[hex.EncodeToString(vin.Txhash)]
		txCopy.Vins[inID].ScriptSig = nil
		txCopy.Vins[inID].PublicKey = prevTx.Vouts[vin.Vout].PubKeyHash
		txCopy.TxHash = txCopy.Hash()
		txCopy.Vins[inID].PublicKey = nil

		// 私钥 ID
		r := big.Int{}
		s := big.Int{}
		sigLen := len(vin.ScriptSig)
		r.SetBytes(vin.ScriptSig[:(sigLen / 2)])
		s.SetBytes(vin.ScriptSig[(sigLen / 2):])

		x := big.Int{}
		y := big.Int{}
		keyLen := len(vin.PublicKey)
		x.SetBytes(vin.PublicKey[:(keyLen / 2)])
		y.SetBytes(vin.PublicKey[(keyLen / 2):])

		rawPubKey := ecdsa.PublicKey{curve, &x, &y}
		if ecdsa.Verify(&rawPubKey, txCopy.TxHash, &r, &s) == false {
			return false
		}
	}

	return true
}


func (tx *Transaction) Hash() []byte {

	txCopy := tx

	txCopy.TxHash = []byte{}

	hash := sha256.Sum256(txCopy.Serialize())
	return hash[:]
}

func (tx *Transaction) Serialize() []byte {
	var encoded bytes.Buffer

	enc := gob.NewEncoder(&encoded)
	err := enc.Encode(tx)
	if err != nil {
		log.Panic(err)
	}

	return encoded.Bytes()
}
