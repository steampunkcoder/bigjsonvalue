package bigjsonvalue

// Kind enumerates the kind of underlying value being held
// by the wrapped interface{} value
type Kind uint

// Kind enumeration constants
const (
	Nil Kind = iota
	Bool
	String
	Int64
	Uint64
	Float64
	BigInt
	BigFloat
	lastKind
	// insert new enums before lastKind, lastKind MUST ALWAYS BE LAST
)

var kindNames = [...]string{
	"Nil",
	"Bool",
	"String",
	"Int64",
	"Uint64",
	"Float64",
	"BigInt",
	"BigFloat",
}

// String implements fmt.Stringer interface for Kind
func (k Kind) String() string {
	return kindNames[k]
}
