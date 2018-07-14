package blc

import (
	"github.com/boltdb/bolt"
	"log"
	"fmt"
	"math/big"
	"time"
	"os"
	"encoding/hex"
	"strconv"
	"bytes"
	"crypto/ecdsa"
)

// 数据库名字
const dbName = "blockchain_%s.db"

// 表的名字
const blockTableName = "blocks"

type Blockchain struct {
	Tip []byte //最新的区块的Hash
	DB  *bolt.DB
}

// 迭代器
func (blockchain *Blockchain) Iterator() *BlockchainIterator {

	return &BlockchainIterator{blockchain.Tip, blockchain.DB}
}

// 判断数据库是否存在
func YX_dbExists(dbName string) bool {
	if _, err := os.Stat(dbName); os.IsNotExist(err) {
		return false
	}

	return true
}


// 遍历输出所有区块的信息
func (blc *Blockchain) Printchain() {

	fmt.Println("PrintchainPrintchainPrintchainPrintchain")
	blockchainIterator := blc.Iterator()

	for {
		block := blockchainIterator.Next()

		fmt.Printf("Height：%d\n", block.YX_Height)
		fmt.Printf("PrevBlockHash：%x\n", block.YX_PrevBlockHash)

		fmt.Printf("Timestamp：%s\n", time.Unix(block.YX_Timestamp, 0).Format("2006-01-02 03:04:05 PM"))
		fmt.Printf("Hash：%x\n", block.YX_Hash)
		fmt.Printf("Nonce：%d\n", block.YX_Nonce)
		fmt.Printf("-------------------------------------------\n")
		for  _,tx := range  block.YX_Txs{
			fmt.Println("当前交易的HASH值")
			fmt.Println(hex.EncodeToString(tx.TxHash))
			fmt.Println("已花费的TXO的相关记录")
			for _,in := range tx.Vins {
				fmt.Printf("Hash: %s Vout: %d ScriptSig: %s\n",in.Txhash,in.Vout,in.ScriptSig)
			}
			fmt.Println("未花费的TXO的相关记录")
			for _,out := range tx.Vouts {
				fmt.Printf("Value: %d ScriptPubKey: %s\n",out.Value,out.PubKeyHash)
			}
		}
		fmt.Println()

		var hashInt big.Int
		hashInt.SetBytes(block.YX_PrevBlockHash)

		// Cmp compares x and y and returns:
		//
		//   -1 if x <  y
		//    0 if x == y
		//   +1 if x >  y

		if big.NewInt(0).Cmp(&hashInt) == 0 {
			break;
		}
	}

}



//// 增加区块到区块链里面
func (blc *Blockchain) AddBlockToBlockchain(data string) {

	err := blc.DB.Update(func(tx *bolt.Tx) error {

		//1. 获取表
		b := tx.Bucket([]byte(blockTableName))
		//2. 创建新区块
		if b != nil {

			// ⚠️，先获取最新区块
			blockBytes := b.Get(blc.Tip)
			// 反序列化
			block := Deserialization(blockBytes)

			//3. 将区块序列化并且存储到数据库中
			newBlock := NewBlock(block.YX_Height+1, block.YX_Hash, []*Transaction{})
			err := b.Put(newBlock.YX_Hash, newBlock.Serializtion())
			if err != nil {
				log.Panic(err)
			}
			//4. 更新数据库里面"l"对应的hash
			err = b.Put([]byte("l"), newBlock.YX_Hash)
			if err != nil {
				log.Panic(err)
			}
			//5. 更新blockchain的Tip
			blc.Tip = newBlock.YX_Hash
		}

		return nil
	})

	if err != nil {
		log.Panic(err)
	}
}

