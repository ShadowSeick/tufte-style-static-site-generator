package flags

import (
	"flag"
	"fmt"
	"slices"
	"sync"
)

var (
	initialize sync.Once
	flags []Flag
)

type Flag uint8

const (
	Debug Flag = iota
	FlagCount
)

var flagsString = [FlagCount]string{
	Debug: "debug",
}

var flagType = [FlagCount]FlagType{
	Debug: FlagTypeBool,
}

func (f Flag) String() string {
	if f >= FlagCount {
		panic("invalid flag")
	}
	return flagsString[f]
}

func (f Flag) Type() FlagType {
	if f >= FlagCount {
		panic("invalid flag")
	}
	return flagType[f]
}

type FlagType uint8

const (
	FlagTypeBool FlagType = iota
	FlagTypeCount
)

// Think of a better way of doing this
func Init() {
	initialize.Do(func() {
		for i := range FlagCount {
			switch i.Type() {
			case FlagTypeBool:
				boolean := flag.Bool(i.String(), false, "")
				flag.Parse()
				if *boolean {
					flags = append(flags, i)
				}
			default:
				panic("flag not handled")
			}
		}
	})
}

func IsSet(flag Flag) bool {
	return slices.Contains(flags, flag)
}
