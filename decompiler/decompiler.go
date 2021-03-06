package decompiler

import (
	"encoding/binary"
	"github.com/byxor/NeverScript/compiler"
	"log"
)

type Arguments struct {
	ByteCode   []byte
	RootNode   compiler.AstNode
	SourceCode string
	NameTable  map[uint32]string
}

func Decompile(arguments *Arguments) {
	err := ParseByteCode(arguments)
	if err != nil {
		log.Fatal(err)
	}

	{ // scrape name table entries
		arguments.NameTable = make(map[uint32]string, 0)
		rootData := arguments.RootNode.Data.(compiler.AstData_Root)
		for _, bodyNode := range rootData.BodyNodes {
			if bodyNode.Kind == compiler.AstKind_NameTableEntry {
				data := bodyNode.Data.(compiler.AstData_NameTableEntry)
				checksum := binary.LittleEndian.Uint32(data.ChecksumBytes)
				arguments.NameTable[checksum] = data.Name
			}
		}
	}

	nsCode, err := DecompileAstNode(arguments.RootNode, 0, arguments.NameTable)
	if err != nil {
		log.Fatalf("\n%s", err)
	}
	arguments.SourceCode = nsCode
}
