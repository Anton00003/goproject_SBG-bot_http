package api

import (
	//	"encoding/csv"
	"log"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Api struct {
	Service service
}

type service interface {
	ChekAvtorisation(chatID int64) bool
	OutList(chatID int64) string
	AddName(chatID int64) error
	AddNameWork(text_vvod string, chatID int64) error
	DeleteName(chatID int64) (string, error)
	DeleteNameWork(text_vvod string, chatID int64) error
	Сancel(chatID int64) error
	GetWorker() ([][]int64, []string)
	GetPrevious(chatID int64) string
	EnterName(text string, chatID int64) error
	EnterDate(text string, chatID int64) error
}

func New(s service) *Api {
	return &Api{Service: s}
}

func (a *Api) Run(s service) {

	bot, err := tgbotapi.NewBotAPI("6825004600:AAE-h0KLA3tJ4AjIyxH6PU0cCVIrPDFZrsM")
	if err != nil {
		log.Panic(err)
	}

	log.Printf("Authorized on account %s", bot.Self.UserName)

	//++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
	go a.worker(bot)
	//++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil { // Есть новое сообщение
			text := update.Message.Text      // Текст сообщения
			chatID := update.Message.Chat.ID //  ID чата
			userID := update.Message.From.ID // ID пользователя
			var replyMsg string

			log.Printf("[%s](%d) %s", update.Message.From.UserName, userID, text)

			// Анализируем текст сообщения и записываем ответ в переменную

			replyMsg = a.Distribution_answers(text, chatID)

			// Отправляем ответ
			msg := tgbotapi.NewMessage(chatID, replyMsg) // Создаем новое сообщение
			//			msg.ReplyToMessageID = update.Message.MessageID // Указываем сообщение, на которое нужно ответить

			bot.Send(msg)
		}
	}
}

func (a *Api) worker(bot *tgbotapi.BotAPI) {

	for {
		slice_ID, to_name := a.Service.GetWorker()

		for i, name := range to_name {
			for _, v_ID := range slice_ID[i] {
				replyMsg := "скоро ДР у " + name
				msg := tgbotapi.NewMessage(v_ID, replyMsg)
				bot.Send(msg)
			}
		}
		time.Sleep(time.Second * 10)
	}
}

func (a *Api) Distribution_answers(text_vvod string, chatID int64) string {
	var replyMsg string

	if text_vvod == "/cancel" {
		//		replyMsg = a.Service.Сancel(chatID)
		replyMsg = a.Сancel(chatID)
		return replyMsg
	}
	if a.Service.ChekAvtorisation(chatID) == false {
		replyMsg = a.Autorisation(text_vvod, chatID)
	} else {

		if text_vvod == "/list_name" {
			replyMsg = a.OutList(chatID)

		} else if text_vvod == "/add_name" {
			replyMsg = a.AddName(chatID)

		} else if text_vvod == "/delete_name" {
			replyMsg = a.DeleteName(chatID)

		} else if a.Service.GetPrevious(chatID) == "Add" {
			replyMsg = a.AddNameWork(text_vvod, chatID)

		} else if a.Service.GetPrevious(chatID) == "Delete" {
			replyMsg = a.DeleteNameWork(text_vvod, chatID)

		} else {
			replyMsg = Menu()
		}
	}

	return replyMsg
}

func (a *Api) Сancel(chatID int64) string {
	err := a.Service.Сancel(chatID)
	if err != nil {
		return "Ошибка отмены"
	}
	return "Отменено"
}

func (a *Api) Autorisation(text_vvod string, chatID int64) string {

	err := a.Service.EnterName(text_vvod, chatID)
	if err != nil {
		return "Авторизация\n\nИмя пользователя не введео или введено неверно\nВведите имя пользователя\nили отмена действия /cancel"
	}

	err = a.Service.EnterDate(text_vvod, chatID)
	if err != nil {
		return "Авторизация\n\nДата рождения пользователя не введена или введена неверно\nВведите дату рождения пользователя в формате:\nГГГГ-ММ-ДД\nнапример:\n2003-10-28\nили отмена действия /cancel"
	}

	return "Авторизация завершена" + "\n" + "\n" + Menu()
}

func (a *Api) OutList(chatID int64) string {
	return a.Service.OutList(chatID)
}

func (a *Api) AddName(chatID int64) string {
	err := a.Service.AddName(chatID)
	if err != nil {
		return "Поользватель не прошел регистрацию"
	}
	return "Введите имя пользователя\nили отмена действия /cancel"
}

func (a *Api) AddNameWork(text_vvod string, chatID int64) string {
	err := a.Service.AddNameWork(text_vvod, chatID)
	if err != nil {
		return "Пользователь не существует или уже добавлен\n" +
			"Введите другое имя\n" +
			"или отмена действия /cancel"
	}
	return "Пользователь успешно добавлен"
}

func (a *Api) DeleteName(chatID int64) string {
	listName, err := a.Service.DeleteName(chatID)
	if err != nil {
		return "Список на кого подписан пуст"
	}
	return "Введите имя пользователя из списка подписанных:" + listName + "\nили отмена действия /cancel"
}

func (a *Api) DeleteNameWork(text_vvod string, chatID int64) string {
	err := a.Service.DeleteNameWork(text_vvod, chatID)
	if err != nil {
		return "Пользователь с указанным именем не зарегестрирован либо\n" +
			"Вы не подписаны на пользователя с указанным именем\n" +
			"Введите другое имя\n" +
			"или отмена действия /cancel"
	}
	return "Пользователь успешно удален"
}

func Menu() string {
	text := "Доступны следующие действия:\n" +
		"получить список пользователей - /list_name\n" +
		"подписаться на пользователя - /add_name\n" +
		"отписаться от пользователя - /delete_name\n"

	return text
}
