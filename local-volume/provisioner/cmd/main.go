/*
Copyright 2017 The Kubernetes Authors.

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

package main

import (
	"flag"
	"os"

	"github.com/golang/glog"
	"github.com/kubernetes-incubator/external-storage/local-volume/provisioner/pkg/controller"
	"github.com/kubernetes-incubator/external-storage/local-volume/provisioner/pkg/types"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

func setupClient() *kubernetes.Clientset {
	config, err := rest.InClusterConfig()
	if err != nil {
		glog.Fatalf("Error creating InCluster config: %v\n", err)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		glog.Fatalf("Error creating clientset: %v\n", err)
	}
	return clientset
}

func main() {
	flag.Set("logtostderr", "true")
	flag.Parse()

	node := os.Getenv("MY_NODE_NAME")
	if node == "" {
		glog.Fatalf("MY_NODE_NAME environment variable not set\n")
	}

	client := setupClient()

	glog.Info("Starting controller\n")
	controller.StartLocalController(client, &types.UserConfig{
		NodeName:     node,
		HostDir:      "/mnt/disks",
		MountDir:     "/local-disks",
		DiscoveryMap: createDiscoveryMap(),
	})
}

func createDiscoveryMap() map[string]string {
	m := make(map[string]string)
	// Default setting
	m["local-storage"] = ""
	return m
}
