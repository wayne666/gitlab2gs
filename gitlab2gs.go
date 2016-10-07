package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/bugagazavr/go-gitlab-client"
	"github.com/gogits/go-gogs-client"
	"io/ioutil"
	"os"
)

var usage = `Usage: gitlab2gs [options...]

Options:
	  -config Specify your config file
`

type Config struct {
	GitlabHost     string   `json:"gitlabHost"`
	GitlabApiPath  string   `json:"gitlabApiPath"`
	GitlabToken    string   `json:"gitlabToken"`
	GitlabUser     string   `json:"gitlabUser"`
	GitlabPassword string   `json:"gitlabPassword"`
	GitlabProjects []string `json:"gitlabProjects"`
	GogsUrl        string   `json:"gogsUrl"`
	GogsToken      string   `json:"gogsToken"`
	GogsApiPath    string   `json:"gogsApiPath"`
}

var configFile = flag.String("config", "./config.json", "json config file")

var userMap = make(map[string]*gogs.Organization)
var gitlab *gogitlab.Gitlab
var gs *gogs.Client

func main() {
	flag.Usage = func() {
		fmt.Fprint(os.Stderr, usage)
	}
	flag.Parse()

	if *configFile == "" {
		usageAndExit("Please specify your json config file")
	}

	file, err := ioutil.ReadFile(*configFile)
	checkErr("Config file error", err)

	var config Config
	json.Unmarshal(file, &config)

	//judge config file args is legal or not
	isLegalConfigArg(&config)

	gitlab = gogitlab.NewGitlab(config.GitlabHost, config.GitlabApiPath, config.GitlabToken)
	gs = gogs.NewClient(config.GogsUrl, config.GogsToken)

	migrateProjects := getProjects(config.GitlabProjects)
	for _, project := range migrateProjects {
		doMigrate(project, &config)
	}
}

func isLegalConfigArg(config *Config) {
	if config.GogsUrl == "" || config.GogsToken == "" || config.GitlabHost == "" ||
		config.GitlabToken == "" || config.GogsApiPath == "" || config.GitlabApiPath == "" {
		usageAndExit("json config file field error, please check your config file")
	}

}

// get all migrate projects
func getProjects(projects []string) []*gogitlab.Project {
	if len(projects) != 0 {
		var projectss []*gogitlab.Project
		for _, projectID := range projects {
			project, err := gitlab.Project(projectID)
			checkErr("Get project from ID error", err)
			projectss = append(projectss, project)
		}
		return projectss
	}

	allProjects, err := gitlab.AllProjects()
	checkErr("Get gitlab projects error", err)

	return allProjects
}

func getGogsUID(name string) int {
	org, ok := userMap[name]
	if ok {
		return int(org.ID)
	}
	orgUserName, err := gs.GetOrg(name)
	if err == nil {
		return int(orgUserName.ID)
	}

	createOrg := gogs.CreateOrgOption{
		UserName: name,
	}

	org, err = gs.AdminCreateOrg("root", createOrg)
	checkErr("Failed to create org", err)
	userMap[name] = org
	return int(org.ID)
}

func doMigrate(project *gogitlab.Project, config *Config) {
	t, err := gs.GetRepo(project.Namespace.Name, project.Name)
	fmt.Printf("t: %+v\n", t)

	if err == nil {
		fmt.Printf("%s # %s already in your gogs\n", project.Namespace.Name, project.Name)
	} else {
		orgID := getGogsUID(project.Namespace.Name)

		// api is reserved
		var name string = ""
		if project.Name == "api" {
			name = "myapi"
		} else {
			name = project.Name
		}

		fmt.Printf("%s # %s migrating to gogs # %s #...\n",
			project.Namespace.Name, project.Name, name)

		opts := gogs.MigrateRepoOption{
			CloneAddr:    project.HttpRepoUrl,
			AuthUsername: config.GitlabUser,
			AuthPassword: config.GitlabPassword,
			UID:          orgID,
			RepoName:     name,
			Private:      !project.Public,
			Description:  project.Description,
		}

		_, err := gs.MigrateRepo(opts)
		checkErr("", err)
	}
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

func checkErr(errMessage string, err error) {
	if err != nil {
		if errMessage != "" {
			fmt.Fprintf(os.Stderr, errMessage)
			fmt.Fprint(os.Stderr, "\n\n")
		}

		fmt.Fprintf(os.Stderr, "%+v\n", err)
		os.Exit(1)
	}
}
