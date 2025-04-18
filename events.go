package discordgo

import "encoding/json"

type Connect struct{}

type Disconnect struct{}

type RateLimit struct {
	*TooManyRequests
	URL string
}

type Event struct {
	Operation int             `json:"op"`
	Sequence  int64           `json:"s"`
	Type      string          `json:"t"`
	RawData   json.RawMessage `json:"d"`
	Struct    interface{}     `json:"-"`
}

type Ready struct {
	Version         int          `json:"v"`
	SessionID       string       `json:"session_id"`
	User            *User        `json:"user"`
	Shard           *[2]int      `json:"shard"`
	Application     *Application `json:"application"`
	Guilds          []*Guild     `json:"guilds"`
	PrivateChannels []*Channel   `json:"private_channels"`
}

type ChannelCreate struct {
	*Channel
}

type ChannelUpdate struct {
	*Channel
	BeforeUpdate *Channel `json:"-"`
}

type ChannelDelete struct {
	*Channel
}

type ChannelPinsUpdate struct {
	LastPinTimestamp string `json:"last_pin_timestamp"`
	ChannelID        string `json:"channel_id"`
	GuildID          string `json:"guild_id,omitempty"`
}

type ThreadCreate struct {
	*Channel
	NewlyCreated bool `json:"newly_created"`
}

type ThreadUpdate struct {
	*Channel
	BeforeUpdate *Channel `json:"-"`
}

type ThreadDelete struct {
	*Channel
}

type ThreadListSync struct {
	GuildID    string          `json:"guild_id"`
	ChannelIDs []string        `json:"channel_ids"`
	Threads    []*Channel      `json:"threads"`
	Members    []*ThreadMember `json:"members"`
}

type ThreadMemberUpdate struct {
	*ThreadMember
	GuildID string `json:"guild_id"`
}

type ThreadMembersUpdate struct {
	ID             string              `json:"id"`
	GuildID        string              `json:"guild_id"`
	MemberCount    int                 `json:"member_count"`
	AddedMembers   []AddedThreadMember `json:"added_members"`
	RemovedMembers []string            `json:"removed_member_ids"`
}

type GuildCreate struct {
	*Guild
}

type GuildUpdate struct {
	*Guild
}

type GuildDelete struct {
	*Guild
	BeforeDelete *Guild `json:"-"`
}

type GuildBanAdd struct {
	User    *User  `json:"user"`
	GuildID string `json:"guild_id"`
}

type GuildBanRemove struct {
	User    *User  `json:"user"`
	GuildID string `json:"guild_id"`
}

type GuildMemberAdd struct {
	*Member
}

type GuildMemberUpdate struct {
	*Member
	BeforeUpdate *Member `json:"-"`
}

type GuildMemberRemove struct {
	*Member
}

type GuildRoleCreate struct {
	*GuildRole
}

type GuildRoleUpdate struct {
	*GuildRole
}

type GuildRoleDelete struct {
	RoleID  string `json:"role_id"`
	GuildID string `json:"guild_id"`
}

type GuildEmojisUpdate struct {
	GuildID string   `json:"guild_id"`
	Emojis  []*Emoji `json:"emojis"`
}

type GuildMembersChunk struct {
	GuildID    string      `json:"guild_id"`
	Members    []*Member   `json:"members"`
	ChunkIndex int         `json:"chunk_index"`
	ChunkCount int         `json:"chunk_count"`
	NotFound   []string    `json:"not_found,omitempty"`
	Presences  []*Presence `json:"presences,omitempty"`
	Nonce      string      `json:"nonce,omitempty"`
}

type GuildIntegrationsUpdate struct {
	GuildID string `json:"guild_id"`
}

type StageInstanceEventCreate struct {
	*StageInstance
}

type StageInstanceEventUpdate struct {
	*StageInstance
}

type StageInstanceEventDelete struct {
	*StageInstance
}

type GuildScheduledEventCreate struct {
	*GuildScheduledEvent
}

type GuildScheduledEventUpdate struct {
	*GuildScheduledEvent
}

type GuildScheduledEventDelete struct {
	*GuildScheduledEvent
}

type GuildScheduledEventUserAdd struct {
	GuildScheduledEventID string `json:"guild_scheduled_event_id"`
	UserID                string `json:"user_id"`
	GuildID               string `json:"guild_id"`
}

type GuildScheduledEventUserRemove struct {
	GuildScheduledEventID string `json:"guild_scheduled_event_id"`
	UserID                string `json:"user_id"`
	GuildID               string `json:"guild_id"`
}

type IntegrationCreate struct {
	*Integration
	GuildID string `json:"guild_id"`
}

