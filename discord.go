package discordgo

import (
	"fmt"
	"net/http"
	"runtime"
	"time"

	"github.com/gorilla/websocket"
)

const VERSION = "0.30.5"
const clientTimeout = 20 * time.Second

func New(token string) (*Session, error) {
	session := &Session{
		State:                              NewState(),
		Ratelimiter:                        NewRatelimiter(),
		StateEnabled:                       true,
		Compress:                           true,
		Token:                              token,
		Client:                             &http.Client{Timeout: clientTimeout},
		Dialer:                             websocket.DefaultDialer,
		UserAgent:                          fmt.Sprintf("discordgo (https://github.com/jacobbernoulli/discordgo, v%s)", VERSION),
		sequence:                           new(int64),
		LastHeartbeatAck:                   time.Now().UTC(),
		ShardID:                            0,
		ShardCount:                         1,
		MaxRestRetries:                     3,
		ShouldReconnectOnError:             true,
		ShouldReconnectVoiceOnSessionError: true,
		ShouldRetryOnRateLimit:             true,
	}

	session.Identify = Identify{
		Token:          token,
		Compress:       true,
		LargeThreshold: 250,
		Intents:        IntentsAllWithoutPrivileged,
		Properties: IdentifyProperties{
			OS:      runtime.GOOS,
			Browser: fmt.Sprintf("discordgo v%s", VERSION),
		},
	}

	return session, nil
}
