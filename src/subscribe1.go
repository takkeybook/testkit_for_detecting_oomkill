/*
 * Original code is written by fits-san
 *   https://fits.hatenablog.com/entry/2020/12/06/211946
 */
package main

import (
	"context"
	"log"

	"github.com/containerd/containerd"
	"github.com/containerd/containerd/events"
	"github.com/containerd/typeurl"
	_ "github.com/containerd/containerd/api/events"
)

const (
	address = "/run/containerd/containerd.sock"
	namespace = "k8s.io"
)

func printEnvelope(env *events.Envelope) {
	event, err := typeurl.UnmarshalAny(env.Event)

	if err != nil {
		log.Printf("ERROR unmarshal %v", err)
	}

	log.Printf(
		"topic = %s, namespace = %s, event.typeurl = %s, event = %v",
		env.Topic, env.Namespace, env.Event.TypeUrl, event,
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

	defer client.Close()

	ctx := context.Background()

	ch, errs := client.Subscribe(ctx)

	for {
		select {
		case env := <-ch:
			printEnvelope(env)
		case e := <-errs:
			log.Fatal(e)
		}
	}
}
