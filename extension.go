package entgqlplus

import (
	"embed"
	"path"

	"entgo.io/ent/entc"
	"entgo.io/ent/entc/gen"
)

//go:embed templates
var assets embed.FS

type Extension struct {
	entc.DefaultExtension
	hooks  []gen.Hook
	config *Config
}

func (e *Extension) Generate(next gen.Generator) gen.Generator {
	return gen.GenerateFunc(func(g *gen.Graph) error {
		files := []File{}
		schemaDir := path.Join(path.Dir(e.config.GqlGenPath), path.Dir(e.config.GqlGen.Schema[0]))
		resolverDir := path.Join(path.Dir(e.config.GqlGenPath), e.config.GqlGen.Resolver.Dir)

		data := &TemplateData{
			Config:  e.config,
			Package: ejectPackage(g),
			Nodes:   ejectNodes(g),
		}

		files = append(files,
			File{
				Path:   path.Join(resolverDir, "schema.resolvers.go"),
				Buffer: parseTemplate("schema.resolvers.go.tmpl", data),
			},
			File{
				Path:   path.Join(resolverDir, "resolver.go"),
				Buffer: parseTemplate("resolver.go.tmpl", data),
			},

			File{
				Path:   path.Join(schemaDir, "types.graphqls"),
				Buffer: parseTemplate("types.graphqls.go.tmpl", nil),
			},
		)

		if e.config.Mutation {
			for i := range data.Nodes {
				data.Node = data.Nodes[i]
				files = append(files,
					File{
						Path:   path.Join(resolverDir, snake(data.Node)+".resolvers.go"),
						Buffer: parseTemplate("node.resolvers.go.tmpl", data),
					},
					File{
						Path:   path.Join(schemaDir, snake(data.Node)+".graphqls"),
						Buffer: parseTemplate("node.graphqls.go.tmpl", data),
					},
				)
			}
		}

		if e.config.Echo {
			files = append(files,
				File{
					Path:   path.Join(path.Dir(e.config.GqlGenPath), "routes/routes.go"),
					Buffer: parseTemplate("routes.go.tmpl", data),
				},
				File{
					Path:   path.Join(path.Dir(e.config.GqlGenPath), "handlers/handlers.go"),
					Buffer: parseTemplate("handlers.go.tmpl", data),
				},
				File{
					Path:   path.Join(path.Dir(e.config.GqlGenPath), "server.go"),
					Buffer: parseTemplate("server.go.tmpl", data),
				},
			)
		}

		if len(e.config.Database) > 0 {
			files = append(files,
				File{
					Path:   path.Join(path.Dir(e.config.GqlGenPath), "db/db.go"),
					Buffer: parseTemplate("db.go.tmpl", data),
				},
			)
		}

		if e.config.FileUpload {
			files = append(files,
				File{
					Path:   path.Join(resolverDir, "upload.resolvers.go"),
					Buffer: parseTemplate("upload.resolvers.go.tmpl", data),
				},
				File{
					Path:   path.Join(schemaDir, "upload.graphqls"),
					Buffer: parseTemplate("upload.graphqls.go.tmpl", nil),
				},
			)
		}

		writeFiles(files)

		return next.Generate(g)
	})
}

func (e *Extension) Hooks() []gen.Hook {
	return e.hooks
}

func NewExtension(opts ...ExtensionOption) *Extension {
	ex := &Extension{
		// Default Config
		config: &Config{
			FileUpload:   false,
			Subscription: false,
			Mutation:     false,
			Database:     SQLite,
			Echo:         false,
			GqlGenPath:   "../gqlgen.yaml",
			JWT:          false,
		},
		hooks: []gen.Hook{},
	}

	for i := range opts {
		opts[i](ex)
	}

	ex.config.GqlGen = readGqlGen(ex.config.GqlGenPath)
	ex.hooks = append(ex.hooks, ex.Generate)
	return ex
}
