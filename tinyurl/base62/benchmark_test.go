package base62

import (
	"crypto/rand"
	"encoding/base64"
	"math"
	"math/big"
	"testing"
)

var testRandBytes = make([]byte, 16)
var testEncodedBytes []byte
var testEncodedBase64 []byte
var testInteger = uint64(math.MaxInt64)
var testEncodedInteger = []byte("V8qRkBGKRiP")

func init() {
	if _, err := rand.Read(testRandBytes); err != nil {
		panic(err)
	}
	testEncodedBytes = Encode(testRandBytes)

	testEncodedBase64 = make([]byte, base64.RawStdEncoding.EncodedLen(len(testRandBytes)))
	base64.RawStdEncoding.Encode(testEncodedBase64, testRandBytes)
}

func encodeWithBigInt(b []byte) []byte {
	base := big.NewInt(base)
	num := new(big.Int).SetBytes(b)
	mod := new(big.Int)

	ret := make([]byte, 0, len(b)*8/5+1)
	for num.BitLen() > 0 {
		num.DivMod(num, base, mod)
		ret = append(ret, encodeStd[mod.Int64()])
	}
	return ret
}

func BenchmarkEncodeV2(bb *testing.B) {
	for i := 0; i < bb.N; i++ {
		_ = StdEncoding._encodeV2(testRandBytes)
	}
}

func BenchmarkEncodeLegacyV1(bb *testing.B) {
	for i := 0; i < bb.N; i++ {
		_ = StdEncoding._encodeV1(testRandBytes)
	}
}

func BenchmarkEncodeToString(bb *testing.B) {
	for i := 0; i < bb.N; i++ {
		_ = EncodeToString(testRandBytes)
	}
}

func BenchmarkEncodeToBuf(bb *testing.B) {
	buf := make([]byte, 0, 1000)
	for i := 0; i < bb.N; i++ {
		_ = EncodeToBuf(buf, testRandBytes)
	}
}

func BenchmarkDecode(bb *testing.B) {
	for i := 0; i < bb.N; i++ {
		_, _ = Decode(testEncodedBytes)
	}
}

func BenchmarkDecodeString(bb *testing.B) {
	s := string(testEncodedBytes)
	for i := 0; i < bb.N; i++ {
		_, _ = DecodeString(s)
	}
}

func BenchmarkDecodeToBuf(bb *testing.B) {
	buf := make([]byte, 0, 1000)
	for i := 0; i < bb.N; i++ {
		_, _ = DecodeToBuf(buf, testRandBytes)
	}
}

func BenchmarkEncodeBigInt(bb *testing.B) {
	for i := 0; i < bb.N; i++ {
		_ = encodeWithBigInt(testRandBytes)
	}
}

func BenchmarkBase64EncodeToString(bb *testing.B) {
	for i := 0; i < bb.N; i++ {
		_ = base64.RawStdEncoding.EncodeToString(testRandBytes)
	}
}

func BenchmarkBase64Encode(bb *testing.B) {
	buf := make([]byte, 1000)
	for i := 0; i < bb.N; i++ {
		base64.RawStdEncoding.Encode(buf, testRandBytes)
	}
}

func BenchmarkBase64DecodeString(bb *testing.B) {
	s := string(testEncodedBase64)
	for i := 0; i < bb.N; i++ {
		_, _ = base64.RawStdEncoding.DecodeString(s)
	}
}

func BenchmarkBase64Decode(bb *testing.B) {
	buf := make([]byte, 1000)
	for i := 0; i < bb.N; i++ {
		_, _ = base64.RawStdEncoding.Decode(buf, testEncodedBase64)
	}
}

func BenchmarkEncodeInteger(bb *testing.B) {
	for i := 0; i < bb.N; i++ {
		_ = FormatUint(testInteger)
	}
}

func BenchmarkDecodeInteger(bb *testing.B) {
	for i := 0; i < bb.N; i++ {
		_, _ = ParseUint(testEncodedInteger)
	}
}
