package core

import (
	"github.com/speedata/boxesandglue/backend/bag"
	"github.com/speedata/boxesandglue/backend/font"
	"github.com/speedata/boxesandglue/document"
	"github.com/speedata/boxesandglue/pdfbackend/pdf"
	lua "github.com/yuin/gopher-lua"
)

const (
	luaFaceTypeName       = "face"
	luaFontTypeName       = "font"
	luaFontFamilyTypeName = "fontfamily"
)

func documentLoadFace(doc *document.Document) lua.LGFunction {
	return func(l *lua.LState) int {
		fn := l.CheckTable(1)
		nameValue := fn.RawGetString("name")
		srcValue := fn.RawGetString("source")
		if nameValue.Type() != lua.LTString {
			return lerr(l, "the value of name must be a string")
		}
		if srcValue.Type() != lua.LTString {
			return lerr(l, "the value of source must be a string")
		}
		fs := document.FontSource{
			Name:   nameValue.String(),
			Source: srcValue.String(),
		}
		f, err := doc.LoadFace(&fs)
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
			glyphtbl.RawSetString("isspace", lua.LBool(glyph.IsSpace))
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

// Font families
func documentNewFontfamily(doc *document.Document) lua.LGFunction {
	return func(l *lua.LState) int {
		familyname := l.CheckString(1)
		ff := doc.NewFontFamily(familyname)
		l.Push(newUserdataFontfamily(l, ff))
		return 1
	}
}

func checkFontfamily(l *lua.LState, argpos int) *document.FontFamily {
	ud := l.CheckUserData(argpos)
	if v, ok := ud.Value.(*document.FontFamily); ok {
		return v
	}
	l.ArgError(argpos, "fontfamily expected")
	return nil
}

func newUserdataFontfamily(l *lua.LState, ff *document.FontFamily) *lua.LUserData {
	ud := l.NewUserData()
	ud.Value = ff
	mt := l.NewTypeMetatable(luaFontFamilyTypeName)
	l.SetField(mt, "__index", l.NewFunction(fontfamilyIndex))
	l.SetMetatable(ud, mt)
	return ud
}

func fontfamilyIndex(l *lua.LState) int {
	ff := checkFontfamily(l, 1)
	switch l.CheckString(2) {
	case "addmember":
		l.Push(l.NewFunction(fontfamilyaddmember(ff)))
		return 1
	case "id":
		l.Push(lua.LNumber(ff.ID))
		return 1
	}
	return 0
}

func fontfamilyaddmember(p *document.FontFamily) lua.LGFunction {
	return func(l *lua.LState) int {
		fn := l.CheckTable(1)
		nameValue := fn.RawGetString("name")
		srcValue := fn.RawGetString("source")
		if nameValue.Type() != lua.LTString {
			return lerr(l, "the value of name must be a string")
		}
		if srcValue.Type() != lua.LTString {
			return lerr(l, "the value of source must be a string")
		}
		fs := &document.FontSource{
			Name:   nameValue.String(),
			Source: srcValue.String(),
		}

		weight := l.CheckInt(2)
		stylestring := l.CheckString(3)
		var style document.FontStyle
		switch stylestring {
		case "regular", "normal":
			style = document.FontStyleNormal
		case "italic":
			style = document.FontStyleItalic
		}
		p.AddMember(fs, weight, style)
		return 0
	}
}
