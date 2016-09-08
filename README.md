[![Build Status](https://travis-ci.com/kindermoumoute/snap-plugin-collector-slack.svg?branch=master)](https://travis-ci.com/kindermoumoute/snap-plugin-collector-slack)

# snap collector plugin - Slack
This plugin collects Slack events from the Slack API.


It's used in the [snap framework](http://github.com:intelsdi-x/snap).

1. [Getting Started](#getting-started)
  * [System Requirements](#system-requirements)
  * [Installation](#installation)
  * [Configuration and Usage](configuration-and-usage)
2. [Documentation](#documentation)
  * [Collected Metrics](#collected-metrics)
  * [Examples](#examples)
  * [Roadmap](#roadmap)
3. [Community Support](#community-support)
4. [Contributing](#contributing)
5. [License](#license-and-authors)
6. [Acknowledgements](#acknowledgements)

## Getting Started
### System Requirements
* [golang 1.5+](https://golang.org/dl/) (needed only for building)

### Operating systems
All OSs currently supported by snap:
* Linux/amd64
* Darwin/amd64

### Installation
#### Download the Slack plugin binary:
You can get the pre-built binaries for your OS and architecture at snap's [GitHub Releases](https://github.com/intelsdi-x/snap/releases) page.

#### To build the plugin binary:
Download the repository with Go:

```
$ go get -d github.com/kindermoumoute/snap-plugin-collector-slack
```

Build the plugin by running make within the cloned repo:
```
$ make
```
This builds the plugin in `build/rootfs/`

### Configuration and Usage
* Set up the [snap framework](https://github.com/intelsdi-x/snap/blob/master/README.md#getting-started)
* Ensure `$SNAP_PATH` is exported  
`export SNAP_PATH=$GOPATH/src/github.com/intelsdi-x/snap/build`

## Documentation
There are a number of other resources you can review to learn to use this plugin:

* [Slack API](https://api.slack.com)
* [Slack Event documentation](https://api.slack.com/events-api#event_types)
* [snap Slack unit test](https://github.com/kindermoumoute/snap-plugin-collector-slack/blob/master/slack/slack_small_test.go)
* [snap Slack examples](#examples)


### Collected Metrics
This plugin has the ability to gather a Slack event as a metric.
* Namespace: `/kindermoumoute/slack/event`
* Data: data is the event type (see [Slack event documentation](https://api.slack.com/events-api#event_types)).
* Tags: generate a tag for each existing data returned by the event.

### Examples
Before running a task with this plugin, you need to generate a token from the Slack API. This token can be associated to a *bot* or to a *bot user*.

##### Creating a bot user 
You can create a bot user [here](https://api.slack.com/docs/oauth-test-tokens). With a *bot user* token, you will collect all event visible from your *Slack user*.

##### Creating a bot
You can create a bot at [https://your-org.slack.com/services/new/bot](https://your-org.slack.com/services/new/bot). With a *bot* token you will collect all event visible from the *bot*. As a *bot* does not have access to a channel by default, you have to invite this *bot* in the channel you want to collect the event from.

#### Create a task
In one terminal window, open the snap daemon:
```bash
$ $SNAP_PATH/bin/snapd -t 0
```

In another terminal window:
Load Slack plugin
```bash
$ $SNAP_PATH/bin/snapctl plugin load <path to built slack plugin>
```
See available metrics for your system
```bash
$ $SNAP_PATH/bin/snapctl metric list
```

Create a task manifest file (e.g. `task.json`):
```json
    "schedule": {
        "type": "simple",
        "interval": "1s"
    },
    "max-failures": 10,
    "workflow": {
        "collect": {
            "metrics": {
                "/kindermoumoute/slack/event": {}
            },
            "config": {
                "/kindermoumoute/slack": {
                    "token": "INSERT_YOUR_TOKEN_HERE"
                }
            },
             "publish": [
                        {
                            "plugin_name": "file",                            
                            "config": {
                                "file": "/tmp/published_slack.json"
                            }
                        }
                    ] 
      }
   }
}
```

Load file plugin for publishing:
```bash
$ $SNAP_PATH/bin/snapctl plugin load build/plugin/snap-publisher-file
Plugin loaded
Name: file
Version: 3
Type: publisher
Signed: false
Loaded Time: Fri, 20 Nov 2015 11:41:39 PST
```

Create task:
```bash
$ $SNAP_PATH/bin/snapctl task create -t examples/task.json
Using task manifest to create task
Task created
ID: 02dd7ff4-8106-47e9-8b86-70067dc0a850
Name: Task-02dd7ff4-8106-47e9-8b86-70067dc0a850
State: Running
```

See file output (this is just part of the file):
```json
null
null
[ {"timestamp":"2016-08-24T18:55:31.926246288-07:00","namespace":"/kindermoumoute/slack/event","data":"presence_change","unit":"","tags":{"plugin_running_on":"ocano-mac01.jf.intel.com","presence":"active","type":"presence_change","user":"U063H8QUB"},"version":0,"last_advertised_time":"0001-01-01T00:00:00Z"},
  {"timestamp":"2016-08-24T18:55:32.579676428-07:00","namespace":"/kindermoumoute/slack/event","data":"user_typing","unit":"","tags":{"channel":"C2136PUKT","plugin_running_on":"ocano-mac01.jf.intel.com","type":"user_typing","user":"U12LSMG85"},"version":0,"last_advertised_time":"0001-01-01T00:00:00Z"}]
null
[ {"timestamp":"2016-08-24T18:55:34.509201456-07:00","namespace":"/kindermoumoute/slack/event","data":"message","unit":"","tags":{"msg-channel":"C2136PUKT","msg-isstarred":"false","msg-text":"plop","msg-timestamp":"1472090134.000049","msg-type":"message","msg-user":"U12LSMG85","plugin_running_on":"ocano-mac01.jf.intel.com"},"version":0,"last_advertised_time":"0001-01-01T00:00:00Z"}]
```

Stop task:
```bash
$ $SNAP_PATH/bin/snapctl task stop 02dd7ff4-8106-47e9-8b86-70067dc0a850
Task stopped:
ID: 02dd7ff4-8106-47e9-8b86-70067dc0a850
```

### Roadmap
There isn't a current roadmap for this plugin, but it is in active development. As I launch this plugin, I do not have any outstanding requirements for the next release. If you have a feature request, please add it as an [issue](https://github.com/kindermoumoute/snap-plugin-collector-slack/issues/new) and/or submit a [pull request](https://github.com/kindermoumoute/snap-plugin-collector-slack/pulls).

## Community Support
This repository is one of **many** plugins in **snap**, a powerful telemetry framework. See the full project at http://github.com/intelsdi-x/snap To reach out to other users, head to the [main framework](https://github.com/intelsdi-x/snap#community-support)

## Contributing
I love contributions!

There's more than one way to give back, from examples to blogs to code updates. 

## License
[snap](http://github.com:intelsdi-x/snap), along with this plugin, is an Open Source software released under the Apache 2.0 [License](LICENSE).

## Acknowledgements
* Author: [@kindermoumoute](https://github.com/kindermoumoute/)

And **thank you!** Your contribution, through code and participation, is incredibly important to us.