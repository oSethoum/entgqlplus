package entgqlplus

import (
	"embed"
	"os"
	"path"
	"strings"

	"entgo.io/ent/entc"
	"entgo.io/ent/entc/gen"
)

//go:embed templates
var assets embed.FS

type extension struct {
	entc.DefaultExtension
	hooks  []gen.Hook
	config *config
}

func (e *extension) generate(next gen.Generator) gen.Generator {
	return gen.GenerateFunc(func(g *gen.Graph) error {
		files := []file{}
		schemaDir := path.Join(path.Dir(e.config.GqlGenPath), path.Dir(e.config.GqlGen.Schema[0]))
		resolverDir := path.Join(path.Dir(e.config.GqlGenPath), e.config.GqlGen.Resolver.Dir)
		dbDir := path.Join(path.Dir(e.config.GqlGenPath), "db")
		authDir := path.Join(path.Dir(e.config.GqlGenPath), "auth")

		data := &templateData{
			Config: e.config,
		}

		data.parse(g)

		files = append(files,
			file{
				Path:   path.Join(resolverDir, "schema.resolvers.go"),
				Buffer: parseTemplate("schema.resolvers.go.tmpl", data),
			},
			file{
				Path:   path.Join(resolverDir, "resolver.go"),
				Buffer: parseTemplate("resolver.go.tmpl", data),
			},

			file{
				Path:   path.Join(schemaDir, "types.graphqls"),
				Buffer: parseTemplate("types.graphqls.go.tmpl", data),
			},
		)

		if e.config.Mutation != nil {
			if *e.config.Mutation {
				for i := range data.Nodes {
					data.Node = data.Nodes[i]
					files = append(files,
						file{
							Path:   path.Join(resolverDir, snake(data.Node.Name)+".resolvers.go"),
							Buffer: parseTemplate("node.resolvers.go.tmpl", data),
						},
						file{
							Path:   path.Join(schemaDir, snake(data.Node.Name)+".graphqls"),
							Buffer: parseTemplate("node.graphqls.go.tmpl", data),
						},
					)
				}
			} else {
				// delete the files
				for _, n := range data.Nodes {
					os.Remove(path.Join(schemaDir, snake(data.Node.Name)+".graphqls"))
					os.Remove(path.Join(resolverDir, snake(n.Name)+".resolvers.go"))
				}
			}
		}

		if e.config.Echo != nil {
			if *e.config.Echo {
				files = append(files,
					file{
						Path:   path.Join(path.Dir(e.config.GqlGenPath), "routes/routes.go"),
						Buffer: parseTemplate("routes.go.tmpl", data),
					},
					file{
						Path:   path.Join(path.Dir(e.config.GqlGenPath), "handlers/handlers.go"),
						Buffer: parseTemplate("handlers.go.tmpl", data),
					},
					file{
						Path:   path.Join(path.Dir(e.config.GqlGenPath), "server.go"),
						Buffer: parseTemplate("server.go.tmpl", data),
					},
				)
			} else {
				// delete the files
				os.Remove(path.Join(path.Dir(e.config.GqlGenPath), "routes/routes.go"))
				os.Remove(path.Join(path.Dir(e.config.GqlGenPath), "handlers/handlers.go"))
				os.Remove(path.Join(path.Dir(e.config.GqlGenPath), "server.go"))
			}
		}
		if len(e.config.Database) > 0 {
			files = append(files,
				file{
					Path:   path.Join(dbDir, "db.go"),
					Buffer: parseTemplate("db.go.tmpl", data),
				},
			)
		} else {
			// delete the files
			os.RemoveAll(dbDir)
		}

		if e.config.FileUpload != nil {
			if *e.config.FileUpload {
				files = append(files,
					file{
						Path:   path.Join(resolverDir, "upload.resolvers.go"),
						Buffer: parseTemplate("upload.resolvers.go.tmpl", data),
					},
					file{
						Path:   path.Join(schemaDir, "upload.graphqls"),
						Buffer: parseTemplate("upload.graphqls.go.tmpl", data),
					},
				)
			} else {
				// delete the files
				os.Remove(path.Join(resolverDir, "upload.resolvers.go"))
				os.Remove(path.Join(schemaDir, "upload.graphqls"))
			}
		}

		if e.config.Subscription != nil {
			if *e.config.Subscription && data.HasSubscription {
				files = append(files,
					file{
						Path:   path.Join(resolverDir, "notifiers.go"),
						Buffer: parseTemplate("notifiers.go.tmpl", data),
					},
					file{
						Path:   path.Join(resolverDir, "types.go"),
						Buffer: parseTemplate("types.go.tmpl", data),
					},
				)
			} else {
				// delete the files
				os.Remove(path.Join(resolverDir, "notifiers.go"))
				os.Remove(path.Join(resolverDir, "types.go"))
			}
		}

		if e.config.JWT != nil {
			if *e.config.JWT {
				files = append(files,
					file{
						Path:   path.Join(authDir, "login.go"),
						Buffer: parseTemplate("auth.login.go.tmpl", data),
					},
					file{
						Path:   path.Join(authDir, "middleware.go"),
						Buffer: parseTemplate("auth.middleware.go.tmpl", data),
					},
					file{
						Path:   path.Join(authDir, "types.go"),
						Buffer: parseTemplate("auth.types.go.tmpl", data),
					},
				)
			} else {
				// delete the files
				os.Remove(path.Join(authDir, "login.go"))
				os.Remove(path.Join(authDir, "middleware.go"))
				os.Remove(path.Join(authDir, "types.go"))
			}
		}

		if e.config.Privacy != nil {
			if *e.config.Privacy {
				files = append(files, file{
					Path:   path.Join(authDir, "privacy.go"),
					Buffer: parseTemplate("auth.privacy.go.tmpl", data),
				})
				for _, n := range g.Nodes {
					fpath := path.Join(path.Dir(e.config.GqlGenPath), "ent/schema/", lower(n.Name)+".go")

					data.Node = node{
						Name: n.Name,
					}
					addedLines := []string{
						"\t\"" + path.Join(data.Package, "/auth") + "\"",
						"\t\"" + path.Join(data.Package, "/ent/privacy") + "\"",
					}

					check := []string{
						path.Join(data.Package, "/auth"),
						path.Join(data.Package, "/ent/privacy"),
					}
					str := appendLines(fpath, addedLines, 4, beforeMode, check)
					privacybuffer := parseTemplate("node.privacy.go.tmpl", data)
					if !strings.Contains(str, "// Policy defines the privacy policy of the") {
						str += "\n" + privacybuffer
					}
					writeFile(file{
						Path:   fpath,
						Buffer: str,
					})
				}
			} else {
				os.Remove(path.Join(authDir, "privacy.go"))
				for _, n := range g.Nodes {
					fpath := path.Join(path.Dir(e.config.GqlGenPath), "ent/schema/", lower(n.Name)+".go")
					buffer, err := os.ReadFile(fpath)
					if err != nil {
						panic(err)
					}
					removedLines := []removeLine{
						{
							substr: "// Policy defines the privacy policy",
							end:    true,
						},
						{substr: path.Join(data.Package, "/auth")},
						{substr: path.Join(data.Package, "/ent/privacy")},
					}
					newBuffer := removeLines(string(buffer), removedLines)
					writeFile(file{
						Path:   fpath,
						Buffer: newBuffer,
					})
				}
			}
		}

		cleanFiles(resolverDir, schemaDir)
		writeFiles(files)

		return next.Generate(g)
	})
}

func (e *extension) Hooks() []gen.Hook {
	return e.hooks
}

func NewExtension(opts ...extensionOption) *extension {
	ex := &extension{
		// Default Config
		config: &config{
			FileUpload:   nil,
			Subscription: nil,
			Mutation:     nil,
			Database:     "",
			Echo:         nil,
			GqlGenPath:   "../gqlgen.yaml",
			JWT:          nil,
		},
		hooks: []gen.Hook{},
	}

	for i := range opts {
		opts[i](ex)
	}

	gen.Funcs["camel"] = func(s string) string { return camel(snake(s)) }
	gen.Funcs["lower"] = strings.ToLower

	ex.config.GqlGen = readGqlGen(ex.config.GqlGenPath)
	ex.hooks = append(ex.hooks, ex.generate)
	return ex
}
