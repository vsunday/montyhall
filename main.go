package main

import (
	"fmt"
	"math/rand"
	"time"
	"sync"
	"os"
	"strconv"
)

type game struct {
	win, player, master int
	score map[string]bool
}

func (g *game) init() *game {
	rand.Seed(time.Now().UnixNano())
	g.win = rand.Intn(3)
	return g
}

func (g *game) playerChoice() *game {
	rand.Seed(time.Now().UnixNano())
	g.player = rand.Intn(3)
	return g
}

func (g *game) masterChoice() *game {
	alts := map[int]bool{
		0: true,
		1: true,
		2: true,
	}
	delete(alts, g.player)
	delete(alts, g.win)
	rand.Seed(time.Now().UnixNano())
	choice := rand.Intn(len(alts))
	l := []int{}
	for k,_ := range alts {
		l = append(l, k)
	}
	g.master = l[choice]
	return g
}

func (g *game) judge() *game {
	g.score = map[string]bool{
		"stayed": g.win == g.player,
		"changed": g.win != g.player,
	}
	return g
}

func (g *game) playGame() *game {
	return g.init().playerChoice().masterChoice().judge()
}

type SafeCounter struct {
	c int
	m sync.Mutex
}

func (sc *SafeCounter) count() {
	sc.m.Lock()
	sc.c++
	sc.m.Unlock()
}

func playAndCount(sc_stayed *SafeCounter, sc_changed *SafeCounter) {
	g := game{}
	g.playGame()
	if g.score["stayed"] {
		sc_stayed.count()
	} else {
		sc_changed.count()
	}
}

func main() {
	n := 100
	if len(os.Args) < 2 {
	} else if m, err := strconv.Atoi(os.Args[1]); err == nil {
		n = m
	}
	sc_stayed := SafeCounter{}
	sc_changed := SafeCounter{}

	var wg sync.WaitGroup
	start := time.Now()
	for i:=0; i<n; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			playAndCount(&sc_stayed, &sc_changed)
		}()
	}
	wg.Wait()
	elapsed := time.Since(start)
	fmt.Printf("with goroutine...%v stayed: %d, changed: %d\n", elapsed, sc_stayed.c, sc_changed.c)
}
