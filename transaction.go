package main

import (
    "fmt"
    "log"
    "bytes"
    "time"
    "math/big"
    "encoding/gob"
    "encoding/hex"
    "encoding/binary"
    "crypto/sha256"
    "crypto/ecdsa"
    "crypto/rand"
    "crypto/elliptic"
)

const subsidy = 26

type Transaction struct {
    ID   []byte
    Vin  []TXInput
    Vout []TXOutput
}

// 即区块的奖励交易
func NewRewardTx(to, data string) *Transaction {
    // 奖励交易没有输入 也不会被校验
    // 因此 TXInput.Signature = nil, TXInput.PubKey 随机生成
    // 根据 当前时间 和 随机数 生成 PubKey
    if data == "" {
        var ts bytes.Buffer
        binary.Write(&ts, binary.BigEndian, time.Now().UnixNano())

        randData := make([]byte, 20)
        _, err := rand.Read(randData)
        if err != nil { log.Panic(err) }

        randData = append(ts.Bytes(), randData...)
        data = fmt.Sprintf("%v", randData)
    }

    txin := TXInput{[]byte{}, -1, nil, []byte(data)}

    txout := NewTXOutput(subsidy, to)
    tx := Transaction{nil, []TXInput{txin}, []TXOutput{*txout}}
    tx.ID = tx.Hash()

    return &tx
}

// 发起交易
func NewUTXOTransaction(from, to string, amount int, UTXOSet *UTXOSet) *Transaction {
    var inputs []TXInput
    var outputs []TXOutput

    // 获取的未花费输出 总额 & UTXOs
    wallets, err := NewWallets()
    if err != nil { log.Panic(err) }

    // 从wallet获取address对应的pubKeyHash
    wallet := wallets.GetWallet(from)
    pubKeyHash := HashPubKey(wallet.PublicKey)

    acc, validOutputs := UTXOSet.FindSpendableOutputs(pubKeyHash, amount)

    if acc < amount {
        log.Panic("ERROR: Not enough funds")
    }

    // 花费：将获取到的每一个输出都引用并创建一个新的输入
    for txid, outs := range validOutputs {
        txID, err := hex.DecodeString(txid)
        if err != nil { log.Panic(err) }

        for _, out := range outs {
            input := TXInput{txID, out, nil, wallet.PublicKey}
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

    // 交易签名
    tx.ID = tx.Hash()
    UTXOSet.Blockchain.SignTransaction(&tx, wallet.PrivateKey)

    return &tx
}

// if the transaction is rewared to miner.
func (tx *Transaction) IsCoinbase() bool {
	return len(tx.Vin) == 1 && len(tx.Vin[0].Txid) == 0 && tx.Vin[0].Vout == -1
}

// 签名
// 一个私钥和一个之前交易的 map
func (tx *Transaction) Sign(privKey ecdsa.PrivateKey, prevTXs map[string]Transaction) {
    if tx.IsCoinbase() {
        return
    }

    txCopy := tx.TrimmedCopy()

    // 遍历输入
    for inID, vin := range txCopy.Vin {

        // 获取 当前交易输入 对应的上一笔交易
        prevTx := prevTXs[hex.EncodeToString(vin.Txid)]

        // 仅仅是一个双重检验
        txCopy.Vin[inID].Signature = nil
        txCopy.Vin[inID].PubKey = prevTx.Vout[vin.Vout].PubKeyHash
        txCopy.ID = txCopy.Hash()

        // 重置 PubKey 不影响后面的遍历
        txCopy.Vin[inID].PubKey = nil

        // 用privKey 对 txCopy.ID 进行签名
        // ecdsa.Sign f func(rand io.Reader, priv *ecdsa.PrivateKey, hash []byte) (r *big.Int, s *big.Int, err error)
        r, s, err := ecdsa.Sign(rand.Reader, &privKey, txCopy.ID)
        if err != nil { log.Panic(err) }

        signature := append(r.Bytes(), s.Bytes()...)

        tx.Vin[inID].Signature = signature
    }
}

// Verify verifies signatures of Transaction inputs
func (tx *Transaction) Verify(prevTXs map[string]Transaction) bool {
    txCopy := tx.TrimmedCopy()

    for inID, vin := range tx.Vin {
        prevTx := prevTXs[hex.EncodeToString(vin.Txid)]
        txCopy.Vin[inID].Signature = nil
        txCopy.Vin[inID].PubKey = prevTx.Vout[vin.Vout].PubKeyHash
        txCopy.ID = txCopy.Hash()
        txCopy.Vin[inID].PubKey = nil

        var r, s big.Int
        sigLen := len(vin.Signature)

        // big.Int.SetBytes f func(buf []byte) *big.Int
        r.SetBytes(vin.Signature[:(sigLen / 2)])
        s.SetBytes(vin.Signature[(sigLen / 2):])

        var x, y big.Int
        keyLen := len(vin.PubKey)
        x.SetBytes(vin.PubKey[:(keyLen / 2)])
        y.SetBytes(vin.PubKey[(keyLen / 2):])

        // elliptic.P256  f func() elliptic.Curve
        rawPubKey := ecdsa.PublicKey{elliptic.P256(), &x, &y}
        if ecdsa.Verify(&rawPubKey, txCopy.ID, &r, &s) == false {
            return false
        }
    }
    return true
}

// 返回 input.Signature 和 input.PubKey 被设置为nil 的交易副本
// 因为不需要对存储在输入里面的公钥签名
// 即不需要校验当前交易输入的发送方
func (tx *Transaction) TrimmedCopy() Transaction {
    var inputs []TXInput
    var outputs []TXOutput

    for _, vin := range tx.Vin {
        inputs = append(inputs, TXInput{vin.Txid, vin.Vout, nil, nil})
    }

    for _, vout := range tx.Vout {
        outputs = append(outputs, TXOutput{vout.Value, vout.PubKeyHash})
    }

    txCopy := Transaction{tx.ID, inputs, outputs}
    return txCopy
}

// Hash returns the hash of the Transaction
func (tx *Transaction) Hash() []byte {
    var hash [32]byte

    txCopy := *tx
    txCopy.ID = []byte{}

    hash = sha256.Sum256(txCopy.Serialize())

    return hash[:]
}

// Serialize returns a serialized Transaction
func (tx Transaction) Serialize() []byte {
    var encoded bytes.Buffer

    // gob.NewEncoder f func(w io.Writer) *gob.Encoder
    enc := gob.NewEncoder(&encoded)
    err := enc.Encode(tx)
    if err != nil { log.Panic(err) }

    return encoded.Bytes()
}
