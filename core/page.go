package core

import (
	"fmt"

	"github.com/speedata/boxesandglue/document"
	lua "github.com/yuin/gopher-lua"
)

const luaPageTypeName = "page"

type documentPage struct {
	page *document.Page
}

func checkPage(l *lua.LState, argpos int) *documentPage {
	ud := l.CheckUserData(argpos)
	if v, ok := ud.Value.(*documentPage); ok {
		return v
	}
	fmt.Printf("ud.Value %#v\n", ud.Value)
	l.ArgError(argpos, "page expected")
	return nil
}

func newUserdataPage(l *lua.LState, p *document.Page) *lua.LUserData {
	dp := &documentPage{page: p}
	ud := l.NewUserData()
	ud.Value = dp
	mt := l.NewTypeMetatable(luaPageTypeName)
	l.SetField(mt, "__index", l.NewFunction(pageIndex))
	l.SetMetatable(ud, mt)
	return ud
}

func pageIndex(l *lua.LState) int {
	p := checkPage(l, 1)
	switch l.CheckString(2) {
	case "shipout":
		l.Push(l.NewFunction(pageShipoutFunc(p)))
		return 1
	}
	return 0
}

func pageShipoutFunc(p *documentPage) lua.LGFunction {
	return func(l *lua.LState) int {
		p.page.Shipout()
		return 0
	}
}

func documentCurrentPage(doc *document.Document) lua.LGFunction {
	return func(l *lua.LState) int {
		l.Push(newUserdataPage(l, doc.CurrentPage))
		return 1
	}
}

func documentNewPage(l *lua.LState) int {
	d := checkDocument(l, 1)
	p := d.d.NewPage()
	l.Push(newUserdataPage(l, p))
	return 1
}
