package discordgo

import (
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"regexp"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type RecurrenceRuleFrequency int

const (
	YEARLY  RecurrenceRuleFrequency = 0
	MONTHLY RecurrenceRuleFrequency = 1
	WEEKLY  RecurrenceRuleFrequency = 2
	DAILY   RecurrenceRuleFrequency = 3
)

type RecurrenceRuleWeekday int

const (
	MONDAY    RecurrenceRuleWeekday = 0
	TUESDAY   RecurrenceRuleWeekday = 1
	WEDNESDAY RecurrenceRuleWeekday = 2
	THURSDAY  RecurrenceRuleWeekday = 3
	FRIDAY    RecurrenceRuleWeekday = 4
	SATURDAY  RecurrenceRuleWeekday = 5
	SUNDAY    RecurrenceRuleWeekday = 6
)

type RecurrenceRuleNWeekDay struct {
	N   int                   `json:"n"`
	Day RecurrenceRuleWeekday `json:"day"`
}

type RecurrenceRuleMonth int

const (
	JANUARY   RecurrenceRuleMonth = 1
	FEBRUARY  RecurrenceRuleMonth = 2
	MARCH     RecurrenceRuleMonth = 3
	APRIL     RecurrenceRuleMonth = 4
	MAY       RecurrenceRuleMonth = 5
	JUNE      RecurrenceRuleMonth = 6
	JULY      RecurrenceRuleMonth = 7
	AUGUST    RecurrenceRuleMonth = 8
	SEPTEMBER RecurrenceRuleMonth = 9
	OCTOBER   RecurrenceRuleMonth = 10
	NOVEMBER  RecurrenceRuleMonth = 11
	DECEMBER  RecurrenceRuleMonth = 12
)

type RecurrenceRule struct {
	Start      time.Time                `json:"start"`
	End        *time.Time               `json:"end,omitempty"`
	Frequency  RecurrenceRuleFrequency  `json:"frequency"`
	Interval   int                      `json:"interval"`
	ByWeekday  []RecurrenceRuleWeekday  `json:"by_weekday,omitempty"`
	ByNWeekday []RecurrenceRuleNWeekDay `json:"by_n_weekday,omitempty"`
	ByMonth    []RecurrenceRuleMonth    `json:"by_month,omitempty"`
	ByMonthDay []int                    `json:"by_month_day,omitempty"`
	ByYearDay  []int                    `json:"by_year_day,omitempty"`
	Count      int                      `json:"count,omitempty"`
}

type Session struct {
	sync.RWMutex
	Token                              string
	MFA                                bool
	Debug                              bool
	LogLevel                           int
	ShouldReconnectOnError             bool
	ShouldReconnectVoiceOnSessionError bool
	ShouldRetryOnRateLimit             bool
	Identify                           Identify
	Compress                           bool
	ShardID                            int
	ShardCount                         int
	StateEnabled                       bool
	SyncEvents                         bool
	DataReady                          bool
	MaxRestRetries                     int
	VoiceReady                         bool
	UDPReady                           bool
	VoiceConnections                   map[string]*VoiceConnection
	State                              *State
	Client                             *http.Client
	Dialer                             *websocket.Dialer
	UserAgent                          string
	LastHeartbeatAck                   time.Time
	LastHeartbeatSent                  time.Time
	Ratelimiter                        *RateLimiter
	handlersMu                         sync.RWMutex
	handlers                           map[string][]*eventHandlerInstance
	onceHandlers                       map[string][]*eventHandlerInstance
	wsConn                             *websocket.Conn
	listening                          chan interface{}
	sequence                           *int64
	gateway                            string
	sessionID                          string
	wsMutex                            sync.Mutex
}

type ApplicationIntegrationType uint

const (
	ApplicationIntegrationGuildInstall ApplicationIntegrationType = 0
	ApplicationIntegrationUserInstall  ApplicationIntegrationType = 1
)

type ApplicationInstallParams struct {
	Scopes      []string `json:"scopes"`
	Permissions int64    `json:"permissions,string"`
}

type ApplicationIntegrationTypeConfig struct {
	OAuth2InstallParams *ApplicationInstallParams `json:"oauth2_install_params,omitempty"`
}

type Application struct {
	ID                     string                                                           `json:"id,omitempty"`
	Name                   string                                                           `json:"name"`
	Icon                   string                                                           `json:"icon,omitempty"`
	Description            string                                                           `json:"description,omitempty"`
	RPCOrigins             []string                                                         `json:"rpc_origins,omitempty"`
	BotPublic              bool                                                             `json:"bot_public,omitempty"`
	BotRequireCodeGrant    bool                                                             `json:"bot_require_code_grant,omitempty"`
	TermsOfServiceURL      string                                                           `json:"terms_of_service_url"`
	PrivacyProxyURL        string                                                           `json:"privacy_policy_url"`
	Owner                  *User                                                            `json:"owner"`
	Summary                string                                                           `json:"summary"`
	VerifyKey              string                                                           `json:"verify_key"`
	Team                   *Team                                                            `json:"team"`
	GuildID                string                                                           `json:"guild_id"`
	PrimarySKUID           string                                                           `json:"primary_sku_id"`
	Slug                   string                                                           `json:"slug"`
	CoverImage             string                                                           `json:"cover_image"`
	Flags                  int                                                              `json:"flags,omitempty"`
	IntegrationTypesConfig map[ApplicationIntegrationType]*ApplicationIntegrationTypeConfig `json:"integration_types,omitempty"`
}

type ApplicationRoleConnectionMetadataType int

const (
	ApplicationRoleConnectionMetadataIntegerLessThanOrEqual     ApplicationRoleConnectionMetadataType = 1
	ApplicationRoleConnectionMetadataIntegerGreaterThanOrEqual  ApplicationRoleConnectionMetadataType = 2
	ApplicationRoleConnectionMetadataIntegerEqual               ApplicationRoleConnectionMetadataType = 3
	ApplicationRoleConnectionMetadataIntegerNotEqual            ApplicationRoleConnectionMetadataType = 4
	ApplicationRoleConnectionMetadataDatetimeLessThanOrEqual    ApplicationRoleConnectionMetadataType = 5
	ApplicationRoleConnectionMetadataDatetimeGreaterThanOrEqual ApplicationRoleConnectionMetadataType = 6
	ApplicationRoleConnectionMetadataBooleanEqual               ApplicationRoleConnectionMetadataType = 7
	ApplicationRoleConnectionMetadataBooleanNotEqual            ApplicationRoleConnectionMetadataType = 8
)

type ApplicationRoleConnectionMetadata struct {
	Type                     ApplicationRoleConnectionMetadataType `json:"type"`
	Key                      string                                `json:"key"`
	Name                     string                                `json:"name"`
	NameLocalizations        map[Locale]string                     `json:"name_localizations"`
	Description              string                                `json:"description"`
	DescriptionLocalizations map[Locale]string                     `json:"description_localizations"`
}

type ApplicationRoleConnection struct {
	PlatformName     string            `json:"platform_name"`
	PlatformUsername string            `json:"platform_username"`
	Metadata         map[string]string `json:"metadata"`
}

type UserConnection struct {
	ID           string         `json:"id"`
	Name         string         `json:"name"`
	Type         string         `json:"type"`
	Revoked      bool           `json:"revoked"`
	Integrations []*Integration `json:"integrations"`
}

type Integration struct {
	ID                string             `json:"id"`
	Name              string             `json:"name"`
	Type              string             `json:"type"`
	Enabled           bool               `json:"enabled"`
	Syncing           bool               `json:"syncing"`
	RoleID            string             `json:"role_id"`
	EnableEmoticons   bool               `json:"enable_emoticons"`
	ExpireBehavior    ExpireBehavior     `json:"expire_behavior"`
	ExpireGracePeriod int                `json:"expire_grace_period"`
	User              *User              `json:"user"`
	Account           IntegrationAccount `json:"account"`
	SyncedAt          time.Time          `json:"synced_at"`
}

type ExpireBehavior int

const (
	ExpireBehaviorRemoveRole ExpireBehavior = 0
	ExpireBehaviorKick       ExpireBehavior = 1
)

type IntegrationAccount struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type VoiceRegion struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	Optimal    bool   `json:"optimal"`
	Deprecated bool   `json:"deprecated"`
	Custom     bool   `json:"custom"`
}

type InviteTargetType uint8

const (
	InviteTargetStream              InviteTargetType = 1
	InviteTargetEmbeddedApplication InviteTargetType = 2
)

type Invite struct {
	Guild             *Guild           `json:"guild"`
	Channel           *Channel         `json:"channel"`
	Inviter           *User            `json:"inviter"`
	Code              string           `json:"code"`
	CreatedAt         time.Time        `json:"created_at"`
	MaxAge            int              `json:"max_age"`
	Uses              int              `json:"uses"`
	MaxUses           int              `json:"max_uses"`
	Revoked           bool             `json:"revoked"`
	Temporary         bool             `json:"temporary"`
	Unique            bool             `json:"unique"`
	TargetUser        *User            `json:"target_user"`
	TargetType        InviteTargetType `json:"target_type"`
	TargetApplication *Application     `json:"target_application"`

	ApproximatePresenceCount int `json:"approximate_presence_count"`
	ApproximateMemberCount   int `json:"approximate_member_count"`

	ExpiresAt *time.Time `json:"expires_at"`
}

type ChannelType int

const (
	ChannelTypeGuildText          ChannelType = 0
	ChannelTypeDM                 ChannelType = 1
	ChannelTypeGuildVoice         ChannelType = 2
	ChannelTypeGroupDM            ChannelType = 3
	ChannelTypeGuildCategory      ChannelType = 4
	ChannelTypeGuildNews          ChannelType = 5
	ChannelTypeGuildStore         ChannelType = 6
	ChannelTypeGuildNewsThread    ChannelType = 10
	ChannelTypeGuildPublicThread  ChannelType = 11
	ChannelTypeGuildPrivateThread ChannelType = 12
	ChannelTypeGuildStageVoice    ChannelType = 13
	ChannelTypeGuildDirectory     ChannelType = 14
	ChannelTypeGuildForum         ChannelType = 15
	ChannelTypeGuildMedia         ChannelType = 16
)

type ChannelFlags int

const (
	ChannelFlagPinned     ChannelFlags = 1 << 1
	ChannelFlagRequireTag ChannelFlags = 1 << 4
)

type ForumSortOrderType int

const (
	ForumSortOrderLatestActivity ForumSortOrderType = 0
	ForumSortOrderCreationDate   ForumSortOrderType = 1
)

type ForumLayout int

const (
	ForumLayoutNotSet      ForumLayout = 0
	ForumLayoutListView    ForumLayout = 1
	ForumLayoutGalleryView ForumLayout = 2
)

