package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sort"
	"time"

	"github.com/gorilla/mux"
	redisInfra "demo/backend-go/internal/infrastructure/redis"
	"github.com/redis/go-redis/v9"
)

type BotStatsHandler struct {
	redisService *redisInfra.RedisService
}

func NewBotStatsHandler(redisService *redisInfra.RedisService) *BotStatsHandler {
	return &BotStatsHandler{
		redisService: redisService,
	}
}

func (h *BotStatsHandler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/api/bot/stats", h.GetBotStats).Methods("GET")
	router.HandleFunc("/api/bot/stats/reset", h.ResetDailyStats).Methods("POST")
}

type BotStatsResponse struct {
	Commands       CommandStats `json:"commands"`
	Messages       MessageStats `json:"messages"`
	ActiveChatters int          `json:"activeChatters"`
	TopCommands    []TopCommand `json:"topCommands,omitempty"`
}

type CommandStats struct {
	Total   int `json:"total"`
	Today   int `json:"today"`
	Success int `json:"success"`
	Failed  int `json:"failed"`
}

type MessageStats struct {
	Total int `json:"total"`
	Today int `json:"today"`
}

type TopCommand struct {
	Name  string `json:"name"`
	Count int    `json:"count"`
}

func (h *BotStatsHandler) GetBotStats(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	if h.redisService == nil {
		log.Printf("❌ Redis Service nicht verfügbar")
		http.Error(w, "Service nicht verfügbar", http.StatusServiceUnavailable)
		return
	}

	stats, err := h.getBotStatisticsFromRedis(ctx)
	if err != nil {
		log.Printf("❌ Fehler beim Laden der Stats: %v", err)
		stats = &BotStatsResponse{
			Commands:       CommandStats{},
			Messages:       MessageStats{},
			ActiveChatters: 0,
			TopCommands:    []TopCommand{},
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats)
}

func (h *BotStatsHandler) getBotStatisticsFromRedis(ctx context.Context) (*BotStatsResponse, error) {
	client := h.redisService.GetClient()
	if client == nil {
		return nil, nil
	}

	// ✅ Dynamisches Datum
	today := time.Now().Format("2006-01-02")

	pipe := client.Pipeline()

	cmdTotal := pipe.Get(ctx, "bot:stats:commands:total")
	cmdToday := pipe.Get(ctx, fmt.Sprintf("bot:stats:commands:today:%s", today))
	cmdSuccess := pipe.Get(ctx, "bot:stats:commands:success")
	cmdFailed := pipe.Get(ctx, "bot:stats:commands:failed")

	msgTotal := pipe.Get(ctx, "bot:stats:messages:total")
	msgToday := pipe.Get(ctx, fmt.Sprintf("bot:stats:messages:today:%s", today))

	statsData := pipe.Get(ctx, "bot:stats:data")

	topCmds := pipe.HGetAll(ctx, "bot:stats:commands:by_name")

	_, err := pipe.Exec(ctx)
	// ✅ FIX: redis.Nil ist OK (bedeutet Keys existieren noch nicht)
	if err != nil && err != redis.Nil {
		log.Printf("⚠️ Redis Pipeline Fehler: %v", err)
		return nil, err
	}

	response := &BotStatsResponse{
		Commands: CommandStats{
			Total:   parseInt(cmdTotal.Val()),
			Today:   parseInt(cmdToday.Val()),
			Success: parseInt(cmdSuccess.Val()),
			Failed:  parseInt(cmdFailed.Val()),
		},
		Messages: MessageStats{
			Total: parseInt(msgTotal.Val()),
			Today: parseInt(msgToday.Val()),
		},
		ActiveChatters: 0,
		TopCommands:    parseTopCommands(topCmds.Val()),
	}

	if statsDataVal := statsData.Val(); statsDataVal != "" {
		var statsJson struct {
			ActiveChatters int `json:"active_chatters"`
		}
		if err := json.Unmarshal([]byte(statsDataVal), &statsJson); err == nil {
			response.ActiveChatters = statsJson.ActiveChatters
		}
	}

	log.Printf("📊 Stats geladen: Commands=%d, Messages=%d, Chatters=%d",
		response.Commands.Total, response.Messages.Total, response.ActiveChatters)

	return response, nil
}

func (h *BotStatsHandler) ResetDailyStats(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	if h.redisService == nil {
		http.Error(w, "Service nicht verfügbar", http.StatusServiceUnavailable)
		return
	}

	client := h.redisService.GetClient()
	if client == nil {
		http.Error(w, "Redis nicht verbunden", http.StatusServiceUnavailable)
		return
	}

	today := time.Now().Format("2006-01-02")

	err := client.Del(ctx,
		"bot:stats:commands:today:"+today,
		"bot:stats:messages:today:"+today,
	).Err()

	if err != nil {
		log.Printf("❌ Fehler beim Reset: %v", err)
		http.Error(w, "Fehler beim Reset", http.StatusInternalServerError)
		return
	}

	log.Printf("✅ Daily Stats zurückgesetzt")

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]bool{"success": true})
}

func parseInt(s string) int {
	var val int
	if s != "" {
		_, _ = fmt.Sscanf(s, "%d", &val)
	}
	return val
}

func parseTopCommands(data map[string]string) []TopCommand {
	if len(data) == 0 {
		return []TopCommand{}
	}

	commands := make([]TopCommand, 0, len(data))
	for name, countStr := range data {
		count := parseInt(countStr)
		commands = append(commands, TopCommand{
			Name:  name,
			Count: count,
		})
	}

	sort.Slice(commands, func(i, j int) bool {
		return commands[i].Count > commands[j].Count
	})

	if len(commands) > 10 {
		commands = commands[:10]
	}

	return commands
}