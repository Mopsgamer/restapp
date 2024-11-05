package internal

import (
	"bufio"
	"bytes"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"slices"
	"strings"

	"github.com/gofiber/fiber/v3/log"
)

// Creates the .env file, if provided the '--init' option.
func InitProjectFiles() {

	optionInit := "--init"
	optionForce := "--force"

	path := ".env"
	isInit := slices.Contains(os.Args, optionInit)
	if !isInit {
		return
	}

	force := slices.Contains(os.Args, optionForce)
	_, errstat := os.Stat(path)
	exists := !os.IsNotExist(errstat)
	if exists && !force {
		log.Fatal("Failed to write " + path + " - already exists, use " + optionForce)
	}

	err := os.WriteFile(path, []byte(
		"DB_PASSWORD=\n"+
			"JWT_KEY=\n"+
			"DB_NAME=restapp\n"+
			"DB_USER=root\n"+
			"DB_HOST=localhost\n"+
			"DB_PORT=3306\n",
	), fs.ModeDevice)

	if err != nil {
		log.Fatal("Failed to write " + path)
	}

	log.Info("Writed " + path)

	log.Info("Executing deno build task...")
	ExecDeno("task", "build")
	log.Info("Deno tasks completed successfully.")
	os.Exit(0)
}

// Creates the instance. Checks if the deno command exists: exits with 1 if it doesn't.
// Uses default system's output.
func ExecDeno(arg ...string) *exec.Cmd {
	_, err := exec.LookPath("deno")
	if err != nil {
		fmt.Println("Deno is not installed or not available in PATH.")
		os.Exit(1)
	}

	cmd := exec.Command("deno", arg...)

	return cmd
}

func WaitForBundleWatch() {

	optionWatch := "--watch"
	optionBuild := "--build"

	isBuild := slices.Contains(os.Args, optionBuild)
	isWatch := slices.Contains(os.Args, optionWatch)

	if !isBuild && !isWatch {
		log.Info("You can use --build or --watch option to bundle js, css and assets before running server.")
		return
	}

	if !isBuild && !isWatch {
		log.Fatal("Use --build or --watch, if you want to bundle while running the server. You have used both.")
		return
	}

	var deno *exec.Cmd
	if isWatch {
		log.Info("Creating file listeners for bundling js, css and assets...")
		deno = ExecDeno("task", "build", "--watch")
	} else {
		log.Info("Bundling js, css and assets...")
		deno = ExecDeno("task", "build")
	}

	reader, err := deno.StdoutPipe()
	if err != nil {
		log.Fatal(err)
	}

	scanner := bufio.NewScanner(reader)
	var buffer bytes.Buffer
	deno.Start()

	for scanner.Scan() {
		line := scanner.Text()

		buffer.WriteString(line + "\n")

		if isWatch {
			// see ./web/build.ts file
			if strings.Contains(line, "watching") {
				break
			}

			continue
		}

		// see ./web/build.ts file
		if strings.Contains(line, "bundled") {
			break
		}
	}
	log.Info("Watching for file changes while starting the server...")
}
