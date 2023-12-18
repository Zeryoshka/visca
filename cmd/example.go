package main

import (
	"context"
	"fmt"
	"github.com/Zeryoshka/visca"
	"log"
	"time"
)

func main() {
	interval := 3 * time.Second
	controller, err := visca.NewController()
	camera, err := controller.AddCamera(
		"172.18.191.245:52381", 1, 1*time.Second,
	)
	if err != nil {
		log.Fatalln(err)
	}
	for {
		ctx := context.Background()
		fmt.Println("wait for, error:", err)
		time.Sleep(interval)
		err = camera.SendCommand(ctx, visca.CamWbCommand(visca.WbIndoor))
		fmt.Println("gg", err)
		time.Sleep(interval)
		err = camera.SendCommand(ctx, visca.CamWbCommand(visca.WbAuto1))
		fmt.Println("gg", err)
		time.Sleep(interval)
		err = camera.SendCommand(ctx, visca.CamWbCommand(visca.WbOutdoor))
		fmt.Println("gg", err)
		time.Sleep(interval)
		err = camera.SendCommand(ctx, visca.CamWbCommand(visca.WbAuto1))
		fmt.Println("gg", err)
	}
}

//func main() {
//	conn, err := net.Dial("udp", "172.18.191.245:52381")
//	fmt.Println(err)
//	header := []byte{
//		0x01, 0x00, 0x00, 0x06,
//		0x00, 0x00, 0x00, 0x03,
//		0x81,
//		0x01, 0x04, 0x35, 0x01,
//		0xFF,
//	}
//	n, err := conn.Write(header)
//	fmt.Println(n, err)
//}

//Data: 01000006000000008101043501ff
//Data: 01000006000000038101043501ff

//RESP
//Data: 011100031e000000 9051 ff

//Data: 0100 0006 da000000 8101043502ff
//Data: 0111 0003 da000000 9041ff
