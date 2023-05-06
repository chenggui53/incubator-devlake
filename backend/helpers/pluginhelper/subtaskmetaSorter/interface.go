package subtaskmetaSorter

import "github.com/apache/incubator-devlake/core/plugin"

type SubTaskMetaSorter interface {
	Sort() ([]plugin.SubTaskMeta, error)
	DetectLoop() error
}
