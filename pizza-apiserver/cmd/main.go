package main

import (
	"github.com/ZhengHe-MD/programming-kubernetes/pizza-apiserver/pkg/server"
	genericapiserver "k8s.io/apiserver/pkg/server"
	"k8s.io/component-base/logs"
	"k8s.io/klog/v2"
)

func main() {
	logs.InitLogs()
	defer logs.FlushLogs()

	stopCh := genericapiserver.SetupSignalHandler()
	options := server.NewCustomServerOptions()
	cmd := server.NewCommandStartCustomServer(options, stopCh)
	if err := cmd.Execute(); err != nil {
		klog.Fatal(err)
	}
}
