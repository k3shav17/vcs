package main

import (
	"bytes"
	"compress/zlib"
	"crypto/sha1"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"text/tabwriter"
)

const (
	hashLength       int    = 40
	catFileCmdLength int    = 4
	parentDir        string = ".vcs"
)

func main() {

	args := os.Args
	command(args[1], args)
	// later try to create or replicate the easy functionalities of git in this.
	// use cobra module to get the args from command line

}

func command(arg string, args []string) {

	switch arg {
	case "init":
		initVcs()
	case "add":
		add(args)
	case "status":
		status()
	case "cat-file":
		commandCatFile(args)
	case "hash-object":
		hash, _ := hashContent(args)
		fmt.Println(hash)
	default:
		help()
	}
}

func status() {
	fmt.Println("status of the vcs work tree")
}

func initVcs() {
	for _, dir := range []string{parentDir, ".vcs/objects", ".vcs/objects/pack", ".vcs/objects/info", ".vcs/refs"} {
		err := os.Mkdir(dir, 0755)
		check(err)
	}

	contentsOfHeadFile := []byte("ref: refs/heads/main\n")
	if err := os.WriteFile(".vcs/HEAD", contentsOfHeadFile, 0644); err != nil {
		fmt.Fprintf(os.Stderr, "Error writing files -> %s\n", err)
	}
	fmt.Println("Initialized vcs dir")

}

// this is equivalent to git's hash-object function
// due to lack of creativity using same flag name in here.
func hashContent(args []string) (string, []byte) {

	contentsOfFile, err := os.ReadFile(args[2])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to add file -> %s\n", err)
		os.Exit(1)
	}

	// creating a blob to store with a header which contains length of the content blob
	// then append the file content to the header
	header := fmt.Sprintf("blob %d\x00", len(contentsOfFile))
	blobStore := append([]byte(header), contentsOfFile...)

	// create hash for the blob store which is combination of header and file content
	h := sha1.New()
	h.Write(blobStore)
	return fmt.Sprintf("%x", h.Sum(nil)), blobStore

}

// whenever we add a file to git it will create a hash from the content of the file
// and stores it in the objects folder with dir/filename
// dir -> first two chars of the hash and filename -> rest of the hash
func add(args []string) {
	fmt.Println(args[2:])

	for _, val := range args[2:] {
		if _, err := os.Stat(val); err != nil {
			fmt.Printf("%s Not a valid dir or file\n", val)
		}
	}

	// before hashing the file with args[2] which the file path
	// read the contents of the file and then create a hash
	hashBs, blobStore := hashContent(args)

	fullDirPath := filepath.Join(parentDir, "objects", hashBs[:2])
	fullFilePath := filepath.Join(fullDirPath, hashBs[2:])

	if err := os.MkdirAll(fullDirPath, 0755); err != nil {
		fmt.Fprintf(os.Stderr, "Error creating object dir: %s\n", err)
		os.Exit(1)
	}

	// compress and store the blob with zlib compressor
	var compressed bytes.Buffer
	w := zlib.NewWriter(&compressed)
	w.Write(blobStore)
	w.Close()

	if err := os.WriteFile(fullFilePath, compressed.Bytes(), 0644); err != nil {
		fmt.Fprintf(os.Stderr, "Error writing object: %s\n", err)
		os.Exit(1)
	}

	// same happens when we do the commit as well
	// but the hash generated is different than the add command

}

// this is what it spits out when cat-file is done for the hash generated after commit command
// tree fcb545d5746547a597811b7441ed8eba307be1ff
// author Keshava Kommaraju <keshavrao0489@gmail.com> 1742390280 +0530
// committer Keshava Kommaraju <keshavrao0489@gmail.com> 1742390280 +0530

// when created a new branch a new file is created in the .git/refs/ folder with branch name.
// this will be used to check where the reference is right at.
// as there will be a hash created for each commit, so cherry-pick is picking the changes by decoding the hash and applying them.
// Things to following, when hashing will the contents of the file be hashed or only the metadata.
// If only the metadata then how to get the content back.

func commandCatFile(args []string) {

	if len(args) != catFileCmdLength {
		fmt.Fprint(os.Stderr, "usage: vcs cat-file -p <object-hash>\n")
	}

	flag := args[2]
	hash := args[3]

	if flag != "-p" && len(hash) != hashLength {
		fmt.Fprintf(os.Stderr, "usage: vcs cat-file -p <object-hash>\n")
		os.Exit(1)
	}

	dirName := hash[:2]
	fileName := hash[2:]
	filePath := filepath.Join(parentDir, "objects", dirName, fileName)

	fmt.Printf("file path %s", filePath)

	file, err := os.Open(filePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to read file -> %s\n", err)
		os.Exit(1)
	}
	defer file.Close()

	zr, err := zlib.NewReader(file)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error decompressing object: %s\n", err)
		os.Exit(1)
	}
	defer zr.Close()

	contents, err := io.ReadAll(zr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading decompressed content: %s\n", err)
		os.Exit(1)
	}

	// Split header and content
	nullIndex := bytes.IndexByte(contents, 0)
	if nullIndex == -1 {
		fmt.Fprintln(os.Stderr, "Malformed object")
		os.Exit(1)
	}

	content := contents[nullIndex+1:]
	fmt.Print(string(content))

}

func help() {
	fmt.Println("Usage:")
	fmt.Println("  vcs <command> [options]")

	fmt.Println("Available Commands:")
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "init\tInitialize a new VCS repository")
	fmt.Fprintln(w, "add\tAdd file contents to the index")
	fmt.Fprintln(w, "cat-file\tProvide content or type of repository objects")
	fmt.Fprintln(w, "status\tShow the working tree status")
	fmt.Fprintln(w, "commit\tRecord changes to the repository")
	fmt.Fprintln(w, "ignore\tSpecify intentionally untracked files")
	w.Flush()

	fmt.Println("\nUse \"vcs <command> --help\" for more information about a command.")
}

func check(e error) {
	if e != nil {
		if !os.IsExist(e) {
			fmt.Fprintf(os.Stderr, "Error creating directory for vcs -> %s\n", e)
		}
	}
}
