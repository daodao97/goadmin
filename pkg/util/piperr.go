package util

type PieFun = func(input interface{}) (interface{}, error)

func NewPipErr() *pipErr {
	return &pipErr{
		part: []PieFun{},
	}
}

type pipErr struct {
	part []PieFun
}

func (p *pipErr) Wrap(f PieFun) *pipErr {
	p.part = append(p.part, f)
	return p
}

func (p *pipErr) Run() (interface{}, error) {
	var result interface{}
	var err error
	for _, f := range p.part {
		result, err = f(result)
		if err != nil {
			return nil, err
		}
	}
	return result, nil
}
