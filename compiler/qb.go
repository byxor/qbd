package compiler

import (
	"encoding/binary"
	"fmt"
	"math"
	"strconv"
)

type CustomByteBuffer struct {
	Buffer []byte
}

type BytecodeCompiler struct {
	RootAstNode AstNode
	Bytes       []byte
}

func GenerateBytecode(compiler *BytecodeCompiler) {

	write := func(bytes ...byte) {
		compiler.Bytes = append(compiler.Bytes, bytes...)
	}

	writeIndex := func(index int, bytes ...byte) {
		i := index
		j := 0
		for {
			if i >= len(compiler.Bytes) || j >= len(bytes) {
				break
			}
			compiler.Bytes[i] = bytes[j]
			i++
			j++
		}
	}

	writeLittleUint32 := func(n uint32) {
		bytes := make([]byte, 4)
		binary.LittleEndian.PutUint32(bytes, n)
		write(bytes...)
	}

	writeLittleUint32Index := func(n uint32, index int) {
		bytes := make([]byte, 4)
		binary.LittleEndian.PutUint32(bytes, n)
		writeIndex(index, bytes...)
	}

	writeLittleUint16 := func(n uint16) {
		bytes := make([]byte, 2)
		binary.LittleEndian.PutUint16(bytes, n)
		write(bytes...)
	}

	writeLittleUint16Index := func(n uint16, index int) {
		bytes := make([]byte, 2)
		binary.LittleEndian.PutUint16(bytes, n)
		writeIndex(index, bytes...)
	}

	nameTable := make(map[string]uint32)

	var writeBytecodeForNode func(node AstNode)
	var writeBytecodeForIf func(node AstNode)
	var writeBytecodeForIfElse func(conditionNode AstNode, bodyNodes []AstNode, elseNodes []AstNode, hasElse bool)

	writeBytecodeForNode = func(node AstNode) {
		switch node.Kind {
		case AstKind_Root:
			data := node.Data.(AstData_Root)
			for _, rootNode := range data.BodyNodes {
				writeBytecodeForNode(rootNode)
			}
		case AstKind_NewLine:
			write(1)
		case AstKind_Comma:
			write(9)
		case AstKind_Break:
			write(0x22)
		case AstKind_AllArguments:
			write(0x2C)
		case AstKind_LocalReference:
			write(0x2D)
			writeBytecodeForNode(node.Data.(AstData_LocalReference).Node)
		case AstKind_Checksum:
			write(0x16)
			data := node.Data.(AstData_Checksum)

			var checksum uint32
			if data.IsRawChecksum {
				temp, _ := strconv.ParseInt(data.ChecksumToken.Data[1:], 16, 32)
				checksum = uint32(temp)
			} else {
				name := data.ChecksumToken.Data
				checksum = StringToChecksum(name)
				nameTable[name] = checksum
			}

			writeLittleUint32(checksum)
		case AstKind_Integer:
			write(0x17)
			intValue, _ := strconv.ParseInt(node.Data.(AstData_Integer).IntegerToken.Data, 10, 32)
			writeLittleUint32(uint32(intValue))
		case AstKind_Float:
			write(0x1A)
			floatValue, _ := strconv.ParseFloat(node.Data.(AstData_Float).FloatToken.Data, 32)
			var bytesValue [4]byte
			binary.LittleEndian.PutUint32(bytesValue[:], math.Float32bits(float32(floatValue)))
			write(bytesValue[:]...)
		case AstKind_String:
			write(0x1B)
			stringData := node.Data.(AstData_String).StringToken.Data
			stringData = stringData[1 : len(stringData)-1]
			writeLittleUint32(uint32(len(stringData) + 1))
			write([]byte(stringData)...)
			write(0)
		case AstKind_Pair:
			write(0x1F)
			floatValueA, _ := strconv.ParseFloat(node.Data.(AstData_Pair).FloatNodeA.Data.(AstData_Float).FloatToken.Data, 32)
			var bytesValueA [4]byte
			binary.LittleEndian.PutUint32(bytesValueA[:], math.Float32bits(float32(floatValueA)))
			write(bytesValueA[:]...)
			floatValueB, _ := strconv.ParseFloat(node.Data.(AstData_Pair).FloatNodeB.Data.(AstData_Float).FloatToken.Data, 32)
			var bytesValueB [4]byte
			binary.LittleEndian.PutUint32(bytesValueB[:], math.Float32bits(float32(floatValueB)))
			write(bytesValueB[:]...)
		case AstKind_Vector:
			write(0x1E)
			floatValueA, _ := strconv.ParseFloat(node.Data.(AstData_Vector).FloatNodeA.Data.(AstData_Float).FloatToken.Data, 32)
			var bytesValueA [4]byte
			binary.LittleEndian.PutUint32(bytesValueA[:], math.Float32bits(float32(floatValueA)))
			write(bytesValueA[:]...)
			floatValueB, _ := strconv.ParseFloat(node.Data.(AstData_Vector).FloatNodeB.Data.(AstData_Float).FloatToken.Data, 32)
			var bytesValueB [4]byte
			binary.LittleEndian.PutUint32(bytesValueB[:], math.Float32bits(float32(floatValueB)))
			write(bytesValueB[:]...)
			floatValueC, _ := strconv.ParseFloat(node.Data.(AstData_Vector).FloatNodeC.Data.(AstData_Float).FloatToken.Data, 32)
			var bytesValueC [4]byte
			binary.LittleEndian.PutUint32(bytesValueC[:], math.Float32bits(float32(floatValueC)))
			write(bytesValueC[:]...)
		case AstKind_UnaryExpression:
			write(0xE)
			writeBytecodeForNode(node.Data.(AstData_UnaryExpression).Node)
			write(0xF)
		case AstKind_AdditionExpression:
			data := node.Data.(AstData_BinaryExpression)
			write(0xE)
			writeBytecodeForNode(data.LeftNode)
			write(0xB)
			writeBytecodeForNode(data.RightNode)
			write(0xF)
		case AstKind_SubtractionExpression:
			data := node.Data.(AstData_BinaryExpression)
			write(0xE)
			writeBytecodeForNode(data.LeftNode)
			write(0xA)
			writeBytecodeForNode(data.RightNode)
			write(0xF)
		case AstKind_MultiplicationExpression:
			data := node.Data.(AstData_BinaryExpression)
			write(0xE)
			writeBytecodeForNode(data.LeftNode)
			write(0xD)
			writeBytecodeForNode(data.RightNode)
			write(0xF)
		case AstKind_DivisionExpression:
			data := node.Data.(AstData_BinaryExpression)
			write(0xE)
			writeBytecodeForNode(data.LeftNode)
			write(0xC)
			writeBytecodeForNode(data.RightNode)
			write(0xF)
		case AstKind_GreaterThanExpression:
			data := node.Data.(AstData_BinaryExpression)
			write(0xE)
			writeBytecodeForNode(data.LeftNode)
			write(0x14)
			writeBytecodeForNode(data.RightNode)
			write(0xF)
		case AstKind_LessThanExpression:
			data := node.Data.(AstData_BinaryExpression)
			write(0xE)
			writeBytecodeForNode(data.LeftNode)
			write(0x12)
			writeBytecodeForNode(data.RightNode)
			write(0xF)
		case AstKind_LessThanEqualsExpression:
			fmt.Println("Warning: <= does not work in THUG Pro")
			data := node.Data.(AstData_BinaryExpression)
			write(0xE)
			writeBytecodeForNode(data.LeftNode)
			write(0x13)
			writeBytecodeForNode(data.RightNode)
			write(0xF)
		case AstKind_GreaterThanEqualsExpression:
			fmt.Println("Warning: >= does not work in THUG Pro")
			data := node.Data.(AstData_BinaryExpression)
			write(0xE)
			writeBytecodeForNode(data.LeftNode)
			write(0x15)
			writeBytecodeForNode(data.RightNode)
			write(0xF)
		case AstKind_EqualsExpression:
			data := node.Data.(AstData_BinaryExpression)
			write(0xE)
			writeBytecodeForNode(data.LeftNode)
			write(7)
			// write(0x11)
			writeBytecodeForNode(data.RightNode)
			write(0xF)
		case AstKind_DotExpression:
			data := node.Data.(AstData_BinaryExpression)
			writeBytecodeForNode(data.LeftNode)
			write(8)
			writeBytecodeForNode(data.RightNode)
		case AstKind_LogicalNot:
			write(0x39)
			writeBytecodeForNode(node.Data.(AstData_UnaryExpression).Node)
		case AstKind_LogicalAnd:
			data := node.Data.(AstData_BinaryExpression)
			writeBytecodeForNode(data.LeftNode)
			write(0x33)
			writeBytecodeForNode(data.RightNode)
		case AstKind_LogicalOr:
			data := node.Data.(AstData_BinaryExpression)
			writeBytecodeForNode(data.LeftNode)
			write(0x32)
			writeBytecodeForNode(data.RightNode)
		case AstKind_Comment:
			//writeBytecodeForNode(AstNode{
			//	Kind: AstKind_String,
			//	Data: AstData_String{
			//		StringToken: Token{
			//			Kind: TokenKind_String,
			//			Data: "\"" + node.Data.(AstData_Comment).CommentToken.Data + "\"",
			//		},
			//	},
			//})
		case AstKind_Script:
			data := node.Data.(AstData_Script)
			write(0x23)
			writeBytecodeForNode(data.NameNode)
			for _, defaultParameterNode := range data.DefaultParameterNodes {
				writeBytecodeForNode(defaultParameterNode)
			}
			for _, bodyNode := range data.BodyNodes {
				writeBytecodeForNode(bodyNode)
			}
			write(0x24)
		case AstKind_IfStatement:
			writeBytecodeForIf(node)
		case AstKind_Random:
			data := node.Data.(AstData_Random)

			numBranches := len(data.Branches)

			write(0x2F)
			writeLittleUint32(uint32(numBranches))

			// write branch weights
			for i := 0; i < numBranches; i++ {
				branchWeightAsInt, _ := strconv.ParseInt(data.BranchWeights[i].Data.(AstData_Integer).IntegerToken.Data, 10, 32)
				writeLittleUint16(uint16(branchWeightAsInt))
			}

			branchOffsetsIndex := len(compiler.Bytes)

			// write dummy branch offsets (populate later)
			for i := 0; i < numBranches; i++ {
				var offset uint32 = 2
				writeLittleUint32(offset)
			}

			// write branches (record sizes for offset calculations, record longjump positions)
			branchSizes := make([]int, numBranches)
			longJumpPositions := make([]int, numBranches-1)
			for i := 0; i < numBranches; i++ {
				start := len(compiler.Bytes)
				for _, branchNode := range data.Branches[i] {
					writeBytecodeForNode(branchNode)
				}
				if i < (numBranches - 1) {
					// write dummy longjump offset (populate later)
					longJumpPositions[i] = len(compiler.Bytes)
					write(0x2e)
					writeLittleUint32(0)
				}
				end := len(compiler.Bytes)
				branchSizes[i] = end - start
			}

			finalIndex := len(compiler.Bytes)

			// update branch offsets with real values
			for i := 0; i < numBranches; i++ {
				offsetIndex := branchOffsetsIndex + (4 * i)

				offsetValue := 0

				// include next branch offsets in offsetValue
				for j := i + 1; j < numBranches; j++ {
					offsetValue += 4
				}

				// include previous branch sizes in offsetValue too
				for j := 0; j < i; j++ {
					offsetValue += branchSizes[j]
				}

				writeLittleUint32Index(uint32(offsetValue), offsetIndex)
			}

			// update longjump offsets with real values
			for i := 0; i < numBranches-1; i++ {
				realOffset := finalIndex - longJumpPositions[i] - 5
				writeLittleUint32Index(uint32(realOffset), longJumpPositions[i]+1)
			}

		case AstKind_WhileLoop:
			//writeBytecodeForNode(AstNode{
			//	Kind: AstKind_String,
			//	Data: AstData_String{
			//		StringToken: Token{
			//			Kind: TokenKind_String,
			//			Data: "\"le while loop\"",
			//		},
			//	},
			//})
			compilerGeneratedChecksum := AstNode{
				Kind: AstKind_Checksum,
				Data: AstData_Checksum{
					ChecksumToken: Token{
						Kind: TokenKind_Identifier,
						Data: "__COMPILER__infinite_loop_bypasser",
					},
				},
			}
			constantIntegerNode := AstNode{
				Kind: AstKind_Integer,
				Data: AstData_Integer{
					IntegerToken: Token{
						Kind: TokenKind_Integer,
						Data: "0",
					},
				},
			}
			writeBytecodeForNode(AstNode{
				Kind: AstKind_Assignment,
				Data: AstData_Assignment{
					NameNode:  compilerGeneratedChecksum,
					ValueNode: constantIntegerNode,
				},
			})
			write(1)
			write(0x20)
			writeBytecodeForNode(AstNode{
				Kind: AstKind_IfStatement,
				Data: AstData_IfStatement{
					Conditions: []AstNode{
						{
							Kind: AstKind_GreaterThanExpression,
							Data: AstData_BinaryExpression{
								LeftNode: AstNode{
									Kind: AstKind_LocalReference,
									Data: AstData_LocalReference{
										Node: compilerGeneratedChecksum,
									},
								},
								RightNode: constantIntegerNode,
							},
						},
					},
					Bodies: [][]AstNode{
						{
							{
								Kind: AstKind_NewLine,
								Data: AstData_Empty{},
							},
							{
								Kind: AstKind_Break,
								Data: AstData_Empty{},
							},
							{
								Kind: AstKind_NewLine,
								Data: AstData_Empty{},
							},
						},
					},
				},
			})

			for _, bodyNode := range node.Data.(AstData_WhileLoop).BodyNodes {
				writeBytecodeForNode(bodyNode)
			}
			write(0x21)
		case AstKind_Return:
			data := node.Data.(AstData_UnaryExpression)
			invocationData := data.Node.Data.(AstData_Invocation)
			write(0x29)
			for _, parameterNode := range invocationData.ParameterNodes {
				writeBytecodeForNode(parameterNode)
			}
		case AstKind_Invocation:
			data := node.Data.(AstData_Invocation)
			writeBytecodeForNode(data.ScriptIdentifierNode)
			for _, parameterNode := range data.ParameterNodes {
				writeBytecodeForNode(parameterNode)
			}
		case AstKind_Assignment:
			data := node.Data.(AstData_Assignment)
			writeBytecodeForNode(data.NameNode)
			write(7)
			writeBytecodeForNode(data.ValueNode)
		case AstKind_Struct:
			write(3)
			for _, elementNode := range node.Data.(AstData_Struct).ElementNodes {
				writeBytecodeForNode(elementNode)
			}
			write(4)
		case AstKind_Array:
			write(5)
			for _, elementNode := range node.Data.(AstData_Array).ElementNodes {
				writeBytecodeForNode(elementNode)
			}
			write(6)
		default:
			fmt.Printf("Warning: no bytecode generated for AstNode of type '%s'\n", node.Kind.String())
		}
	}

	var transformIfElseIf func(ifElseIfNode AstNode) AstNode
	transformIfElseIf = func(ifElseIfNode AstNode) AstNode {
		ifElseIfData := ifElseIfNode.Data.(AstData_IfStatement)

		if len(ifElseIfData.Conditions) == 1 {
			return ifElseIfNode
		}

		bodies := [][]AstNode{
			ifElseIfData.Bodies[0],
		}

		if len(ifElseIfData.Conditions) > 1 {
			bodies = append(bodies,
				[]AstNode{
					{
						Kind: AstKind_NewLine,
						Data: AstData_Empty{},
					},
					transformIfElseIf(AstNode{
						Kind: AstKind_IfStatement,
						Data: AstData_IfStatement{
							Conditions: ifElseIfData.Conditions[1:],
							Bodies:     ifElseIfData.Bodies[1:],
						},
					}),
					{
						Kind: AstKind_NewLine,
						Data: AstData_Empty{},
					},
				},
			)
		}

		return AstNode{
			Kind: AstKind_IfStatement,
			Data: AstData_IfStatement{
				Conditions: []AstNode{
					ifElseIfData.Conditions[0],
				},
				Bodies: bodies,
			},
		}
	}

	writeBytecodeForIf = func(node AstNode) {
		transformedIf := transformIfElseIf(node)
		transformedIfData := transformedIf.Data.(AstData_IfStatement)
		hasElse := len(transformedIfData.Bodies) > 1
		var elseNodes []AstNode
		if hasElse {
			elseNodes = transformedIfData.Bodies[1]
		}
		writeBytecodeForIfElse(
			transformedIfData.Conditions[0],
			transformedIfData.Bodies[0],
			elseNodes,
			hasElse,
		)
	}

	writeBytecodeForIfElse = func(conditionNode AstNode, bodyNodes []AstNode, elseNodes []AstNode, hasElse bool) {
		{
			start := len(compiler.Bytes)
			write(0x47)
			write(0x00) // 2 temporary bytes for branch size
			write(0x00)
			writeBytecodeForNode(conditionNode)
			for _, bodyNode := range bodyNodes {
				writeBytecodeForNode(bodyNode)
			}
			end := len(compiler.Bytes)
			size := end - start
			if hasElse {
				size += 2
			}
			writeLittleUint16Index(uint16(size), start+1)
		}
		if hasElse {
			start := len(compiler.Bytes)
			write(0x48)
			write(0x00) // 2 temporary bytes for branch size
			write(0x00)
			for _, bodyNode := range elseNodes {
				writeBytecodeForNode(bodyNode)
			}
			end := len(compiler.Bytes)
			writeLittleUint16Index(uint16(end-start), start+1)
		}
		write(0x28)
	}

	writeBytecodeForNode(compiler.RootAstNode)

	writeNameTableEntry := func(checksum uint32, name string) {
		write(0x2B)
		writeLittleUint32(checksum)
		write([]byte(name)...)
		write(0)
	}

	for name, checksum := range nameTable {
		writeNameTableEntry(checksum, name)
	}
	write(0)
}