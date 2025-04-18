package discordgo

import (
	"bytes"
	"compress/zlib"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync/atomic"
	"time"

	"github.com/gorilla/websocket"
)

type resumePacket struct {
	Op   int `json:"op"`
	Data struct {
		Token     string `json:"token"`
		SessionID string `json:"session_id"`
		Sequence  int64  `json:"seq"`
	} `json:"d"`
}

func (s *Session) Open() error {
	s.log(LogInformational, "called")

	s.Lock()
	defer s.Unlock()

	if s.wsConn != nil {
		return ErrWSAlreadyOpen
	}

	if s.gateway == "" {
		var err error
		s.gateway, err = s.Gateway()
		if err != nil {
			return err
		}
		s.gateway = s.gateway + "?v=" + APIVersion + "&encoding=json"
	}

	s.log(LogInformational, "connecting to gateway %s", s.gateway)
	header := http.Header{}
	header.Add("accept-encoding", "zlib")
	var err error
	s.wsConn, _, err = s.Dialer.Dial(s.gateway, header)
	if err != nil {
		s.log(LogError, "error connecting to gateway %s, %s", s.gateway, err)
		s.gateway = ""
		s.wsConn = nil
		return err
	}

	s.wsConn.SetCloseHandler(func(code int, text string) error {
		return nil
	})

	defer func() {
		if err != nil {
			s.wsConn.Close()
			s.wsConn = nil
		}
	}()

	mt, m, err := s.wsConn.ReadMessage()
	if err != nil {
		return err
	}

	e, err := s.onEvent(mt, m)
	if err != nil {
		return err
	}

	if e.Operation != 10 {
		return fmt.Errorf("expecting Op 10, got Op %d instead", e.Operation)
	}
	s.log(LogInformational, "Op 10 Hello Packet received from Discord.")
	s.LastHeartbeatAck = time.Now().UTC()

	var h helloOp
	if err := Unmarshal(e.RawData, &h); err != nil {
		return fmt.Errorf("error unmarshalling helloOp: %s", err)
	}

	sequence := atomic.LoadInt64(s.sequence)
	if s.sessionID == "" && sequence == 0 {
		if err := s.identify(); err != nil {
			return fmt.Errorf("error sending identify packet to gateway: %s - %s", s.gateway, err)
		}
	} else {
		p := resumePacket{
			Op: 6,
			Data: struct {
				Token     string `json:"token"`
				SessionID string `json:"session_id"`
				Sequence  int64  `json:"seq"`
			}{
				Token:     s.Token,
				SessionID: s.sessionID,
				Sequence:  sequence,
			},
		}

		s.log(LogInformational, "Sending resume packet to gateway...")
		s.wsMutex.Lock()
		err := s.wsConn.WriteJSON(p)
		s.wsMutex.Unlock()
		if err != nil {
			return fmt.Errorf("error sending gateway resume packet: %s - %s", s.gateway, err)
		}
	}

	if s.State == nil {
		s.State = NewState()
		s.State.TrackChannels = false
		s.State.TrackEmojis = false
		s.State.TrackMembers = false
		s.State.TrackRoles = false
		s.State.TrackVoice = false
	}

	mt, m, err = s.wsConn.ReadMessage()
	if err != nil {
		return err
	}

	e, err = s.onEvent(mt, m)
	if err != nil {
		return err
	}

	if e.Type != `READY` && e.Type != `RESUMED` {
		s.log(LogWarning, "Expected READY/RESUMED, instead got:\n%#v\n", e)
	}
	s.log(LogInformational, "First Packet:\n%#v\n", e)

	s.log(LogInformational, "Connected to Discord, emitting connect event...")
	s.handleEvent(connectEventType, &Connect{})

	if s.VoiceConnections == nil {
		s.log(LogInformational, "Creating new VoiceConnections map...")
		s.VoiceConnections = make(map[string]*VoiceConnection)
	}

	s.listening = make(chan interface{})
	go s.heartbeat(s.wsConn, s.listening, h.HeartbeatInterval)
	go s.listen(s.wsConn, s.listening)

	s.log(LogInformational, "exiting")
	return nil
}

