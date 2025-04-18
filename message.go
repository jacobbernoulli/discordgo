package discordgo

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"regexp"
	"runtime"
	"strings"
	"time"
)

type MessageType int

const (
	MessageTypeDefault                               MessageType = 0
	MessageTypeRecipientAdd                          MessageType = 1
	MessageTypeRecipientRemove                       MessageType = 2
	MessageTypeCall                                  MessageType = 3
	MessageTypeChannelNameChange                     MessageType = 4
	MessageTypeChannelIconChange                     MessageType = 5
	MessageTypeChannelPinnedMessage                  MessageType = 6
	MessageTypeGuildMemberJoin                       MessageType = 7
	MessageTypeUserPremiumGuildSubscription          MessageType = 8
	MessageTypeUserPremiumGuildSubscriptionTierOne   MessageType = 9
	MessageTypeUserPremiumGuildSubscriptionTierTwo   MessageType = 10
	MessageTypeUserPremiumGuildSubscriptionTierThree MessageType = 11
	MessageTypeChannelFollowAdd                      MessageType = 12
	MessageTypeGuildDiscoveryDisqualified            MessageType = 14
	MessageTypeGuildDiscoveryRequalified             MessageType = 15
	MessageTypeThreadCreated                         MessageType = 18
	MessageTypeReply                                 MessageType = 19
	MessageTypeChatInputCommand                      MessageType = 20
	MessageTypeThreadStarterMessage                  MessageType = 21
	MessageTypeContextMenuCommand                    MessageType = 23
)

type Message struct {
	ID                  string                      `json:"id"`
	ChannelID           string                      `json:"channel_id"`
	GuildID             string                      `json:"guild_id,omitempty"`
	Content             string                      `json:"content"`
	Timestamp           time.Time                   `json:"timestamp"`
	EditedTimestamp     *time.Time                  `json:"edited_timestamp"`
	MentionRoles        []string                    `json:"mention_roles"`
	TTS                 bool                        `json:"tts"`
	MentionEveryone     bool                        `json:"mention_everyone"`
	Author              *User                       `json:"author"`
	Attachments         []*MessageAttachment        `json:"attachments"`
	Components          []MessageComponent          `json:"-"`
	Embeds              []*MessageEmbed             `json:"embeds"`
	Mentions            []*User                     `json:"mentions"`
	Reactions           []*MessageReactions         `json:"reactions"`
	Pinned              bool                        `json:"pinned"`
	Type                MessageType                 `json:"type"`
	WebhookID           string                      `json:"webhook_id"`
	Member              *Member                     `json:"member"`
	MentionChannels     []*Channel                  `json:"mention_channels"`
	Activity            *MessageActivity            `json:"activity"`
	Application         *MessageApplication         `json:"application"`
	MessageReference    *MessageReference           `json:"message_reference"`
	ReferencedMessage   *Message                    `json:"referenced_message"`
	MessageSnapshots    []MessageSnapshot           `json:"message_snapshots"`
	Interaction         *MessageInteraction         `json:"interaction"`
	InteractionMetadata *MessageInteractionMetadata `json:"interaction_metadata"`
	Flags               MessageFlags                `json:"flags"`
	Thread              *Channel                    `json:"thread,omitempty"`
	StickerItems        []*StickerItem              `json:"sticker_items"`
	Poll                *Poll                       `json:"poll"`
}

func (m *Message) UnmarshalJSON(data []byte) error {
	type Alias Message
	var temp struct {
		Alias
		RawComponents []unmarshalableMessageComponent `json:"components"`
	}

	if err := json.Unmarshal(data, &temp); err != nil {
		return err
	}

	*m = Message(temp.Alias)

	if len(temp.RawComponents) > 0 {
		m.Components = make([]MessageComponent, len(temp.RawComponents))
		for i, component := range temp.RawComponents {
			m.Components[i] = component.MessageComponent
		}
	}

	return nil
}

func (m *Message) GetCustomEmojis() []*Emoji {
	var toReturn []*Emoji
	emojis := EmojiRegex.FindAllString(m.Content, -1)
	if len(emojis) < 1 {
		return toReturn
	}
	for _, em := range emojis {
		parts := strings.Split(em, ":")
		toReturn = append(toReturn, &Emoji{
			ID:       parts[2][:len(parts[2])-1],
			Name:     parts[1],
			Animated: strings.HasPrefix(em, "<a:"),
		})
	}
	return toReturn
}

