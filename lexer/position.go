package lexer

type position struct {
	origin        int
	currentLine   int
	currentColumn int
	nextLine      int
	nextColumn    int
	prevLine      int
	prevColumn    int
}

func newPosition(origin int) *position {
	return &position{origin, origin, origin, origin, origin, origin, origin}
}

func (p *position) next() {
	p.prevLine = p.currentLine
	p.prevColumn = p.currentColumn
	p.currentLine = p.nextLine
	p.currentColumn = p.nextColumn
	p.nextColumn++
}

func (p *position) newline() {
	p.next()
	p.nextLine++
	p.nextColumn = p.origin
}

func (p *position) back() {
	p.nextLine = p.currentLine
	p.nextColumn = p.currentColumn
	p.currentLine = p.prevLine
	p.currentColumn = p.prevColumn
}
