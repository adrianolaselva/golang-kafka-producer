package cmd

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"gopkg.in/confluentinc/confluent-kafka-go.v1/kafka"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"
)

var producerCommand = &cobra.Command{
	Use:   "produce",
	Short: "Producer events",
	Args: func(cmd *cobra.Command, args []string) error {

		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		files, err := cmd.Flags().GetStringArray("files")
		if err != nil {
			fmt.Printf("files not defined: %v\n", err.Error())
			os.Exit(1)
		}

		iterations, _ := cmd.Flags().GetInt("iterations")

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

		fmt.Printf("iterations: %v\n", iterations)

		for _, config := range configurations {
			for _, bridge := range config.Bridges {
				fmt.Printf("%v\n", config)
				p, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": bridge.Parameters.BootstrapServers})
				if err != nil {
					panic(err)
				}

				defer p.Close()

				for i := 1; i <= iterations; i++ {
					value, _ := json.Marshal(map[string]interface{}{
						"userId":      9999999999,
						"anonymousId": fmt.Sprintf("%x", md5.Sum([]byte(string(i*1000)))),
						"keyType":     "CONSUMER",
						"integrations": map[string]interface{}{
							"all":       false,
							"mixpanel":  false,
							"appsflyer": false,
							"s3":        false,
						},
						"traits": map[string]interface{}{
							"appsflyer_id":   fmt.Sprintf("%x", md5.Sum([]byte(string(i)))),
							"advertising_id": fmt.Sprintf("%x", md5.Sum([]byte(string(time.Now().Format(time.RFC3339))))),
							"os":             "ios",
							"document":       "333.111.222-32",
							"cnpj":           "33.111.222/0001-32",
							"email":          "adrianolaselva@gmail.com",
							"cvv":            "335",
							"credit_card":    "5232 9609 5260 4191",
							"card":           "5232960952604191",
						},
						"context": map[string]interface{}{
							"ip":             "127.0.0.1",
							"appsflyer_id":   fmt.Sprintf("%x", md5.Sum([]byte(string(i)))),
							"advertising_id": fmt.Sprintf("%x", md5.Sum([]byte(string(time.Now().Format(time.RFC3339))))),
							"idfa":           fmt.Sprintf("%x", md5.Sum([]byte(string(time.Now().Format(time.RFC3339))))),
						},
						"timestamp": time.Now().Format(time.RFC3339),
					})

					err = p.Produce(&kafka.Message{
						TopicPartition: kafka.TopicPartition{Topic: &bridge.Topic, Partition: kafka.PartitionAny},
						Key:            []byte(fmt.Sprintf("%x", md5.Sum(value))),
						Value:          value,
					}, nil)

					if err != nil {
						fmt.Printf("failed to produce events: %v\n", err.Error())
						os.Exit(1)
					}

					fmt.Printf("produce event: %v\n", string(value))
				}

				p.Flush(-1)
			}
		}
	},
}

type FileInfo struct {
	os.FileInfo
	Path string
}

func loadYmlFile(filePaths []string) ([]FileInfo, error) {
	var files []FileInfo
	for _, file := range filePaths {
		for _, extension := range []string{".yml", ".yaml"} {
			if strings.LastIndex(file, extension) > 0 {
				fileInfo, err := os.Stat(file)
				if err != nil {
					continue
				}

				files = append(files, FileInfo{
					FileInfo: fileInfo,
					Path:     path.Dir(file),
				})
			}
		}
	}

	return files, nil
}


type Configuration struct {
	Bridges      []struct{
		Topic      string `yaml:"topic"`
		Parameters struct {
			GroupId          string `yaml:"group.id"`
			BootstrapServers string `yaml:"bootstrap.servers"`
			AutoOffsetReset  string `yaml:"auto.offset.reset"`
		} `yaml:"parameters"`
	} `yaml:"bridges"`
}

