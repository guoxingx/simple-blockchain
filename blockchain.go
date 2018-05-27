package main

import (
    "os"
    "fmt"
    "log"
    "encoding/hex"
    "github.com/boltdb/bolt"
)

const dbFile = "chain.db"
const blocksBucket = "blocks" // means database.
const latestBlockName = "latest"
const genesisCoinbaseData = "Do not go gentle into that good night"

type Blockchain struct {
    tip []byte
    db  *bolt.DB
	// blocks []*Block
}

/*
创建一个新的 Blockchain 实例
设置 Blockchain 实例的 tip 为数据库中存储的最后一个块的哈希
*/
func NewBlockchain() *Blockchain {
    if dbExists() == false {
        fmt.Println("No existing blockchain found. Create one first.")
        os.Exit(1)
    }
    var tip []byte
    db, err := bolt.Open(dbFile, 0600, nil)
    if err != nil { log.Panic(err) }

	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		tip = b.Get([]byte(latestBlockName))

		return nil
	})
    if err != nil { log.Panic(err) }

    bc := Blockchain{tip, db}
    return &bc
}

/*
创建创世块
把奖励交易发送到指定address
存储到数据库
将创世块哈希保存为最后一个块的哈希
创建一个新的 Blockchain 实例，其 tip 指向创世块（tip 有尾部，尖端的意思，在这里 tip 存储的是最后一个块的哈希）
*/
func CreateBlockchain(address string) *Blockchain {
    if dbExists() {
        fmt.Println("Blockchain already exists.")
        os.Exit(1)
    }

    var tip []byte
    db, err := bolt.Open(dbFile, 0600, nil)
    if err != nil { log.Panic(err) }

    err = db.Update(func(tx *bolt.Tx) error {
        coinbase := NewCoinbaseTX(address, genesisCoinbaseData)
        genesis := NewGenesisBlock(coinbase)

        b, err := tx.CreateBucket([]byte(blocksBucket))
        if err != nil { log.Panic(err) }

        err = b.Put(genesis.Hash, genesis.Serialize())
        if err != nil { log.Panic(err) }

        err = b.Put([]byte(latestBlockName), genesis.Hash)
        if err != nil { log.Panic(err) }

        tip = genesis.Hash

        return nil
    })
    if err != nil { log.Panic(err) }

    bc := Blockchain{tip, db}
    return &bc
}

// 判断数据库是否已经存在
func dbExists() bool {
    // os.IsNotExist f func(err error) bool
    if _, err := os.Stat(dbFile); os.IsNotExist(err) { return false }
    return true
}

// 添加一个区块
// func (bc *Blockchain) AddBlock(data string) {
func (bc *Blockchain) MineBlock(transactions []*Transaction) {
    var lastEncodedBlock []byte

    err := bc.db.View(func(tx *bolt.Tx) error {
        b := tx.Bucket([]byte(blocksBucket))
        lastHash := b.Get([]byte(latestBlockName))
        lastEncodedBlock = b.Get(lastHash)

        return nil
    })
    if err != nil { log.Panic(err) }

    // load last block by lastHash
    newBlock := NewBlock(transactions, DeserializeBlock(lastEncodedBlock))

    err = bc.db.Update(func(tx *bolt.Tx) error {
        b := tx.Bucket([]byte(blocksBucket))
        err := b.Put(newBlock.Hash, newBlock.Serialize())
        if err != nil { log.Panic(err) }

        err = b.Put([]byte(latestBlockName), newBlock.Hash)
        if err != nil { log.Panic(err) }

        bc.tip = newBlock.Hash

        return nil
    })
    if err != nil { log.Panic(err) }
}

func (bc *Blockchain) Iterator() *BlockchainIterator {
    bci := &BlockchainIterator{bc.tip, bc.db}

    return bci
}

// 找到包含未花费输出的交易
func (bc *Blockchain) FindUnspentTransactions(pubKeyHash []byte) []Transaction {
    var unspentTXs []Transaction
    spentTXOs := make(map[string][]int)
    bci := bc.Iterator()

    for {
        block := bci.Next()

        // 遍历区块中全部交易
        for _, tx := range block.Transactions {
            // hex.EncodeToString f func(src []byte) string
            txID := hex.EncodeToString(tx.ID)

        // 多层循环 continue作用于指定的循环
        Outputs:
            // 遍历交易的所有输出
            // 因为区块是从最新往前遍历的
            // 所以可以先检查输出，再检查输入
            for outIdx, out := range tx.Vout {
                // 如果输出已经被包含在某个输入内 即已被花费 则跳过
                if spentTXOs[txID] != nil {
                    for _, spentOut := range spentTXOs[txID] {
                        if spentOut == outIdx {
                            continue Outputs
                        }
                    }
                }

                // 如果该输出可以被解锁，即可被花费
                if out.IsLockedWithKey(pubKeyHash) {
                    unspentTXs = append(unspentTXs, *tx)
                }
            }

            // 遍历交易的所有输入
            if tx.IsCoinbase() == false {
                for _, in := range tx.Vin {
                    // 如果该输入可以被解锁，记录该输入
                    if in.UsesKey(pubKeyHash) {
                        inTxID := hex.EncodeToString(in.Txid)
                        spentTXOs[inTxID] = append(spentTXOs[inTxID], in.Vout)
                    }
                }
            }
        }

        // 循环至创世块
        if len(block.PrevBlockHash) == 0 { break }
    }

    return unspentTXs
}

// 找到所有未花费的输出
func (bc *Blockchain) FindUTXO(pubKeyHash []byte) []TXOutput {
    var UTXOs []TXOutput
    unspentTransactions := bc.FindUnspentTransactions(pubKeyHash)

    for _, tx := range unspentTransactions {
        for _, out := range tx.Vout {
            if out.IsLockedWithKey(pubKeyHash) {
                UTXOs = append(UTXOs, out)
            }
        }
    }

    return UTXOs
}

// 找到总额大于amount的足够的未花费输出
func (bc *Blockchain) FindSpendableOutputs(pubKeyHash []byte, amount int) (int, map[string][]int) {
    unspentOutputs := make(map[string][]int)
    unspentTXs := bc.FindUnspentTransactions(pubKeyHash)
    accumulated := 0

    Work:
        // 遍历未花费输出，直至总额大于 amount
        for _, tx := range unspentTXs {
            txID := hex.EncodeToString(tx.ID)

            for outIdx, out := range tx.Vout {
                if out.IsLockedWithKey(pubKeyHash) && accumulated < amount {
                    accumulated += out.Value
                    unspentOutputs[txID] = append(unspentOutputs[txID], outIdx)

                    if accumulated >= amount { break Work }
                }
            }
        }
    return accumulated, unspentOutputs
}
