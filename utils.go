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
	snake = gen.Funcs["snake"].(func(string) string)
)

func writeFiles(files []file) {
	for i := range files {
		writeFile(files[i])
	}
}

func writeFile(f file) {
	err := os.MkdirAll(path.Dir(f.Path), 0666)
	catch(err)
	os.WriteFile(f.Path, []byte(f.Buffer), 0666)
}

func ejectNodes(g *gen.Graph) []string {
	nodes := []string{}
	for i := range g.Nodes {
		nodes = append(nodes, g.Nodes[i].Name)
	}
	return nodes
}

func ejectPackage(g *gen.Graph) string {
	return strings.ReplaceAll(g.Package, "/ent", "")
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
	return out
}
