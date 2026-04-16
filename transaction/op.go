package transaction

import (
	"ecc"
	"fmt"
	"math/big"
)

const (
	OP_0 = 0
)
const (
	OP_1NEGATE = iota + 79
)
const (
	OP_1 = iota + 81
	OP_2
	OP_3
	OP_4
	OP_5
	OP_6
	OP_7
	OP_8
	OP_9
	OP_10
	OP_11
	OP_12
	OP_13
	OP_14
	OP_15
	OP_16
	OP_NOP
)

const (
	OP_IF = iota + 99
	OP_NOTIf
)

const (
	OP_VERIFY = iota + 105
	OP_RETURN
	OP_TOTALSTACK
	OP_FROMALTSTACK
	OP_2DROP
	OP_2DUP
	OP_3DUP
	OP_2OVER
	OP_2ROT
	OP_2SWAP
	OP_IFDUP
	OP_DEPTH
	OP_DROP
	OP_DUP
	OP_NIP
	OP_OVER
	OP_PICK
	OP_ROLL
	OP_ROT
	OP_SWAP
	OP_TUCK
)

const (
	OP_SIZE = iota + 130
)

const (
	OP_EQUAL = iota + 135
	OP_EQUALVERIFY
)

const (
	OP_1ADD = iota + 139
	OP_1SUB
)

const (
	OP_NEGATE = iota + 143
	OP_ABS
	OP_NOT
	OP_0NOTEQUAL
	OP_ADD
	OP_SUB
	OP_MUL
)

const (
	OP_BOOLAND = iota + 154
	OP_BOOLOR
	OP_NUMEQUAL
	OP_NUMEQUALVERIFY
	OP_NUMNOTEQUAL
	OP_LESSTHAN
	OP_GREATERTHAN
	OP_LESSTHANOREQUAL
	OP_GREATERTHANOREQUAL
	OP_MIN
	OP_MAX
	OP_WITHIN
	OP_RIPEMD160
	OP_SHA1
	OP_SHA256
	OP_HASH160
	OP_HASH256
)

const (
	OP_CHECKSIG = iota + 172
	OP_HECKSIGVERIFY
	OP_CHECKMULTISIG
	OP_CHECKMULTISIGVERIFY
	OP_NOP1
	OP_CHECKLOGTIMEVERIFY
	OP_CHECKSEQUENCEVERIFY
	OP_NOP4
	OP_NOP5
	OP_NOP6
	OP_NOP7
	OP_NOP8
	OP_NOP9
	OP_NOP10
)

type BitcoinOpCode struct {
	opCodeNames map[int]string
	stack       [][]byte
	altStack    [][]byte
	cmds        [][]byte
}

func NewBitcoinOpCode() *BitcoinOpCode {
	opCodeNames := map[int]string{
		0:   "OP_0",
		76:  "OP_PUSHDATA1",
		77:  "OP_PUSHDATA2",
		78:  "OP_PUSHDATA4",
		79:  "OP_1NEGATE",
		81:  "OP_1",
		82:  "OP_2",
		83:  "OP_3",
		84:  "OP_4",
		85:  "OP_5",
		86:  "OP_6",
		87:  "OP_7",
		88:  "OP_8",
		89:  "OP_9",
		90:  "OP_10",
		91:  "OP_11",
		92:  "OP_12",
		93:  "OP_13",
		94:  "OP_14",
		95:  "OP_15",
		96:  "OP_16",
		97:  "OP_NOP",
		99:  "OP_IF",
		100: "OP_NOTIF",
		103: "OP_ELSE",
		104: "OP_ENDIF",
		105: "OP_VERIFY",
		106: "OP_RETURN",
		107: "OP_TOALTSTACK",
		108: "OP_FROMALTSTACK",
		109: "OP_2DROP",
		110: "OP_2DUP",
		111: "OP_3DUP",
		112: "OP_2OVER",
		113: "OP_2ROT",
		114: "OP_2SWAP",
		115: "OP_IFDUP",
		116: "OP_DEPTH",
		117: "OP_DROP",
		118: "OP_DUP",
		119: "OP_NIP",
		120: "OP_OVER",
		121: "OP_PICK",
		122: "OP_ROLL",
		123: "OP_ROT",
		124: "OP_SWAP",
		125: "OP_TUCK",
		130: "OP_SIZE",
		135: "OP_EQUAL",
		136: "OP_EQUALVERIFY",
		139: "OP_1ADD",
		140: "OP_1SUB",
		143: "OP_NEGATE",
		144: "OP_ABS",
		145: "OP_NOT",
		146: "OP_0NOTEQUAL",
		147: "OP_ADD",
		148: "OP_SUB",
		149: "OP_MUL",
		154: "OP_BOOLAND",
		155: "OP_BOOLOR",
		156: "OP_NUMEQUAL",
		157: "OP_NUMEQUALVERIFY",
		158: "OP_NUMNOTEQUAL",
		159: "OP_LESSTHAN",
		160: "OP_GREATERTHAN",
		161: "OP_LESSTHANOREQUAL",
		162: "OP_GREATERTHANOREQUAL",
		163: "OP_MIN",
		164: "OP_MAX",
		165: "OP_WITHIN",
		166: "OP_RIPEMD160",
		167: "OP_SHA1",
		168: "OP_SHA256",
		169: "OP_HASH160",
		170: "OP_HASH256",
		171: "OP_CODESEPARATOR",
		172: "OP_CHECKSIG",
		173: "OP_CHECKSIGVERIFY",
		174: "OP_CHECKMULTISIG",
		175: "OP_CHECKMULTISIGVERIFY",
		176: "OP_NOP1",
		177: "OP_CHECKLOCKTIMEVERIFY",
		178: "OP_CHECKSEQUENCEVERIFY",
		179: "OP_NOP4",
		180: "OP_NOP5",
		181: "OP_NOP6",
		182: "OP_NOP7",
		183: "OP_NOP8",
		184: "OP_NOP9",
		185: "OP_NOP10",
	}
	return &BitcoinOpCode{
		opCodeNames: opCodeNames,
		stack:       make([][]byte, 0),
		altStack:    make([][]byte, 0),
		cmds:        make([][]byte, 0),
	}
}

