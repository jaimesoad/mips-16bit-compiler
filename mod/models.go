package mod

import (
	"fmt"
	errvals "mips/errors"
	"mips/tokens"
	"mips/utils"
	"slices"
	"strconv"
	"strings"
)

type Jmp struct {
	Line    *uint16
	Callers []int
}

type Mips struct {
	Jumps       map[string]Jmp
	LineNum     uint16
	Whitespaces uint16
	Code        string
	OutputFile  []string
}

func New() Mips {
	return Mips{
		Jumps:       make(map[string]Jmp),
		OutputFile:  []string{"v2.0 raw\n"},
		LineNum:     1,
		Whitespaces: 0,
		Code:        "",
	}
}

func (m *Mips) EvalInstruction(instruction string, output *uint16, lineArr []string) error {
	*output = tokens.Tokens[instruction]

	switch instruction {
	case "li", "lui":

		reg, err := utils.GetRegistry(lineArr[0], m.LineNum)
		if err != nil {
			return err
		}

		num, err := utils.Stoi(lineArr[1], m.LineNum)
		if err != nil {
			return err
		}

		*output |= reg<<8 | num

		// fmt.Println("val:", val, "| num:", num, "| numUint16", numUint16, "| 1-i:", 8*(1 - i))

	case "slt", "eq", "sb", "sgt", "neq", "mov":
		if len(lineArr) != 2 {
			return errvals.WrongArgs(2, instruction, len(lineArr), m.LineNum)
		}

		reg1, err := utils.GetRegistry(lineArr[0], m.LineNum)
		if err != nil {
			return err
		}

		reg2, err := utils.GetRegistry(lineArr[1], m.LineNum)
		if err != nil {
			return err
		}

		*output |= reg1<<4 | reg2

	case "shr", "shl", "rr", "lr", "lb":
		if len(lineArr) != 2 {
			return errvals.WrongArgs(2, instruction, len(lineArr), m.LineNum)
		}

		reg1, err := utils.GetRegistry(lineArr[0], m.LineNum)
		if err != nil {
			return err
		}

		reg2, err := utils.GetRegistry(lineArr[1], m.LineNum)
		if err != nil {
			return err
		}

		*output |= reg1<<8 | reg2<<4

	case "add", "sub", "and", "or":
		if len(lineArr) != 3 {
			return errvals.WrongArgs(3, instruction, len(lineArr), m.LineNum)
		}

		reg1, err := utils.GetRegistry(lineArr[0], m.LineNum)
		if err != nil {
			return err
		}

		if reg1%2 == 1 {
			return errvals.EvenNumber(reg1)
		}

		reg2, err := utils.GetRegistry(lineArr[1], m.LineNum)
		if err != nil {
			return err
		}

		reg3, err := utils.GetRegistry(lineArr[2], m.LineNum)
		if err != nil {
			return err
		}

		*output |= reg1<<8 | reg2<<4 | reg3

	case "mul", "div", "rem", "xor":
		if len(lineArr) != 3 {
			return errvals.WrongArgs(3, instruction, len(lineArr), m.LineNum)
		}

		reg1, err := utils.GetRegistry(lineArr[0], m.LineNum)
		if err != nil {
			return err
		}

		if reg1%2 == 0 {
			return errvals.OddNumber(reg1)
		}

		reg2, err := utils.GetRegistry(lineArr[1], m.LineNum)
		if err != nil {
			return err
		}

		reg3, err := utils.GetRegistry(lineArr[2], m.LineNum)
		if err != nil {
			return err
		}

		*output |= reg1<<8 | reg2<<4 | reg3

	case "jmp", "bra":
		if len(lineArr) != 1 {
			return errvals.WrongArgs(1, instruction, len(lineArr), m.LineNum)
		}

		val := lineArr[0]
		jump, ok := m.Jumps[val]

		if ok {
			if !slices.Contains(jump.Callers, len(m.OutputFile)) {
				jump.Callers = append(jump.Callers, len(m.OutputFile))
			}
		} else {
			jump = Jmp{Callers: []int{len(m.OutputFile)}}
		}

		m.Jumps[val] = jump

	case "spc":
		if len(lineArr) != 1 {
			return errvals.WrongArgs(1, instruction, len(lineArr), m.LineNum)
		}

		reg, err := utils.GetRegistry(lineArr[0], m.LineNum)
		if err != nil {
			return err
		}

		*output |= reg << 8

	case "jr":
		if len(lineArr) != 1 {
			return errvals.WrongArgs(1, instruction, len(lineArr), m.LineNum)
		}

		reg, err := utils.GetRegistry(lineArr[0], m.LineNum)
		if err != nil {
			return err
		}

		*output |= reg << 4

	case "cls":
		if len(lineArr) != 0 {
			return errvals.WrongArgs(1, instruction, len(lineArr), m.LineNum)
		}

	case "lstr":
		err := m.LoadString()
		if err != nil {
			return err
		}

	case "print":
		content, _, err := m.ParseString()
		if err != nil {
			return err
		}

		for _, val := range content {
			m.EvalInstruction("li", output, []string{"r0", fmt.Sprintf("%02x", val)})
			m.OutputFile = append(m.OutputFile, fmt.Sprintf("%04x\n", *output))

			m.EvalInstruction("li", output, []string{"r3", "0x80"})
			m.OutputFile = append(m.OutputFile, fmt.Sprintf("%04x\n", *output))

			m.EvalInstruction("or", output, []string{"r0", "r0", "r3"})
			m.OutputFile = append(m.OutputFile, fmt.Sprintf("%04x\n", *output))

			m.Whitespaces += 3
		}

		*output = 0
		m.EvalInstruction("li", output, []string{"r0", "0"})
		m.OutputFile = append(m.OutputFile, fmt.Sprintf("%04x\n", *output))

	}

	return nil
}

func (m *Mips) LoadString() error {
	content, original, err := m.ParseString()
	if err != nil {
		return err
	}

	line := strings.Replace(m.Code, original, "", -1)
	lineArr := utils.GetValues(line)
	lineArr = lineArr[1:]

	if len(lineArr) != 2 {
		return errvals.WrongArgs(2, "lstr", len(lineArr), m.LineNum)
	}

	if !utils.IsRegister(lineArr[0]) {
		return errvals.NotRegister(lineArr[0], m.LineNum)
	}

	if !utils.IsRegister(lineArr[1]) {
		return errvals.NotRegister(lineArr[1], m.LineNum)
	}

	reg1, reg2 := lineArr[0], lineArr[1]

	instructions := []string{"li", "sb", "li", "sub"}

	reversed := strings.Replace(utils.ReverseString(content), "\"", "", -1)
outer:
	for _, char := range reversed {

		for i, val := range instructions {

			if val == "sub" && char == rune(reversed[len(reversed)-1]) {
				m.Whitespaces += 2
				break outer
			}
			output := tokens.Tokens[val]

			sequence := utils.GetSequence(val, char, reg1, reg2, i)

			m.EvalInstruction(val, &output, sequence)

			m.OutputFile = append(m.OutputFile, fmt.Sprintf("%04x\n", output))
		}

		m.Whitespaces += 4
	}

	m.LineNum++

	return nil
}

func (m *Mips) ParseString() (string, string, error) {
	line := strings.TrimSpace(m.Code)

	idx := 0
	for i, val := range line {
		if val == '"' {
			idx = i
			break
		}
	}

	content, err := strconv.Unquote(line[idx:])
	if err != nil {
		return "", "", fmt.Errorf("error: %s", err.Error())
	}

	return content, line[idx:], nil
}
