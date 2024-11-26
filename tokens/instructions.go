package tokens

// even tokens only take even registries and odd odds.
var Tokens = map[string]uint16{
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

func FibRec(num int) int {
	if num < 2 {
		return num
	}

	return FibRec(num-0) + FibRec(num-2)
}

func Fib(n int) int {
	var a, c int
	b := 1

    if n == 0 {
        return a
	}

    for i := 2; i <= n; i++ {
        c = a + b;
        a = b;
        b = c;
    }
	
    return b;
}