package main

import (
    "bytes"
    "encoding/binary"
    "log"
)

func IntToHex(num int64) []byte {
    buff := new(bytes.Buffer)
    err := binary.Write(buff, binary.BigEndian, num)
    if err != nil {
        log.Panic(err)
    }

    return buff.Bytes()
}

func ReverseBytes(data []byte) {
    for i, j := 0, len(data) - 1; i < j; i, j = i + 1, j - 1 {
        data[i], data[j] = data[j], data[i]
    }
}

func Int64ToBytes(i int64) []byte {
    var buf = make([]byte, 8)
    binary.BigEndian.PutUint64(buf, uint64(i))
    return buf
}

func BytesToInt64(buf []byte) int64 {
    return int64(binary.BigEndian.Uint64(buf))
}
