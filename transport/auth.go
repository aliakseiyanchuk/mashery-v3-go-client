package transport

type Authorizer interface {
	Authorization() (map[string]string, error)
	Close()
}
