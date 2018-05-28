package main

/*
将一个公钥转换成一个 Base58 地址：
    使用 RIPEMD160(SHA256(PubKey)) 哈希算法，取公钥并对其哈希两次
    给哈希加上地址生成算法版本的前缀
    对于第二步生成的结果，使用 SHA256(SHA256(payload)) 再哈希，计算校验和。校验和是结果哈希的前四个字节。
    将校验和附加到 version+PubKeyHash 的组合中。
    使用 Base58 对 version+PubKeyHash+checksum 组合进行编码。
*/

import (
    "log"
    "bytes"
    "crypto/ecdsa"
    "crypto/sha256"
    "crypto/elliptic"
    "crypto/rand"

    "golang.org/x/crypto/ripemd160"
)

const addressChecksumLen = 4
const version = byte(0x00)

type Wallet struct {
    PrivateKey ecdsa.PrivateKey
    PublicKey  []byte
}

func NewWallet() *Wallet {
    private, public := newKeyPair()
    wallet := Wallet{private, public}

    return &wallet
}

// 生成新的公私钥
func newKeyPair() (ecdsa.PrivateKey, []byte) {
    curve := elliptic.P256()
    private, err := ecdsa.GenerateKey(curve, rand.Reader)
    if err != nil { log.Panic(err) }

    // ... 在这里类似于python的 *list
    // 将Y.Bytes() 逐个append到 X.Bytes()
    pubKey := append(private.PublicKey.X.Bytes(), private.PublicKey.Y.Bytes()...)

    return *private, pubKey
}

/*
地址由三个部分组成
    1. Version
    2. Public key hash
    3. Checksum 来自sha256(sha256(PublicKeyHash))
*/
func (w Wallet) GetAddress() []byte {
    pubKeyHash := HashPubKey(w.PublicKey)

    versionedPayload := append([]byte{ version }, pubKeyHash...)
    checksum := checksum(versionedPayload)

    fullPayload := append(versionedPayload, checksum...)
    address := Base58Encode(fullPayload)

    return address
}

// 将address string 转换成pubKeyHash []byte
func AddressToPubKeyHash(address string) (pubKeyHash []byte) {
    pubKeyHash = Base58Decode([]byte(address))
    pubKeyHash = pubKeyHash[1 : len(pubKeyHash) - addressChecksumLen]
    return
}

func ValidateAddress(address string) bool {
    pubKeyHash := Base58Decode([]byte(address))
    actualChecksum := pubKeyHash[len(pubKeyHash) - addressChecksumLen:]
    version := pubKeyHash[0]
    pubKeyHash = pubKeyHash[1:len(pubKeyHash) - addressChecksumLen]

    targetChecksum := checksum(append([]byte{version}, pubKeyHash...))

    return bytes.Compare(actualChecksum, targetChecksum) == 0
}

// hash PublicKey
func HashPubKey(pubKey []byte) []byte {
    publicSHA256 := sha256.Sum256(pubKey)

    RIPEMD160Hasher := ripemd160.New()
    _, err := RIPEMD160Hasher.Write(publicSHA256[:])
    if err != nil { log.Panic(err) }

    publicRIPEMD160 := RIPEMD160Hasher.Sum(nil)

	return publicRIPEMD160
}

// just sha256.Sum256 twice
func checksum(payload []byte) []byte {
    firstSHA := sha256.Sum256(payload)
    secondSHA := sha256.Sum256(firstSHA[:])

    return secondSHA[:addressChecksumLen]
}
