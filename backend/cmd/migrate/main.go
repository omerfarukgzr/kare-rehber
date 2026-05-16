package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/golang-migrate/migrate/v4"

	"github.com/koc-luk/backend/internal/config"
	"github.com/koc-luk/backend/internal/dbmigrate"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("config: %v", err)
	}

	flag.Parse()
	args := flag.Args()
	cmd := "up"
	if len(args) > 0 {
		cmd = args[0]
	}

	m, err := dbmigrate.New(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("migrate init: %v", err)
	}
	defer func() {
		srcErr, dbErr := m.Close()
		if srcErr != nil {
			log.Printf("source close: %v", srcErr)
		}
		if dbErr != nil {
			log.Printf("db close: %v", dbErr)
		}
	}()

	switch cmd {
	case "up":
		if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
			log.Fatalf("up: %v", err)
		}
		fmt.Println("migrations: up")
	case "down":
		steps := 1
		if len(args) > 1 {
			n, err := strconv.Atoi(args[1])
			if err != nil {
				log.Fatalf("down step: %v", err)
			}
			steps = n
		}
		if err := m.Steps(-steps); err != nil && !errors.Is(err, migrate.ErrNoChange) {
			log.Fatalf("down: %v", err)
		}
		fmt.Printf("migrations: down %d\n", steps)
	case "drop":
		if err := m.Drop(); err != nil {
			log.Fatalf("drop: %v", err)
		}
		fmt.Println("migrations: dropped")
	case "version":
		v, dirty, err := m.Version()
		if err != nil {
			log.Fatalf("version: %v", err)
		}
		fmt.Printf("version=%d dirty=%v\n", v, dirty)
	case "force":
		if len(args) < 2 {
			log.Fatal("force requires version arg")
		}
		v, err := strconv.Atoi(args[1])
		if err != nil {
			log.Fatalf("force: %v", err)
		}
		if err := m.Force(v); err != nil {
			log.Fatalf("force: %v", err)
		}
		fmt.Printf("forced to version %d\n", v)
	default:
		fmt.Fprintf(os.Stderr, "unknown command: %s\n", cmd)
		os.Exit(2)
	}
}
