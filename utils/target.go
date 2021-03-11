package utils

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	yaml "gopkg.in/yaml.v2"
)

// Target : Target struct
type Target struct {
	FollowersNextCursor string       `yaml:"followersNextCursor"`
	FriendsNextCursor   string       `yaml:"friendsNextCursor"`
	StatusesNextCursor  string       `yaml:"statusesextCursor"`
	User                string       `yaml:"user"`
	Credentials         *Credentials `yaml:"credentials"`
	Path                string       `yaml:"path"`
}

// NewTarget : Create Folders & Files
func NewTarget(user string, conf string) *Target {
	u := &Target{
		User:                user,
		FollowersNextCursor: "-1",
		FriendsNextCursor:   "-1",
		StatusesNextCursor:  "-1",
	}
	confFile, err := ioutil.ReadFile(conf)
	ExitOnError(err)
	yaml.Unmarshal(confFile, &u)

	basePath := "./targets/" + u.User

	if _, err = os.Stat(basePath + "/img"); os.IsNotExist(err) {
		os.MkdirAll(basePath+"/img", os.ModePerm)
		fmt.Printf("Initialized new target: '%s', under: './targets/%s'\n", u.User, u.User)
	} else {
		fmt.Printf("Target: '%s', allready initialized under: './targets/%s'\n", u.User, u.User)
	}

	u.Path, _ = filepath.Abs(basePath)

	c := CreateFile(basePath + "/" + u.User + ".yaml")
	// CreateFile(basePath + "/" + "statuses.csv")
	// CreateFile(basePath + "/" + "likes.csv")
	// CreateFile(basePath + "/" + "friends.csv")
	// CreateFile(basePath + "/" + "followers.csv")
	data, _ := yaml.Marshal(&u)
	ioutil.WriteFile(c, data, 0644)

	return u
}

// GetTarget : Returns Target Struct
func GetTarget(user string) *Target {
	basePath := "./targets/" + user + "/" + user + ".yaml"

	u := &Target{}

	confFile, err := ioutil.ReadFile(basePath)
	ExitOnError(err)
	yaml.Unmarshal(confFile, &u)

	return u
}

// SaveTarget : Updates conf file
func SaveTarget(target *Target) {
	data, _ := yaml.Marshal(target)
	ioutil.WriteFile(target.Path+"/"+target.User+".yaml", data, 0644)
}
