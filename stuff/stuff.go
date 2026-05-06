package stuff

import (
	"fmt"
)

// Stuff is like a class (exported because capital S)
type Stuff struct {
    name string   // private field (lowercase)
}

// Constructor – exported function (capital N)
func NewStuff(name string) *Stuff {
    return &Stuff{name: name}
}

// Methods
func (s *Stuff) GetName() string {
    return s.name
}

func (s *Stuff) SetName(newName string) {
    s.name = newName
}

func (s *Stuff) SayHello() {
    fmt.Println("Hello, I'm", s.name)  // note: need import "fmt" above
}