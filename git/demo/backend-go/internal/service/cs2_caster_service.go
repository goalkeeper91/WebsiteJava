package service

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"regexp"
	"sort"
	"strings"
	"sync"

	"demo/backend-go/internal/domain"
	"demo/backend-go/internal/infrastructure/redis"
	"demo/backend-go/internal/repository"
	"demo/backend-go/internal/twitch"
)

// cs2ScoreTagPattern matches a score tag this service previously inserted
// into a stream title (e.g. " [3:2]"), so it can be stripped before inserting
// an updated one - keeps the streamer's own title intact otherwise. Not
// end-anchored since the tag sits right after the matchup segment, not
// necessarily at the very end of the title.
var cs2ScoreTagPattern = regexp.MustCompile(`\s\[\d+:\d+\]`)

// twitchTitleMaxLength is Twitch's documented character limit for a channel title.
const twitchTitleMaxLength = 140

// cs2TitleSegmentSeparator is the convention this feature assumes for
// splitting a title into "matchup | rest" (e.g.
// "[GER] TeamA vs. TeamB | Tournament Info | #hashtag") - the score tag is
// inserted right after the matchup segment, before the first separator.
const cs2TitleSegmentSeparator = " | "

// insertScoreTag places scoreTag right after the matchup segment of a title
// (before the first " | "), or appends it at the end if the title doesn't
// follow that "matchup | rest" convention.
func insertScoreTag(baseTitle, scoreTag string) string {
	if idx := strings.Index(baseTitle, cs2TitleSegmentSeparator); idx != -1 {
		return baseTitle[:idx] + scoreTag + baseTitle[idx:]
	}
	return baseTitle + scoreTag
}

// CS2CasterService ingests CS2 Game State Integration payloads and turns
// specific state transitions (match start / multi-kill round / map end /
// score change) into Twitch Predictions, chat announcements and stream-title
// updates. Live match state is kept purely in-memory (like StreamIngestor) -
// nothing about a match is persisted once it ends.
type CS2CasterService struct {
	casterRepo    repository.CS2CasterRepository
	authService   *AuthService
	channelClient *twitch.ChannelClient
	redisService  *redis.RedisService

	mu      sync.Mutex
	matches map[string]*cs2MatchState // key: user_twitch_id
}

func NewCS2CasterService(
	casterRepo repository.CS2CasterRepository,
	authService *AuthService,
	channelClient *twitch.ChannelClient,
	redisService *redis.RedisService,
) *CS2CasterService {
	return &CS2CasterService{
		casterRepo:    casterRepo,
		authService:   authService,
		channelClient: channelClient,
		redisService:  redisService,
		matches:       map[string]*cs2MatchState{},
	}
}

// cs2MatchState tracks one channel's currently-live map. GSI ticks for a
// given channel arrive serially from a single CS2 client, so fields on an
// already-retrieved *cs2MatchState are never mutated concurrently - only
// creating/looking up the pointer in the matches map needs the mutex.
type cs2MatchState struct {
	mapPhase                string
	mapRound                int
	announcedKillsThisRound map[string]int // steamid -> highest kill count already announced this round

	predictionActive    bool
	predictionID        string
	outcomeIDByTeamName map[string]string // lowercased team name -> Twitch outcome id

	lastTitleScoreKey string // dedup marker so title updates only fire on an actual score change

	// For the dashboard's live-status panel / note highlighting.
	observedPlayerName string
	teamCTName         string
	teamTName          string
	scoreCT            int
	scoreT             int
	mapName            string
	rosterCT           []string // player names currently on the CT side, sorted
	rosterT            []string // player names currently on the T side, sorted
}

// --- GSI payload shapes -----------------------------------------------
//
// Subscription keys (provider/map/round/player_id/player_state/allplayers_id/
// allplayers_state/allplayers_match_stats) are confirmed against a real,
// working GSI config (Overwolf's own CS2 integration). The exact payload
// field names below (team_ct.name/.score, round_kills, ...) still reflect
// the best-available understanding of CS2's GSI JSON structure, cross-checked
// against multiple third-party GSI libraries - Valve's own wiki could not be
// fetched directly during research. Verify against a real captured payload
// if ingestion silently no-ops in practice.

type gsiPayload struct {
	Map        *gsiMap                 `json:"map"`
	Player     *gsiPlayer              `json:"player"`
	AllPlayers map[string]gsiAllPlayer `json:"allplayers"`
}

