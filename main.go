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
	"time"

	"golang.org/x/sync/errgroup"
)

var TELNET_ADDR = "koukoku.shadan.open.ad.jp:992"
var BOM = [...]byte{0xEF, 0xBB, 0xBF}

func main() {
	// if write/read fails, the current loop will be end and new one will run
	// i.e. reconnecting will be triggered.
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

		var eg errgroup.Group

		// read
		eg.Go(func() error {
			scanner := bufio.NewScanner(conn)
			for scanner.Scan() {
				fmt.Println(scanner.Text())
			}

			if err := scanner.Err(); err != nil {
				return fmt.Errorf("error while reading: %v", err)
			}

			return nil
		})

		// write
		eg.Go(func() error {
			scanner := bufio.NewScanner(os.Stdin)
			for scanner.Scan() {
				fmt.Fprintln(conn, scanner.Text())
			}

			if err := scanner.Err(); err != nil {
				return fmt.Errorf("error while writing: %v", err)
			}

			return nil
		})

		if err := eg.Wait(); err != nil {
			log.Println(err)
			time.Sleep(time.Second)
		}
	}
}
