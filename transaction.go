package main

import (
    "fmt"
    "log"
    "encoding/hex"
)

type Transaction struct {
    ID   []byte
    Vin  []TXInput
    Vout []TXOutput
}

type TXOutput struct {
    Value        int
    ScriptPubKey string
}

type TXInput struct {
    Txid      []byte
    Vout      int
    ScriptSig string
}

const subsidy = 26

// coinbase 交易
// 即区块的奖励交易
func NewCoinbaseTX(to, data string) *Transaction {
    if data == "" {
        data = fmt.Sprintf("Reward to '%s'", to)
    }

    txin := TXInput{[]byte{}, -1, data}
    txout := TXOutput{subsidy, to}
    tx := Transaction{nil, []TXInput{txin}, []TXOutput{txout}}
    tx.SetID()

    return &tx
}

func NewTx(from, to, data string) *Transaction {
    return nil
}

func (tx *Transaction) SetID() {

}

func NewUTXOTransaction(from, to string, amount int, bc *Blockchain) *Transaction {
    var inputs []TXInput
    var outputs []TXOutput

    // 获取的未花费输出 总额 & UTXOs
    acc, validOutputs := bc.FindSpendableOutputs(from, amount)

    if acc < amount {
        log.Panic("ERROR: Not enough funds")
    }

    // 花费：将获取到的每一个输出都引用并创建一个新的输入
    for txid, outs := range validOutputs {
        txID, err := hex.DecodeString(txid)
        if err != nil { log.Panic(err) }

        for _, out := range outs {
            input := TXInput{txID, out, from}
            inputs = append(inputs, input)
        }
    }

    // 转账 amount 到 to 的输出
    outputs = append(outputs, TXOutput{amount, to})

    // 转账 acc - amount 的 from 的输出，即找零
    if acc > amount {
        outputs = append(outputs, TXOutput{acc - amount, from}) // a change
    }

    tx := Transaction{nil, inputs, outputs}
    tx.SetID()

    return &tx
}

func (tx *Transaction) IsCoinbase() bool {
    return false
}

func (in *TXInput) CanUnlockOutputWith(unlockingData string) bool {
    return in.ScriptSig == unlockingData
}

func (out *TXOutput) CanBeUnlockedWith(unlockingData string) bool {
    return out.ScriptPubKey == unlockingData
}
