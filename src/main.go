package main

import (
	eventHandle "event"
	"fmt"
	docker "github.com/fsouza/go-dockerclient"
	"os"
	"time"
)

func main() {
	// 设置endpoint
	//endpoint := "unix:///var/run/docker.sock"
	//	endpoint := common.GetConfig("Section", "endpoint")
	endpoint := os.Getenv("DNS_ENDPOINT")
	if endpoint == "" {
		fmt.Println("env DNS_ENDPOINT is empty")
		return
	}
	fmt.Println("DNS_ENDPOINT", endpoint)

	client, err := docker.NewClient(endpoint)
	if err != nil {
		fmt.Println(err)
	}

	eventChan := make(chan *docker.APIEvents, 100)
	defer close(eventChan)

	watching := false
	for {

		if client == nil {
			break
		}
		err := client.Ping()
		if err != nil {
			fmt.Printf("Unable to ping docker daemon: %s", err)
			if watching {
				client.RemoveEventListener(eventChan)
				watching = false
				client = nil
			}
			time.Sleep(10 * time.Second)
			break

		}

		if !watching {
			err = client.AddEventListener(eventChan)
			if err != nil && err != docker.ErrListenerAlreadyExists {
				fmt.Printf("Error registering docker event listener: %s", err)
				time.Sleep(10 * time.Second)
				continue
			}
			watching = true
			fmt.Println("Watching docker events")
		}

		select {
		case event := <-eventChan:
			if event == nil {
				if watching {
					client.RemoveEventListener(eventChan)
					watching = false
					client = nil
				}
				break
			}

			if event.Status == "start" {
				eventHandle.Start(client, event)
			} else if event.Status == "die" {
				eventHandle.Die(client, event)
			}

		case <-time.After(10 * time.Second):
			// check for docker liveness
		}
	}
}
