// MIT License
//
// (C) Copyright [2021,2025] Hewlett Packard Enterprise Development LP
//
// Permission is hereby granted, free of charge, to any person obtaining a
// copy of this software and associated documentation files (the "Software"),
// to deal in the Software without restriction, including without limitation
// the rights to use, copy, modify, merge, publish, distribute, sublicense,
// and/or sell copies of the Software, and to permit persons to whom the
// Software is furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included
// in all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL
// THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR
// OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE,
// ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR
// OTHER DEALINGS IN THE SOFTWARE.

package main

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	tkafka "github.com/Cray-HPE/hms-trs-kafkalib/v2/pkg/trs-kafkalib"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/sirupsen/logrus"
)

const (
	SYNC_TOKEN = "STARTSYNC"
)


func main() {

	var envstr string

	brokerSpec := "kafka:9092"
	ctx, _ := context.WithCancel(context.Background())
	kinst := &tkafka.TRSKafka{}
	rspChan := make(chan *kafka.Message)
	useQuit := false

	//SETUP LOGGER
	//It is important to do the NEW else, the pointer is nil, and when you TRY to set the level, that will fail.  If you DONT do NEW, just dont try to access the logger!
	logy := logrus.New()
	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})

	//Check if we're a sender or echo server

	ksender := true
	envstr = os.Getenv("KTYPE")
	if envstr != "" {
		if strings.ToLower(envstr) == "sender" {
			ksender = true
		} else {
			ksender = false
		}
	}
	envstr = os.Getenv("USE_QUIT")
	if envstr != "" {
		useQuit = true
	}

	envstr = os.Getenv("LOG_LEVEL")
	if envstr != "" {
		logLevel := strings.ToUpper(envstr)
		logrus.Infof("Setting log level to: %s", envstr)

		switch logLevel {

		case "TRACE":
			logrus.SetLevel(logrus.TraceLevel)
		case "DEBUG":
			logrus.SetLevel(logrus.DebugLevel)
		case "INFO":
			logrus.SetLevel(logrus.InfoLevel)
		case "WARN":
			logrus.SetLevel(logrus.WarnLevel)
		case "ERROR":
			logrus.SetLevel(logrus.ErrorLevel)
		case "FATAL":
			logrus.SetLevel(logrus.FatalLevel)
		case "PANIC":
			logrus.SetLevel(logrus.PanicLevel)
		default:
			logrus.SetLevel(logrus.ErrorLevel)
		}

		//Set the kafka level to the same level.
		logy.SetLevel(logrus.GetLevel())
	}

	//Note: this doesn't seem to be necessary but does make the
	//runs cleaner.
	logrus.Infof("Giving kafka time to warm up...")
	time.Sleep(10 * time.Second)

	envstr = os.Getenv("BROKER")
	if envstr != "" {
		brokerSpec = envstr
	}

	if ksender {
		logrus.Infof("** FUNCTION: SENDER **")

		st := "test"
		rt := []string{"test"}
		cg := ""
		err := kinst.Init(ctx, rt, cg, brokerSpec, rspChan, logy)

		if err != nil {
			logrus.Info("Init() ERROR:", err)
			os.Exit(1)
		}

		// Echo back what the receiver has ricocheted back to us.
		//
		// We'll send a special token to the receiver over and over
		// which the receiver will eventually echo back.  Once we
		// receive it, set a flag so we know kafka is operational.

		syncup := false

		go func() {
			for {
				pld := <-rspChan
				if (strings.Contains(string(pld.Value),SYNC_TOKEN)) {
					logrus.Infof("Start/Syncp message received.")
					syncup = true
					continue
				}
				logrus.Infof("RECEIVED: '%s'", string(pld.Value))
			}
		}()

		//Send messages

		time.Sleep(2 * time.Second)
		logrus.Infof("Send topics: '%s'", kinst.RcvTopicNames)

		numMsgs := 5
		envstr := os.Getenv("NUMMSGS")
		if envstr != "" {
			numMsgs, _ = strconv.Atoi(envstr)
		}

		logrus.Infof("Waiting for channels and components to be ready...")

		for ii := 0; ii < 60 ; ii ++ {
			//Send a startup token which will get echo'd back to verify
			//everyone is ready to play.
			kinst.Write(st, []byte(SYNC_TOKEN))
			if (syncup == true) {
				break
			}
			time.Sleep(time.Second)
		}

		if (!syncup) {
			logrus.Errorf("ERROR: never got start/sync message, exiting.")
			os.Exit(1)
		}

		logrus.Infof("Sending %d messages...", numMsgs)

		for ii := 0; ii < numMsgs; ii++ {
			pld := fmt.Sprintf("Test message iteration %d", ii)
			logrus.Infof("Sending message # %d", ii)
			kinst.Write(st, []byte(pld))
			time.Sleep(10 * time.Millisecond)
			if (ii == (numMsgs / 3)) {
				//Minor test of topic changing on the fly
				err := kinst.SetTopics([]string{"test","newtest"})
				if (err != nil) {
					logrus.Errorf("Error setting new topics: %v",err)
				}
			}
		}
		logrus.Infof("Sending quit message...")
		kinst.Write(st, []byte("quit"))
		time.Sleep(1 * time.Second)
		logrus.Infof("Closing channel...")
		kinst.Shutdown()
		time.Sleep(10 * time.Second)

		logrus.Infof("All done...")
		os.Exit(0)
	} else {
		logrus.Infof("** FUNCTION: ECHO **")

		st := "test"
		rt := []string{"test"}
		cg := ""
		err := kinst.Init(ctx, rt, cg, brokerSpec, rspChan, logy)
		if err != nil {
			logrus.Info("Init() ERROR:", err)
			os.Exit(1)
		}

		kinst.Client.Consumer.Logger.SetLevel(logrus.ErrorLevel)
		logrus.Infof("Rcv topics: '%s'", kinst.RcvTopicNames)
		time.Sleep(2 * time.Second)
		msgIX := 0

		for {
			pld := <-kinst.Client.Consumer.Responses
			msgIX ++
			if (msgIX == 10) {
				//Minor test of topic changing on the fly
				err := kinst.SetTopics([]string{"test","newtest"})
				if (err != nil) {
					logrus.Errorf("Error setting new topics: %v",err)
				}
			}

			if string(pld.Value) == "quit" {
				if useQuit {
					logrus.Infof("Received: '%s' -- exiting.",
						string(pld.Value))
					break
				} else {
					logrus.Infof("Received: '%s' -- continuing.",
						string(pld.Value))
				}
			} else {
				//Only echo back once.

				if (!strings.Contains(string(pld.Value),"Echo:") &&
					!strings.Contains(string(pld.Value),SYNC_TOKEN)) {
					logrus.Infof("Received: '%s' -- echoing back.",
							string(pld.Value))
					eval := "Echo: " + string(pld.Value)
					kinst.Write(st, []byte(eval))
				} else {
					logrus.Infof("Received: '%s' .", string(pld.Value))
				}
			}
		}

		logrus.Infof("Closing channel...")
		time.Sleep(2 * time.Second)
		kinst.Shutdown()
		os.Exit(0)
	}
}
