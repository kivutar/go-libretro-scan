package main

import (
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

func parseRDB(path string) []Game {
	rdb, _ := ioutil.ReadFile(path)

	var output []Game

	pos := 0x10

	for int(rdb[pos]) != 192 {
		fmt.Println("POSITION:", int(rdb[pos]))
		g := Game{ROM: ROM{}}

		nfields := int(rdb[pos]) - 0x80
		fmt.Println("Number of fields: ", nfields)
		pos++

		for i := 0; i < nfields; i++ {

			len := int(rdb[pos]) - 0xA0
			pos++
			key := rdb[pos : pos+len]
			fmt.Println("KEY:", string(key[:]))
			pos += len

			switch string(key[:]) {
			case "name":
				if int(rdb[pos]) == 0xD9 {
					fmt.Println("String")
					pos++
					len := int(rdb[pos])
					pos++
					value := rdb[pos : pos+len]
					fmt.Println(string(value[:]))
					g.Name = string(value[:])
					pos += len
				} else {
					fmt.Println("Raw")
					len := int(rdb[pos]) - 0xA0
					pos++
					value := rdb[pos : pos+len]
					fmt.Println(string(value[:]))
					g.Name = string(value[:])
					pos += len
				}
			case "description":
				if int(rdb[pos]) == 0xD9 {
					fmt.Println("String")
					pos++
					len := int(rdb[pos])
					pos++
					value := rdb[pos : pos+len]
					fmt.Println(string(value[:]))
					g.Description = string(value[:])
					pos += len
				} else {
					fmt.Println("Raw")
					len := int(rdb[pos]) - 0xA0
					pos++
					value := rdb[pos : pos+len]
					fmt.Println(string(value[:]))
					g.Description = string(value[:])
					pos += len
				}
			case "genre":
				if int(rdb[pos]) == 0xD9 {
					fmt.Println("String")
					pos++
					len := int(rdb[pos])
					pos++
					value := rdb[pos : pos+len]
					fmt.Println(string(value[:]))
					g.Genre = string(value[:])
					pos += len
				} else {
					fmt.Println("Raw")
					len := int(rdb[pos]) - 0xA0
					pos++
					value := rdb[pos : pos+len]
					fmt.Println(string(value[:]))
					g.Genre = string(value[:])
					pos += len
				}
			case "developer":
				if int(rdb[pos]) == 0xD9 {
					fmt.Println("String")
					pos++
					len := int(rdb[pos])
					pos++
					value := rdb[pos : pos+len]
					fmt.Println(string(value[:]))
					g.Developer = string(value[:])
					pos += len
				} else {
					fmt.Println("Raw")
					len := int(rdb[pos]) - 0xA0
					pos++
					value := rdb[pos : pos+len]
					fmt.Println(string(value[:]))
					g.Developer = string(value[:])
					pos += len
				}
			case "publisher":
				if int(rdb[pos]) == 0xD9 {
					fmt.Println("String")
					pos++
					len := int(rdb[pos])
					pos++
					value := rdb[pos : pos+len]
					fmt.Println(string(value[:]))
					g.Publisher = string(value[:])
					pos += len
				} else {
					fmt.Println("Raw")
					len := int(rdb[pos]) - 0xA0
					pos++
					value := rdb[pos : pos+len]
					fmt.Println(string(value[:]))
					g.Publisher = string(value[:])
					pos += len
				}
			case "rom_name":
				if int(rdb[pos]) == 0xD9 {
					fmt.Println("String")
					pos++
					len := int(rdb[pos])
					pos++
					value := rdb[pos : pos+len]
					fmt.Println(string(value[:]))
					g.ROM.Name = string(value[:])
					pos += len
				} else {
					len := int(rdb[pos]) - 0xA0
					pos++
					value := rdb[pos : pos+len]
					fmt.Println(string(value[:]))
					g.ROM.Name = string(value[:])
					pos += len
				}
			case "size":
				len := int(rdb[pos]) - 0xCA // CE -> 4
				fmt.Println(len)
				pos++
				pos += len
			case "releaseyear":
				len := int(rdb[pos]) - 0xCB // CD -> 2
				fmt.Println(len)
				pos++
				pos += len
			case "releasemonth":
				len := int(rdb[pos]) - 0xCB // CC -> 1
				fmt.Println(len)
				pos++
				pos += len
			case "users":
				len := int(rdb[pos]) - 0xCB // CC -> 1
				fmt.Println(len)
				pos++
				pos += len
			case "crc":
				if int(rdb[pos]) == 0xC4 {
					pos++
					len := int(rdb[pos])
					pos++
					value := rdb[pos : pos+len]
					str := fmt.Sprintf("%#x", string(value[:]))
					u64, _ := strconv.ParseUint(str, 16, 32)
					g.ROM.CRC32 = uint32(u64)
					pos += len
				}
			case "md5":
				if int(rdb[pos]) == 0xC4 {
					pos++
					len := int(rdb[pos])
					pos++
					value := rdb[pos : pos+len]
					fmt.Println(value)
					pos += len
				}
			case "sha1":
				if int(rdb[pos]) == 0xC4 {
					pos++
					len := int(rdb[pos])
					fmt.Println(len)
					pos++
					value := rdb[pos : pos+len]
					fmt.Println(value)
					pos += len
				}
			}
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
	for _, f := range files[43:44] {
		dat := parseRDB(dir + f.Name())
		DB = append(DB, dat)
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
			fmt.Printf("Found %v\n", game.Name)
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
	fmt.Println(len(DB), len(roms))

	elapsed := time.Since(start)
	fmt.Println("Loading DB took ", elapsed)
	// scanstart := time.Now()

	// for _, f := range roms {
	// 	z, _ := zip.OpenReader(f)
	// 	for _, rom := range z.File {
	// 		go findInDB(DB, rom.CRC32)
	// 	}
	// }

	// elapsed2 := time.Since(scanstart)
	// fmt.Println("Scanning ROMs took ", elapsed2)
}