type gsiMap struct {
	Name   string  `json:"name"`
	Phase  string  `json:"phase"` // "warmup" | "live" | "gameover"
	Round  int     `json:"round"`
	TeamCT gsiTeam `json:"team_ct"`
	TeamT  gsiTeam `json:"team_t"`
}

type gsiTeam struct {
	Score int    `json:"score"`
	Name  string `json:"name"`
}

type gsiPlayer struct {
	SteamID string `json:"steamid"`
	Name    string `json:"name"`
}

type gsiAllPlayer struct {
	Name  string            `json:"name"`
	Team  string            `json:"team"` // "CT" | "T"
	State gsiAllPlayerState `json:"state"`
}

type gsiAllPlayerState struct {
	RoundKills int `json:"round_kills"`
}

const predictionWindowSeconds = 1800 // Twitch's own maximum - resolution happens on map-end, not on window expiry

// IngestGSIPayload is the single entry point the HTTP handler calls for
// every incoming GSI POST.
func (s *CS2CasterService) IngestGSIPayload(ctx context.Context, token string, raw []byte) error {
	settings, err := s.casterRepo.GetByGSIToken(ctx, token)
	if err != nil {
		return fmt.Errorf("fehler beim Auflösen des GSI-Tokens: %w", err)
	}
	if settings == nil {
		return domain.ErrCS2InvalidGSIToken
	}

	var payload gsiPayload
	if err := json.Unmarshal(raw, &payload); err != nil {
		return fmt.Errorf("fehler beim Parsen des GSI-Payloads: %w", err)
	}
	if payload.Map == nil {
		return nil // menu/warmup-only tick without an active map - nothing to do
	}

	userTwitchID := settings.UserTwitchID

	s.mu.Lock()
	state, ok := s.matches[userTwitchID]
	if !ok {
		state = &cs2MatchState{announcedKillsThisRound: map[string]int{}}
		s.matches[userTwitchID] = state
	}
	s.mu.Unlock()

	previousPhase := state.mapPhase
	state.mapPhase = payload.Map.Phase
	state.teamCTName = payload.Map.TeamCT.Name
	state.teamTName = payload.Map.TeamT.Name
	state.scoreCT = payload.Map.TeamCT.Score
	state.scoreT = payload.Map.TeamT.Score
	state.mapName = payload.Map.Name
	if payload.Player != nil {
		state.observedPlayerName = payload.Player.Name
	}
	updateRoster(state, payload.AllPlayers)

	if payload.Map.Round != state.mapRound {
		state.mapRound = payload.Map.Round
		state.announcedKillsThisRound = map[string]int{}
	}

	if previousPhase != "live" && payload.Map.Phase == "live" {
		s.handleMatchStart(ctx, userTwitchID, settings, state, payload.Map)
	}

	if settings.MultikillAnnounceEnabled {
		s.handleMultikills(userTwitchID, state, payload.AllPlayers)
	}

	if settings.TitleUpdateEnabled && payload.Map.Phase == "live" {
		s.handleTitleUpdate(ctx, userTwitchID, state, payload.Map)
	}

	if previousPhase != "gameover" && payload.Map.Phase == "gameover" {
		s.handleMapEnd(ctx, userTwitchID, settings, state, payload.Map)
	}

	return nil
}

// updateRoster rebuilds the per-side player-name lists from the current GSI
// tick's allplayers section - used by the dashboard's match-notes popup to
// match saved player notes against whoever is actually in the current match.
func updateRoster(state *cs2MatchState, allPlayers map[string]gsiAllPlayer) {
	var ct, t []string
	for _, p := range allPlayers {
		switch strings.ToUpper(p.Team) {
		case "CT":
			ct = append(ct, p.Name)
		case "T":
			t = append(t, p.Name)
		}
	}
	sort.Strings(ct)
	sort.Strings(t)
	state.rosterCT = ct
	state.rosterT = t
}

