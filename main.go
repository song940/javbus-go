package main

import (
	"fmt"

	"github.com/song940/javbus-go/javbus"
)

func main() {
	client := javbus.New(&javbus.Config{
		Token: "99afDkPHVFe%2FEo8acl0yEpSD3GEjvFRCqtDFA1z3Ns8M6aX7BIl3qXiuJw",
	})
	detail, err := client.GetDetail("IPX-404")
	if err != nil {
		panic(err)
	}
	fmt.Println(detail.Title)
	fmt.Println(detail.Cover)
	fmt.Println(detail.Thumb)
	fmt.Println(detail.Director)
	fmt.Println(detail.Date)
	fmt.Println(detail.Duration)
	fmt.Println(detail.Studio)
	fmt.Println(detail.Label)
	fmt.Println(detail.Genres)
	fmt.Println(detail.Stars)
}
