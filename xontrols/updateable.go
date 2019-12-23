package xontrols

import "github.com/louisevanderlith/squareroot/ctx"

//Store handles controls that can Update [PUT]
type Updateable interface {
	Nomad
	Update(ctx.Requester) (int, interface{})
}
