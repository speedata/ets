package core

import (
	lua "github.com/yuin/gopher-lua"
)

// Dothings opens the Lua file and executes it
func Dothings(luafile string) error {

	l := lua.NewState()
	registerDocumentType(l)
	registerNodeType(l)
	if err := l.DoFile(luafile); err != nil {
		return err
	}

	return nil
}
