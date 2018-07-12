package blc

import (
	"math/big"
	"bytes"
	"crypto/sha256"
	"fmt"
)

type ProofOfWork struct {
	block *YX_Block
	target *big.Int
}

const targetbit = 4

func NewProofOfWork(block *YX_Block) *ProofOfWork {
	newdata := big.NewInt(1)
	newdata = newdata.Lsh(newdata,256-targetbit)
	return &ProofOfWork{block,newdata}
}

func (pow *ProofOfWork)run() ([]byte,int64){
	nonce := 0
	var hash [32]byte
	var hashInt big.Int
	for {
		data := pow.prepareData(nonce)
		//fmt.Println(data)
		// 生成hash
		hash = sha256.Sum256(data)
		fmt.Printf("\r%x",hash)

		hashInt.SetBytes(hash[:])
		if pow.target.Cmp(&hashInt) == 1{
			break
		}
		nonce++


	}
	return hash[:],int64(nonce)
}

func (pow *ProofOfWork) prepareData(nonce int) []byte {
	//fmt.Println(nonce)
	// 生成hash
	data := bytes.Join([][]byte{pow.block.YX_PrevBlockHash,
								pow.block.HashTransaction(),
								Int2Byte(int64(pow.block.YX_Height)),
								Int2Byte(pow.block.YX_Timestamp),
								Int2Byte(int64(targetbit)),
								Int2Byte(int64(nonce))},[]byte{})
	//fmt.Println(data)
	return data
}