type Channel struct {
	ID                            string                 `json:"id"`
	GuildID                       string                 `json:"guild_id"`
	Name                          string                 `json:"name"`
	Topic                         string                 `json:"topic"`
	Type                          ChannelType            `json:"type"`
	LastMessageID                 string                 `json:"last_message_id"`
	LastPinTimestamp              *time.Time             `json:"last_pin_timestamp"`
	MessageCount                  int                    `json:"message_count"`
	MemberCount                   int                    `json:"member_count"`
	NSFW                          bool                   `json:"nsfw"`
	Icon                          string                 `json:"icon"`
	Position                      int                    `json:"position"`
	Bitrate                       int                    `json:"bitrate"`
	Recipients                    []*User                `json:"recipients"`
	Messages                      []*Message             `json:"-"`
	PermissionOverwrites          []*PermissionOverwrite `json:"permission_overwrites"`
	UserLimit                     int                    `json:"user_limit"`
	ParentID                      string                 `json:"parent_id"`
	RateLimitPerUser              int                    `json:"rate_limit_per_user"`
	OwnerID                       string                 `json:"owner_id"`
	ApplicationID                 string                 `json:"application_id"`
	ThreadMetadata                *ThreadMetadata        `json:"thread_metadata,omitempty"`
	Member                        *ThreadMember          `json:"thread_member"`
	Members                       []*ThreadMember        `json:"-"`
	Flags                         ChannelFlags           `json:"flags"`
	AvailableTags                 []ForumTag             `json:"available_tags"`
	AppliedTags                   []string               `json:"applied_tags"`
	DefaultReactionEmoji          ForumDefaultReaction   `json:"default_reaction_emoji"`
	DefaultThreadRateLimitPerUser int                    `json:"default_thread_rate_limit_per_user"`
	DefaultSortOrder              *ForumSortOrderType    `json:"default_sort_order"`
	DefaultForumLayout            ForumLayout            `json:"default_forum_layout"`
}

func (c *Channel) Mention() string {
	return fmt.Sprintf("<#%s>", c.ID)
}

func (c *Channel) IsThread() bool {
	return c.Type == ChannelTypeGuildPublicThread || c.Type == ChannelTypeGuildPrivateThread || c.Type == ChannelTypeGuildNewsThread
}

type ThreadArchiveDuration int

const (
	ThreadArchiveDurationOneHour   = 60
	ThreadArchiveDurationOneDay    = 1440
	ThreadArchiveDurationThreeDays = 4320
	ThreadArchiveDurationOneWeek   = 10080
)

type ChannelEdit struct {
	Name                          string                 `json:"name,omitempty"`
	Topic                         string                 `json:"topic,omitempty"`
	NSFW                          *bool                  `json:"nsfw,omitempty"`
	Position                      *int                   `json:"position,omitempty"`
	Bitrate                       int                    `json:"bitrate,omitempty"`
	UserLimit                     int                    `json:"user_limit,omitempty"`
	PermissionOverwrites          []*PermissionOverwrite `json:"permission_overwrites,omitempty"`
	ParentID                      string                 `json:"parent_id,omitempty"`
	RateLimitPerUser              *int                   `json:"rate_limit_per_user,omitempty"`
	Flags                         *ChannelFlags          `json:"flags,omitempty"`
	DefaultThreadRateLimitPerUser *int                   `json:"default_thread_rate_limit_per_user,omitempty"`
	Archived                      *bool                  `json:"archived,omitempty"`
	AutoArchiveDuration           ThreadArchiveDuration  `json:"auto_archive_duration,omitempty"`
	Locked                        *bool                  `json:"locked,omitempty"`
	Invitable                     *bool                  `json:"invitable,omitempty"`
	AvailableTags                 *[]ForumTag            `json:"available_tags,omitempty"`
	DefaultReactionEmoji          *ForumDefaultReaction  `json:"default_reaction_emoji,omitempty"`
	DefaultSortOrder              *ForumSortOrderType    `json:"default_sort_order,omitempty"`
	DefaultForumLayout            *ForumLayout           `json:"default_forum_layout,omitempty"`
	AppliedTags                   *[]string              `json:"applied_tags,omitempty"`
}

type ChannelFollow struct {
	ChannelID string `json:"channel_id"`
	WebhookID string `json:"webhook_id"`
}

type PermissionOverwriteType int

const (
	PermissionOverwriteTypeRole   PermissionOverwriteType = 0
	PermissionOverwriteTypeMember PermissionOverwriteType = 1
)

type PermissionOverwrite struct {
	ID    string                  `json:"id"`
	Type  PermissionOverwriteType `json:"type"`
	Deny  int64                   `json:"deny,string"`
	Allow int64                   `json:"allow,string"`
}

type ThreadStart struct {
	Name                string                `json:"name"`
	AutoArchiveDuration ThreadArchiveDuration `json:"auto_archive_duration,omitempty"`
	Type                ChannelType           `json:"type,omitempty"`
	Invitable           bool                  `json:"invitable"`
	RateLimitPerUser    int                   `json:"rate_limit_per_user,omitempty"`
	AppliedTags         []string              `json:"applied_tags,omitempty"`
}

type ThreadMetadata struct {
	Archived            bool                  `json:"archived"`
	AutoArchiveDuration ThreadArchiveDuration `json:"auto_archive_duration"`
	ArchiveTimestamp    time.Time             `json:"archive_timestamp"`
	Locked              bool                  `json:"locked"`
	Invitable           bool                  `json:"invitable"`
}

type ThreadMember struct {
	ID            string    `json:"id,omitempty"`
	UserID        string    `json:"user_id,omitempty"`
	JoinTimestamp time.Time `json:"join_timestamp"`
	Flags         int       `json:"flags"`
	Member        *Member   `json:"member,omitempty"`
}

type ThreadsList struct {
	Threads []*Channel      `json:"threads"`
	Members []*ThreadMember `json:"members"`
	HasMore bool            `json:"has_more"`
}

type AddedThreadMember struct {
	*ThreadMember
	Member   *Member   `json:"member"`
	Presence *Presence `json:"presence"`
}

type ForumDefaultReaction struct {
	EmojiID   string `json:"emoji_id,omitempty"`
	EmojiName string `json:"emoji_name,omitempty"`
}

type ForumTag struct {
	ID        string `json:"id,omitempty"`
	Name      string `json:"name"`
	Moderated bool   `json:"moderated"`
	EmojiID   string `json:"emoji_id,omitempty"`
	EmojiName string `json:"emoji_name,omitempty"`
}

type Emoji struct {
	ID            string   `json:"id"`
	Name          string   `json:"name"`
	Roles         []string `json:"roles"`
	User          *User    `json:"user"`
	RequireColons bool     `json:"require_colons"`
	Managed       bool     `json:"managed"`
	Animated      bool     `json:"animated"`
	Available     bool     `json:"available"`
}

var (
	EmojiRegex = regexp.MustCompile(`<(a|):[A-Za-z0-9_~]+:[0-9]{18,20}>`)
)

func (e *Emoji) MessageFormat() string {
	if e.ID != "" && e.Name != "" {
		if e.Animated {
			return "<a:" + e.APIName() + ">"
		}

		return "<:" + e.APIName() + ">"
	}

	return e.APIName()
}

func (e *Emoji) APIName() string {
	if e.ID != "" && e.Name != "" {
		return e.Name + ":" + e.ID
	}
	if e.Name != "" {
		return e.Name
	}
	return e.ID
}

type EmojiParams struct {
	Name  string   `json:"name,omitempty"`
	Image string   `json:"image,omitempty"`
	Roles []string `json:"roles,omitempty"`
}

type StickerFormat int

const (
	StickerFormatTypePNG    StickerFormat = 1
	StickerFormatTypeAPNG   StickerFormat = 2
	StickerFormatTypeLottie StickerFormat = 3
	StickerFormatTypeGIF    StickerFormat = 4
)

type StickerType int

const (
	StickerTypeStandard StickerType = 1
	StickerTypeGuild    StickerType = 2
)

type Sticker struct {
	ID          string        `json:"id"`
	PackID      string        `json:"pack_id"`
	Name        string        `json:"name"`
	Description string        `json:"description"`
	Tags        string        `json:"tags"`
	Type        StickerType   `json:"type"`
	FormatType  StickerFormat `json:"format_type"`
	Available   bool          `json:"available"`
	GuildID     string        `json:"guild_id"`
	User        *User         `json:"user"`
	SortValue   int           `json:"sort_value"`
}

type StickerItem struct {
	ID         string        `json:"id"`
	Name       string        `json:"name"`
	FormatType StickerFormat `json:"format_type"`
}

type StickerPack struct {
	ID             string     `json:"id"`
	Stickers       []*Sticker `json:"stickers"`
	Name           string     `json:"name"`
	SKUID          string     `json:"sku_id"`
	CoverStickerID string     `json:"cover_sticker_id"`
	Description    string     `json:"description"`
	BannerAssetID  string     `json:"banner_asset_id"`
}

type VerificationLevel int

const (
	VerificationLevelNone     VerificationLevel = 0
	VerificationLevelLow      VerificationLevel = 1
	VerificationLevelMedium   VerificationLevel = 2
	VerificationLevelHigh     VerificationLevel = 3
	VerificationLevelVeryHigh VerificationLevel = 4
)

type ExplicitContentFilterLevel int

const (
	ExplicitContentFilterDisabled            ExplicitContentFilterLevel = 0
	ExplicitContentFilterMembersWithoutRoles ExplicitContentFilterLevel = 1
	ExplicitContentFilterAllMembers          ExplicitContentFilterLevel = 2
)

type GuildNSFWLevel int

const (
	GuildNSFWLevelDefault       GuildNSFWLevel = 0
	GuildNSFWLevelExplicit      GuildNSFWLevel = 1
	GuildNSFWLevelSafe          GuildNSFWLevel = 2
	GuildNSFWLevelAgeRestricted GuildNSFWLevel = 3
)

type MfaLevel int

const (
	MfaLevelNone     MfaLevel = 0
	MfaLevelElevated MfaLevel = 1
)

type PremiumTier int

const (
	PremiumTierNone PremiumTier = 0
	PremiumTier1    PremiumTier = 1
	PremiumTier2    PremiumTier = 2
	PremiumTier3    PremiumTier = 3
)

