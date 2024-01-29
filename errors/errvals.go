package errvals

import "fmt"

func EvenNumber(val uint16) error {
	return fmt.Errorf("error: Registry %d is not an even registry", val)
}

func OddNumber(val uint16) error {
	return fmt.Errorf("error: Registry %d is not an odd registry", val)
}

func WrongArgs(number int, instruction string, length int, lineNum uint16) error {
	return fmt.Errorf("error: Instruction \"%s\" expected %d arguments\nGot %d at line %d", instruction, number, length, lineNum)
}

func NotRegister(val string, lineNum uint16) error {
	return fmt.Errorf("error: %s is not a valid registry\nAt line %d", val, lineNum)
}