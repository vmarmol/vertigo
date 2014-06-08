package main

import (
	"log"

	"github.com/vmarmol/vertigo/pulet"
)

func main() {
	id := "e0fdd7e222041ec15e41d596708f384d6084ef598f6aee709c7423b38b5f30fb"
	pul := pulet.NewPulet()
	importSpec := &pulet.ImportSpec{
		SourceHost: "107.178.222.32",
		SourcePort: 8080,
		SourceId:   id,
	}
	img, err := pul.Import(importSpec)
	if err != nil {
		log.Fatal(err)
	}
	err = pul.RunImage(img, nil, []string{
		"touch",
		"/target.txt",
	})
	if err != nil {
		log.Fatal(err)
	}
}
