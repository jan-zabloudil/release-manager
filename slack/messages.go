package slack

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	"github.com/slack-go/slack"
)

func (s *Slack) PostTestMessage(ctx context.Context, channelID string) error {
	text := ":rocket: *Post message test* :rocket:\n\nHey, it is working!"

	return s.sendTextMessage(ctx, channelID, text)
}

func (s *Slack) PostReleaseMessage2(ctx context.Context, channelID string) error {

	projectField := slack.AttachmentField{
		Title: "Project",
		Value: "Amazing project",
		Short: true,
	}
	appField := slack.AttachmentField{
		Title: "App",
		Value: "Backend",
		Short: true,
	}
	releaseField := slack.AttachmentField{
		Title: "Release",
		Value: "v0.6.1",
		Short: true,
	}
	deploymentsField := slack.AttachmentField{
		Title: "Deployments",
		Value: "Staging: <https://example.com|example.com>\nProd: <https://example.com|example.com>",
		Short: false,
	}
	githubField := slack.AttachmentField{
		Title: "GitHub",
		Value: "<https://github.com|Release> ",
		Short: false,
	}
	textField := slack.AttachmentField{
		Title: "Changelog",
		Value: "Here is the detailed information about the latest release:\n• Addresses multiple customer feedback points\n• Introduces several new features\n• Includes numerous performance improvements\n• Fixes bugs to enhance user experience\n\nWe encourage everyone to update to the latest version to benefit from these improvements.",
		Short: true,
	}

	// Create the attachment using those fields
	attachment := slack.Attachment{
		Color:  "#36a64f", // Green color, you can choose what you prefer
		Fields: []slack.AttachmentField{projectField, appField, releaseField, deploymentsField, githubField, textField},
		Ts:     json.Number(fmt.Sprintf("%d", time.Now().Unix())),
	}

	// Text Blocks
	headerText := slack.NewTextBlockObject("mrkdwn", "*Hi everyone,*", false, false)
	headerSection := slack.NewSectionBlock(headerText, nil, nil)
	announceText := slack.NewTextBlockObject("mrkdwn", "I'm thrilled to announce a new release! :rocket:", false, false)
	announceSection := slack.NewSectionBlock(announceText, nil, nil)

	// Long Paragraph as the last part of the original message
	longParagraphText := slack.NewTextBlockObject("mrkdwn", "*App*\nBackend\n\n*Release*\nv0.6.1\n\n*Changelog*\nHere is the detailed information about the latest release:\n• Addresses multiple customer feedback points\n• Introduces several new features\n• Includes numerous performance improvements\n• Fixes bugs to enhance user experience\n\nWe encourage everyone to update to the latest version to benefit from these improvements.\n", false, false)
	longParagraphSection := slack.NewSectionBlock(longParagraphText, nil, nil)

	// Combine all blocks and attachment
	messageBlocks := []slack.Block{
		headerSection,
		announceSection,
		slack.NewDividerBlock(),
		longParagraphSection, // Long paragraph added here as part of the original message
	}

	// Post the message with all blocks and attachments
	channelID, _, err := s.client.PostMessage(channelID,
		slack.MsgOptionBlocks(messageBlocks...),
		slack.MsgOptionAttachments(attachment),
		slack.MsgOptionText("", false), // Use MsgOptionText to avoid defaulting to attachment fallback text
	)
	if err != nil {
		fmt.Printf("%s\n", err)
		return err
	}

	return nil
}

func (s *SilentSlack) PostTestMessage(_ context.Context, channelID string) error {
	slog.Debug("no message send to slack because silent slack is used", "channelId", channelID)

	return nil
}
