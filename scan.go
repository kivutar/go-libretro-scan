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

// DB is a database that contains many RDB, mapped to their system name
type DB map[string]RDB

// RDB contains all the game descriptions for a system
type RDB []Game

// Game represents a game in the libretro database
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

// Number of ROMs matched
var matched uint64

// loadDB loops over the RDBs in a given directory and parses them
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

// allFilesIn recursively builds a list of the files in a given directory
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

// writePlaylistEntry writes a playlist entry
func writePlaylistEntry(rompath string, romname string, gamename string, CRC32 uint32, system string) {
	CRC32Str := strconv.FormatUint(uint64(CRC32), 10)
	lpl, _ := os.OpenFile("playlists/"+system+".lpl", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	lpl.WriteString(rompath + "#" + romname + "\n")
	lpl.WriteString(gamename + "\n")
	lpl.WriteString("DETECT\n")
	lpl.WriteString("DETECT\n")
	lpl.WriteString(CRC32Str + "|crc\n")
	lpl.WriteString(system + ".lpl\n")
	lpl.Close()
}

// findInDB loops over the RDBs in the DB and concurrently matches CRC32 checksums.
func findInDB(db DB, rompath string, romname string, CRC32 uint32) {
	var wg sync.WaitGroup
	wg.Add(len(db))
	// For every RDB in the DB
	for system, rdb := range db {
		go func(rdb RDB, CRC32 uint32, system string) {
			// For each game in the RDB
			for _, game := range rdb {
				// If the checksums match
				if CRC32 == game.CRC32 {
					// Write the playlist entry
					writePlaylistEntry(rompath, romname, game.Name, CRC32, system)
					matched++
				}
			}
			wg.Done()
		}(rdb, CRC32, system)
	}
	// Synchronize all the goroutines
	wg.Wait()
}

func main() {
	rompath := flag.String("roms", "", "Path to the folder you want to scan.")
	rdbpath := flag.String("rdbs", "", "Path to the folder containing the RDB files.")
	//lplpath := flag.String("playlists", "", "Path to the folder where playlists will be generated")
	flag.Parse()

	start := time.Now()

	db := loadDB(*rdbpath)
	roms := allFilesIn(*rompath)
	fmt.Println(len(db), "RDB files and", len(roms), "zips to scan.")

	elapsed := time.Since(start)
	fmt.Println("Loading DB took ", elapsed)
	scanstart := time.Now()

	// For every rom found in the rom path
	for _, f := range roms {
		ext := filepath.Ext(f)
		switch ext {
		case ".zip":
			// Open the ZIP archive
			z, _ := zip.OpenReader(f)
			for _, rom := range z.File {
				if rom.CRC32 > 0 {
					// Look for a matching game entry in the database
					findInDB(db, f, rom.Name, rom.CRC32)
				}
			}
			z.Close()
		}
	}

	elapsed2 := time.Since(scanstart)
	fmt.Println("Scanning ROMs took", elapsed2)
	fmt.Println("Matched:", matched)
}
