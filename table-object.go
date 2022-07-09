package luajit

type Table struct {
	L   *State
	ref Ref
}

// ToTable must later be .Release, otherwise the registry ref will leak.
// TODO: define a runtime.SetFinalizer() on *Table
func (L *State) ToTable(index Index) (t *Table) {
	t = &Table{
		L:   L,
		ref: L.RegistryRefIndex(index),
	}
	return
}

func (t *Table) Release() {
	t.L.RegistryUnref(t.ref)
	t.ref = RefNil
	t.L = nil
}

func (t *Table) GetAssocAll() (rows []map[string]any) {
	var (
		L = t.L
	)
	L.RegistryGet(t.ref)
	defer L.Pop(1) // remove table ref

	if L.IsTable(-1) == false {
		return nil
	}

	L._table_iterate(-1, func() {
		var (
			rowId = L.ToIntSafe(-2)
			_     = rowId
			row   = make(map[string]any)
		)
		L._table_iterate(-1, func() {
			var (
				colName = L.ToStringSafe(-2)
				val     = L.ToAny(-1)
			)
			row[colName] = val
		})
		rows = append(rows, row)
	})

	return
}
