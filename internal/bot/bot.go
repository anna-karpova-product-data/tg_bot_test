package bot

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"tg_bot/internal/speech"
)

type Bot struct {
	api *tgbotapi.BotAPI
	// Map to store conversation states
	conversationStates map[int64]string
	// Speech recognition client
	speechClient *speech.DeepgramClient
}

func New(telegramToken, deepgramToken string) (*Bot, error) {
	api, err := tgbotapi.NewBotAPI(telegramToken)
	if err != nil {
		return nil, err
	}

	return &Bot{
		api:                api,
		conversationStates: make(map[int64]string),
		speechClient:       speech.NewDeepgramClient(deepgramToken),
	}, nil
}

// analyzeMood analyzes the text and returns the detected mood
func analyzeMood(text string) string {
	text = strings.ToLower(text)

	// Positive mood indicators
	positiveWords := []string{
		"радостно", "весело", "прекрасно", "замечательно", "чудесно",
		"восхитительно", "потрясающе", "изумительно", "великолепно",
		"блестяще", "превосходно", "идеально", "совершенно",
		"прекрасный день", "замечательный день", "чудесный день",
		"в восторге", "в восхищении", "в эйфории",
		"на седьмом небе", "на вершине счастья", "полон радости",
		"полна радости", "счастливый", "счастливая",
		"доволен", "довольна", "удовлетворен", "удовлетворена",
		"в хорошем настроении", "в отличном настроении",
		"в прекрасном настроении", "в чудесном настроении",
		"в восхитительном настроении", "в потрясающем настроении",
		"хорошо", "хорошая", "хороший", "хорошее",
	}

	// Negative mood indicators
	negativeWords := []string{
		"грустно", "печально", "тоскливо", "мрачно", "уныло",
		"депрессивно", "подавленно", "разбито", "разбита",
		"опустошен", "опустошена", "разочарован", "разочарована",
		"в отчаянии", "в унынии", "в депрессии",
		"в плохом настроении", "в ужасном настроении",
		"в отвратительном настроении", "в мерзком настроении",
		"в паршивом настроении", "в скверном настроении",
		"в дурном настроении", "в гадком настроении",
		"в мерзопакостном настроении", "в отвратном настроении",
		"в ужасном состоянии", "в плохом состоянии",
		"в отвратительном состоянии", "в мерзком состоянии",
		"в паршивом состоянии", "в скверном состоянии",
		"в дурном состоянии", "в гадком состоянии",
		"в мерзопакостном состоянии", "в отвратном состоянии",
	}

	// Tired state indicators
	tiredWords := []string{
		"устал", "устала", "утомлен", "утомлена",
		"нет сил", "нет энергии", "упадок сил",
		"хочу спать", "сонный", "сонная",
		"вымотан", "вымотана", "измотан", "измотана",
		"нет настроения", "усталость", "утомление",
		"хочу отдохнуть", "нужен отдых", "нужен сон",
		"изнурен", "изнурена", "истощен", "истощена",
		"нет бодрости", "вялый", "вялая",
	}

	// Energized state indicators
	energizedWords := []string{
		"энергичн", "бодр", "бодра", "полон сил", "полна сил",
		"готов", "готова", "все могу", "все смогу",
		"отличное настроение", "прекрасное настроение",
		"полон энергии", "полна энергии", "много энергии",
		"активен", "активна", "бодрость", "энергия",
		"готов к работе", "готова к работе",
		"все по плечу", "все под силу",
		"отличное самочувствие", "прекрасное самочувствие",
		"полон энтузиазма", "полна энтузиазма",
	}

	// Check for energized words first
	for _, word := range energizedWords {
		if strings.Contains(text, word) {
			return "energized"
		}
	}

	// Then check for tired words
	for _, word := range tiredWords {
		if strings.Contains(text, word) {
			return "tired"
		}
	}

	// Then check for positive words
	for _, word := range positiveWords {
		if strings.Contains(text, word) {
			return "positive"
		}
	}

	// Finally check for negative words
	for _, word := range negativeWords {
		if strings.Contains(text, word) {
			return "negative"
		}
	}

	return "neutral"
}