//1. 创建带有创世区块的区块链
func CreateGenesisBlockChainWithBlock(addr string,nodeID string) *Blockchain {

	//组装blockchain name
	dbName := fmt.Sprintf(dbName,nodeID)
	// 判断数据库是否存在
	if YX_dbExists(dbName) {
		fmt.Println("创世区块已经存在.......")
		os.Exit(1)
	}


	fmt.Println("正在创建创世区块.......")

	// 创建或者打开数据库
	db, err := bolt.Open(dbName, 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	var genesisHash []byte
	err = db.Update(func(tx *bolt.Tx) error {

		// 创建数据库表
		b, err := tx.CreateBucket([]byte(blockTableName))

		if err != nil {
			log.Panic(err)
		}

		if b != nil {
			// 创建创世区块

			transaction := NewCoinbaseTransaction(addr)

			genesisBlock := CreateGenesisBlock([]*Transaction{transaction})
			seriral := genesisBlock.Serializtion()

			// 将创世区块存储到表中
			err := b.Put(genesisBlock.YX_Hash, seriral)
			if err != nil {
				log.Panic(err)
			}
			fmt.Println(seriral)
			// 存储最新的区块的hash
			err = b.Put([]byte("l"), genesisBlock.YX_Hash)
			if err != nil {
				log.Panic(err)
			}
			genesisHash = genesisBlock.YX_Hash
		}

		return nil
	})
	return &Blockchain{genesisHash, db}
}


// 返回Blockchain对象
func BlockchainObject(nodeID string) *Blockchain {
	//组装blockchain name
	dbName := fmt.Sprintf(dbName,nodeID)
	// 判断数据库是否存在
	if YX_dbExists(dbName) {
		fmt.Println("创世区块不存在.......")
		os.Exit(1)
	}
	db, err := bolt.Open(dbName, 0600, nil)
	if err != nil {
		log.Fatal(err)
	}

	var tip []byte

	err = db.View(func(tx *bolt.Tx) error {

		b := tx.Bucket([]byte(blockTableName))

		if b != nil {
			// 读取最新区块的Hash
			tip = b.Get([]byte("l"))

		}


		return nil
	})

	return &Blockchain{tip,db}
}

func getunUTXO(addr string, txs []*Transaction) []*UTXO {
	hasspentoutputs := make(map[string][]int)
	PubHash := Base58Decode([]byte(addr))
	Get160Hash := PubHash[1:len(PubHash)-4]
	var utxo []*UTXO
	//get has be spented money
	for _,tx := range txs {
		for _, in := range tx.Vins {
			if in.UnLockRipemd160Hash(Get160Hash) {
				key := hex.EncodeToString(in.Txhash)
				//add vout to hasspentoutputs[key]
				hasspentoutputs[key] = append(hasspentoutputs[key],in.Vout)
			}
		}
	}

	//calculate different person money
	for _,tx := range txs {
work:
		for index, out := range tx.Vouts {
			if out.UnLockScriptPubKeyWithAddress(addr) {
				if len(hasspentoutputs) != 0{
					for inhash,spentint := range hasspentoutputs {
						tmphash := hex.EncodeToString(tx.TxHash)
						if inhash == tmphash {
							for _,value := range spentint {
								if value == index {
									continue work
								} else {
									tmp_utox := &UTXO{tx.TxHash,index,out}
									utxo = append(utxo,tmp_utox)
								}
							}
						} else {
							tmp_utox := &UTXO{tx.TxHash,index,out}
							utxo = append(utxo,tmp_utox)
						}
					}
				} else {
					tmp_utox := &UTXO{tx.TxHash,index,out}
					utxo = append(utxo,tmp_utox)
				}
			}
		}
	}
	return utxo
}

func (blockchain *Blockchain) UnUTXOs(address string,txs []*Transaction) []*UTXO {
	//
	//var unUTXOs []*UTXO
	//
	spentTXOutputs := make(map[string][]int)
	//
	////{hash:[0]}
	//
	//for _,tx := range txs {
	//
	//	if tx.IsCoinbaseTransaction() == false {
	//		for _, in := range tx.Vins {
	//			//是否能够解锁
	//			if in.UnLockWithAddress(address) {
	//
	//				key := hex.EncodeToString(in.Txhash)
	//
	//				spentTXOutputs[key] = append(spentTXOutputs[key], in.Vout)
	//			}
	//
	//		}
	//	}
	//}
	//
	//
	//for _,tx := range txs {
	//
	//Work1:
	//	for index,out := range tx.Vouts {
	//
	//		if out.UnLockScriptPubKeyWithAddress(address) {
	//			fmt.Println("看看是否是俊诚...")
	//			fmt.Println(address)
	//
	//			fmt.Println(spentTXOutputs)
	//
	//			if len(spentTXOutputs) == 0 {
	//				utxo := &UTXO{tx.TxHash, index, out}
	//				unUTXOs = append(unUTXOs, utxo)
	//			} else {
	//				for hash,indexArray := range spentTXOutputs {
	//
	//					txHashStr := hex.EncodeToString(tx.TxHash)
	//
	//					if hash == txHashStr {
	//
	//						var isUnSpentUTXO bool
	//
	//						for _,outIndex := range indexArray {
	//
	//							if index == outIndex {
	//								isUnSpentUTXO = true
	//								continue Work1
	//							}
	//
	//							if isUnSpentUTXO == false {
	//								utxo := &UTXO{tx.TxHash, index, out}
	//								unUTXOs = append(unUTXOs, utxo)
	//							}
	//						}
	//					} else {
	//						utxo := &UTXO{tx.TxHash, index, out}
	//						unUTXOs = append(unUTXOs, utxo)
	//					}
	//				}
	//			}
	//
	//		}
	//
	//	}
	//
	//}
	//
	unUTXOs := getunUTXO(address,txs)

	PubHash := Base58Decode([]byte(address))
	Get160Hash := PubHash[1:len(PubHash)-4]

	blockIterator := blockchain.Iterator()

	for {

		block := blockIterator.Next()

		//fmt.Println(block)
		//fmt.Println()

		for i := len(block.YX_Txs) - 1; i >= 0 ; i-- {

			tx := block.YX_Txs[i]
			// txHash
			// Vins
			if tx.IsCoinbaseTransaction() == false {
				for _, in := range tx.Vins {
					//是否能够解锁
					if in.UnLockRipemd160Hash(Get160Hash) {

						key := hex.EncodeToString(in.Txhash)

						spentTXOutputs[key] = append(spentTXOutputs[key], in.Vout)
					}

				}
			}

			// Vouts

		work:
			for index, out := range tx.Vouts {

				if out.UnLockScriptPubKeyWithAddress(address) {

					//fmt.Println(out)
					//fmt.Println(spentTXOutputs)

					//&{2 zhangqiang}
					//map[]

					if spentTXOutputs != nil {

						//map[cea12d33b2e7083221bf3401764fb661fd6c34fab50f5460e77628c42ca0e92b:[0]]

						if len(spentTXOutputs) != 0 {

							var isSpentUTXO bool

							for txHash, indexArray := range spentTXOutputs {

								for _, i := range indexArray {
									if index == i && txHash == hex.EncodeToString(tx.TxHash) {
										isSpentUTXO = true
										continue work
									}
								}
							}

							if isSpentUTXO == false {

								utxo := &UTXO{tx.TxHash, index, out}
								unUTXOs = append(unUTXOs, utxo)

							}
						} else {
							utxo := &UTXO{tx.TxHash, index, out}
							unUTXOs = append(unUTXOs, utxo)
						}

					}
				}

			}

		}

		//fmt.Println(spentTXOutputs)

		var hashInt big.Int
		hashInt.SetBytes(block.YX_PrevBlockHash)

		// Cmp compares x and y and returns:
		//
		//   -1 if x <  y
		//    0 if x == y
		//   +1 if x >  y
		if hashInt.Cmp(big.NewInt(0)) == 0 {
			break;
		}

	}

	return unUTXOs
}

// 转账时查找可用的UTXO
func (blockchain *Blockchain) FindSpendableUTXOS(from string, amount int,txs []*Transaction) (int64, map[string][]int) {
	//1. 现获取所有的UTXO
	utxos := blockchain.UnUTXOs(from,txs)
	spendableUTXO := make(map[string][]int)
	//2. 遍历utxos
	var value int64
	for _, utxo := range utxos {
		value = value + utxo.Output.Value
		hash := hex.EncodeToString(utxo.TxHash)
		spendableUTXO[hash] = append(spendableUTXO[hash], utxo.Index)

		if value >= int64(amount) {
			break
		}
	}
	if value < int64(amount) {
		fmt.Printf("%s's fund is 不足\n", from)
		os.Exit(1)
	}
	return value, spendableUTXO
}

func (blockchain *Blockchain) CreateNewBlockWithTransaction(from []string, to []string, amount []string,nodeID string) {
	////debug for send
	//fmt.Println("-------------CreateNewBlockWithTransaction--------------")
	//fmt.Println(from)
	//fmt.Println(to)
	//fmt.Println(amount)
	//fmt.Println("-------------CreateNewBlockWithTransaction--------------")

	var txs []*Transaction
	utxoSet := &UTXOSet{blockchain}
	for index,address := range from {
		value, _ := strconv.Atoi(amount[index])
		tx := NewSimpleTransaction(address, to[index], int64(value), utxoSet,txs,nodeID)
		txs = append(txs, tx)
		fmt.Println(tx)
	}

	//奖励
	tx := NewCoinbaseTransaction(from[0])
	txs = append(txs,tx)
	// 找到所有我要删除的数据

	//1. 通过相关算法建立Transaction数组
	var block *YX_Block

	blockchain.DB.View(func(tx *bolt.Tx) error {

		b := tx.Bucket([]byte(blockTableName))
		if b != nil {

			hash := b.Get([]byte("l"))

			blockBytes := b.Get(hash)

			block = Deserialization(blockBytes)

		}

		return nil
	})

	// 在建立新区块之前对txs进行签名验证

	//for _,tx := range txs  {
	//
	//	if blockchain.VerifyTransaction(tx) != true {
	//		log.Panic("ERROR: Invalid transaction")
	//	}
	//}

	//2. 建立新的区块
	block = NewBlock( block.YX_Height+1, block.YX_Hash, txs)

	//将新区块存储到数据库
	blockchain.DB.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blockTableName))
		if b != nil {

			b.Put(block.YX_Hash, block.Serializtion())

			b.Put([]byte("l"), block.YX_Hash)

			blockchain.Tip = block.YX_Hash

		}
		return nil
	})

}

