package core

import (
	"fmt"

	bagnode "github.com/speedata/boxesandglue/backend/node"
	lua "github.com/yuin/gopher-lua"
)

func lerr(l *lua.LState, errormessage string) int {
	l.SetTop(0)
	l.Push(lua.LFalse)
	l.Push(lua.LString(errormessage))
	return 2
}

// for debugging
func stackDump(l *lua.LState) {
	fmt.Println("-------stack------")
	top := l.GetTop()
	fmt.Println("Top", top)
	for i := 1; i <= top; i++ {
		obj := l.Get(i)
		fmt.Printf("%2d: %-8.8s", i, obj.Type())
		switch obj.Type() {
		case lua.LTString:
			fmt.Println(obj)
		case lua.LTUserData:
			fmt.Println()
			ud := obj.(*lua.LUserData)
			if hl, ok := ud.Value.(*bagnode.HList); ok {
				fmt.Println("hlist id:", hl.ID)
			}
		}
		fmt.Println()
	}
}
