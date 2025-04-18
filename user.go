package discordgo

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/textproto"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type UserFlags int

const (
	UserFlagDiscordEmployee           UserFlags = 1 << 0
	UserFlagDiscordPartner            UserFlags = 1 << 1
	UserFlagHypeSquadEvents           UserFlags = 1 << 2
	UserFlagBugHunterLevel1           UserFlags = 1 << 3
	UserFlagHouseBravery              UserFlags = 1 << 6
	UserFlagHouseBrilliance           UserFlags = 1 << 7
	UserFlagHouseBalance              UserFlags = 1 << 8
	UserFlagEarlySupporter            UserFlags = 1 << 9
	UserFlagTeamUser                  UserFlags = 1 << 10
	UserFlagSystem                    UserFlags = 1 << 12
	UserFlagBugHunterLevel2           UserFlags = 1 << 14
	UserFlagVerifiedBot               UserFlags = 1 << 16
	UserFlagVerifiedBotDeveloper      UserFlags = 1 << 17
	UserFlagDiscordCertifiedModerator UserFlags = 1 << 18
	UserFlagBotHTTPInteractions       UserFlags = 1 << 19
	UserFlagSpammer                   UserFlags = 1 << 20
	UserFlagActiveBotDeveloper        UserFlags = 1 << 22
)

type UserPremiumType int

const (
	UserPremiumTypeNone         UserPremiumType = 0
	UserPremiumTypeNitroClassic UserPremiumType = 1
	UserPremiumTypeNitro        UserPremiumType = 2
	UserPremiumTypeNitroBasic   UserPremiumType = 3
)

type User struct {
	ID            string          `json:"id"`
	Email         string          `json:"email"`
	Username      string          `json:"username"`
	Avatar        string          `json:"avatar"`
	Locale        string          `json:"locale"`
	Discriminator string          `json:"discriminator"`
	GlobalName    string          `json:"global_name"`
	Token         string          `json:"token"`
	Verified      bool            `json:"verified"`
	MFAEnabled    bool            `json:"mfa_enabled"`
	Banner        string          `json:"banner"`
	AccentColor   int             `json:"accent_color"`
	Bot           bool            `json:"bot"`
	PublicFlags   UserFlags       `json:"public_flags"`
	PremiumType   UserPremiumType `json:"premium_type"`
	System        bool            `json:"system"`
	Flags         int             `json:"flags"`
}

func (u *User) String() string {
	if u.Discriminator == "0" {
		return u.Username
	}

	return u.Username + "#" + u.Discriminator
}

func (u *User) Mention() string {
	return "<@" + u.ID + ">"
}

func (u *User) AvatarURL(size string) string {
	return avatarURL(
		u.Avatar,
		EndpointDefaultUserAvatar(u.DefaultAvatarIndex()),
		EndpointUserAvatar(u.ID, u.Avatar),
		EndpointUserAvatarAnimated(u.ID, u.Avatar),
		size,
	)
}

func (u *User) BannerURL(size string) string {
	return bannerURL(u.Banner, EndpointUserBanner(u.ID, u.Banner), EndpointUserBannerAnimated(u.ID, u.Banner), size)
}

func (u *User) Int64ID() (uint64, error) {
	return strconv.ParseUint(u.ID, 10, 64)
}

func (u *User) DefaultAvatarIndex() int {
	if u.Discriminator == "0" {
		id, _ := strconv.ParseUint(u.ID, 10, 64)
		return int((id >> 22) % 6)
	}

	id, _ := strconv.Atoi(u.Discriminator)
	return id % 5
}

func (u *User) DisplayName() string {
	if u.GlobalName != "" {
		return u.GlobalName
	}
	return u.Username
}

type Locale string

func (l Locale) String() string {
	if name, ok := Locales[l]; ok {
		return name
	}
	return Unknown.String()
}

const (
	EnglishUS    Locale = "en-US"
	EnglishGB    Locale = "en-GB"
	Bulgarian    Locale = "bg"
	ChineseCN    Locale = "zh-CN"
	ChineseTW    Locale = "zh-TW"
	Croatian     Locale = "hr"
	Czech        Locale = "cs"
	Danish       Locale = "da"
	Dutch        Locale = "nl"
	Finnish      Locale = "fi"
	French       Locale = "fr"
	German       Locale = "de"
	Greek        Locale = "el"
	Hindi        Locale = "hi"
	Hungarian    Locale = "hu"
	Italian      Locale = "it"
	Japanese     Locale = "ja"
	Korean       Locale = "ko"
	Lithuanian   Locale = "lt"
	Norwegian    Locale = "no"
	Polish       Locale = "pl"
	PortugueseBR Locale = "pt-BR"
	Romanian     Locale = "ro"
	Russian      Locale = "ru"
	SpanishES    Locale = "es-ES"
	SpanishLATAM Locale = "es-419"
	Swedish      Locale = "sv-SE"
	Thai         Locale = "th"
	Turkish      Locale = "tr"
	Ukrainian    Locale = "uk"
	Vietnamese   Locale = "vi"
	Unknown      Locale = ""
)

