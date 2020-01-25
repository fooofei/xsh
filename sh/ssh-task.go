package sh

import (
	"fmt"
	. "github.com/xied5531/xsh/base"
	"gopkg.in/yaml.v2"
	"html/template"
	"io/ioutil"
	"os"
)

type SshTask struct {
	Name       string      `yaml:"name"`
	SshActions []SshAction `yaml:"actions"`
}

type SshTaskResult struct {
	Name             string
	SshActionResults []SshActionResult
	Err              error
}

var xTask SshTask

func (s SshTask) applyValue() error {
	if *Task == "" {
		return fmt.Errorf("task file not found")
	}

	var values map[string]interface{}
	if *Value != "" {
		vf, err := ioutil.ReadFile(*Value)
		if err != nil {
			return fmt.Errorf("can not read value file[%s], err: %v", *Value, err)
		}

		err = yaml.Unmarshal(vf, &values)
		if err != nil {
			return fmt.Errorf("value file[%s] unmarshal error: %v", *Value, err)
		}
	}

	if len(values) == 0 {
		tf, err := ioutil.ReadFile(*Task)
		if err != nil {
			return fmt.Errorf("can not read task file[%s], err: %v", *Task, err)
		}

		err = yaml.Unmarshal(tf, &xTask)
		if err != nil {
			return fmt.Errorf("task file[%s] unmarshal error: %v", *Task, err)
		}
	} else {
		tf, err := template.ParseFiles(*Task)
		if err != nil {
			return fmt.Errorf("parse template task file[%s] error: %v", *Task, err)
		}

		tmp, err := ioutil.TempFile(TempPath, "xsh-task-*.yaml")
		if err != nil {
			return fmt.Errorf("create temp file in [%s] error: %v", TempPath, err)
		}
		defer func() {
			tmp.Close()
			os.Remove(tmp.Name())
		}()

		if err := tf.Execute(tmp, values); err != nil {
			return fmt.Errorf("execute template task file[%s] by value file[%s] task error: %v", *Task, *Value, err)
		}

		rtf, err := ioutil.ReadFile(tmp.Name())
		if err != nil {
			return fmt.Errorf("can not read task file[%s], err: %v", tmp.Name(), err)
		}

		err = yaml.Unmarshal(rtf, &xTask)
		if err != nil {
			return fmt.Errorf("task file[%s] unmarshal error: %v", tmp.Name(), err)
		}
	}

	if len(xTask.SshActions) == 0 {
		return fmt.Errorf("Task file[%s] action empty", *Task)
	}

	return nil
}

func (s SshTask) checkTask() error {
	if xTask.Name == "" {
		return fmt.Errorf("task name empty")
	}
	if len(xTask.SshActions) == 0 {
		return fmt.Errorf("task actions empty")
	}
	for _, action := range xTask.SshActions {
		if action.Name == "" {
			return fmt.Errorf("action name empty")
		}
		if _, ok := XHostMap[action.Group]; !ok {
			return fmt.Errorf("action[%s] group[%s] not found", action.Name, action.Group)
		}
		if len(action.Steps) == 0 {
			return fmt.Errorf("action[%s] steps empty", action.Name)
		}
		for _, step := range action.Steps {
			if step.Type != "command" && step.Type != "copy" {
				return fmt.Errorf("action[%s] step type illegal", action.Name)
			}
			if step.Type == "command" {
				if len(step.Commands) == 0 {
					return fmt.Errorf("action[%s] step commands empty", action.Name)
				}
			}
			if step.Type == "copy" {
				if step.Direction != "upload" && step.Direction != "download" {
					return fmt.Errorf("action[%s] step direction illegal", action.Name)
				}
				if step.Local == "" || step.Remote == "" {
					return fmt.Errorf("action[%s] step local or remote empty", action.Name)
				}
			}
		}
	}

	return nil
}

func (s SshTask) Do() SshTaskResult {
	result := SshTaskResult{}
	if err := s.applyValue(); err != nil {
		result.Err = err
		return result
	}

	if err := s.checkTask(); err != nil {
		result.Err = err
		return result
	}

	result.Name = xTask.Name

	for _, action := range xTask.SshActions {
		actionResult := action.Do()
		result.SshActionResults = append(result.SshActionResults, actionResult)
		if actionResult.Err != nil {
			result.Err = fmt.Errorf("Interrupted due to an error: %s\n", actionResult.Err.Error())
			break
		}
	}

	return result
}
