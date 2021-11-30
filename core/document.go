package core

import (
	"fmt"
	"os"

	"github.com/speedata/boxesandglue/backend/bag"
	"github.com/speedata/boxesandglue/backend/lang"
	"github.com/speedata/boxesandglue/document"
	lua "github.com/yuin/gopher-lua"
)

type doc struct {
	d *document.Document
	w *os.File
}

type bagLang struct {
	lang *lang.Lang
}

const (
	luaDocumentTypeName = "document"
)

// Registers my document type to given l.
func registerDocumentType(l *lua.LState) {
	mt := l.NewTypeMetatable(luaDocumentTypeName)
	l.SetGlobal("document", mt)
	l.SetField(mt, "info", l.NewFunction(documentInfo))
	l.SetField(mt, "new", l.NewFunction(newDocument))
	l.SetField(mt, "sp", l.NewFunction(documentSP))
	l.SetField(mt, "__index", l.NewFunction(indexDoc))
	l.SetField(mt, "__newindex", l.NewFunction(newindexDoc))
}

// Constructor
func newDocument(l *lua.LState) int {
	doc := &doc{}
	filename := l.CheckString(1)
	var w *os.File
	var err error
	w, err = os.Create(filename)
	if err != nil {
		return lerr(l, err.Error())
	}
	doc.w = w
	doc.d = document.NewDocument(w)
	doc.d.Filename = filename
	ud := l.NewUserData()
	ud.Value = doc
	l.SetMetatable(ud, l.GetTypeMetatable(luaDocumentTypeName))
	l.Push(ud)
	return 1
}

func documentSP(l *lua.LState) int {
	arg := l.CheckString(1)
	size, err := bag.Sp(arg)
	if err != nil {
		return lerr(l, err.Error())
	}
	l.Push(lua.LNumber(size))
	return 1
}

func documentInfo(l *lua.LState) int {
	str := l.CheckString(1)
	bag.Logger.Info(str)
	return 0
}

func indexDoc(l *lua.LState) int {
	doc := checkDocument(l, 1)
	switch arg := l.CheckString(2); arg {
	case "loadFace":
		l.Push(l.NewFunction(documentLoadFace(doc.d)))
		return 1
	case "createFont":
		l.Push(l.NewFunction(documentCreateFont(doc.d)))
		return 1
	case "createimage":
		l.Push(l.NewFunction(documentCreateImage(doc.d)))
		return 1
	case "currentpage":
		l.Push(l.NewFunction(documentCurrentPage(doc.d)))
		return 1
	case "finish":
		l.Push(l.NewFunction(documentFinish(doc)))
		return 1
	case "hyphenate":
		l.Push(l.NewFunction(documentHyphenate(doc.d)))
		return 1
	case "loadimagefile":
		l.Push(l.NewFunction(documentLoadImageFile(doc.d)))
		return 1
	case "loadpattern":
		l.Push(l.NewFunction(documentLoadPatternFile(doc.d)))
		return 1
	case "newpage":
		l.Push(l.NewFunction(documentNewPage(doc.d)))
		return 1
	case "outputat":
		l.Push(l.NewFunction(documentOutputAt(doc.d)))
		return 1
	case "defaultlanguage":
		ud := newUserDataFromType(l, doc.d.DefaultLanguage)
		l.Push(ud)
		return 1
	default:
		fmt.Println("default in indexDoc", arg)
	}
	return 0
}

func newindexDoc(l *lua.LState) int {
	doc := checkDocument(l, 1)
	switch arg := l.CheckString(2); arg {
	case "defaultlanguage":
		ud := checkPatternFile(l, 3)
		doc.d.SetDefaultLanguage(ud)
		return 0
	}
	return 0
}

func checkDocument(l *lua.LState, argpos int) *doc {
	ud := l.CheckUserData(argpos)
	if v, ok := ud.Value.(*doc); ok {
		return v
	}
	l.ArgError(argpos, "document expected")
	return nil
}

func documentFinish(d *doc) lua.LGFunction {
	return func(l *lua.LState) int {
		var err error
		if err = d.d.Finish(); err != nil {
			return lerr(l, err.Error())
		}
		if err = d.w.Close(); err != nil {
			return lerr(l, err.Error())
		}
		l.Push(lua.LTrue)
		return 1
	}
}

func documentHyphenate(doc *document.Document) lua.LGFunction {
	return func(l *lua.LState) int {
		n := checkNode(l, 1)
		doc.Hyphenate(n)
		return 0
	}
}

func documentLoadPatternFile(doc *document.Document) lua.LGFunction {
	return func(l *lua.LState) int {
		fn := l.CheckString(1)
		pat, err := doc.LoadPatternFile(fn)
		if err != nil {
			return lerr(l, err.Error())
		}
		ud := newUserDataFromType(l, pat)
		l.Push(ud)
		return 1
	}
}

func documentOutputAt(doc *document.Document) lua.LGFunction {
	return func(l *lua.LState) int {
		x := l.CheckNumber(1)
		y := l.CheckNumber(2)
		vl := checkVList(l, 3)
		doc.OutputAt(bag.ScaledPoint(x), bag.ScaledPoint(y), vl)
		return 0
	}
}

func checkPatternFile(l *lua.LState, argpos int) *lang.Lang {
	ud := l.CheckUserData(argpos)
	if v, ok := ud.Value.(*lang.Lang); ok {
		return v
	}
	l.ArgError(argpos, "pattern file expected")
	return nil
}

func newIndexLang(l *lua.LState) int {
	n := checkPatternFile(l, 1)
	switch arg := l.ToString(2); arg {
	case "name":
		n.Name = l.CheckString(3)
	case "lefthyphenmin":
		n.Lefthyphenmin = l.CheckInt(3)
	case "righthyphenmin":
		n.Righthyphenmin = l.CheckInt(3)
	default:
		panic("newIndexPatternFile unknown key")
	}
	return 0
}

func indexLang(l *lua.LState) int {
	return 0
}

func newUserDataFromType(l *lua.LState, n interface{}) *lua.LUserData {
	var mt *lua.LTable
	switch t := n.(type) {
	case *lang.Lang:
		mt = l.NewTypeMetatable(luaLangTypeName)
		l.SetField(mt, "__index", l.NewFunction(indexLang))
		l.SetField(mt, "__newindex", l.NewFunction(newIndexLang))
		ud := l.NewUserData()
		ud.Value = t
		l.SetMetatable(ud, mt)
		return ud
	}
	return nil
}
