package instruction

import (
    "os"    
    "fmt"
    . "../emulator"
    . "../emulator_function"
    . "../modrm"
)


type Instruction_func_t func(*Emulator)
var Instructions = [256]Instruction_func_t{}

func Mov_r32_imm32(emu *Emulator) {
    reg := Get_code8(emu, 0) - 0xb8;
    value := Get_code32(emu, 1);
    emu.Registers[reg] = value;
    emu.Eip += 5
}

func Mov_r32_rm32(emu *Emulator) {
    emu.Eip += 1
    var modrm ModRM
    Parse_modrm(emu, &modrm)
    rm32 := Get_rm32(emu, &modrm)
    Set_r32(emu, &modrm, rm32)
}

func Mov_rm32_r32(emu *Emulator) {
    emu.Eip += 1
    var modrm ModRM
    Parse_modrm(emu, &modrm)
    r32 := Get_r32(emu, &modrm)
    Set_rm32(emu, &modrm, r32)
}

func Mov_rm32_imm32(emu *Emulator) {
    emu.Eip += 1
    var modrm ModRM
    Parse_modrm(emu, &modrm)
    value := Get_code32(emu, 0)
    emu.Eip += 4
    Set_rm32(emu, &modrm, value)
}

func Sub_rm32_imm8(emu *Emulator, modrm *ModRM) {
    rm32 := int32(Get_rm32(emu, modrm))
    imm8 := int32(Get_sign_code8(emu, 0))
    emu.Eip += 1
    Set_rm32(emu, modrm, uint32(rm32 - imm8))
}

func Code_83(emu *Emulator) {
    emu.Eip += 1
    var modrm ModRM
    Parse_modrm(emu, &modrm)

    switch modrm.Opecode {
    case 5:
        Sub_rm32_imm8(emu, &modrm)
    default:
        fmt.Printf("not implemented: 83 /%d\n", modrm.Opecode)
        os.Exit(1)
    }
}

func Inc_rm32(emu *Emulator, modrm *ModRM) {
    value := Get_rm32(emu, modrm)
    Set_rm32(emu, modrm, value + 1)
}

func Code_ff(emu *Emulator) {
    emu.Eip += 1
    var modrm ModRM
    Parse_modrm(emu, &modrm)

    switch modrm.Opecode {
    case 0:
        Inc_rm32(emu, &modrm);
    default:
        fmt.Printf("not implemented: FF /%d\n", modrm.Opecode)
        os.Exit(1)
    }
}

func Add_rm32_r32(emu *Emulator) {
    emu.Eip += 1
    var modrm ModRM
    Parse_modrm(emu, &modrm)
    r32 := Get_r32(emu, &modrm)
    rm32 := Get_rm32(emu, &modrm)
    Set_rm32(emu, &modrm, rm32 + r32)
}

func Short_jump(emu *Emulator) {
    diff := int8(Get_sign_code8(emu, 1))
    emu.Eip += uint32(diff + 2)
}

func Near_jump(emu *Emulator) {
    diff := Get_sign_code32(emu, 1)
    emu.Eip += uint32(diff + 5)
}


func Init_instructions() {

    Instructions[0x01] = Add_rm32_r32
    Instructions[0x83] = Code_83
    Instructions[0x89] = Mov_rm32_r32
    Instructions[0x8b] = Mov_r32_rm32
    
    // 0xb8 + regsiter number
    for i := 0; i < 8; i++ {
        Instructions[0xb8 + i] = Mov_r32_imm32
    }

    Instructions[0xc7] = Mov_rm32_imm32
    Instructions[0xe9] = Near_jump
    Instructions[0xeb] = Short_jump
    Instructions[0xff] = Code_ff
}
