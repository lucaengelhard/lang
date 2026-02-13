package parser

type Type interface {
	_type()
}

func parse_type(p *Parser, bp binding_power) Type {
	panic("Types not implemented yet")
}
