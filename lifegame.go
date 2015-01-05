package main

import (
  "flag"
  "fmt"
  "math/rand"
  "os"
  "runtime"
  "strconv"
  "strings"
  "sync"
  "time"
)

type Lifegame struct {
  Width int
  Height int
  Field []int
  nextField []int
}

func NewLifegame(w int, h int, n int) *Lifegame {
  res := &Lifegame{w, h, make([]int, w*h), make([]int, w*h)}
  rand.Seed(time.Now().Unix())
  for i:=0; i<n; i++ {
    x := rand.Intn(w)
    y := rand.Intn(h)
    res.Set(x, y, 1)
  }
  return res
}

func (self *Lifegame) Get(x int, y int) int {
  return self.Field[x + y*self.Width]
}

func (self *Lifegame) Set(x int, y int, val int) {
  self.Field[x + y*self.Width] = val
}

func (self *Lifegame) Print() {
  for y:=0; y<self.Height; y++ {
    for x:=0; x<self.Width; x++ {
      var s string
      if self.Get(x,y) == 0 {
        s = "__"
      } else {
        s = "[]"
      }
      fmt.Printf("%s", s)
    }
    fmt.Printf("\n")
  }
}

func (self *Lifegame) CellNext(idx int) int {
  w := self.Width
  delta := [8]int{-w-1, -w, -w+1, -1, 1, w-1, w, w+1}

  numAlive := 0
  for _, d := range delta {
    i := idx + d
    if 0 <= i && i < len(self.Field) && self.Field[i] != 0 {
      numAlive++
    }
  }
  life := self.Field[idx]
  if (life==0 && numAlive==3) || (life==1 && (numAlive==2 || numAlive==3)) {
    return 1
  }
  return 0
}

func (self *Lifegame) Update(numRoutine int) {
  m := len(self.Field)/numRoutine
  var wg sync.WaitGroup
  for i:=0; i<len(self.Field); i+=m {
    wg.Add(1)
    go func(first int, last int) {
      for j:=first; j<last; j++ {
        self.nextField[j] = self.CellNext(j)
      }
      wg.Done()
    }(i, i+m)
  }
  wg.Wait()
  self.Field, self.nextField = self.nextField, self.Field
}

func parseSize(size string) (int, int) {
  sep := strings.Index(size, "x")
  w, errw := strconv.Atoi(size[:sep])
  h, errh := strconv.Atoi(size[sep+1:])
  if errw != nil || errh != nil {
    fmt.Printf("Invalid flag --size='%s'\n", size)
    os.Exit(1)
  }
  return w, h
}

func main() {
  size  := flag.String("size", "10x10", "Size of field. e.g. 10x10")
  steps := flag.Int("steps", 50, "Number of steps")
  procs := flag.Int("procs", 1, "Number of procs to run")
  flag.Parse()

  w, h := parseSize(*size)

  fmt.Printf("Procs: %d\n", *procs)
  fmt.Printf("Field: %dx%d\n", w, h)
  runtime.GOMAXPROCS(*procs)

  lifegame := NewLifegame(w, h, w*h/4)
  lifegame.Print()

  for t:=0; t<*steps; t++ {
    fmt.Println("\n", t)
    lifegame.Update(*procs)
    lifegame.Print()
    time.Sleep(time.Second/10)
  }
}
