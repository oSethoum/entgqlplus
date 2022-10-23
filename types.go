package entgqlplus

type TemplateData struct {
	Package string
	Nodes   []string
	Node    string
	Config  *Config
}

type Config struct {
	Database     Database
	Echo         bool
	JWT          bool
	Mutation     bool
	FileUpload   bool
	Subscription bool
	GqlGenPath   string
	GqlGen       GqlGen
}

type File struct {
	Path   string
	Buffer string
}

type GqlGen struct {
	Schema []string `yaml:"schema"`
	Exec   struct {
		FileName string `yaml:"filename"`
		Dir      string
		Package  string `yaml:"package"`
	} `yaml:"exec"`
	Model struct {
		FileName string `yaml:"filename"`
		Package  string `yaml:"package"`
	} `yaml:"model"`
	Resolver struct {
		Dir     string `yaml:"dir"`
		Package string `yaml:"package"`
	} `yaml:"resolver"`
}