type MessageFlags int

const (
	MessageFlagsCrossPosted                      MessageFlags = 1 << 0
	MessageFlagsIsCrossPosted                    MessageFlags = 1 << 1
	MessageFlagsSuppressEmbeds                   MessageFlags = 1 << 2
	MessageFlagsSupressEmbeds                    MessageFlags = 1 << 2
	MessageFlagsSourceMessageDeleted             MessageFlags = 1 << 3
	MessageFlagsUrgent                           MessageFlags = 1 << 4
	MessageFlagsHasThread                        MessageFlags = 1 << 5
	MessageFlagsEphemeral                        MessageFlags = 1 << 6
	MessageFlagsLoading                          MessageFlags = 1 << 7
	MessageFlagsFailedToMentionSomeRolesInThread MessageFlags = 1 << 8
	MessageFlagsSuppressNotifications            MessageFlags = 1 << 12
	MessageFlagsIsVoiceMessage                   MessageFlags = 1 << 13
)

type File struct {
	Name        string
	ContentType string
	Reader      io.Reader
}

type MessageSend struct {
	Content         string                  `json:"content,omitempty"`
	Embeds          []*MessageEmbed         `json:"embeds"`
	TTS             bool                    `json:"tts"`
	Components      []MessageComponent      `json:"components"`
	Files           []*File                 `json:"-"`
	AllowedMentions *MessageAllowedMentions `json:"allowed_mentions,omitempty"`
	Reference       *MessageReference       `json:"message_reference,omitempty"`
	StickerIDs      []string                `json:"sticker_ids"`
	Flags           MessageFlags            `json:"flags,omitempty"`
	Poll            *Poll                   `json:"poll,omitempty"`
	File            *File                   `json:"-"`
	Embed           *MessageEmbed           `json:"-"`
}

type MessageEdit struct {
	Content         *string                 `json:"content,omitempty"`
	Components      *[]MessageComponent     `json:"components,omitempty"`
	Embeds          *[]*MessageEmbed        `json:"embeds,omitempty"`
	AllowedMentions *MessageAllowedMentions `json:"allowed_mentions,omitempty"`
	Flags           MessageFlags            `json:"flags,omitempty"`
	Files           []*File                 `json:"-"`
	Attachments     *[]*MessageAttachment   `json:"attachments,omitempty"`
	ID              string
	Channel         string
	Embed           *MessageEmbed `json:"-"`
}

func NewMessageEdit(channelID string, messageID string) *MessageEdit {
	return &MessageEdit{
		Channel: channelID,
		ID:      messageID,
	}
}

func (m *MessageEdit) SetContent(str string) *MessageEdit {
	m.Content = &str
	return m
}

func (m *MessageEdit) SetEmbed(embed *MessageEmbed) *MessageEdit {
	m.Embeds = &[]*MessageEmbed{embed}
	return m
}

func (m *MessageEdit) SetEmbeds(embeds []*MessageEmbed) *MessageEdit {
	m.Embeds = &embeds
	return m
}

type AllowedMentionType string

const (
	AllowedMentionTypeRoles    AllowedMentionType = "roles"
	AllowedMentionTypeUsers    AllowedMentionType = "users"
	AllowedMentionTypeEveryone AllowedMentionType = "everyone"
)

type MessageAllowedMentions struct {
	Parse       []AllowedMentionType `json:"parse"`
	Roles       []string             `json:"roles,omitempty"`
	Users       []string             `json:"users,omitempty"`
	RepliedUser bool                 `json:"replied_user"`
}

type MessageAttachment struct {
	ID          string `json:"id"`
	URL         string `json:"url"`
	ProxyURL    string `json:"proxy_url"`
	Filename    string `json:"filename"`
	ContentType string `json:"content_type"`
	Width       int    `json:"width"`
	Height      int    `json:"height"`
	Size        int    `json:"size"`
	Ephemeral   bool   `json:"ephemeral"`
}

type MessageEmbedFooter struct {
	Text         string `json:"text,omitempty"`
	IconURL      string `json:"icon_url,omitempty"`
	ProxyIconURL string `json:"proxy_icon_url,omitempty"`
}

