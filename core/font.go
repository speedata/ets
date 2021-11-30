package core

import (
	"github.com/speedata/boxesandglue/backend/bag"
	"github.com/speedata/boxesandglue/backend/font"
	"github.com/speedata/boxesandglue/document"
	"github.com/speedata/boxesandglue/pdfbackend/pdf"
	lua "github.com/yuin/gopher-lua"
)

const (
	luaFaceTypeName = "face"
	luaFontTypeName = "font"
)

func documentLoadFace(doc *document.Document) lua.LGFunction {
	return func(l *lua.LState) int {
		fn := l.CheckString(1)
		f, err := doc.LoadFace(fn, 0)
		if err != nil {
			return lerr(l, err.Error())
		}

		mt := l.NewTypeMetatable(luaFaceTypeName)
		ud := l.NewUserData()
		ud.Value = f
		l.SetMetatable(ud, mt)
		l.Push(ud)
		return 1
	}
}

func checkFace(l *lua.LState, argpos int) *pdf.Face {
	ud := l.CheckUserData(argpos)
	if v, ok := ud.Value.(*pdf.Face); ok {
		return v
	}
	l.ArgError(argpos, "face expected")
	return nil
}

func checkFont(l *lua.LState, argpos int) *font.Font {
	ud := l.CheckUserData(argpos)
	if v, ok := ud.Value.(*font.Font); ok {
		return v
	}
	l.ArgError(argpos, "font expected")
	return nil
}

// func indexFace(l *lua.LState) int {
// 	f := checkFace(l, 1)
// 	arg := l.CheckString(2)
// 	switch arg {
// 	case "font":

// 	}
// }

func documentCreateFont(doc *document.Document) lua.LGFunction {
	return func(l *lua.LState) int {
		face := checkFace(l, 1)
		size := l.CheckNumber(2)
		fnt := doc.CreateFont(face, bag.ScaledPoint(size))
		mt := l.NewTypeMetatable(luaFontTypeName)
		l.SetField(mt, "__index", l.NewFunction(indexFont))
		ud := l.NewUserData()
		ud.Value = fnt
		l.SetMetatable(ud, mt)
		l.Push(ud)
		return 1
	}
}

func fontShape(fnt *font.Font, fntObj lua.LValue) lua.LGFunction {
	return func(l *lua.LState) int {
		str := l.CheckString(1)
		tbl := l.NewTable()
		for _, glyph := range fnt.Shape(str) {
			glyphtbl := l.NewTable()
			glyphtbl.RawSetString("codepoint", lua.LNumber(glyph.Codepoint))
			glyphtbl.RawSetString("advance", lua.LNumber(glyph.Advance))
			glyphtbl.RawSetString("components", lua.LString(glyph.Components))
			glyphtbl.RawSetString("glyph", lua.LNumber(glyph.Glyph))
			glyphtbl.RawSetString("hyphenate", lua.LBool(glyph.Hyphenate))
			glyphtbl.RawSetString("font", fntObj)

			tbl.Append(glyphtbl)
		}
		l.Push(tbl)
		return 1
	}
}

func indexFont(l *lua.LState) int {
	f := checkFont(l, 1)
	fontObj := l.Get(1)
	arg := l.CheckString(2)
	switch arg {
	case "size":
		l.Push(lua.LNumber(f.Size))
		return 1
	case "space":
		l.Push(lua.LNumber(f.Space))
		return 1
	case "stretch":
		l.Push(lua.LNumber(f.SpaceStretch))
		return 1
	case "shrink":
		l.Push(lua.LNumber(f.SpaceShrink))
		return 1
	case "shape":
		l.Push(l.NewFunction(fontShape(f, fontObj)))
		return 1
	}
	return 0
}
