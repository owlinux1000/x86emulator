package modrm

import (
    "os"
    "fmt"
    . "../emulator"
    . "../emulator_function"
)

type ModRM struct {
    Mod uint8
    Rm uint8
    Opecode uint8
    Reg_index uint8
    Sib uint8
    Disp8 int8
    Disp32 uint32
};

func Parse_modrm(emu *Emulator, modrm *ModRM) {
    
    var code uint8

    if emu == nil || modrm == nil {
        
    }
    
    code = uint8(Get_code8(emu, 0))
    
    modrm.Mod = ((code & 0xc0) >> 6)
    modrm.Opecode = ((code & 0x38) >> 3)
    modrm.Reg_index = modrm.Opecode
    modrm.Rm = code & 0x7

    emu.Eip += 1

    if modrm.Mod != 3 && modrm.Rm == 4 {
        modrm.Sib = uint8(Get_code8(emu, 0))
        emu.Eip += 1
    }

    if (modrm.Mod == 0 && modrm.Rm == 5) || modrm.Mod == 2 {
        modrm.Disp32 = uint32(Get_sign_code32(emu, 0))
        modrm.Disp8 = int8(modrm.Disp8)
        emu.Eip += 4
    } else if modrm.Mod == 1 {
        modrm.Disp8 = int8(Get_sign_code8(emu, 0))
        modrm.Disp32 = uint32(modrm.Disp8)
        emu.Eip += 1
    }
}

func Calc_memory_address(emu *Emulator, modrm *ModRM) (result uint32) {
    
    if modrm.Mod == 0 {
        if modrm.Rm == 4 {
            fmt.Println("not implemented ModRM mod = 0, rm = 4")
            os.Exit(0)
        } else if modrm.Rm == 5 {
            result = modrm.Disp32
        } else {
            result = uint32(Get_register32(emu, int(modrm.Rm)))
        }
    } else if modrm.Mod == 1 {
        if modrm.Rm == 4 {
            fmt.Println("not implemented ModRM mod = 2, rm = 4")
            os.Exit(0)
        } else {
            result = uint32(Get_register32(emu, int(modrm.Rm))) + modrm.Disp32
        }
    } else {
        fmt.Println("not implemented ModRM mod = 3")
        os.Exit(0)
    }
    return result
}

func Set_rm32(emu *Emulator, modrm *ModRM, value uint32) {
    if modrm.Mod == 3 {
        Set_register32(emu, uint32(modrm.Rm), value)
    } else {
        address := Calc_memory_address(emu, modrm)
        Set_memory32(emu, address, value)
    }
}

func Get_rm32(emu *Emulator, modrm *ModRM) (result uint32) {
    if modrm.Mod == 3 {
        result = Get_register32(emu, int(modrm.Rm))
    } else {
        address := Calc_memory_address(emu, modrm)
        result = Get_memory32(emu, address)
    }
    return result
}

func Set_r32(emu *Emulator, modrm *ModRM, value uint32) {
    Set_register32(emu, uint32(modrm.Reg_index), value)
}

func Get_r32(emu *Emulator, modrm *ModRM) uint32 {
    return Get_register32(emu, int(modrm.Reg_index))
}
