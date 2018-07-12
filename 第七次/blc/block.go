package blc

import (
	"time"
	"strconv"
	"bytes"
	"fmt"
	"crypto/sha256"
	"encoding/gob"
	"log"
)

type YX_Block struct {
	YX_Height int64
	YX_PrevBlockHash []byte
	//Data []byte
	YX_Txs []*Transaction
	YX_Timestamp int64
	YX_Hash []byte
	YX_Nonce int64
}

func setHash(block *YX_Block){
	//1.Height to byte
	height_byte := Int2Byte(block.YX_Height)
	//2.timestamp to byte
	/*
	fmt.Println(" = ",block.Timestamp)
	timestamp_byte := Int2Byte(block.Timestamp)
	fmt.Println("timestamp_byte = ",timestamp_byte)
	timestamp_byte1 := strconv.FormatInt(block.Timestamp,2)
	fmt.Println("timestamp_byte1 = ",timestamp_byte1)*/
    timestamo_byte2 := []byte(strconv.FormatInt(block.YX_Timestamp,2))
	//fmt.Println("timestamo_byte2 = ",timestamo_byte2)



	joint_byte := bytes.Join([][]byte{height_byte,block.YX_PrevBlockHash,block.HashTransaction(),timestamo_byte2},[]byte{})

	//生成hash
	hashbyte := sha256.Sum256(joint_byte)
	fmt.Println("hashbyte = ",hashbyte)
	block.YX_Hash = hashbyte[:]
}

func setHashwithProofOfWork(block *YX_Block){
	pow := NewProofOfWork(block)
	hash, nonce := pow.run()

	block.YX_Hash = hash[:]
	block.YX_Nonce = nonce

}

func NewBlock(height int64,prehash []byte,txs []*Transaction) *YX_Block{
	block := &YX_Block{height,prehash,txs,time.Now().Unix(),nil,0}
	setHashwithProofOfWork(block)
	return block
}

func CreateGenesisBlock(txs []*Transaction) *YX_Block {
	block := NewBlock(1,[]byte{0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0},txs)
	return block
}

//serialization
func (block *YX_Block)Serializtion() []byte {
	var result bytes.Buffer
	enc := gob.NewEncoder(&result)
	err := enc.Encode(block)

	if err != nil {
		log.Panic(err)
	}

	return result.Bytes()

}

//deserialization
func Deserialization(blockBytes []byte) *YX_Block {
	var block YX_Block
	decode := gob.NewDecoder(bytes.NewReader(blockBytes))
	err := decode.Decode(&block)

	if err != nil {
		log.Panic(err)
	}

	return &block
}

func (block *YX_Block)HashTransaction() []byte {

	var transactions [][]byte

	for _, tx := range block.YX_Txs {
		transactions = append(transactions, tx.Serialize())
	}
	mTree := NewMerkleTree(transactions)

	return mTree.RootNode.Data
}