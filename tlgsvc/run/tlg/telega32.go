// Run telegram-bot as service on Windows
// Copyright (c) 2023 Georgii Batanov gbatanov@yandex.rupackage tlg32
package tlg

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func init() {
	//	fmt.Println("init in telega 32")
}

// Структура сообщения
type Message struct {
	ChatId int64  // Идентификатор клиента
	Msg    string // Строка сообщения
}

// We inherit the bot to rewrite the function for receiving updates
type Tlg32 struct {
	botApi    *tgbotapi.BotAPI // АПИ из библиотеки Go
	mode      string           // режим "prod" | "test"
	MyId      int64            // Мой идентификатор
	botName   string           // Имя бота, зарегистрированное в телеграм
	chatIds   []int64          // Список идентификаторов, с которыми работает бот
	Flag      bool             // Признак разрешения работы бота
	tokenPath string           // Путь к файлу с токеном
	token     string           // Токен полученный при регистрации бота
	MsgChan   chan Message     // Канал для обмен сообщением с основной программой
	wg        *sync.WaitGroup
}

// Создание экземпляра бота
func Tlg32Create(botName string, mode string, tokenPath string, myId int64, msgChan chan Message, wg *sync.WaitGroup) *Tlg32 {
	bot := Tlg32{}
	bot.mode = mode
	bot.tokenPath = tokenPath
	bot.botName = botName //your bot name
	bot.MyId = myId
	bot.chatIds = append(bot.chatIds, myId)
	bot.Flag = true
	bot.MsgChan = msgChan
	bot.wg = wg
	return &bot
}

// Получение токена телеграм (должен храниться локально, не включать в репозиторий!)
func (bot *Tlg32) get_token() error {

	token, err := os.ReadFile(bot.tokenPath)
	if err != nil {
		return errors.New("incorrect file with token")
	}

	// remove trailing CR LF SPACE
	l := len(token)
	c := token[l-1]
	for c == 10 || c == 13 || c == 32 {
		token = token[:l-1]
		l = len(token)
		c = token[l-1]
	}
	bot.token = string(token)
	return nil

}

// Останов  бота
func (bot *Tlg32) Stop() {
	bot.Flag = false
	// толкнем очередь сообщений
	bot.MsgChan <- Message{ChatId: 0, Msg: "0"}
}

// Запуск бесконечных процессов получения новых сообщений и отправки готовых
func (bot *Tlg32) Run() error {
	var err error
	bot.get_token()
	if err != nil {
		log.Println(err.Error())
		return err
	}

	bot.botApi, err = tgbotapi.NewBotAPI(string(bot.token))
	if err != nil {
		return errors.New("incorrect token")
	}
	bot.botApi.Debug = bot.mode == "test"

	bot.wg.Add(1)
	go func() {
		bot.send_msg()
		bot.wg.Done()
	}()
	bot.wg.Add(1)
	go func() {
		uConfig := tgbotapi.NewUpdate(0)
		uConfig.Timeout = 60
		updates := bot.MyGetUpdatesChan(uConfig)
		for bot.Flag {

			for update := range updates { // чтобы прекратить этот прием, источник должен закрыть канал
				if !bot.Flag {
					bot.botApi.StopReceivingUpdates()
					bot.wg.Done()
					return
				}

				if update.Message != nil { // Есть новые входящие сообщения
					var outMsg string
					chatId := update.Message.Chat.ID
					msgIn := update.Message.Text
					firstName := update.Message.From.FirstName
					outMsg, err := bot.handle_msg_in(msgIn, chatId, firstName)
					if err != nil {
						outMsg = "I'm understand you"
					}
					if bot.Flag {
						bot.MsgChan <- Message{ChatId: chatId, Msg: outMsg}
					}
				}
			}
		}
		bot.wg.Done()
	}()
	return nil
}

// Функция отправки сообщений.
// Сообщение для отправки получаем из канала
func (bot *Tlg32) send_msg() {
	// Этот код включаем, если нужно цитирование принятого сообщения
	//			msg.ReplyToMessageID = update.Message.MessageID
	for bot.Flag {
		inMsg := <-bot.MsgChan
		chatId := inMsg.ChatId
		text := inMsg.Msg
		if chatId > 0 {
			msg := tgbotapi.NewMessage(chatId, text)
			bot.botApi.Send(msg)
		}
	}
}

// GetUpdatesChan возвращает канал для принятых сообщений и запускает горутину приема
func (bot *Tlg32) MyGetUpdatesChan(config tgbotapi.UpdateConfig) tgbotapi.UpdatesChannel {
	ch := make(chan tgbotapi.Update, bot.botApi.Buffer)
	go func() {
		for bot.Flag {

			updates, err := bot.botApi.GetUpdates(config)
			if err != nil {
				if strings.Contains(err.Error(), "Conflict:") {
					bot.Flag = false
					close(ch)
					return
				} else {
					continue
				}
			}

			// Отправляем новые сообщения в канал
			for _, update := range updates {
				if update.UpdateID >= config.Offset {
					config.Offset = update.UpdateID + 1
					ch <- update
				}
			}
		}
		close(ch) // Закрываем канал
	}()

	return ch
}

// Обработчик входящих сообщений
func (bot *Tlg32) handle_msg_in(msg string, chatId int64, firstName string) (string, error) {
	// ID входит в список разрешенных
	found := false
	for _, acc := range bot.chatIds {
		if chatId == acc {
			found = true
			break
		}
	}

	if strings.Contains(msg, "/start") {
		return fmt.Sprintf("Привет, %s!", firstName), nil
	}
	if strings.Contains(msg, "/stop") {
		bot.Flag = false
		return fmt.Sprintf("Good bye, %s!", firstName), nil
	}

	// "/mvreg" add in chatIds
	if strings.Contains(msg, "/mvreg") {
		bot.chatIds = append(bot.chatIds, chatId)
		return "Ok", nil
	}

	if found {
		//"/mvunreg" remove from chatIds
		if strings.Contains(msg, "/mvreg") {
			for i, acc := range bot.chatIds {
				if acc == chatId {
					copy(bot.chatIds[i:], bot.chatIds[i+1:])
					bot.chatIds = bot.chatIds[:len(bot.chatIds)-1]
					break
				}
			}
		}
		return "Ok", nil
	} else {
		return "Запрос не принят!", errors.New("запрос не принят")
	}
}
