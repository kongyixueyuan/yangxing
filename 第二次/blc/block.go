package blc

import (
	"time"
	"strconv"
	"bytes"
	"fmt"
	"crypto/sha256"
)

type Block struct {
	Height int64
	PrevBlockHash []byte
	Data []byte
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



	joint_byte := bytes.Join([][]byte{height_byte,block.PrevBlockHash,block.Data,timestamo_byte2},[]byte{})

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

func NewBlock(height int64,prehash []byte,data string) *Block{
	block := &Block{height,prehash,[]byte(data),time.Now().Unix(),nil,0}
	setHashwithProofOfWork(block)
	return block
}

func CreateGenesisBlock(data string) *Block {
	block := NewBlock(1,[]byte{0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0},data)
	return block
}