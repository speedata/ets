package core

import (
	"fmt"

	"github.com/speedata/boxesandglue/backend/bag"
	bagnode "github.com/speedata/boxesandglue/backend/node"
	lua "github.com/yuin/gopher-lua"
)

const (
	luaNodeTypeName        = "node"
	luaDiscNodeTypeName    = "discnode"
	luaGlueNodeTypeName    = "gluenode"
	luaGlyphNodeTypeName   = "glyphnode"
	luaHlistNodeTypeName   = "hlistnode"
	luaImageNodeTypeName   = "imagenode"
	luaLangNodeTypeName    = "langnode"
	luaPenaltyNodeTypeName = "penaltynode"
	luaVlistNodeTypeName   = "vlistnode"
)

/*
	Common data structures and methods for all kind of nodes
*/

func checkNode(l *lua.LState, argpos int) bagnode.Node {
	ud := l.CheckUserData(argpos)
	if v, ok := ud.Value.(bagnode.Node); ok {
		return v
	}
	l.ArgError(argpos, "node expected")
	return nil
}

type node struct {
}

// Registers my node type to given l.
func registerNodeType(l *lua.LState) {
	mt := l.NewTypeMetatable(luaNodeTypeName)
	l.SetGlobal("node", mt)
	l.SetField(mt, "new", l.NewFunction(newNode))
	l.SetField(mt, "debug", l.NewFunction(debugNode))
	l.SetField(mt, "hpack", l.NewFunction(nodeHpack))
	l.SetField(mt, "insertafter", l.NewFunction(nodeInsertAfter))
	l.SetField(mt, "insertbefore", l.NewFunction(nodeInsertBefore))
	l.SetField(mt, "simplelinebreak", l.NewFunction(nodeSimpleLinebreak))
}

func debugNode(l *lua.LState) int {
	n := checkNode(l, 1)
	bagnode.Debug(n)
	return 0
}

func nodeHpack(l *lua.LState) int {
	n := checkNode(l, 1)
	hl := bagnode.Hpack(n)
	ud := newUserDataFromNode(l, hl)
	l.Push(ud)
	return 1
}

func nodeInsertAfter(l *lua.LState) int {
	var head, cur bagnode.Node
	if l.Get(1) == lua.LNil {
	} else {
		head = checkNode(l, 1)
	}
	if l.Get(2) == lua.LNil {
	} else {
		cur = checkNode(l, 2)
	}
	ins := checkNode(l, 3)
	newhead := bagnode.InsertAfter(head, cur, ins)
	l.Push(newUserDataFromNode(l, newhead))
	return 1
}

func nodeInsertBefore(l *lua.LState) int {
	var head, cur bagnode.Node
	if l.Get(1) == lua.LNil {
	} else {
		head = checkNode(l, 1)
	}
	if l.Get(2) == lua.LNil {
	} else {
		cur = checkNode(l, 2)
	}
	ins := checkNode(l, 3)
	newhead := bagnode.InsertBefore(head, cur, ins)
	l.Push(newUserDataFromNode(l, newhead))
	return 1
}

func nodeSimpleLinebreak(l *lua.LState) int {
	n := checkNode(l, 1)
	tbl := l.CheckTable(2)

	settings := bagnode.LinebreakSettings{}
	l.Push(tbl.RawGetString("hsize"))
	hsize := l.CheckNumber(3)

	l.Push(tbl.RawGetString("lineheight"))
	linehight := l.CheckNumber(4)
	settings.HSize = bag.ScaledPoint(hsize)
	settings.LineHeight = bag.ScaledPoint(linehight)
	vl := bagnode.SimpleLinebreak(n.(*bagnode.HList), settings)
	l.Push(newUserDataFromNode(l, vl))
	return 1
}