type Guild struct {
	ID                          string                     `json:"id"`
	Name                        string                     `json:"name"`
	Icon                        string                     `json:"icon"`
	Region                      string                     `json:"region"`
	AfkChannelID                string                     `json:"afk_channel_id"`
	OwnerID                     string                     `json:"owner_id"`
	Owner                       bool                       `json:"owner"`
	JoinedAt                    time.Time                  `json:"joined_at"`
	DiscoverySplash             string                     `json:"discovery_splash"`
	Splash                      string                     `json:"splash"`
	AfkTimeout                  int                        `json:"afk_timeout"`
	MemberCount                 int                        `json:"member_count"`
	VerificationLevel           VerificationLevel          `json:"verification_level"`
	Large                       bool                       `json:"large"`
	DefaultMessageNotifications MessageNotifications       `json:"default_message_notifications"`
	Roles                       []*Role                    `json:"roles"`
	Emojis                      []*Emoji                   `json:"emojis"`
	Stickers                    []*Sticker                 `json:"stickers"`
	Members                     []*Member                  `json:"members"`
	Presences                   []*Presence                `json:"presences"`
	MaxPresences                int                        `json:"max_presences"`
	MaxMembers                  int                        `json:"max_members"`
	Channels                    []*Channel                 `json:"channels"`
	Threads                     []*Channel                 `json:"threads"`
	VoiceStates                 []*VoiceState              `json:"voice_states"`
	Unavailable                 bool                       `json:"unavailable"`
	ExplicitContentFilter       ExplicitContentFilterLevel `json:"explicit_content_filter"`
	NSFWLevel                   GuildNSFWLevel             `json:"nsfw_level"`
	Features                    []GuildFeature             `json:"features"`
	MfaLevel                    MfaLevel                   `json:"mfa_level"`
	ApplicationID               string                     `json:"application_id"`
	WidgetEnabled               bool                       `json:"widget_enabled"`
	WidgetChannelID             string                     `json:"widget_channel_id"`
	SystemChannelID             string                     `json:"system_channel_id"`
	SystemChannelFlags          SystemChannelFlag          `json:"system_channel_flags"`
	RulesChannelID              string                     `json:"rules_channel_id"`
	VanityURLCode               string                     `json:"vanity_url_code"`
	Description                 string                     `json:"description"`
	Banner                      string                     `json:"banner"`
	PremiumTier                 PremiumTier                `json:"premium_tier"`
	PremiumSubscriptionCount    int                        `json:"premium_subscription_count"`
	PreferredLocale             string                     `json:"preferred_locale"`
	PublicUpdatesChannelID      string                     `json:"public_updates_channel_id"`
	MaxVideoChannelUsers        int                        `json:"max_video_channel_users"`
	ApproximateMemberCount      int                        `json:"approximate_member_count"`
	ApproximatePresenceCount    int                        `json:"approximate_presence_count"`
	Permissions                 int64                      `json:"permissions,string"`
	StageInstances              []*StageInstance           `json:"stage_instances"`
}

type GuildPreview struct {
	ID                       string   `json:"id"`
	Name                     string   `json:"name"`
	Icon                     string   `json:"icon"`
	Splash                   string   `json:"splash"`
	DiscoverySplash          string   `json:"discovery_splash"`
	Emojis                   []*Emoji `json:"emojis"`
	Features                 []string `json:"features"`
	ApproximateMemberCount   int      `json:"approximate_member_count"`
	ApproximatePresenceCount int      `json:"approximate_presence_count"`
	Description              string   `json:"description"`
}

func (g *GuildPreview) IconURL(size string) string {
	return iconURL(g.Icon, EndpointGuildIcon(g.ID, g.Icon), EndpointGuildIconAnimated(g.ID, g.Icon), size)
}

type GuildScheduledEvent struct {
	ID                 string                            `json:"id"`
	GuildID            string                            `json:"guild_id"`
	ChannelID          string                            `json:"channel_id"`
	CreatorID          string                            `json:"creator_id"`
	Name               string                            `json:"name"`
	Description        string                            `json:"description"`
	ScheduledStartTime time.Time                         `json:"scheduled_start_time"`
	ScheduledEndTime   *time.Time                        `json:"scheduled_end_time"`
	PrivacyLevel       GuildScheduledEventPrivacyLevel   `json:"privacy_level"`
	Status             GuildScheduledEventStatus         `json:"status"`
	EntityType         GuildScheduledEventEntityType     `json:"entity_type"`
	EntityID           string                            `json:"entity_id"`
	EntityMetadata     GuildScheduledEventEntityMetadata `json:"entity_metadata"`
	Creator            *User                             `json:"creator"`
	UserCount          int                               `json:"user_count"`
	Image              string                            `json:"image"`
	RecurrenceRule     RecurrenceRule                    `json:"recurrence_rule,omitempty"`
}

type GuildScheduledEventParams struct {
	ChannelID          string                             `json:"channel_id,omitempty"`
	Name               string                             `json:"name,omitempty"`
	Description        string                             `json:"description,omitempty"`
	ScheduledStartTime *time.Time                         `json:"scheduled_start_time,omitempty"`
	ScheduledEndTime   *time.Time                         `json:"scheduled_end_time,omitempty"`
	PrivacyLevel       GuildScheduledEventPrivacyLevel    `json:"privacy_level,omitempty"`
	Status             GuildScheduledEventStatus          `json:"status,omitempty"`
	EntityType         GuildScheduledEventEntityType      `json:"entity_type,omitempty"`
	EntityMetadata     *GuildScheduledEventEntityMetadata `json:"entity_metadata,omitempty"`
	Image              string                             `json:"image,omitempty"`
}

func (p GuildScheduledEventParams) MarshalJSON() ([]byte, error) {
	type guildScheduledEventParams GuildScheduledEventParams

	if p.EntityType == GuildScheduledEventEntityTypeExternal && p.ChannelID == "" {
		return Marshal(struct {
			guildScheduledEventParams
			ChannelID json.RawMessage `json:"channel_id"`
		}{
			guildScheduledEventParams: guildScheduledEventParams(p),
			ChannelID:                 json.RawMessage("null"),
		})
	}

	return Marshal(guildScheduledEventParams(p))
}

type GuildScheduledEventEntityMetadata struct {
	Location string `json:"location"`
}

type GuildScheduledEventPrivacyLevel int

const (
	GuildScheduledEventPrivacyLevelGuildOnly GuildScheduledEventPrivacyLevel = 2
)

type GuildScheduledEventStatus int

const (
	GuildScheduledEventStatusScheduled GuildScheduledEventStatus = 1
	GuildScheduledEventStatusActive    GuildScheduledEventStatus = 2
	GuildScheduledEventStatusCompleted GuildScheduledEventStatus = 3
	GuildScheduledEventStatusCanceled  GuildScheduledEventStatus = 4
)

type GuildScheduledEventEntityType int

const (
	GuildScheduledEventEntityTypeStageInstance GuildScheduledEventEntityType = 1
	GuildScheduledEventEntityTypeVoice         GuildScheduledEventEntityType = 2
	GuildScheduledEventEntityTypeExternal      GuildScheduledEventEntityType = 3
)

type GuildScheduledEventUser struct {
	GuildScheduledEventID string  `json:"guild_scheduled_event_id"`
	User                  *User   `json:"user"`
	Member                *Member `json:"member"`
}

type GuildOnboardingMode int

const (
	GuildOnboardingModeDefault  GuildOnboardingMode = 0
	GuildOnboardingModeAdvanced GuildOnboardingMode = 1
)

type GuildOnboarding struct {
	GuildID           string                   `json:"guild_id,omitempty"`
	Prompts           *[]GuildOnboardingPrompt `json:"prompts,omitempty"`
	DefaultChannelIDs []string                 `json:"default_channel_ids,omitempty"`
	Enabled           *bool                    `json:"enabled,omitempty"`
	Mode              *GuildOnboardingMode     `json:"mode,omitempty"`
}

type GuildOnboardingPromptType int

const (
	GuildOnboardingPromptTypeMultipleChoice GuildOnboardingPromptType = 0
	GuildOnboardingPromptTypeDropdown       GuildOnboardingPromptType = 1
)

type GuildOnboardingPrompt struct {
	ID           string                        `json:"id,omitempty"`
	Type         GuildOnboardingPromptType     `json:"type"`
	Options      []GuildOnboardingPromptOption `json:"options"`
	Title        string                        `json:"title"`
	SingleSelect bool                          `json:"single_select"`
	Required     bool                          `json:"required"`
	InOnboarding bool                          `json:"in_onboarding"`
}

type GuildOnboardingPromptOption struct {
	ID            string   `json:"id,omitempty"`
	ChannelIDs    []string `json:"channel_ids"`
	RoleIDs       []string `json:"role_ids"`
	Emoji         *Emoji   `json:"emoji,omitempty"`
	Title         string   `json:"title"`
	Description   string   `json:"description"`
	EmojiID       string   `json:"emoji_id,omitempty"`
	EmojiName     string   `json:"emoji_name,omitempty"`
	EmojiAnimated *bool    `json:"emoji_animated,omitempty"`
}

type GuildTemplate struct {
	Code                  string    `json:"code"`
	Name                  string    `json:"name,omitempty"`
	Description           *string   `json:"description,omitempty"`
	UsageCount            int       `json:"usage_count"`
	CreatorID             string    `json:"creator_id"`
	Creator               *User     `json:"creator"`
	CreatedAt             time.Time `json:"created_at"`
	UpdatedAt             time.Time `json:"updated_at"`
	SourceGuildID         string    `json:"source_guild_id"`
	SerializedSourceGuild *Guild    `json:"serialized_source_guild"`
	IsDirty               bool      `json:"is_dirty"`
}

type GuildTemplateParams struct {
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
}

type MessageNotifications int

const (
	MessageNotificationsAllMessages  MessageNotifications = 0
	MessageNotificationsOnlyMentions MessageNotifications = 1
)

type SystemChannelFlag int

const (
	SystemChannelFlagsSuppressJoinNotifications          SystemChannelFlag = 1 << 0
	SystemChannelFlagsSuppressPremium                    SystemChannelFlag = 1 << 1
	SystemChannelFlagsSuppressGuildReminderNotifications SystemChannelFlag = 1 << 2
	SystemChannelFlagsSuppressJoinNotificationReplies    SystemChannelFlag = 1 << 3
)

func (g *Guild) IconURL(size string) string {
	return iconURL(g.Icon, EndpointGuildIcon(g.ID, g.Icon), EndpointGuildIconAnimated(g.ID, g.Icon), size)
}

func (g *Guild) BannerURL(size string) string {
	return bannerURL(g.Banner, EndpointGuildBanner(g.ID, g.Banner), EndpointGuildBannerAnimated(g.ID, g.Banner), size)
}

type UserGuild struct {
	ID                       string         `json:"id"`
	Name                     string         `json:"name"`
	Icon                     string         `json:"icon"`
	Owner                    bool           `json:"owner"`
	Permissions              int64          `json:"permissions,string"`
	Features                 []GuildFeature `json:"features"`
	ApproximateMemberCount   int            `json:"approximate_member_count"`
	ApproximatePresenceCount int            `json:"approximate_presence_count"`
}

type GuildFeature string

