package main

import (
	"bufio"
	"fmt"
	"mips/mod"
	"mips/tokens"
	"mips/utils"
	"os"
	"strconv"
	"strings"
)

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

	generated := mod.New()

	for scanner.Scan() {
		generated.Code = strings.Split(scanner.Text(), "#")[0]
		line := strings.ToLower(strings.TrimSpace(generated.Code))
		lineArr := utils.GetValues(line)

		if line == "" {
			/* whitespaces++
			lineNum++
			outputFile = append(outputFile, "0000\n") */
			generated.LineNum++
			continue
		}

		if line[len(line)-1] == ':' {
			tag := line[:len(line)-1]
			jump, ok := generated.Jumps[tag]

			if ok {
				if jump.Line != nil {
					fmt.Printf("Error: \"%s\" already defined\n", tag)
					return
				}

				jump.Line = new(uint16)
				*jump.Line = generated.Whitespaces

			} else {
				jump = mod.Jmp{Line: new(uint16), Callers: []int{}}
				*jump.Line = generated.Whitespaces
			}

			generated.Jumps[tag] = jump

			generated.LineNum++

			continue
		}

		if strings.Contains(line, "*") {
			newLines, err := strconv.Atoi(strings.Split(line, "*")[0])
			if err != nil {
				fmt.Println("Syntax error at line:", generated.LineNum, "\nInstruction", line, "not recognized")
				return
			}

			generated.Whitespaces += uint16(newLines)
			generated.LineNum++
			generated.OutputFile = append(generated.OutputFile, line+"\n")

			continue
		}

		instruction := lineArr[0]
		lineArr = lineArr[1:]
		output, ok := tokens.Tokens[instruction]

		if !ok {
			fmt.Printf("Syntax error: instruction %s, \nAt line: %d\n%s not recognized.\n", instruction, generated.LineNum, line)
			return
		}

		err := generated.EvalInstruction(instruction, &output, lineArr)
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		if !(instruction == "lstr" || instruction == "print") {
			generated.OutputFile = append(generated.OutputFile, fmt.Sprintf("%04x\n", output))

		}

		generated.LineNum++
		generated.Whitespaces++
	}

	//var jumps map[string]uint16

	for key, value := range generated.Jumps {

		if value.Line == nil {
			times := "time"
			if len(value.Callers) > 1 {
				times += "s"
			}

			fmt.Printf("error: tag \"%s\" referenced %d %s was never declared\n", key, len(value.Callers), times)
			return
		}

		if len(value.Callers) == 0 {
			fmt.Printf("error: tag \"%s\" never used.\n", key)
			return
		}

		for _, i := range value.Callers {
			newVal := string(generated.OutputFile[i][0])

			jp := fmt.Sprintf("%02x\n", *value.Line)

			if len(jp) <= 3 {
				newVal += "0" + jp

			} else {
				newVal += jp
			}

			generated.OutputFile[i] = newVal
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

	for _, line := range generated.OutputFile {
		f.WriteString(line)
	}

	f.Sync()
}
