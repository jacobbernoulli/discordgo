package discordgo

import (
	"bytes"
	"crypto/ed25519"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

const InteractionDeadline = time.Second * 3

type ApplicationCommandType uint8

const (
	ChatApplicationCommand    ApplicationCommandType = 1
	UserApplicationCommand    ApplicationCommandType = 2
	MessageApplicationCommand ApplicationCommandType = 3
)

type ApplicationCommand struct {
	ID                       string                        `json:"id,omitempty"`
	ApplicationID            string                        `json:"application_id,omitempty"`
	GuildID                  string                        `json:"guild_id,omitempty"`
	Version                  string                        `json:"version,omitempty"`
	Type                     ApplicationCommandType        `json:"type,omitempty"`
	Name                     string                        `json:"name"`
	NameLocalizations        *map[Locale]string            `json:"name_localizations,omitempty"`
	DefaultPermission        *bool                         `json:"default_permission,omitempty"`
	DefaultMemberPermissions *int64                        `json:"default_member_permissions,string,omitempty"`
	NSFW                     *bool                         `json:"nsfw,omitempty"`
	DMPermission             *bool                         `json:"dm_permission,omitempty"`
	Contexts                 *[]InteractionContextType     `json:"contexts,omitempty"`
	IntegrationTypes         *[]ApplicationIntegrationType `json:"integration_types,omitempty"`
	Description              string                        `json:"description,omitempty"`
	DescriptionLocalizations *map[Locale]string            `json:"description_localizations,omitempty"`
	Options                  []*ApplicationCommandOption   `json:"options"`
}

type ApplicationCommandOptionType uint8

const (
	ApplicationCommandOptionSubCommand      ApplicationCommandOptionType = 1
	ApplicationCommandOptionSubCommandGroup ApplicationCommandOptionType = 2
	ApplicationCommandOptionString          ApplicationCommandOptionType = 3
	ApplicationCommandOptionInteger         ApplicationCommandOptionType = 4
	ApplicationCommandOptionBoolean         ApplicationCommandOptionType = 5
	ApplicationCommandOptionUser            ApplicationCommandOptionType = 6
	ApplicationCommandOptionChannel         ApplicationCommandOptionType = 7
	ApplicationCommandOptionRole            ApplicationCommandOptionType = 8
	ApplicationCommandOptionMentionable     ApplicationCommandOptionType = 9
	ApplicationCommandOptionNumber          ApplicationCommandOptionType = 10
	ApplicationCommandOptionAttachment      ApplicationCommandOptionType = 11
)

func (t ApplicationCommandOptionType) String() string {
	switch t {
	case ApplicationCommandOptionSubCommand:
		return "SubCommand"
	case ApplicationCommandOptionSubCommandGroup:
		return "SubCommandGroup"
	case ApplicationCommandOptionString:
		return "String"
	case ApplicationCommandOptionInteger:
		return "Integer"
	case ApplicationCommandOptionBoolean:
		return "Boolean"
	case ApplicationCommandOptionUser:
		return "User"
	case ApplicationCommandOptionChannel:
		return "Channel"
	case ApplicationCommandOptionRole:
		return "Role"
	case ApplicationCommandOptionMentionable:
		return "Mentionable"
	case ApplicationCommandOptionNumber:
		return "Number"
	case ApplicationCommandOptionAttachment:
		return "Attachment"
	}
	return fmt.Sprintf("ApplicationCommandOptionType(%d)", t)
}

type ApplicationCommandOption struct {
	Type                     ApplicationCommandOptionType      `json:"type"`
	Name                     string                            `json:"name"`
	NameLocalizations        map[Locale]string                 `json:"name_localizations,omitempty"`
	Description              string                            `json:"description,omitempty"`
	DescriptionLocalizations map[Locale]string                 `json:"description_localizations,omitempty"`
	ChannelTypes             []ChannelType                     `json:"channel_types"`
	Required                 bool                              `json:"required"`
	Options                  []*ApplicationCommandOption       `json:"options"`
	Autocomplete             bool                              `json:"autocomplete"`
	Choices                  []*ApplicationCommandOptionChoice `json:"choices"`
	MinValue                 *float64                          `json:"min_value,omitempty"`
	MaxValue                 float64                           `json:"max_value,omitempty"`
	MinLength                *int                              `json:"min_length,omitempty"`
	MaxLength                int                               `json:"max_length,omitempty"`
}

type ApplicationCommandOptionChoice struct {
	Name              string            `json:"name"`
	NameLocalizations map[Locale]string `json:"name_localizations,omitempty"`
	Value             interface{}       `json:"value"`
}

type ApplicationCommandPermissions struct {
	ID         string                           `json:"id"`
	Type       ApplicationCommandPermissionType `json:"type"`
	Permission bool                             `json:"permission"`
}

func GuildAllChannelsID(guild string) (id string, err error) {
	var v uint64
	v, err = strconv.ParseUint(guild, 10, 64)
	if err != nil {
		return
	}

	return strconv.FormatUint(v-1, 10), nil
}

type ApplicationCommandPermissionsList struct {
	Permissions []*ApplicationCommandPermissions `json:"permissions"`
}

type GuildApplicationCommandPermissions struct {
	ID            string                           `json:"id"`
	ApplicationID string                           `json:"application_id"`
	GuildID       string                           `json:"guild_id"`
	Permissions   []*ApplicationCommandPermissions `json:"permissions"`
}

type ApplicationCommandPermissionType uint8

const (
	ApplicationCommandPermissionTypeRole    ApplicationCommandPermissionType = 1
	ApplicationCommandPermissionTypeUser    ApplicationCommandPermissionType = 2
	ApplicationCommandPermissionTypeChannel ApplicationCommandPermissionType = 3
)

type InteractionType uint8

const (
	InteractionPing                           InteractionType = 1
	InteractionApplicationCommand             InteractionType = 2
	InteractionMessageComponent               InteractionType = 3
	InteractionApplicationCommandAutocomplete InteractionType = 4
	InteractionModalSubmit                    InteractionType = 5
)

func (t InteractionType) String() string {
	switch t {
	case InteractionPing:
		return "Ping"
	case InteractionApplicationCommand:
		return "ApplicationCommand"
	case InteractionMessageComponent:
		return "MessageComponent"
	case InteractionModalSubmit:
		return "ModalSubmit"
	}
	return fmt.Sprintf("InteractionType(%d)", t)
}

type InteractionContextType uint

const (
	InteractionContextGuild          InteractionContextType = 0
	InteractionContextBotDM          InteractionContextType = 1
	InteractionContextPrivateChannel InteractionContextType = 2
)

type Interaction struct {
	ID                           string                                `json:"id"`
	AppID                        string                                `json:"application_id"`
	Type                         InteractionType                       `json:"type"`
	Data                         InteractionData                       `json:"data"`
	GuildID                      string                                `json:"guild_id"`
	ChannelID                    string                                `json:"channel_id"`
	Message                      *Message                              `json:"message"`
	AppPermissions               int64                                 `json:"app_permissions,string"`
	Member                       *Member                               `json:"member"`
	User                         *User                                 `json:"user"`
	Locale                       Locale                                `json:"locale"`
	GuildLocale                  *Locale                               `json:"guild_locale"`
	Context                      InteractionContextType                `json:"context"`
	AuthorizingIntegrationOwners map[ApplicationIntegrationType]string `json:"authorizing_integration_owners"`
	Token                        string                                `json:"token"`
	Version                      int                                   `json:"version"`
	Entitlements                 []*Entitlement                        `json:"entitlements"`
}

type interaction Interaction

type rawInteraction struct {
	interaction
	Data json.RawMessage `json:"data"`
}

func (i *Interaction) UnmarshalJSON(raw []byte) error {
	var tmp rawInteraction
	if err := json.Unmarshal(raw, &tmp); err != nil {
		return err
	}

	*i = Interaction(tmp.interaction)

	var parsed InteractionData

	switch tmp.Type {
	case InteractionApplicationCommand, InteractionApplicationCommandAutocomplete:
		var v ApplicationCommandInteractionData
		if err := json.Unmarshal(tmp.Data, &v); err != nil {
			return err
		}
		parsed = v
	case InteractionMessageComponent:
		var v MessageComponentInteractionData
		if err := json.Unmarshal(tmp.Data, &v); err != nil {
			return err
		}
		parsed = v
	case InteractionModalSubmit:
		var v ModalSubmitInteractionData
		if err := json.Unmarshal(tmp.Data, &v); err != nil {
			return err
		}
		parsed = v
	default:
		return nil
	}

	i.Data = parsed
	return nil
}

func (i Interaction) MessageComponentData() (data MessageComponentInteractionData) {
	if i.Type != InteractionMessageComponent {
		panic("MessageComponentData called on interaction of type " + i.Type.String())
	}
	return i.Data.(MessageComponentInteractionData)
}

func (i Interaction) ApplicationCommandData() (data ApplicationCommandInteractionData) {
	if i.Type != InteractionApplicationCommand && i.Type != InteractionApplicationCommandAutocomplete {
		panic("ApplicationCommandData called on interaction of type " + i.Type.String())
	}
	return i.Data.(ApplicationCommandInteractionData)
}

func (i Interaction) ModalSubmitData() (data ModalSubmitInteractionData) {
	if i.Type != InteractionModalSubmit {
		panic("ModalSubmitData called on interaction of type " + i.Type.String())
	}
	return i.Data.(ModalSubmitInteractionData)
}

type InteractionData interface {
	Type() InteractionType
}

type ApplicationCommandInteractionData struct {
	ID          string                                     `json:"id"`
	Name        string                                     `json:"name"`
	CommandType ApplicationCommandType                     `json:"type"`
	Resolved    *ApplicationCommandInteractionDataResolved `json:"resolved"`
	Options     []*ApplicationCommandInteractionDataOption `json:"options"`
	TargetID    string                                     `json:"target_id"`
}

func (d ApplicationCommandInteractionData) GetOption(name string) (option *ApplicationCommandInteractionDataOption) {
	for _, opt := range d.Options {
		if opt.Name == name {
			option = opt
			break
		}
	}

	return
}

type ApplicationCommandInteractionDataResolved struct {
	Users       map[string]*User              `json:"users"`
	Members     map[string]*Member            `json:"members"`
	Roles       map[string]*Role              `json:"roles"`
	Channels    map[string]*Channel           `json:"channels"`
	Messages    map[string]*Message           `json:"messages"`
	Attachments map[string]*MessageAttachment `json:"attachments"`
}

func (ApplicationCommandInteractionData) Type() InteractionType {
	return InteractionApplicationCommand
}

type MessageComponentInteractionData struct {
	CustomID      string                                  `json:"custom_id"`
	ComponentType ComponentType                           `json:"component_type"`
	Resolved      MessageComponentInteractionDataResolved `json:"resolved"`
	Values        []string                                `json:"values"`
}

type MessageComponentInteractionDataResolved struct {
	Users    map[string]*User    `json:"users"`
	Members  map[string]*Member  `json:"members"`
	Roles    map[string]*Role    `json:"roles"`
	Channels map[string]*Channel `json:"channels"`
}

func (MessageComponentInteractionData) Type() InteractionType {
	return InteractionMessageComponent
}

type ModalSubmitInteractionData struct {
	CustomID   string             `json:"custom_id"`
	Components []MessageComponent `json:"-"`
}

func (ModalSubmitInteractionData) Type() InteractionType {
	return InteractionModalSubmit
}

func (d *ModalSubmitInteractionData) UnmarshalJSON(data []byte) error {
	type modalSubmitInteractionData ModalSubmitInteractionData
	var v struct {
		modalSubmitInteractionData
		RawComponents []unmarshalableMessageComponent `json:"components"`
	}
	err := json.Unmarshal(data, &v)
	if err != nil {
		return err
	}
	*d = ModalSubmitInteractionData(v.modalSubmitInteractionData)
	d.Components = make([]MessageComponent, len(v.RawComponents))
	for i, v := range v.RawComponents {
		d.Components[i] = v.MessageComponent
	}
	return err
}

type ApplicationCommandInteractionDataOption struct {
	Name    string                                     `json:"name"`
	Type    ApplicationCommandOptionType               `json:"type"`
	Value   interface{}                                `json:"value,omitempty"`
	Options []*ApplicationCommandInteractionDataOption `json:"options,omitempty"`
	Focused bool                                       `json:"focused,omitempty"`
}

func (o ApplicationCommandInteractionDataOption) GetOption(name string) (option *ApplicationCommandInteractionDataOption) {
	for _, opt := range o.Options {
		if opt.Name == name {
			option = opt
			break
		}
	}

	return
}

func (o ApplicationCommandInteractionDataOption) IntValue() int64 {
	if o.Type != ApplicationCommandOptionInteger {
		panic("IntValue called on data option of type " + o.Type.String())
	}
	return int64(o.Value.(float64))
}

func (o ApplicationCommandInteractionDataOption) UintValue() uint64 {
	if o.Type != ApplicationCommandOptionInteger {
		panic("UintValue called on data option of type " + o.Type.String())
	}
	return uint64(o.Value.(float64))
}

func (o ApplicationCommandInteractionDataOption) FloatValue() float64 {
	if o.Type != ApplicationCommandOptionNumber {
		panic("FloatValue called on data option of type " + o.Type.String())
	}
	return o.Value.(float64)
}

func (o ApplicationCommandInteractionDataOption) StringValue() string {
	if o.Type != ApplicationCommandOptionString {
		panic("StringValue called on data option of type " + o.Type.String())
	}
	return o.Value.(string)
}

func (o ApplicationCommandInteractionDataOption) BoolValue() bool {
	if o.Type != ApplicationCommandOptionBoolean {
		panic("BoolValue called on data option of type " + o.Type.String())
	}
	return o.Value.(bool)
}

func (o ApplicationCommandInteractionDataOption) ChannelValue(s *Session) *Channel {
	if o.Type != ApplicationCommandOptionChannel {
		panic("ChannelValue called on data option of type " + o.Type.String())
	}
	chanID := o.Value.(string)

	if s == nil {
		return &Channel{ID: chanID}
	}

	ch, err := s.State.Channel(chanID)
	if err != nil {
		ch, err = s.Channel(chanID)
		if err != nil {
			return &Channel{ID: chanID}
		}
	}

	return ch
}

func (o ApplicationCommandInteractionDataOption) RoleValue(s *Session, gID string) *Role {
	if o.Type != ApplicationCommandOptionRole && o.Type != ApplicationCommandOptionMentionable {
		panic("RoleValue called on data option of type " + o.Type.String())
	}
	roleID := o.Value.(string)

	if s == nil || gID == "" {
		return &Role{ID: roleID}
	}

	r, err := s.State.Role(gID, roleID)
	if err != nil {
		roles, err := s.GuildRoles(gID)
		if err == nil {
			for _, r = range roles {
				if r.ID == roleID {
					return r
				}
			}
		}
		return &Role{ID: roleID}
	}

	return r
}

func (o ApplicationCommandInteractionDataOption) UserValue(s *Session) *User {
	if o.Type != ApplicationCommandOptionUser && o.Type != ApplicationCommandOptionMentionable {
		panic("UserValue called on data option of type " + o.Type.String())
	}
	userID := o.Value.(string)

	if s == nil {
		return &User{ID: userID}
	}

	u, err := s.User(userID)
	if err != nil {
		return &User{ID: userID}
	}

	return u
}

type InteractionResponseType uint8

const (
	InteractionResponsePong                             InteractionResponseType = 1
	InteractionResponseChannelMessageWithSource         InteractionResponseType = 4
	InteractionResponseDeferredChannelMessageWithSource InteractionResponseType = 5
	InteractionResponseDeferredMessageUpdate            InteractionResponseType = 6
	InteractionResponseUpdateMessage                    InteractionResponseType = 7
	InteractionApplicationCommandAutocompleteResult     InteractionResponseType = 8
	InteractionResponseModal                            InteractionResponseType = 9
)

type InteractionResponse struct {
	Type InteractionResponseType  `json:"type,omitempty"`
	Data *InteractionResponseData `json:"data,omitempty"`
}

type InteractionResponseData struct {
	TTS             bool                              `json:"tts"`
	Content         string                            `json:"content"`
	Components      []MessageComponent                `json:"components"`
	Embeds          []*MessageEmbed                   `json:"embeds"`
	AllowedMentions *MessageAllowedMentions           `json:"allowed_mentions,omitempty"`
	Files           []*File                           `json:"-"`
	Attachments     *[]*MessageAttachment             `json:"attachments,omitempty"`
	Poll            *Poll                             `json:"poll,omitempty"`
	Flags           MessageFlags                      `json:"flags,omitempty"`
	Choices         []*ApplicationCommandOptionChoice `json:"choices,omitempty"`
	CustomID        string                            `json:"custom_id,omitempty"`
	Title           string                            `json:"title,omitempty"`
}

func VerifyInteraction(r *http.Request, key ed25519.PublicKey) bool {
	var msg bytes.Buffer

	signature := r.Header.Get("X-Signature-Ed25519")
	if signature == "" {
		return false
	}

	sig, err := hex.DecodeString(signature)
	if err != nil {
		return false
	}

	if len(sig) != ed25519.SignatureSize {
		return false
	}

	timestamp := r.Header.Get("X-Signature-Timestamp")
	if timestamp == "" {
		return false
	}

	msg.WriteString(timestamp)

	defer r.Body.Close()
	var body bytes.Buffer

	defer func() {
		r.Body = ioutil.NopCloser(&body)
	}()

	_, err = io.Copy(&msg, io.TeeReader(r.Body, &body))
	if err != nil {
		return false
	}

	return ed25519.Verify(key, msg.Bytes(), sig)
}
