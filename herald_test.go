package herald

import (
	"errors"
	"testing"

	"github.com/duckbunny/service"
)

func TestSetFlagEnv(t *testing.T) {
	heraldBoth = "test"
	setFlagEnv()
	if heraldPool != "test" && heraldDeclare != "test" {
		t.Error("expecting pool and declare to eb set to test")
	}
}

func TestThis(t *testing.T) {
	_, err := This()
	if err.Error() != "Attempt to utilize unrecognized Pooling mechanism test" {
		t.Error("expected no pooling mechanism error")
	}
	AddPool("test", StubHerald{})
	_, err = This()
	if err.Error() != "Attempt to utilize unrecognized Declaration mechanism test" {
		t.Error("expected no pooling mechanism error")
	}
	AddDeclaration("test", StubHerald{})
	_, err = This()
}

func TestInit(t *testing.T) {
	h := Herald{Pool: StubHerald{Error: true}, Declaration: StubHerald{}}
	err := h.Init()
	if err == nil {
		t.Error("expected init error for pool")
	}
	h = Herald{Pool: StubHerald{}, Declaration: StubHerald{Error: true}}
	err = h.Init()
	if err == nil {
		t.Error("expected init error for declaration")
	}
	h = Herald{Declaration: StubHerald{Error: true}}
	err = h.Init()
	if err == nil {
		t.Error("expected init error for declaration with no pool")
	}
	h = Herald{Pool: StubHerald{}, Declaration: StubHerald{}}
	err = h.Init()
	if err != nil {
		t.Error("expected no error")
	}
}

func TestStartPool(t *testing.T) {
	h := Herald{}
	err := h.StartPool()
	if err != nil {
		t.Error("expecting no error for nil pool")
	}
	h = Herald{Pool: StubHerald{Error: true}}
	err = h.StartPool()
	if err == nil {
		t.Error("expecting error from start pool")
	}
}

func TestStopPool(t *testing.T) {
	h := Herald{}
	err := h.StopPool()
	if err != nil {
		t.Error("expecting no error for nil pool")
	}
	h = Herald{Pool: StubHerald{Error: true}}
	err = h.StopPool()
	if err == nil {
		t.Error("expecting error from start pool")
	}
}

func TestDeclare(t *testing.T) {
	h := Herald{}
	err := h.Declare()
	if err != nil {
		t.Error("expecting no error for nil declaration")
	}
	h = Herald{Declaration: StubHerald{Error: true}}
	err = h.Declare()
	if err == nil {
		t.Error("expecting error from declare")
	}
}

func TestGetService(t *testing.T) {
	s := service.Service{}
	h := Herald{Declaration: StubHerald{}}
	h.GetService(&s)
	if s.Domain != "test" {
		t.Error("expected new service")
	}

}

type StubHerald struct {
	Error bool
}

func (s StubHerald) Start(*service.Service) error {
	if s.Error {
		return errors.New("an error")
	}
	return nil
}

func (s StubHerald) Stop(*service.Service) error {
	if s.Error {
		return errors.New("an error")
	}
	return nil
}

func (s StubHerald) Init() error {
	if s.Error {
		return errors.New("an error")
	}
	return nil
}

func (s StubHerald) Declare(*service.Service) error {
	if s.Error {
		return errors.New("an error")
	}
	return nil
}

func (s StubHerald) GetService(service *service.Service) error {
	service.Domain = "test"
	return nil
}
