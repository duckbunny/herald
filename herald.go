// Copyright © 2015 Jason Smith <jasonrichardsmith@gmail.com>.
//
// Use of this source code is governed by the LGPL-3
// license that can be found in the LICENSE file.

package herald

import (
	"flag"
	"github.com/duckbunny/service"
	"os"
)

var PoolTypes map[string]Pool = make(map[string]Pool)

var DeclarationTypes map[string]Declaration = make(map[string]Declaration)

// Flags
var heraldPool string
var heraldDeclare string
var heraldBoth string

func init() {
	flag.StringVar(&heraldPool, "herald-p", "", "Herald to microservice pool.")
	flag.StringVar(&heraldDeclare, "herald-d", "", "Herald to microservice declaration.")
	flag.StringVar(&heraldBoth, "herald", "", "Herald to handle both declaration and pooling.")
}

type Herald struct {
	Pool
	Declaration
	Service *service.Service
}

func setFlagEnv() {
	if !flag.Parsed() {
		flag.Parse()
	}
	// Check flag precedence.
	// Invalidates heraldBoth
	if heraldPool != "" || heraldDeclare != "" {
		if heraldDeclare == "" {
			heraldDeclare = os.Getenv("HERALDDECLARE")
		}
		if heraldPool == "" {
			heraldPool = os.Getenv("HERALDPOOL")
		}
		return
	}
	if heraldBoth == "" {
		heraldBoth = os.Getenv("HERALD")
	}
	heraldPool = heraldBoth
	heraldDeclare = heraldBoth
	return
}

func New(s *service.Service) *Herald {
	return &Herald{Service: s}
}

func This() (*Herald, error) {
	setFlagEnv()
	h := &Herald{}
	if pt, ok := PoolTypes[heraldPool]; ok {
		h.Pool = pt
	}
	if dt, ok := DeclarationTypes[heraldDeclare]; ok {
		h.Declaration = dt
	}
	s, err := service.This()
	if err != nil {
		return h, err
	}
	h.Service = s
	return h, err
}

func (h *Herald) Init() error {
	if h.Pool != nil {
		err := h.Pool.Init()
		if err != nil {
			return err
		}
	}
	if h.Declaration.Init() != h.Pool.Init() {
		return h.Declaration.Init()
	}
	return nil
}

func (h *Herald) StartPool() error {
	return h.Pool.Start(h.Service)
}

func (h *Herald) StopPool() error {
	return h.Pool.Stop(h.Service)
}

func (h *Herald) Declare() error {
	return h.Declaration.Declare(h.Service)
}

type Pool interface {
	Start(*service.Service) error
	Stop(*service.Service) error
	Init() error
}

type Declaration interface {
	Declare(*service.Service) error
	Init() error
}

func AddPool(key string, p Pool) {
	PoolTypes[key] = p
}

func AddDeclaration(key string, d Declaration) {
	DeclarationTypes[key] = d
}