func (s *Session) listen(wsConn *websocket.Conn, listening <-chan interface{}) {
	s.log(LogInformational, "called")

	for {
		messageType, message, err := wsConn.ReadMessage()
		if err != nil {
			s.RLock()
			sameConnection := s.wsConn == wsConn
			s.RUnlock()

			if sameConnection {
				s.log(LogWarning, "error reading from gateway %s websocket, %s", s.gateway, err)
				err := s.Close()
				if err != nil {
					s.log(LogWarning, "error closing session connection, %s", err)
				}

				s.log(LogInformational, "calling reconnect() now")
				s.reconnect()
			}

			return
		}

		select {
		case <-listening:
			return

		default:
			s.onEvent(messageType, message)

		}
	}
}

type heartbeatOp struct {
	Op   int   `json:"op"`
	Data int64 `json:"d"`
}

type helloOp struct {
	HeartbeatInterval time.Duration `json:"heartbeat_interval"`
}

const FailedHeartbeatAcks time.Duration = 5 * time.Millisecond

func (s *Session) HeartbeatLatency() time.Duration {
	return s.LastHeartbeatAck.Sub(s.LastHeartbeatSent)

}

func (s *Session) heartbeat(wsConn *websocket.Conn, listening <-chan interface{}, heartbeatInterval time.Duration) {
	s.log(LogInformational, "called")

	if listening == nil || wsConn == nil {
		return
	}

	var err error
	ticker := time.NewTicker(heartbeatInterval * time.Millisecond)
	defer ticker.Stop()

	for {
		s.RLock()
		last := s.LastHeartbeatAck
		s.RUnlock()
		sequence := atomic.LoadInt64(s.sequence)
		s.log(LogDebug, "sending gateway websocket heartbeat seq %d", sequence)
		s.wsMutex.Lock()
		s.LastHeartbeatSent = time.Now().UTC()
		err = wsConn.WriteJSON(heartbeatOp{1, sequence})
		s.wsMutex.Unlock()
		if err != nil || time.Now().UTC().Sub(last) > (heartbeatInterval*FailedHeartbeatAcks) {
			if err != nil {
				s.log(LogError, "error sending heartbeat to gateway %s, %s", s.gateway, err)
			} else {
				s.log(LogError, "haven't gotten a heartbeat ACK in %v, triggering a reconnection", time.Now().UTC().Sub(last))
			}
			s.Close()
			s.reconnect()
			return
		}
		s.Lock()
		s.DataReady = true
		s.Unlock()

		select {
		case <-ticker.C:
		case <-listening:
			return
		}
	}
}

type UpdateStatusData struct {
	IdleSince  *int        `json:"since"`
	Activities []*Activity `json:"activities"`
	AFK        bool        `json:"afk"`
	Status     string      `json:"status"`
}

type updateStatusOp struct {
	Op   int              `json:"op"`
	Data UpdateStatusData `json:"d"`
}

func newUpdateStatusData(idle int, activityType ActivityType, name, url string) *UpdateStatusData {
	usd := &UpdateStatusData{
		Status: "online",
	}

	if idle > 0 {
		usd.IdleSince = &idle
	}

	if name != "" {
		usd.Activities = []*Activity{{
			Name: name,
			Type: activityType,
			URL:  url,
		}}
	}

	return usd
}

func (s *Session) UpdateGameStatus(idle int, name string) (err error) {
	return s.UpdateStatusComplex(*newUpdateStatusData(idle, ActivityTypeGame, name, ""))
}

func (s *Session) UpdateWatchStatus(idle int, name string) (err error) {
	return s.UpdateStatusComplex(*newUpdateStatusData(idle, ActivityTypeWatching, name, ""))
}

func (s *Session) UpdateStreamingStatus(idle int, name string, url string) (err error) {
	gameType := ActivityTypeGame
	if url != "" {
		gameType = ActivityTypeStreaming
	}
	return s.UpdateStatusComplex(*newUpdateStatusData(idle, gameType, name, url))
}

