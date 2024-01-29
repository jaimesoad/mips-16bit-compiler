package main

import (
	"bufio"
	"fmt"
	errvals "mips/errors"
	"mips/utils"
	"os"
	"slices"
	"strconv"
	"strings"
)

type jmp struct {
	line    *uint16
	callers []int
}

// even tokens only take even registries and odd odds.
var tokens = map[string]uint16{
	"add":   0x0000, // even
	"mul":   0x0000, // odd
	"sub":   0x1000, // even
	"div":   0x1000, // odd
	"and":   0x2000, // even
	"rem":   0x2000, // odd
	"or":    0x3000, // even
	"xor":   0x3000, // odd
	"shr":   0x4000, // even
	"rr":    0x4000, // odd
	"shl":   0x5000, // even
	"lr":    0x5000, // odd
	"slt":   0x6000,
	"sgt":   0x6100,
	"eq":    0x7000,
	"neq":   0x7100,
	"sb":    0x8000,
	"mov":   0x8100,
	"lb":    0x9000,
	"li":    0xA000,
	"lui":   0xB000,
	"jmp":   0xC000,
	"bra":   0xD000,
	"jr":    0xE000,
	"spc":   0xF000,
	"cls":   0xF100,
	"lstr":  0,
	"print": 0,
}

var jumps = make(map[string]jmp)
var outputFile = []string{"v2.0 raw\n"}
var lineNum = uint16(1)
var whitespaces = uint16(0)
var code = ""

func main() {

	if len(os.Args) != 2 {
		fmt.Printf("error: expected 1 argument but got %d\n", len(os.Args)-1)
		return
	}

	file, err := os.Open(os.Args[1])
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fileArr := strings.Split(os.Args[1], ".")

	if len(fileArr) <= 1 {
		fmt.Printf("error: File \"%s\" has no extension\nConsirder renaming it to: \"%s.s\"\n", os.Args[1], os.Args[1])
		return
	}

	if fileArr[len(fileArr)-1] != "s" {
		fmt.Printf("error: File extension \"%s\" not recognized\n", fileArr[len(fileArr)-1])
		return
	}

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		code = strings.Split(scanner.Text(), "#")[0]
		line := strings.ToLower(strings.TrimSpace(code))
		lineArr := getValues(line)

		if line == "" {
			/* whitespaces++
			lineNum++
			outputFile = append(outputFile, "0000\n") */
			lineNum++
			continue
		}

		if line[len(line)-1] == ':' {
			tag := line[:len(line)-1]
			jump, ok := jumps[tag]

			if ok {
				if jump.line != nil {
					fmt.Printf("Error: \"%s\" already defined\n", tag)
					return
				}

				jump.line = new(uint16)
				*jump.line = whitespaces

			} else {
				jump = jmp{new(uint16), []int{}}
				*jump.line = whitespaces
			}

			jumps[tag] = jump

			lineNum++

			continue
		}

		if strings.Contains(line, "*") {
			newLines, err := strconv.Atoi(strings.Split(line, "*")[0])
			if err != nil {
				fmt.Println("Syntax error at line:", lineNum, "\nInstruction", line, "not recognized")
				return
			}

			whitespaces += uint16(newLines)
			lineNum++
			outputFile = append(outputFile, line+"\n")

			continue
		}

		instruction := lineArr[0]
		lineArr = lineArr[1:]
		output, ok := tokens[instruction]

		if !ok {
			fmt.Printf("Syntax error: instruction %s, \nAt line: %d\n%s not recognized.\n", instruction, lineNum, line)
			return
		}

		err := EvalInstruction(instruction, &output, lineArr)
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		if !(instruction == "lstr" || instruction == "print") {
			outputFile = append(outputFile, fmt.Sprintf("%04x\n", output))

		}

		lineNum++
		whitespaces++
	}

	//var jumps map[string]uint16

	for key, value := range jumps {

		if value.line == nil {
			times := "time"
			if len(value.callers) > 1 {
				times += "s"
			}

			fmt.Printf("error: tag \"%s\" referenced %d %s was never declared\n", key, len(value.callers), times)
			return
		}

		if len(value.callers) == 0 {
			fmt.Printf("error: tag \"%s\" never used.\n", key)
			return
		}

		for _, i := range value.callers {

			newVal := string(outputFile[i][0])

			jp := fmt.Sprintf("%02x\n", *value.line)

			if len(jp) <= 3 {
				newVal += "0" + jp

			} else {
				newVal += jp
			}

			outputFile[i] = newVal
		}
	}

	/* for _, value := range outputFile {
		fmt.Print(value)
	} */

	fileName := strings.Join(fileArr[:len(fileArr)-1], "") + ".out"

	f, err := os.Create(fileName)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	defer f.Close()

	for _, line := range outputFile {
		f.WriteString(line)
	}

	f.Sync()
}

func getValues(value string) []string {
	var out []string

	arrStr := strings.Split(strings.Replace(value, ",", "", -1), " ")

	for _, val := range arrStr {
		if val != "" {
			out = append(out, val)
		}
	}

	return out
}

