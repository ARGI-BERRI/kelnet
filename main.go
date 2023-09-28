// Copyright (C) 2023 Yasuhiro Matsumoto (a.k.a. mattn)
// Copyright (C) 2023 Aoi Asagi (ARGI-BERRI)
// The software is redistributable under the condition of MIT License.

package main

import (
	"bufio"
	"crypto/tls"
	"fmt"
	"sync"
	"time"
)

var TELNET_ADDR = "koukoku.shadan.open.ad.jp:992"
var BOM = [...]byte{0xEF, 0xBB, 0xBF}

func main() {
	for {
		conn, err := tls.Dial("tcp", TELNET_ADDR, nil)

		if err != nil {
			fmt.Printf("[%s] Failed to connect: %v\n", timestamp(), err)
			return
		}

		defer conn.Close()

		// Send nobody command to suppress the notice
		fmt.Fprintln(conn, "nobody")

		var wg sync.WaitGroup

		wg.Add(1)
		// Receives bytes from the server
		go func() {
			defer wg.Done()

			scanner := bufio.NewScanner(conn)

			for scanner.Scan() {
				text := scanner.Text()
				bytes := scanner.Bytes()

				// Ignore the line if it includes BEL character
				if text == "\a" {
					continue
				}

				// Ignore the line if it includes BOM character
				if len(bytes) >= len(BOM) &&
					bytes[0] == BOM[0] &&
					bytes[1] == BOM[1] &&
					bytes[2] == BOM[2] {
					continue
				}

				fmt.Println(text)
			}
		}()

		// Attempt to write a byte into the connection to check if it's still open
		go func() {
			defer wg.Done()

			for {
				conn.SetWriteDeadline(time.Now().Add(time.Second))
				_, err := conn.Write([]byte("..."))

				if err != nil {
					fmt.Printf("[%s] Disconnected: %v\n", timestamp(), err)
					break
				}

				fmt.Printf("[%s] Healthcheck is done.\n", timestamp())
				time.Sleep(time.Second * 3)
			}
		}()

		wg.Wait()
	}
}

func timestamp() string {
	return time.Now().Format("15:04:05")
}
