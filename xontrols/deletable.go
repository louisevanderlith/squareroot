package xontrols

import "github.com/louisevanderlith/squareroot/ctx"

//Deleteable handles controls that can delete content [DELETE]
type Deleteable interface {
	Nomad
	Delete(ctx.Requester) (int, interface{})
}
