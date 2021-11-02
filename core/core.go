package core

import (
	"os"
	"path/filepath"

	"github.com/speedata/boxesandglue/backend/bag"
	lua "github.com/yuin/gopher-lua"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Runs the Lua file if it exists. If the exename is "foo", it looks for a Lua file called "foo.lua"
func runDefaultLua(l *lua.LState, exename string) error {
	var extension = filepath.Ext(exename)
	var name = exename[0:len(exename)-len(extension)] + ".lua"
	var err error
	if _, err = os.Stat(name); err != nil {
		return nil
	}
	return l.DoFile(name)
}

func newZapLogger() *zap.SugaredLogger {
	logger, _ := zap.Config{
		Level:       zap.NewAtomicLevelAt(zapcore.DebugLevel),
		Encoding:    "console",
		OutputPaths: []string{"stdout"},
		EncoderConfig: zapcore.EncoderConfig{
			EncodeLevel: zapcore.LowercaseColorLevelEncoder,
			LevelKey:    "level",
			MessageKey:  "message",
		},
	}.Build()
	return logger.Sugar()
}

// Dothings opens the Lua file and executes it
func Dothings(luafile string, exename string) error {
	l := lua.NewState()
	bag.Logger = newZapLogger()
	registerDocumentType(l)
	registerNodeType(l)

	if err := runDefaultLua(l, exename); err != nil {
		return err
	}

	if err := l.DoFile(luafile); err != nil {
		return err
	}

	return nil
}
