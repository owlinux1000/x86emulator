package emulator

const MEMORY_SIZE = 1024 * 1024
type Register int
const (
    EAX Register = iota
    ECX
    EDX
    EBX
    ESP
    EBP
    ESI
    EDI
    REGISTERS_COUNT
)

type Emulator struct {
    Registers [REGISTERS_COUNT]uint32
    Eflags uint32
    Memory [MEMORY_SIZE]byte
    Eip uint32
}
