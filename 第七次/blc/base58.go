package blc

import (
	"math/big"
	"bytes"
)

var YX_b58Alphabet = []byte("123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz")

func Base58Encode(input []byte) []byte {
	var result []byte
	x := big.NewInt(0).SetBytes(input)
	base := big.NewInt(int64(len(YX_b58Alphabet)))
	mod := &big.Int{}
	zero := big.NewInt(0)

	for x.Cmp(zero) != 0{
		x.DivMod(x,base,mod)
		result = append(result,YX_b58Alphabet[mod.Int64()])
	}
	ReverseByte(result)
	for by := range input {
		if by == 0 {
			result = append([]byte{YX_b58Alphabet[0]},result...)
		} else {
			break
		}
	}
	return result
}

func Base58Decode(input []byte) []byte {
	result := big.NewInt(0)
	zeroBytes := 0

	for b := range input {
		if b == 0x00 {
			zeroBytes++
		}
	}
	payload := input[zeroBytes:]
	for _, b := range payload {
		charIndex := bytes.IndexByte(YX_b58Alphabet, b)
		result.Mul(result, big.NewInt(58))
		result.Add(result, big.NewInt(int64(charIndex)))
	}

	decoded := result.Bytes()
	decoded = append(bytes.Repeat([]byte{byte(0x00)}, zeroBytes), decoded...)

	return decoded
}