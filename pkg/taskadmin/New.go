package taskadmin

import (
	"errors"
	"log/slog"
	"os"

	"github.com/dracory/taskstore"
)

type LayoutOptions struct {
	Title      string
	Content    string
	Styles     []string
	Scripts    []string
	ScriptURLs []string
	StyleURLs  []string
}

type Options struct {
	Layout     func(options LayoutOptions) string
	LogHandler slog.Handler
	TaskStore  taskstore.StoreInterface
}

func New(options Options) (*admin, error) {
	if options.TaskStore == nil {
		return nil, errors.New("taskadmin > taskstore is required")
	}

	if options.LogHandler == nil {
		options.LogHandler = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelDebug,
		})
	}

	logger := slog.New(options.LogHandler)

	return &admin{
		externalLayout: options.Layout,
		logger:         logger,
		taskStore:      options.TaskStore,
	}, nil
}
