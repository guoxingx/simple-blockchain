package main

import (
    "bytes"
    "time"
    "log"
    "encoding/gob"
    "crypto/sha256"
)

type Block struct {
	Timestamp	  int64
    Transactions  []*Transaction
	PrevBlockHash []byte
	Hash		  []byte
    Nonce         int
    Number        int
}

// 获取一个新区块
// @param: transactions: []*Transaction: 待写入的交易
// @param: prevBlock: *Block: 上一个区块
// @return: *Block
func NewBlock(transactions []*Transaction, prevBlock *Block) *Block {
    prevBlockHash := []byte{}
    blockNumber := 0
    if prevBlock != nil {
        prevBlockHash = prevBlock.Hash
        blockNumber = prevBlock.Number + 1
    }

    block := &Block{ time.Now().Unix(), transactions, prevBlockHash, []byte{}, 0, blockNumber }

    pow := NewProofOfWork(block)
    nonce, hash := pow.Run()

    block.Nonce = nonce
    block.Hash = hash[:]

    return block
}

// 获取创世块
// coinbase 矿工的奖励交易，不需要引用之前交易。
// 与以太坊的coinbase账户意义不同
// @return: *Block
func NewGenesisBlock(coinbase *Transaction) *Block {
    return NewBlock([]*Transaction{coinbase}, nil)
}

// 将一个区块序列化
// @param: b: *Block: 区块
// @return: []byte
func (b *Block) Serialize() []byte {
    var result bytes.Buffer
    // encodind/gob.NewEncoder(w io.Writer)
    encoder := gob.NewEncoder(&result)

    err := encoder.Encode(b)
    if err != nil { log.Panic(err) }

    return result.Bytes()
}

func DeserializeBlock(d []byte) *Block {
    var block Block

    decoder := gob.NewDecoder(bytes.NewReader(d))
    err := decoder.Decode(&block)
    if err != nil { log.Panic(err) }

    return &block
}

// 一个区块所有交易的hash
// 暂时没有使用merkle tree
func (b *Block) HashTransactions() []byte {
    var txHashes [][]byte // list of tx.ID
    var txHash [32]byte

    for _, tx := range b.Transactions {
        txHashes = append(txHashes, tx.ID)
    }
    txHash = sha256.Sum256(bytes.Join(txHashes, []byte{}))

    return txHash[:]
}
