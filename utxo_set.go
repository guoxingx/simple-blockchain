package main

import (
    "log"
    "encoding/hex"
    "github.com/boltdb/bolt"
)

const utxoBucket = "chainstate"

type UTXOSet struct {
    Blockchain *Blockchain
}

// 重新加载 utxo bucket
func (u UTXOSet) Reindex() {
    db := u.Blockchain.db
    bucketName := []byte(utxoBucket)

    // 重新加载
    err := db.Update(func(tx *bolt.Tx) error {
        err := tx.DeleteBucket(bucketName)
        _, err = tx.CreateBucket(bucketName)
        if err != nil { log.Panic(err) }
        return nil
    })
    if err != nil { log.Panic(err) }

    // 从链中获取所有 utxo
    // blockchain.FindUTXO return map[txID]TXOutputs
    UTXO := u.Blockchain.FindUTXO()

    // 保存 utxo 到 utxoBucket
    err = db.Update(func(tx *bolt.Tx) error {
        b := tx.Bucket(bucketName)

        for txID, out := range UTXO {
            // hex.DecodeString  f func(s string) ([]byte, error)
            key, err := hex.DecodeString(txID)
            if err != nil { log.Panic(err) }

            err = b.Put(key, out.Serialize())
            if err != nil { log.Panic(err) }
        }
        return nil
    })
    if err != nil { log.Panic(err) }
}

// 找到总额大于 amount 的足够的未花费输出
func (u UTXOSet) FindSpendableOutputs(pubKeyHash []byte, amount int) (int, map[string][]int) {
    unspentOutputs := make(map[string][]int)
    accumulated := 0
    db := u.Blockchain.db

    err := db.View(func(tx *bolt.Tx) error {
        b := tx.Bucket([]byte(utxoBucket))
        c := b.Cursor()

        for k, v := c.First(); k != nil; k, v = c.Next() {
            txID := hex.EncodeToString(k)
            outs := DeserializeOutputs(v)

            for outIdx, out := range outs.Outputs {
                if out.IsLockedWithKey(pubKeyHash) && accumulated < amount {
                    accumulated += out.Value
                    unspentOutputs[txID] = append(unspentOutputs[txID], outIdx)
                }
            }
        }
        return nil
    })
    if err != nil { log.Panic(err) }

    return accumulated, unspentOutputs
}

// 找到所有未花费输出
func (u UTXOSet) FindUTXO(pubKeyHash []byte) []TXOutput {
    var UTXOs []TXOutput
    db := u.Blockchain.db

    err := db.View(func(tx *bolt.Tx) error {
        b := tx.Bucket([]byte(utxoBucket))
        c := b.Cursor()

        for k, v := c.First(); k != nil; k, v = c.Next() {
            outs := DeserializeOutputs(v)

            for _, out := range outs.Outputs {
                if out.IsLockedWithKey(pubKeyHash) {
                    UTXOs = append(UTXOs, out)
                }
            }
        }
        return nil
    })
    if err != nil { log.Panic(err) }

    return UTXOs
}

// 当挖出一个新块时，更新 utxo Bucket
// 移除已花费输出，并从新挖出来的交易中加入未花费输出
func (u UTXOSet) Update(block *Block) {
    db := u.Blockchain.db

    err := db.Update(func(tx *bolt.Tx) error {
        b := tx.Bucket([]byte(utxoBucket))

        // 遍历区块中的交易
        for _, tx := range block.Transactions {
            if tx.IsCoinbase() == false {

                // 遍历交易的输入
                for _, vin := range tx.Vin {
                    updatedOut := TXOutputs{}

                    // 当前交易输入的上一笔输出
                    outsBytes := b.Get(vin.Txid)
                    outs := DeserializeOutputs(outsBytes)

                    for outIdx, out := range outs.Outputs {
                        if outIdx != vin.Vout {
                            updatedOut.Outputs = append(updatedOut.Outputs, out)
                        }
                    }

                    if len(updatedOut.Outputs) == 0 {
                        err := b.Delete(vin.Txid)
                        if err != nil { log.Panic(err) }
                    } else {
                        err := b.Put(vin.Txid, updatedOut.Serialize())
                        if err != nil { log.Panic(err) }
                    }
                }
            }

            newOutputs := TXOutputs{}
            for _, out := range tx.Vout {
                newOutputs.Outputs = append(newOutputs.Outputs, out)
            }
            err := b.Put(tx.ID, newOutputs.Serialize())
            if err != nil { log.Panic(err) }
        }
        return nil
    })
    if err != nil { log.Panic(err) }
}
