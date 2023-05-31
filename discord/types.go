package discord

const weeaaImage = "https://images-ext-2.discordapp.net/external/VeKohQ7LpG4m2Xubc8m1Um7-32RVIsEJSBqB9OAA9kA/%3Fsize%3D4096/https/cdn.discordapp.com/avatars/363975551393988620/3153e87ecb026728055a7e61433c2594.png?width=1132&height=1132"
const weeaaFooterText = "@weeaa – 🦦"

type ExchangeArtWebhook struct {
	Name        string
	Description string
	Image       string
	MintLink    string
	CMID        string
	Supply      string
	ReleaseType string
	Minted      int
	MintCap     int
	Artist      string
	Edition     interface{}
	EditionBool interface{}
	Price       string
	LiveAt      string
	ToSend      bool
}

/*🍀 DISCORD TYPES 🍀*/
type Webhook struct {
	Content   string  `json:"content,omitempty"`
	Username  string  `json:"username,omitempty"`
	AvatarUrl string  `json:"avatar_url,omitempty"`
	Tts       bool    `json:"tts,omitempty"`
	Embeds    []Embed `json:"embeds,omitempty"`
}

type Embed struct {
	Title       string         `json:"title,omitempty"`
	Description string         `json:"description,omitempty"`
	Url         string         `json:"url,omitempty"`
	Timestamp   string         `json:"timestamp,omitempty"`
	Color       int            `json:"color,omitempty"`
	Footer      EmbedFooter    `json:"footer,omitempty"`
	Image       EmbedImage     `json:"image,omitempty"`
	Thumbnail   EmbedThumbnail `json:"thumbnail,omitempty"`
	Video       EmbedVideo     `json:"video,omitempty"`
	Provider    EmbedProvider  `json:"provider,omitempty"`
	Author      EmbedAuthor    `json:"author,omitempty"`
	Fields      []EmbedFields  `json:"fields,omitempty"`
}

type EmbedFooter struct {
	Text         string `json:"text,omitempty"`
	IconUrl      string `json:"icon_url,omitempty"`
	ProxyIconUrl string `json:"proxy_icon_url,omitempty"`
}

type EmbedImage struct {
	Url      string `json:"url,omitempty"`
	ProxyUrl string `json:"proxy_url,omitempty"`
	Height   int    `json:"height,omitempty"`
	Width    int    `json:"width,omitempty"`
}

type EmbedThumbnail struct {
	Url      string `json:"url,omitempty"`
	ProxyUrl string `json:"proxy_url,omitempty"`
	Height   int    `json:"height,omitempty"`
	Width    int    `json:"width,omitempty"`
}

type EmbedVideo struct {
	Url      string `json:"url,omitempty"`
	ProxyUrl string `json:"proxy_url,omitempty"`
	Height   int    `json:"height,omitempty"`
	Width    int    `json:"width,omitempty"`
}

type EmbedProvider struct {
	Name string `json:"name,omitempty"`
	Url  string `json:"url,omitempty"`
}

type EmbedAuthor struct {
	Name         string `json:"name,omitempty"`
	Url          string `json:"url,omitempty"`
	IconUrl      string `json:"icon_url,omitempty"`
	ProxyIconUrl string `json:"proxy_icon_url,omitempty"`
}

type EmbedFields struct {
	Name   string `json:"name,omitempty"`
	Value  string `json:"value,omitempty"`
	Inline bool   `json:"inline,omitempty"`
}
