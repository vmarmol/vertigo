package main

import (
	"flag"
	"log"

	"github.com/vmarmol/vertigo/pulet"
)

var argContainerId = flag.String("id", "13ccf48b7cec704c1cc441e87a4c66102a1ead1fa3dbcec6ac255c53e802b631", "container id")
var argCommand = flag.String("cmd", "/bin/true", "command")

func main() {
	id := *argContainerId
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
	/*
		err = pul.RunImage(img, []string{
			"-d",
			"-p",
			"3306:3306",
			"-e",
			"MYSQL_PASS=\"mypass\"",
		}, []string{
			"/run.sh",
		})
	*/
	err = pul.RunImage(img, nil, []string{"/bin/true"})
	if err != nil {
		log.Fatal(err)
	}
}
