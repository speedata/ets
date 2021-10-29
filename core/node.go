package core

import (
	"fmt"

	"github.com/speedata/boxesandglue/backend/bag"
	bagnode "github.com/speedata/boxesandglue/backend/node"
	lua "github.com/yuin/gopher-lua"
)

const (
	luaNodeTypeName      = "node"
	luaHlistNodeTypeName = "hlistnode"
	luaVlistNodeTypeName = "vlistnode"
	luaGlyphNodeTypeName = "glyphnode"
	luaImageNodeTypeName = "imagenode"
	luaLangNodeTypeName  = "langnode"
	luaGlueNodeTypeName  = "gluenode"
	luaDiscNodeTypeName  = "discnode"
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
		l.Push(newUserDataFromNode(l, &bagnode.Image{}))
		return 1
	case "lang":
		l.Push(newUserDataFromNode(l, bagnode.NewLang()))
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
	case *bagnode.Image:
		mt = l.NewTypeMetatable(luaImageNodeTypeName)
		l.SetField(mt, "__index", l.NewFunction(imageNodeIndex))
		l.SetField(mt, "__newindex", l.NewFunction(imageNodeNewIndex))
	case *bagnode.Lang:
		mt = l.NewTypeMetatable(luaLangNodeTypeName)
		l.SetField(mt, "__index", l.NewFunction(langNodeIndex))
		l.SetField(mt, "__newindex", l.NewFunction(langNodeNewIndex))
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
func discIndex(l *lua.LState) int {
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
	case "next":
		nd := checkNode(l, 3)
		n.SetNext(nd)
	case "prev":
		nd := checkNode(l, 3)
		n.SetPrev(nd)
	default:
		fmt.Println("newindex", arg)
		_ = n
	}
	return 0
}
func glyphIndex(l *lua.LState) int {
	n := checkGlyph(l, 1)
	switch arg := l.ToString(2); arg {
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
	return 0
}

func glueNewIndex(l *lua.LState) int {
	n := checkGlue(l, 1)
	switch arg := l.ToString(2); arg {
	case "width":
		arg := l.CheckNumber(3)
		n.Width = bag.ScaledPoint(arg)
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
	n := checkHlist(l, 1)
	switch arg := l.ToString(2); arg {
	case "list":
		if n.List == nil {
			return 0
		}
		l.Push(newUserDataFromNode(l, n.List))
		return 1
	case "set_head":
		l.Push(l.NewFunction(setHead))
		return 1
	default:
		fmt.Println("arg", arg)
	}
	return 0
}

func setHead(l *lua.LState) int {
	self := checkHlist(l, 1)
	self.List = checkNode(l, 2)
	return 0
}

/*

	Image nodes

*/
func imageNodeIndex(l *lua.LState) int {
	return 0
}
func imageNodeNewIndex(l *lua.LState) int {
	n := checkImageNode(l, 1)
	switch arg := l.ToString(2); arg {
	case "img":
		n.Img = checkImage(l, 3)
	case "width":
		n.Width = bag.ScaledPoint(l.CheckNumber(3))
	case "height":
		n.Height = bag.ScaledPoint(l.CheckNumber(3))
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
	case "lang":
		pf := checkPatternFile(l, 3)
		n.Lang = pf
	}
	return 0
}

func langNodeIndex(l *lua.LState) int {
	n := checkLangNode(l, 1)
	switch arg := l.ToString(2); arg {
	case "name":
		l.Push(lua.LString(n.Lang.Name))
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
	case "list":
		newnode := checkNode(l, 3)
		n.List = newnode
	}
	return 0
}

func vlistIndex(l *lua.LState) int {
	n := checkVList(l, 1)
	switch arg := l.ToString(2); arg {
	case "list":
		if n.List == nil {
			return 0
		}
		l.Push(newUserDataFromNode(l, n.List))
		return 1
	case "sethead":
		l.Push(l.NewFunction(setVHead))
		return 1
	default:
		fmt.Println("arg", arg)
	}
	return 0
}

func setVHead(l *lua.LState) int {
	self := checkHlist(l, 1)
	self.List = checkNode(l, 2)
	return 0
}
