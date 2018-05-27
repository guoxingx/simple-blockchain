package main

import (
    "fmt"
    "log"
    "bytes"
    "encoding/gob"
    "encoding/hex"
    "crypto/sha256"
)

const subsidy = 26

type Transaction struct {
    ID   []byte
    Vin  []TXInput
    Vout []TXOutput
}

// coinbase 交易
// 即区块的奖励交易
func NewCoinbaseTX(to, data string) *Transaction {
    if data == "" {
        data = fmt.Sprintf("Reward to '%s'", to)
    }

    txin := TXInput{[]byte{}, -1, nil, []byte(data)}
    txout := NewTXOutput(subsidy, to)
    tx := Transaction{nil, []TXInput{txin}, []TXOutput{*txout}}
    tx.SetID()

    return &tx
}

func (tx *Transaction) SetID() {
	var encoded bytes.Buffer
	var hash [32]byte

	enc := gob.NewEncoder(&encoded)
	err := enc.Encode(tx)
	if err != nil {
		log.Panic(err)
	}
	hash = sha256.Sum256(encoded.Bytes())
	tx.ID = hash[:]
}

func NewUTXOTransaction(from, to string, amount int, bc *Blockchain) *Transaction {
    var inputs []TXInput
    var outputs []TXOutput

    // 获取的未花费输出 总额 & UTXOs
    acc, validOutputs := bc.FindSpendableOutputs([]byte(from), amount)

    if acc < amount {
        log.Panic("ERROR: Not enough funds")
    }

    // 花费：将获取到的每一个输出都引用并创建一个新的输入
    for txid, outs := range validOutputs {
        txID, err := hex.DecodeString(txid)
        if err != nil { log.Panic(err) }

        for _, out := range outs {
            input := TXInput{txID, out, nil, []byte(from)}
            inputs = append(inputs, input)
        }
    }

    // 转账 amount 到 to 的输出
    outputs = append(outputs, *NewTXOutput(amount, to))

    // 转账 acc - amount 的 from 的输出，即找零
    if acc > amount {
        outputs = append(outputs, *NewTXOutput(acc - amount, from)) // a change
    }

    tx := Transaction{nil, inputs, outputs}
    tx.SetID()

    return &tx
}

func (tx *Transaction) IsCoinbase() bool {
	return len(tx.Vin) == 1 && len(tx.Vin[0].Txid) == 0 && tx.Vin[0].Vout == -1
}
