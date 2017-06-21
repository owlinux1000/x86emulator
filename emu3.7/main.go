package main

import (
    "fmt"
    "os"
    . "./emulator"
    . "./instruction"
    . "./emulator_function"
)

var registers_name = [8]string{"EAX", "ECX", "EDX", "EBX", "ESP", "EBP", "ESI", "EDI"}

func read_binary(emu *Emulator, filename string) {
    
    f, err := os.Open(filename)
    
    if err != nil {
        fmt.Printf("%sファイルが開けません\n", os.Args[1])
        os.Exit(1)
    }
    
    buf := make([]byte, 0x200)
    n, _ := f.Read(buf)
    
    for i := 0; i < n; i++ {
        emu.Memory[i+0x7c00] = buf[i]
    }
    f.Close()
}

func create_emu(eip uint32, esp uint32) (emu *Emulator) {

    emu = &Emulator{}
    emu.Eip = eip
    emu.Registers[ESP] = esp
    
    return emu
}

func dump_registers(emu *Emulator) {
    
    for i, v := range emu.Registers {
        fmt.Printf("%s = %08x\n", registers_name[i], v)
    }
    
    fmt.Printf("EIP = %08x\n", emu.Eip)
}

func main() {

    if len(os.Args) != 2 {
        fmt.Println("usage: ./px86 filename")
        return
    }
    
    emu := create_emu(uint32(0x7c00), uint32(0x7c00))

    read_binary(emu, os.Args[1])
    
    Init_instructions()

    for i := emu.Eip; i < MEMORY_SIZE; i++ {
        
        code := uint8(Get_code8(emu, 0))
        fmt.Printf("EIP = %X, Code = %02X\n", emu.Eip, code)

        if Instructions[code] == nil {
            fmt.Printf("\n\nNot Implemented: %x\n", code)
            break
        }

        Instructions[code](emu)
        
        if emu.Eip == 0x00 {
            fmt.Printf("\n\nend of program.\n\n")
            break
        }
    }
    dump_registers(emu)
}
