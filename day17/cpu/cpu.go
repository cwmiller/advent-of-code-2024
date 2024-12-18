package cpu

import (
	"fmt"
	"math"
)

type Opcode int

const (
	Adv Opcode = 0
	Bxl Opcode = 1
	Bst Opcode = 2
	Jnz Opcode = 3
	Bxc Opcode = 4
	Out Opcode = 5
	Bdv Opcode = 6
	Cdv Opcode = 7
)

type Cpu struct {
	halted      bool
	a, b, c, ip int
	rom         Rom
	output      OutputHandler
}

func (cpu *Cpu) Halted() bool {
	return cpu.halted
}

func (cpu *Cpu) SetA(a int) {
	cpu.a = a
}

func (cpu *Cpu) SetB(b int) {
	cpu.b = b
}

func (cpu *Cpu) SetC(c int) {
	cpu.c = c
}

func (cpu *Cpu) Step() {
	// CPU HALTs once it tries to execute something outside of ROM
	if cpu.ip >= len(cpu.rom) {
		cpu.halted = true
		return
	}

	if cpu.halted {
		return
	}

	opcode := Opcode(cpu.rom[cpu.ip])
	operand := cpu.rom[cpu.ip+1]

	cpu.ip += 2

	switch opcode {
	case Adv:
		cpu.a /= int(math.Pow(2, float64(cpu.comboOperand(operand))))
	case Bxl:
		cpu.b ^= operand
	case Bst:
		cpu.b = cpu.comboOperand(operand) % 8
	case Jnz:
		if cpu.a > 0 {
			cpu.ip = operand
		}
	case Bxc:
		cpu.b ^= cpu.c
	case Out:
		cpu.output(cpu.comboOperand(operand) % 8)
	case Bdv:
		cpu.b = cpu.a / int(math.Pow(2, float64(cpu.comboOperand(operand))))
	case Cdv:
		cpu.c = cpu.a / int(math.Pow(2, float64(cpu.comboOperand(operand))))
	}
}

func (cpu *Cpu) Reset() {
	cpu.ip = 0
	cpu.halted = false
}

// Return value of a combo operand
func (cpu *Cpu) comboOperand(operand int) int {
	switch operand {
	case 0, 1, 2, 3:
		return operand
	case 4:
		return cpu.a
	case 5:
		return cpu.b
	case 6:
		return cpu.c
	case 7:
		panic("Combo operand 7 not implemented")
	default:
		panic(fmt.Sprintf("Unhandled combo operand: %d", operand))
	}

}

type OutputHandler func(val int)

func NewCpu(rom Rom, output OutputHandler) *Cpu {
	return &Cpu{
		output: output,
		rom:    rom,
	}
}

type Rom []int