func newNode(l *lua.LState) int {
	switch l.CheckString(1) {
	case "disc":
		l.Push(newUserDataFromNode(l, bagnode.NewDisc()))
		return 1
	case "glue":
		l.Push(newUserDataFromNode(l, bagnode.NewGlue()))
		return 1
	case "glyph":
		l.Push(newUserDataFromNode(l, bagnode.NewGlyph()))
		return 1
	case "hlist":
		l.Push(newUserDataFromNode(l, bagnode.NewHList()))
		return 1
	case "image":
		l.Push(newUserDataFromNode(l, bagnode.NewImage()))
		return 1
	case "lang":
		l.Push(newUserDataFromNode(l, bagnode.NewLang()))
		return 1
	case "penalty":
		l.Push(newUserDataFromNode(l, bagnode.NewPenalty()))
		return 1
	case "vlist":
		l.Push(newUserDataFromNode(l, bagnode.NewVList()))
		return 1
	default:
		panic("nyi (newNode)")
	}
}

func checkIsNode(l *lua.LState) bool {
	ud := l.CheckUserData(1)
	return bagnode.IsNode(ud.Value)
}

func newUserDataFromNode(l *lua.LState, n bagnode.Node) *lua.LUserData {
	var mt *lua.LTable
	switch n.(type) {
	case *bagnode.Disc:
		mt = l.NewTypeMetatable(luaDiscNodeTypeName)
		l.SetField(mt, "__index", l.NewFunction(discIndex))
		l.SetField(mt, "__newindex", l.NewFunction(discNewIndex))
	case *bagnode.Glue:
		mt = l.NewTypeMetatable(luaGlueNodeTypeName)
		l.SetField(mt, "__index", l.NewFunction(glueIndex))
		l.SetField(mt, "__newindex", l.NewFunction(glueNewIndex))
	case *bagnode.Glyph:
		mt = l.NewTypeMetatable(luaGlyphNodeTypeName)
		l.SetField(mt, "__index", l.NewFunction(glyphIndex))
		l.SetField(mt, "__newindex", l.NewFunction(glyphNewIndex))
	case *bagnode.HList:
		mt = l.NewTypeMetatable(luaHlistNodeTypeName)
		l.SetField(mt, "__index", l.NewFunction(hlistIndex))
		l.SetField(mt, "__newindex", l.NewFunction(hlistNewIndex))
	case *bagnode.Image:
		mt = l.NewTypeMetatable(luaImageNodeTypeName)
		l.SetField(mt, "__index", l.NewFunction(imageNodeIndex))
		l.SetField(mt, "__newindex", l.NewFunction(imageNodeNewIndex))
	case *bagnode.Lang:
		mt = l.NewTypeMetatable(luaLangNodeTypeName)
		l.SetField(mt, "__index", l.NewFunction(langNodeIndex))
		l.SetField(mt, "__newindex", l.NewFunction(langNodeNewIndex))
	case *bagnode.Penalty:
		mt = l.NewTypeMetatable(luaPenaltyNodeTypeName)
		l.SetField(mt, "__index", l.NewFunction(penaltyNodeIndex))
		l.SetField(mt, "__newindex", l.NewFunction(penaltyNodeNewIndex))
	case *bagnode.VList:
		mt = l.NewTypeMetatable(luaVlistNodeTypeName)
		l.SetField(mt, "__index", l.NewFunction(vlistIndex))
		l.SetField(mt, "__newindex", l.NewFunction(vlistNewIndex))
	default:
		panic("nyi newUserDataFromNode")
	}
	ud := l.NewUserData()
	ud.Value = n
	l.SetMetatable(ud, mt)
	return ud
}

/*

	Disc nodes

*/

func checkDisc(l *lua.LState, argpos int) *bagnode.Disc {
	ud := l.CheckUserData(argpos)
	if v, ok := ud.Value.(*bagnode.Disc); ok {
		return v
	}
	l.ArgError(argpos, "disc node expected")
	return nil
}

func discIndex(l *lua.LState) int {
	n := checkDisc(l, 1)
	switch arg := l.ToString(2); arg {
	case "next":
		var other bagnode.Node
		if other = n.Next(); other == nil {
			return 0
		}
		l.Push(newUserDataFromNode(l, other))
		return 1
	case "prev":
		var other bagnode.Node
		if other = n.Next(); other == nil {
			return 0
		}
		l.Push(newUserDataFromNode(l, other))
		return 1
	}
	return 0
}

