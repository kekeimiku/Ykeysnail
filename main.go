package main

import (
	"fmt"
	"time"
	"ykeysnail/window/wayland"
)

func main() {
	go func() {
		for {
			select {
			case event := <-wayland.SubscriptionSwayIPC.Events:
				if event.Container.AppID == "Alacritty" && event.Container.Name == "Alacritty" {
					//不同窗口的按键逻辑在这
					fmt.Println("hello alacritty")
				}
			case err := <-wayland.SubscriptionSwayIPC.Errors:
				fmt.Println("Error:", err)
			}
		}
	}()

	for {
		fmt.Println("adadad")
		time.Sleep(time.Second)
	}
}
