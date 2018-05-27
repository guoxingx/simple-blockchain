package main

import (
    "fmt"
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

fVunc NewTx(from, to, data string) *Transaction {
    return nil
}

func (tx *Transaction) SetID() {

}

func NewUTXOTransaction(from, to string, amount int, bc *Blockchain) *Transaction {
    var inputs []TXInput
    var outputs []TXOutput

    acc, validOutputs := bc.FindSpendableOutputs(from, amount)

    if acc < amount {
        log.Panic("ERROR: Not enough funds")
    }

    // Build a list of inputs
    for txid, outs := range validOutputs {
        txID, err := hex.DecodeString(txid)

        for _, out := range outs {
            input := TXInput{txID, out, from}
            inputs = append(inputs, input)
        }
    }

    // Build a list of outputs
    outputs = append(outputs, TXOutput{amount, to})
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
