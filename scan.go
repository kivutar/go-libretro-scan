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

var found uint64

func loadDB(dir string) [][]Game {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		log.Fatal(err)
	}

	var DB = [][]Game{}
	for _, f := range files {
		rdb := parseRDB(dir + f.Name())
		DB = append(DB, rdb)
	}

	return DB
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

func findInDB(DB [][]Game, CRC32 uint32) {
	var wg sync.WaitGroup
	wg.Add(len(DB))
	for _, rdb := range DB {
		go func(rdb []Game, CRC32 uint32) {
			for _, game := range rdb {
				if CRC32 == game.ROM.CRC32 {
					//fmt.Printf("Found %s\n", game.Name)
					found++
				}
			}
			wg.Done()
		}(rdb, CRC32)
	}
	wg.Wait()
}

func main() {
	rompath := flag.String("roms", "", "Path to the folder you want to scan.")
	rdbpath := flag.String("rdbs", "", "Path to the folder containing the RDB files.")
	flag.Parse()

	start := time.Now()

	DB := loadDB(*rdbpath)
	roms := allFilesIn(*rompath)
	fmt.Println(len(DB), "RDB files and", len(roms), "zips to scan.")

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
					findInDB(DB, rom.CRC32)
				}
			}
			z.Close()
		}
	}

	elapsed2 := time.Since(scanstart)
	fmt.Println("Scanning ROMs took", elapsed2)
	fmt.Println("Found", found)
}