func (s *Session) UpdateListeningStatus(name string) (err error) {
	return s.UpdateStatusComplex(*newUpdateStatusData(0, ActivityTypeListening, name, ""))
}

func (s *Session) UpdateCustomStatus(state string) (err error) {
	data := UpdateStatusData{
		Status: "online",
	}

	if state != "" {
		data.Activities = []*Activity{{
			Name:  "Custom Status",
			Type:  ActivityTypeCustom,
			State: state,
		}}
	}

	return s.UpdateStatusComplex(data)
}

func (s *Session) UpdateStatusComplex(usd UpdateStatusData) (err error) {
	if usd.Activities == nil {
		usd.Activities = make([]*Activity, 0)
	}

	s.RLock()
	defer s.RUnlock()
	if s.wsConn == nil {
		return ErrWSNotFound
	}

	s.wsMutex.Lock()
	err = s.wsConn.WriteJSON(updateStatusOp{3, usd})
	s.wsMutex.Unlock()

	return
}

type requestGuildMembersData struct {
	GuildIDs  []string  `json:"guild_id"`
	Query     *string   `json:"query,omitempty"`
	UserIDs   *[]string `json:"user_ids,omitempty"`
	Limit     int       `json:"limit"`
	Nonce     string    `json:"nonce,omitempty"`
	Presences bool      `json:"presences"`
}

type requestGuildMembersOp struct {
	Op   int                     `json:"op"`
	Data requestGuildMembersData `json:"d"`
}

func (s *Session) RequestGuildMembers(guildID, query string, limit int, nonce string, presences bool) error {
	return s.RequestGuildMembersBatch([]string{guildID}, query, limit, nonce, presences)
}

func (s *Session) RequestGuildMembersList(guildID string, userIDs []string, limit int, nonce string, presences bool) error {
	return s.RequestGuildMembersBatchList([]string{guildID}, userIDs, limit, nonce, presences)
}

func (s *Session) RequestGuildMembersBatch(guildIDs []string, query string, limit int, nonce string, presences bool) (err error) {
	data := requestGuildMembersData{
		GuildIDs:  guildIDs,
		Query:     &query,
		Limit:     limit,
		Nonce:     nonce,
		Presences: presences,
	}
	err = s.requestGuildMembers(data)
	return
}

func (s *Session) RequestGuildMembersBatchList(guildIDs []string, userIDs []string, limit int, nonce string, presences bool) (err error) {
	data := requestGuildMembersData{
		GuildIDs:  guildIDs,
		UserIDs:   &userIDs,
		Limit:     limit,
		Nonce:     nonce,
		Presences: presences,
	}
	err = s.requestGuildMembers(data)
	return
}

func (s *Session) GatewayWriteStruct(data interface{}) (err error) {
	s.RLock()
	defer s.RUnlock()
	if s.wsConn == nil {
		return ErrWSNotFound
	}

	s.wsMutex.Lock()
	err = s.wsConn.WriteJSON(data)
	s.wsMutex.Unlock()

	return err
}

func (s *Session) requestGuildMembers(data requestGuildMembersData) (err error) {
	s.log(LogInformational, "called")

	s.RLock()
	defer s.RUnlock()
	if s.wsConn == nil {
		return ErrWSNotFound
	}

	s.wsMutex.Lock()
	err = s.wsConn.WriteJSON(requestGuildMembersOp{8, data})
	s.wsMutex.Unlock()

	return
}

