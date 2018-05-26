package main

import (
    "bytes"
    "time"
    "encoding/gob"
)

type Block struct {
	Timestamp	  int64
	Data		  []byte
	PrevBlockHash []byte
	Hash		  []byte
    Nonce         int
    Number        int
}

// 获取一个新区块
// @param: data: string: 区块data
// @param: prevBlock: *Block: 上一个区块
// @return: *Block
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

// 获取创世块
// @return: *Block
func NewGenesisBlock() *Block {
    return NewBlock("genesis", nil)
}

// 将一个区块序列化
// @param: b: *Block: 区块
// @return: []byte
func (b *Block) Serialize() []byte {
    var result bytes.Buffer
    // encodind/gob.NewEncoder(w io.Writer)
    encoder := gob.NewEncoder(&result)

    err := encoder.Encode(b)

    return result.Bytes()
}

func DeserializeBlock(d []byte) *Block {
    var block Block

    decoder := gob.NewDecoder(bytes.NewReader(d))
    err := decoder.Decode(&block)

    return &block
}
