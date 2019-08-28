package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path"
	"regexp"
	"strings"
	"time"
)

var running = map[string]*exec.Cmd{}
var dummyDir = "./tmp"

var dummyStorage = flag.String("dummies", "./.", "The location of where the dummies will be stored")
var scanWait = flag.Int64("wait", 60, "duration in seconds between scans for .exe processes")
var firejail = flag.Bool("use-firejail", false, "Set if this discord is running in a firejail")
var firejailName = flag.String("firejail-name", "discord", "the name of the firejail discord is running in.")
var firejailPathOverride = flag.String("firejail-path-override", "", "override the dummy storage passed to firejail 'firejail --join=discord override/something.dummy'")

var names = map[string]bool{}

func main() {
	flag.Parse()
	dummyFile, err := os.Open("./dummy")
	defer dummyFile.Close()
	dummyDir = *dummyStorage
	var cmd *exec.Cmd
	var ok bool
	if err == nil {

		for {
			pIDsAndName, dummies := getExePidNameListAndDummyMap()
			names = map[string]bool{}
			for _, pn := range pIDsAndName {
				pid, name := pn[1], path.Base(strings.Replace(pn[2], "\\", "/", -1))
				names[name] = true
				log.Printf("PID: %s, Name: %s\n", pid, name)

				if file, err := os.Open(path.Join(dummyDir, fmt.Sprint(name, ".dummy"))); os.IsNotExist(err) {
					file.Close()

					if file, err := os.Create(path.Join(dummyDir, fmt.Sprint(name, ".dummy"))); err == nil {
						_, err := dummyFile.Seek(0, 0)
						if err != nil {
							file.Close()
							continue
						}
						_, err = io.Copy(file, dummyFile)
						file.Close()

						if err != nil {
							log.Println("Error copying dummy for", name)

							continue
						}

						err = os.Chmod(file.Name(), 500)
						if err != nil {
							log.Println("Error chmoding dummy for", name)
							continue
						}
						log.Println("Successfully created dummy for", name)

					} else {
						file.Close()
						log.Println("Error creating dummy for", name)
						continue
					}

				}

				if _, ok = dummies[name]; !ok {

					if !(*firejail) {

						cmd = exec.Command(path.Join(dummyDir, fmt.Sprint(name, ".dummy")))
						cmd.Stdout = os.Stdout
						err = cmd.Start()
						if err != nil {
							log.Fatalln(err)
						}

					} else {

						dir := dummyDir
						if *firejailPathOverride != "" {
							dir = *firejailPathOverride
						}

						cmd = exec.Command("firejail", fmt.Sprintf("--join=%s", *firejailName), path.Join(dir, fmt.Sprint(name, ".dummy")))

						err = cmd.Start()
						if err != nil {
							log.Fatalln(err)
						}
					}

				}
				dummies[name] = ""

			}

			for _, pid := range dummies {
				if pid != "" {
					cmd = exec.Command("kill", pid)
					var out bytes.Buffer
					var errOut bytes.Buffer
					cmd.Stdout = &out
					cmd.Stderr = &errOut

					err := cmd.Run()
					if err != nil {
						log.Println("Error killing dummy", pid, err, errOut.String())
					}
				}
			}

			time.Sleep(time.Duration(*scanWait) * time.Second)
		}

	} else {

		log.Println("Error opening dummy", err)

	}

}

func getExePidNameListAndDummyMap() ([][]string, map[string]string) {

	log.Println("Checking for .exe and .dummy processes")

	cmd := exec.Command("ps", "aux")

	var out bytes.Buffer

	cmd.Stdout = &out

	err := cmd.Run()

	if err != nil {

		log.Fatal(err)

	}

	compile := regexp.MustCompile(`(?m)^\S+\s+(\d+)(?:\s+\S+){8}\s+(.+\.exe)[$\s*]`)

	pidsAndNames := compile.FindAllStringSubmatch(out.String(), -1)

	compile = regexp.MustCompile(`(?m)^\S+\s+(\d+)(?:\s+\S+){8}\s+(.+)\.dummy[$\s*]`)

	all := compile.FindAllStringSubmatch(out.String(), -1)

	dummies := make(map[string]string)
	for _, v := range all {
		split := strings.Split(v[2], "/")
		dummies[split[len(split)-1]] = v[1]
	}

	return pidsAndNames, dummies

}
