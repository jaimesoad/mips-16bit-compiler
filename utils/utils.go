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
	val = strings.Replace(val, "0x", "", -1)
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
	val = strings.Replace(val, "0x", "", -1)
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