package entgqlplus

type (
	templateData struct {
		Package string
		Nodes   []string
		Node    string
		Config  *config
	}

	config struct {
		Database     Database
		Echo         bool
		JWT          bool
		Mutation     bool
		FileUpload   bool
		Subscription bool
		GqlGenPath   string
		GqlGen       gqlGen
	}

	file struct {
		Path   string
		Buffer string
	}

	gqlGen struct {
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
)
