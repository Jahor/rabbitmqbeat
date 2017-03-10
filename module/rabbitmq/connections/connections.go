package connections

import (
	"fmt"
	"github.com/bernielomax/rabbitmqbeat/module/rabbitmq"
	"github.com/elastic/beats/libbeat/common"
	"github.com/elastic/beats/metricbeat/mb"
	"github.com/pkg/errors"
	"time"
)

const (
	beatType = "rabbitmq.connections"
)

// init registers the MetricSet with the central registry.
// The New method will be called after the setup of the module and before starting to fetch data
func init() {
	if err := mb.Registry.AddMetricSet("rabbitmq", "connections", New); err != nil {
		panic(err)
	}
}

// MetricSet type defines all fields of the MetricSet
// As a minimum it must inherit the mb.BaseMetricSet fields, but can be extended with
// additional entries. These variables can be used to persist data or configuration between
// multiple fetch calls.
type MetricSet struct {
	mb.BaseMetricSet
	api rabbitmq.Api
}

// New create a new instance of the MetricSet
// Part of new is also setting up the configuration by processing additional
// configuration entries if needed.
func New(base mb.BaseMetricSet) (mb.MetricSet, error) {

	client, err := rabbitmq.CreateClient(base)

	if err != nil {
		return nil, err
	}

	return &MetricSet{
		BaseMetricSet: base,
		api:           client,
	}, nil
}

// Fetch methods implements the data gathering and data conversion to the right format
// It returns the event which is then forward to the output. In case of an error, a
// descriptive error must be returned.
func (m *MetricSet) Fetch() ([]common.MapStr, error) {
	connections, err := m.api.Connections()
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("%v api fetch failed.", beatType))
	}
	now := common.Time(time.Now())
	for _, connection := range connections {
		connection["@timestamp"] = now
		connection["type"] = beatType
	}
	return connections, nil
}
