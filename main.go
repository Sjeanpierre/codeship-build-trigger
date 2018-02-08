package main

import (
	"context"
	"fmt"
	codeship "github.com/codeship/codeship-go"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
)

type buildConfig struct {
	OrgName  string    `yaml:"org_name"`
	Projects []project `yaml:"projects"`
}

type project struct {
	Name   string
	Branch string
	UUID   string
}

func readConfig() (m buildConfig) {
	data, err := ioutil.ReadFile("build_trigger.yml")
	if err != nil {
		log.Fatal("could not read build_trigger.yml in current directory")
	}
	err = yaml.Unmarshal(data, &m)
	if err != nil {
		log.Fatal("could not unmarshall config file")
	}
	return
}

func findLatestSHA(p project, org *codeship.Organization) (string, error) {
	opt := codeship.PerPage(50)
	ctx := context.Background()
	builds, _, err := org.ListBuilds(ctx, p.UUID, opt)
	if err != nil {
		log.Fatalf("Could not list builds %s", err)
	}
	for _, build := range builds.Builds {
		if build.Ref == p.Branch && build.Status == "success" {
			return build.CommitSha, nil
		}
	}
	return "", fmt.Errorf("No prior builds to pull from for %s on branch %s", p.Name, p.Branch)
}

func trigger_build(org *codeship.Organization, p project) {
	ctx := context.Background()
	//list past builds, trigger last build
	latestBuildSha, err := findLatestSHA(p, org)
	if err != nil {
		log.Printf("Encountered error triggering builf for repo %s, error %s", p.Name, err.Error())
		return
	}
	success, resp, err := org.CreateBuild(ctx, p.UUID, p.Branch, latestBuildSha)
	if err != nil {
		log.Fatalf("Could not trigger build for %s\n response details:\n %s", p.Name, resp)
	}
	if success == true {
		log.Printf("Build for %s successfully triggered", p.Name)
	}
}

func main() {
	config := readConfig()
	ctx := context.Background()
	auth := codeship.NewBasicAuth(os.Getenv("CODESHIP_USERNAME"), os.Getenv("CODESHIP_PASSWORD"))
	client, err := codeship.New(auth)
	if err != nil {
		log.Fatalln("encountered issue authenticating")
	}
	org, err := client.Organization(ctx, config.OrgName)
	if err != nil {
		log.Fatalf("encountered issue selecting organization %s", config.OrgName)
	}
	println("ORG UUID =", org.UUID)
	for _, project := range config.Projects {
		log.Printf("Triggering build for %s", project.Name)
		trigger_build(org, project)
	}
}