func (b *Bot) Run() error {
	log.Printf("Authorized on account %s", b.api.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := b.api.GetUpdatesChan(u)

	for update := range updates {
		// Handle callback queries (button presses)
		if update.CallbackQuery != nil {
			callback := update.CallbackQuery
			chatID := callback.Message.Chat.ID

			var response string
			switch callback.Data {
			case "exercise1":
				response = "Упражнение 1: Глубокое дыхание\n\n" +
					"1. Сядьте удобно и расслабьтесь\n" +
					"2. Сделайте глубокий вдох через нос на 4 счета\n" +
					"3. Задержите дыхание на 4 счета\n" +
					"4. Медленно выдохните через рот на 4 счета\n" +
					"5. Повторите 5-7 раз\n\n" +
					"Это упражнение поможет снять напряжение и восстановить энергию."
			case "exercise2":
				response = "Упражнение 2: Растяжка шеи\n\n" +
					"1. Сядьте прямо\n" +
					"2. Медленно наклоните голову вправо, задержитесь на 10 секунд\n" +
					"3. Вернитесь в исходное положение\n" +
					"4. Повторите влево\n" +
					"5. Сделайте по 3-4 раза в каждую сторону\n\n" +
					"Это упражнение поможет снять напряжение в шее и плечах."
			case "exercise3":
				response = "Упражнение 3: Мини-прогулка\n\n" +
					"1. Встаньте и пройдитесь по комнате 2-3 минуты\n" +
					"2. Делайте это в спокойном темпе\n" +
					"3. Следите за дыханием\n" +
					"4. Можно выйти на свежий воздух, если есть возможность\n\n" +
					"Это упражнение поможет разогнать кровь и взбодриться."
			case "exercise4":
				response = "Упражнение 4: Гимнастика для глаз\n\n" +
					"1. Закройте глаза на 10 секунд\n" +
					"2. Откройте и посмотрите вдаль 10 секунд\n" +
					"3. Сделайте круговые движения глазами по часовой стрелке\n" +
					"4. Повторите против часовой стрелки\n" +
					"5. Сделайте 3-4 подхода\n\n" +
					"Это упражнение поможет снять напряжение с глаз и улучшить концентрацию."
			}

			msg := tgbotapi.NewMessage(chatID, response)
			if _, err := b.api.Send(msg); err != nil {
				log.Printf("Error sending message: %v", err)
			}

			// Answer callback query to remove loading state
			callbackConfig := tgbotapi.NewCallback(callback.ID, "")
			if _, err := b.api.Request(callbackConfig); err != nil {
				log.Printf("Error answering callback query: %v", err)
			}

			continue
		}

		if update.Message == nil {
			continue
		}

		chatID := update.Message.Chat.ID
		state := b.conversationStates[chatID]

		// Handle voice messages
		if update.Message.Voice != nil {
			log.Printf("Received voice message from user %d", chatID)
			
			// Download the voice message
			voice := update.Message.Voice
			file, err := b.api.GetFile(tgbotapi.FileConfig{FileID: voice.FileID})
			if err != nil {
				log.Printf("Error getting file: %v", err)
				continue
			}
			log.Printf("Got file info: %+v", file)

			// Create temp directory if it doesn't exist
			if err := os.MkdirAll("temp", 0755); err != nil {
				log.Printf("Error creating temp directory: %v", err)
				continue
			}

			// Download the file
			resp, err := http.Get(file.Link(b.api.Token))
			if err != nil {
				log.Printf("Error downloading file: %v", err)
				continue
			}
			defer resp.Body.Close()
			log.Printf("Downloaded file successfully")

			// Save the file
			audioPath := filepath.Join("temp", fmt.Sprintf("%s.ogg", voice.FileID))
			out, err := os.Create(audioPath)
			if err != nil {
				log.Printf("Error creating file: %v", err)
				continue
			}
			defer out.Close()

			_, err = io.Copy(out, resp.Body)
			if err != nil {
				log.Printf("Error saving file: %v", err)
				continue
			}
			log.Printf("Saved file to %s", audioPath)

			// Log file info
			fileInfo, err := os.Stat(audioPath)
			if err != nil {
				log.Printf("Error getting file info: %v", err)
			} else {
				log.Printf("Audio file size: %d bytes", fileInfo.Size())
			}

			// Convert OGG to WAV
			wavPath, err := speech.ConvertOggToWav(audioPath)
			if err != nil {
				log.Printf("Error converting audio: %v", err)
				continue
			}
			defer speech.CleanupAudioFiles(audioPath, wavPath)
			log.Printf("Converted to WAV: %s", wavPath)

			// Transcribe the audio
			text, err := b.speechClient.TranscribeAudio(wavPath)
			if err != nil {
				log.Printf("Error transcribing audio: %v", err)
				msg := tgbotapi.NewMessage(chatID, "Извините, не удалось распознать голосовое сообщение.")
				if _, err := b.api.Send(msg); err != nil {
					log.Printf("Error sending message: %v", err)
				}
				continue
			}
			log.Printf("Transcribed text: %s", text)

			// Send the transcribed text back to the user
			msg := tgbotapi.NewMessage(chatID, fmt.Sprintf("Распознанный текст: %s", text))
			if _, err := b.api.Send(msg); err != nil {
				log.Printf("Error sending transcribed text: %v", err)
			}

			// Process the transcribed text as if it was a text message
			text = strings.ToLower(text)
			log.Printf("Processing mood for text: %s", text)
			
			// Анализируем настроение сразу после получения голосового сообщения
			mood := analyzeMood(text)
			log.Printf("Detected mood: %s", mood)
			var response string

			switch mood {
			case "energized":
				response = "Отлично! 💪 Такая энергия - это здорово! Держи этот настрой и используй его для достижения своих целей!"
			case "tired":
				response = "Сожалею, что ты сейчас устал. Давай я предложу тебе 4 упражнения, которые помогут восстановиться."
				// Create keyboard with exercise buttons
				var keyboard = tgbotapi.NewInlineKeyboardMarkup(
					tgbotapi.NewInlineKeyboardRow(
						tgbotapi.NewInlineKeyboardButtonData("Упражнение 1", "exercise1"),
						tgbotapi.NewInlineKeyboardButtonData("Упражнение 2", "exercise2"),
					),
					tgbotapi.NewInlineKeyboardRow(
						tgbotapi.NewInlineKeyboardButtonData("Упражнение 3", "exercise3"),
						tgbotapi.NewInlineKeyboardButtonData("Упражнение 4", "exercise4"),
					),
				)
				msg := tgbotapi.NewMessage(chatID, response)
				msg.ReplyMarkup = keyboard
				if _, err := b.api.Send(msg); err != nil {
					log.Printf("Error sending tired response: %v", err)
				}
				continue
			case "positive":
				response = "Рад слышать, что у тебя всё хорошо! 😊 Давай сохраним это настроение!"
			case "negative":
				response = "Мне жаль, что тебе сейчас нелегко. Давай я предложу тебе несколько упражнений, которые помогут улучшить настроение."
				// Create keyboard with exercise buttons
				var keyboard = tgbotapi.NewInlineKeyboardMarkup(
					tgbotapi.NewInlineKeyboardRow(
						tgbotapi.NewInlineKeyboardButtonData("Упражнение 1", "exercise1"),
						tgbotapi.NewInlineKeyboardButtonData("Упражнение 2", "exercise2"),
					),
					tgbotapi.NewInlineKeyboardRow(
						tgbotapi.NewInlineKeyboardButtonData("Упражнение 3", "exercise3"),
						tgbotapi.NewInlineKeyboardButtonData("Упражнение 4", "exercise4"),
					),
				)
				msg := tgbotapi.NewMessage(chatID, response)
				msg.ReplyMarkup = keyboard
				if _, err := b.api.Send(msg); err != nil {
					log.Printf("Error sending negative response: %v", err)
				}
				continue
			default:
				response = "Понятно. Как я могу тебе помочь?"
			}

			msg = tgbotapi.NewMessage(chatID, response)
			if _, err := b.api.Send(msg); err != nil {
				log.Printf("Error sending final response: %v", err)
			}
			continue
		}

		// Handle text messages
		if !update.Message.IsCommand() {
			text := strings.ToLower(update.Message.Text)

			switch {
			case strings.Contains(text, "привет") && state == "":
				// Send greeting
				msg := tgbotapi.NewMessage(chatID, "Привет! 👋")
				if _, err := b.api.Send(msg); err != nil {
					log.Printf("Error sending message: %v", err)
				}

				// Ask how are you
				howAreYouMsg := tgbotapi.NewMessage(chatID, "Как ты сейчас?")
				if _, err := b.api.Send(howAreYouMsg); err != nil {
					log.Printf("Error sending message: %v", err)
				}

				b.conversationStates[chatID] = "waiting_for_mood"

			case state == "waiting_for_mood":
				mood := analyzeMood(text)
				var response string

				switch mood {
				case "energized":
					response = "Отлично! 💪 Такая энергия - это здорово! Держи этот настрой и используй его для достижения своих целей!"
					msg := tgbotapi.NewMessage(chatID, response)
					if _, err := b.api.Send(msg); err != nil {
						log.Printf("Error sending message: %v", err)
					}
					delete(b.conversationStates, chatID)

				case "tired":
					response = "Сожалею, что ты сейчас устал. Давай я предложу тебе 4 упражнения, которые помогут восстановиться."
					msg := tgbotapi.NewMessage(chatID, response)

					// Create keyboard with exercise buttons
					var keyboard = tgbotapi.NewInlineKeyboardMarkup(
						tgbotapi.NewInlineKeyboardRow(
							tgbotapi.NewInlineKeyboardButtonData("Упражнение 1", "exercise1"),
							tgbotapi.NewInlineKeyboardButtonData("Упражнение 2", "exercise2"),
						),
						tgbotapi.NewInlineKeyboardRow(
							tgbotapi.NewInlineKeyboardButtonData("Упражнение 3", "exercise3"),
							tgbotapi.NewInlineKeyboardButtonData("Упражнение 4", "exercise4"),
						),
					)

					msg.ReplyMarkup = keyboard
					if _, err := b.api.Send(msg); err != nil {
						log.Printf("Error sending message: %v", err)
					}
					b.conversationStates[chatID] = "waiting_for_exercise"

				case "positive":
					response = "Рад слышать, что у тебя всё хорошо! 😊 Давай сохраним это настроение!"
					msg := tgbotapi.NewMessage(chatID, response)
					if _, err := b.api.Send(msg); err != nil {
						log.Printf("Error sending message: %v", err)
					}
					delete(b.conversationStates, chatID)

				case "negative":
					response = "Понимаю, что сейчас не лучший момент. 🌟 Надеюсь, скоро всё наладится! Может, стоит сделать что-то приятное для себя?"
					msg := tgbotapi.NewMessage(chatID, response)
					if _, err := b.api.Send(msg); err != nil {
						log.Printf("Error sending message: %v", err)
					}
					delete(b.conversationStates, chatID)

				default:
					response = "Спасибо за ответ! Надеюсь, у тебя будет хороший день! 🌞"
					msg := tgbotapi.NewMessage(chatID, response)
					if _, err := b.api.Send(msg); err != nil {
						log.Printf("Error sending message: %v", err)
					}
					delete(b.conversationStates, chatID)
				}
			}
		}
	}

	return nil
}
