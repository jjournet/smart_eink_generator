package main

import (
	"io/ioutil"
	"log"
	"os"

	"github.com/golang/freetype/truetype"
)

type preloaded_fonts struct {
	text   *truetype.Font
	icons  *truetype.Font
	custom *truetype.Font
}

func load_font(fonts *preloaded_fonts) {
	//fontFile, err := os.Open("fonts/Roboto-Regular.ttf")
	//fontFile, err := os.Open("fonts/IBMPlexSans-Medium.ttf")
	fontFile, err := os.Open("fonts/AvenirNextLTPro-Bold.ttf")

	if err != nil {
		log.Fatal(err)
	}
	defer fontFile.Close()
	fontBytes, err := ioutil.ReadAll(fontFile)
	if err != nil {
		log.Fatal(err)
	}
	// load font from ttf file
	font, err := truetype.Parse(fontBytes)
	if err != nil {
		log.Fatal(err)
	}
	fonts.text = font
	fontFile2, err := os.Open("fonts/weathericons-regular-webfont.ttf")
	if err != nil {
		log.Fatal(err)
	}
	defer fontFile2.Close()
	fontBytes2, err := ioutil.ReadAll(fontFile2)
	if err != nil {
		log.Fatal(err)
	}
	// load font from ttf file
	font2, err := truetype.Parse(fontBytes2)
	if err != nil {
		log.Fatal(err)
	}
	fonts.icons = font2

	// font 3 : icomoon.ttf
	fontFile3, err := os.Open("fonts/icomoon.ttf")
	if err != nil {
		log.Fatal(err)
	}
	defer fontFile3.Close()
	fontBytes3, err := ioutil.ReadAll(fontFile3)
	if err != nil {
		log.Fatal(err)
	}
	// load font from ttf file
	font3, err := truetype.Parse(fontBytes3)
	if err != nil {
		log.Fatal(err)
	}
	fonts.custom = font3
}
