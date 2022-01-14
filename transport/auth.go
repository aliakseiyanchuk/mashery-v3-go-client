package transport

type Authorizer interface {
	HeaderAuthorization() (map[string]string, error)
	QueryStringAuthorization() (map[string]string, error)
	Close()
}
