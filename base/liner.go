package base

import (
	"github.com/peterh/liner"
	"os"
	"path"
	"strings"
)

func keywordCompleter(line string) (c []string) {
	for _, value := range Keywords {
		if strings.HasPrefix(value, strings.TrimLeft(line, " ")) {
			c = append(c, value)
		}
	}
	return
}

func setoptsCompleter(line string) (c []string) {
	for _, value := range setopts {
		if strings.HasPrefix(value, strings.TrimLeft(line, " ")) {
			c = append(c, value)
		}
	}
	return
}

func showoptsCompleter(line string) (c []string) {
	for _, value := range showopts {
		if strings.HasPrefix(value, strings.TrimLeft(line, " ")) {
			c = append(c, value)
		}
	}
	return
}

func cmdCompleter(line string) (c []string) {
	for command := range XCommonCommandSet {
		if strings.HasPrefix(command, strings.TrimLeft(line, " ")) {
			c = append(c, command)
		}
	}
	return
}

func groupCompleter(line string) (c []string) {
	for key := range XHostMap {
		if strings.HasPrefix(key, strings.TrimLeft(line, " ")) {
			c = append(c, key)
		}
	}
	return
}

func wordCompleter(line string, pos int) (head string, completions []string, tail string) {
	head = string([]rune(line)[:pos])
	tail = string([]rune(line)[pos:])

	if strings.HasPrefix(head, ":") {
		fields := strings.Fields(head)
		switch fields[0] {
		case ":set":
			if len(fields) == 1 {
				if strings.HasSuffix(head, " ") {
					return head, setoptsCompleter(""), ""
				}
			}

			if len(fields) > 1 {
				if strings.HasSuffix(head, " ") {
					return head, setoptsCompleter(""), ""
				}

				newHead := strings.Join(fields[:len(fields)-1], " ")
				lastField := fields[len(fields)-1]
				if !strings.Contains(lastField, "=") {
					return newHead + " ", setoptsCompleter(lastField), ""
				} else {
					lastFields := strings.Split(lastField, "=")
					switch lastFields[0] {
					case "group":
						return newHead + " group=", groupCompleter(lastFields[1]), tail
					}
				}
			}
		case ":show":
			if len(fields) == 1 {
				if strings.HasSuffix(head, " ") {
					return head, showoptsCompleter(""), ""
				}
			}

			if len(fields) == 2 {
				if !strings.HasSuffix(head, " ") {
					return fields[0] + " ", showoptsCompleter(fields[1]), ""
				}
			}
		default:
			if len(fields) == 1 {
				if !strings.HasSuffix(head, " ") {
					return "", keywordCompleter(head), ""
				}
			}
		}
	}

	i := strings.LastIndex(head, XConfig.CommandSep)
	if i > 0 {
		return head[:i] + XConfig.CommandSep, cmdCompleter(head[i+1:]), ""
	} else {
		return "", cmdCompleter(head), ""
	}
}

func NewLiner() (*liner.State, error) {
	line := liner.NewLiner()
	line.SetWordCompleter(wordCompleter)
	line.SetTabCompletionStyle(liner.TabPrints)

	if f, err := os.Open(path.Join(RootPath, HisFile)); err == nil {
		line.ReadHistory(f)
		f.Close()
	}

	return line, nil
}

func Prompt(l *liner.State) (string, error) {
	name, err := l.Prompt(CurEnv.Prompt)
	if err == nil {
		l.AppendHistory(name)
	} else if err == liner.ErrPromptAborted {
		return name, PromptAborted
	} else {
		return "", PromptAborted
	}

	if f, err := os.Create(path.Join(RootPath, HisFile)); err != nil {
		Warn.Print("Error writing history file: ", err)
		return name, PromptHisErr
	} else {
		l.WriteHistory(f)
		f.Close()
	}

	return name, nil
}