func discNewIndex(l *lua.LState) int {
	n := checkDisc(l, 1)
	switch arg := l.ToString(2); arg {
	case "next":
		if l.Get(3) == lua.LNil {
			n.SetNext(nil)
		} else {
			n.SetNext(checkNode(l, 3))
		}
		return 0
	case "prev":
		if l.Get(3) == lua.LNil {
			n.SetNext(nil)
		} else {
			n.SetNext(checkNode(l, 3))
		}
		return 0
	default:
		fmt.Println("newindex", arg)
	}
	return 0
}

/*

   Glyph nodes

*/
func checkGlyph(l *lua.LState, argpos int) *bagnode.Glyph {
	ud := l.CheckUserData(argpos)
	if v, ok := ud.Value.(*bagnode.Glyph); ok {
		return v
	}
	l.ArgError(argpos, "glyph expected")
	return nil
}

func glyphNewIndex(l *lua.LState) int {
	n := checkGlyph(l, 1)
	switch arg := l.ToString(2); arg {
	case "next":
		if l.Get(3) == lua.LNil {
			n.SetNext(nil)
		} else {
			n.SetNext(checkNode(l, 3))
		}
		return 0
	case "prev":
		if l.Get(3) == lua.LNil {
			n.SetNext(nil)
		} else {
			n.SetNext(checkNode(l, 3))
		}
		return 0
	case "codepoint":
		arg := l.CheckNumber(3)
		n.Codepoint = int(arg)
	case "components":
		arg := l.CheckString(3)
		n.Components = arg
	case "font":
		arg := checkFont(l, 3)
		n.Font = arg
	case "width":
		wd := l.CheckNumber(3)
		n.Width = bag.ScaledPoint(wd)
	default:
		fmt.Println("newindex", arg)
		_ = n
	}
	return 0
}
func glyphIndex(l *lua.LState) int {
	n := checkGlyph(l, 1)
	switch arg := l.ToString(2); arg {
	case "next":
		var other bagnode.Node
		if other = n.Next(); other == nil {
			return 0
		}
		l.Push(newUserDataFromNode(l, other))
		return 1
	case "prev":
		var other bagnode.Node
		if other = n.Next(); other == nil {
			return 0
		}
		l.Push(newUserDataFromNode(l, other))
		return 1
	case "width":
		l.Push(lua.LNumber(n.Width))
		return 1
	case "components":
		l.Push(lua.LString(n.Components))
		return 1
	default:
		fmt.Println("arg", arg)
	}
	return 0
}

/*

	Glue nodes

*/

func checkGlue(l *lua.LState, argpos int) *bagnode.Glue {
	ud := l.CheckUserData(argpos)
	if v, ok := ud.Value.(*bagnode.Glue); ok {
		return v
	}
	l.ArgError(argpos, "glue expected")
	return nil
}

func glueIndex(l *lua.LState) int {
	n := checkGlue(l, 1)
	switch arg := l.ToString(2); arg {
	case "next":
		var other bagnode.Node
		if other = n.Next(); other == nil {
			return 0
		}
		l.Push(newUserDataFromNode(l, other))
		return 1
	case "prev":
		var other bagnode.Node
		if other = n.Next(); other == nil {
			return 0
		}
		l.Push(newUserDataFromNode(l, other))
		return 1
	case "width":
		l.Push(lua.LNumber(n.Width))
		return 1
	default:
		l.ArgError(2, fmt.Sprintf("unknown field %s in glue", arg))
		return 0
	}
}

func glueNewIndex(l *lua.LState) int {
	n := checkGlue(l, 1)
	switch arg := l.ToString(2); arg {
	case "next":
		if l.Get(3) == lua.LNil {
			n.SetNext(nil)
		} else {
			n.SetNext(checkNode(l, 3))
		}
		return 0
	case "prev":
		if l.Get(3) == lua.LNil {
			n.SetNext(nil)
		} else {
			n.SetNext(checkNode(l, 3))
		}
		return 0
	case "width":
		arg := l.CheckNumber(3)
		n.Width = bag.ScaledPoint(arg)
	default:
		l.ArgError(2, fmt.Sprintf("unknown field %s in glue", arg))
		return 0

	}
	return 0
}