func (s *Session) onEvent(messageType int, message []byte) (*Event, error) {
	var err error
	var reader io.Reader
	reader = bytes.NewBuffer(message)

	if messageType == websocket.BinaryMessage {

		z, err2 := zlib.NewReader(reader)
		if err2 != nil {
			s.log(LogError, "error uncompressing websocket message, %s", err)
			return nil, err2
		}

		defer func() {
			err3 := z.Close()
			if err3 != nil {
				s.log(LogWarning, "error closing zlib, %s", err)
			}
		}()

		reader = z
	}

	var e *Event
	decoder := json.NewDecoder(reader)
	if err = decoder.Decode(&e); err != nil {
		s.log(LogError, "error decoding websocket message, %s", err)
		return e, err
	}

	s.log(LogDebug, "Op: %d, Seq: %d, Type: %s, Data: %s\n\n", e.Operation, e.Sequence, e.Type, string(e.RawData))

	if e.Operation == 1 {
		s.log(LogInformational, "sending heartbeat in response to Op1")
		s.wsMutex.Lock()
		err = s.wsConn.WriteJSON(heartbeatOp{1, atomic.LoadInt64(s.sequence)})
		s.wsMutex.Unlock()
		if err != nil {
			s.log(LogError, "error sending heartbeat in response to Op1")
			return e, err
		}

		return e, nil
	}

	if e.Operation == 7 {
		s.log(LogInformational, "Closing and reconnecting in response to Op7")
		s.CloseWithCode(websocket.CloseServiceRestart)
		s.reconnect()
		return e, nil
	}

	if e.Operation == 9 {

		s.log(LogInformational, "sending identify packet to gateway in response to Op9")

		err = s.identify()
		if err != nil {
			s.log(LogWarning, "error sending gateway identify packet, %s, %s", s.gateway, err)
			return e, err
		}

		return e, nil
	}

	if e.Operation == 10 {
		return e, nil
	}

	if e.Operation == 11 {
		s.Lock()
		s.LastHeartbeatAck = time.Now().UTC()
		s.Unlock()
		s.log(LogDebug, "got heartbeat ACK")
		return e, nil
	}

	if e.Operation != 0 {
		s.log(LogWarning, "unknown Op: %d, Seq: %d, Type: %s, Data: %s, message: %s", e.Operation, e.Sequence, e.Type, string(e.RawData), string(message))
		return e, nil
	}

	atomic.StoreInt64(s.sequence, e.Sequence)

	if eh, ok := registeredInterfaceProviders[e.Type]; ok {
		e.Struct = eh.New()

		if err = Unmarshal(e.RawData, e.Struct); err != nil {
			s.log(LogError, "error unmarshalling %s event, %s", e.Type, err)
		}
		s.handleEvent(e.Type, e.Struct)
	} else {
		s.log(LogWarning, "unknown event: Op: %d, Seq: %d, Type: %s, Data: %s", e.Operation, e.Sequence, e.Type, string(e.RawData))
	}

	s.handleEvent(eventEventType, e)

	return e, nil
}

type voiceChannelJoinData struct {
	GuildID   *string `json:"guild_id"`
	ChannelID *string `json:"channel_id"`
	SelfMute  bool    `json:"self_mute"`
	SelfDeaf  bool    `json:"self_deaf"`
}

type voiceChannelJoinOp struct {
	Op   int                  `json:"op"`
	Data voiceChannelJoinData `json:"d"`
}

func (s *Session) ChannelVoiceJoin(gID, cID string, mute, deaf bool) (voice *VoiceConnection, err error) {
	s.log(LogInformational, "called")

	s.RLock()
	voice = s.VoiceConnections[gID]
	s.RUnlock()

	if voice == nil {
		voice = &VoiceConnection{}
		s.Lock()
		s.VoiceConnections[gID] = voice
		s.Unlock()
	}

	voice.Lock()
	voice.GuildID = gID
	voice.ChannelID = cID
	voice.deaf = deaf
	voice.mute = mute
	voice.session = s
	voice.Unlock()

	err = s.ChannelVoiceJoinManual(gID, cID, mute, deaf)
	if err != nil {
		return
	}

	err = voice.waitUntilConnected()
	if err != nil {
		s.log(LogWarning, "error waiting for voice to connect, %s", err)
		voice.Close()
		return
	}

	return
}

func (s *Session) ChannelVoiceJoinManual(gID, cID string, mute, deaf bool) (err error) {
	s.log(LogInformational, "called")

	var channelID *string
	if cID == "" {
		channelID = nil
	} else {
		channelID = &cID
	}

	data := voiceChannelJoinOp{4, voiceChannelJoinData{&gID, channelID, mute, deaf}}
	s.wsMutex.Lock()
	err = s.wsConn.WriteJSON(data)
	s.wsMutex.Unlock()
	return
}

