// Copyright © 2015 Jason Smith <jasonrichardsmith@gmail.com>.
//
// Use of this source code is governed by the LGPL-3
// license that can be found in the LICENSE file.

package herald

import (
	"flag"
	"fmt"
	"os"

	"github.com/duckbunny/service"
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
	flag.StringVar(&heraldPool, "herald-p", os.Getenv("HERALD_POOL"), "Herald to microservice pool.")
	flag.StringVar(&heraldDeclare, "herald-d", os.Getenv("HERALD_DECLARE"), "Herald to microservice declaration.")
	flag.StringVar(&heraldBoth, "herald", os.Getenv("HERALD"), "Herald to handle both declaration and pooling.")
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
	if heraldPool == "" {
		heraldPool = heraldBoth
	}
	if heraldDeclare == "" {
		heraldDeclare = heraldBoth
	}
	return
}

// Helper function to build herald for the currently running microservice.
func This() (*Herald, error) {
	setFlagEnv()
	h := &Herald{}
	if heraldPool != "" {
		if _, ok := PoolTypes[heraldPool]; !ok {
			err := fmt.Errorf(
				"Attempt to utilize unrecognized Pooling mechanism %v",
				heraldPool)
			return h, err
		}
		h.Pool = PoolTypes[heraldPool]
	}
	if heraldDeclare != "" {
		if _, ok := DeclarationTypes[heraldDeclare]; !ok {
			err := fmt.Errorf(
				"Attempt to utilize unrecognized Declaration mechanism %v",
				heraldPool)
			return h, err
		}
		h.Declaration = DeclarationTypes[heraldDeclare]
	}
	s, err := service.This()
	h.Service = s
	return h, err
}

// Wrapper to init Pool and Declare if they are set.  This allows them to
// consume flags if necessary.
func (h *Herald) Init() error {
	if h.Pool != nil {
		err := h.Pool.Init()
		if err != nil {
			return err
		}
	}
	if h.Declaration != nil {
		if h.Declaration.Init() != h.Pool.Init() {
			return h.Declaration.Init()
		}
	}
	return nil
}

// Add this service to pool of microservices.
func (h *Herald) StartPool() error {
	if h.Pool != nil {
		return h.Pool.Start(h.Service)
	}
	return nil
}

// Remove this service from pool of microservices.
func (h *Herald) StopPool() error {
	if h.Pool != nil {
		return h.Pool.Stop(h.Service)
	}
	return nil
}

// Declare this microservice definition.
func (h *Herald) Declare() error {
	if h.Declaration != nil {
		return h.Declaration.Declare(h.Service)
	}
	return nil
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
		DeclarationTypes[key] = dec
	}
}
