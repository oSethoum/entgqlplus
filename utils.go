package entgqlplus

import (
	"log"
	"os"
	"path"
	"strings"

	"entgo.io/ent/entc/gen"
	"gopkg.in/yaml.v3"
)

var (
	lower = strings.ToLower
	camel = gen.Funcs["camel"].(func(string) string)
	snake = gen.Funcs["snake"].(func(string) string)

	beforeMode appendMode = "before"
	afterMode  appendMode = "after"
)

func writeFiles(files []file) {
	for i := range files {
		writeFile(files[i])
	}
}

func writeFile(f file) {
	err := os.MkdirAll(path.Dir(f.Path), 0777)
	catch(err)
	os.WriteFile(f.Path, []byte(f.Buffer), 0777)
}

func readFile(filePath string) string {
	buffer, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatalln(err)
	}
	return string(buffer)
}

func catch(err error) {
	if err != nil {
		log.Fatalln("entgqlplus:", err)
	}
}

func readGqlGen(fpath string) gqlGen {
	buffer, err := os.ReadFile(fpath)
	catch(err)
	out := gqlGen{}
	err = yaml.Unmarshal(buffer, &out)
	catch(err)
	out.Exec.Dir = path.Dir(out.Exec.FileName)
	out.Model.Dir = path.Dir(out.Model.FileName)
	return out
}

func inArray[T string | int | uint](array []T, value T) bool {
	for _, v := range array {
		if v == value {
			return true
		}
	}
	return false
}

func cleanFiles(resolverDir, schemaDir string) {
	os.RemoveAll(resolverDir)
	os.RemoveAll(schemaDir)
}

func in(lines []string, line string) bool {
	for _, l := range lines {
		if strings.Contains(l, line) {
			return true
		}
	}
	return false
}

func appendLines(filePath string, addedLines []string, pos int, mode appendMode, check []string) string {
	if mode == beforeMode {
		pos -= 2
	} else if mode == afterMode {
		pos -= 1
	}

	lines := strings.Split(readFile(filePath), "\n")

	newLines := []string{}
	for i, l := range lines {
		newLines = append(newLines, l)
		if i == pos {
			for ni, nl := range addedLines {
				if !in(lines, check[ni]) {
					newLines = append(newLines, nl)
				}
			}
		}
	}
	return strings.Join(newLines, "\n")
}

type removeLine struct {
	substr string
	end    bool
}

func removeLines(buffer string, rlines []removeLine) string {
	lines := strings.Split(buffer, "\n")
	newLines := []string{}

	for _, rl := range rlines {
		inRange := false
		for _, l := range lines {
			if strings.Contains(l, rl.substr) {
				inRange = true
				continue
			}
			if !inRange || !rl.end {
				newLines = append(newLines, l)
			}
		}
		lines = newLines
		newLines = []string{}
	}

	return strings.Join(lines, "\n")
}
