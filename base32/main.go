package main

import (
	"errors"
	"fmt"
	"math/rand/v2"
	"slices"
	"strings"
	"time"
)

// base32EncodedLengthFor64Bits 是一個 64-bit Id 編碼成 Base32 後的固定長度。
//
// Base32 的 "32" 代表每一個字元可以表示 32 種不同值。
//
// 因為 32 = 2^5 所以每一個 Base32 字元需要 5 bits 表示。
//
// 例如:
//
//	00000 -> 0
//	00001 -> 1
//	...
//	11111 -> 31
//
// 每次可以從 binary Id 取出 5 bits，轉換成一個 Base32 字元。
// 最終需要 13 (64 / 5 = 12.8) 個 Base32 字元表示 64-bit Id
const base32EncodedLengthFor64Bits = 13

// crockfordAlphabet 是 Crockford Base32 的編碼字元表，encode/decode 共用同一份。
const crockfordAlphabet = "0123456789ABCDEFGHJKMNPQRSTVWXYZ"

// EncodeCrockfordBase32 將 id 編碼成 Crockford Base32 字串。
// 若有熱路徑優化需求，改用 [AppendCrockfordBase32]。
func EncodeCrockfordBase32(id uint64) string {
	var buffer [base32EncodedLengthFor64Bits]byte
	return string(AppendCrockfordBase32(id, buffer[:0]))
}

// AppendCrockfordBase32
// 把 id 編碼後附加到 dst 尾端，buffer 由呼叫端擁有並重複利用，省下逐筆的記憶體配置成本。
// dst 是為效能而存在的 optional 參數，傳 nil 也能正常運作。
func AppendCrockfordBase32(id uint64, dst []byte) []byte {
	dst = slices.Grow(dst, base32EncodedLengthFor64Bits)

	// 伸展 dst 的長度，buf 直接指向 dst 尾端那段記憶體，寫入 buf 即寫入 dst 本體。
	n := len(dst)
	dst = dst[:n+base32EncodedLengthFor64Bits]
	buf := dst[n:]

	for i := base32EncodedLengthFor64Bits - 1; i >= 0; i-- {
		buf[i] = crockfordAlphabet[id&31] // 取最低 5 bits (31 = 0b11111)
		id >>= 5                          // 消耗掉已處理的 5 bits
	}

	return dst
}

// decodeTable 把任一 byte 對應到它在 crockfordAlphabet 中的值 (0-31)。
// 不合法的字元對應到 -1，方便 decode 時判斷驗證失敗。
//
// Crockford Base32 規定 I、L、O 是容易與 1、1、0 搞混的字元，
// 解碼時要把它們正規化成對應的合法字元，並且不分大小寫。
var decodeTable [256]int8

func init() {
	for i := range decodeTable {
		decodeTable[i] = -1
	}
	for value, ch := range []byte(crockfordAlphabet) {
		decodeTable[ch] = int8(value)
	}

	// Crockford 建議的易混淆字元正規化。
	decodeTable['O'] = decodeTable['0']
	decodeTable['I'] = decodeTable['1']
	decodeTable['L'] = decodeTable['1']

	// 小寫字元對應到跟大寫相同的值。
	for ch := byte('a'); ch <= 'z'; ch++ {
		decodeTable[ch] = decodeTable[ch-'a'+'A']
	}
}

var (
	// ErrInvalidLength 代表輸入字串長度不是 13。
	ErrInvalidLength = errors.New("base32: 長度必須為 13 個字元")

	// ErrInvalidChar 代表輸入字串包含不合法的字元。
	ErrInvalidChar = errors.New("base32: 包含不合法的字元")

	// ErrOverflow 代表輸入字串解碼後超出 64-bit 範圍。
	//
	// 13 個字元總共可以表示 65 bits (13 * 5)，
	// 比 64-bit 多 1 bit，所以第一個字元只能使用最低 4 bits (0-15)，
	// 否則解碼結果會溢位。
	ErrOverflow = errors.New("base32: 數值超出 64-bit 範圍")

	// ErrRoundTripMismatch 代表 decode 後再 encode 回去，
	// 結果與正規化後的輸入不一致 (雙重驗證失敗)。
	ErrRoundTripMismatch = errors.New("base32: 雙重驗證失敗，encode(decode(s)) != s")
)

// DecodeCrockfordBase32 將 Crockford Base32 字串解碼成 64-bit Id。
//
// 驗證分成兩層 (雙重驗證):
//
//  1. 字元驗證: 長度必須是 13，且每個字元都必須是合法的
//     Crockford Base32 字元 (含大小寫、易混淆字元正規化)，
//     並且解碼後的數值不能超出 64-bit 範圍。
//  2. 往返驗證 (round-trip): 把 decode 出來的 id 重新 encode，
//     結果必須與正規化後的輸入完全相同，
//     否則代表輸入雖然「合法」但不是由 EncodeCrockfordBase32 產生的正規表示法。
func DecodeCrockfordBase32(s string) (uint64, error) {
	if len(s) != base32EncodedLengthFor64Bits {
		return 0, ErrInvalidLength
	}

	var id uint64

	for i := range base32EncodedLengthFor64Bits {
		value := decodeTable[s[i]]
		if value < 0 {
			return 0, ErrInvalidChar
		}

		// 第一個字元只允許使用最低 4 bits，
		// 否則代表數值會超出 64-bit 範圍。
		if i == 0 && value > 15 {
			return 0, ErrOverflow
		}

		id = (id << 5) | uint64(value)
	}

	// 第二層驗證: 往返驗證。
	if EncodeCrockfordBase32(id) != strings.ToUpper(s) {
		return 0, ErrRoundTripMismatch
	}

	return id, nil
}

func main() {
	start := uint64(time.Now().UnixNano())
	end := start + 10
	for i := start; i < end; i++ {
		id := uint64(i)
		// id := rand.Uint64()
		encoded := string(AppendCrockfordBase32(id, nil))

		decoded, err := DecodeCrockfordBase32(encoded)
		if err != nil {
			panic(fmt.Sprintf("decode 失敗: %v (encoded=%s)", err, encoded))
		}

		if decoded != id {
			panic(fmt.Sprintf("雙重驗證失敗: id=%d encoded=%s decoded=%d", id, encoded, decoded))
		}

		println(i, encoded, "-> decode OK:", decoded == id)
	}

	// 驗證不合法輸入會被正確擋下。
	if _, err := DecodeCrockfordBase32("invalid-len"); !errors.Is(err, ErrInvalidLength) {
		panic("應該回傳 ErrInvalidLength")
	}
	if _, err := DecodeCrockfordBase32("UUUUUUUUUUUUU"); !errors.Is(err, ErrInvalidChar) {
		panic("應該回傳 ErrInvalidChar")
	}
	if _, err := DecodeCrockfordBase32("ZZZZZZZZZZZZZ"); !errors.Is(err, ErrOverflow) {
		panic("應該回傳 ErrOverflow")
	}

	println("全部驗證通過")

	// AppendCrockfordBase32 的使用範例:
	// 重複利用同一塊 buffer，把多個 id 用逗號組成一行輸出，
	// 過程中不需要為每個 id 各自產生一個 string。
	var line []byte
	for i := range 5 {
		if i > 0 {
			line = append(line, ',')
		}
		line = AppendCrockfordBase32(rand.Uint64(), line)
	}
	println("batch:", string(line))
}