func (boc *BitcoinOpCode) ExecuteOperation(
	stack [][]byte,
	altstack [][]byte,
	cmd int,
	cmds [][]byte,
	z []byte,
) bool {
	// If the operation executed successfully,then return true, otherwise return false
	switch cmd {
	case OP_CHECKSIG:
		return boc.opCheckSig(z)
	default:
		panic("Not implemented yet")
	}
}

/*
OP_CHECKSIG:
	- Verify validity of message z
	- DER binary data of the signature and the uncompressed sec public key are top two elements of the stack
	- If the verification success, push 1 to the stack, otherwise push 0 to the stack
	- If the script is using uncopressed public key, then script is called p2pk (pay to public key)
*/
func (boc *BitcoinOpCode) opCheckSig(zBin []byte) bool {
	// Notice!! we need to remove the last byte of the der binary data, because this byte is used to indicate the hash type
	if len(boc.stack) < 2 {
		fmt.Println("Current stack data is invalid.")
		return false
	}

	// Get the public key
	pubKey := boc.stack[len(boc.stack)-1]
	boc.stack = boc.stack[0 : len(boc.stack)-1]
	// Get the DER encoded signature
	derSig := boc.stack[len(boc.stack)-1]
	boc.stack = boc.stack[0 : len(boc.stack)-1]
	point := ecc.ParseSEC(pubKey)
	signature := ecc.ParseSigBin(derSig)

	z := new(big.Int)
	z.SetBytes(zBin)
	n := ecc.BitcoinN()
	zField := ecc.NewFieldElement(n, z)
	if point.Verify(zField, signature) == true {
		boc.stack = append(boc.stack, boc.EncodeNum(1))
	} else {
		boc.stack = append(boc.stack, boc.EncodeNum(0))
	}

	return true
}

func (boc *BitcoinOpCode) DecodeNum(element []byte) int64 {
	bigEndian := reverseByteSlice(element)
	negative := false
	result := int64(0)

	if (bigEndian[0] & 0x80) != 0 {
		negative = true
		// reset the msb to 0
		// 0x7f is 0111 1111
		result = int64(bigEndian[0] & 0x7f)
	} else {
		negative = false
		result = int64(bigEndian[0])
	}

	for i := 1; i < len(bigEndian); i++ {
		result <<= 8 // 8 bit => 1 byte, prepare space for the next byte
		result += int64(bigEndian[i])
	}

	if negative {
		return -result
	}

	return result
}

// Encode number using two complement algorithm, use the most significant byte as sign byte
func (boc *BitcoinOpCode) EncodeNum(num int64) []byte {
	if num == 0 {
		return []byte("")
	}

	result := []byte{}
	absNum := num
	negative := false

	if num < 0 {
		absNum = -num
		negative = true
	}

	if absNum > 0 {
		// Result will be in little endian format so we append the last byte first (actually we're reversing the byte array)
		result = append(result, byte(absNum&0xff))
		absNum >>= 8
	}

	// Check the most significant byte, notice the most significant byte is at the end of the result
	// 0x8080 -> 32896 -32896
	// 0x80 -> 1000 0000
	// 0.... & 1000 0000 -> 0000 0000
	// 1.... & 1000 0000 -> 1000 0000

	if (result[len(result)-1] & 0x80) != 0 {
		if negative {
			// insert 0x80 in the head, most significant byte is currently at the end of result
			result = append(result, 0x80)
		} else {
			result = append(result, 0x00)
		}
	} else {
		// set the most significant bit to one
		// 1000 0000 || 0... ....
		result[len(result)-1] |= 0x80
	}

	return result
}
