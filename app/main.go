package main

import (
	"fmt"
	"math/rand"
	"neuro/logger"
	"time"
)

func main() {
	RunApp()

	defer logger.Log.Shutdown()

	actions := []func(){
		func() {
			logger.Log.ChatLog("$2Пользователь$ $6ЮЗЕРНЕЙМ$ $2отправил сообщение пользователю$ $6Другой юзернейм$", 6)
		},
		func() {
			logger.Log.ErrorLog("$4Произошла ошибка во время запуска$ $3сетевого$ $4модуля$", 4)
		},
		func() {
			logger.Log.InfoLog("$3Сервер успешно запушен на$ $5IP:$$6 1.1.1.1$ $3и$ $6порте:$$6 65355$", 3)
		},
		func() {
			logger.Log.InfoLog("$2Пользователь$ $6USER1$ $2передал$$5 500$$ $2игроку$ $6USER2$", 2)
		},
		func() {
			logger.Log.InfoLog("$3Магазин:$ $6ITEM_NAME$ $3был куплен пользователем$ $6BUYER$ $3за$$5 1000$ $3монет$", 3)
		},
		func() {
			logger.Log.ChatLog("$4Администратор$ $6MODERATOR$ $4заблокировал пользователя$ $6SPAMMER$ $4на$ $2 24 часа$", 4)
		},
		func() {
			logger.Log.ChatLog("$2Хелпер$ $6STAFF_NAME$ $2разморозил игрока$ $6PLAYER$", 2)
		},
		func() {
			logger.Log.ErrorLog("$4Критическая ошибка:$ $3Не удалось подключиться к базе данных$ $6MySQL$ $4по адресу$$6 127.0.0.1$", 4)
		},
		func() {
			logger.Log.ErrorLog("$4Потеряно соединение с$ $3центральным$ $4узлом. Попытка реконнекта...$", 4)
		},
		func() {
			logger.Log.InfoLog("$3Загрузка конфигурации$ $6config.json$ $3завершена за$$5 145мс$", 3)
		},
		func() {
			logger.Log.InfoLog("$2Система безопасности$ $5Anticheat$ $2успешно инициализирована$", 2)
		},
	}

	start := time.Now()
	for range 10 {
		actions[rand.Intn(11)]()
	}
	fmt.Printf("Время работы программы: %v\n", time.Since(start))
}
