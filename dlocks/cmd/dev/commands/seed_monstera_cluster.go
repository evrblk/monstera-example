package commands

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/evrblk/monstera"
	"github.com/samber/lo/mutable"
	"github.com/spf13/cobra"
)

var seedMonsteraClusterCmd = &cobra.Command{
	Use:   "seed-monstera-cluster",
	Short: "Seed Monstera cluster config file",
	Run: func(cmd *cobra.Command, args []string) {
		if err := os.MkdirAll("./", os.ModePerm); err != nil {
			log.Fatal(err)
		}

		clusterConfig := monstera.CreateEmptyConfig()

		nodeIds := createNodes(clusterConfig, 3)

		for _, a := range applications {
			createApplication(clusterConfig, a, nodeIds)
		}

		data, err := monstera.WriteConfigToProto(clusterConfig)
		if err != nil {
			log.Fatal(err)
		}

		if err := os.WriteFile(filepath.Join("./", "cluster_config.pb"), data, 0666); err != nil {
			log.Fatal(err)
		}

		data, err = monstera.WriteConfigToJson(clusterConfig)
		if err != nil {
			log.Fatal(err)
		}

		if err := os.WriteFile(filepath.Join("./", "cluster_config.json"), data, 0666); err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(seedMonsteraClusterCmd)
}

type Application struct {
	Name              string
	Implementation    string
	ReplicationFactor int
	ShardsCount       int
}

var (
	applications = []Application{
		{
			Name:              "Accounts",
			Implementation:    "Accounts",
			ReplicationFactor: 3,
			ShardsCount:       1,
		},
		{
			Name:              "Locks",
			Implementation:    "Locks",
			ReplicationFactor: 3,
			ShardsCount:       16,
		},
		{
			Name:              "Namespaces",
			Implementation:    "Namespaces",
			ReplicationFactor: 3,
			ShardsCount:       16,
		},
	}
)

func createApplication(clusterConfig *monstera.ClusterConfig, application Application, nodeIds []string) {
	_, err := clusterConfig.CreateApplication(application.Name, application.Implementation, int32(application.ReplicationFactor))
	if err != nil {
		log.Fatal(err)
	}

	step := 256 / application.ShardsCount
	for i := 0; i < application.ShardsCount; i++ {
		shard, err := clusterConfig.CreateShard(application.Name, []byte{byte(step * i), 0x00, 0x00, 0x00}, []byte{byte(step*(i+1)) - 1, 0xff, 0xff, 0xff}, "")
		if err != nil {
			log.Fatal(err)
		}

		shuffledIds := make([]string, len(nodeIds))
		copy(shuffledIds, nodeIds)
		mutable.Shuffle(shuffledIds)

		for j := 0; j < application.ReplicationFactor; j++ {
			_, err := clusterConfig.CreateReplica(application.Name, shard.Id, shuffledIds[j])
			if err != nil {
				log.Fatal(err)
			}
		}
	}
}

func createNodes(clusterConfig *monstera.ClusterConfig, count int) []string {
	result := make([]string, count)
	for i := 0; i < count; i++ {
		node, err := clusterConfig.CreateNode(fmt.Sprintf("localhost:%d", i+7000))
		if err != nil {
			log.Fatal(err)
		}
		result[i] = node.Id
	}
	return result
}
