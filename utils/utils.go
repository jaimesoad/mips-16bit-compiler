package utils

import (
	"fmt"
	"log"
	errvals "mips/errors"
	"strconv"
	"strings"
)

var registries = func() map[string]uint16 {
	regs := make(map[string]uint16)

	for i := uint16(0); i <= 15; i++ {
		regs[fmt.Sprintf("r%d", i)] = i
	}

	return regs
}()

func Stoi(val string, lineNum uint16) (uint16, error) {
	val = strings.ReplaceAll(val, "0x", "")
	var num uint64
	var err error

	if val[0] == 'r' {
		num = uint64(regToNum(val))
	} else {
		num, err = strconv.ParseUint(val, 16, 16)
	}

	if err != nil || num > 0xff {
		return 0, fmt.Errorf("error: number \"%d\" overflows 2-Byte number\nAt line %d", num, lineNum)
	}

	return uint16(num), nil
}

/*
Transforms strings into 8-bit integers, example:
r12 -> 12
*/
func regToNum(val string) uint16 {
	val = strings.ReplaceAll(val, "0x", "")
	var num int
	var err error

	if val[0] == 'r' {
		num, err = strconv.Atoi(val[1:])
	} else {
		num, err = strconv.Atoi(val)
	}

	if err != nil || num > 0xf {
		log.Fatal(err)
	}

	return uint16(num)
}

func GetRegistry(val string, lineNum uint16) (uint16, error) {
	out, ok := registries[val]

	if !ok {
		return 0, errvals.NotRegister(val, lineNum)
	}

	return out, nil
}

func ReverseString(val string) string {
	var out string

	for _, char := range val {
		out = string(char) + out
	}

	return out
}

func IsRegister(val string) bool {
	_, ok := registries[val]

	return ok
}

func GetSequence(instruction string, char rune, reg1, reg2 string, i int) []string {
	switch instruction {
	case "li":
		if i != 0 {
			return []string{reg2, "0x1"}
		}

		return []string{reg2, fmt.Sprintf("%02x", char)}

	case "sb":
		return []string{reg1, reg2}

	case "sub":
		return []string{reg1, reg1, reg2}
	}

	return []string{}
}

func GetValues(value string) []string {
	var out []string

	arrStr := strings.Split(strings.ReplaceAll(value, ",", ""), " ")

	for _, val := range arrStr {
		if val != "" {
			out = append(out, val)
		}
	}

	return out
}

func RemoveComments(input string) string {
	var result []rune
	inQuotes := false
	escapeNext := false

	for _, char := range input {
		// Handle escaping within quotes
		if inQuotes && char == '\\' && !escapeNext {
			escapeNext = true
			result = append(result, char)
			continue
		}

		// Toggle inQuotes state
		if char == '"' && !escapeNext {
			inQuotes = !inQuotes
			result = append(result, char)
			continue
		}

		escapeNext = false

		// Handle comments: stop processing after unquoted `#`
		if char == '#' && !inQuotes {
			break
		}

		// Add character to result
		result = append(result, char)
	}

	return string(result)
}
