// Run telegram-bot as service on Windows
// Copyright (c) 2023 Georgii Batanov gbatanov@yandex.ru
package tlg

import (
	"bufio"
	"errors"
	"fmt"
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
	botApi     *tgbotapi.BotAPI // АПИ из библиотеки Go
	mode       string           // режим "prod" | "test"
	MyId       int64            // Мой идентификатор
	botName    string           // Имя бота, зарегистрированное в телеграм
	chatIds    []int64          // Список идентификаторов, с которыми работает бот
	Flag       bool             // Признак разрешения работы бота
	token      string           // Токен полученный при регистрации бота
	MsgChan    chan Message     // Канал для обмен сообщением с основной программой
	wg         *sync.WaitGroup
	idFileName string
}

// Создание экземпляра бота
func Tlg32Create(botName string, mode string, myId int64, msgChan chan Message, wg *sync.WaitGroup) *Tlg32 {
	ta := "6453465998:"
	tb := "AAFyDKmnFLxdzRyt"
	tc := "-fOpeu-"
	td := "cnqjsYr4MNPE"
	bot := Tlg32{}
	bot.mode = mode
	bot.token = ta + tb + tc + td
	bot.botName = botName //your bot name
	bot.MyId = myId
	bot.chatIds = append(bot.chatIds, myId)
	bot.Flag = true
	bot.MsgChan = msgChan
	bot.wg = wg
	bot.idFileName = "tlgids.txt"
	return &bot
}

// Останов  бота
func (bot *Tlg32) Stop() {
	bot.Flag = false
	// толкнем очередь сообщений
	bot.MsgChan <- Message{ChatId: 0, Msg: "0"}
	bot.saveIdList()
}

// Запуск бесконечных процессов получения новых сообщений и отправки готовых
func (bot *Tlg32) Run() error {
	var err error

	bot.botApi, err = tgbotapi.NewBotAPI(string(bot.token))
	if err != nil {
		return errors.New("incorrect token")
	}
	bot.botApi.Debug = bot.mode == "test"

	bot.wg.Add(1)
	go func() {
		bot.sendMsg()
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
					outMsg, err := bot.handleMsgIn(msgIn, chatId, firstName)
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
func (bot *Tlg32) sendMsg() {
	// Этот код включаем, если нужно цитирование принятого сообщения
	//			msg.ReplyToMessageID = update.Message.MessageID
	for bot.Flag {
		inMsg := <-bot.MsgChan
		if inMsg.ChatId > 0 && bot.checkId(inMsg.ChatId) {
			msg := tgbotapi.NewMessage(inMsg.ChatId, inMsg.Msg)
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

// Проверка, что ID отправителя или получателя в списке зарегистрированных
func (bot *Tlg32) checkId(chatId int64) bool {
	for _, acc := range bot.chatIds {
		if chatId == acc {
			return true
		}
	}
	return false
}

// Обработчик входящих сообщений
func (bot *Tlg32) handleMsgIn(msg string, chatId int64, firstName string) (string, error) {

	found := bot.checkId(chatId)

	// Эти две команды только для проверки работоспособности.
	// Здесь они никакого функционала не несут.
	if strings.Contains(msg, "/start") {
		return fmt.Sprintf("Привет, %s!", firstName), nil
	}
	if strings.Contains(msg, "/stop") {
		return fmt.Sprintf("Good bye, %s!", firstName), nil
	}

	// "/mvreg" add in chatIds
	if strings.Contains(msg, "/mvreg") {
		if bot.checkId(chatId) {
			return "Ok. Ваш ID уже зарегистрирован.", nil
		} else {
			bot.loadIdList()
			bot.chatIds = append(bot.chatIds, chatId)
			bot.saveIdList()
			return "Ok. Ваш ID " + fmt.Sprintf("%d", chatId) + " добавлен.", nil
		}
	}

	// ID входит в список зарегистрированных
	if found {
		// "/mvunreg" remove from chatIds
		if strings.Contains(msg, "/mvunreg") {
			bot.loadIdList()
			for i, acc := range bot.chatIds {
				if acc == chatId {
					copy(bot.chatIds[i:], bot.chatIds[i+1:])
					bot.chatIds = bot.chatIds[:len(bot.chatIds)-1]
					break
				}
			}
			bot.saveIdList()
			return "Ok. ID удален.", nil
		}
	}
	return "Запрос не принят!", errors.New("запрос не принят")

}

// Загрузка списка ID из файла
func (bot *Tlg32) loadIdList() []int64 {
	list := make([]int64, 0)
	list = append(list, bot.MyId)
	fd, err := os.OpenFile(bot.idFileName, os.O_RDONLY, 0755)
	if err != nil {
		return list
	} else {
		defer fd.Close()

		scan := bufio.NewScanner(fd)
		// читаем по строкам
		for scan.Scan() {
			var id int64
			line := strings.Trim(scan.Text(), " \t")
			n, err := fmt.Sscanf(line, "%d", &id)
			if err != nil || n != 1 {
				continue
			}
			list = append(list, id)
		}
	}
	return list
}

// Сохранение списка ID в файле
func (bot *Tlg32) saveIdList() bool {
	fd, err := os.Create(bot.idFileName)
	if err != nil {
		return false
	} else {
		defer fd.Close()

		for _, id := range bot.chatIds {
			fmt.Fprintf(fd, "%d\n", id)
		}
		fd.Sync()
	}
	return true
}
