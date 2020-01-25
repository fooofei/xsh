package base

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"path"
)

type Auth struct {
	Name       string `yaml:"name"`
	Username   string `yaml:"username"`
	Password   string `yaml:"password,omitempty"`
	PrivateKey string `yaml:"private_key,omitempty"`
	Passphrase string `yaml:"passphrase,omitempty"`
	SuType     string `yaml:"su_type,omitempty"`
	SuPass     string `yaml:"su_pass,omitempty"`
}

type xAuth struct {
	Auths []Auth `yaml:"auths"`
}

var XAuth = xAuth{}
var XAuthMap = make(map[string]Auth)

func InitXAuth() {
	var filePath = path.Join(RootPath, AuthFile)

	a, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Fatalf("Can not read auth file[%s].", filePath)
	}
	err = yaml.Unmarshal(a, &XAuth)
	if err != nil {
		log.Fatalf("Auths[%s] unmarshal error: %v", filePath, err)
	}

	if len(XAuth.Auths) == 0 {
		Warn.Printf("The auth empty.")
	}

	for _, value := range XAuth.Auths {
		if !CheckName(value.Name) {
			log.Fatalf("Auth name [%s] illegal", value.Name)
		}
		value.Password = GetPlainPassword(value.Password)
		value.Passphrase = GetPlainPassword(value.Passphrase)
		value.SuPass = GetPlainPassword(value.SuPass)
		XAuthMap[value.Name] = value
	}

	if len(XAuth.Auths) != len(XAuthMap) {
		log.Fatal("Auth duplicate")
	}
}
