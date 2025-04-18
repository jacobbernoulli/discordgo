package discordgo

import (
	"fmt"
	"net/http"
	"runtime"
	"time"

	"github.com/gorilla/websocket"
)

const VERSION = "0.28.1"

const clientTimeout = 20 * time.Second

func New(token string) (s *Session, err error) {
	s = &Session{
		State:                              NewState(),
		Ratelimiter:                        NewRatelimiter(),
		StateEnabled:                       true,
		Compress:                           true,
		ShouldReconnectOnError:             true,
		ShouldReconnectVoiceOnSessionError: true,
		ShouldRetryOnRateLimit:             true,
		ShardID:                            0,
		ShardCount:                         1,
		MaxRestRetries:                     3,
		Client:                             &http.Client{Timeout: clientTimeout},
		Dialer:                             websocket.DefaultDialer,
		UserAgent:                          fmt.Sprintf("discordgo (https://github.com/jacobbernoulli/discordgo, v%s)", VERSION),
		sequence:                           new(int64),
		LastHeartbeatAck:                   time.Now().UTC(),
	}

	s.Identify.Compress = true
	s.Identify.LargeThreshold = 250
	s.Identify.Intents = IntentsAllWithoutPrivileged
	s.Identify.Properties.OS = runtime.GOOS
	s.Identify.Properties.Browser = fmt.Sprintf("discordgo v%s", VERSION)
	s.Identify.Token = token
	s.Token = token

	return
}
