// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// Contributor: Julien Vehent jvehent@mozilla.com [:ulfr]
package main

import (
	"code.google.com/p/gcfg"
	"flag"
	"fmt"
	"mig/event"
	"mig/worker"
	"os"
)

const workerName = "agent_verif_worker"

type Config struct {
	Mq worker.MqConf
}

func main() {
	var (
		err  error
		conf Config
	)
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "%s - a worker verifying agents that fail to authenticate\n", os.Args[0])
		flag.PrintDefaults()
	}
	var configPath = flag.String("c", "/etc/mig/agent_verif_worker.cfg", "Load configuration from file")
	flag.Parse()
	err = gcfg.ReadFileInto(&conf, *configPath)
	if err != nil {
		panic(err)
	}
	// set a binding to route events from event.Q_Agt_Auth_Fail into the queue named after the worker
	// and return a channel that consumes the queue
	workerQueue := "migevent." + workerName
	consumerChan, err := worker.InitMqWithConsumer(conf.Mq, workerQueue, event.Q_Agt_Auth_Fail)
	if err != nil {
		panic(err)
	}
	fmt.Println("started worker", workerName, "consuming queue", workerQueue, "from key", event.Q_Agt_Auth_Fail)
	for event := range consumerChan {
		fmt.Printf("%s\n", event.Body)
	}
	return
}
