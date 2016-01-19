package main

import (
	kingpin "gopkg.in/alecthomas/kingpin.v2"
	"log"
	"os"
	"os/signal"
	"strings"
)

var (
	zookeeper        = kingpin.Flag("zookeeper", "zookeeper connection string. can be comma-separated.").Default("localhost:2181").String()
	lockPath         = kingpin.Flag("lockPath", "zookeeper path for lock, should identify task").Required().String()
	cronSchedule     = kingpin.Flag("schedule", "cron expression for task schedule").Required().String()
	command          = kingpin.Flag("command", "command to execute").Required().String()
	workingDirectory = kingpin.Flag("cwd", "change to directory when executing").String()
)

func taskFromArgs() *cronyTask {
	commandAndArgs := *command
	splits := strings.Split(commandAndArgs, " ")
	command := splits[0]
	args := make([]string, 0)
	if len(splits) > 1 {
		args = splits[1:]
	}

	task := &cronyTask{
		cronSchedule: *cronSchedule,
		command:      command,
		args:         args,
		directory:    *workingDirectory,
	}

	return task
}

func main() {
	kingpin.Parse()
	task := taskFromArgs()
	app, err := newApp(strings.Split(*zookeeper, ","), *lockPath, task)
	if err != nil {
		log.Fatal(err)
	}

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)
	go func() {
		<-signals
		app.stop()
	}()
	app.run()
}
