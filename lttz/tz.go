package tz

import (
	"time"
)

type Lttz struct {
	LocationName string
	TimeLayout   string
	Time         time.Time
}

type ILttz interface {
	SetLN(ln string) *Lttz
	SetTL(tl string) *Lttz
	SetTime(t time.Time) *Lttz
	Apply() time.Time
	ApplyWithLayout() string
	Modify(time *time.Time)
}

func (l *Lttz) SetLN(ln string) *Lttz {
	l.LocationName = ln
	return l
}

func (l *Lttz) SetTL(lt string) *Lttz {
	l.TimeLayout = lt
	return l
}

func (l *Lttz) SetTime(t time.Time) *Lttz {
	l.Time = t
	return l
}

func (l *Lttz) Apply() time.Time {
	local, _ := time.LoadLocation(l.LocationName)
	return l.Time.In(local)
}

func (l *Lttz) Modify(times *time.Time) {
	local, _ := time.LoadLocation(l.LocationName)
	*times = l.Time.In(local)
}

func (l *Lttz) ApplyWithLayout() string {
	local, _ := time.LoadLocation(l.LocationName)
	if l.TimeLayout != "" {
		return l.Time.In(local).Format(l.TimeLayout)
	}
	return ""
}

func NewLttz() *Lttz {
	return &Lttz{}
}
