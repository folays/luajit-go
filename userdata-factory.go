package luajit

// factory is a user-friendly wrapper to bring together metadata + userdata.
type factory struct {
	mt *metatable
}

func (L *State) factory_new() (fac *factory) {
	fac = &factory{}

	fac.mt = L.metatable_newRef()

	L.metatable_cb_set_new(fac.mt, L._userdata_new)
	L.metatable_cb_set_gc(fac.mt, L._userdata_gc)

	return
}

func (L *State) factory_produce(fac *factory, values ...any) {
	L.userdata_new(fac.mt, values...)
}

func (L *State) factory_produceRef(fac *factory, values ...any) (u *userdata) {
	return L.userdata_newRef(fac.mt, values...)
}

func (L *State) factory_checkFatal(fac *factory, index Index) (v any) {
	return L.userdata_checkFatal(index, fac.mt)
}
