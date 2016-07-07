package main

import (
	"flag"
	"log"
	"net/http"
	_ "net/http/pprof"
	"sync"
	"time"

	"golang.org/x/net/context"
	"google.golang.org/cloud/bigtable"
)

var para = flag.Int("para", 100, "spawn para routine")

func main() {
	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()

	start := time.Now()
	defer func() {
		if time.Since(start) > 2*time.Minute {
			log.Println("**************************************************")
			log.Printf("%s elapsed!  looks like a failure!", time.Since(start))
			log.Println("**************************************************")
		}
	}()

	flag.Parse()

	ctx := context.Background()
	client, err := bigtable.NewClient(ctx, Project, Instance)
	if err != nil {
		log.Fatalf("Failed to create bigtable client: %v", err)
	}

	tbl := client.Open(Table)

	m := bigtable.NewMutation()
	m.Set(Family, Column, bigtable.Now(), []byte("foobar"))
	if err := tbl.Apply(ctx, Row, m); err != nil {
		log.Fatalf("failed to Set: %s", err)
	}

	var wg sync.WaitGroup

	ctx2, cancel := context.WithTimeout(ctx, 800*time.Millisecond)

	go func() {
		time.Sleep(600 * time.Millisecond)
		cancel()
	}()

	for i := 0; i < *para; i++ {
		wg.Add(1)
		go func() {
			_, err = tbl.ReadRow(ctx2, Row)
			if err != nil {
				log.Println(err)
			}
			wg.Done()
		}()
	}
	wg.Wait()

	log.Println("1st pass ok waiting 5s")
	time.Sleep(5 * time.Second)
	log.Println("2nd pass")

	for i := 0; i < *para; i++ {
		wg.Add(1)
		go func() {
			_, err = tbl.ReadRow(ctx2, Column)
			if err != nil {
				log.Println(err)
			} else {
				log.Println("ok")
			}
			wg.Done()
		}()
	}
	wg.Wait()

	log.Println("starting last call on ctx")
	_, err = tbl.ReadRow(ctx, Row)
	if err != nil {
		log.Println("last", err)
	} else {
		log.Println("last ok")
	}
}
