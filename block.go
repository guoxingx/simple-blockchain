package main

import (
    "bytes"
    "time"
    "log"
    "encoding/gob"
    "math/big"
    "encoding/binary"

    "github.com/guoxingx/simple-blockchain/common"
)

type BlockNonce [8]byte

func EncodeNonce(i uint64) BlockNonce {
    var n BlockNonce
    binary.BigEndian.PutUint64(n[:], i)
    return n
}

type Header struct {
    ParentHash    common.Hash
    Miner         common.Address
    TxHash        common.Hash
    Number        *big.Int
    Timestamp     *big.Int
    Nonce         BlockNonce
}

type Block struct {
    header        *Header
    transactions  []*Transaction
    hash          common.Hash
}

func (block *Block) ParentHash() common.Hash      { return block.header.ParentHash }
func (block *Block) Miner() common.Address        { return block.header.Miner }
func (block *Block) TxHash() common.Hash          { return block.header.TxHash }
func (block *Block) Number() *big.Int             { return new(big.Int).Set(block.header.Number) }
func (block *Block) Timestamp() *big.Int          { return new(big.Int).Set(block.header.Timestamp) }
func (block *Block) Nonce() uint64                { return binary.BigEndian.Uint64(block.header.Nonce[:]) }
func (block *Block) Transactions() []*Transaction { return block.transactions }
func (block *Block) Hash() common.Hash            { return block.hash }

// 获取一个新区块
// @param: miner: []byte: 挖出区块的矿工
// @param: parent: *Block: 上一个区块
// @param: transactions: []*Transaction: 待写入的交易
// @return: *Block
func NewBlock(miner string, parent *Block, transactions []*Transaction) *Block {
    var parentHash common.Hash
    var blockNumber big.Int
    if parent != nil {
        parentHash = parent.hash
        blockNumber = *new(big.Int).Add(parent.Number(), big.NewInt(1))
    }

    header := &Header{parentHash, common.HexToAddress(miner), common.Hash{}, &blockNumber, big.NewInt(time.Now().Unix()), BlockNonce{}}
    block := &Block{header, transactions, common.Hash{}}

    block.HashTransactions()
    pow := NewProofOfWork(block)
    nonce, hash := pow.Run()
    block.header.Nonce = EncodeNonce(uint64(nonce))
    block.hash.SetBytes(hash)

    return block
}

// 获取创世块
// rewardTx 矿工的奖励交易，不需要引用之前交易。
// @return: *Block
// func NewGenesisBlock(miner common.Address, rewardTx *Transaction) *Block {
func NewGenesisBlock(miner string) *Block {
    return NewBlock(miner, nil, []*Transaction{})
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
func (b *Block) HashTransactions() {
    var transactions [][]byte

    for _, tx := range b.transactions {
        transactions = append(transactions, tx.Serialize())
    }
    mTree := NewMerkleTree(transactions)

    b.hash.SetBytes(mTree.RootNode.Data)
}
