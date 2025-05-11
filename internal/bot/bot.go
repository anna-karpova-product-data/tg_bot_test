package bot

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"tg_bot/internal/logger"
	"tg_bot/internal/speech"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Bot struct {
	api *tgbotapi.BotAPI
	// Map to store conversation states
	conversationStates map[int64]string
	// Speech recognition client
	speechClient *speech.DeepgramClient
	// Map to store mood recognition attempts
	moodAttempts map[int64]int
	// Logger
	logger *logger.Logger
}

func New(telegramToken, deepgramToken string) (*Bot, error) {
	api, err := tgbotapi.NewBotAPI(telegramToken)
	if err != nil {
		return nil, err
	}

	// Инициализируем логгер
	logger, err := logger.New("logs/bot.log")
	if err != nil {
		return nil, fmt.Errorf("failed to initialize logger: %v", err)
	}

	return &Bot{
		api:                api,
		conversationStates: make(map[int64]string),
		speechClient:       speech.NewDeepgramClient(deepgramToken),
		moodAttempts:       make(map[int64]int),
		logger:             logger,
	}, nil
}

// analyzeMood analyzes the text and returns the detected mood
func analyzeMood(text string) string {
	text = strings.ToLower(text)

	// Positive mood indicators
	positiveWords := []string{
		// Базовые положительные состояния
		"кайф", "охуенн", "заеб", "пиздат", "огонь", "ахуенн", "волшебн", "балдеж", "душевн", "чум", "кайфец", "кайфушк", "сладк", "красот", "тепл", "милот", "лампов", "трепетн", "пушечн", "праздник",

		// Эмоциональные реакции
		"раду", "мурашк", "приятн", "трогательн", "крут", "слез", "красив", "классн", "спокойн", "глубин", "прослез", "щем", "счаст", "любл", "обожа", "сердечк", "зашл", "тема",

		// Усилители и сравнения
		"как", "будто", "словно", "точно", "прям", "уж", "вот", "ну", "аж", "через", "край", "слож", "надо",

		// Базовые эмоции
		"радостн", "спокойн", "легк", "приятн", "тепл", "уютн", "светл", "хорош", "мягк", "вдохновл", "трогательн", "умиротворен", "благодарн", "довольн", "счаст", "восхищен", "нежн", "любов", "уверен", "забот", "интерес", "любопытн",

		// Глубокие состояния
		"полнот", "смысл", "волнен", "принят", "наслажден", "восторг", "удовлетворен", "гармони", "ясн", "открыт", "довер", "легк", "поко", "надежд", "искрен", "целост", "благ", "благополуч", "признательн", "очарован",

		// Физические ощущения
		"тепл", "свет", "обня", "сердц", "поет", "внутр", "место",

		// Действия и состояния
		"улыба", "получил", "чувству", "дума", "тронул", "произошл", "доволен", "довольн",

		// Существующие слова
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
		// Базовые негативные состояния
		"тосклив", "тревожн", "пуст", "обидн", "тяжел", "больн", "одинок", "горьк", "несправедлив", "страшн", "неловк", "стыдн", "злост", "безысходн", "уныл", "мучительн", "раздража", "разочарован", "нудн", "мерзк", "мерзост", "отвращен", "тревог", "скук", "апати", "ненавиж", "отчаян", "беспомощн",

		// Матерные и разговорные выражения
		"хуев", "паршив", "дерьмов", "говен", "бес", "жоп", "пизд", "еба", "надоел", "чертик", "ад", "бляд", "хренов", "хуйн", "сран", "сук", "хуяр", "херн", "черт", "жоп", "больн", "херов", "нахуй", "надежд", "выт", "скреб", "сдох", "заеб",

		// Эмоциональные состояния
		"понима", "раздража", "не так", "успоко", "плака", "застря", "дело", "лишн", "невыносим", "хоч", "испорт", "отпуст", "почему", "валит", "смысл", "дыр", "сер", "раду", "говор",

		// Отрицания и усилители
		"не", "ни", "вс", "как", "будто", "словно", "точно", "опят", "внутр", "ничего", "никому", "ни с кем",

		// Существующие слова
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
		// Базовые состояния
		"устал", "устал", "вымотан", "выжат", "опустошен", "изможден", "разбит", "истощен", "перегруз", "перегор", "сонн", "мутн", "напряжен", "предел",

		// Физические ощущения
		"ватн", "голов", "тяжел", "шум", "плыв", "туп", "засыпа", "перегрев", "замедл", "туман", "тело", "диван", "леж", "стен", "поезд", "навалил", "тян", "одеял",

		// Ментальные состояния
		"сообража", "вар", "мозг", "ресурс", "сил", "автопилот", "зомб", "провал", "существу", "ком", "тряпк", "говн", "лошад", "паш",

		// Эмоциональные состояния
		"выгоран", "нетерпим", "эмоциональн", "нахуй", "заеб", "задолб", "вымота", "еба", "пиздец", "сдох", "бляд", "охует", "говн", "говор",

		// Отрицания и усилители
		"не", "нет", "никак", "больш", "последн", "вс", "просто", "как", "будто", "хоть", "уже", "больш", "всё", "все", "ничего", "ничего", "никакой", "никакая",

		// Действия и состояния
		"лечь", "лежать", "исчез", "выспат", "кончит", "встават", "полз", "встава", "тян", "вар", "плыв", "провалива", "лез", "работа", "выжра", "высос", "заеба", "задолб", "вымота", "еба", "сдох", "говн", "говор",

		// Сравнения
		"как", "будто", "словно", "точно", "похож", "напомина", "подобн", "такой", "такая", "такое", "такие",

		// Существующие слова
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
		// Ментальная бодрость
		"ясн", "собран", "сконцентрирован", "сфокусирован", "внимательн", "включен", "волн", "соображаю", "поток", "четк", "структурн", "остр", "голов", "гибк", "мышлен",
		"мозг", "работает", "соображаю", "раз", "два", "решаю", "налету", "мысл", "ясн", "полочк", "проснул", "порядок",
		"схватываю", "лету", "башк", "варит", "голов", "тормозит", "врубаюсь", "полуслов", "соображаю", "мозг", "тупит", "фигачу", "шерлок",

		// Физическая бодрость
		"бодр", "легк", "свеж", "заряжен", "жив", "подвижн", "гибк", "пружин", "энерг", "прет", "ход", "активн", "летиш", "теле", "огонь",
		"легкост", "теле", "крыл", "выросл", "могу", "заряд", "полн", "двигаться", "усидеть", "тело", "радуется",
		"пр", "прет", "огурчик", "бегаю", "заведен", "хрен", "догониш", "ног", "несут", "ебашу", "спортзал", "качаю", "энерг", "охуенно", "теле", "пляшет", "бодрячком", "остановить", "заткнеш",

		// Эмоциональный подъём
		"ресурс", "вдохновлен", "стабильн", "радостн", "наполнен", "поток", "уверенн", "спокойн", "баланс", "душ", "цельн", "интерес",
		"делиться", "плечу", "хорошо", "нравится", "жив", "возможн", "вдохновляюсь", "процесс",
		"заебись", "кайфую", "жизн", "охуенно", "душ", "ебать", "прет", "добро", "летиш", "улыбаеш", "балдежн", "состояние", "хуярю", "удовольстви", "идет", "надо", "жизн", "огонь", "аплодирую", "светится", "позитив",

		// Существующие слова
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
	defer b.logger.Close()

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := b.api.GetUpdatesChan(u)

	for update := range updates {
		// Handle callback queries (button presses)
		if update.CallbackQuery != nil {
			callback := update.CallbackQuery
			chatID := callback.Message.Chat.ID
			username := callback.From.UserName
			if username == "" {
				username = fmt.Sprintf("User%d", chatID)
			}

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

			// Логируем ответ на callback
			if err := b.logger.Log(chatID, username, "callback", callback.Data, response, ""); err != nil {
				log.Printf("Error logging callback: %v", err)
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
		username := update.Message.From.UserName
		if username == "" {
			username = fmt.Sprintf("User%d", chatID)
		}
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

			// Process the transcribed text as if it was a text message
			text = strings.ToLower(text)
			log.Printf("Processing mood for text: %s", text)

			// Анализируем настроение сразу после получения голосового сообщения
			mood := analyzeMood(text)
			log.Printf("Detected mood: %s", mood)
			var response string
			var msg tgbotapi.MessageConfig

			// Увеличиваем счетчик попыток
			b.moodAttempts[chatID]++
			attempts := b.moodAttempts[chatID]

			// Если после 3 попыток не удалось определить настроение, считаем его нейтральным
			if mood == "neutral" && attempts >= 3 {
				mood = "neutral_final"
			}

			switch mood {
			case "energized":
				response = "Отлично! 💪 Такая энергия - это здорово! Держи этот настрой и используй его для достижения своих целей!"
				delete(b.moodAttempts, chatID)
			case "tired":
				response = "Сожалею, что ты сейчас устал. Давай я предложу тебе 4 упражнения, которые помогут восстановиться."
				delete(b.moodAttempts, chatID)
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
				msg = tgbotapi.NewMessage(chatID, response)
				msg.ReplyMarkup = keyboard
				if _, err := b.api.Send(msg); err != nil {
					log.Printf("Error sending tired response: %v", err)
				}
				continue
			case "positive":
				response = "Рад слышать, что у тебя всё хорошо! 😊 Давай сохраним это настроение!"
				delete(b.moodAttempts, chatID)
			case "negative":
				response = "Мне жаль, что тебе сейчас нелегко. Давай я предложу тебе несколько упражнений, которые помогут улучшить настроение."
				delete(b.moodAttempts, chatID)
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
				msg = tgbotapi.NewMessage(chatID, response)
				msg.ReplyMarkup = keyboard
				if _, err := b.api.Send(msg); err != nil {
					log.Printf("Error sending negative response: %v", err)
				}
				continue
			case "neutral":
				response = "Расскажи мне побольше."
				msg = tgbotapi.NewMessage(chatID, response)
				if _, err := b.api.Send(msg); err != nil {
					log.Printf("Error sending neutral response: %v", err)
				}
				continue
			case "neutral_final":
				response = "Спасибо за ответ! Надеюсь, у тебя будет хороший день! 🌞"
				delete(b.moodAttempts, chatID)
			default:
				response = "Расскажи мне побольше."
			}

			msg = tgbotapi.NewMessage(chatID, response)
			if _, err := b.api.Send(msg); err != nil {
				log.Printf("Error sending final response: %v", err)
			}

			// Логируем голосовое сообщение и ответ
			if err := b.logger.Log(chatID, username, "voice", text, response, mood); err != nil {
				log.Printf("Error logging voice message: %v", err)
			}

			continue
		}

		// Handle text messages
		if !update.Message.IsCommand() {
			text := strings.ToLower(update.Message.Text)

			switch {
			case strings.Contains(text, "привет") && state == "":
				// Reset mood attempts counter
				delete(b.moodAttempts, chatID)

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

				// Логируем приветствие
				if err := b.logger.Log(chatID, username, "text", text, "Привет! 👋\nКак ты сейчас?", ""); err != nil {
					log.Printf("Error logging greeting: %v", err)
				}

				b.conversationStates[chatID] = "waiting_for_mood"

			case state == "waiting_for_mood":
				mood := analyzeMood(text)
				var response string

				// Увеличиваем счетчик попыток
				b.moodAttempts[chatID]++
				attempts := b.moodAttempts[chatID]

				// Если после 3 попыток не удалось определить настроение, считаем его нейтральным
				if mood == "neutral" && attempts >= 3 {
					mood = "neutral_final"
				}

				switch mood {
				case "energized":
					response = "Отлично! 💪 Такая энергия - это здорово! Держи этот настрой и используй его для достижения своих целей!"
					msg := tgbotapi.NewMessage(chatID, response)
					if _, err := b.api.Send(msg); err != nil {
						log.Printf("Error sending message: %v", err)
					}
					delete(b.conversationStates, chatID)
					delete(b.moodAttempts, chatID)

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
					delete(b.moodAttempts, chatID)

				case "positive":
					response = "Рад слышать, что у тебя всё хорошо! 😊 Давай сохраним это настроение!"
					msg := tgbotapi.NewMessage(chatID, response)
					if _, err := b.api.Send(msg); err != nil {
						log.Printf("Error sending message: %v", err)
					}
					delete(b.conversationStates, chatID)
					delete(b.moodAttempts, chatID)

				case "negative":
					response = "Понимаю, что сейчас не лучший момент. 🌟 Надеюсь, скоро всё наладится! Может, стоит сделать что-то приятное для себя?"
					msg := tgbotapi.NewMessage(chatID, response)
					if _, err := b.api.Send(msg); err != nil {
						log.Printf("Error sending message: %v", err)
					}
					delete(b.conversationStates, chatID)
					delete(b.moodAttempts, chatID)

				case "neutral":
					response = "Расскажи мне побольше."
					msg := tgbotapi.NewMessage(chatID, response)
					if _, err := b.api.Send(msg); err != nil {
						log.Printf("Error sending message: %v", err)
					}

				case "neutral_final":
					response = "Спасибо за ответ! Надеюсь, у тебя будет хороший день! 🌞"
					msg := tgbotapi.NewMessage(chatID, response)
					if _, err := b.api.Send(msg); err != nil {
						log.Printf("Error sending message: %v", err)
					}
					delete(b.conversationStates, chatID)
					delete(b.moodAttempts, chatID)

				default:
					response = "Расскажи мне побольше."
					msg := tgbotapi.NewMessage(chatID, response)
					if _, err := b.api.Send(msg); err != nil {
						log.Printf("Error sending message: %v", err)
					}
				}

				// Логируем текстовое сообщение и ответ
				if err := b.logger.Log(chatID, username, "text", text, response, mood); err != nil {
					log.Printf("Error logging text message: %v", err)
				}
			}
		}
	}

	return nil
}
