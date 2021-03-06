package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"time"

	"ajesler/playbulb"
)

var (
	effectFlag = flag.String("effect", "", "[flash|pulse|rainbow|fade|candle|solid]")
	colourFlag = flag.String("colour", "", "6 or 8 character hex code. If 8 characters, the first byte is the brightness, with 0 being off. Defaults to 00FF0000")
	speedFlag  = flag.Int("speed", 0, "a value from 0 - 255")
)

func effectMode() (byte, error) {
	switch *effectFlag {
	case "flash":
		return playbulb.FLASH, nil
	case "pulse":
		return playbulb.PULSE, nil
	case "rainbow":
		return playbulb.RAINBOW, nil
	case "fade":
		return playbulb.FADE, nil
	case "candle":
		return playbulb.CANDLE, nil
	case "solid":
		return playbulb.SOLID, nil
	case "":
		return playbulb.SOLID, nil
	default:
		return 0, errors.New(fmt.Sprintf("Unsupported effect '%s'", *effectFlag))
	}
}

func main() {
	flag.Usage = func() {
		fmt.Printf("Usage: %s [<options>] <candleID ...>\n", os.Args[0])
		fmt.Println("")
		flag.PrintDefaults()
	}

	flag.Parse()

	if *colourFlag == "" && *effectFlag == "" {
		flag.Usage()
		return
	}

	if *colourFlag == "" {
		*colourFlag = "00FF0000"
	}

	colour, err := playbulb.ColourFromHexString(*colourFlag)
	if err != nil {
		fmt.Println(err)
		return
	}

	eM, err := effectMode()
	if err != nil {
		fmt.Println(err)
		return
	}

	speed := byte(0)
	if *speedFlag > 255 || *speedFlag < 0 {
		fmt.Println("Speed must be between 0 and 255")
		return
	} else {
		speed = byte(*speedFlag)
	}

	candleIDs := flag.Args()
	if len(candleIDs) == 0 {
		fmt.Println("No candle IDs given")
		return
	}

	effect := playbulb.NewEffect(eM, colour, speed)

	candle := (playbulb.Candle)(nil)
	if len(candleIDs) == 1 {
		candle = playbulb.NewCandle(candleIDs[0])
	} else {
		candles := make([]playbulb.Candle, len(candleIDs))
		for i, cID := range candleIDs {
			c := playbulb.NewCandle(cID)
			candles[i] = c
		}
		candle = playbulb.NewCandleGroup(candles)
	}

	err = candle.Connect()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	candle.SetEffect(effect)

	// Required to give the SetEffect time to send before disconnection
	time.Sleep(time.Second)

	candle.Disconnect()
}
