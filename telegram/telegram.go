package telegram

import (
	"context"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/nmeilick/go-notify/message"
	"github.com/nmeilick/go-notify/resolve"
	"github.com/nmeilick/go-tools"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

const (
	MaxCaptionSize = 1024
	MaxMessageSize = 4096
)

type Telegram struct{}

func (t *Telegram) Notify(ctx context.Context, recipient string, m *message.Message) error {
	var log zerolog.Logger
	if l, ok := ctx.Value("log").(zerolog.Logger); ok {
		log = l
	}

	uri, err := url.Parse(recipient)
	if err != nil {
		return nil
	}

	switch strings.ToLower(uri.Scheme) {
	case "tg", "telegram":
	default:
		return nil
	}

	query := uri.Query()
	token, err := resolve.Value(query.Get("bot"))
	if err != nil {
		return errors.Wrap(err, "failed to resolve bot token")
	} else if token == "" {
		return errors.New("bot token is empty")
	}

	ids, err := resolve.Value(strings.TrimPrefix(uri.Path, "/"))
	if err != nil {
		return errors.Wrap(err, "failed to resolve channels")
	}

	var channels []int64
	for _, c := range tools.Tokens(ids) {
		if n, err := strconv.ParseInt(c, 10, 64); err != nil || n == 0 {
			return errors.New("not a numeric channel id: " + c)
		} else {
			channels = append(channels, n)
		}
	}
	channels = tools.Unique(channels)
	if len(channels) == 0 {
		return errors.New("no channels defined")
	}

	bot, err := tg.NewBotAPI(token)
	if err != nil {
		return err
	}
	bot.Debug = tools.IsOn(os.Getenv("DEBUG"), true)

	for _, id := range channels {
		if err := ctx.Err(); err != nil {
			return err
		}

		text := m.Text
		summary := m.Summary
		file := m.File

		var lastID int
		if file != "" {
			var doc tg.Chattable
			var caption string
			if text != "" && len(text) <= MaxCaptionSize {
				caption = text
				text, summary = "", ""
			} else if summary != "" {
				caption = summary
				summary = ""
			}

			switch strings.ToLower(filepath.Ext(file)) {
			case ".mp3", ".wav", ".flac", ".aac", ".ogg", ".m4a":
				audio := tg.NewAudio(id, tg.FilePath(file))
				audio.Caption = caption
				audio.ReplyToMessageID = lastID
				doc = audio
			case ".mp4", ".avi", ".mov", ".flv", ".wmv", ".mkv":
				video := tg.NewVideo(id, tg.FilePath(file))
				video.ReplyToMessageID = lastID
				video.Caption = caption
				doc = video
			case ".jpg", ".jpeg", ".png", ".gif", ".bmp", ".tiff":
				photo := tg.NewPhoto(id, tg.FilePath(file))
				photo.ReplyToMessageID = lastID
				photo.Caption = caption
				doc = photo
			default:
				unknown := tg.NewDocument(id, tg.FilePath(file))
				unknown.ReplyToMessageID = lastID
				unknown.Caption = caption
				doc = unknown
			}

			log.Info().Int64("channel", id).Msg("Sending media file")
			msg, err := bot.Send(doc)
			if err != nil {
				log.Err(err).Int64("channel", id).Msg("Error sending media file")
				continue
			}
			lastID = msg.MessageID
		}

		var part string
		for text != "" {
			if len(text) > MaxMessageSize {
				part = text[:MaxMessageSize]
			} else {
				part = text
			}
			text = text[len(part):]
			m := tg.NewMessage(id, part)
			m.ReplyToMessageID = lastID
			msg, err := bot.Send(m)
			if err != nil {
				break
			}
			lastID = msg.MessageID
		}
	}

	return nil
}
