package boml

import (
	"bufio"
	"errors"
	"os"
	"reflect"
	"strconv"
	"strings"
)

var BomlComments []string = []string{"//"} /* # ; не допускаются, потому что теоретически могут быть в пароле */

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

			config.RemoveComments(&line)

			if len(line) < 3 { //пустая строка - минимальная строка a=b (3 символа)
				continue
			}

			if strings.HasPrefix(line, "[") {
				section := line[1 : len(line)-1]
				inSection = section == mode
				continue // Пропускаем строку с именем секции
			}
			if !inSection {
				// Пропускаем строки внутри ненужной секции
				continue
			}
			values := strings.Split(line, "=")
			if len(values) < 2 {
				continue
			}
			vKey := strings.Trim(strings.ToLower(values[0]), " \t")
			vVal := strings.Trim(values[1], " \t")

			s := reflect.ValueOf(config).Elem()
			typeOfT := s.Type()
			if s.Kind() == reflect.Struct {
				for i := 0; i < s.NumField(); i++ {
					f := s.Field(i)
					tag := typeOfT.Field(i).Tag.Get("boml")
					if tag == "" { // тэг не прописан для этого поля, берем имя поля в LowerCase
						tag = strings.ToLower(typeOfT.Field(i).Name)
					}
					if tag == vKey && f.CanSet() {
						switch f.Type().Name() {
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

// Обнуление строки с комментарием
// Обрезка хвоста после комментария (Символ комментария выбирается самый левый)
func (conf *BomlConfig) RemoveComments(line *string) {
	seps := BomlComments
	pos := len(*line)
	sep := ""
	for _, sp := range seps {
		if strings.HasPrefix(*line, sp) {
			*line = ""
			break
		}
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
