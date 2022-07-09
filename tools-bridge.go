package luajit

import (
	"log"
	"reflect"
)

type Bridge struct {
	tBridge reflect.Type
	fac     *factory
}

type IBridgePayload interface {
	_lua_payload_new(L *State) []any
}

func (L *State) bridge_prepare() {
	L.all_br_by_type_bridge = make(map[reflect.Type]*Bridge)
	L.all_br_by_type_shared = make(map[reflect.Type]*Bridge)
}

func (L *State) BridgeCreate(oBridge any) (br *Bridge) {
	var (
		tBridge = reflect.TypeOf(oBridge)
	)

	if tBridge.Kind() != reflect.Struct && tBridge.Kind() != reflect.Pointer {
		log.Fatalf("(*State).BridgeCreate() error : not a [possibly pointer] struct")
	}

	if br = L.all_br_by_type_bridge[tBridge]; br != nil {
		log.Printf("\033[33mWARNING\033[39m BridgeCreate() : bridge already exists")
		return
	}

	br = L.bridge_newRef(tBridge)

	L.all_br_by_type_bridge[tBridge] = br
	return
}

func (L *State) BridgeGet(oBridge any) (br *Bridge) {
	var (
		tBridge = reflect.TypeOf(oBridge)
	)

	br = L.all_br_by_type_bridge[tBridge]
	return
}

// BridgeAdd prepares exporting all `reflect.TypeOf(t)` as __index metamethods
func (L *State) BridgeAdd(oBridge any) {
	L.BridgeAddShared(oBridge, oBridge)
}

// BridgeAddShared
//
// - oBridge is the custom object type having `Lua_*`  methods
// - oShared is the target object type which will be "bridged"
func (L *State) BridgeAddShared(oBridge any, oShared any) {
	var (
		br *Bridge
	)

	if br = L.BridgeGet(oBridge); br == nil {
		br = L.BridgeCreate(oBridge)
	}

	var (
		tShared = reflect.TypeOf(oShared)
	)

	if tShared.Kind() != reflect.Struct && tShared.Kind() != reflect.Pointer {
		log.Fatalf("(*State).BridgeAddShared() error : not a [possibly pointer] struct")
	}

	if L.all_br_by_type_shared[tShared] != nil {
		log.Printf("\033[33mWARNING\033[39m BridgeAddShared() : bridge already shared")
		return
	}

	L.all_br_by_type_shared[tShared] = br
}

func (L *State) bridge_newRef(tBridge reflect.Type) (br *Bridge) {
	br = &Bridge{
		tBridge: tBridge,
	}

	{
		br.fac = L.factory_new()

		{
			L.NewTable()
			L.FuncsAddFilter_tableIndex(-1, tBridge, FilterFuncsPrefix("Lua_"))
			L.metatable_setField(br.fac.mt, "__index")
		}
	}
	return
}

func (L *State) pushGoStruct(v reflect.Value) {
	var (
		vFinal reflect.Value
	)

	for vFinal = v; ; vFinal = vFinal.Elem() {
		if br := L.all_br_by_type_shared[vFinal.Type()]; br != nil {
			L._pushGoStruct_bridge(br, v)
			return
		}

		switch vFinal.Kind() {
		case reflect.Pointer:
			continue
		case reflect.Struct:
			log.Fatalf("(*State).pushGoStruct() error : [.Elem...].Kind() struct unknown : %s", v.Type())
		default:
			log.Fatalf("(*State).pushGoStruct() error : [.Elem...].Kind() is not a Struct")
		}
	}
}

func (L *State) _pushGoStruct_bridge(br *Bridge, v reflect.Value) {
	var (
		values = []any{v.Interface()}
	)

	if vPayload, ok := v.Convert(br.tBridge).Interface().(IBridgePayload); ok == true {
		values = append(values, vPayload._lua_payload_new(L)...)

	}

	L.factory_produce(br.fac, values...)
}
