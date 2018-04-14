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
		fmt.Println(int(rdb[pos]))
		g := Game{ROM: ROM{}}

		nfields := int(rdb[pos]) - 0x80
		fmt.Println("nfields: ", nfields)
		pos++

		for i := 0; i <= nfields; i++ {

			len := int(rdb[pos]) - 0xA0
			pos++
			key := rdb[pos : pos+len]
			fmt.Println(string(key[:]))
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
					g.Name = string(value[:])
					pos += len
				} else {
					len := int(rdb[pos]) - 0xA0
					pos++
					value := rdb[pos : pos+len]
					fmt.Println(string(value[:]))
					g.Name = string(value[:])
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
					g.Name = string(value[:])
					pos += len
				} else {
					len := int(rdb[pos]) - 0xA0
					pos++
					value := rdb[pos : pos+len]
					fmt.Println(string(value[:]))
					g.Name = string(value[:])
					pos += len
				}
			case "size":
				pos++
				len := 4
				pos += len
			case "releaseyear":
				pos++
				len := 2
				pos += len
			case "crc":
				pos++
				pos++
				len := 4
				value := rdb[pos : pos+len]
				str := fmt.Sprintf("%#x", string(value[:]))
				u64, _ := strconv.ParseUint(str, 16, 32)
				g.ROM.CRC32 = uint32(u64)
				pos += len
			case "md5":
				pos++
				pos++
				len := 16
				pos += len
			case "sha1":
				pos++
				pos++
				len := 20
				pos += len
			}
		}
		output = append(output, g)
	}

	//	<83>
	//		<A4>name<AF>Lutris Launcher
	//		<AB>description<AF>Lutris Launcher
	//		<A8>rom_name<B1>love-lutris.lutro
	//	<83>
	//		<A4>name<A6>Tetris
	//		<AB>description<A6>Tetris
	//		<A8>rom_name<B2>lutro-tetris.lutro
	//	<84>
	//		<A4>name<A9>Spaceship
	//		<AB>description<A9>Spaceship
	//		<A8>rom_name<B5>lutro-spaceship.lutro
	//		<A9>developer<B3>Jean-André Santoni
	//	<84>
	//		<A4>name<A5>Snake
	//		<AB>description<A5>Snake
	//		<A8>rom_name<B1>lutro-snake.lutro
	//		<A9>developer<B3>Jean-André Santoni

	return output
}

func loadDB(dir string) [][]Game {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		log.Fatal(err)
	}

	var DB = [][]Game{}
	for _, f := range files[11:12] {
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
