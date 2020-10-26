package main

import (
	"context"
	"fmt"
	"github.com/Shopify/sarama"
	"github.com/sirupsen/logrus"
	"os"
	tkafka "stash.us.cray.com/HMS/hms-trs-kafkalib/pkg/trs-kafkalib"
	"strconv"
	"strings"
	"time"
)

func main() {

	var envstr string

	brokerSpec := "kafka:9092"
	ctx, _ := context.WithCancel(context.Background())
	kinst := &tkafka.TRSKafka{}
	rspChan := make(chan *sarama.ConsumerMessage)
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
		logrus.Infof("Setting log level to: %d\n", envstr)

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

	envstr = os.Getenv("BROKER")
	if envstr != "" {
		brokerSpec = envstr
	}

	if ksender {
		logrus.Infof("** FUNCTION: SENDER **\n")

		//st,rt,cg := tkafka.GenerateSendReceiveConsumerGroupName("KTestApp","REST","")
		st := "test"
		rt := []string{"test"}
		cg := ""
		err := kinst.Init(ctx, rt, cg, brokerSpec, rspChan, logy)

		if err != nil {
			logrus.Info("Init() ERROR:", err)
			os.Exit(1)
		}

		time.Sleep(2 * time.Second)

		go func() {
			for {
				pld := <-rspChan
				logrus.Infof("RECEIVED: '%s'\n", string(pld.Value))
			}
		}()

		//Send messages

		time.Sleep(1 * time.Second)
		logrus.Infof("Rcv topics: '%s'\n", kinst.RcvTopicNames)

		numMsgs := 5
		envstr := os.Getenv("NUMMSGS")
		if envstr != "" {
			numMsgs, _ = strconv.Atoi(envstr)
		}
		logrus.Infof("Sending %d messages...\n", numMsgs)

		for ii := 0; ii < numMsgs; ii++ {
			pld := fmt.Sprintf("Test message iteration %d", ii)
			logrus.Infof("Sending message # %d\n", ii)
			kinst.Write(st, []byte(pld))
			time.Sleep(100 * time.Millisecond)
		}
		logrus.Infof("Sending quit message...\n")
		kinst.Write(st, []byte("quit"))
		time.Sleep(1 * time.Second)
		logrus.Infof("Closing channel...\n")
		kinst.Shutdown()
		time.Sleep(2 * time.Second)

		logrus.Infof("All done...\n")
		os.Exit(0)
	} else {
		logrus.Infof("** FUNCTION: ECHO **\n\n")

		//st,rt,cg := tkafka.GenerateSendReceiveConsumerGroupName("KTestApp","REST","")
		st := "test"
		rt := []string{"test"}
		cg := ""
		err := kinst.Init(ctx, rt, cg, brokerSpec, rspChan, logy)
		if err != nil {
			logrus.Info("Init() ERROR:", err)
			os.Exit(1)
		}

		kinst.Client.Consumer.Logger.SetLevel(logrus.ErrorLevel)
		logrus.Infof("Rcv topics: '%s'\n", kinst.RcvTopicNames)
		for {
			pld := <-kinst.Client.Consumer.Responses
			if string(pld.Value) == "quit" {
				if useQuit {
					logrus.Infof("Received: '%s' -- exiting.\n", string(pld.Value))
					break
				} else {
					logrus.Infof("Received: '%s' -- continuing.\n", string(pld.Value))
				}
			} else {
				logrus.Infof("Received: '%s' -- echoing back.\n", string(pld.Value))
				kinst.Write(st, []byte(pld.Value))
			}
		}

		logrus.Infof("Closing channel...\n")
		kinst.Shutdown()
		time.Sleep(2 * time.Second)
		os.Exit(0)
	}
}
