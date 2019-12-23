package xontrols

import "github.com/louisevanderlith/squareroot/ctx"

//NomadController is the simplest form of controller. [GET]
type Nomad interface {
	Get(ctx.Requester) (int, interface{})
}
