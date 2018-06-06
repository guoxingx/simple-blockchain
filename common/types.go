package common

const (
    HashLength    = 32
    AddressLength = 20
)

type Hash [HashLength]byte

func (h *Hash) SetBytes(b []byte) {
    if len(b) > len(h) {
        b = b[len(b) - HashLength:]
    }

    copy(h[HashLength - len(b):], b)
}

func (h Hash) Bytes() []byte { return h[:] }

type Address [AddressLength]byte

func (a *Address) SetBytes(b []byte) {
    if len(b) > len(a) {
        b = b[len(b) - HashLength:]
    }

    copy(a[AddressLength - len(b):], b)
}

func HexToAddress(s string) Address { return BytesToAddress(FromHex(s)) }

func BytesToAddress(b []byte) Address {
    var a Address
    a.SetBytes(b)
    return a
}
