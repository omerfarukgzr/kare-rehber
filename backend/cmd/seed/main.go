package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/koc-luk/backend/internal/config"
	"github.com/koc-luk/backend/internal/seed"
)

func main() {
	wipe := flag.Bool("wipe", false, "drop existing data before inserting (DANGER: silinir)")
	flag.Parse()

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("config: %v", err)
	}
	db, err := config.OpenDB(cfg)
	if err != nil {
		log.Fatalf("db: %v", err)
	}

	if *wipe {
		if err := seed.WipeAll(db); err != nil {
			log.Fatalf("wipe: %v", err)
		}
		fmt.Println("wipe: all data cleared")
	}

	if err := seed.SeedAdmin(db, cfg); err != nil {
		log.Fatalf("seed admin: %v", err)
	}
	fmt.Println("seed admin: ok")

	if err := seed.SeedWeeks(db); err != nil {
		log.Fatalf("seed weeks: %v", err)
	}
	fmt.Println("seed weeks: ok")

	if err := seed.SeedTestUsers(db, cfg); err != nil {
		log.Fatalf("seed test users: %v", err)
	}
	fmt.Println("seed test users: ok")

	fmt.Println("\n=== seed: done ===")
	printQuickRef()
}

func printQuickRef() {
	fmt.Print(`
Hızlı referans (login: telefon + şifre):

  Admin       05000000000 / Admin123!
  Koç         05311111111 / Test123!  (Ahmet Koç)
              05311111112 / Test123!  (Zeynep Koç)
              05311111113 / Test123!  (Mehmet Koç)
  Öğrenci     05322222221..05322222226 / Test123!
  Veli        05333333331..05333333333 / Test123!
  Koordinatör 05344444441..05344444442 / Test123!
`)
}