type MessageEmbedImage struct {
	URL      string `json:"url"`
	ProxyURL string `json:"proxy_url,omitempty"`
	Width    int    `json:"width,omitempty"`
	Height   int    `json:"height,omitempty"`
}

type MessageEmbedThumbnail struct {
	URL      string `json:"url"`
	ProxyURL string `json:"proxy_url,omitempty"`
	Width    int    `json:"width,omitempty"`
	Height   int    `json:"height,omitempty"`
}

type MessageEmbedVideo struct {
	URL    string `json:"url,omitempty"`
	Width  int    `json:"width,omitempty"`
	Height int    `json:"height,omitempty"`
}

type MessageEmbedProvider struct {
	URL  string `json:"url,omitempty"`
	Name string `json:"name,omitempty"`
}

type MessageEmbedAuthor struct {
	URL          string `json:"url,omitempty"`
	Name         string `json:"name"`
	IconURL      string `json:"icon_url,omitempty"`
	ProxyIconURL string `json:"proxy_icon_url,omitempty"`
}

type MessageEmbedField struct {
	Name   string `json:"name"`
	Value  string `json:"value"`
	Inline bool   `json:"inline,omitempty"`
}

type MessageEmbed struct {
	URL         string                 `json:"url,omitempty"`
	Type        EmbedType              `json:"type,omitempty"`
	Title       string                 `json:"title,omitempty"`
	Description string                 `json:"description,omitempty"`
	Timestamp   string                 `json:"timestamp,omitempty"`
	Color       int                    `json:"color,omitempty"`
	Footer      *MessageEmbedFooter    `json:"footer,omitempty"`
	Image       *MessageEmbedImage     `json:"image,omitempty"`
	Thumbnail   *MessageEmbedThumbnail `json:"thumbnail,omitempty"`
	Video       *MessageEmbedVideo     `json:"video,omitempty"`
	Provider    *MessageEmbedProvider  `json:"provider,omitempty"`
	Author      *MessageEmbedAuthor    `json:"author,omitempty"`
	Fields      []*MessageEmbedField   `json:"fields,omitempty"`
}

type EmbedType string

const (
	EmbedTypeRich    EmbedType = "rich"
	EmbedTypeImage   EmbedType = "image"
	EmbedTypeVideo   EmbedType = "video"
	EmbedTypeGifv    EmbedType = "gifv"
	EmbedTypeArticle EmbedType = "article"
	EmbedTypeLink    EmbedType = "link"
)

type MessageReactions struct {
	Count int    `json:"count"`
	Me    bool   `json:"me"`
	Emoji *Emoji `json:"emoji"`
}

type MessageActivity struct {
	Type    MessageActivityType `json:"type"`
	PartyID string              `json:"party_id"`
}

type MessageActivityType int

const (
	MessageActivityTypeJoin        MessageActivityType = 1
	MessageActivityTypeSpectate    MessageActivityType = 2
	MessageActivityTypeListen      MessageActivityType = 3
	MessageActivityTypeJoinRequest MessageActivityType = 5
)

type MessageApplication struct {
	ID          string `json:"id"`
	CoverImage  string `json:"cover_image"`
	Description string `json:"description"`
	Icon        string `json:"icon"`
	Name        string `json:"name"`
}

type MessageSnapshot struct {
	Message *Message `json:"message"`
}

type MessageReferenceType int

const (
	MessageReferenceTypeDefault MessageReferenceType = 0
	MessageReferenceTypeForward MessageReferenceType = 1
)

type MessageReference struct {
	Type            MessageReferenceType `json:"type,omitempty"`
	MessageID       string               `json:"message_id"`
	ChannelID       string               `json:"channel_id,omitempty"`
	GuildID         string               `json:"guild_id,omitempty"`
	FailIfNotExists *bool                `json:"fail_if_not_exists,omitempty"`
}

func (m *Message) reference(refType MessageReferenceType, failIfNotExists bool) *MessageReference {
	return &MessageReference{
		Type:            refType,
		GuildID:         m.GuildID,
		ChannelID:       m.ChannelID,
		MessageID:       m.ID,
		FailIfNotExists: &failIfNotExists,
	}
}

