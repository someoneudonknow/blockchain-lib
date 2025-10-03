package transaction

import "bufio"

/*
- one kind for data operation -> move a chunk of data to stack
- one kind for data processing -> get data of the top of stack and do some computation and push the result on to the stack

mov eax, 0x1234
stack: [] 0x1234

it is not allowed for loop -> turing incomplete

byte value -> n, n [0x1, 0x4b] -> data operation, n is length of the chunk of data we need to put on the top of stack [0x4, 0x1, 0x2, 0x3, 0x4]
0x04 -> stack [0x01020304]

4b -> 75, how can we move more than 75 bytes to stack?
n = 0x4c -> OP_PUSHDATA1, the following 1 byte is the length of the chunk of data [0x4c, 0xfe, 0x1, ...]
*/

const (
	// [0x1, 0x4b] -> [1, 75] -> data operation, 75 is length of the chunk of data we need to put on the top of stack
	SCRIPT_DATA_LENGTH_BEGIN = 1
	SCRIPT_DATA_LENGTH_END   = 75
	OP_PUSHDATA1             = 76
	OP_PUSHDATA2             = 77
)

type ScriptSig struct {
	cmds [][]byte
}

func NewScriptSig(reader *bufio.Reader) *ScriptSig {
	cmds := [][]byte{}

	// Read the script length to know how many bytes to read
	scriptLen := ReadVarint(reader).Int64()
	count := int64(0)
	current := make([]byte, 1)

	for count < scriptLen {
		reader.Read(current)

		count++
		currentByte := current[0]
		if currentByte >= SCRIPT_DATA_LENGTH_BEGIN && currentByte <= SCRIPT_DATA_LENGTH_END {
			// push the following byte to stack
			data := make([]byte, currentByte)
			reader.Read(data)
			cmds = append(cmds, data)
			count += int64(currentByte)
		} else if currentByte == OP_PUSHDATA1 {
			// read the next byte as the length of the chunk of data
			length := make([]byte, 1)
			reader.Read(length)

			data := make([]byte, length[0])
			reader.Read(data)

			cmds = append(cmds, data)
			count += int64(length[0] + 1)
		} else if currentByte == OP_PUSHDATA2 {
			// read the next 2 byte as the length of the chunk of data (two byte in little endian format so we have to convert it to big endian)
			lenBuf := make([]byte, 2)
			reader.Read(lenBuf)

			length := LittleEndianToBigInt(lenBuf, LITTLE_ENDIAN_2_BYTES).Int64()
			data := make([]byte, length)
			reader.Read(data)

			cmds = append(cmds, data)
			count += length + 2
		} else {
			// Data processing operation such as OP_DUP, OP_EQUALVERIFY,...
			cmds = append(cmds, []byte{currentByte})
		}

		if count != scriptLen {
			panic("parsing script field failed")
		}

		return &ScriptSig{
			cmds,
		}
	}

	return &ScriptSig{
		cmds,
	}
}
