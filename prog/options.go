package prog

import (
	"github.com/ambientsound/visp/log"
	"github.com/ambientsound/visp/options"
	"github.com/ambientsound/visp/pkg/library"
	"github.com/ambientsound/visp/topbar"
)

func (v *Visp) optionChanged(key string) {
	switch key {
	case options.LogFile:
		logFile := options.GetString(options.LogFile)
		overwrite := options.GetBool(options.LogOverwrite)
		if len(logFile) == 0 {
			break
		}
		err := log.Configure(logFile, overwrite)
		if err != nil {
			log.Errorf("log configuration: %s", err)
			break
		}
		log.Infof("Note: log file will be backfilled with existing log")
		log.Infof("Writing debug log to %s", logFile)

	case options.Topbar:
		config := options.GetString(options.Topbar)
		matrix, err := topbar.Parse(v, config)
		if err == nil {
			v.Termui.Widgets.Topbar.SetMatrix(matrix)
			v.Termui.Resize()
		} else {
			log.Errorf("topbar configuration: %s", err)
		}

	case options.Database:
		const optionMemory = "memory"
		const optionFilesystem = "filesystem"
		var idx library.Index
		var err error

		value := options.GetString(options.Database)

		if v.index != nil {
			err = v.index.Close()
			if err != nil {
				panic(err)
			}
		}

		switch value {
		case optionMemory:
			idx, err = library.NewInMemory()
		case optionFilesystem:
			idx, err = library.New()
		default:
			log.Errorf("unsupported value '%s', try one of '%s' or '%s'", value, optionMemory, optionFilesystem)
			return
		}

		if err != nil {
			panic(err)
		}

		v.index = idx

	case options.ExpandColumns:
		// Re-render columns
		v.UI().TableWidget().SetColumns(v.UI().TableWidget().ColumnNames())
	}
}
