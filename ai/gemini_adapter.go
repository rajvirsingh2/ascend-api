package ai

import (
	"context"
	"encoding/json"
	_ "encoding/json"
	"fmt"
	"github.com/google/generative-ai-go/genai"
	"github.com/rajvirsingh2/ascend-api/models"
	"google.golang.org/api/option"
	"log"
	_ "log"
	"os"
)

type GeminiAdapter struct {
	client *genai.GenerativeModel
}

func NewGeminiAdapter() (*GeminiAdapter, error) {
	ctx := context.Background()
	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("GEMINI_API_KEY environment variable not set")
	}
	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		return nil, fmt.Errorf("failed to create new Gemini client: %v", err)
	}

	model := client.GenerativeModel("gemini-1.5-flash")
	model.ResponseMIMEType = "application/json"
	model.ResponseSchema = &genai.Schema{
		Type: genai.TypeArray,
		Items: &genai.Schema{
			Type: genai.TypeObject,
			Properties: map[string]*genai.Schema{
				"title":            {Type: genai.TypeString, Description: "Title for the quest"},
				"description":      {Type: genai.TypeString, Description: "Description for the quest"},
				"type":             {Type: genai.TypeString, Enum: []string{"DAILY", "WEEKLY"}, Description: "The type of the quest."},
				"xp":               {Type: genai.TypeInteger, Description: "The XP value awarded(harder quest more the xp)"},
				"attribute_reward": {Type: genai.TypeString, Enum: []string{"STRENGTH", "AGILITY", "INTELLIGENCE", "VITALITY", "SENSE"}, Description: "The reward of the quest."},
			},
			Required: []string{"title", "description", "type", "xp", "attribute_reward"},
		},
	}
	return &GeminiAdapter{client: model}, nil
}

func (a *GeminiAdapter) GenerateQuests(ctx context.Context, goal string, profile models.PlayerProfile) ([]GeneratedQuest, error) {
	prompt := genai.Text(fmt.Sprintf(`
        <CONTEXT>
        You are "The Architect," a hyper-intelligent AI system from the webtoon Solo Leveling, tasked with designing a growth plan for a "Player."
        The Player is at Level %d.
        Their primary long-term goal is: "%s".
        Their current attributes are: Strength %d, Agility %d, Intelligence %d, Vitality %d, Sense %d.
        </CONTEXT>

        <TASK>
        Generate a set of 5 daily quests and 2 weekly quests.
        The quests must be progressive, actionable steps towards the Player's primary goal.
        Tailor the difficulty and type of quests to the Player's current attributes. For example, a player with low 'Intelligence' should receive more foundational knowledge-based quests for a technical goal.
        </TASK>

        <CONSTRAINTS>
        - Daily quests must be tasks completable in 15-60 minutes.
        - Weekly quests must be larger milestones.
        - The 'attribute_reward' for each quest must logically align with the task.
        - Adhere strictly to the provided JSON schema for the response. Do not add any extra text or explanations.
        </CONSTRAINTS>
    `, profile.Level, goal, profile.Strength, profile.Agility, profile.Intelligence, profile.Vitality, profile.Sense))

	resp, err := a.client.GenerateContent(ctx, prompt)
	if err != nil {
		return []GeneratedQuest{}, fmt.Errorf("failed to generate quests: %w", err)
	}

	var quests []GeneratedQuest
	if len(resp.Candidates) > 0 && len(resp.Candidates[0].Content.Parts) > 0 {
		rawJson := resp.Candidates[0].Content.Parts[0].(genai.Text)
		err = json.Unmarshal([]byte(rawJson), &quests)
		if err != nil {
			log.Printf("failed to unmarshal JSON: %v. Raw response: %s", err, rawJson)
			return []GeneratedQuest{}, fmt.Errorf("failed to unmarshal AI response: %w", err)
		}
	} else {
		return []GeneratedQuest{}, fmt.Errorf("failed to generate quests: %w", err)
	}
	return quests, nil
}
