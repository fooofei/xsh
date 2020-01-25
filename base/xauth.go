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
	PrivateKey string `yaml:"privatekey,omitempty"`
	Passphrase string `yaml:"passphrase,omitempty"`
	SuType     string `yaml:"sutype,omitempty"`
	SuPass     string `yaml:"supass,omitempty"`
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
		if value.Password, err = GetPlainPassword(value.Password); err != nil {
			Error.Printf("auth[%s] password decrypt error: %v", value.Name, err)
		}
		if value.Passphrase, err = GetPlainPassword(value.Passphrase); err != nil {
			Error.Printf("auth[%s] passphrase decrypt error: %v", value.Name, err)
		}
		if value.SuPass, err = GetPlainPassword(value.SuPass); err != nil {
			Error.Printf("auth[%s] suPass decrypt error: %v", value.Name, err)
		}
		XAuthMap[value.Name] = value
	}

	if len(XAuth.Auths) != len(XAuthMap) {
		log.Fatal("Auth duplicate")
	}
}
