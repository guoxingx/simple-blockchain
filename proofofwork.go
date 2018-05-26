package main

import (
    "math"
    "math/big"
    "bytes"
    "crypto/sha256"
    "fmt"
)

const targetBits = 22

type ProofOfWork struct {
    block *Block
    target *big.Int
}

func NewProofOfWork(b *Block) *ProofOfWork {
    // big.NewInt(1) 左移 256 - targetBits 位 (即 2 的 256 - targetBits - 1 次方)
    // 找到一个nonce使得区块的hash小于该数
    target := big.NewInt(1)
    target.Lsh(target, uint(256 - targetBits))

    pow := &ProofOfWork{b, target}

    return pow
}

func (pow *ProofOfWork) prepareData(nonce int) []byte {
    // 综合区块信息和nonce
    data := bytes.Join(
        [][]byte{
            pow.block.PrevBlockHash,
            pow.block.Data,
            IntToHex(pow.block.Timestamp),
            IntToHex(int64(targetBits)),
            IntToHex(int64(nonce)),
        },
        []byte{},
    )

    return data
}

func (pow *ProofOfWork) Run() (int, []byte) {
    var hashInt big.Int
    var hash [32]byte
    nonce := 0

    fmt.Printf("Mining the block containing \"%s\"\n", pow.block.Data)
    for nonce < math.MaxInt64 {
        // 从0开始累加nonce，反复计算直至区块hash值小于目标值
        data := pow.prepareData(nonce)
        hash = sha256.Sum256(data)
        hashInt.SetBytes(hash[:])

        // big.Int.Cmp    f func(y *big.Int) (r int)
        if hashInt.Cmp(pow.target) == -1 {
            fmt.Printf("\r%x", hash)
            break
        } else {
            nonce ++
        }
    }
    fmt.Printf("\n\n")

    return nonce, hash[:]
}

// Validate block's Pow
func (pow *ProofOfWork) Validate() bool {
	var hashInt big.Int

	data := pow.prepareData(pow.block.Nonce)
	hash := sha256.Sum256(data)
	hashInt.SetBytes(hash[:])

	isValid := hashInt.Cmp(pow.target) == -1

	return isValid
}
