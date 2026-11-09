package interp

type Var interface {
}

type VarInt struct {
	i   int
}

type VarNumber struct {
	f float32
}

type VarBool struct {
	b bool
}

type VarString struct {
	s string
}