func (m *Message) Reference() *MessageReference {
	return m.reference(MessageReferenceTypeDefault, true)
}

func (m *Message) SoftReference() *MessageReference {
	return m.reference(MessageReferenceTypeDefault, false)
}

func (m *Message) Forward() *MessageReference {
	return m.reference(MessageReferenceTypeForward, true)
}

func (m *Message) ContentWithMentionsReplaced() (content string) {
	content = m.Content

	for _, user := range m.Mentions {
		content = strings.NewReplacer(
			fmt.Sprintf("<@%s>", user.ID), fmt.Sprintf("@%s", user.Username),
			fmt.Sprintf("<@!%s>", user.ID), fmt.Sprintf("@%s", user.Username),
		).Replace(content)
	}
	return
}

var patternChannels = regexp.MustCompile("<#[^>]*>")

func (m *Message) ContentWithMoreMentionsReplaced(s *Session) (content string, err error) {
	content = m.Content

	if !s.StateEnabled {
		content = m.ContentWithMentionsReplaced()
		return
	}

	channel, err := s.State.Channel(m.ChannelID)
	if err != nil {
		content = m.ContentWithMentionsReplaced()
		return
	}

	for _, user := range m.Mentions {
		nick := user.Username

		member, err := s.State.Member(channel.GuildID, user.ID)
		if err == nil && member.Nick != "" {
			nick = member.Nick
		}

		content = strings.NewReplacer(
			fmt.Sprintf("<@%s>", user.ID), fmt.Sprintf("@%s", user.Username),
			fmt.Sprintf("<@!%s>", user.ID), fmt.Sprintf("@%s", nick),
		).Replace(content)
	}
	for _, roleID := range m.MentionRoles {
		role, err := s.State.Role(channel.GuildID, roleID)
		if err != nil || !role.Mentionable {
			continue
		}

		content = strings.Replace(content, fmt.Sprintf("<@&%s>", role.ID), fmt.Sprintf("@%s", role.Name), -1)
	}

	content = patternChannels.ReplaceAllStringFunc(content, func(mention string) string {
		channel, err := s.State.Channel(mention[2 : len(mention)-1])
		if err != nil || channel.Type == ChannelTypeGuildVoice {
			return mention
		}

		return "#" + channel.Name
	})
	return
}

type MessageInteraction struct {
	ID     string          `json:"id"`
	Type   InteractionType `json:"type"`
	Name   string          `json:"name"`
	User   *User           `json:"user"`
	Member *Member         `json:"member"`
}

type MessageInteractionMetadata struct {
	ID                            string                                `json:"id"`
	Type                          InteractionType                       `json:"type"`
	User                          *User                                 `json:"user"`
	AuthorizingIntegrationOwners  map[ApplicationIntegrationType]string `json:"authorizing_integration_owners"`
	OriginalResponseMessageID     string                                `json:"original_response_message_id,omitempty"`
	InteractedMessageID           string                                `json:"interacted_message_id,omitempty"`
	TriggeringInteractionMetadata *MessageInteractionMetadata           `json:"triggering_interaction_metadata,omitempty"`
}

const (
	LogError int = iota
	LogWarning
	LogInformational
	LogDebug
)

var Logger func(msgL, caller int, format string, a ...interface{})

func msglog(msgL, caller int, format string, a ...interface{}) {
	if Logger != nil {
		Logger(msgL, caller, format, a...)
	} else {

		pc, file, line, _ := runtime.Caller(caller)

		files := strings.Split(file, "/")
		file = files[len(files)-1]

		name := runtime.FuncForPC(pc).Name()
		fns := strings.Split(name, ".")
		name = fns[len(fns)-1]

		msg := fmt.Sprintf(format, a...)

		log.Printf("[discordgo%d] %s:%d:%s() %s\n", msgL, file, line, name, msg)
	}
}

func (s *Session) log(msgL int, format string, a ...interface{}) {
	if msgL > s.LogLevel {
		return
	}

	msglog(msgL, 2, format, a...)
}

func (v *VoiceConnection) log(msgL int, format string, a ...interface{}) {
	if msgL > v.LogLevel {
		return
	}

	msglog(msgL, 2, format, a...)
}
