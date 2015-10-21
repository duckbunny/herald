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

// PoolTypes are a collection of Pool Interfaces identified by strings
// associated with a pooling mechanism.
var PoolTypes map[string]Pool = make(map[string]Pool)

// DeclarationTypes are a collection of Declaration Interfaces identified by
// strings associated with a declaration mechanism.
var DeclarationTypes map[string]Declaration = make(map[string]Declaration)

// Flags to be parsed for setting pool and declaration interfaces.
var heraldPool string
var heraldDeclare string
var heraldBoth string

func init() {
	flag.StringVar(&heraldPool, "herald-p", "", "Herald to microservice pool.")
	flag.StringVar(&heraldDeclare, "herald-d", "", "Herald to microservice declaration.")
	flag.StringVar(&heraldBoth, "herald", "", "Herald to handle both declaration and pooling.")
}

// Herald is a wrapper structure of the pooling and declaration interfaces that will be implemented
// based on flags or environement variables.
type Herald struct {
	Pool
	Declaration
	Service *service.Service
}

// Set declared pool and declaration for this instance.
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

// Helper function to build herald for the currently running microservice.
func This() (*Herald, error) {
	setFlagEnv()
	h := &Herald{}
	if heraldPool != "" {
		if pt, ok := PoolTypes[heraldPool]; !ok {
			err := fmt.Errorf("Attempt to utilize unrecognized Pooling mechanism %v", heraldPool)
			return h, err
		}
		h.Pool = PoolTypes[heraldPool]
	}
	if heraldDeclare != "" {
		if dt, ok := DeclarationTypes[heraldDeclare]; !ok {
			err := fmt.Errorf("Attempt to utilize unrecognized Declaration mechanism %v", heraldPool)
			return h, err
		}
		h.Declaration = DeclarationTypes[heraldDeclare]
	}
	s, err := service.This()
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

// Add this service to pool of microservices.
func (h *Herald) StartPool() error {
	return h.Pool.Start(h.Service)
}

// Remove this service from pool of microservices.
func (h *Herald) StopPool() error {
	return h.Pool.Stop(h.Service)
}

// Declare this microservice definition.
func (h *Herald) Declare() error {
	return h.Declaration.Declare(h.Service)
}

// Get a foreign microservice definition.
func (h *Herald) Get(s *service.Service) error {
	return h.Declaration.Get(s)
}

// Pool defines an interface that will add and remove a microservice
// to a pool of microservices.
type Pool interface {
	Start(*service.Service) error
	Stop(*service.Service) error
	Init() error
}

// Add a single pool.
func AddPool(key string, p Pool) {
	PoolTypes[key] = p
}

// Add Multiple pools at once.
func AddPools(ps map[string]Pool) {
	for key, pool := range ps {
		PoolTypes[key] = pool
	}
}

// Declaration defines an interface that will broadcast the microservice
// definition, for other services to digest.
type Declaration interface {
	Declare(*service.Service) error
	Get(*service.Service) error
	Init() error
}

// Add a single declaration.
func AddDeclaration(key string, d Declaration) {
	DeclarationTypes[key] = d
}

// Add Multiple declarations at once.
func AddDeclarations(ds map[string]Declaration) {
	for key, dec := range ds {
		DeclarationTypes[key] = d
	}
}