type IntegrationUpdate struct {
	*Integration
	GuildID string `json:"guild_id"`
}

type IntegrationDelete struct {
	ID            string `json:"id"`
	GuildID       string `json:"guild_id"`
	ApplicationID string `json:"application_id,omitempty"`
}

type MessageCreate struct {
	*Message
}

func (m *MessageCreate) UnmarshalJSON(b []byte) error {
	return json.Unmarshal(b, &m.Message)
}

type MessageUpdate struct {
	*Message
	BeforeUpdate *Message `json:"-"`
}

func (m *MessageUpdate) UnmarshalJSON(b []byte) error {
	return json.Unmarshal(b, &m.Message)
}

type MessageDelete struct {
	*Message
	BeforeDelete *Message `json:"-"`
}

func (m *MessageDelete) UnmarshalJSON(b []byte) error {
	return json.Unmarshal(b, &m.Message)
}

type MessageReactionAdd struct {
	*MessageReaction
	Member *Member `json:"member,omitempty"`
}

type MessageReactionRemove struct {
	*MessageReaction
}

type MessageReactionRemoveAll struct {
	*MessageReaction
}

type PresencesReplace []*Presence

type PresenceUpdate struct {
	Presence
	GuildID string `json:"guild_id"`
}

type Resumed struct {
	Trace []string `json:"_trace"`
}

type TypingStart struct {
	UserID    string `json:"user_id"`
	ChannelID string `json:"channel_id"`
	GuildID   string `json:"guild_id,omitempty"`
	Timestamp int    `json:"timestamp"`
}

type UserUpdate struct {
	*User
}

type VoiceServerUpdate struct {
	Token    string `json:"token"`
	GuildID  string `json:"guild_id"`
	Endpoint string `json:"endpoint"`
}

type VoiceStateUpdate struct {
	*VoiceState
	BeforeUpdate *VoiceState `json:"-"`
}

type MessageDeleteBulk struct {
	Messages  []string `json:"ids"`
	ChannelID string   `json:"channel_id"`
	GuildID   string   `json:"guild_id"`
}

type WebhooksUpdate struct {
	GuildID   string `json:"guild_id"`
	ChannelID string `json:"channel_id"`
}

type InteractionCreate struct {
	*Interaction
}

func (i *InteractionCreate) UnmarshalJSON(b []byte) error {
	return json.Unmarshal(b, &i.Interaction)
}

type InviteCreate struct {
	*Invite
	ChannelID string `json:"channel_id"`
	GuildID   string `json:"guild_id"`
}

type InviteDelete struct {
	ChannelID string `json:"channel_id"`
	GuildID   string `json:"guild_id"`
	Code      string `json:"code"`
}

type ApplicationCommandPermissionsUpdate struct {
	*GuildApplicationCommandPermissions
}

type AutoModerationRuleCreate struct {
	*AutoModerationRule
}

type AutoModerationRuleUpdate struct {
	*AutoModerationRule
}

type AutoModerationRuleDelete struct {
	*AutoModerationRule
}

type AutoModerationActionExecution struct {
	GuildID              string                        `json:"guild_id"`
	Action               AutoModerationAction          `json:"action"`
	RuleID               string                        `json:"rule_id"`
	RuleTriggerType      AutoModerationRuleTriggerType `json:"rule_trigger_type"`
	UserID               string                        `json:"user_id"`
	ChannelID            string                        `json:"channel_id"`
	MessageID            string                        `json:"message_id"`
	AlertSystemMessageID string                        `json:"alert_system_message_id"`
	Content              string                        `json:"content"`
	MatchedKeyword       string                        `json:"matched_keyword"`
	MatchedContent       string                        `json:"matched_content"`
}

type GuildAuditLogEntryCreate struct {
	*AuditLogEntry
	GuildID string `json:"guild_id"`
}

type MessagePollVoteAdd struct {
	UserID    string `json:"user_id"`
	ChannelID string `json:"channel_id"`
	MessageID string `json:"message_id"`
	GuildID   string `json:"guild_id,omitempty"`
	AnswerID  int    `json:"answer_id"`
}

type MessagePollVoteRemove struct {
	UserID    string `json:"user_id"`
	ChannelID string `json:"channel_id"`
	MessageID string `json:"message_id"`
	GuildID   string `json:"guild_id,omitempty"`
	AnswerID  int    `json:"answer_id"`
}

type EntitlementCreate struct {
	*Entitlement
}

type EntitlementUpdate struct {
	*Entitlement
}

type EntitlementDelete struct {
	*Entitlement
}

type SubscriptionCreate struct {
	*Subscription
}

type SubscriptionUpdate struct {
	*Subscription
}

type SubscriptionDelete struct {
	*Subscription
}