func (s *CS2CasterService) handleMatchStart(ctx context.Context, userTwitchID string, settings *domain.CS2CasterSettings, state *cs2MatchState, m *gsiMap) {
	state.announcedKillsThisRound = map[string]int{}
	state.predictionActive = false
	state.predictionID = ""
	state.outcomeIDByTeamName = nil
	state.lastTitleScoreKey = ""

	if !settings.PredictionsEnabled {
		return
	}

	teamAName := m.TeamCT.Name
	if teamAName == "" {
		teamAName = "CT-Team"
	}
	teamBName := m.TeamT.Name
	if teamBName == "" {
		teamBName = "T-Team"
	}

	accessToken, err := s.authService.GetFreshAccessToken(ctx, userTwitchID)
	if err != nil {
		log.Printf("⚠️ CS2: kein gültiger Token für Prediction-Erstellung (%s): %v", userTwitchID, err)
		return
	}

	predictionID, outcomeIDs, err := s.channelClient.CreatePrediction(ctx, userTwitchID, accessToken, "Wer gewinnt die Map?", []string{teamAName, teamBName}, predictionWindowSeconds)
	if err != nil {
		log.Printf("⚠️ CS2: Prediction konnte nicht erstellt werden (%s): %v", userTwitchID, err)
		return
	}
	if len(outcomeIDs) != 2 {
		log.Printf("⚠️ CS2: unerwartete Anzahl Outcomes von Twitch (%s): %d", userTwitchID, len(outcomeIDs))
		return
	}

	state.predictionID = predictionID
	state.predictionActive = true
	state.outcomeIDByTeamName = map[string]string{
		strings.ToLower(teamAName): outcomeIDs[0],
		strings.ToLower(teamBName): outcomeIDs[1],
	}
}

func (s *CS2CasterService) handleMultikills(userTwitchID string, state *cs2MatchState, allPlayers map[string]gsiAllPlayer) {
	for steamID, p := range allPlayers {
		kills := p.State.RoundKills
		if kills < 3 {
			continue
		}
		if kills <= state.announcedKillsThisRound[steamID] {
			continue
		}
		state.announcedKillsThisRound[steamID] = kills

		var msg string
		switch {
		case kills >= 5:
			msg = fmt.Sprintf("🏆 ACE! %s holt alle 5 Kills in dieser Runde!", p.Name)
		case kills == 4:
			msg = fmt.Sprintf("🔥🔥 4 Kills von %s in dieser Runde!", p.Name)
		default:
			msg = fmt.Sprintf("🔥 3 Kills von %s in dieser Runde!", p.Name)
		}
		if err := s.redisService.SendBotAnnounce(msg, userTwitchID); err != nil {
			log.Printf("⚠️ CS2: Multikill-Ankündigung fehlgeschlagen (%s): %v", userTwitchID, err)
		}
	}
}

func (s *CS2CasterService) handleTitleUpdate(ctx context.Context, userTwitchID string, state *cs2MatchState, m *gsiMap) {
	scoreKey := fmt.Sprintf("%s:%d-%d:%s", m.TeamCT.Name, m.TeamCT.Score, m.TeamT.Score, m.TeamT.Name)
	if scoreKey == state.lastTitleScoreKey {
		return
	}
	state.lastTitleScoreKey = scoreKey

	accessToken, err := s.authService.GetFreshAccessToken(ctx, userTwitchID)
	if err != nil {
		log.Printf("⚠️ CS2: kein gültiger Token für Titel-Update (%s): %v", userTwitchID, err)
		return
	}

	// Aktuellen Titel frisch abrufen (nicht zwischenspeichern) - so bleibt ein
	// Titel, den der Streamer selbst mitten im Match noch anpasst, respektiert.
	// Nur ein von uns selbst zuvor angehängter Score-Tag wird entfernt, bevor
	// der aktuelle Stand neu angehängt wird.
	info, err := s.channelClient.GetChannelInfo(ctx, userTwitchID, accessToken)
	if err != nil {
		log.Printf("⚠️ CS2: aktueller Titel konnte nicht geladen werden (%s): %v", userTwitchID, err)
		return
	}

	baseTitle := cs2ScoreTagPattern.ReplaceAllString(info.Title, "")
	scoreTag := fmt.Sprintf(" [%d:%d]", m.TeamCT.Score, m.TeamT.Score)
	title := insertScoreTag(baseTitle, scoreTag)

	// Titel zu lang für Twitch (140 Zeichen)? Am Ende kürzen - das lässt den
	// wichtigeren vorderen Teil (Matchup + Score) unangetastet.
	if len(title) > twitchTitleMaxLength {
		title = strings.TrimSpace(title[:twitchTitleMaxLength])
	}

	if title == info.Title {
		return
	}

	if err := s.channelClient.ModifyChannelInfo(ctx, userTwitchID, accessToken, &title, nil); err != nil {
		log.Printf("⚠️ CS2: Titel-Update fehlgeschlagen (%s): %v", userTwitchID, err)
	}
}

