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
	camel = gen.Funcs["camel"].(func(string) string)
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

func ejectNodes(g *gen.Graph) []node {
	nodes := []node{}
	for i := range g.Nodes {
		a := &annotation{}
		a.decode(g.Nodes[i].Annotations[entgqlSchemaKey])
		n := node{Name: g.Nodes[i].Name}
		if a.SchemaOptions != nil {
			n.Subscription = inArray(a.SchemaOptions, Subscription)
		}
		nodes = append(nodes, n)
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