func (s *Session) onVoiceStateUpdate(st *VoiceStateUpdate) {
	if st.ChannelID == "" {
		return
	}

	s.RLock()
	voice, exists := s.VoiceConnections[st.GuildID]
	s.RUnlock()
	if !exists {
		return
	}

	if s.State.User.ID != st.UserID {
		return
	}

	voice.Lock()
	voice.UserID = st.UserID
	voice.sessionID = st.SessionID
	voice.ChannelID = st.ChannelID
	voice.Unlock()
}

func (s *Session) onVoiceServerUpdate(st *VoiceServerUpdate) {
	s.log(LogInformational, "called")

	s.RLock()
	voice, exists := s.VoiceConnections[st.GuildID]
	s.RUnlock()

	if !exists {
		return
	}

	voice.Close()

	voice.Lock()
	voice.token = st.Token
	voice.endpoint = st.Endpoint
	voice.GuildID = st.GuildID
	voice.Unlock()

	err := voice.open()
	if err != nil {
		s.log(LogError, "onVoiceServerUpdate voice.open, %s", err)
	}
}

type identifyOp struct {
	Op   int      `json:"op"`
	Data Identify `json:"d"`
}

func (s *Session) identify() error {
	s.log(LogDebug, "called")

	if !s.Compress {
		s.Identify.Compress = false
	}

	if s.Token != "" && s.Identify.Token == "" {
		s.Identify.Token = s.Token
	}

	if s.ShardCount > 1 {

		if s.ShardID >= s.ShardCount {
			return ErrWSShardBounds
		}

		s.Identify.Shard = &[2]int{s.ShardID, s.ShardCount}
	}

	op := identifyOp{2, s.Identify}
	s.log(LogDebug, "Identify Packet: \n%#v", op)
	s.wsMutex.Lock()
	err := s.wsConn.WriteJSON(op)
	s.wsMutex.Unlock()

	return err
}

func (s *Session) reconnect() {
	s.log(LogInformational, "called")

	var err error

	if s.ShouldReconnectOnError {

		wait := time.Duration(1)

		for {
			s.log(LogInformational, "trying to reconnect to gateway")

			err = s.Open()
			if err == nil {
				s.log(LogInformational, "successfully reconnected to gateway")
				if s.ShouldReconnectVoiceOnSessionError {
					s.RLock()
					defer s.RUnlock()
					for _, v := range s.VoiceConnections {
						s.log(LogInformational, "reconnecting voice connection to guild %s", v.GuildID)
						go v.reconnect()
						time.Sleep(1 * time.Second)
					}
				}
				return
			}

			if err == ErrWSAlreadyOpen {
				s.log(LogInformational, "Websocket already exists, no need to reconnect")
				return
			}

			s.log(LogError, "error reconnecting to gateway, %s", err)

			<-time.After(wait * time.Second)
			wait *= 2
			if wait > 600 {
				wait = 600
			}
		}
	}
}

func (s *Session) Close() error {
	return s.CloseWithCode(websocket.CloseNormalClosure)
}

func (s *Session) CloseWithCode(closeCode int) (err error) {
	s.log(LogInformational, "called")
	s.Lock()

	s.DataReady = false

	if s.listening != nil {
		s.log(LogInformational, "closing listening channel")
		close(s.listening)
		s.listening = nil
	}

	if s.wsConn != nil {

		s.log(LogInformational, "sending close frame")
		s.wsMutex.Lock()
		err := s.wsConn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(closeCode, ""))
		s.wsMutex.Unlock()
		if err != nil {
			s.log(LogInformational, "error closing websocket, %s", err)
		}

		time.Sleep(1 * time.Second)

		s.log(LogInformational, "closing gateway websocket")
		err = s.wsConn.Close()
		if err != nil {
			s.log(LogInformational, "error closing websocket, %s", err)
		}

		s.wsConn = nil
	}

	s.Unlock()

	s.log(LogInformational, "emit disconnect event")
	s.handleEvent(disconnectEventType, &Disconnect{})

	return
}
