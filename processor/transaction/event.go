package transaction

import (
	"time"

	m "github.com/elastic/apm-server/model"
	"github.com/elastic/apm-server/utility"
	"github.com/elastic/beats/libbeat/common"
)

type Event struct {
	Id        string        `json:"id"`
	Name      string        `json:"name"`
	Type      string        `json:"type"`
	Result    *string       `json:"result"`
	Duration  float64       `json:"duration"`
	Timestamp time.Time     `json:"timestamp"`
	Context   common.MapStr `json:"context"`
	Traces    []Trace       `json:"traces"`
}

func (t *Event) DocType() string {
	return "transaction"
}

func (t *Event) Transform() common.MapStr {
	enh := utility.NewMapStrEnhancer()
	tx := common.MapStr{"id": t.Id}
	enh.Add(tx, "name", t.Name)
	enh.Add(tx, "duration", utility.MillisAsMicros(t.Duration))
	enh.Add(tx, "type", t.Type)
	enh.Add(tx, "result", t.Result)
	return tx
}

func (t *Event) Mappings(pa *payload) (time.Time, []m.DocMapping) {
	return t.Timestamp,
		[]m.DocMapping{
			{Key: "processor", Apply: func() common.MapStr {
				return common.MapStr{"name": processorName, "event": t.DocType()}
			}},
			{Key: t.DocType(), Apply: t.Transform},
			{Key: "context", Apply: func() common.MapStr { return t.Context }},
			{Key: "context.app", Apply: pa.App.Transform},
			{Key: "context.system", Apply: pa.System.Transform},
		}
}
