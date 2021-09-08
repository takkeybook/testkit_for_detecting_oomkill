package main

import (
    "context"
    "fmt"
    "log"
    "os"
    "os/signal"
    "syscall"

    "github.com/containerd/containerd"
    "github.com/containerd/containerd/events"
    "github.com/containerd/typeurl"
    apievents "github.com/containerd/containerd/api/events"
)

const (
	address = "/run/containerd/containerd.sock"
	namespace = "k8s.io"
)

func printEnvelope(client *containerd.Client, env *events.Envelope) {
    event, err := typeurl.UnmarshalAny(env.Event)

    if err != nil {
        log.Printf("ERROR unmarshal %v", err)
    }

    var s string

    switch ev := event.(type) {
    case *apievents.TaskOOM:
        ctx := context.Background()
        container, _ := client.ContainerService().Get(ctx, ev.ContainerID)
        labels := container.Labels
        container_name := labels["io.kubernetes.container.name"]
        pod_name := labels["io.kubernetes.pod.name"]
        pod_namespace := labels["io.kubernetes.pod.namespace"]
        if container_name == "" || pod_name == "" {
            return
        }
        s = fmt.Sprintf(
            "{ container_id = %s, name = %s, namespace = %s }", 
            ev.ContainerID, pod_name, pod_namespace,
        )
    default:
        return
    }

    log.Printf(
        "topic = %s, namespace = %s, event.typeurl = %s, event = %v", 
        env.Topic, env.Namespace, env.Event.TypeUrl, s,
    )
}

func main() {
    client, err := containerd.New(
        address,
        containerd.WithDefaultNamespace(namespace),
    )

    if err != nil {
        log.Fatal(err)
    }

    //defer client.Close()
    defer func() {
        log.Print("close")
        client.Close()
    }()

    sig := make(chan os.Signal, 1)
    signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

    ctx, cancel := context.WithCancel(context.Background())

    go func() {
        s := <-sig
        log.Printf("syscall: %v", s)
        cancel()
    }()

    ch, errs := client.Subscribe(ctx)

    for {
        select {
        case env := <-ch:
            printEnvelope(client, env)
        case e := <-errs:
            log.Printf("ERROR %v", e)
            return
        }
    }
}
