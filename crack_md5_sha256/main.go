package main

import (
	"bufio"
	"crypto/md5"
	"crypto/sha256"
	"flag"
	"fmt"
	"log"
	"os"
)

var (
	_sha256hash string // "95a5e1547df7dabdd4781b6c9e55f3377c15d08884b11738c2727dbd887d4ced"
	_md5hash    string // "77f62e3524cdd83d698d51fa24fdff4f"
	_wordlist   string
	_nb_workers int
	_channels   chan string
)

func init() {
	flag.StringVar(&_wordlist, "f", "wordlist.txt", "wordlist is the name of the file")
	flag.StringVar(&_sha256hash, "sha256", "", "sha256 hash")
	flag.StringVar(&_md5hash, "md5", "", "md5 hash")
	flag.IntVar(&_nb_workers, "w", 10, "number of workers")
	flag.Parse()
}

func crack_md_and_sha_256(password string) {
	if _md5hash != "" {
		check_password_md5(password)
	}
	if _sha256hash != "" {
		check_password_sha_256(password)
	}
}

func worker() {
	for p := range _channels {
		crack_md_and_sha_256(p)
	}
}

func open_file(path string) *os.File {
	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	return file
}

func check_password_md5(password string) {
	hash := fmt.Sprintf("%x", md5.Sum([]byte(password)))
	if hash == _md5hash {
		fmt.Printf("[+] Password found (MD5): %s\n", password)
		os.Exit(0)
	}
}

func check_password_sha_256(password string) {
	hash := fmt.Sprintf("%x", sha256.Sum256([]byte(password)))
	if hash == _sha256hash {
		fmt.Printf("[+] Password found (SHA-256): %s\n", password)
		os.Exit(0)
	}
}

func main() {
	// open file
	file := open_file(_wordlist)
	defer file.Close()

	// check
	if _nb_workers <= 0 {
		_nb_workers = 10
	}
	if _sha256hash == "" && _md5hash == "" {
		log.Println("[!] No hash provided")
		os.Exit(84)
	}

	// Create a channel to communicate with the workers
	_channels = make(chan string, _nb_workers)

	// Start the workers
	for i := 0; i < cap(_channels); i++ {
		go worker()
	}

	// loop wordlist
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		password := scanner.Text()
		_channels <- password
		if err := scanner.Err(); err != nil {
			log.Fatalln(err)
		}
	}
	fmt.Println("[!] No Password found")
}
