package blc

import (
	"github.com/boltdb/bolt"
	"log"
)

type BlockchainIterator struct {
	CurrentHash []byte
	DB  *bolt.DB
}

func (blockchainIterator *BlockchainIterator) Next() *YX_Block {

	var block *YX_Block

	err := blockchainIterator.DB.View(func(tx *bolt.Tx) error{

		b := tx.Bucket([]byte(blockTableName))

		if b != nil {
			currentBloclBytes := b.Get(blockchainIterator.CurrentHash)
			//  获取到当前迭代器里面的currentHash所对应的区块
			block = Deserialization(currentBloclBytes)

			// 更新迭代器里面CurrentHash
			blockchainIterator.CurrentHash = block.YX_PrevBlockHash
		}

		return nil
	})

	if err != nil {
		log.Panic(err)
	}


	return block

}