var Locales = map[Locale]string{
	EnglishUS:    "English (United States)",
	EnglishGB:    "English (Great Britain)",
	Bulgarian:    "Bulgarian",
	ChineseCN:    "Chinese (China)",
	ChineseTW:    "Chinese (Taiwan)",
	Croatian:     "Croatian",
	Czech:        "Czech",
	Danish:       "Danish",
	Dutch:        "Dutch",
	Finnish:      "Finnish",
	French:       "French",
	German:       "German",
	Greek:        "Greek",
	Hindi:        "Hindi",
	Hungarian:    "Hungarian",
	Italian:      "Italian",
	Japanese:     "Japanese",
	Korean:       "Korean",
	Lithuanian:   "Lithuanian",
	Norwegian:    "Norwegian",
	Polish:       "Polish",
	PortugueseBR: "Portuguese (Brazil)",
	Romanian:     "Romanian",
	Russian:      "Russian",
	SpanishES:    "Spanish (Spain)",
	SpanishLATAM: "Spanish (LATAM)",
	Swedish:      "Swedish",
	Thai:         "Thai",
	Turkish:      "Turkish",
	Ukrainian:    "Ukrainian",
	Vietnamese:   "Vietnamese",
	Unknown:      "unknown",
}

func SnowflakeTimestamp(ID string) (time.Time, error) {
	const epoch = 1420070400000

	i, err := strconv.ParseInt(ID, 10, 64)
	if err != nil {
		return time.Time{}, err
	}

	ms := (i >> 22) + epoch
	return time.UnixMilli(ms), nil
}

func MultipartBodyWithJSON(data interface{}, files []*File) (string, []byte, error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	defer writer.Close()

	payload, err := Marshal(data)
	if err != nil {
		return "", nil, err
	}

	jsonHeader := make(textproto.MIMEHeader)
	jsonHeader.Set("Content-Disposition", `form-data; name="payload_json"`)
	jsonHeader.Set("Content-Type", "application/json")

	jsonPart, err := writer.CreatePart(jsonHeader)
	if err != nil {
		return "", nil, err
	}
	if _, err = jsonPart.Write(payload); err != nil {
		return "", nil, err
	}

	for i, file := range files {
		fileHeader := make(textproto.MIMEHeader)
		fileHeader.Set("Content-Disposition", fmt.Sprintf(`form-data; name="files[%d]"; filename="%s"`, i, quoteEscaper.Replace(file.Name)))

		contentType := file.ContentType
		if contentType == "" {
			contentType = "application/octet-stream"
		}
		fileHeader.Set("Content-Type", contentType)

		filePart, err := writer.CreatePart(fileHeader)
		if err != nil {
			return "", nil, err
		}
		if _, err = io.Copy(filePart, file.Reader); err != nil {
			return "", nil, err
		}
	}

	return writer.FormDataContentType(), body.Bytes(), nil
}

func avatarURL(avatarHash, defaultAvatarURL, staticAvatarURL, animatedAvatarURL, size string) string {
	if avatarHash == "" {
		return defaultAvatarURL
	}

	var URL string
	if strings.HasPrefix(avatarHash, "a_") {
		URL = animatedAvatarURL
	} else {
		URL = staticAvatarURL
	}

	if size != "" {
		URL = fmt.Sprintf("%s?size=%s", URL, url.QueryEscape(size))
	}

	return URL
}

func bannerURL(bannerHash, staticBannerURL, animatedBannerURL, size string) string {
	if bannerHash == "" {
		return ""
	}

	var URL string
	if strings.HasPrefix(bannerHash, "a_") {
		URL = animatedBannerURL
	} else {
		URL = staticBannerURL
	}

	if size != "" {
		return fmt.Sprintf("%s?size=%s", URL, url.QueryEscape(size))
	}

	return URL
}

func iconURL(iconHash, staticIconURL, animatedIconURL, size string) string {
	if iconHash == "" {
		return ""
	}

	var URL string
	if strings.HasPrefix(iconHash, "a_") {
		URL = animatedIconURL
	} else {
		URL = staticIconURL
	}

	if size != "" {
		return fmt.Sprintf("%s?size=%s", URL, url.QueryEscape(size))
	}

	return URL
}
