/*
http://www.apache.org/licenses/LICENSE-2.0.txt



Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package slack

import (
	"errors"
	"sync"

	"github.com/intelsdi-x/snap/control/plugin"
	"github.com/intelsdi-x/snap/control/plugin/cpolicy"
	"github.com/intelsdi-x/snap/core"
	"github.com/intelsdi-x/snap/core/ctypes"
	"github.com/nlopes/slack"
)

const (
	// Name of plugin
	name = "slack"
	//vendor namespace part
	vendor = "kindermoumoute"
	// Version of plugin
	version = 1
	// Type of plugin
	pluginType = plugin.CollectorPluginType
)

var (
	errBadMaxElt = errors.New("failed to parse given max_element")
	errBadToken  = errors.New("failed to parse given Slack token")
)

func NewSlackCollector() *Slack {
	return &Slack{}
}

type Slack struct {
	sync.Mutex
	api   *slack.Client
	rtm   *slack.RTM // Slack Real Time Messaging object
	queue metricQueue
	errs  []error
}

// CollectMetrics returns metrics from Slack
func (s *Slack) CollectMetrics(metrics []plugin.MetricType) ([]plugin.MetricType, error) {
	if s.api == nil {
		token, err := getTokenAndSetMaxEltFromConfig(metrics[0].Config().Table())
		if err != nil {
			return nil, err
		}
		s.api = slack.New(token)
		s.rtm = s.api.NewRTM()
		s.queue = newQueue(defaultQueueSize)
		go s.rtm.ManageConnection()
		go s.rtmCollectEvent()
	}
	var mts []plugin.MetricType

	s.Lock()
	for {
		elt, err := s.queue.Pop()
		if err != nil {
			break
		}
		if elt.err != nil {
			return nil, err
		}
		mts = append(mts, elt.Metric)
	}
	s.Unlock()
	return mts, nil
}

// get token and set defaultQueueSize from config
func getTokenAndSetMaxEltFromConfig(table map[string]ctypes.ConfigValue) (string, error) {
	token, ok := table["token"].(ctypes.ConfigValueStr)
	if !ok {
		return "", errBadToken
	}
	maxeltcfg, ok := table["max_element"]
	if !ok {
		return token.Value, nil
	}
	maxelt, ok := maxeltcfg.(ctypes.ConfigValueInt)
	if !ok {
		return "", errBadMaxElt
	}
	defaultQueueSize = maxelt.Value

	return token.Value, nil
}

// GetMetricTypes returns the metric types exposed by gopsutil
func (s *Slack) GetMetricTypes(config plugin.ConfigType) (metricTypes []plugin.MetricType, err error) {
	return []plugin.MetricType{plugin.MetricType{
		Namespace_: core.NewNamespace(vendor, name, "event"),
	}}, nil
}

//GetConfigPolicy returns a ConfigPolicy
func (s *Slack) GetConfigPolicy() (*cpolicy.ConfigPolicy, error) {
	c := cpolicy.New()
	policy := cpolicy.NewPolicyNode()

	token, _ := cpolicy.NewStringRule("token", true)
	maxElt, _ := cpolicy.NewIntegerRule("max_element", false)

	policy.Add(token, maxElt)
	c.Add([]string{vendor, name}, policy)

	return c, nil
}

//Meta returns meta data for testing
func Meta() *plugin.PluginMeta {
	return plugin.NewPluginMeta(
		name,
		version,
		pluginType,
		[]string{plugin.SnapGOBContentType},
		[]string{plugin.SnapGOBContentType},
		plugin.RoutingStrategy(plugin.StickyRouting),
	)
}
