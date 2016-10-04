/*
Package include imports all Module and MetricSet packages so that they register
their factories with the global registry. This package can be imported in the
main package to automatically register all of the standard supported Metricbeat
modules.
*/
package include

import (
	// This list is automatically generated by `make imports`
	_ "github.com/bernielomax/rabbitmqbeat/module/rabbitmq"
	_ "github.com/bernielomax/rabbitmqbeat/module/rabbitmq/nodes"
	_ "github.com/bernielomax/rabbitmqbeat/module/rabbitmq/overview"
	_ "github.com/bernielomax/rabbitmqbeat/module/rabbitmq/queues"
)