package luajit

import (
	"fmt"
	"log"
	"unsafe"
)

type LuaCTypeId uint16

func (L *State) _push_cdata_int_str(code string) {
	{
		L._stackTops_push(L.GetTop())
		{
			if err := L.LoadString(code); err != nil {
				log.Fatalf("(*State)._push_cdata_int_str(%s) ERROR : %s", code, err)
			}
			if _, err := L._run(0, 1); err != nil {
				log.Fatalf("(*State)._push_cdata_int_str(%s) ERROR : %s", code, err)
			}
		}
		L._stackTops_pop()
	}
	L._stackTops_pop()
}

func (L *State) ToCtype(index Index) (v any) {
	//     GCcdata *cd = cdataV(L->base + idx - 1);
	// #define cdataV(o)	check_exp(tviscdata(o), &gcval(o)->cd)
	//  #define gcval(o)	(gcref((o)->gcr))
	//  #define gcval(o)	((GCobj *)(gcrefu((o)->gcr) & LJ_GCVMASK))
	//

	var (
		ptr  = L.ToPointer(index)
		addr = uintptr(ptr)

		ctypeid = *(*LuaCTypeId)(unsafe.Pointer(addr - 6))
	)

	//{
	//	fmt.Printf("\033[2mPTR: %p\033[22m\n", ptr)
	//
	//	fmt.Printf("%8x | %8x\n",
	//		*(*uint32)(unsafe.Pointer(addr - 16)),
	//		*(*uint32)(unsafe.Pointer(addr - 12)),
	//	)
	//
	//	fmt.Printf("%2x marked | %2x gct | %4x ctypeid | %4x offset | %4x extra\n",
	//		*(*uint8)(unsafe.Pointer(addr - 8)),
	//		*(*uint8)(unsafe.Pointer(addr - 7)),
	//		*(*uint16)(unsafe.Pointer(addr - 6)),
	//		*(*uint16)(unsafe.Pointer(addr - 4)),
	//		*(*uint16)(unsafe.Pointer(addr - 2)),
	//	)
	//
	//	fmt.Printf("%x\n", unsafe.Slice((*byte)(unsafe.Pointer(addr+0)), 8))
	//	fmt.Printf("%x\n", unsafe.Slice((*byte)(unsafe.Pointer(addr+8)), 8))
	//}

	switch {
	case ctypeid == 0x16:
		v = *(*LuaCTypeId)(unsafe.Pointer(addr))
		fmt.Printf("(*State).ToCtype() ctypeid/[0x%x %d] -> v/[0x%x %x]\n", ctypeid, ctypeid, v, v)
	case ctypeid == 0x0b: // CTID_INT64
		v = *(*int64)(unsafe.Pointer(addr))
	case ctypeid == 0x0c: // CTID_UINT64
		v = *(*uint64)(unsafe.Pointer(addr))
	case ctypeid >= 0x65:
		// ok
	default:
		log.Fatalf("(*State).ToCtype() error : unhandled ctypeid 0x%x", ctypeid)
	}

	return
}
