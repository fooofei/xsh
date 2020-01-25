package base

import (
	"gopkg.in/yaml.v2"
	"html/template"
	"io/ioutil"
	"log"
	"os"
)

// Action type : ssh/sftp/scp
type Action struct {
	Name     string   `yaml:"name"`
	Type     string   `yaml:"type"`
	Commands []string `yaml:"commands,omitempty"`
	Local    string   `yaml:"local,omitempty"`
	Remote   string   `yaml:"remote,omitempty"`
}

type Task struct {
	Name      string    `yaml:"name"`
	HostGroup HostGroup `yaml:"hostgroup"`
	Actions   []Action  `yaml:"actions"`
}

type xTask struct {
	Tasks []Task `yaml:"tasks"`
}

var XTask = xTask{}

func initXTask() {
	var taskFile = *taskFile

	if taskFile != "" {
		if *valueFile != "" {
			initValues(*valueFile)
		}
		if len(values) > 0 {
			templateTask(taskFile)
			taskFile = tempScript
		}

		tf, err := ioutil.ReadFile(taskFile)
		if err != nil {
			log.Fatalf("Can not read task file[%s], err: %v", taskFile, err)
		}

		err = yaml.Unmarshal(tf, &XTask)
		if err != nil {
			log.Fatalf("Task file[%s] unmarshal error: %v", taskFile, err)
		}

		if len(XTask.Tasks) == 0 {
			log.Fatalf("Task file[%s] empty", taskFile)
		}
	}
}

var values map[string]interface{}

func initValues(valueFile string) {
	p, err := ioutil.ReadFile(valueFile)
	if err != nil {
		log.Fatalf("Can not read value file[%s], err: %v", valueFile, err)
	}

	err = yaml.Unmarshal(p, &values)
	if err != nil {
		log.Fatalf("Value file[%s] unmarshal error: %v", valueFile, err)
	}
}

var tempScript = ConfigRootPath + "/" + TempFile

func templateTask(taskFile string) {
	t, err := template.ParseFiles(taskFile)
	if err != nil {
		log.Fatalf("Parse template task file[%s] failed", taskFile)
	}

	f, err := os.Create(tempScript)
	if err != nil {
		log.Fatalf("Create temp task file[%s] failed", tempScript)
	}
	defer f.Close()

	if err := t.Execute(f, values); err != nil {
		log.Fatalf("Execute template task file[%s] by value file[%s] task failed", taskFile, *valueFile)
	}
}
