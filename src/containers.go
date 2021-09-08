package main

import (
    "context"
    "encoding/json"
    "log"
    "fmt"

    "github.com/containerd/containerd"
)

const (
    address = "/run/containerd/containerd.sock"
    namespace = "k8s.io"
)

func main() {
    client, err := containerd.New(
        address,
        containerd.WithDefaultNamespace(namespace),
    )

    if err != nil {
        log.Fatal(err)
    }

    defer client.Close()

    ctx := context.Background()

     containers, _ := client.Containers(ctx)
     for _, c := range containers {
        ctxInfo := context.Background()
        ci, _ := c.Info(ctxInfo)
        l, _ := json.Marshal(ci.Labels)
        fmt.Printf("ID = %s, Labels = %s\n", ci.ID, string(l))
     }
}
