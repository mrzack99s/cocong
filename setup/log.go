package setup

import (
	"io"
	"log"
	"os"
	"path"

	"github.com/mrzack99s/cocong/vars"
	"gopkg.in/natefinch/lumberjack.v2"
)

type LoggingConfig struct {
	Directory  string
	Filename   string
	MaxSize    int
	MaxBackups int
	MaxAge     int
}

func (c *LoggingConfig) Configure() {
	e, err := os.OpenFile(path.Join(c.Directory, c.Filename), os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}

	vars.SystemLog = log.New(e, "", log.Ldate|log.Ltime)
	vars.SystemLog.SetOutput(c.newRollingFile())

}

func (c *LoggingConfig) newRollingFile() io.Writer {
	if err := os.MkdirAll(c.Directory, 0600); err != nil {
		panic("can't create log directory")
	}

	return &lumberjack.Logger{
		Filename:   path.Join(c.Directory, c.Filename),
		MaxBackups: c.MaxBackups, // files
		MaxSize:    c.MaxSize,    // megabytes
		MaxAge:     c.MaxAge,     // days
	}
}
