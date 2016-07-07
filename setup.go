package main

import (
	"log"

	"golang.org/x/net/context"
	"google.golang.org/cloud/bigtable"
)

func main() {
	ctx := context.Background()

	client, err := bigtable.NewAdminClient(ctx, Project, Instance)
	if err != nil {
		log.Fatalf("failed to create bigtable admin client: %s", err)
	}

	famExists := false

	if info, err := client.TableInfo(ctx, Table); err != nil {
		// failed to read table info - assume this means table does not exist
		if err := client.CreateTable(ctx, Table); err != nil {
			log.Fatalf("failed to create table %q: %s", Table, err)
		}
		log.Printf("created table %q", Table)
	} else {
		log.Printf("table %q already exists", Table)

		for _, fam := range info.Families {
			if fam == Family {
				famExists = true
				break
			}
		}
	}

	if famExists {
		log.Printf("CF %q already exists", Family)
	} else {
		if err := client.CreateColumnFamily(ctx, Table, Family); err != nil {
			log.Fatalf("failed to create CF %q: %s", Family, err)
		}
		log.Printf("created CF %q", Family)
	}

	log.Printf("good to go")
}
