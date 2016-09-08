package slack

import (
	"errors"
	"time"

	"github.com/intelsdi-x/snap/control/plugin"
	"github.com/intelsdi-x/snap/core"
)

type metricQueue struct {
	elt                      []MetricElt
	start, currentSize, size int
}

type MetricElt struct {
	Metric plugin.MetricType
	err    error
}

var (
	errEmptyQueue    = errors.New("the metric queue is empty")
	defaultQueueSize = 1000
)

// Create a circular queue for collecting the metrics coming from the stream
func newQueue(size int) metricQueue {
	return metricQueue{elt: make([]MetricElt, size), start: 0, currentSize: 0, size: size}
}

// Pop an element from the queue
func (q *metricQueue) Pop() (MetricElt, error) {
	var elt MetricElt
	if q.currentSize > 0 {
		elt = q.elt[q.start]
		q.start = (q.start + 1) % q.size
		q.currentSize--
		return elt, nil
	}
	return elt, errEmptyQueue
}

// Add an element to the queue
func (q *metricQueue) Push(elt MetricElt) {
	// Overwrite the element if the queue is full
	currentIndex := (q.start + q.currentSize) % q.size
	q.elt[currentIndex] = elt

	if q.currentSize >= q.size {
		q.start = (q.start + 1) % q.size
	} else {
		q.currentSize++
	}
}

// RTM will keep a stream connection between Slack and the plugin. Every time Slack
// send a new event the Slack client will push in in the IncomingEvents channel. This
// function will convert incoming events into Snap metrics, and push it in the metric
// queue.
func (s *Slack) rtmCollectEvent() {
	for event := range s.rtm.IncomingEvents {
		tags, err := GetTags(event.Data, event.Type, "-")
		metric := MetricElt{Metric: plugin.MetricType{
			Namespace_: core.NewNamespace(vendor, name, "event"),
			Data_:      event.Type,
			Timestamp_: time.Now(),
			Tags_:      tags,
		}, err: err}
		s.Lock()
		s.queue.Push(metric)
		s.Unlock()
	}
}
