package main

import (
	"archive/zip"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type ROM struct {
	Name  string
	CRC32 uint32
	Size  uint64
}

type Game struct {
	Name        string
	Description string
	Genre       string
	Developer   string
	Publisher   string
	Franchise   string
	Serial      string
	ROM         ROM
}

type RDB []Game
type DB map[string]RDB

var found uint64

func loadDB(dir string) DB {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		log.Fatal(err)
	}

	db := make(DB)
	for _, f := range files {
		filename := f.Name()
		system := filename[0 : len(filename)-4]
		db[system] = parseRDB(dir + f.Name())
	}

	return db
}

func allFilesIn(dir string) []string {
	roms := []string{}
	filepath.Walk(dir, func(path string, f os.FileInfo, err error) error {
		if !f.IsDir() {
			roms = append(roms, path)
		}
		return nil
	})
	return roms
}

func findInDB(db DB, rompath string, romname string, CRC32 uint32) {
	var wg sync.WaitGroup
	wg.Add(len(db))
	for system, rdb := range db {
		go func(rdb RDB, CRC32 uint32, system string) {
			for _, game := range rdb {
				if CRC32 == game.ROM.CRC32 {
					fmt.Printf("%s#%s\n", rompath, romname)
					fmt.Printf("%s\n", game.Name)
					fmt.Printf("DETECT\n")
					fmt.Printf("DETECT\n")
					fmt.Printf("%d|crc\n", CRC32)
					fmt.Printf("%s.lpl\n", system)
					found++
				}
			}
			wg.Done()
		}(rdb, CRC32, system)
	}
	wg.Wait()
}

func main() {
	rompath := flag.String("roms", "", "Path to the folder you want to scan.")
	rdbpath := flag.String("rdbs", "", "Path to the folder containing the RDB files.")
	flag.Parse()

	start := time.Now()

	db := loadDB(*rdbpath)
	roms := allFilesIn(*rompath)
	fmt.Println(len(db), "RDB files and", len(roms), "zips to scan.")

	elapsed := time.Since(start)
	fmt.Println("Loading DB took ", elapsed)
	scanstart := time.Now()

	for _, f := range roms {
		ext := filepath.Ext(f)
		switch ext {
		case ".zip":
			z, _ := zip.OpenReader(f)
			for _, rom := range z.File {
				if rom.CRC32 > 0 {
					findInDB(db, f, rom.Name, rom.CRC32)
				}
			}
			z.Close()
		}
	}

	elapsed2 := time.Since(scanstart)
	fmt.Println("Scanning ROMs took", elapsed2)
	fmt.Println("Found", found)
}
