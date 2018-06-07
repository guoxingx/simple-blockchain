package main

import (
    "os"
    "fmt"
    "log"
    "bytes"
    "errors"
    "crypto/ecdsa"
    "encoding/hex"

    "github.com/boltdb/bolt"
    "github.com/guoxingx/simple-blockchain/common"
)

const dbFile = "data/chain.db"
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
        rewardTx := NewRewardTx(address, genesisCoinbaseData)
        genesis := NewGenesisBlock(address, rewardTx)

        b, err := tx.CreateBucket([]byte(blocksBucket))
        if err != nil { log.Panic(err) }

        err = b.Put(genesis.Hash.Bytes(), genesis.Serialize())
        if err != nil { log.Panic(err) }

        err = b.Put([]byte(latestBlockName), genesis.Hash.Bytes())
        if err != nil { log.Panic(err) }

        tip = genesis.Hash.Bytes()

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
func (bc *Blockchain) MineBlock(miner string, transactions []*Transaction) *Block {
    var lastEncodedBlock []byte

    // 校验将被写入区块的所有交易
    for _, tx := range transactions {
        if bc.VerifyTransaction(tx) != true {
            log.Panic("ERROR: Invalid transaction")
        }
    }

    err := bc.db.View(func(tx *bolt.Tx) error {
        b := tx.Bucket([]byte(blocksBucket))
        lastHash := b.Get([]byte(latestBlockName))
        lastEncodedBlock = b.Get(lastHash)

        return nil
    })
    if err != nil { log.Panic(err) }

    // load last block by lastHash
    transactions = append([]*Transaction{NewRewardTx(miner, "")}, transactions...)
    newBlock := NewBlock(miner, DeserializeBlock(lastEncodedBlock), transactions)
    // transactions = append(transactions, NewRewardTx(miner, ""))

    err = bc.db.Update(func(tx *bolt.Tx) error {
        b := tx.Bucket([]byte(blocksBucket))

        err := b.Put(newBlock.Hash.Bytes(), newBlock.Serialize())
        if err != nil { log.Panic(err) }

        err = b.Put([]byte(latestBlockName), newBlock.Hash.Bytes())
        if err != nil { log.Panic(err) }

        bc.tip = newBlock.Hash.Bytes()

        return nil
    })
    if err != nil { log.Panic(err) }

    return newBlock
}

func (bc *Blockchain) Iterator() *BlockchainIterator {
    bci := &BlockchainIterator{bc.tip, bc.db}

    return bci
}

// 找到包含未花费输出的交易
func (bc *Blockchain) FindUnspentTransactions() []Transaction {
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
            for outIdx, _ := range tx.Vout {
                // 如果输出已经被包含在某个输入内 即已被花费 则跳过
                if spentTXOs[txID] != nil {
                    for _, spentOut := range spentTXOs[txID] {
                        if spentOut == outIdx {
                            continue Outputs
                        }
                    }
                }
                unspentTXs = append(unspentTXs, *tx)
            }

            // 遍历交易的所有输入
            if tx.IsCoinbase() == false {
                for _, in := range tx.Vin {
                    inTxID := hex.EncodeToString(in.Txid)
                    spentTXOs[inTxID] = append(spentTXOs[inTxID], in.Vout)
                }
            }
        }

        // 循环至创世块
        if (block.ParentHash() == common.Hash{}) { break }
    }

    return unspentTXs
}

// 找到所有未花费的输出
// func (bc *Blockchain) FindUTXO(pubKeyHash []byte) []TXOutput {
func (bc *Blockchain) FindUTXO() map[string]TXOutputs {
    var UTXO = make(map[string]TXOutputs)
    unspentTransactions := bc.FindUnspentTransactions()

    for _, tx := range unspentTransactions {
        txID := hex.EncodeToString(tx.ID)
        outs := UTXO[txID]
        outs.Outputs = append(outs.Outputs, tx.Vout...)
        UTXO[txID] = outs
    }

    return UTXO
}

// 根据 tx.ID 找到交易
func (bc *Blockchain) FindTransaction(ID []byte) (Transaction, error) {
    bci := bc.Iterator()

    for {
        block := bci.Next()

        for _, tx := range block.Transactions {
            if bytes.Compare(tx.ID, ID) == 0 { return *tx, nil }
        }

        if len(block.ParentHash()) == 0 { break }
    }

    return Transaction{}, errors.New("Transaction is not found")
}

// 交易签名
func (bc *Blockchain) SignTransaction(tx *Transaction, privKey ecdsa.PrivateKey) {
    prevTXs := make(map[string]Transaction)

    for _, vin := range tx.Vin {
        prevTX, err := bc.FindTransaction(vin.Txid)
        if err != nil { log.Panic(err) }

        prevTXs[hex.EncodeToString(prevTX.ID)] = prevTX
    }

    tx.Sign(privKey, prevTXs)
}

// 验证交易
func (bc *Blockchain) VerifyTransaction(tx *Transaction) bool {
    if tx.IsCoinbase() { return true }

    prevTXs := make(map[string]Transaction)

    for _, vin := range tx.Vin {
        prevTX, err := bc.FindTransaction(vin.Txid)
        if err != nil { log.Panic(err) }

        prevTXs[hex.EncodeToString(prevTX.ID)] = prevTX
    }
    return tx.Verify(prevTXs)
}
