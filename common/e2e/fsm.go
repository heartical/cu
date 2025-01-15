package e2e

import (
	"context"
	"fmt"
	"log"

	"github.com/looplab/fsm"
)

// E2EE управляет процессом E2EE (End-to-End Encryption) с использованием конечного автомата (FSM).
func E2EE(ctx context.Context) {
	client, err := initClient()
	if err != nil {
		log.Fatalf("Ошибка при инициализации клиента: %v", err)
	}

	// Запуск FSM
	if err := createFSM(client, ctx).Event(ctx, "init"); err != nil {
		log.Fatalf("Ошибка при инициализации FSM: %v", err)
	}
}

// initClient инициализирует клиент E2EE.
func initClient() (*Client, error) {
	client, err := NewClient("http://localhost:8010")
	if err != nil {
		return nil, fmt.Errorf("не удалось создать клиент: %w", err)
	}
	return client, nil
}

// createFSM создает и настраивает конечный автомат (FSM) для управления процессом E2EE.
func createFSM(client *Client, ctx context.Context) *fsm.FSM {
	return fsm.NewFSM(
		"initial", // Начальное состояние
		fsm.Events{
			{Name: "init", Src: []string{"initial"}, Dst: "session_loaded"},                                      // Инициализация клиента
			{Name: "exchange_keys", Src: []string{"session_loaded"}, Dst: "keys_exchanged"},                      // Обмен ключами
			{Name: "ready", Src: []string{"keys_exchanged", "session_loaded"}, Dst: "ready"},                     // Готов к отправке сообщений
			{Name: "send_message", Src: []string{"ready"}, Dst: "message_sent"},                                  // Отправка сообщения
			{Name: "error", Src: []string{"initial", "session_loaded", "keys_exchanged", "ready"}, Dst: "error"}, // Обработка ошибок
		},
		fsm.Callbacks{
			// "enter_state": func(_ context.Context, e *fsm.Event) {
			// 	logStateTransition(e.FSM.Current())
			// },
			"init": func(_ context.Context, e *fsm.Event) {
				handleInit(client, e, ctx)
			},
			"exchange_keys": func(_ context.Context, e *fsm.Event) {
				handleExchangeKeys(client, e, ctx)
			},
			"ready": func(_ context.Context, e *fsm.Event) {
				fmt.Println("Клиент готов к отправке сообщений")
				if err := e.FSM.Event(ctx, "send_message"); err != nil {
					handleFSMError(e, ctx, "Ошибка при отправке сообщения", err)
				}
			},
			"send_message": func(_ context.Context, e *fsm.Event) {
				handleSendMessage(client, e, ctx)
			},
			"error": func(_ context.Context, e *fsm.Event) {
				handleError(e.FSM.Current())
			},
		},
	)
}

// logStateTransition логирует переход в новое состояние.
// func logStateTransition(state string) {
// 	fmt.Printf("Переход в состояние: %s\n", state)
// }

// handleInit обрабатывает инициализацию клиента.
func handleInit(client *Client, e *fsm.Event, ctx context.Context) {
	// fmt.Println("Инициализация клиента...")
	if !client.LoadSessionFromLocalStorage() {
		// fmt.Println("Сессия не найдена, требуется обмен ключами")
		if err := e.FSM.Event(ctx, "exchange_keys"); err != nil {
			handleFSMError(e, ctx, "Ошибка при обмене ключами", err)
		}
	} else {
		// fmt.Println("Сессия загружена из локального хранилища")
		// После загрузки сессии переходим в состояние ready
		if err := e.FSM.Event(ctx, "ready"); err != nil {
			handleFSMError(e, ctx, "Ошибка при переходе в состояние ready", err)
		}
	}
}

// handleExchangeKeys обрабатывает обмен ключами с сервером.
func handleExchangeKeys(client *Client, e *fsm.Event, ctx context.Context) {
	// fmt.Println("Обмен ключами с сервером...")
	if err := client.ExchangeKeysWithServer(); err != nil {
		handleFSMError(e, ctx, "Ошибка при обмене ключами", err)
	} else {
		// fmt.Println("Ключи успешно обменяны")
		// После обмена ключами переходим в состояние ready
		if err := e.FSM.Event(ctx, "ready"); err != nil {
			handleFSMError(e, ctx, "Ошибка при переходе в состояние ready", err)
		}
	}
}

// handleSendMessage обрабатывает отправку сообщения на сервер.
func handleSendMessage(client *Client, e *fsm.Event, ctx context.Context) {
	// fmt.Println("Отправка сообщения на сервер...")
	response, err := client.SendMessageToServer("ping")
	if err != nil {
		handleFSMError(e, ctx, "Ошибка при отправке сообщения", err)
	} else {
		fmt.Println("Ответ от сервера:", response)
	}
}

// handleError обрабатывает ошибки и логирует текущее состояние.
func handleError(currentState string) {
	fmt.Printf("Произошла ошибка. Текущее состояние: %s\n", currentState)
}

// handleFSMError обрабатывает ошибки FSM и переводит автомат в состояние ошибки.
func handleFSMError(e *fsm.Event, ctx context.Context, message string, err error) {
	fmt.Printf("%s: %v\n", message, err)
	if err := e.FSM.Event(ctx, "error"); err != nil {
		log.Fatalf("Ошибка при переходе в состояние ошибки: %v", err)
	}
}
