// Copyright © 2018 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (

	"github.com/spf13/cobra"
	"github.com/yhaobj/docker/agent"
)

// agentCmd represents the agent command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "cdn ats agent",
	Long: `cdn ats agent for ats config restart, pull docker images`,
	Run: func(cmd *cobra.Command, args []string) {
		agent.Run(cmd, args)
	},
}

func init() {
	rootCmd.AddCommand(runCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// agentCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// agentCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	runCmd.Flags().StringP("address", "a", "127.0.0.1:8080", "listen ip and port")
	//agentCmd.Flags().StringP("port", "p", "8080", "listen port")
	runCmd.Flags().StringP("endpoint", "e", "127.0.0.1:80", "remote server")
	runCmd.Flags().StringP("user", "u", "admin", "registry user")
	runCmd.Flags().StringP("password", "p", "admin", "registry password")
	runCmd.Flags().StringP("certfile", "c", "./ca.pem", "cert file")
	runCmd.Flags().StringP("keyfile", "k", "./key.pem", "key file")
}
