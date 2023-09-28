// Copyright (C) 2023 Yasuhiro Matsumoto (a.k.a. mattn)
// Copyright (C) 2023 Aoi Asagi (ARGI-BERRI)
// The software is redistributable under the condition of MIT License.

package main

import (
	"bufio"
	"crypto/tls"
	"fmt"
	"log"
	"os"
	"sync"
	"time"
)

var TELNET_ADDR = "koukoku.shadan.open.ad.jp:992"
var BOM = [...]byte{0xEF, 0xBB, 0xBF}

func main() {
	for {
		conn, err := tls.Dial("tcp", TELNET_ADDR, nil)

		if err != nil {
			log.Printf("Failed to connect: %v\n", err)
			return
		}

		defer conn.Close()

		log.Printf("connected to %v\n", TELNET_ADDR)

		// Send nobody command to suppress the notice
		fmt.Fprintln(conn, "nobody")

		var wg sync.WaitGroup

		// Read
		wg.Add(1)
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

		// Write
		wg.Add(1)
		go func() {
			defer wg.Done()

			scanner := bufio.NewScanner(os.Stdin)
			for scanner.Scan() {
				fmt.Fprintln(conn, scanner.Text())
			}
		}()

		// Attempt to write a byte into the connection to check if it's still open
		wg.Add(1)
		go func() {
			defer wg.Done()

			for {
				conn.SetDeadline(time.Now().Add(time.Second))
				_, err := conn.Read(make([]byte, 0))
				if err != nil {
					log.Printf("error: %v\n", err)
					break
				}

				time.Sleep(time.Second)
			}
		}()

		wg.Wait()
		time.Sleep(time.Second * 1)
	}
}