const (
	GuildFeatureAnimatedBanner                GuildFeature = "ANIMATED_BANNER"
	GuildFeatureAnimatedIcon                  GuildFeature = "ANIMATED_ICON"
	GuildFeatureAutoModeration                GuildFeature = "AUTO_MODERATION"
	GuildFeatureBanner                        GuildFeature = "BANNER"
	GuildFeatureCommunity                     GuildFeature = "COMMUNITY"
	GuildFeatureDiscoverable                  GuildFeature = "DISCOVERABLE"
	GuildFeatureFeaturable                    GuildFeature = "FEATURABLE"
	GuildFeatureInviteSplash                  GuildFeature = "INVITE_SPLASH"
	GuildFeatureMemberVerificationGateEnabled GuildFeature = "MEMBER_VERIFICATION_GATE_ENABLED"
	GuildFeatureMonetizationEnabled           GuildFeature = "MONETIZATION_ENABLED"
	GuildFeatureMoreStickers                  GuildFeature = "MORE_STICKERS"
	GuildFeatureNews                          GuildFeature = "NEWS"
	GuildFeaturePartnered                     GuildFeature = "PARTNERED"
	GuildFeaturePreviewEnabled                GuildFeature = "PREVIEW_ENABLED"
	GuildFeaturePrivateThreads                GuildFeature = "PRIVATE_THREADS"
	GuildFeatureRoleIcons                     GuildFeature = "ROLE_ICONS"
	GuildFeatureTicketedEventsEnabled         GuildFeature = "TICKETED_EVENTS_ENABLED"
	GuildFeatureVanityURL                     GuildFeature = "VANITY_URL"
	GuildFeatureVerified                      GuildFeature = "VERIFIED"
	GuildFeatureVipRegions                    GuildFeature = "VIP_REGIONS"
	GuildFeatureWelcomeScreenEnabled          GuildFeature = "WELCOME_SCREEN_ENABLED"
)

type GuildParams struct {
	Name                        string             `json:"name,omitempty"`
	Region                      string             `json:"region,omitempty"`
	VerificationLevel           *VerificationLevel `json:"verification_level,omitempty"`
	DefaultMessageNotifications int                `json:"default_message_notifications,omitempty"`
	ExplicitContentFilter       int                `json:"explicit_content_filter,omitempty"`
	AfkChannelID                string             `json:"afk_channel_id,omitempty"`
	AfkTimeout                  int                `json:"afk_timeout,omitempty"`
	Icon                        string             `json:"icon,omitempty"`
	OwnerID                     string             `json:"owner_id,omitempty"`
	Splash                      string             `json:"splash,omitempty"`
	DiscoverySplash             string             `json:"discovery_splash,omitempty"`
	Banner                      string             `json:"banner,omitempty"`
	SystemChannelID             string             `json:"system_channel_id,omitempty"`
	SystemChannelFlags          SystemChannelFlag  `json:"system_channel_flags,omitempty"`
	RulesChannelID              string             `json:"rules_channel_id,omitempty"`
	PublicUpdatesChannelID      string             `json:"public_updates_channel_id,omitempty"`
	PreferredLocale             Locale             `json:"preferred_locale,omitempty"`
	Features                    []GuildFeature     `json:"features,omitempty"`
	Description                 string             `json:"description,omitempty"`
	PremiumProgressBarEnabled   *bool              `json:"premium_progress_bar_enabled,omitempty"`
}

type Role struct {
	ID           string    `json:"id"`
	Name         string    `json:"name"`
	Managed      bool      `json:"managed"`
	Mentionable  bool      `json:"mentionable"`
	Hoist        bool      `json:"hoist"`
	Color        int       `json:"color"`
	Position     int       `json:"position"`
	Permissions  int64     `json:"permissions,string"`
	Icon         string    `json:"icon"`
	UnicodeEmoji string    `json:"unicode_emoji"`
	Flags        RoleFlags `json:"flags"`
}

type RoleFlags int

const (
	RoleFlagInPrompt RoleFlags = 1 << 0
)

func (r *Role) Mention() string {
	return fmt.Sprintf("<@&%s>", r.ID)
}

func (r *Role) IconURL(size string) string {
	if r.Icon == "" {
		return ""
	}

	URL := EndpointRoleIcon(r.ID, r.Icon)

	if size != "" {
		return URL + "?size=" + size
	}
	return URL
}

type RoleParams struct {
	Name         string  `json:"name,omitempty"`
	Color        *int    `json:"color,omitempty"`
	Hoist        *bool   `json:"hoist,omitempty"`
	Permissions  *int64  `json:"permissions,omitempty,string"`
	Mentionable  *bool   `json:"mentionable,omitempty"`
	UnicodeEmoji *string `json:"unicode_emoji,omitempty"`
	Icon         *string `json:"icon,omitempty"`
}

type Roles []*Role

func (r Roles) Len() int {
	return len(r)
}

func (r Roles) Less(i, j int) bool {
	return r[i].Position > r[j].Position
}

func (r Roles) Swap(i, j int) {
	r[i], r[j] = r[j], r[i]
}

type VoiceState struct {
	GuildID                 string     `json:"guild_id"`
	ChannelID               string     `json:"channel_id"`
	UserID                  string     `json:"user_id"`
	Member                  *Member    `json:"member"`
	SessionID               string     `json:"session_id"`
	Deaf                    bool       `json:"deaf"`
	Mute                    bool       `json:"mute"`
	SelfDeaf                bool       `json:"self_deaf"`
	SelfMute                bool       `json:"self_mute"`
	SelfStream              bool       `json:"self_stream"`
	SelfVideo               bool       `json:"self_video"`
	Suppress                bool       `json:"suppress"`
	RequestToSpeakTimestamp *time.Time `json:"request_to_speak_timestamp"`
}

type Presence struct {
	User         *User        `json:"user"`
	Status       Status       `json:"status"`
	Activities   []*Activity  `json:"activities"`
	Since        *int         `json:"since"`
	ClientStatus ClientStatus `json:"client_status"`
}

type TimeStamps struct {
	EndTimestamp   int64 `json:"end,omitempty"`
	StartTimestamp int64 `json:"start,omitempty"`
}

