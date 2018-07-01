package blc

import "fmt"

type Blockchain struct {
	blocks []*Block
} 

func CreateGenesisBlockChainWithBlock() *Blockchain {

	block:=CreateGenesisBlock("Genesis Block")
	fmt.Println(block)
	blockchain := &Blockchain{[]*Block{block}}

	return blockchain
}


func (blc *Blockchain)AddNewBlock(data string){
	preheight := blc.blocks[len(blc.blocks)-1].Height
	prehash := blc.blocks[len(blc.blocks)-1].Hash
	block := NewBlock(preheight + 1,prehash,data)
	blc.blocks = append(blc.blocks,block)
	fmt.Println(blc.blocks[len(blc.blocks)-1])
}

