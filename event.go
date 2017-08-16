package event

type Event interface {
	isevent()
}

type Insert struct {
	ID     int
	P      []byte
	Q0, Q1 int64
}
type Delete struct {
	ID     int
	Q0, Q1 int64
}
type Select struct {
	ID     int
	Q0, Q1 int64
}
type SetOrigin struct {
	ID    int
	Q0    int64
	Exact bool
}
type Fill struct {
}
type Scroll struct {
}
type Redraw struct {
}
type Sweep struct {
}
type Move struct {
}

func (Insert) isevent()    {}
func (Delete) isevent()    {}
func (Select) isevent()    {}
func (SetOrigin) isevent() {}
func (Fill) isevent()      {}
func (Scroll) isevent()    {}
func (Redraw) isevent()    {}
func (Sweep) isevent()     {}
func (Move) isevent()      {}
