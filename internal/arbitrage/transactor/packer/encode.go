package packer

import (
	"encoding/hex"
	"math/big"
)

const (
	swapV2Sig = "9b604860"
	swapV3Sig = "ac7d4304"
	token0Sig = "0dfe1681"
	token1Sig = "d21220a7"
)

var (
	one       = big.NewInt(1)
	zero      = big.NewInt(0)
	oneBytes  = one.Bytes()
	zeroBytes = zero.Bytes()
	tt256     = new(big.Int).Lsh(big.NewInt(1), 256)   // 2 ** 256
	tt256m1   = new(big.Int).Sub(tt256, big.NewInt(1)) // 2 ** 256 - 1
)

func encodeBool(val bool) []byte {
	if val {
		return leftPad32(oneBytes)
	}
	return leftPad32(zeroBytes)
}

func encodeAddress(addressBytes []byte) []byte {
	return leftPad32(addressBytes)
}

func encodeUint256(n *big.Int) []byte {
	b := new(big.Int)
	b = b.Set(n)

	if b.Sign() < 0 || b.BitLen() > 256 {
		b.And(b, tt256m1)
	}

	return leftPad32(b.Bytes())
}

// use only for hex parameters without 0x
func mustDecodeHex(src string) []byte {
	dst := make([]byte, hex.DecodedLen(len([]byte(src))))
	n, err := hex.Decode(dst, []byte(src))
	if err != nil {
		panic(err)
	}

	return dst[:n]
}

func leftPad32(b []byte) []byte {
	return padBytes(b, 32, true)
}

func padBytes(b []byte, size int, left bool) []byte {
	l := len(b)
	if l == size {
		return b
	}
	if l > size {
		return b[l-size:]
	}
	tmp := make([]byte, size)
	if left {
		copy(tmp[size-l:], b)
	} else {
		copy(tmp, b)
	}
	return tmp
}