func (t *TimeStamps) UnmarshalJSON(data []byte) error {
	var aux struct {
		End   float64 `json:"end,omitempty"`
		Start float64 `json:"start,omitempty"`
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	t.EndTimestamp = int64(aux.End)
	t.StartTimestamp = int64(aux.Start)
	return nil
}

type Assets struct {
	LargeImageID string `json:"large_image,omitempty"`
	SmallImageID string `json:"small_image,omitempty"`
	LargeText    string `json:"large_text,omitempty"`
	SmallText    string `json:"small_text,omitempty"`
}

type MemberFlags int

const (
	MemberFlagDidRejoin            MemberFlags = 1 << 0
	MemberFlagCompletedOnboarding  MemberFlags = 1 << 1
	MemberFlagBypassesVerification MemberFlags = 1 << 2
	MemberFlagStartedOnboarding    MemberFlags = 1 << 3
)

type Member struct {
	GuildID                    string      `json:"guild_id"`
	JoinedAt                   time.Time   `json:"joined_at"`
	Nick                       string      `json:"nick"`
	Deaf                       bool        `json:"deaf"`
	Mute                       bool        `json:"mute"`
	Avatar                     string      `json:"avatar"`
	Banner                     string      `json:"banner"`
	User                       *User       `json:"user"`
	Roles                      []string    `json:"roles"`
	PremiumSince               *time.Time  `json:"premium_since"`
	Flags                      MemberFlags `json:"flags"`
	Pending                    bool        `json:"pending"`
	Permissions                int64       `json:"permissions,string"`
	CommunicationDisabledUntil *time.Time  `json:"communication_disabled_until"`
}

func (m *Member) Mention() string {
	return "<@!" + m.User.ID + ">"
}

func (m *Member) AvatarURL(size string) string {
	if m.Avatar == "" {
		return m.User.AvatarURL(size)
	}
	return avatarURL(m.Avatar, "", EndpointGuildMemberAvatar(m.GuildID, m.User.ID, m.Avatar),
		EndpointGuildMemberAvatarAnimated(m.GuildID, m.User.ID, m.Avatar), size)

}

func (m *Member) BannerURL(size string) string {
	if m.Banner == "" {
		return m.User.BannerURL(size)
	}
	return bannerURL(
		m.Banner,
		EndpointGuildMemberBanner(m.GuildID, m.User.ID, m.Banner),
		EndpointGuildMemberBannerAnimated(m.GuildID, m.User.ID, m.Banner),
		size,
	)
}

func (m *Member) DisplayName() string {
	if m.Nick != "" {
		return m.Nick
	}
	return m.User.DisplayName()
}

type ClientStatus struct {
	Desktop Status `json:"desktop"`
	Mobile  Status `json:"mobile"`
	Web     Status `json:"web"`
}

type Status string

const (
	StatusOnline       Status = "online"
	StatusIdle         Status = "idle"
	StatusDoNotDisturb Status = "dnd"
	StatusInvisible    Status = "invisible"
	StatusOffline      Status = "offline"
)

type TooManyRequests struct {
	Bucket     string        `json:"bucket"`
	Message    string        `json:"message"`
	RetryAfter time.Duration `json:"retry_after"`
}

func (t *TooManyRequests) UnmarshalJSON(data []byte) error {
	var aux struct {
		Bucket     string  `json:"bucket"`
		Message    string  `json:"message"`
		RetryAfter float64 `json:"retry_after"`
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	t.Bucket = aux.Bucket
	t.Message = aux.Message
	secs, frac := math.Modf(aux.RetryAfter)
	t.RetryAfter = time.Duration(secs)*time.Second + time.Duration(frac*1e3)*time.Millisecond

	return nil
}

type ReadState struct {
	MentionCount  int    `json:"mention_count"`
	LastMessageID string `json:"last_message_id"`
	ID            string `json:"id"`
}

type GuildRole struct {
	Role    *Role  `json:"role"`
	GuildID string `json:"guild_id"`
}

type GuildBan struct {
	Reason string `json:"reason"`
	User   *User  `json:"user"`
}

type AutoModerationRule struct {
	ID              string                         `json:"id,omitempty"`
	GuildID         string                         `json:"guild_id,omitempty"`
	Name            string                         `json:"name,omitempty"`
	CreatorID       string                         `json:"creator_id,omitempty"`
	EventType       AutoModerationRuleEventType    `json:"event_type,omitempty"`
	TriggerType     AutoModerationRuleTriggerType  `json:"trigger_type,omitempty"`
	TriggerMetadata *AutoModerationTriggerMetadata `json:"trigger_metadata,omitempty"`
	Actions         []AutoModerationAction         `json:"actions,omitempty"`
	Enabled         *bool                          `json:"enabled,omitempty"`
	ExemptRoles     *[]string                      `json:"exempt_roles,omitempty"`
	ExemptChannels  *[]string                      `json:"exempt_channels,omitempty"`
}

type AutoModerationRuleEventType int

const (
	AutoModerationEventMessageSend AutoModerationRuleEventType = 1
)

type AutoModerationRuleTriggerType int

const (
	AutoModerationEventTriggerKeyword       AutoModerationRuleTriggerType = 1
	AutoModerationEventTriggerHarmfulLink   AutoModerationRuleTriggerType = 2
	AutoModerationEventTriggerSpam          AutoModerationRuleTriggerType = 3
	AutoModerationEventTriggerKeywordPreset AutoModerationRuleTriggerType = 4
)

type AutoModerationKeywordPreset uint

const (
	AutoModerationKeywordPresetProfanity     AutoModerationKeywordPreset = 1
	AutoModerationKeywordPresetSexualContent AutoModerationKeywordPreset = 2
	AutoModerationKeywordPresetSlurs         AutoModerationKeywordPreset = 3
)

type AutoModerationTriggerMetadata struct {
	KeywordFilter     []string                      `json:"keyword_filter,omitempty"`
	RegexPatterns     []string                      `json:"regex_patterns,omitempty"`
	Presets           []AutoModerationKeywordPreset `json:"presets,omitempty"`
	AllowList         *[]string                     `json:"allow_list,omitempty"`
	MentionTotalLimit int                           `json:"mention_total_limit,omitempty"`
}

type AutoModerationActionType int

const (
	AutoModerationRuleActionBlockMessage     AutoModerationActionType = 1
	AutoModerationRuleActionSendAlertMessage AutoModerationActionType = 2
	AutoModerationRuleActionTimeout          AutoModerationActionType = 3
)

type AutoModerationActionMetadata struct {
	ChannelID     string `json:"channel_id,omitempty"`
	Duration      int    `json:"duration_seconds,omitempty"`
	CustomMessage string `json:"custom_message,omitempty"`
}

type AutoModerationAction struct {
	Type     AutoModerationActionType      `json:"type"`
	Metadata *AutoModerationActionMetadata `json:"metadata,omitempty"`
}

type GuildEmbed struct {
	Enabled   *bool  `json:"enabled,omitempty"`
	ChannelID string `json:"channel_id,omitempty"`
}

type GuildAuditLog struct {
	Webhooks        []*Webhook       `json:"webhooks,omitempty"`
	Users           []*User          `json:"users,omitempty"`
	AuditLogEntries []*AuditLogEntry `json:"audit_log_entries"`
	Integrations    []*Integration   `json:"integrations"`
}

type AuditLogEntry struct {
	TargetID   string            `json:"target_id"`
	Changes    []*AuditLogChange `json:"changes"`
	UserID     string            `json:"user_id"`
	ID         string            `json:"id"`
	ActionType *AuditLogAction   `json:"action_type"`
	Options    *AuditLogOptions  `json:"options"`
	Reason     string            `json:"reason"`
}

type AuditLogChange struct {
	NewValue interface{}        `json:"new_value"`
	OldValue interface{}        `json:"old_value"`
	Key      *AuditLogChangeKey `json:"key"`
}

type AuditLogChangeKey string

const (
	AuditLogChangeKeyAfkChannelID               AuditLogChangeKey = "afk_channel_id"
	AuditLogChangeKeyAfkTimeout                 AuditLogChangeKey = "afk_timeout"
	AuditLogChangeKeyAllow                      AuditLogChangeKey = "allow"
	AuditLogChangeKeyApplicationID              AuditLogChangeKey = "application_id"
	AuditLogChangeKeyArchived                   AuditLogChangeKey = "archived"
	AuditLogChangeKeyAsset                      AuditLogChangeKey = "asset"
	AuditLogChangeKeyAutoArchiveDuration        AuditLogChangeKey = "auto_archive_duration"
	AuditLogChangeKeyAvailable                  AuditLogChangeKey = "available"
	AuditLogChangeKeyAvatarHash                 AuditLogChangeKey = "avatar_hash"
	AuditLogChangeKeyBannerHash                 AuditLogChangeKey = "banner_hash"
	AuditLogChangeKeyBitrate                    AuditLogChangeKey = "bitrate"
	AuditLogChangeKeyChannelID                  AuditLogChangeKey = "channel_id"
	AuditLogChangeKeyCode                       AuditLogChangeKey = "code"
	AuditLogChangeKeyColor                      AuditLogChangeKey = "color"
	AuditLogChangeKeyCommunicationDisabledUntil AuditLogChangeKey = "communication_disabled_until"
	AuditLogChangeKeyDeaf                       AuditLogChangeKey = "deaf"
	AuditLogChangeKeyDefaultAutoArchiveDuration AuditLogChangeKey = "default_auto_archive_duration"
	AuditLogChangeKeyDefaultMessageNotification AuditLogChangeKey = "default_message_notifications"
	AuditLogChangeKeyDeny                       AuditLogChangeKey = "deny"
	AuditLogChangeKeyDescription                AuditLogChangeKey = "description"
	AuditLogChangeKeyDiscoverySplashHash        AuditLogChangeKey = "discovery_splash_hash"
	AuditLogChangeKeyEnableEmoticons            AuditLogChangeKey = "enable_emoticons"
	AuditLogChangeKeyEntityType                 AuditLogChangeKey = "entity_type"
	AuditLogChangeKeyExpireBehavior             AuditLogChangeKey = "expire_behavior"
	AuditLogChangeKeyExpireGracePeriod          AuditLogChangeKey = "expire_grace_period"
	AuditLogChangeKeyExplicitContentFilter      AuditLogChangeKey = "explicit_content_filter"
	AuditLogChangeKeyFormatType                 AuditLogChangeKey = "format_type"
	AuditLogChangeKeyGuildID                    AuditLogChangeKey = "guild_id"
	AuditLogChangeKeyHoist                      AuditLogChangeKey = "hoist"
	AuditLogChangeKeyIconHash                   AuditLogChangeKey = "icon_hash"
	AuditLogChangeKeyID                         AuditLogChangeKey = "id"
	AuditLogChangeKeyInvitable                  AuditLogChangeKey = "invitable"
	AuditLogChangeKeyInviterID                  AuditLogChangeKey = "inviter_id"
	AuditLogChangeKeyLocation                   AuditLogChangeKey = "location"
	AuditLogChangeKeyLocked                     AuditLogChangeKey = "locked"
	AuditLogChangeKeyMaxAge                     AuditLogChangeKey = "max_age"
	AuditLogChangeKeyMaxUses                    AuditLogChangeKey = "max_uses"
	AuditLogChangeKeyMentionable                AuditLogChangeKey = "mentionable"
	AuditLogChangeKeyMfaLevel                   AuditLogChangeKey = "mfa_level"
	AuditLogChangeKeyMute                       AuditLogChangeKey = "mute"
	AuditLogChangeKeyName                       AuditLogChangeKey = "name"
	AuditLogChangeKeyNick                       AuditLogChangeKey = "nick"
	AuditLogChangeKeyNSFW                       AuditLogChangeKey = "nsfw"
	AuditLogChangeKeyOwnerID                    AuditLogChangeKey = "owner_id"
	AuditLogChangeKeyPermissionOverwrite        AuditLogChangeKey = "permission_overwrites"
	AuditLogChangeKeyPermissions                AuditLogChangeKey = "permissions"
	AuditLogChangeKeyPosition                   AuditLogChangeKey = "position"
	AuditLogChangeKeyPreferredLocale            AuditLogChangeKey = "preferred_locale"
	AuditLogChangeKeyPrivacylevel               AuditLogChangeKey = "privacy_level"
	AuditLogChangeKeyPruneDeleteDays            AuditLogChangeKey = "prune_delete_days"
	AuditLogChangeKeyPublicUpdatesChannelID     AuditLogChangeKey = "public_updates_channel_id"
	AuditLogChangeKeyRateLimitPerUser           AuditLogChangeKey = "rate_limit_per_user"
	AuditLogChangeKeyRegion                     AuditLogChangeKey = "region"
	AuditLogChangeKeyRulesChannelID             AuditLogChangeKey = "rules_channel_id"
	AuditLogChangeKeySplashHash                 AuditLogChangeKey = "splash_hash"
	AuditLogChangeKeyStatus                     AuditLogChangeKey = "status"
	AuditLogChangeKeySystemChannelID            AuditLogChangeKey = "system_channel_id"
	AuditLogChangeKeyTags                       AuditLogChangeKey = "tags"
	AuditLogChangeKeyTemporary                  AuditLogChangeKey = "temporary"
	AuditLogChangeKeyTempoary                                     = AuditLogChangeKeyTemporary
	AuditLogChangeKeyTopic                      AuditLogChangeKey = "topic"
	AuditLogChangeKeyType                       AuditLogChangeKey = "type"
	AuditLogChangeKeyUnicodeEmoji               AuditLogChangeKey = "unicode_emoji"
	AuditLogChangeKeyUserLimit                  AuditLogChangeKey = "user_limit"
	AuditLogChangeKeyUses                       AuditLogChangeKey = "uses"
	AuditLogChangeKeyVanityURLCode              AuditLogChangeKey = "vanity_url_code"
	AuditLogChangeKeyVerificationLevel          AuditLogChangeKey = "verification_level"
	AuditLogChangeKeyWidgetChannelID            AuditLogChangeKey = "widget_channel_id"
	AuditLogChangeKeyWidgetEnabled              AuditLogChangeKey = "widget_enabled"
	AuditLogChangeKeyRoleAdd                    AuditLogChangeKey = "$add"
	AuditLogChangeKeyRoleRemove                 AuditLogChangeKey = "$remove"
)

type AuditLogOptions struct {
	DeleteMemberDays              string               `json:"delete_member_days"`
	MembersRemoved                string               `json:"members_removed"`
	ChannelID                     string               `json:"channel_id"`
	MessageID                     string               `json:"message_id"`
	Count                         string               `json:"count"`
	ID                            string               `json:"id"`
	Type                          *AuditLogOptionsType `json:"type"`
	RoleName                      string               `json:"role_name"`
	ApplicationID                 string               `json:"application_id"`
	AutoModerationRuleName        string               `json:"auto_moderation_rule_name"`
	AutoModerationRuleTriggerType string               `json:"auto_moderation_rule_trigger_type"`
	IntegrationType               string               `json:"integration_type"`
}

type AuditLogOptionsType string

const (
	AuditLogOptionsTypeRole   AuditLogOptionsType = "0"
	AuditLogOptionsTypeMember AuditLogOptionsType = "1"
)

type AuditLogAction int

const (
	AuditLogActionGuildUpdate                             AuditLogAction = 1
	AuditLogActionChannelCreate                           AuditLogAction = 10
	AuditLogActionChannelUpdate                           AuditLogAction = 11
	AuditLogActionChannelDelete                           AuditLogAction = 12
	AuditLogActionChannelOverwriteCreate                  AuditLogAction = 13
	AuditLogActionChannelOverwriteUpdate                  AuditLogAction = 14
	AuditLogActionChannelOverwriteDelete                  AuditLogAction = 15
	AuditLogActionMemberKick                              AuditLogAction = 20
	AuditLogActionMemberPrune                             AuditLogAction = 21
	AuditLogActionMemberBanAdd                            AuditLogAction = 22
	AuditLogActionMemberBanRemove                         AuditLogAction = 23
	AuditLogActionMemberUpdate                            AuditLogAction = 24
	AuditLogActionMemberRoleUpdate                        AuditLogAction = 25
	AuditLogActionMemberMove                              AuditLogAction = 26
	AuditLogActionMemberDisconnect                        AuditLogAction = 27
	AuditLogActionBotAdd                                  AuditLogAction = 28
	AuditLogActionRoleCreate                              AuditLogAction = 30
	AuditLogActionRoleUpdate                              AuditLogAction = 31
	AuditLogActionRoleDelete                              AuditLogAction = 32
	AuditLogActionInviteCreate                            AuditLogAction = 40
	AuditLogActionInviteUpdate                            AuditLogAction = 41
	AuditLogActionInviteDelete                            AuditLogAction = 42
	AuditLogActionWebhookCreate                           AuditLogAction = 50
	AuditLogActionWebhookUpdate                           AuditLogAction = 51
	AuditLogActionWebhookDelete                           AuditLogAction = 52
	AuditLogActionEmojiCreate                             AuditLogAction = 60
	AuditLogActionEmojiUpdate                             AuditLogAction = 61
	AuditLogActionEmojiDelete                             AuditLogAction = 62
	AuditLogActionMessageDelete                           AuditLogAction = 72
	AuditLogActionMessageBulkDelete                       AuditLogAction = 73
	AuditLogActionMessagePin                              AuditLogAction = 74
	AuditLogActionMessageUnpin                            AuditLogAction = 75
	AuditLogActionIntegrationCreate                       AuditLogAction = 80
	AuditLogActionIntegrationUpdate                       AuditLogAction = 81
	AuditLogActionIntegrationDelete                       AuditLogAction = 82
	AuditLogActionStageInstanceCreate                     AuditLogAction = 83
	AuditLogActionStageInstanceUpdate                     AuditLogAction = 84
	AuditLogActionStageInstanceDelete                     AuditLogAction = 85
	AuditLogActionStickerCreate                           AuditLogAction = 90
	AuditLogActionStickerUpdate                           AuditLogAction = 91
	AuditLogActionStickerDelete                           AuditLogAction = 92
	AuditLogGuildScheduledEventCreate                     AuditLogAction = 100
	AuditLogGuildScheduledEventUpdate                     AuditLogAction = 101
	AuditLogGuildScheduledEventDelete                     AuditLogAction = 102
	AuditLogActionThreadCreate                            AuditLogAction = 110
	AuditLogActionThreadUpdate                            AuditLogAction = 111
	AuditLogActionThreadDelete                            AuditLogAction = 112
	AuditLogActionApplicationCommandPermissionUpdate      AuditLogAction = 121
	AuditLogActionAutoModerationRuleCreate                AuditLogAction = 140
	AuditLogActionAutoModerationRuleUpdate                AuditLogAction = 141
	AuditLogActionAutoModerationRuleDelete                AuditLogAction = 142
	AuditLogActionAutoModerationBlockMessage              AuditLogAction = 143
	AuditLogActionAutoModerationFlagToChannel             AuditLogAction = 144
	AuditLogActionAutoModerationUserCommunicationDisabled AuditLogAction = 145
	AuditLogActionCreatorMonetizationRequestCreated       AuditLogAction = 150
	AuditLogActionCreatorMonetizationTermsAccepted        AuditLogAction = 151
	AuditLogActionOnboardingPromptCreate                  AuditLogAction = 163
	AuditLogActionOnboardingPromptUpdate                  AuditLogAction = 164
	AuditLogActionOnboardingPromptDelete                  AuditLogAction = 165
	AuditLogActionOnboardingCreate                        AuditLogAction = 166
	AuditLogActionOnboardingUpdate                        AuditLogAction = 167
	AuditLogActionHomeSettingsCreate                                     = 190
	AuditLogActionHomeSettingsUpdate                                     = 191
)

type GuildMemberParams struct {
	Nick                       string     `json:"nick,omitempty"`
	Roles                      *[]string  `json:"roles,omitempty"`
	ChannelID                  *string    `json:"channel_id,omitempty"`
	Mute                       *bool      `json:"mute,omitempty"`
	Deaf                       *bool      `json:"deaf,omitempty"`
	CommunicationDisabledUntil *time.Time `json:"communication_disabled_until,omitempty"`
}

func (p GuildMemberParams) MarshalJSON() (res []byte, err error) {
	type guildMemberParams GuildMemberParams
	v := struct {
		guildMemberParams
		ChannelID                  json.RawMessage `json:"channel_id,omitempty"`
		CommunicationDisabledUntil json.RawMessage `json:"communication_disabled_until,omitempty"`
	}{guildMemberParams: guildMemberParams(p)}

	if p.ChannelID != nil {
		if *p.ChannelID == "" {
			v.ChannelID = json.RawMessage(`null`)
		} else {
			res, err = json.Marshal(p.ChannelID)
			if err != nil {
				return
			}
			v.ChannelID = res
		}
	}

	if p.CommunicationDisabledUntil != nil {
		if p.CommunicationDisabledUntil.IsZero() {
			v.CommunicationDisabledUntil = json.RawMessage(`null`)
		} else {
			res, err = json.Marshal(p.CommunicationDisabledUntil)
			if err != nil {
				return
			}
			v.CommunicationDisabledUntil = res
		}
	}

	return json.Marshal(v)
}

type GuildMemberAddParams struct {
	AccessToken string   `json:"access_token"`
	Nick        string   `json:"nick,omitempty"`
	Roles       []string `json:"roles,omitempty"`
	Mute        bool     `json:"mute,omitempty"`
	Deaf        bool     `json:"deaf,omitempty"`
}

type APIErrorMessage struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type MessageReaction struct {
	UserID    string `json:"user_id"`
	MessageID string `json:"message_id"`
	Emoji     Emoji  `json:"emoji"`
	ChannelID string `json:"channel_id"`
	GuildID   string `json:"guild_id,omitempty"`
}

type GatewayBotResponse struct {
	URL               string             `json:"url"`
	Shards            int                `json:"shards"`
	SessionStartLimit SessionInformation `json:"session_start_limit"`
}

type SessionInformation struct {
	Total          int `json:"total,omitempty"`
	Remaining      int `json:"remaining,omitempty"`
	ResetAfter     int `json:"reset_after,omitempty"`
	MaxConcurrency int `json:"max_concurrency,omitempty"`
}

type GatewayStatusUpdate struct {
	Since  int      `json:"since"`
	Game   Activity `json:"game"`
	Status string   `json:"status"`
	AFK    bool     `json:"afk"`
}

type Activity struct {
	Name          string       `json:"name"`
	Type          ActivityType `json:"type"`
	URL           string       `json:"url,omitempty"`
	CreatedAt     time.Time    `json:"created_at"`
	ApplicationID string       `json:"application_id,omitempty"`
	State         string       `json:"state,omitempty"`
	Details       string       `json:"details,omitempty"`
	Timestamps    TimeStamps   `json:"timestamps,omitempty"`
	Emoji         Emoji        `json:"emoji,omitempty"`
	Party         Party        `json:"party,omitempty"`
	Assets        Assets       `json:"assets,omitempty"`
	Secrets       Secrets      `json:"secrets,omitempty"`
	Instance      bool         `json:"instance,omitempty"`
	Flags         int          `json:"flags,omitempty"`
	SyncID        string       `json:"sync_id,omitempty"`
}

func (a *Activity) UnmarshalJSON(data []byte) error {
	var aux struct {
		Name          string       `json:"name"`
		Type          ActivityType `json:"type"`
		URL           string       `json:"url,omitempty"`
		CreatedAt     int64        `json:"created_at"`
		ApplicationID json.Number  `json:"application_id,omitempty"`
		State         string       `json:"state,omitempty"`
		Details       string       `json:"details,omitempty"`
		Timestamps    TimeStamps   `json:"timestamps,omitempty"`
		Emoji         Emoji        `json:"emoji,omitempty"`
		Party         Party        `json:"party,omitempty"`
		Assets        Assets       `json:"assets,omitempty"`
		Secrets       Secrets      `json:"secrets,omitempty"`
		Instance      bool         `json:"instance,omitempty"`
		Flags         int          `json:"flags,omitempty"`
		SyncID        string       `json:"sync_id,omitempty"`
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	*a = Activity{
		Name:          aux.Name,
		Type:          aux.Type,
		URL:           aux.URL,
		CreatedAt:     time.Unix(0, aux.CreatedAt*int64(time.Millisecond)),
		ApplicationID: aux.ApplicationID.String(),
		State:         aux.State,
		Details:       aux.Details,
		Timestamps:    aux.Timestamps,
		Emoji:         aux.Emoji,
		Party:         aux.Party,
		Assets:        aux.Assets,
		Secrets:       aux.Secrets,
		Instance:      aux.Instance,
		Flags:         aux.Flags,
		SyncID:        aux.SyncID,
	}

	return nil
}

type Party struct {
	ID   string `json:"id,omitempty"`
	Size []int  `json:"size,omitempty"`
}

type Secrets struct {
	Join     string `json:"join,omitempty"`
	Spectate string `json:"spectate,omitempty"`
	Match    string `json:"match,omitempty"`
}

type ActivityType int

const (
	ActivityTypeGame      ActivityType = 0
	ActivityTypeStreaming ActivityType = 1
	ActivityTypeListening ActivityType = 2
	ActivityTypeWatching  ActivityType = 3
	ActivityTypeCustom    ActivityType = 4
	ActivityTypeCompeting ActivityType = 5
)

type Identify struct {
	Token          string              `json:"token"`
	Properties     IdentifyProperties  `json:"properties"`
	Compress       bool                `json:"compress"`
	LargeThreshold int                 `json:"large_threshold"`
	Shard          *[2]int             `json:"shard,omitempty"`
	Presence       GatewayStatusUpdate `json:"presence,omitempty"`
	Intents        Intent              `json:"intents"`
}

type IdentifyProperties struct {
	OS              string `json:"$os"`
	Browser         string `json:"$browser"`
	Device          string `json:"$device"`
	Referer         string `json:"$referer"`
	ReferringDomain string `json:"$referring_domain"`
}

type StageInstance struct {
	ID                    string                    `json:"id"`
	GuildID               string                    `json:"guild_id"`
	ChannelID             string                    `json:"channel_id"`
	Topic                 string                    `json:"topic"`
	PrivacyLevel          StageInstancePrivacyLevel `json:"privacy_level"`
	DiscoverableDisabled  bool                      `json:"discoverable_disabled"`
	GuildScheduledEventID string                    `json:"guild_scheduled_event_id"`
}

type StageInstanceParams struct {
	ChannelID             string                    `json:"channel_id,omitempty"`
	Topic                 string                    `json:"topic,omitempty"`
	PrivacyLevel          StageInstancePrivacyLevel `json:"privacy_level,omitempty"`
	SendStartNotification bool                      `json:"send_start_notification,omitempty"`
}

type StageInstancePrivacyLevel int

const (
	StageInstancePrivacyLevelPublic    StageInstancePrivacyLevel = 1
	StageInstancePrivacyLevelGuildOnly StageInstancePrivacyLevel = 2
)

type PollLayoutType int

const (
	PollLayoutTypeDefault PollLayoutType = 1
)

type PollMedia struct {
	Text  string          `json:"text,omitempty"`
	Emoji *ComponentEmoji `json:"emoji,omitempty"`
}

type PollAnswer struct {
	AnswerID int        `json:"answer_id,omitempty"`
	Media    *PollMedia `json:"poll_media"`
}

type PollAnswerCount struct {
	ID      int  `json:"id"`
	Count   int  `json:"count"`
	MeVoted bool `json:"me_voted"`
}

type PollResults struct {
	Finalized    bool               `json:"is_finalized"`
	AnswerCounts []*PollAnswerCount `json:"answer_counts"`
}

type Poll struct {
	Question         PollMedia      `json:"question"`
	Answers          []PollAnswer   `json:"answers"`
	AllowMultiselect bool           `json:"allow_multiselect"`
	LayoutType       PollLayoutType `json:"layout_type,omitempty"`
	Duration         int            `json:"duration,omitempty"`
	Results          *PollResults   `json:"results,omitempty"`
	Expiry           *time.Time     `json:"expiry,omitempty"`
}

type SKUType int

const (
	SKUTypeDurable           SKUType = 2
	SKUTypeConsumable        SKUType = 3
	SKUTypeSubscription      SKUType = 5
	SKUTypeSubscriptionGroup SKUType = 6
)

type SKUFlags int

const (
	SKUFlagAvailable         SKUFlags = 1 << 2
	SKUFlagGuildSubscription SKUFlags = 1 << 7
	SKUFlagUserSubscription  SKUFlags = 1 << 8
)

type SKU struct {
	ID            string   `json:"id"`
	Type          SKUType  `json:"type"`
	ApplicationID string   `json:"application_id"`
	Name          string   `json:"name"`
	Slug          string   `json:"slug"`
	Flags         SKUFlags `json:"flags"`
}

type Subscription struct {
	ID                 string             `json:"id"`
	UserID             string             `json:"user_id"`
	SKUIDs             []string           `json:"sku_ids"`
	EntitlementIDs     []string           `json:"entitlement_ids"`
	RenewalSKUIDs      []string           `json:"renewal_sku_ids,omitempty"`
	CurrentPeriodStart time.Time          `json:"current_period_start"`
	CurrentPeriodEnd   time.Time          `json:"current_period_end"`
	Status             SubscriptionStatus `json:"status"`
	CanceledAt         *time.Time         `json:"canceled_at,omitempty"`
	Country            string             `json:"country,omitempty"`
}

type SubscriptionStatus int

const (
	SubscriptionStatusActive   = 0
	SubscriptionStatusEnding   = 1
	SubscriptionStatusInactive = 2
)

type EntitlementType int

const (
	EntitlementTypePurchase                = 1
	EntitlementTypePremiumSubscription     = 2
	EntitlementTypeDeveloperGift           = 3
	EntitlementTypeTestModePurchase        = 4
	EntitlementTypeFreePurchase            = 5
	EntitlementTypeUserGift                = 6
	EntitlementTypePremiumPurchase         = 7
	EntitlementTypeApplicationSubscription = 8
)

type Entitlement struct {
	ID             string          `json:"id"`
	SKUID          string          `json:"sku_id"`
	ApplicationID  string          `json:"application_id"`
	UserID         string          `json:"user_id,omitempty"`
	Type           EntitlementType `json:"type"`
	Deleted        bool            `json:"deleted"`
	StartsAt       *time.Time      `json:"starts_at,omitempty"`
	EndsAt         *time.Time      `json:"ends_at,omitempty"`
	GuildID        string          `json:"guild_id,omitempty"`
	Consumed       *bool           `json:"consumed,omitempty"`
	SubscriptionID string          `json:"subscription_id,omitempty"`
}

type EntitlementOwnerType int

const (
	EntitlementOwnerTypeGuildSubscription EntitlementOwnerType = 1
	EntitlementOwnerTypeUserSubscription  EntitlementOwnerType = 2
)

type EntitlementTest struct {
	SKUID     string               `json:"sku_id"`
	OwnerID   string               `json:"owner_id"`
	OwnerType EntitlementOwnerType `json:"owner_type"`
}

type EntitlementFilterOptions struct {
	UserID       string
	SkuIDs       []string
	Before       *time.Time
	After        *time.Time
	Limit        int
	GuildID      string
	ExcludeEnded bool
}

const (
	PermissionReadMessages           = 1 << 10
	PermissionSendMessages           = 1 << 11
	PermissionSendTTSMessages        = 1 << 12
	PermissionManageMessages         = 1 << 13
	PermissionEmbedLinks             = 1 << 14
	PermissionAttachFiles            = 1 << 15
	PermissionReadMessageHistory     = 1 << 16
	PermissionMentionEveryone        = 1 << 17
	PermissionUseExternalEmojis      = 1 << 18
	PermissionUseSlashCommands       = 1 << 31
	PermissionUseApplicationCommands = 1 << 31
	PermissionManageThreads          = 1 << 34
	PermissionCreatePublicThreads    = 1 << 35
	PermissionCreatePrivateThreads   = 1 << 36
	PermissionUseExternalStickers    = 1 << 37
	PermissionSendMessagesInThreads  = 1 << 38
	PermissionSendVoiceMessages      = 1 << 46
	PermissionSendPolls              = 1 << 49
	PermissionUseExternalApps        = 1 << 50
)

const (
	PermissionVoicePrioritySpeaker  = 1 << 8
	PermissionVoiceStreamVideo      = 1 << 9
	PermissionVoiceConnect          = 1 << 20
	PermissionVoiceSpeak            = 1 << 21
	PermissionVoiceMuteMembers      = 1 << 22
	PermissionVoiceDeafenMembers    = 1 << 23
	PermissionVoiceMoveMembers      = 1 << 24
	PermissionVoiceUseVAD           = 1 << 25
	PermissionVoiceRequestToSpeak   = 1 << 32
	PermissionUseActivities         = 1 << 39
	PermissionUseEmbeddedActivities = 1 << 39
	PermissionUseSoundboard         = 1 << 42
	PermissionUseExternalSounds     = 1 << 45
)
const (
	PermissionChangeNickname                   = 1 << 26
	PermissionManageNicknames                  = 1 << 27
	PermissionManageRoles                      = 1 << 28
	PermissionManageWebhooks                   = 1 << 29
	PermissionManageEmojis                     = 1 << 30
	PermissionManageGuildExpressions           = 1 << 30
	PermissionManageEvents                     = 1 << 33
	PermissionViewCreatorMonetizationAnalytics = 1 << 41
	PermissionCreateGuildExpressions           = 1 << 43
	PermissionCreateEvents                     = 1 << 44
)

const (
	PermissionCreateInstantInvite = 1 << 0
	PermissionKickMembers         = 1 << 1
	PermissionBanMembers          = 1 << 2
	PermissionAdministrator       = 1 << 3
	PermissionManageChannels      = 1 << 4
	PermissionManageServer        = 1 << 5
	PermissionManageGuild         = 1 << 5
	PermissionAddReactions        = 1 << 6
	PermissionViewAuditLogs       = 1 << 7
	PermissionViewChannel         = 1 << 10
	PermissionViewGuildInsights   = 1 << 19
	PermissionModerateMembers     = 1 << 40
	PermissionAllText             = PermissionViewChannel |
		PermissionSendMessages |
		PermissionSendTTSMessages |
		PermissionManageMessages |
		PermissionEmbedLinks |
		PermissionAttachFiles |
		PermissionReadMessageHistory |
		PermissionMentionEveryone
	PermissionAllVoice = PermissionViewChannel |
		PermissionVoiceConnect |
		PermissionVoiceSpeak |
		PermissionVoiceMuteMembers |
		PermissionVoiceDeafenMembers |
		PermissionVoiceMoveMembers |
		PermissionVoiceUseVAD |
		PermissionVoicePrioritySpeaker
	PermissionAllChannel = PermissionAllText |
		PermissionAllVoice |
		PermissionCreateInstantInvite |
		PermissionManageRoles |
		PermissionManageChannels |
		PermissionAddReactions |
		PermissionViewAuditLogs
	PermissionAll = PermissionAllChannel |
		PermissionKickMembers |
		PermissionBanMembers |
		PermissionManageServer |
		PermissionAdministrator |
		PermissionManageWebhooks |
		PermissionManageEmojis
)

const (
	ErrCodeGeneralError                                                     = 0
	ErrCodeUnknownAccount                                                   = 10001
	ErrCodeUnknownApplication                                               = 10002
	ErrCodeUnknownChannel                                                   = 10003
	ErrCodeUnknownGuild                                                     = 10004
	ErrCodeUnknownIntegration                                               = 10005
	ErrCodeUnknownInvite                                                    = 10006
	ErrCodeUnknownMember                                                    = 10007
	ErrCodeUnknownMessage                                                   = 10008
	ErrCodeUnknownOverwrite                                                 = 10009
	ErrCodeUnknownProvider                                                  = 10010
	ErrCodeUnknownRole                                                      = 10011
	ErrCodeUnknownToken                                                     = 10012
	ErrCodeUnknownUser                                                      = 10013
	ErrCodeUnknownEmoji                                                     = 10014
	ErrCodeUnknownWebhook                                                   = 10015
	ErrCodeUnknownWebhookService                                            = 10016
	ErrCodeUnknownSession                                                   = 10020
	ErrCodeUnknownBan                                                       = 10026
	ErrCodeUnknownSKU                                                       = 10027
	ErrCodeUnknownStoreListing                                              = 10028
	ErrCodeUnknownEntitlement                                               = 10029
	ErrCodeUnknownBuild                                                     = 10030
	ErrCodeUnknownLobby                                                     = 10031
	ErrCodeUnknownBranch                                                    = 10032
	ErrCodeUnknownStoreDirectoryLayout                                      = 10033
	ErrCodeUnknownRedistributable                                           = 10036
	ErrCodeUnknownGiftCode                                                  = 10038
	ErrCodeUnknownStream                                                    = 10049
	ErrCodeUnknownPremiumServerSubscribeCooldown                            = 10050
	ErrCodeUnknownGuildTemplate                                             = 10057
	ErrCodeUnknownDiscoveryCategory                                         = 10059
	ErrCodeUnknownSticker                                                   = 10060
	ErrCodeUnknownInteraction                                               = 10062
	ErrCodeUnknownApplicationCommand                                        = 10063
	ErrCodeUnknownApplicationCommandPermissions                             = 10066
	ErrCodeUnknownStageInstance                                             = 10067
	ErrCodeUnknownGuildMemberVerificationForm                               = 10068
	ErrCodeUnknownGuildWelcomeScreen                                        = 10069
	ErrCodeUnknownGuildScheduledEvent                                       = 10070
	ErrCodeUnknownGuildScheduledEventUser                                   = 10071
	ErrUnknownTag                                                           = 10087
	ErrCodeBotsCannotUseEndpoint                                            = 20001
	ErrCodeOnlyBotsCanUseEndpoint                                           = 20002
	ErrCodeExplicitContentCannotBeSentToTheDesiredRecipients                = 20009
	ErrCodeYouAreNotAuthorizedToPerformThisActionOnThisApplication          = 20012
	ErrCodeThisActionCannotBePerformedDueToSlowmodeRateLimit                = 20016
	ErrCodeOnlyTheOwnerOfThisAccountCanPerformThisAction                    = 20018
	ErrCodeMessageCannotBeEditedDueToAnnouncementRateLimits                 = 20022
	ErrCodeChannelHasHitWriteRateLimit                                      = 20028
	ErrCodeTheWriteActionYouArePerformingOnTheServerHasHitTheWriteRateLimit = 20029
	ErrCodeStageTopicContainsNotAllowedWordsForPublicStages                 = 20031
	ErrCodeGuildPremiumSubscriptionLevelTooLow                              = 20035
	ErrCodeMaximumGuildsReached                                             = 30001
	ErrCodeMaximumPinsReached                                               = 30003
	ErrCodeMaximumNumberOfRecipientsReached                                 = 30004
	ErrCodeMaximumGuildRolesReached                                         = 30005
	ErrCodeMaximumNumberOfWebhooksReached                                   = 30007
	ErrCodeMaximumNumberOfEmojisReached                                     = 30008
	ErrCodeTooManyReactions                                                 = 30010
	ErrCodeMaximumNumberOfGuildChannelsReached                              = 30013
	ErrCodeMaximumNumberOfAttachmentsInAMessageReached                      = 30015
	ErrCodeMaximumNumberOfInvitesReached                                    = 30016
	ErrCodeMaximumNumberOfAnimatedEmojisReached                             = 30018
	ErrCodeMaximumNumberOfServerMembersReached                              = 30019
	ErrCodeMaximumNumberOfGuildDiscoverySubcategoriesReached                = 30030
	ErrCodeGuildAlreadyHasATemplate                                         = 30031
	ErrCodeMaximumNumberOfThreadParticipantsReached                         = 30033
	ErrCodeMaximumNumberOfBansForNonGuildMembersHaveBeenExceeded            = 30035
	ErrCodeMaximumNumberOfBansFetchesHasBeenReached                         = 30037
	ErrCodeMaximumNumberOfUncompletedGuildScheduledEventsReached            = 30038
	ErrCodeMaximumNumberOfStickersReached                                   = 30039
	ErrCodeMaximumNumberOfPruneRequestsHasBeenReached                       = 30040
	ErrCodeMaximumNumberOfGuildWidgetSettingsUpdatesHasBeenReached          = 30042
	ErrCodeMaximumNumberOfEditsToMessagesOlderThanOneHourReached            = 30046
	ErrCodeMaximumNumberOfPinnedThreadsInForumChannelHasBeenReached         = 30047
	ErrCodeMaximumNumberOfTagsInForumChannelHasBeenReached                  = 30048
	ErrCodeUnauthorized                                                     = 40001
	ErrCodeActionRequiredVerifiedAccount                                    = 40002
	ErrCodeOpeningDirectMessagesTooFast                                     = 40003
	ErrCodeSendMessagesHasBeenTemporarilyDisabled                           = 40004
	ErrCodeRequestEntityTooLarge                                            = 40005
	ErrCodeFeatureTemporarilyDisabledServerSide                             = 40006
	ErrCodeUserIsBannedFromThisGuild                                        = 40007
	ErrCodeTargetIsNotConnectedToVoice                                      = 40032
	ErrCodeMessageAlreadyCrossposted                                        = 40033
	ErrCodeAnApplicationWithThatNameAlreadyExists                           = 40041
	ErrCodeInteractionHasAlreadyBeenAcknowledged                            = 40060
	ErrCodeTagNamesMustBeUnique                                             = 40061
	ErrCodeMissingAccess                                                    = 50001
	ErrCodeInvalidAccountType                                               = 50002
	ErrCodeCannotExecuteActionOnDMChannel                                   = 50003
	ErrCodeEmbedDisabled                                                    = 50004
	ErrCodeGuildWidgetDisabled                                              = 50004
	ErrCodeCannotEditFromAnotherUser                                        = 50005
	ErrCodeCannotSendEmptyMessage                                           = 50006
	ErrCodeCannotSendMessagesToThisUser                                     = 50007
	ErrCodeCannotSendMessagesInVoiceChannel                                 = 50008
	ErrCodeChannelVerificationLevelTooHigh                                  = 50009
	ErrCodeOAuth2ApplicationDoesNotHaveBot                                  = 50010
	ErrCodeOAuth2ApplicationLimitReached                                    = 50011
	ErrCodeInvalidOAuthState                                                = 50012
	ErrCodeMissingPermissions                                               = 50013
	ErrCodeInvalidAuthenticationToken                                       = 50014
	ErrCodeTooFewOrTooManyMessagesToDelete                                  = 50016
	ErrCodeCanOnlyPinMessageToOriginatingChannel                            = 50019
	ErrCodeInviteCodeWasEitherInvalidOrTaken                                = 50020
	ErrCodeCannotExecuteActionOnSystemMessage                               = 50021
	ErrCodeCannotExecuteActionOnThisChannelType                             = 50024
	ErrCodeInvalidOAuth2AccessTokenProvided                                 = 50025
	ErrCodeMissingRequiredOAuth2Scope                                       = 50026
	ErrCodeInvalidWebhookTokenProvided                                      = 50027
	ErrCodeInvalidRole                                                      = 50028
	ErrCodeInvalidRecipients                                                = 50033
	ErrCodeMessageProvidedTooOldForBulkDelete                               = 50034
	ErrCodeInvalidFormBody                                                  = 50035
	ErrCodeInviteAcceptedToGuildApplicationsBotNotIn                        = 50036
	ErrCodeInvalidAPIVersionProvided                                        = 50041
	ErrCodeFileUploadedExceedsTheMaximumSize                                = 50045
	ErrCodeInvalidFileUploaded                                              = 50046
	ErrCodeInvalidGuild                                                     = 50055
	ErrCodeInvalidMessageType                                               = 50068
	ErrCodeCannotDeleteAChannelRequiredForCommunityGuilds                   = 50074
	ErrCodeInvalidStickerSent                                               = 50081
	ErrCodePerformedOperationOnArchivedThread                               = 50083
	ErrCodeBeforeValueIsEarlierThanThreadCreationDate                       = 50085
	ErrCodeCommunityServerChannelsMustBeTextChannels                        = 50086
	ErrCodeThisServerIsNotAvailableInYourLocation                           = 50095
	ErrCodeThisServerNeedsMonetizationEnabledInOrderToPerformThisAction     = 50097
	ErrCodeThisServerNeedsMoreBoostsToPerformThisAction                     = 50101
	ErrCodeTheRequestBodyContainsInvalidJSON                                = 50109
	ErrCodeNoUsersWithDiscordTagExist                                       = 80004
	ErrCodeReactionBlocked                                                  = 90001
	ErrCodeAPIResourceIsCurrentlyOverloaded                                 = 130000
	ErrCodeTheStageIsAlreadyOpen                                            = 150006
	ErrCodeCannotReplyWithoutPermissionToReadMessageHistory                 = 160002
	ErrCodeThreadAlreadyCreatedForThisMessage                               = 160004
	ErrCodeThreadIsLocked                                                   = 160005
	ErrCodeMaximumNumberOfActiveThreadsReached                              = 160006
	ErrCodeMaximumNumberOfActiveAnnouncementThreadsReached                  = 160007
	ErrCodeInvalidJSONForUploadedLottieFile                                 = 170001
	ErrCodeUploadedLottiesCannotContainRasterizedImages                     = 170002
	ErrCodeStickerMaximumFramerateExceeded                                  = 170003
	ErrCodeStickerFrameCountExceedsMaximumOfOneThousandFrames               = 170004
	ErrCodeLottieAnimationMaximumDimensionsExceeded                         = 170005
	ErrCodeStickerFrameRateOutOfRange                                       = 170006
	ErrCodeStickerAnimationDurationExceedsMaximumOfFiveSeconds              = 170007
	ErrCodeCannotUpdateAFinishedEvent                                       = 180000
	ErrCodeFailedToCreateStageNeededForStageEvent                           = 180002
	ErrCodeCannotEnableOnboardingRequirementsAreNotMet                      = 350000
	ErrCodeCannotUpdateOnboardingWhileBelowRequirements                     = 350001
)

type Intent int

const (
	IntentGuilds                      Intent = 1 << 0
	IntentGuildMembers                Intent = 1 << 1
	IntentGuildModeration             Intent = 1 << 2
	IntentGuildEmojis                 Intent = 1 << 3
	IntentGuildIntegrations           Intent = 1 << 4
	IntentGuildWebhooks               Intent = 1 << 5
	IntentGuildInvites                Intent = 1 << 6
	IntentGuildVoiceStates            Intent = 1 << 7
	IntentGuildPresences              Intent = 1 << 8
	IntentGuildMessages               Intent = 1 << 9
	IntentGuildMessageReactions       Intent = 1 << 10
	IntentGuildMessageTyping          Intent = 1 << 11
	IntentDirectMessages              Intent = 1 << 12
	IntentDirectMessageReactions      Intent = 1 << 13
	IntentDirectMessageTyping         Intent = 1 << 14
	IntentMessageContent              Intent = 1 << 15
	IntentGuildScheduledEvents        Intent = 1 << 16
	IntentAutoModerationConfiguration Intent = 1 << 20
	IntentAutoModerationExecution     Intent = 1 << 21
	IntentGuildMessagePolls           Intent = 1 << 24
	IntentDirectMessagePolls          Intent = 1 << 25
	IntentGuildBans                   Intent = IntentGuildModeration
	IntentsGuilds                     Intent = 1 << 0
	IntentsGuildMembers               Intent = 1 << 1
	IntentsGuildBans                  Intent = 1 << 2
	IntentsGuildEmojis                Intent = 1 << 3
	IntentsGuildIntegrations          Intent = 1 << 4
	IntentsGuildWebhooks              Intent = 1 << 5
	IntentsGuildInvites               Intent = 1 << 6
	IntentsGuildVoiceStates           Intent = 1 << 7
	IntentsGuildPresences             Intent = 1 << 8
	IntentsGuildMessages              Intent = 1 << 9
	IntentsGuildMessageReactions      Intent = 1 << 10
	IntentsGuildMessageTyping         Intent = 1 << 11
	IntentsDirectMessages             Intent = 1 << 12
	IntentsDirectMessageReactions     Intent = 1 << 13
	IntentsDirectMessageTyping        Intent = 1 << 14
	IntentsMessageContent             Intent = 1 << 15
	IntentsGuildScheduledEvents       Intent = 1 << 16

	IntentsAllWithoutPrivileged = IntentGuilds |
		IntentGuildBans |
		IntentGuildEmojis |
		IntentGuildIntegrations |
		IntentGuildWebhooks |
		IntentGuildInvites |
		IntentGuildVoiceStates |
		IntentGuildMessages |
		IntentGuildMessageReactions |
		IntentGuildMessageTyping |
		IntentDirectMessages |
		IntentDirectMessageReactions |
		IntentDirectMessageTyping |
		IntentGuildScheduledEvents |
		IntentAutoModerationConfiguration |
		IntentAutoModerationExecution

	IntentsAll = IntentsAllWithoutPrivileged |
		IntentGuildMembers |
		IntentGuildPresences |
		IntentMessageContent

	IntentsNone Intent = 0
)

func MakeIntent(intents Intent) Intent {
	return intents
}