/*

	HList nodes

*/
func checkHlist(l *lua.LState, argpos int) *bagnode.HList {
	ud := l.CheckUserData(argpos)
	if v, ok := ud.Value.(*bagnode.HList); ok {
		return v
	}
	l.ArgError(argpos, "hlist expected")
	return nil
}

func hlistIndex(l *lua.LState) int {
	var other bagnode.Node
	n := checkHlist(l, 1)
	switch arg := l.ToString(2); arg {
	case "list":
		if other = n.List; other == nil {
			return 0
		}
		l.Push(newUserDataFromNode(l, other))
		return 1
	case "next":
		if other = n.Next(); other == nil {
			return 0
		}
		l.Push(newUserDataFromNode(l, other))
		return 1
	case "prev":
		if other = n.Next(); other == nil {
			return 0
		}
		l.Push(newUserDataFromNode(l, other))
		return 1
	default:
		l.ArgError(2, fmt.Sprintf("unknown field %s in hlist", arg))
		return 0
	}
}

func hlistNewIndex(l *lua.LState) int {
	n := checkHlist(l, 1)
	switch arg := l.ToString(2); arg {
	case "next":
		if l.Get(3) == lua.LNil {
			n.SetNext(nil)
		} else {
			n.SetNext(checkNode(l, 3))
		}
		return 0
	case "prev":
		if l.Get(3) == lua.LNil {
			n.SetNext(nil)
		} else {
			n.SetNext(checkNode(l, 3))
		}
		return 0
	case "list":
		n.List = checkNode(l, 3)
		return 0
	default:
		l.ArgError(2, fmt.Sprintf("unknown field %s in hlist", arg))
		return 0
	}
}

/*

	Image nodes

*/
func imageNodeIndex(l *lua.LState) int {
	n := checkImageNode(l, 1)
	switch arg := l.ToString(2); arg {
	case "next":
		var other bagnode.Node
		if other = n.Next(); other == nil {
			return 0
		}
		l.Push(newUserDataFromNode(l, other))
		return 1
	case "prev":
		var other bagnode.Node
		if other = n.Next(); other == nil {
			return 0
		}
		l.Push(newUserDataFromNode(l, other))
		return 1
	default:
		l.ArgError(2, fmt.Sprintf("unknown field %s in image", arg))
		return 0
	}
}
func imageNodeNewIndex(l *lua.LState) int {
	n := checkImageNode(l, 1)
	switch arg := l.ToString(2); arg {
	case "next":
		if l.Get(3) == lua.LNil {
			n.SetNext(nil)
		} else {
			n.SetNext(checkNode(l, 3))
		}
		return 0
	case "prev":
		if l.Get(3) == lua.LNil {
			n.SetNext(nil)
		} else {
			n.SetNext(checkNode(l, 3))
		}
		return 0
	case "img":
		n.Img = checkImage(l, 3)
	case "width":
		n.Width = bag.ScaledPoint(l.CheckNumber(3))
	case "height":
		n.Height = bag.ScaledPoint(l.CheckNumber(3))
	default:
		l.ArgError(2, fmt.Sprintf("unknown field %s in image", arg))
		return 0
	}
	return 0
}

func checkImageNode(l *lua.LState, argpos int) *bagnode.Image {
	ud := l.CheckUserData(argpos)
	if v, ok := ud.Value.(*bagnode.Image); ok {
		return v
	}
	l.ArgError(argpos, "image expected")
	return nil
}

/*

	Lang nodes

*/

func checkLangNode(l *lua.LState, argpos int) *bagnode.Lang {
	ud := l.CheckUserData(argpos)
	if v, ok := ud.Value.(*bagnode.Lang); ok {
		return v
	}
	l.ArgError(argpos, "lang expected")
	return nil
}

func langNodeNewIndex(l *lua.LState) int {
	n := checkLangNode(l, 1)
	switch arg := l.ToString(2); arg {
	case "next":
		if l.Get(3) == lua.LNil {
			n.SetNext(nil)
		} else {
			n.SetNext(checkNode(l, 3))
		}
		return 0
	case "prev":
		if l.Get(3) == lua.LNil {
			n.SetNext(nil)
		} else {
			n.SetNext(checkNode(l, 3))
		}
		return 0
	case "lang":
		pf := checkPatternFile(l, 3)
		n.Lang = pf
	}
	return 0
}

