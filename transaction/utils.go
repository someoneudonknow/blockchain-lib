package transaction

import (
	"bufio"
	"math/big"

	"github.com/tsuna/endian"
)

type LITTLE_ENDIAN_LENGTH int

const (
	LITTLE_ENDIAN_2_BYTES = iota
	LITTLE_ENDIAN_4_BYTES
	LITTLE_ENDIAN_8_BYTES
)

func ReadVarint(bufReader *bufio.Reader) *big.Int {
	i := make([]byte, 1)
	bufReader.Read(i)

	v := new(big.Int)
	v.SetBytes(i)

	if v.Cmp(big.NewInt(int64(0xfd))) < 0 {
		return v
	}

	if v.Cmp(big.NewInt(int64(0xfd))) == 0 {
		i1 := make([]byte, 2)
		bufReader.Read(i1)
		return LittleEndianToBigInt(i1, LITTLE_ENDIAN_2_BYTES)
	}

	if v.Cmp(big.NewInt(int64(0xfe))) == 0 {
		i2 := make([]byte, 4)
		bufReader.Read(i2)
		return LittleEndianToBigInt(i2, LITTLE_ENDIAN_4_BYTES)
	}

	if v.Cmp(big.NewInt(int64(0xff))) == 0 {
		i3 := make([]byte, 8)
		bufReader.Read(i3)
		return LittleEndianToBigInt(i3, LITTLE_ENDIAN_8_BYTES)
	}

	panic("Invalid varint")
}

func BigIntToLittleEndian(bigInt *big.Int, length LITTLE_ENDIAN_LENGTH) []byte {
	switch length {
	case LITTLE_ENDIAN_2_BYTES:
		val := bigInt.Int64()
		littleEndian := endian.HostToNetUint16(uint16(val))
		p := big.NewInt(int64(littleEndian))
		return p.Bytes()
	case LITTLE_ENDIAN_4_BYTES:
		val := bigInt.Int64()
		littleEndian := endian.HostToNetUint32(uint32(val))
		p := big.NewInt(int64(littleEndian))
		return p.Bytes()
	case LITTLE_ENDIAN_8_BYTES:
		val := bigInt.Int64()
		littleEndian := endian.HostToNetUint64(uint64(val))
		p := big.NewInt(int64(littleEndian))
		return p.Bytes()
	}
	panic("Not implemented error")
}

func LittleEndianToBigInt(bytes []byte, length LITTLE_ENDIAN_LENGTH) *big.Int {
	switch length {
	case LITTLE_ENDIAN_2_BYTES:
		p := new(big.Int)
		p.SetBytes(bytes)
		val := endian.NetToHostUint16(uint16(p.Int64()))
		return big.NewInt(int64(val))
	case LITTLE_ENDIAN_4_BYTES:
		p := new(big.Int)
		p.SetBytes(bytes)
		val := endian.NetToHostUint32(uint32(p.Int64()))
		return big.NewInt(int64(val))
	case LITTLE_ENDIAN_8_BYTES:
		p := new(big.Int)
		p.SetBytes(bytes)
		val := endian.NetToHostUint64(uint64(p.Int64()))
		return big.NewInt(int64(val))
	}
	panic("Not implemented error")
}
