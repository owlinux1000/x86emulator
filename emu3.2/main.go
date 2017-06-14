package main

import (
    "fmt"
    "os"
)

type Register int

type instruction_func_t func(*Emulator)

type Emulator struct {
    registers [REGISTERS_COUNT]uint32
    eflags uint32
    memory [MEMORY_SIZE]byte
    eip uint32
}

const MEMORY_SIZE = 1024 * 1024

const (
    EAX Register = iota
    ECX
    EDX
    EBX
    ESP
    ESI
    EDI
    REGISTERS_COUNT
)

var registers_name = [8]string{"EAX", "ECX", "EDX", "EBX", "ESP", "EBP", "ESI", "EDI"}
var instructions = [256]instruction_func_t{}

func create_emu(eip uint32, esp uint32) (emu *Emulator) {

    emu = &Emulator{}
    emu.eip = eip
    emu.registers[ESP] = esp
    
    return emu
}


func dump_registers(emu *Emulator) {
    
    for i, v := range emu.registers {
        fmt.Printf("%s = %08x\n", registers_name[i], v)
    }
    fmt.Printf("EIP = %08x\n", emu.eip)
    
}

func get_code8(emu *Emulator, index int) uint32{
    return uint32(emu.memory[int(emu.eip) + index])
}

func get_sign_code8(emu *Emulator, index int) int32{
    return int32(emu.memory[int(emu.eip) + index])
}

func get_code32(emu *Emulator, index int) (ret uint32) {
    var i uint
    for i = 0; i < 4; i++ {
        ret |= get_code8(emu, index + int(i)) << (i * 8);
    }
    return ret;
}

func get_sign_code32(emu *Emulator, index int) int32 {
    return int32(get_code32(emu, index))
}

func mov_r32_imm32(emu *Emulator) {
    reg := get_code8(emu, 0) - 0xb8;
    value := get_code32(emu, 1);
    emu.registers[reg] = value;
    emu.eip += 5
}

func short_jump(emu *Emulator) {
    diff := int8(get_sign_code8(emu, 1))
    emu.eip += uint32(diff + 2)
}

func near_jump(emu *Emulator) {
    diff := get_sign_code32(emu, 1)
    emu.eip += uint32(diff + 5)
}

func init_instructions() {
    // 0xb8 + regsiter number()
    for i := 0; i < 8; i++ {
        instructions[0xb8 + i] = mov_r32_imm32;
    }
    instructions[0xe9] = near_jump
    instructions[0xeb] = short_jump

}

func main() {

    if len(os.Args) != 2 {
        fmt.Println("usage: ./px86 filename")
        return
    }

    emu := create_emu(uint32(0x7c00), uint32(0x7c00))
    
    f, err := os.Open(os.Args[1])
    if err != nil {
        fmt.Printf("%sファイルが開けません\n", os.Args[1])
    }
    
    buf := make([]byte, 0x200)
    n, _ := f.Read(buf)
    for i := 0; i < n; i++ {
        emu.memory[i+0x7c00] = buf[i]
    }
    f.Close()
    
    init_instructions()

    for i := emu.eip; i < MEMORY_SIZE; i++ {
        code := uint8(get_code8(emu, 0))
        fmt.Printf("EIP = %X, Code = %02X\n", emu.eip, code)

        if instructions[code] == nil {
            fmt.Printf("\n\nNot Implemented: %x\n", code)
            break
        }

        instructions[code](emu)

        if emu.eip == 0x00 {
            fmt.Printf("\n\nend of program.\n\n")
            break
        }
    }
    dump_registers(emu)
}
