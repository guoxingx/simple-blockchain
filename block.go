package main

import (
    "time"
)


type Block struct {
	Timestamp	  int64
	Data		  []byte
	PrevBlockHash []byte
	Hash		  []byte
    None          int
}

func NewBlock(data string, prevBlockHash []byte) *Block {
    block := &Block{time.Now().Unix(), []byte(data), prevBlockHash, []byte{}, 0}
    pow := NewProofOfWork(block)

    nonce, hash := pow.Run()
    block.Hash = hash[:]
    block.None = nonce
    return block
}

func NewGenesisBlock() *Block {
    return NewBlock("genesis", []byte{})
}
