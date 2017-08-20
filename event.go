package event

import (
	"fmt"
	"bytes"
	"unicode"
)

type Record interface{
	String() string
//	Equal(Record) bool
	Coalesce(Record) bool
//	Invert(Record) Record
}

type Event interface {
	isevent()
}

type Rec struct{
	ID int
	Kind byte
	Q0, Q1 int64
	N int64
	P []byte
}

func (r Rec) String() string{
	return fmt.Sprintf("%d	%x	%d	%d	%d	%q\n",r.ID, r.Kind, r.Q0,r.Q1,r.N,r.P)
}
func (r Rec) Record() (p []byte, err error){
	return []byte(r.String()), nil
}

func (r Insert) Record() (p []byte, err error){
	r.Kind = 'i'
	r.N = int64(len(p))
	return r.Rec.Record()
}
func (r Delete) Record() (p []byte, err error){
	r.Kind = 'd'
	return r.Rec.Record()
}
func (r Select) Record() (p []byte, err error){
	r.Kind = 's'
	return r.Rec.Record()
}

type Insert struct {
	Rec
}
func space(b byte) int{
	if unicode.IsSpace(rune(b)){
		return 1
	}
	return 0
}

func (e *Insert) Coalesce(v Record) bool{
	if v == nil{
		return false
	}
	switch v := v.(type){
	case *Insert:
		if len(v.P) == 0{
			return true
		}
		if v.ID != e.ID{
			return false
		}
		if (space(v.P[0]) != space(e.P[0])) {
			return false
		}
		if v.Q0 == e.Q1 {
			e.Q1 = v.Q1
			e.P = append(e.P, v.P...)
			return true
		}
	case *Delete:
		if v.ID != e.ID{
			return false
		}
		if len(v.P) == 1 && len(e.P) == 1 && (space(v.P[0]) != space(e.P[0])) {
			return false
		}
		if e.Q0 >= v.Q0 && v.Q1 == e.Q1-1{
		   // 0      3        3       4
			e.Q1 -= v.Q1-v.Q0
			e.P = e.P[:e.Q1]
			return true
		}
	case *Select:
		return false
	}
	return false
}

func (e *Delete) Coalesce(v Record) bool{
	switch v :=v.(type){
	case *Insert:
		return false
		// Works better to just return false here
		//
		//
		if v.ID != e.ID{
			return false
		}
		if v.Q0 != e.Q0{
			return false
		}
		if v.Q1 > e.Q1+1{
			return false
		}
		e.Q0 = v.Q1
		if len(e.P) == 0 || bytes.Equal(e.P[:e.Q0], v.P){
			return false
		}
		copy(e.P, e.P[e.Q0:])
		e.P=e.P[:(e.Q1-e.Q0)+1]
		return true
	case *Delete:
		if v.ID != e.ID{
			return false
		}
		if e.Q0 != v.Q1{
			return false
		}
			e.Q0 = v.Q0
			e.P = append(v.P, e.P...)
		return true
	}
	return false
}
func (e *Select) Coalesce(v Record) bool{
	switch v := v.(type){
	case *Select:
		if v.Q0 == e.Q0 && v.Q1 == e.Q1{
			return true
		}
	}
	return false
}

type Delete struct {
	Rec
}
type Select struct {
	Rec
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
