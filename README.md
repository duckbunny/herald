###Herald

Is the core to the declaration and pooling services intended to be used by microservices.

[![GoDoc](https://godoc.org/github.com/duckbunny/herald?status.svg)](https://godoc.org/github.com/duckbunny/herald)


# herald
--
    import "github.com/duckbunny/herald"

The herald package is intended to provide the tools to make your microservice
public to other declared microservices.

This tool facilitates two functions.


### Pool

The pool interface allows the microservice to support many API routers through a
pool registry.


### Declaration

The declaration interface allows the microservice to declare itself in a
registry of services, so other services can consume that configuration and
automate common actions. By utilizing a Declaration registry a microservice can
support many systems.

A package can support one or both of the interfaces, and the interface to be
used will be determined by environment variables or flags.


## Usage

```go
var DeclarationTypes map[string]Declaration = make(map[string]Declaration)
```
DeclarationTypes are a collection of Declaration Interfaces identified by
strings associated with a declaration mechanism.

```go
var PoolTypes map[string]Pool = make(map[string]Pool)
```
PoolTypes are a collection of Pool Interfaces identified by strings associated
with a pooling mechanism.

#### func  AddDeclaration

```go
func AddDeclaration(key string, d Declaration)
```
Add a single declaration.

#### func  AddDeclarations

```go
func AddDeclarations(ds map[string]Declaration)
```
Add Multiple declarations at once.

#### func  AddPool

```go
func AddPool(key string, p Pool)
```
Add a single pool.

#### func  AddPools

```go
func AddPools(ps map[string]Pool)
```
Add Multiple pools at once.

#### type Declaration

```go
type Declaration interface {
	Declare(*service.Service) error
	Get(*service.Service) error
	Init() error
}
```

Declaration defines an interface that will broadcast the microservice
definition, for other services to digest.

#### type Herald

```go
type Herald struct {
	Pool
	Declaration
	Service *service.Service
}
```

Herald is a wrapper structure of the pooling and declaration interfaces that
will be implemented based on flags or environement variables.

#### func  This

```go
func This() (*Herald, error)
```
Helper function to build herald for the currently running microservice.

#### func (*Herald) Declare

```go
func (h *Herald) Declare() error
```
Declare this microservice definition.

#### func (*Herald) Get

```go
func (h *Herald) Get(s *service.Service) error
```
Get a foreign microservice definition.

#### func (*Herald) Init

```go
func (h *Herald) Init() error
```
Wrapper to init Pool and Declare if they are set. This allows them to consume
flags if necessary.

#### func (*Herald) StartPool

```go
func (h *Herald) StartPool() error
```
Add this service to pool of microservices.

#### func (*Herald) StopPool

```go
func (h *Herald) StopPool() error
```
Remove this service from pool of microservices.

#### type Pool

```go
type Pool interface {
	Start(*service.Service) error
	Stop(*service.Service) error
	Init() error
}
```

Pool defines an interface that will add and remove a microservice to a pool of
microservices.
