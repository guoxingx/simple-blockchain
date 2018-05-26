package main

import (
    "fmt"
    "log"
    "github.com/boltdb/bolt"
)

const dbFile = "blockchain.db"
const blocksBucket = "blocks" // means database.

type Blockchain struct {
    tip []byte
    db  *bolt.DB
	// blocks []*Block
}

/*
打开一个数据库文件,检查文件里面是否已经存储了一个区块链
1. 如果已经存储了一个区块链：
    创建一个新的 Blockchain 实例
    设置 Blockchain 实例的 tip 为数据库中存储的最后一个块的哈希
2. 如果没有区块链：
    创建创世块
    存储到数据库
    将创世块哈希保存为最后一个块的哈希
    创建一个新的 Blockchain 实例，其 tip 指向创世块（tip 有尾部，尖端的意思，在这里 tip 存储的是最后一个块的哈希）
*/
func NewBlockchain() *Blockchain {
    var tip []byte
    db, err := bolt.Open(dbFile, 0600, nil)
    if err != nil { log.Panic(err) }

    err = db.Update(func(tx *bolt.Tx) error {
        // bucket: database.
        b := tx.Bucket([]byte(blocksBucket))

        if b == nil {
            fmt.Println("no exising blockchain founded. a new one will be created.")

            genesis := NewGenesisBlock()
            b, err := tx.CreateBucket([]byte(blocksBucket))
            if err != nil { log.Panic(err) }

            err = b.Put(genesis.Hash, genesis.Serialize())
            if err != nil { log.Panic(err) }

            err = b.Put([]byte("latest"), genesis.Hash)
            if err != nil { log.Panic(err) }

            tip = genesis.Hash
        } else {
            tip = b.Get([]byte("latest"))
        }

        return nil
    })
    if err != nil { log.Panic(err) }

    bc := Blockchain{tip, db}

    return &bc
}

// 添加一个区块
func (bc *Blockchain) AddBlock(data string) {
    var last Hash []byte

    err := bc.db.View(func(tx *bolt.Tx) error {
        b := tx.Bucket([]byte(blocksBucket))
        lastHash = b.Get([]byte("l"))
    })
    if err != nil { log.Panic(err) }

    // load last block by lastHash
    newBlock := NewBlock(data, b.Get(lastHash))

    err = bc.db.Update(func(tx *bolt.Tx) error {
        b := tx.Bucket([]byte(blocksBucket))
        err := b.Put(newBlock.Hash, newBlock.Serialize())
        if err != nil { log.Panic(err) }

        err = b.Put([]byte("latest"), newBlock.Hash)
        if err != nil { log.Panic(err) }

        bc.tip = newBlock.Hash

        return nil
    })
    if err != nil { log.Panic(err) }
}

// 区块链迭代
type BlockchainIterator struct {
    currentHash []byte
    db          *bolt.Db
}

func (bc *Blockchain) Iterator() *BlockchainIterator {
    bci := &BlockchainIterator{bc.tip, bc.db}

    return bci
}

// 其实是查找上一个区块
func (i *BlockchainIterator) Next() *Block {
    var block *Block

    err := i.db.View(func(tx *bolt.Tx) error {
        b := tx.Bucket([]byte(blocksBucket))
        encodedBlock := b.Get(i.currentHash)
        block = DeserializeBlock(encodedBlock)

        return nil
    })
    if err != nil { log.Panic(err) }


    i.currentHash = block.PrevBlockHash

    return block
}

