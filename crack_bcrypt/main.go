package main

import (
	"bufio"
	"flag"
	"log"
	"os"

	"golang.org/x/crypto/bcrypt"
)

var (
	//$2a$10$Zs3ZwsjV/nF.KujSUE.51uwtDrK6UVXcBpQrH84V8q3Opg1yddWLu"
	_stored_hash string = "PUT YOUR HASH HERE"
	_wordlist    string
	_nb_workers  int
	_channels    chan string
)

func init() {
	flag.StringVar(&_wordlist, "f", "", "Wordlist")
	flag.IntVar(&_nb_workers, "w", 10, "number of workers")
	flag.Parse()
}

func worker() {
	for p := range _channels {
		crack_bcrypt(p)
	}
}

func bcrypt_hash_string(pwd string) string {
	hash, err := bcrypt.GenerateFromPassword(
		[]byte(pwd),
		bcrypt.DefaultCost,
	)
	if err != nil {
		panic(err)
	}
	return string(hash)
}

func open_file(path string) *os.File {
	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	return file
}

func crack_bcrypt(password string) {
	hash := bcrypt_hash_string(password)
	log.Printf("password = %s | hash = %s\n", password, hash)
	if err := bcrypt.CompareHashAndPassword(
		[]byte(_stored_hash),
		[]byte(password),
	); err != nil {
		return
	} else {
		log.Printf("[+] Password is correct (bcrypt): %s\n", password)
		os.Exit(0)
	}
}

// bcrypt
func main() {
	// open file
	file := open_file(_wordlist)
	defer file.Close()

	if _nb_workers <= 0 {
		_nb_workers = 10
	}
	if _stored_hash == "" || _stored_hash == "PUT YOUR HASH HERE" {
		log.Println("[!] No hash provided! Please set the hash in line 14 of main.go")
		os.Exit(84)
	}

	// Create a channel to communicate with the workers
	_channels = make(chan string, _nb_workers)

	// Start the workers
	for i := 0; i < cap(_channels); i++ {
		go worker()
	}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		password := scanner.Text()
		_channels <- password
	}
	log.Printf("[-] Not password found")
	close(_channels)
}
