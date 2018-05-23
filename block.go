package main

import (
    "time"
)


type Block struct {
	Timestamp	  int64
	Data		  []byte
	PrevBlockHash []byte
	Hash		  []byte
    Nonce         int
    Number        int
}

func NewBlock(data string, prevBlock *Block) *Block {
    prevBlockHash := []byte{}
    blockNumber := 0
    if prevBlock != nil {
        prevBlockHash = prevBlock.Hash
        blockNumber = prevBlock.Number + 1
    }

    block := &Block{ time.Now().Unix(), []byte(data), prevBlockHash, []byte{}, 0, blockNumber }

    if block.Number == 0 { return block }
    pow := NewProofOfWork(block)

    nonce, hash := pow.Run()
    block.Hash = hash[:]
    block.Nonce = nonce
    return block
}

func NewGenesisBlock() *Block {
    return NewBlock("genesis", nil)
}
