/////////////////////////////////////////////////////////////////////////
// WIP-----------------
/////////////////////////////////////////////////////////////////////////

package prototype

import "github.com/ohayouarmaan/proton/lexer"

type Prototype struct {
	Memmory  map[string]lexer.Literal_Value
	Previous *Prototype
}

func New() Prototype {
	return Prototype{
		Memmory:  make(map[string]lexer.Literal_Value),
		Previous: nil,
	}
}
func (p *Prototype) Look(var_name string) *lexer.Literal_Value {
	if p.Previous == nil {
		if val, ok := p.Memmory[var_name]; ok {
			return &val
		} else {
			return nil
		}
	} else {
		val, ok := p.Memmory[var_name]
		if ok {
			return &val
		} else {
			return p.Previous.Look(var_name)
		}
	}
}

func (p *Prototype) Update(var_name string, val lexer.Literal_Value) *lexer.Literal_Value {
	if old_value, ok := p.Memmory[var_name]; ok {
		p.Memmory[var_name] = val
		return &old_value
	} else {
		if p.Previous != nil {
			return p.Previous.Update(var_name, val)
		}
		return nil
	}
}

func (p *Prototype) Set(var_name string, val lexer.Literal_Value) {
	p.Memmory[var_name] = val
}
