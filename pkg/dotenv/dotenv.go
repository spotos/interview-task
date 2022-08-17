package dotenv

import (
	"bufio"
	"os"
	"strings"

	"github.com/ardanlabs/conf/v3"
	"github.com/pkg/errors"
)

func FromEnvFiles(fileNames ...string) conf.Parsers {
	var files []*os.File

	for _, fileName := range fileNames {
		_, err := os.Stat(fileName)

		if err != nil {
			continue
		}

		file, err := os.Open(fileName)
		if err != nil {
			continue
		}

		files = append(files, file)
	}

	return &parser{
		files: files,
	}
}

type parser struct {
	files []*os.File
}

func (p *parser) Process(_ string, _ interface{}) error {
	envMap := map[string]string{}

	for _, file := range p.files {
		scanner := bufio.NewScanner(file)

		for i := 1; scanner.Scan(); i++ {
			if key, value, err := parse(scanner.Text()); err != nil {
				return errors.Errorf("error in file %v, line %v; err: %v", file.Name(), i, err)
			} else if key != "" {
				envMap[key] = value
			}
		}

		if err := scanner.Err(); err != nil {
			return errors.Errorf("error when scanning file %v; err: %v", file.Name(), err)
		}
	}

	currentEnv := map[string]bool{}
	rawEnv := os.Environ()
	for _, rawEnvLine := range rawEnv {
		key := strings.Split(rawEnvLine, "=")[0]
		currentEnv[key] = true
	}

	for key, value := range envMap {
		if !currentEnv[key] {
			_ = os.Setenv(key, value)
		}
	}

	return nil
}

// parse extracts a key/value pair from the given dot env (.env) single line.
func parse(line string) (string, string, error) {
	line = strings.TrimSpace(line)
	keyValue := []string{"", ""}
	position := 0
	iq := false
	qt := "'"

	for i := 0; i < len(line); i++ {
		if string(line[i]) == "#" && position == 0 {
			break
		}

		if string(line[i]) == "#" && position == 1 && iq == false {
			break
		}

		if string(line[i]) == "=" && position == 0 {
			position = 1
			continue
		}

		if string(line[i]) == " " && position == 1 {
			if iq == false && keyValue[position] == "" {
				continue
			}
		}

		if (string(line[i]) == "\"" || string(line[i]) == "'") && position == 1 {
			if keyValue[position] == "" {
				iq = true
				qt = string(line[i])
				continue
			} else if iq == true && qt == string(line[i]) {
				break
			}
		}

		keyValue[position] += string(line[i])
	}

	keyValue[0] = strings.TrimSpace(keyValue[0])
	if iq == false {
		keyValue[1] = strings.TrimSpace(keyValue[1])
	}

	if (position == 0 && keyValue[0] != "") || (position == 1 && keyValue[0] == "") {
		return "", "", errors.New("invalid syntax")
	}

	return keyValue[0], keyValue[1], nil
}