func (blc *Blockchain)getBalance(addr string) int64 {
	var value int64
	utxos := blc.UnUTXOs(addr,[]*Transaction{})
	for _,utxo := range utxos {
		value += utxo.Output.Value
	}
	return value
}




func (bc *Blockchain) FindTransaction(ID []byte) (Transaction, error) {
	//log.Printf("yangxing FindTransaction")
	//log.Println(ID)

	bci := bc.Iterator()

	for {
		block := bci.Next()

		for _, tx := range block.YX_Txs {
			//log.Println(ID)
			//log.Println(tx.TxHash)
			//log.Println(tx)
			if bytes.Compare(tx.TxHash, ID) == 0 {

				return *tx, nil
			}
		}

		var hashInt big.Int
		hashInt.SetBytes(block.YX_PrevBlockHash)


		if big.NewInt(0).Cmp(&hashInt) == 0 {
			break;
		}
	}

	return Transaction{},nil
}


// 验证数字签名
func (bc *Blockchain) VerifyTransaction(tx *Transaction) bool {

	//log.Printf("yangxing VerifyTransaction")
	//log.Println(tx)
	prevTXs := make(map[string]Transaction)

	for _, vin := range tx.Vins {
		prevTX, err := bc.FindTransaction(vin.Txhash)
		if err != nil {
			log.Panic(err)
		}
		prevTXs[hex.EncodeToString(prevTX.TxHash)] = prevTX
		//log.Println(prevTX)
		//log.Println(hex.EncodeToString(prevTX.TxHash))
	}
	//log.Println(prevTXs)
	return tx.Verify(prevTXs)
}

