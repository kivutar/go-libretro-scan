package main

import (
	"fmt"
	"io/ioutil"
	"math"
	"strconv"
)

const (
	MPF_FIXMAP = 0x80
	MPF_MAP16  = 0xde
	MPF_MAP32  = 0xdf

	MPF_FIXARRAY = 0x90
	MPF_ARRAY16  = 0xdc
	MPF_ARRAY32  = 0xdd

	MPF_FIXSTR = 0xa0
	MPF_STR8   = 0xd9
	MPF_STR16  = 0xda
	MPF_STR32  = 0xdb

	MPF_BIN8  = 0xc4
	MPF_BIN16 = 0xc5
	MPF_BIN32 = 0xc6

	MPF_FALSE = 0xc2
	MPF_TRUE  = 0xc3

	MPF_INT8  = 0xd0
	MPF_INT16 = 0xd1
	MPF_INT32 = 0xd2
	MPF_INT64 = 0xd3

	MPF_UINT8  = 0xcc
	MPF_UINT16 = 0xcd
	MPF_UINT32 = 0xce
	MPF_UINT64 = 0xcf

	MPF_NIL = 0xc0
)

func parseRDB(path string) []Game {
	rdb, _ := ioutil.ReadFile(path)

	var output []Game

	pos := 0x10

	for int(rdb[pos]) != 192 && int(rdb[pos]) != 223 {
		//fmt.Println("\nPOSITION:", int(rdb[pos]))
		g := Game{ROM: ROM{}}

		nfields := int(rdb[pos]) - MPF_FIXMAP
		//fmt.Println("Number of fields: ", nfields)
		pos++

		for i := 0; i < nfields; i++ {

			len := int(rdb[pos]) - MPF_FIXSTR
			pos++
			key := rdb[pos : pos+len]
			//fmt.Println("KEY:", string(key[:]))
			pos += len

			fieldtype := int(rdb[pos])

			var value []byte

			if fieldtype < MPF_NIL {
				len := int(rdb[pos]) - MPF_FIXSTR
				pos++
				value = rdb[pos : pos+len]
				//fmt.Println(string(value[:]))
				pos += len
			}

			switch fieldtype {
			case MPF_STR8, MPF_STR16, MPF_STR32:
				pos++
				lenlen := fieldtype - MPF_STR8 + 1
				lenhex := fmt.Sprintf("%x", string(rdb[pos:pos+lenlen]))
				i64, _ := strconv.ParseInt(lenhex, 16, 32)
				len := int(i64)
				pos += lenlen
				value = rdb[pos : pos+len]
				//fmt.Println(string(value[:]))
				pos += len
			case MPF_UINT8, MPF_UINT16, MPF_UINT32, MPF_UINT64:
				pow := float64(rdb[pos]) - 0xC9
				len := int(math.Pow(2, pow)) / 8
				pos++
				value = rdb[pos : pos+len]
				//fmt.Println(value)
				pos += len
			case MPF_BIN8, MPF_BIN16, MPF_BIN32:
				pos++
				len := int(rdb[pos])
				pos++
				value = rdb[pos : pos+len]
				//fmt.Println(value)
				pos += len
			case MPF_MAP16, MPF_MAP32:
				pow := float64(rdb[pos]) - MPF_MAP16
				len := int(math.Pow(2, pow)) / 8
				pos++
				value = rdb[pos : pos+len]
				//fmt.Println(value)
				pos += len
			}

			switch string(key[:]) {
			case "name":
				g.Name = string(value[:])
			case "description":
				g.Description = string(value[:])
			case "genre":
				g.Genre = string(value[:])
			case "developer":
				g.Developer = string(value[:])
			case "publisher":
				g.Publisher = string(value[:])
			case "franchise":
				g.Franchise = string(value[:])
			case "serial":
				g.Serial = string(value[:])
			case "rom_name":
				g.ROM.Name = string(value[:])
			case "size":
				value2 := fmt.Sprintf("%x", string(value[:]))
				u64, _ := strconv.ParseUint(value2, 16, 32)
				g.ROM.Size = u64
			case "crc":
				value2 := fmt.Sprintf("%x", string(value[:]))
				u64, _ := strconv.ParseUint(value2, 16, 32)
				g.ROM.CRC32 = uint32(u64)
				//fmt.Println(uint32(u64))
			}
		}
		output = append(output, g)
	}

	return output
}
