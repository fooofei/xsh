package sh

import (
	"fmt"
	. "github.com/xied5531/xsh/base"
	"gopkg.in/yaml.v2"
	"html/template"
	"io/ioutil"
	"os"
	"path/filepath"
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
			return fmt.Errorf("parse template task file[%s] failed", *Task)
		}

		tmp, err := ioutil.TempFile(TempPath, "xsh-task-*.yaml")
		if err != nil {
			return fmt.Errorf("create temp file in [%s] failed", TempPath)
		}
		defer os.Remove(filepath.Join(TempPath, tmp.Name()))

		if err := tf.Execute(tmp, values); err != nil {
			return fmt.Errorf("execute template task file[%s] by value file[%s] task failed", *Task, *Value)
		}

		realTaskFile := filepath.Join(TempPath, tmp.Name())
		rtf, err := ioutil.ReadFile(realTaskFile)
		if err != nil {
			return fmt.Errorf("can not read task file[%s], err: %v", realTaskFile, err)
		}

		err = yaml.Unmarshal(rtf, &xTask)
		if err != nil {
			return fmt.Errorf("task file[%s] unmarshal error: %v", realTaskFile, err)
		}
	}

	if len(xTask.SshActions) == 0 {
		return fmt.Errorf("Task file[%s] action empty", *Task)
	}

	return nil
}

func (s SshTask) Do() SshTaskResult {
	result := SshTaskResult{}
	result.Name = xTask.Name

	if err := s.applyValue(); err != nil {
		result.Err = err
		return result
	}

	for _, action := range xTask.SshActions {
		result.SshActionResults = append(result.SshActionResults, action.Do())
	}

	return result
}