func (bclockchain *Blockchain) SignTransaction(tx *Transaction,privKey ecdsa.PrivateKey)  {
	//log.Printf("yangxing SignTransaction")
	//log.Println(tx)
	if tx.IsCoinbaseTransaction() {
		return
	}

	prevTXs := make(map[string]Transaction)

	for _, vin := range tx.Vins {
		prevTX, err := bclockchain.FindTransaction(vin.Txhash)
		if err != nil {
			log.Panic(err)
		}
		prevTXs[hex.EncodeToString(prevTX.TxHash)] = prevTX
		//log.Println(prevTX)
		//log.Println(hex.EncodeToString(prevTX.TxHash))
	}
	//log.Println(prevTXs)
	tx.Sign(privKey, prevTXs)
}

// [string]*TXOutputs
func (blc *Blockchain) FindUTXOMap() map[string]*TXOutputs  {
	fmt.Println("FindUTXOMap")
	blcIterator := blc.Iterator()

	// 存储已花费的UTXO的信息
	spentableUTXOsMap := make(map[string][]*TXInput)


	utxoMaps := make(map[string]*TXOutputs)


	for {
		block := blcIterator.Next()


		for i := len(block.YX_Txs) - 1; i >= 0 ;i-- {

			txOutputs := &TXOutputs{[]*UTXO{}}

			tx := block.YX_Txs[i]


			// coinbase
			if tx.IsCoinbaseTransaction() == false {
				for _,txInput := range tx.Vins {

					txHash := hex.EncodeToString(txInput.Txhash)
					spentableUTXOsMap[txHash] = append(spentableUTXOsMap[txHash],txInput)

				}
			}



			txHash := hex.EncodeToString(tx.TxHash)

		WorkOutLoop:
			for index,out := range tx.Vouts  {

				if tx.IsCoinbaseTransaction() {

					fmt.Println("IsCoinbaseTransaction")
					fmt.Println(out)
					fmt.Println(txHash)
					fmt.Println("IsCoinbaseTransaction finish")
				}

				txInputs := spentableUTXOsMap[txHash]
				fmt.Println(len(txInputs))
				if len(txInputs) > 0 {

					isSpent := false

					for _,in := range  txInputs {

						outPublicKey := out.PubKeyHash
						inPublicKey := in.PublicKey

						if bytes.Compare(outPublicKey,Ripemd160Hash(inPublicKey)) == 0{
							if index == in.Vout {
								fmt.Println("isSpent == true")
								isSpent = true
								continue WorkOutLoop
							}
						}

					}

					if isSpent == false {
						fmt.Println("isSpent == false")
						fmt.Println(out)
						fmt.Println(txHash)
						fmt.Println("isSpent == false finish")
						utxo := &UTXO{tx.TxHash,index,out}
						txOutputs.UTXOS = append(txOutputs.UTXOS,utxo)
					}

				} else {
					utxo := &UTXO{tx.TxHash,index,out}
					txOutputs.UTXOS = append(txOutputs.UTXOS,utxo)
				}

			}

			// 设置键值对
			utxoMaps[txHash] = txOutputs

		}


		// 找到创世区块时退出
		var hashInt big.Int
		hashInt.SetBytes(block.YX_PrevBlockHash)

		if hashInt.Cmp(big.NewInt(0)) == 0 {
			break;
		}



	}

	return utxoMaps
}


//----------

func (bc *Blockchain) YX_GetBestHeight() int64 {

	block := bc.Iterator().Next()

	return block.YX_Height
}

func (bc *Blockchain) YX_GetBlockHashes() [][]byte {

	blockIterator := bc.Iterator()

	var blockHashs [][]byte

	for {
		block := blockIterator.Next()

		blockHashs = append(blockHashs,block.YX_Hash)

		var hashInt big.Int
		hashInt.SetBytes(block.YX_PrevBlockHash)

		if hashInt.Cmp(big.NewInt(0)) == 0 {
			break;
		}
	}

	return blockHashs
}