func (s *CS2CasterService) handleMapEnd(ctx context.Context, userTwitchID string, settings *domain.CS2CasterSettings, state *cs2MatchState, m *gsiMap) {
	winnerName := m.TeamT.Name
	if m.TeamCT.Score >= m.TeamT.Score {
		winnerName = m.TeamCT.Name
	}
	if winnerName == "" {
		winnerName = "Unbekanntes Team"
	}

	if settings.MapEndAnnounceEnabled {
		if err := s.redisService.SendBotAnnounce(fmt.Sprintf("GG WP %s!", winnerName), userTwitchID); err != nil {
			log.Printf("⚠️ CS2: GG-WP-Ankündigung fehlgeschlagen (%s): %v", userTwitchID, err)
		}
	}

	if state.predictionActive && settings.PredictionsEnabled {
		outcomeID, found := state.outcomeIDByTeamName[strings.ToLower(winnerName)]
		if found {
			accessToken, err := s.authService.GetFreshAccessToken(ctx, userTwitchID)
			if err != nil {
				log.Printf("⚠️ CS2: kein gültiger Token für Prediction-Auflösung (%s): %v", userTwitchID, err)
			} else if err := s.channelClient.EndPrediction(ctx, userTwitchID, accessToken, state.predictionID, "RESOLVED", &outcomeID); err != nil {
				log.Printf("⚠️ CS2: Prediction konnte nicht aufgelöst werden (%s): %v", userTwitchID, err)
			}
		}
		state.predictionActive = false
	}
}

// CS2LiveStatus is a read-only snapshot for the dashboard's live-status
// panel - nil/zero fields simply mean no GSI session is currently active.
type CS2LiveStatus struct {
	Active             bool     `json:"active"`
	ObservedPlayerName string   `json:"observed_player_name,omitempty"`
	TeamCTName         string   `json:"team_ct_name,omitempty"`
	TeamTName          string   `json:"team_t_name,omitempty"`
	ScoreCT            int      `json:"score_ct"`
	ScoreT             int      `json:"score_t"`
	MapName            string   `json:"map_name,omitempty"`
	TeamCTPlayers      []string `json:"team_ct_players,omitempty"`
	TeamTPlayers       []string `json:"team_t_players,omitempty"`
}

func (s *CS2CasterService) GetLiveStatus(userTwitchID string) *CS2LiveStatus {
	s.mu.Lock()
	state, ok := s.matches[userTwitchID]
	s.mu.Unlock()
	if !ok {
		return &CS2LiveStatus{Active: false}
	}
	return &CS2LiveStatus{
		Active:             state.mapPhase == "live",
		ObservedPlayerName: state.observedPlayerName,
		TeamCTName:         state.teamCTName,
		TeamTName:          state.teamTName,
		ScoreCT:            state.scoreCT,
		ScoreT:             state.scoreT,
		MapName:            state.mapName,
		TeamCTPlayers:      state.rosterCT,
		TeamTPlayers:       state.rosterT,
	}
}

// --- Settings/Notes passthroughs ---------------------------------------

func (s *CS2CasterService) GetOrCreateSettings(ctx context.Context, userTwitchID string) (*domain.CS2CasterSettings, error) {
	return s.casterRepo.GetOrCreateSettings(ctx, userTwitchID)
}

func (s *CS2CasterService) UpdateSettings(ctx context.Context, userTwitchID string, input domain.CS2CasterSettingsUpdateInput) error {
	return s.casterRepo.UpdateSettings(ctx, userTwitchID, input)
}

func (s *CS2CasterService) ListNotes(ctx context.Context, userTwitchID string) ([]*domain.CS2Note, error) {
	return s.casterRepo.ListNotes(ctx, userTwitchID)
}

func (s *CS2CasterService) CreateNote(ctx context.Context, userTwitchID string, input domain.CS2NoteCreateInput) (*domain.CS2Note, error) {
	if err := input.SubjectType.Validate(); err != nil {
		return nil, err
	}
	return s.casterRepo.CreateNote(ctx, userTwitchID, input)
}

func (s *CS2CasterService) UpdateNote(ctx context.Context, userTwitchID string, noteID int64, input domain.CS2NoteUpdateInput) error {
	err := s.casterRepo.UpdateNote(ctx, userTwitchID, noteID, input)
	if err == sql.ErrNoRows {
		return domain.ErrCS2NoteNotFound
	}
	return err
}

func (s *CS2CasterService) DeleteNote(ctx context.Context, userTwitchID string, noteID int64) error {
	err := s.casterRepo.DeleteNote(ctx, userTwitchID, noteID)
	if err == sql.ErrNoRows {
		return domain.ErrCS2NoteNotFound
	}
	return err
}
