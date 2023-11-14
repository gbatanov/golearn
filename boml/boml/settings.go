package boml

import (
	"bufio"
	"errors"
	"os"
	"reflect"
	"strconv"
	"strings"
)

// Пример файла конфигурации
type BomlConfig struct {
	Mode             string `boml:"mode"`
	User_variant     int    `boml:"user_variant"`
	User_count       int    `boml:"user_count"`
	Threads          int    `boml:"threads"`
	Cycles           int    `boml:"cycles"`
	Dir_create       int    `boml:"dir_create"`
	Domain           string `boml:"domain"`
	Server_name      string `boml:"server_name"`
	Share_folder     string `boml:"share_folder"`
	Execute_user     string `boml:"execute_user"`
	Execute_user_pwd string `boml:"execute_user_pwd"`
	Common_user_pwd  string `boml:"common_user_pwd"`
	To_log           int    `boml:"to_log"`
	With_office      int    `boml:"with_office"`
	Port             string `boml:"port"`
}

// Загрузка из секционированного conf-файла
// Параметры в конкретной секции заменяют общие  параметры до секций
func (config *BomlConfig) Load(filename string) error {
	fd, err := os.OpenFile(filename, os.O_RDONLY, 0755)
	if err != nil {
		return errors.New("incorrect file with configuration")
	} else {
		scan := bufio.NewScanner(fd)
		var mode string = ""
		var sectionMode bool = true
		var values []string = []string{}
		// read line by line
		for scan.Scan() {

			line := scan.Text()
			line = strings.Trim(line, " \t")

			if strings.HasPrefix(line, "//") || strings.HasPrefix(line, "#") || strings.HasPrefix(line, ";") { //comment
				continue
			}
			config.trailComments(&line)

			if len(line) < 3 { //empty string
				continue
			}
			if len(mode) == 0 {
				values = strings.Split(line, "=")
				if strings.ToLower(strings.Trim(values[0], " \t")) != "mode" {
					continue
				}
				mode = strings.Trim(line[len(values[0]):], " \t")
				mode = strings.ToLower(strings.Trim(mode, "="))
				config.Mode = mode
				continue
			}

			if strings.HasPrefix(line, "[") {
				section := line[1 : len(line)-1]
				sectionMode = section == mode
				continue
			}
			if !sectionMode { //pass section
				continue
			}
			values := strings.Split(line, "=")
			values1 := strings.Trim(line[len(values[0]):], " \t")
			values1 = strings.Trim(values1, "=")
			values1 = strings.Trim(values1, " \t")
			values0 := strings.ToLower(values[0])

			t := reflect.TypeOf(*config)
			// указатель на структуру - addressable, не имеет свойства Tag, поэтому еще используем reflect.TypeOf
			ps := reflect.ValueOf(config)
			// сама структура
			s := ps.Elem()
			if s.Kind() == reflect.Struct {
				for i := 0; i < s.NumField(); i++ {
					f := s.Field(i)
					field := t.Field(i)
					tag := field.Tag.Get("boml")
					if tag == values0 {
						//log.Printf(" %v (%v), tag: '%v'\n", field.Name, field.Type.Name(), tag)
						switch field.Type.Name() {
						case "string":
							f.SetString(values[1])
						case "int":
							in, err := strconv.Atoi(values1)
							if err == nil {
								// несмотря на то, что на входе и на выходе int, reflect требует int64
								f.SetInt(int64(in))
							}
						}
						break
					}
				}
			}
		}
		fd.Close()
	}
	return nil
}

// Обрезка хвоста после комментария
func (conf *BomlConfig) trailComments(line *string) {
	if strings.Contains(*line, "//") { // tail comment
		spl := strings.Split(*line, "//")
		*line = strings.Trim(spl[0], " \t")
	}
	if strings.Contains(*line, "#") { // tail comment
		spl := strings.Split(*line, "#")
		*line = strings.Trim(spl[0], " \t")
	}
	if strings.Contains(*line, ";") { // tail comment
		spl := strings.Split(*line, ";")
		*line = strings.Trim(spl[0], " \t")
	}

}
