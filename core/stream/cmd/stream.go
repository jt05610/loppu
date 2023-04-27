/*
Copyright © 2023 Jonathan Taylor <jonrtaylor12@gmail.com>
*/

package cmd

import (
	"fmt"
	"github.com/jt05610/loppu/core/stream/redis"
	"github.com/spf13/cobra"
	"log"
	"os"
	"path"
	"time"
)

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "start the stream",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		s := redis.Node{}
		dest := path.Join("nodes", "stream", "node.yaml")
		df, err := os.OpenFile(dest, os.O_RDONLY, os.ModePerm)
		if err != nil {
			log.Fatal(err)
		}
		err = s.Load(df)
		if err != nil {
			log.Fatal(err)
		}
		err = s.Start()
		if err != nil {
			log.Fatal(err)
		}
	},
}

// stopCmd represents the stop command
var stopCmd = &cobra.Command{
	Use:   "stop",
	Short: "stop the stream",
	Long:  `stops stream`,
	Run: func(cmd *cobra.Command, args []string) {
		s := redis.Node{}
		err := s.Stop()
		if err != nil {
			log.Fatal(err)
		}
	},
}

// initCmd represents the stop command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "initializes stream",
	Long:  `use to initialize`,
	Run: func(cmd *cobra.Command, args []string) {
		s := redis.NewRedisNode()
		dest := path.Join("nodes", "stream", "node.yaml")
		df, err := os.Create(dest)
		defer df.Close()
		if err != nil {
			panic(err)
		}
		err = s.Flush(df)
		if err != nil {
			panic(err)
		}
	},
}

var (
	sName string
	sId   string
	delay time.Duration
	reqs  []string
	// addCmd represents the stop command
	addCmd = &cobra.Command{
		Use:   "add",
		Short: "add stream",
		Long:  `adds stream to node.yaml`,
		Run: func(cmd *cobra.Command, args []string) {
			s := redis.NewRedisNode()
			dest := path.Join("nodes", "stream", "node.yaml")
			df, err := os.OpenFile(dest, os.O_RDONLY, os.ModePerm)
			if err != nil {
				panic(err)
			}
			err = s.Load(df)
			if err != nil {
				panic(err)
			}
			_ = df.Close()
			df, err = os.OpenFile(dest, os.O_WRONLY, os.ModePerm)
			if err != nil {
				panic(err)
			}
			reqArray := make([]*redis.Request, 0)
			for i, r := range reqs {
				reqArray = append(reqArray, &redis.Request{
					Name: fmt.Sprintf("req_%v", i),
					Uri:  r,
				})
			}
			newStream := redis.NewStream(sName, sId, delay, reqArray)
			s.(*redis.Node).AddStream(newStream)
			err = s.Flush(df)
			if err != nil {
				panic(err)
			}
			fmt.Println("stream added")
		},
	}
)

func init() {
	rootCmd.AddCommand(initCmd)
	rootCmd.AddCommand(startCmd)
	rootCmd.AddCommand(stopCmd)
	addCmd.PersistentFlags().StringVarP(&sName, "name", "n", "newStream",
		"name of stream")
	addCmd.PersistentFlags().StringVarP(&sId, "id", "i", "steamID",
		"sample id for stream (will probably deprecate)")
	addCmd.PersistentFlags().DurationVarP(&delay, "delay", "d",
		time.Duration(100)*time.Millisecond,
		"Amount of time to delay between requests")
	addCmd.PersistentFlags().StringArrayVarP(&reqs, "reqs", "r",
		[]string{"http://localhost:50000/", "http://localhost:50001/"},
		"node urls to sent Get requests to")
	rootCmd.AddCommand(addCmd)
}
