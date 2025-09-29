package transaction

import "bufio"

type ScriptSig struct {
	reader *bufio.Reader
}

func NewScriptSig(reader *bufio.Reader) *ScriptSig {
	return &ScriptSig{
		reader: reader,
	}
}
