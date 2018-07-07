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

type Block struct {
	Height int64
	PrevBlockHash []byte
	//Data []byte
	Txs []*Transaction
	Timestamp int64
	Hash []byte
	Nonce int64
}

func setHash(block *Block){
	//1.Height to byte
	height_byte := Int2Byte(block.Height)
	//2.timestamp to byte
	/*
	fmt.Println(" = ",block.Timestamp)
	timestamp_byte := Int2Byte(block.Timestamp)
	fmt.Println("timestamp_byte = ",timestamp_byte)
	timestamp_byte1 := strconv.FormatInt(block.Timestamp,2)
	fmt.Println("timestamp_byte1 = ",timestamp_byte1)*/
    timestamo_byte2 := []byte(strconv.FormatInt(block.Timestamp,2))
	//fmt.Println("timestamo_byte2 = ",timestamo_byte2)



	joint_byte := bytes.Join([][]byte{height_byte,block.PrevBlockHash,block.HashTransaction(),timestamo_byte2},[]byte{})

	//生成hash
	hashbyte := sha256.Sum256(joint_byte)
	fmt.Println("hashbyte = ",hashbyte)
	block.Hash = hashbyte[:]
}

func setHashwithProofOfWork(block *Block){
	pow := NewProofOfWork(block)
	hash, nonce := pow.run()

	block.Hash = hash[:]
	block.Nonce = nonce

}

func NewBlock(height int64,prehash []byte,txs []*Transaction) *Block{
	block := &Block{height,prehash,txs,time.Now().Unix(),nil,0}
	setHashwithProofOfWork(block)
	return block
}

func CreateGenesisBlock(txs []*Transaction) *Block {
	block := NewBlock(1,[]byte{0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0},txs)
	return block
}

//serialization
func (block *Block)Serializtion() []byte {
	var result bytes.Buffer
	enc := gob.NewEncoder(&result)
	err := enc.Encode(block)

	if err != nil {
		log.Panic(err)
	}

	return result.Bytes()

}

//deserialization
func Deserialization(blockBytes []byte) *Block {
	var block Block
	decode := gob.NewDecoder(bytes.NewReader(blockBytes))
	err := decode.Decode(&block)

	if err != nil {
		log.Panic(err)
	}

	return &block
}

func (block *Block)HashTransaction() []byte {

	var txHashes [][]byte
	var txHash [32]byte

	for _, tx := range block.Txs {
		txHashes = append(txHashes, tx.TxHash)
	}
	txHash = sha256.Sum256(bytes.Join(txHashes, []byte{}))

	return txHash[:]
}