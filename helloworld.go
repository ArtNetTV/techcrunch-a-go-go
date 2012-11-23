package main

import (
       "fmt"
)

func double(x, y int) (p, q int) {
     p = x * 2
     q = y * 2
     return
}

func main() {
     a, b := double(1,2)
     fmt.Println(a,b)
}