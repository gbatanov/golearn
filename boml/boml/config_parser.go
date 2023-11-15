package boml

import (
	"bufio"
	"errors"
	"log"
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
		var mode string = "" // имя рабочей секции
		var inSection = true
		// read line by line
		for scan.Scan() {

			line := strings.Trim(scan.Text(), " \t")

			if strings.HasPrefix(line, "//") || strings.HasPrefix(line, "#") || strings.HasPrefix(line, ";") { //comment
				continue
			}
			config.trailComments(&line)

			if len(line) < 3 { //пустая строка - минимальная строка a=b (3 символа)
				continue
			}

			if strings.HasPrefix(line, "[") {
				section := line[1 : len(line)-1]
				log.Printf("section %s mode %s\n", section, mode)
				inSection = section == mode
				continue // Пропускаем строку с именем секции
			}
			if !inSection {
				// Пропускаем строки внутри ненужной секции
				continue
			}
			values := strings.Split(line, "=")
			vKey := strings.Trim(strings.ToLower(values[0]), " \t")
			vVal := strings.Trim(line[len(values[0]):], "= \t")

			t := reflect.TypeOf(*config)
			// ps - указатель на структуру - addressable,
			// в нем поле Field не имеет свойства Tag, поэтому еще используем reflect.TypeOf
			ps := reflect.ValueOf(config)
			// сама структура
			s := ps.Elem()
			if s.Kind() == reflect.Struct {
				for i := 0; i < s.NumField(); i++ {
					f := s.Field(i)
					field := t.Field(i)
					tag := field.Tag.Get("boml")
					if tag == vKey {
						//log.Printf(" %v (%v), tag: '%v'\n", field.Name, field.Type.Name(), tag)
						switch field.Type.Name() {
						// В этой конфигурации используются только строки и целые числа
						case "string":
							if vKey == "mode" {
								mode = vVal
							}
							f.SetString(vVal)
						case "int":
							in, err := strconv.Atoi(vVal)
							if err == nil {
								// несмотря на то, что на входе и на выходе int, reflect требует int64
								f.SetInt(int64(in))
							}
						}
						break // прекращаем перебор полей
					}
				}
			}
		}
		fd.Close()
	}
	return nil
}

// Обрезка хвоста после комментария
// Символ комментария выбирается самый левый
func (conf *BomlConfig) trailComments(line *string) {
	seps := []string{"//", "#", ";"}
	pos := len(*line)
	sep := ""
	for _, sp := range seps {
		pos1 := strings.Index(*line, sp)
		if pos1 > 0 && pos1 < pos {
			pos = pos1
			sep = sp
		}
	}
	if sep != "" {
		spl := strings.Split(*line, sep)
		*line = strings.Trim(spl[0], " \t")
	}
}
