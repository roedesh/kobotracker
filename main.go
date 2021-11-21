package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"time"

	"github.com/shermp/go-fbink-v2/gofbink"

	"cryptokobo/screener"
)

var (
	version string
)

func testPrints(fb *gofbink.FBInk, cfg *gofbink.FBInkConfig) {
	// Test the console like printing feature
	for i := 0; i < 10; i++ {
		s := fmt.Sprintf("Test line %d", i)
		fb.Println(s)
		fmt.Println(s)
		time.Sleep(500 * time.Millisecond)
	}
	// Lets test last line replacement next
	fb.PrintLastLn("This should update the last line!")
	time.Sleep(1 * time.Second)
	// And we finish with a nice progress bar :)
	for i := 0; i <= 100; i += 10 {
		fb.PrintProgressBar(uint8(i), cfg)
		fmt.Println("Progress bar @", i, "%")
		time.Sleep(500 * time.Millisecond)
	}
}

func main() {
	exec.Command("killall", "-s", "SIGKILL", "KoboMenu").Run()

	screen := screener.InitScreen()

	if _, err := os.Stat("./config.ini"); errors.Is(err, os.ErrNotExist) {
		// path/to/whatever does not exist
	}

	screen.ClearScreen()
	defer screen.Close()

	screen.Print("CryptoKobo")
	screen.Print(version)

	time.Sleep(5 * time.Second)

	return
}
