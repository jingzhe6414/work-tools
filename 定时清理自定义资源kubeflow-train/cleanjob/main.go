package main

import (
	"cleanjob/server"
)

func main() {
	nsList := server.GetNamespaceList()
	for _, item := range nsList.Items {
		ns := item.Name
		server.CleanTfJob(ns)
		server.MSJob(ns)
		server.PaddleJob(ns)
		server.PyTorchJob(ns)
	}
}