func langNodeIndex(l *lua.LState) int {
	n := checkLangNode(l, 1)
	switch arg := l.ToString(2); arg {
	case "next":
		var other bagnode.Node
		if other = n.Next(); other == nil {
			return 0
		}
		l.Push(newUserDataFromNode(l, other))
		return 1
	case "prev":
		var other bagnode.Node
		if other = n.Next(); other == nil {
			return 0
		}
		l.Push(newUserDataFromNode(l, other))
		return 1
	case "name":
		l.Push(lua.LString(n.Lang.Name))
		return 1
	}
	return 0
}

/*

   Penalty nodes

*/

func checkPenaltyNode(l *lua.LState, argpos int) *bagnode.Penalty {
	ud := l.CheckUserData(argpos)
	if v, ok := ud.Value.(*bagnode.Penalty); ok {
		return v
	}
	l.ArgError(argpos, "penalty node expected")
	return nil
}

func penaltyNodeNewIndex(l *lua.LState) int {
	n := checkPenaltyNode(l, 1)
	switch arg := l.ToString(2); arg {
	case "next":
		if l.Get(3) == lua.LNil {
			n.SetNext(nil)
		} else {
			n.SetNext(checkNode(l, 3))
		}
		return 0
	case "prev":
		if l.Get(3) == lua.LNil {
			n.SetNext(nil)
		} else {
			n.SetNext(checkNode(l, 3))
		}
		return 0
	case "penalty":
		n.Penalty = l.CheckInt(3)
	case "flagged":
		n.Flagged = l.CheckBool(3)
	case "width":
		wd := l.CheckNumber(3)
		n.Width = bag.ScaledPoint(wd)
	}
	return 0
}

func penaltyNodeIndex(l *lua.LState) int {
	n := checkPenaltyNode(l, 1)
	switch arg := l.ToString(2); arg {
	case "next":
		var other bagnode.Node
		if other = n.Next(); other == nil {
			return 0
		}
		l.Push(newUserDataFromNode(l, other))
		return 1
	case "prev":
		var other bagnode.Node
		if other = n.Next(); other == nil {
			return 0
		}
		l.Push(newUserDataFromNode(l, other))
		return 1
	case "penalty":
		l.Push(lua.LNumber(n.Penalty))
		return 1
	case "flagged":
		l.Push(lua.LBool(n.Flagged))
		return 1
	case "width":
		l.Push(lua.LNumber(n.Width))
		return 1
	}
	return 0
}

/*

	VList nodes

*/

func checkVList(l *lua.LState, argpos int) *bagnode.VList {
	ud := l.CheckUserData(argpos)
	if v, ok := ud.Value.(*bagnode.VList); ok {
		return v
	}
	l.ArgError(argpos, "vlist expected")
	return nil
}

func vlistNewIndex(l *lua.LState) int {
	n := checkVList(l, 1)
	switch arg := l.ToString(2); arg {
	case "next":
		if l.Get(3) == lua.LNil {
			n.SetNext(nil)
		} else {
			n.SetNext(checkNode(l, 3))
		}
		return 0
	case "prev":
		if l.Get(3) == lua.LNil {
			n.SetNext(nil)
		} else {
			n.SetNext(checkNode(l, 3))
		}
		return 0
	case "list":
		newnode := checkNode(l, 3)
		n.List = newnode
	}
	return 0
}

func vlistIndex(l *lua.LState) int {
	n := checkVList(l, 1)
	switch arg := l.ToString(2); arg {
	case "next":
		var other bagnode.Node
		if other = n.Next(); other == nil {
			return 0
		}
		l.Push(newUserDataFromNode(l, other))
		return 1
	case "prev":
		var other bagnode.Node
		if other = n.Next(); other == nil {
			return 0
		}
		l.Push(newUserDataFromNode(l, other))
		return 1
	case "list":
		if n.List == nil {
			return 0
		}
		l.Push(newUserDataFromNode(l, n.List))
		return 1
	default:
		fmt.Println("arg", arg)
	}
	return 0
}