func EvalInstruction(instruction string, output *uint16, lineArr []string) error {
	*output = tokens[instruction]

	switch instruction {
	case "li", "lui":

		reg, err := utils.GetRegistry(lineArr[0], lineNum)
		if err != nil {
			return err
		}

		num, err := utils.Stoi(lineArr[1], lineNum)
		if err != nil {
			return err
		}

		*output |= reg<<8 | num

		//fmt.Println("val:", val, "| num:", num, "| numUint16", numUint16, "| 1-i:", 8*(1 - i))

	case "slt", "eq", "sb", "sgt", "neq", "mov":
		if len(lineArr) != 2 {
			return errvals.WrongArgs(2, instruction, len(lineArr), lineNum)
		}

		reg1, err := utils.GetRegistry(lineArr[0], lineNum)
		if err != nil {
			return err
		}

		reg2, err := utils.GetRegistry(lineArr[1], lineNum)
		if err != nil {
			return err
		}

		*output |= reg1<<4 | reg2

	case "shr", "shl", "rr", "lr", "lb":
		if len(lineArr) != 2 {
			return errvals.WrongArgs(2, instruction, len(lineArr), lineNum)
		}

		reg1, err := utils.GetRegistry(lineArr[0], lineNum)
		if err != nil {
			return err
		}

		reg2, err := utils.GetRegistry(lineArr[1], lineNum)
		if err != nil {
			return err
		}

		*output |= reg1<<8 | reg2<<4

	case "add", "sub", "and", "or":
		if len(lineArr) != 3 {
			return errvals.WrongArgs(3, instruction, len(lineArr), lineNum)
		}

		reg1, err := utils.GetRegistry(lineArr[0], lineNum)
		if err != nil {
			return err
		}

		if reg1%2 == 1 {
			return errvals.EvenNumber(reg1)
		}

		reg2, err := utils.GetRegistry(lineArr[1], lineNum)
		if err != nil {
			return err
		}

		reg3, err := utils.GetRegistry(lineArr[2], lineNum)
		if err != nil {
			return err
		}

		*output |= reg1<<8 | reg2<<4 | reg3

	case "mul", "div", "rem", "xor":
		if len(lineArr) != 3 {
			return errvals.WrongArgs(3, instruction, len(lineArr), lineNum)
		}

		reg1, err := utils.GetRegistry(lineArr[0], lineNum)
		if err != nil {
			return err
		}

		if reg1%2 == 0 {
			return errvals.OddNumber(reg1)
		}

		reg2, err := utils.GetRegistry(lineArr[1], lineNum)
		if err != nil {
			return err
		}

		reg3, err := utils.GetRegistry(lineArr[2], lineNum)
		if err != nil {
			return err
		}

		*output |= reg1<<8 | reg2<<4 | reg3

	case "jmp", "bra":
		if len(lineArr) != 1 {
			return errvals.WrongArgs(1, instruction, len(lineArr), lineNum)
		}

		val := lineArr[0]
		jump, ok := jumps[val]

		if ok {
			if !slices.Contains(jump.callers, len(outputFile)) {
				jump.callers = append(jump.callers, len(outputFile))

			}

		} else {
			jump = jmp{callers: []int{len(outputFile)}}
		}

		jumps[val] = jump

	case "spc":
		if len(lineArr) != 1 {
			return errvals.WrongArgs(1, instruction, len(lineArr), lineNum)
		}

		reg, err := utils.GetRegistry(lineArr[0], lineNum)
		if err != nil {
			return err
		}

		*output |= reg << 8

	case "jr":
		if len(lineArr) != 1 {
			return errvals.WrongArgs(1, instruction, len(lineArr), lineNum)
		}

		reg, err := utils.GetRegistry(lineArr[0], lineNum)
		if err != nil {
			return err
		}

		*output |= reg << 4

	case "cls":
		if len(lineArr) != 0 {
			return errvals.WrongArgs(1, instruction, len(lineArr), lineNum)
		}

	case "lstr":
		err := LoadString()
		if err != nil {
			return err
		}

	case "print":
		content, _, err := parseString()
		if err != nil {
			return err
		}

		for _, val := range content {
			EvalInstruction("li", output, []string{"r0", fmt.Sprintf("%02x", val)})
			outputFile = append(outputFile, fmt.Sprintf("%04x\n", *output))

			EvalInstruction("li", output, []string{"r3", "0x80"})
			outputFile = append(outputFile, fmt.Sprintf("%04x\n", *output))

			EvalInstruction("or", output, []string{"r0", "r0", "r3"})
			outputFile = append(outputFile, fmt.Sprintf("%04x\n", *output))

			whitespaces += 3
		}

		*output = 0
		EvalInstruction("li", output, []string{"r0", "0"})
		outputFile = append(outputFile, fmt.Sprintf("%04x\n", *output))

	}

	return nil
}

func LoadString() error {
	content, original, err := parseString()
	if err != nil {
		return err
	}

	line := strings.Replace(code, original, "", -1)
	lineArr := getValues(line)
	lineArr = lineArr[1:]

	//fmt.Println(line, content)

	if !utils.IsRegister(lineArr[0]) {
		return errvals.NotRegister(lineArr[0], lineNum)
	}

	if !utils.IsRegister(lineArr[1]) {
		return errvals.NotRegister(lineArr[1], lineNum)
	}

	reg1, reg2 := lineArr[0], lineArr[1]

	instructions := []string{"li", "sb", "li", "sub"}

	for _, char := range strings.Replace(utils.ReverseString(content), "\"", "", -1) {

		for i, val := range instructions {
			output := tokens[val]

			sequence := getSequence(val, char, reg1, reg2, i)

			EvalInstruction(val, &output, sequence)

			outputFile = append(outputFile, fmt.Sprintf("%04x\n", output))
		}

		whitespaces += 4
	}

	lineNum++

	return nil
}

func getSequence(instruction string, char rune, reg1, reg2 string, i int) []string {
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

func parseString() (string, string, error) {
	line := strings.TrimSpace(code)

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
