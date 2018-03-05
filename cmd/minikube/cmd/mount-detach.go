/*
Copyright 2016 The Kubernetes Authors All rights reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package cmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"github.com/spf13/viper"
	"github.com/golang/glog"
	"github.com/spf13/cobra"
	cmdUtil "k8s.io/minikube/cmd/util"
	cfg "k8s.io/minikube/pkg/minikube/config"
	"k8s.io/minikube/pkg/minikube/constants"
	"k8s.io/minikube/pkg/minikube/machine"
)

// mountDetachCmd represents the mount command
var mountDetachCmd = &cobra.Command{
	Use:   "mount-detach [flags] MOUNT_DIRECTORY(ex:\"/home\")",
	Short: "Mounts the specified directory into minikube",
	Long:  `Mounts the specified directory into minikube.`,
	Run: func(cmd *cobra.Command, args []string) {

		mountString := args[0]
		var debugVal int
		if glog.V(1) {
			debugVal = 1 // ufs.StartServer takes int debug param
		}
		api, err := machine.NewAPIClient()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error getting client: %s\n", err)
			os.Exit(1)
		}
		defer api.Close()
		fmt.Printf("Setting up hostmount on %s...\n", mountString)

		path := os.Args[0]
		debugVal = 0
		if glog.V(8) {
			debugVal = 1
		}
		mountCmd := exec.Command(path, "mount", mountString, fmt.Sprintf("-p=%s", viper.GetString(cfg.MachineProfile)), fmt.Sprintf("--v=%d", debugVal))
		mountCmd.Env = append(os.Environ(), constants.IsMinikubeChildProcess+"=true")
		if glog.V(8) {
			mountCmd.Stdout = os.Stdout
			mountCmd.Stderr = os.Stderr
		}
		err = mountCmd.Start()
		if err != nil {
			glog.Errorf("Error running command minikube mount %s", err)
			cmdUtil.MaybeReportErrorAndExit(err)
		}
		err = ioutil.WriteFile(filepath.Join(constants.GetMinipath(), constants.MountProcessFileName), []byte(strconv.Itoa(mountCmd.Process.Pid)), 0644)
		if err != nil {
			glog.Errorf("Error writing mount process pid to file: %s", err)
			cmdUtil.MaybeReportErrorAndExit(err)
		}
	},
}

func init() {
	RootCmd.AddCommand(mountDetachCmd)
}
