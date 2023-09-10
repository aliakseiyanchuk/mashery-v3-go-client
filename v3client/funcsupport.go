package v3client

type Consumer[T any] func(t1 T)
type BiConsumer[T1, T2 any] func(t1 T1, t2 T2)
type BiConsumerCanErr[T1, T2 any] func(t1 T1, t2 T2) error
type BiConsumerCanErrLocator[TInA, TInB any] func(cl Client) BiConsumerCanErr[TInA, TInB]

type Supplier[T any] func() T
type SupplierCanErr[T any] func() (T, error)
