package main

import (
	"flag"
	"fmt"
	"golang.org/x/exp/slices"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

const filename = "./urlsBrands.csv"

var urls = map[string]string{
	"ORLEN":     "https://platinum-oil.ru/t_import.php?id=9",
	"MEGUIN":    "https://meguin.su/t_import.php?id=3",
	"RUSEFF":    "https://ruseff-auto.ru/t_import.php?id=15",
	"PETROFER":  "https://petrofer.ru/t_import.php?id=2",
	"LIQUIMOLY": "https://liquimoly.ru/t_import.php?id=15",
	"REINWELL":  "https://reinwell.ru/t_import.php?id=1",
	"LUBEX":     "https://lubex-oil.ru/t_import.php?id=1",
	"OPET":      "https://opet.ru/t_import.php?id=1",
}

var separator = ";"

func main() {

	var brands, separatorArgs string

	flag.StringVar(&brands, "brands", "", "add string of brands with \",\" as a separator. "+
		"=ALL for all brands")
	flag.StringVar(&brands, "b", "", "add string of brands with \",\" as a separator (shorthand). "+
		"=ALL for all brands")
	flag.StringVar(&separatorArgs, "s", "", "separator, \";\" by default.")
	flag.Parse()

	if separatorArgs != "" {
		separator = separatorArgs
	}

	if brands == "" {
		err := fmt.Errorf("you have to add at least one brand")
		fmt.Println(err.Error())
		os.Exit(0)
	}

	var allBrands = false

	if brands == "ALL" {
		allBrands = true
	}

	colBrands := strings.Split(brands, ",")
	if !allBrands {
		for i := 0; i < len(colBrands); i++ {
			colBrands[i] = strings.ToUpper(colBrands[i])
		}

		for i := 0; i < len(colBrands); i++ {
			if _, ok := urls[colBrands[i]]; !ok {
				err := fmt.Errorf("Brand %s not found", colBrands[i])
				fmt.Println(err.Error())
				os.Exit(0)
			}
		}
	}

	os.Remove(filename)

	_, err := os.Create(filename)
	if err != nil {
		log.Fatal(err.Error())
	}

	f, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		panic(err)
	}

	defer f.Close()

	for k, v := range urls {
		if !allBrands {
			if !slices.Contains(colBrands, k) {
				continue
			}
		}

		response, err := http.Get(v)
		if err != nil {
			log.Fatal(err)
		}

		defer response.Body.Close()

		responseData, err := io.ReadAll(response.Body)
		if err != nil {
			log.Fatal(err)
		}

		var responseString = string(responseData)
		responseStrings := strings.Split(responseString, "\n")
		for i := 0; i < len(responseStrings); i++ {
			if responseStrings[i] != "" {
				oneString := strings.Split(responseStrings[i], ";")
				s := oneString[0] + separator + oneString[1] + separator + oneString[2] + "\n"
				if _, err := f.WriteString(s); err != nil {
					panic(err)
				}
			}
		}
	}
}
