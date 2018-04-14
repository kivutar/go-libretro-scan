package main

import (
	"archive/zip"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
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
		//fmt.Println(rdb)
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

func findInDat(dat []Game, CRC32 uint32) {
	for _, game := range dat {
		if CRC32 == game.ROM.CRC32 {
			//fmt.Printf("Found %s\n", game.Name)
			found++
		}
	}
}

func findInDB(DB [][]Game, CRC32 uint32) {
	for _, dat := range DB {
		findInDat(dat, CRC32)
	}
}

func main() {

	start := time.Now()

	DB := loadDB("libretro-database/rdb/")
	roms := allFilesIn("../Downloads/No-Intro/")
	fmt.Println(len(DB), "RDB files and", len(roms), "zips to scan.")

	elapsed := time.Since(start)
	fmt.Println("Loading DB took ", elapsed)
	scanstart := time.Now()

	for _, f := range roms {
		z, _ := zip.OpenReader(f)
		for _, rom := range z.File {
			if rom.CRC32 > 0 {
				findInDB(DB, rom.CRC32)
			}
		}
		z.Close()
	}

	elapsed2 := time.Since(scanstart)
	fmt.Println("Scanning ROMs took", elapsed2)
	fmt.Println("Found", found)
}
