package main

import (
	"encoding/xml"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/j6s/lightroom2aftershot/lib"
)

func main() {
	flag.Parse()
	if flag.NArg() != 1 {
		log.Printf("[ERROR] Must specify exactly 1 file to convert. %d specified", flag.NArg())
		log.Printf("[ERROR] Usage: lightroom2aftershot lightroompreset.xml > aftershotpreset.xml")
		os.Exit(1)
	}

	fileContents, err := ioutil.ReadFile(flag.Arg(0))
	if err != nil {
		log.Printf("Error while reading file")
		log.Fatal(err)
	}

	preset := lib.NewLightroomPreset()
	err = xml.Unmarshal(fileContents, &preset)
	if err != nil {
		log.Printf("Error while unmarshalling")
		log.Fatal(err)
	}

	aftershot := lib.NewAftershotPresetFromLightroom(preset)

	xml, err := xml.MarshalIndent(aftershot, "", "    ")
	if err != nil {
		log.Fatal(err)
	}

	// Just a rough estimate of indentation to have the options layed out nicer in the file
	prettyXml := strings.Replace(string(xml), "bopt:", "\n                                        bopt:", -1)

	fmt.Printf("%s", prettyXml)
}
