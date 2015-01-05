package main

import (
  "fmt"
  "math/rand"
  "time"
  "sync"
  "runtime"
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
        s = ".."
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

  num_alive := 0
  for _, d := range delta {
    i := idx + d
    if 0 <= i && i < len(self.Field) && self.Field[i] != 0 {
      num_alive++
    }
  }
  life := self.Field[idx]
  if (life==0 && num_alive==3) || (life==1 && (num_alive==2 || num_alive==3)) {
    return 1
  }
  return 0
}

func (self *Lifegame) Update(num_routine int) {
  m := len(self.Field)/num_routine
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

func main() {
  fmt.Println("CPU:", runtime.NumCPU())
  numRoutine:= runtime.NumCPU()
  runtime.GOMAXPROCS(numRoutine)

  w, h := 10, 10
  lifegame := NewLifegame(w, h, w*h/4)
  lifegame.Print()

  tend := 50
  for t:=0; t<tend; t++ {
    fmt.Println(t)
    lifegame.Update(numRoutine)
    lifegame.Print()
    fmt.Println("\n")
    time.Sleep(time.Second/10)
  }
}
