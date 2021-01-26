package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"gopkg.in/confluentinc/confluent-kafka-go.v1/kafka"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
)

var consumerCommand = &cobra.Command{
	Use:   "consumer",
	Short: "Consumer events",
	Args: func(cmd *cobra.Command, args []string) error {

		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		sigchan := make(chan os.Signal, 1)
		signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

		files, err := cmd.Flags().GetStringArray("files")
		if err != nil {
			fmt.Printf("files not defined: %v\n", err.Error())
			os.Exit(1)
		}

		filesInfo, err := loadYmlFile(files)
		if err != nil {
			fmt.Printf("failed to load files info: %v\n", err.Error())
			os.Exit(1)
		}

		configurations := make([]Configuration, 0)
		for _, fileInfo := range filesInfo {
			fmt.Printf("fount file: %v/%v\n", fileInfo.Path, fileInfo.Name())
			if strings.LastIndex(fileInfo.Name(), ".dist") > 0 {
				continue
			}

			data, err := ioutil.ReadFile(filepath.Join(fileInfo.Path, fileInfo.Name()))
			if err != nil {
				fmt.Printf("failed to load file: %v\n", err.Error())
				os.Exit(1)
			}

			var config Configuration
			if err = yaml.Unmarshal(data, &config); err != nil {
				fmt.Printf("failed to load file: %v\n", err.Error())
				os.Exit(1)
			}

			configurations = append(configurations, config)
		}

		if len(configurations) == 0 {
			fmt.Printf("configuration not defined")
			os.Exit(1)
		}

		for _, config := range configurations {
			for _, bridge := range config.Bridges {
				fmt.Printf("%v\n", config)
				consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
					"group.id": bridge.Parameters.GroupId,
					"bootstrap.servers": bridge.Parameters.BootstrapServers,
				})
				if err != nil {
					fmt.Printf("failed to create consumer: %s\n", err.Error())
					os.Exit(1)
				}

				defer consumer.Close()

				topics := append(strings.Split(bridge.Topic, ","), "^aRegex.*[Tt]opic")

				err = consumer.SubscribeTopics(topics, nil)
				if err != nil {
					fmt.Printf("failed to subscribe: %s\n", err.Error())
					os.Exit(1)
				}

				fmt.Printf("topic successfully subscribed [%v], waiting messages\n", topics)

				run := true
				for run {
					select {
					case sig := <-sigchan:
						run = false
						fmt.Printf("caught signal %v: terminating\n", sig)
						os.Exit(1)
					default:
						ev := consumer.Poll(1000)
						switch e := ev.(type) {
						case *kafka.Message:
							fmt.Printf("message[%v] => %v\n", string(e.Key), string(e.Value))
							_, err = consumer.Commit()
						case kafka.AssignedPartitions:
							fmt.Printf("assigned partitions [%v]", e.Partitions)
							_ = consumer.Assign(e.Partitions)
						case kafka.RevokedPartitions:
							fmt.Printf("revoked partitions [%v]", e.Partitions)
							_ = consumer.Unassign()
						case kafka.PartitionEOF:
							fmt.Printf("partition EOF")
						case kafka.Error:
							fmt.Printf("kafka error: %s", e.Error())
						}
					}
				}
			}
		}
	},
}
