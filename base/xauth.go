package base

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"path"
)

type Authentication struct {
	Name       string `yaml:"name,omitempty"`
	Username   string `yaml:"username,omitempty"`
	Password   string `yaml:"password,omitempty"`
	PrivateKey string `yaml:"privatekey,omitempty"`
	Passphrase string `yaml:"passphrase,omitempty"`
	SuType     string `yaml:"sutype,omitempty"`
	SuPass     string `yaml:"supass,omitempty"`
}

type xAuth struct {
	Authentications []Authentication `yaml:"authentications,omitempty"`
}

var XAuth = xAuth{}
var XAuthMap = make(map[string]Authentication)

func initXAuth() {
	var filePath = path.Join(ConfigRootPath, AuthenticationsFile)

	a, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Fatalf("Can not read authentications file[%s].", filePath)
	}
	err = yaml.Unmarshal(a, &XAuth)
	if err != nil {
		log.Fatalf("Authentications[%s] unmarshal error: %v", filePath, err)
	}

	if len(XAuth.Authentications) == 0 {
		Warn.Printf("The authentications empty.")
	}

	for _, value := range XAuth.Authentications {
		if !CheckName(value.Name) {
			log.Fatalf("Authentication name [%s] illegal", value.Name)
		}
		value.Password = GetPlainPassword(value.Password)
		value.Passphrase = GetPlainPassword(value.Passphrase)
		value.SuPass = GetPlainPassword(value.SuPass)
		XAuthMap[value.Name] = value
	}

	if len(XAuth.Authentications) != len(XAuthMap) {
		log.Fatal("Authentication duplicate")
	}
}
