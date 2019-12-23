package ctx

type ServeFunc func(Requester) (int, interface{})
