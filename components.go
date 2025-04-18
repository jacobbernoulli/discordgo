package discordgo

import (
	"encoding/json"
	"fmt"
)

type ComponentType uint

const (
	ActionsRowComponent            ComponentType = 1
	ButtonComponent                ComponentType = 2
	SelectMenuComponent            ComponentType = 3
	TextInputComponent             ComponentType = 4
	UserSelectMenuComponent        ComponentType = 5
	RoleSelectMenuComponent        ComponentType = 6
	MentionableSelectMenuComponent ComponentType = 7
	ChannelSelectMenuComponent     ComponentType = 8
)

type MessageComponent interface {
	json.Marshaler
	Type() ComponentType
}

type unmarshalableMessageComponent struct {
	MessageComponent
}

func (umc *unmarshalableMessageComponent) UnmarshalJSON(src []byte) error {
	var v struct {
		Type ComponentType `json:"type"`
	}

	if err := Unmarshal(src, &v); err != nil {
		return err
	}

	var component MessageComponent
	switch v.Type {
	case ActionsRowComponent:
		component = &ActionsRow{}
	case ButtonComponent:
		component = &Button{}
	case SelectMenuComponent, ChannelSelectMenuComponent, UserSelectMenuComponent,
		RoleSelectMenuComponent, MentionableSelectMenuComponent:
		component = &SelectMenu{}
	case TextInputComponent:
		component = &TextInput{}
	default:
		return fmt.Errorf("unknown component type: %d", v.Type)
	}

	if err := Unmarshal(src, component); err != nil {
		return err
	}

	umc.MessageComponent = component
	return nil
}

func MessageComponentFromJSON(b []byte) (MessageComponent, error) {
	var u unmarshalableMessageComponent
	if err := u.UnmarshalJSON(b); err != nil {
		return nil, fmt.Errorf("failed to unmarshal into MessageComponent: %w", err)
	}
	return u.MessageComponent, nil
}

type ActionsRow struct {
	Components []MessageComponent `json:"components"`
}

func (r ActionsRow) MarshalJSON() ([]byte, error) {
	type actionsRow ActionsRow

	return Marshal(struct {
		actionsRow
		Type ComponentType `json:"type"`
	}{
		actionsRow: actionsRow(r),
		Type:       r.Type(),
	})
}

func (r *ActionsRow) UnmarshalJSON(data []byte) error {
	var v struct {
		RawComponents []unmarshalableMessageComponent `json:"components"`
	}

	if err := Unmarshal(data, &v); err != nil {
		return err
	}

	r.Components = make([]MessageComponent, len(v.RawComponents))
	for i, component := range v.RawComponents {
		r.Components[i] = component.MessageComponent
	}

	return nil
}

func (r ActionsRow) Type() ComponentType {
	return ActionsRowComponent
}

type ButtonStyle uint

const (
	PrimaryButton   ButtonStyle = 1
	SecondaryButton ButtonStyle = 2
	SuccessButton   ButtonStyle = 3
	DangerButton    ButtonStyle = 4
	LinkButton      ButtonStyle = 5
	PremiumButton   ButtonStyle = 6
)

type ComponentEmoji struct {
	Name     string `json:"name,omitempty"`
	ID       string `json:"id,omitempty"`
	Animated bool   `json:"animated,omitempty"`
}

type Button struct {
	Label    string          `json:"label"`
	Style    ButtonStyle     `json:"style"`
	Disabled bool            `json:"disabled"`
	Emoji    *ComponentEmoji `json:"emoji,omitempty"`
	URL      string          `json:"url,omitempty"`
	CustomID string          `json:"custom_id,omitempty"`
	SKUID    string          `json:"sku_id,omitempty"`
}

func (b Button) MarshalJSON() ([]byte, error) {
	type button Button

	if b.Style == 0 {
		b.Style = PrimaryButton
	}

	return Marshal(struct {
		button
		Type ComponentType `json:"type"`
	}{
		button: button(b),
		Type:   b.Type(),
	})
}

func (Button) Type() ComponentType {
	return ButtonComponent
}

type SelectMenuOption struct {
	Label       string          `json:"label,omitempty"`
	Value       string          `json:"value"`
	Description string          `json:"description"`
	Emoji       *ComponentEmoji `json:"emoji,omitempty"`
	Default     bool            `json:"default"`
}

type SelectMenuDefaultValueType string

const (
	SelectMenuDefaultValueUser    SelectMenuDefaultValueType = "user"
	SelectMenuDefaultValueRole    SelectMenuDefaultValueType = "role"
	SelectMenuDefaultValueChannel SelectMenuDefaultValueType = "channel"
)

type SelectMenuDefaultValue struct {
	ID   string                     `json:"id"`
	Type SelectMenuDefaultValueType `json:"type"`
}

type SelectMenuType ComponentType

const (
	StringSelectMenu      = SelectMenuType(SelectMenuComponent)
	UserSelectMenu        = SelectMenuType(UserSelectMenuComponent)
	RoleSelectMenu        = SelectMenuType(RoleSelectMenuComponent)
	MentionableSelectMenu = SelectMenuType(MentionableSelectMenuComponent)
	ChannelSelectMenu     = SelectMenuType(ChannelSelectMenuComponent)
)

type SelectMenu struct {
	MenuType      SelectMenuType           `json:"type,omitempty"`
	CustomID      string                   `json:"custom_id,omitempty"`
	Placeholder   string                   `json:"placeholder"`
	MinValues     *int                     `json:"min_values,omitempty"`
	MaxValues     int                      `json:"max_values,omitempty"`
	DefaultValues []SelectMenuDefaultValue `json:"default_values,omitempty"`
	Options       []SelectMenuOption       `json:"options,omitempty"`
	Disabled      bool                     `json:"disabled"`
	ChannelTypes  []ChannelType            `json:"channel_types,omitempty"`
}

func (s SelectMenu) Type() ComponentType {
	if s.MenuType != 0 {
		return ComponentType(s.MenuType)
	}
	return SelectMenuComponent
}

func (s SelectMenu) MarshalJSON() ([]byte, error) {
	type selectMenu SelectMenu

	return Marshal(struct {
		selectMenu
		Type ComponentType `json:"type"`
	}{
		selectMenu: selectMenu(s),
		Type:       s.Type(),
	})
}

type TextInput struct {
	CustomID    string         `json:"custom_id"`
	Label       string         `json:"label"`
	Style       TextInputStyle `json:"style"`
	Placeholder string         `json:"placeholder,omitempty"`
	Value       string         `json:"value,omitempty"`
	Required    bool           `json:"required"`
	MinLength   int            `json:"min_length,omitempty"`
	MaxLength   int            `json:"max_length,omitempty"`
}

func (TextInput) Type() ComponentType {
	return TextInputComponent
}

func (m TextInput) MarshalJSON() ([]byte, error) {
	type inputText TextInput

	return Marshal(struct {
		inputText
		Type ComponentType `json:"type"`
	}{
		inputText: inputText(m),
		Type:      m.Type(),
	})
}

type TextInputStyle uint

const (
	TextInputShort     TextInputStyle = 1
	TextInputParagraph TextInputStyle = 2
)
