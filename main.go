package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/markus-wa/demoinfocs-golang/v5/pkg/demoinfocs"
	"github.com/markus-wa/demoinfocs-golang/v5/pkg/demoinfocs/common"
	"github.com/markus-wa/demoinfocs-golang/v5/pkg/demoinfocs/events"
)

func main() {
	filepath := os.Args[1]
	f, err := os.Open(filepath)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	p := demoinfocs.NewParser(f)
	defer p.Close()

	var ctIDs []int
	var tIDs []int

	p.RegisterEventHandler(func(e events.MatchStart) {
		participants := p.GameState().Participants()
		cts := participants.TeamMembers(common.TeamCounterTerrorists)
		ts := participants.TeamMembers(common.TeamTerrorists)

		for _, ct := range cts {
			ctIDs = append(ctIDs, ct.EntityID)
		}
		for _, t := range ts {
			tIDs = append(tIDs, t.EntityID)
		}
	})

	err = p.ParseToEnd()
	if err != nil {
		panic(err)
	}

	ctMask := idsToBitmask(ctIDs)
	tMask := idsToBitmask(tIDs)

	fmt.Println("Command to hear counter-terrorists:")
	fmt.Printf("tv_listen_voice_indices %d\n", ctMask)
	fmt.Printf("tv_listen_voice_indices_h %d\n", ctMask)

	fmt.Println()
	fmt.Println("Command to hear terrorists:")
	fmt.Printf("tv_listen_voice_indices %d\n", tMask)
	fmt.Printf("tv_listen_voice_indices_h %d\n", tMask)

	fmt.Println()
	fmt.Println("Command to hear both:")
	fmt.Println("tv_listen_voice_indices -1")
	fmt.Println("tv_listen_voice_indices_h -1")

	waitForExit()
}

func waitForExit() {
	fmt.Println()
	fmt.Println("Press Enter to exit...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')
}

func idsToBitmask(ids []int) uint32 {
	var mask uint32
	for _, id := range ids {
		if id >= 1 && id <= 32 {
			mask |= 1 << (id - 1)
		}
	}
	return mask
}
