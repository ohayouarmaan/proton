/////////////////////////////////////////////////////////////////////////
// WIP-----------------
/////////////////////////////////////////////////////////////////////////

package prototype

import "github.com/ohayouarmaan/proton/lexer"

type Prototype struct {
	Memmory  map[string]any
	Previous *Prototype
}

func (p *Prototype) Look(var_name string) *any {
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

func (p *Prototype) Update(var_name string, val lexer.Literal_Value) *any {
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
