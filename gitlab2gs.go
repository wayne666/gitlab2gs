package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/bugagazavr/go-gitlab-client"
	//"github.com/gogits/go-gogs-client"
	"io/ioutil"
	"os"
)

var usage = `Usage: gitlab2gs [options...]

Options:
	  -config Specific your config file
`

type Config struct {
	GitlabHost     string   `json:"gitlabHost"`
	GitlabApiPath  string   `json:"gitlabApiPath"`
	GitlabToken    string   `json:"gitlabToken"`
	GitlabProjects []string `json:"gitlabProjects"`
	GogsUrl        string   `json:"gogsUrl"`
	GogsToken      string   `json:"gogsToken"`
	GogsApiPath    string   `json:"gogsApiPath"`
}

var userMap = make(map[string]*gogs.Organization)
var gitlab = gogitlab.NewGitlab()

func main() {
	flag.Usage = func() {
		fmt.Fprint(os.Stderr, usage)
	}
	flag.Parse()

	if *configFile == "" {
		usageAndExit("Please specific your json config file")
	}

	file, err := ioutil.ReadFile(configFile)
	if err != nil {
		fmt.Printf("Config file error: %v\n", e)
		os.Exit(1)
	}

	var config Config
	json.Unmarshal(file, &config)
	fmt.Printf("Results: %+v\n", config)

	if config.GogsUrl == "" || config.GogsToken == "" || config.GitlabHost == "" ||
		config.GitlabToken == "" || config.GogsApiPath == "" || config.GitlabApiPath == "" {
		usageAndExit("test")
	}

	gitlab = gogitlab.NewGitlab(config.Host, config.ApiPath, config.Token)
	fmt.Printf("%+v\n", gitlab)

	migrateProjects := getProjects(config.GitlabProjects)
	for _, project := range migrateProjects {
		fmt.Printf("project %+v\n", project)
		fmt.Printf("org name: %s\n", project.Namespace.Name)
	}

	//org, err := gc.GetOrg(name)
}

// migrate range projects to this function
func getProjects(projects []string) []*gogitlab.Project {
	if projects != "" {
		return nil
	}

	projects, err := gitlab.AllProjects()
	if err != nil {
		fmt.Printf("gitlab gets projects err %v\n", err)
		os.Exit(1)
	}

	return projects
}

func userVerify(name string) *gogs.Organization {

}

func doMigrate(orgArgs *gogs.MigrateRepoOption) {

}

func usageAndExit(message string) {
	if message != "" {
		fmt.Fprintf(os.Stderr, message)
		fmt.Fprintf(os.Stderr, "\n\n")
	}

	flag.Usage()
	fmt.Fprintf(os.Stderr, "\n")
	os.Exit(1)
}
