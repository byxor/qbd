// Code generated by "stringer -type=Token"; DO NOT EDIT.

package tokens

import "strconv"

const _Token_name = "EndOfFileEndOfLineAssignmentLocalReferenceAllLocalReferencesSubtractionAdditionDivisionMultiplicationNotEqualityCheckLessThanCheckLessThanOrEqualCheckGreaterThanCheckGreaterThanOrEqualCheckStartOfExpressionEndOfExpressionStartOfStructEndOfStructStartOfArrayEndOfArrayStartOfSwitchEndOfSwitchSwitchCaseDefaultSwitchCaseStartOfFunctionEndOfFunctionReturnBreakStartOfIfElseElseIfEndOfIfOptimisedIfOptimisedElseIntegerFloatNameShortJumpChecksumTableEntryNamespaceAccessInvalid"

var _Token_index = [...]uint16{0, 9, 18, 28, 42, 60, 71, 79, 87, 101, 104, 117, 130, 150, 166, 189, 206, 221, 234, 245, 257, 267, 280, 291, 301, 318, 333, 346, 352, 357, 366, 370, 376, 383, 394, 407, 414, 419, 423, 432, 450, 465, 472}

func (i Token) String() string {
	if i < 0 || i >= Token(len(_Token_index)-1) {
		return "Token(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _Token_name[_Token_index[i]:_Token_index[i+1]]
}
