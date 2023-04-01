package main

import (
	"bytes"
	"embed"
	"log"
	"os"
	"path"
	"strings"
	"text/template"

	"entgo.io/ent/entc/gen"
)

type (
	templateData struct {
		Pkg string
	}

	file struct {
		Path   string
		Buffer string
	}
)

//go:embed templates
var assets embed.FS

func main() {
	// check teh go.mpd file if found get the package from it
	buff, err := os.ReadFile("go.mod")
	if err != nil {
		log.Fatalln("Cannot find the file go.mod")
	}

	pkg := strings.ReplaceAll(strings.Split(strings.Split(string(buff), "\n")[0], " ")[1], " ", "")
	if pkg == "" {
		log.Fatalln("Unable to read package from go.mod")
	}

	data := &templateData{
		Pkg: pkg,
	}
	files := []file{
		{
			Path:   "gqlgen.yml",
			Buffer: parseTemplate("gqlgen.go.tmpl", data),
		},
		{
			Path:   "ent/generate/entc.go",
			Buffer: parseTemplate("entc.go.tmpl", data),
		},
		{
			Path:   "ent/generate/generate.go",
			Buffer: parseTemplate("generate.go.tmpl", data),
		},
	}
	os.Mkdir("ent/schema", 0777)
	writeFiles(files)
}

func writeFiles(files []file) {
	for _, f := range files {
		err := os.MkdirAll(path.Dir(f.Path), 0777)
		catch(err)
		os.WriteFile(f.Path, []byte(f.Buffer), 0777)
	}
}

func parseTemplate(filename string, data *templateData) string {
	buffer, err := assets.ReadFile("templates/" + filename)
	catch(err)
	t, err := template.New(filename).Funcs(gen.Funcs).Parse(string(buffer))
	catch(err)
	out := bytes.Buffer{}
	err = t.Execute(&out, data)
	catch(err)
	return out.String()
}

func catch(err error) {
	if err != nil {
		log.Fatalln("entgqlplus:", err)
	}
}
