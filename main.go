package main

import (
    "fmt"
    "transcendance/stuff"   // replace 'your_module' with your actual module name
)

func main() {
    s := stuff.NewStuff("Gadget")
    s.SayHello()
    fmt.Println("Name:", s.GetName())
}