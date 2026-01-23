package master

import (
	"encoding/json"

	"github.com/ueebee/tachibanashi/model"
)

const (
	DateInfoFieldDayKey               = "sDayKey"
	DateInfoFieldMaeEigyouDay1        = "sMaeEigyouDay_1"
	DateInfoFieldMaeEigyouDay2        = "sMaeEigyouDay_2"
	DateInfoFieldMaeEigyouDay3        = "sMaeEigyouDay_3"
	DateInfoFieldTheDay               = "sTheDay"
	DateInfoFieldYokuEigyouDay1       = "sYokuEigyouDay_1"
	DateInfoFieldYokuEigyouDay2       = "sYokuEigyouDay_2"
	DateInfoFieldYokuEigyouDay3       = "sYokuEigyouDay_3"
	DateInfoFieldYokuEigyouDay4       = "sYokuEigyouDay_4"
	DateInfoFieldYokuEigyouDay5       = "sYokuEigyouDay_5"
	DateInfoFieldYokuEigyouDay6       = "sYokuEigyouDay_6"
	DateInfoFieldYokuEigyouDay7       = "sYokuEigyouDay_7"
	DateInfoFieldYokuEigyouDay8       = "sYokuEigyouDay_8"
	DateInfoFieldYokuEigyouDay9       = "sYokuEigyouDay_9"
	DateInfoFieldYokuEigyouDay10      = "sYokuEigyouDay_10"
	DateInfoFieldKabuUkewatasiDay     = "sKabuUkewatasiDay"
	DateInfoFieldKabuKariUkewatasiDay = "sKabuKariUkewatasiDay"
	DateInfoFieldBondUkewatasiDay     = "sBondUkewatasiDay"
)

type DateInfo struct {
	Fields model.Attributes
}

func (d *DateInfo) UnmarshalJSON(data []byte) error {
	var raw map[string]json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	d.Fields = make(model.Attributes, len(raw))
	for key, value := range raw {
		d.Fields[key] = jsonString(value)
	}
	return nil
}

func (d DateInfo) DayKey() string {
	return d.Fields.Value(DateInfoFieldDayKey)
}

func (d DateInfo) TheDay() string {
	return d.Fields.Value(DateInfoFieldTheDay)
}
