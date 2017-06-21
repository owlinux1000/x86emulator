package emulator_function

import (
    . "../emulator"
)

func Get_code8(emu *Emulator, index int) uint32{
    return uint32(emu.Memory[int(emu.Eip) + index])
}

func Get_sign_code8(emu *Emulator, index int) int32{
    return int32(emu.Memory[int(emu.Eip) + index])
}

func Get_code32(emu *Emulator, index int) (ret uint32) {
    var i uint
    for i = 0; i < 4; i++ {
        ret |= Get_code8(emu, index + int(i)) << (i * 8);
    }
    return ret;
}

func Get_sign_code32(emu *Emulator, index int) int32 {
    return int32(Get_code32(emu, index))
}

func Get_register32(emu *Emulator, index int) uint32 {
    return emu.Registers[index]
}

func Set_register32(emu *Emulator, index uint32, value uint32) {
    emu.Registers[index] = value
}

func Set_memory8(emu *Emulator, address uint32, value uint32) {
    emu.Memory[address] = byte(value & 0xff)
}

func Set_memory32(emu *Emulator, address uint32, value uint32) {
    for i := 0; i < 4; i++ {
        Set_memory8(emu, address + uint32(i), value >> (uint(i) * 8))
    }
}

func Get_memory8(emu *Emulator, address uint32) uint32 {
    return uint32(emu.Memory[address])
}

func Get_memory32(emu *Emulator, address uint32) (ret uint32){
    
    for i := 0; i < 4; i++ {
        ret |= Get_memory8(emu, address + uint32(i)) << (8 * uint32(i))
    }
    return ret
}


func Push32(emu *Emulator, value uint32) {
    
    address := Get_register32(emu, int(ESP)) - 4
    Set_register32(emu, uint32(ESP), address)
    Set_memory32(emu, address, value)
}

func Pop32(emu *Emulator) (ret uint32) {
    
    address := Get_register32(emu, int(ESP))
    ret = Get_memory32(emu, address)
    Set_register32(emu, uint32(ESP), address + 4)
    return ret
    
}
