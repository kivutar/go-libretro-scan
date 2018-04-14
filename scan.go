package main

import (
	"archive/zip"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
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

func parseDAT(path string) []Game {
	dat, _ := ioutil.ReadFile(path)

	r, _ := regexp.Compile(`(?s)game \((.*?)\n\)\n`)
	games := r.FindAllStringSubmatch(string(dat[:]), -1)

	var output []Game

	for _, game := range games {
		r2, _ := regexp.Compile(`\tname "(.*?)"`)
		name := r2.FindStringSubmatch(game[1])

		r3, _ := regexp.Compile(`(?s)\trom \( (.*?) \)`)
		rom := r3.FindStringSubmatch(game[1])

		r4, _ := regexp.Compile(`crc (\w*?) `)
		crc := r4.FindStringSubmatch(rom[1])

		u64, err := strconv.ParseUint(crc[1], 16, 32)
		if err != nil {
			fmt.Println(err)
		}

		g := Game{
			Name: name[1],
			ROM: ROM{
				CRC32: uint32(u64),
			},
		}

		output = append(output, g)
	}

	return output
}

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
			fmt.Printf("Found %s\n", game.Name)
		}
	}
}

func findInDB(DB [][]Game, CRC32 uint32) {
	for _, dat := range DB {
		go findInDat(dat, CRC32)
	}
}

func main() {

	start := time.Now()

	DB := loadDB("libretro-database/rdb/")
	roms := allFilesIn("../Downloads/No-Intro/")
	fmt.Println(len(DB), len(roms))

	elapsed := time.Since(start)
	fmt.Println("Loading DB took ", elapsed)
	scanstart := time.Now()

	for _, f := range roms {
		z, _ := zip.OpenReader(f)
		for _, rom := range z.File {
			findInDB(DB, rom.CRC32)
		}
		z.Close()
	}

	elapsed2 := time.Since(scanstart)
	fmt.Println("Scanning ROMs took ", elapsed2)
}
