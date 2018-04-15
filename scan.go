package main

import (
	"archive/zip"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"sync"
	"time"
)

type DB map[string]RDB

type RDB []Game

type Game struct {
	Name        string
	Description string
	Genre       string
	Developer   string
	Publisher   string
	Franchise   string
	Serial      string
	ROMName     string
	Size        uint64
	CRC32       uint32
}

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
				if CRC32 == game.CRC32 {
					CRC32Str := strconv.FormatUint(uint64(CRC32), 10)
					lpl, _ := os.OpenFile("playlists/"+system+".lpl", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
					lpl.WriteString(rompath + "#" + romname + "\n")
					lpl.WriteString(game.Name + "\n")
					lpl.WriteString("DETECT\n")
					lpl.WriteString("DETECT\n")
					lpl.WriteString(CRC32Str + "|crc\n")
					lpl.WriteString(system + ".lpl\n")
					lpl.Close()
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
