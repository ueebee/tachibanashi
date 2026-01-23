package master

import "github.com/ueebee/tachibanashi/model"

type EventDownloadComplete struct {
	Fields model.Attributes
}

func (e *EventDownloadComplete) UnmarshalJSON(data []byte) error {
	return unmarshalAttributes(data, &e.Fields)
}
