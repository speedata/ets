package core

import (
	"github.com/speedata/boxesandglue/backend/image"
	"github.com/speedata/boxesandglue/document"
	"github.com/speedata/boxesandglue/pdfbackend/pdf"
	lua "github.com/yuin/gopher-lua"
)

const (
	luaImageFileTypeName = "imagefile"
	luaImageTypeName     = "image"
	luaLangTypeName      = "lang"
)

func checkImage(l *lua.LState, argpos int) *image.Image {
	ud := l.CheckUserData(argpos)
	if v, ok := ud.Value.(*image.Image); ok {
		return v
	}
	l.ArgError(argpos, "imagefile expected")
	return nil
}

func checkImagefile(l *lua.LState, argpos int) *pdf.Imagefile {
	ud := l.CheckUserData(argpos)
	if v, ok := ud.Value.(*pdf.Imagefile); ok {
		return v
	}
	l.ArgError(argpos, "imagefile expected")
	return nil
}

func documentLoadImageFile(doc *document.Document) lua.LGFunction {
	return func(l *lua.LState) int {
		fn := l.CheckString(1)
		dif, err := doc.LoadImageFile(fn)
		if err != nil {
			return lerr(l, err.Error())
		}
		mt := l.NewTypeMetatable(luaImageFileTypeName)
		l.SetField(mt, "__index", l.NewFunction(indexImageFile))
		ud := l.NewUserData()
		ud.Value = dif
		l.SetMetatable(ud, mt)
		l.Push(ud)
		return 1
	}
}

func documentCreateImage(doc *document.Document) lua.LGFunction {
	return func(l *lua.LState) int {
		imgf := checkImagefile(l, 1)
		img := doc.CreateImage(imgf)
		mt := l.NewTypeMetatable(luaImageTypeName)
		l.SetField(mt, "__index", l.NewFunction(indexImage))
		// l.SetField(mt, "__newindex", l.NewFunction(newIndexImage))
		ud := l.NewUserData()
		ud.Value = img
		l.SetMetatable(ud, mt)
		l.Push(ud)
		return 1
	}
}

// func nweindexImage(l *lua.LState) int {
// 	img := checkImage(l, 1)
// 	switch l.CheckString(2) {
// 	case "width":

// 		img.ImageFile.W = bag.ScaledPoint(l.CheckNumber(3))
// 	}
// 	return 0
// }

func indexImage(l *lua.LState) int {
	return 0
}

func indexImageFile(l *lua.LState) int {
	imgf := checkImagefile(l, 1)
	arg := l.ToString(2)
	switch arg {
	case "format":
		l.Push(lua.LString(imgf.Format))
		return 1
	case "numberOfPages":
		l.Push(lua.LNumber(imgf.NumberOfPages))
		return 1
	case "filename":
		l.Push(lua.LString(imgf.Filename))
		return 1
	}
	return 0
}
