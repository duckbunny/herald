// Copyright © 2015 Jason Smith <jasonrichardsmith@gmail.com>.
//
// Use of this source code is governed by the LGPL-3
// license that can be found in the LICENSE file.

/*
The herald package is intended to provide the tools to make your
microservice public to other declared microservices.

This tool facilitates two functions.

Pool

The pool interface allows the microservice to support many API routers through
a pool registry.

Declaration

The declaration interface allows the microservice to declare itself in a registry
of services, so other services can consume that configuration and automate
common actions.  By utilizing a Declaration registry a microservice can support
many systems.


A package can support one or both of the interfaces, and the interface to be
used will be determined by environment variables or flags.

*/
package herald
