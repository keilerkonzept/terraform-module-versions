package main

import (
	"log"
	"os"

	"github.com/sgreben/versions/pkg/semver"
)

func laterCmd(later, earlier string, fail bool) {
	vLater, err := semver.Parse(later)
	if err != nil {
		log.Println(err)
		exit.NonzeroBecause = append(exit.NonzeroBecause, err.Error())
		return
	}
	vEarlier, err := semver.Parse(earlier)
	if err != nil {
		log.Println(err)
		exit.NonzeroBecause = append(exit.NonzeroBecause, err.Error())
		return
	}
	result := vEarlier.LessThan(vLater)
	if fail && !result {
		exit.NonzeroBecause = append(exit.NonzeroBecause, "comparison result is 'false'")
	}
	err = jsonEncode(result, os.Stdout)
	if err != nil {
		log.Println(err)
		exit.NonzeroBecause = append(exit.NonzeroBecause, err.Error())
	}
}
