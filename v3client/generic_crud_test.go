package v3client

// Method to allow override of the transport interactions.

//func (crud *GenericCRUD[TParent, TIdent, T]) GetUsing(f func(ctx context.Context, opCtx transport.ObjectFetchSpec[T], c *transport.HttpTransport) (T, error)) {
//	// TODO: make this stackable.
//	crud.doGet = f
//}
//
//func (crud *GenericCRUD[TParent, TIdent, T]) UpdateUsing(f func(ctx context.Context, opCtx transport.ObjectUpsertSpec[T], c *transport.HttpTransport) (T, error)) {
//	crud.doUpdate = f
//}
