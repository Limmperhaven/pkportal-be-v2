package main

import (
	"flag"
	"log"
	"os"
	"strings"
)

func main() {
	s1 := flag.Bool("step1", false, "hide custom structs from sqlboiler")
	s2 := flag.Bool("step2", false, "return custom structs from tmp")
	flag.Parse()

	switch {
	case *s1:
		step1()
	case *s2:
		step2()
	default:
		log.Fatal("invalid flags")
	}
}

func step1() {
	dir, err := os.ReadDir("internal/models/tpportal")
	if err != nil {
		log.Fatalf("error reading dir: %s", err.Error())
	}

	err = os.Mkdir("internal/models/tmp", 0777)
	if err != nil {
		log.Fatalf("error creating tmp dir: %s", err.Error())
	}

	for _, file := range dir {
		if file.IsDir() {
			log.Fatal("custom should not contain directories")
		}
		if !strings.Contains(file.Name(), "app_") {
			continue
		}

		bytes, err := os.ReadFile("internal/models/tpportal/" + file.Name())
		if err != nil {
			log.Fatalf("error reading file: %s", err.Error())
		}
		err = os.WriteFile("internal/models/tmp/"+file.Name(), bytes, 0777)
		if err != nil {
			log.Fatalf("error writing file: %s", err.Error())
		}
	}
}

func step2() {
	dir, err := os.ReadDir("internal/models/tmp")
	if err != nil {
		log.Fatalf("error reading dir: %s", err.Error())
	}

	for _, file := range dir {
		if file.IsDir() {
			log.Fatal("custom should not contain directories")
		}
		if !strings.Contains(file.Name(), "app_") {
			continue
		}

		bytes, err := os.ReadFile("internal/models/tmp/" + file.Name())
		if err != nil {
			log.Fatalf("error reading file: %s", err.Error())
		}
		err = os.WriteFile("internal/models/tpportal/"+file.Name(), bytes, 0644)
		if err != nil {
			log.Fatalf("error writing file: %s", err.Error())
		}
	}

	err = os.RemoveAll("internal/models/tmp")
	if err != nil {
		log.Fatalf("error removing tmp directory: %s", err.Error())
	}
}
