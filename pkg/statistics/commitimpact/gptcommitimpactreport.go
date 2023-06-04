package commitimpact

import (
	"context"
	"encoding/json"
	"log"
	"math/rand"
	"strings"

	"github.com/claucambra/commit-analysis-tool/pkg/common"
	openai "github.com/sashabaranov/go-openai"
)

const commitLimit = 5
const gptModel = openai.GPT3Dot5Turbo
const maxTokens = 4097

type GPTCommitImpactReport struct {
	Commits common.CommitMap
	Impact  common.CommitMap

	gptClient *openai.Client
}

func NewGPTCommitImpactReport(commits common.CommitMap, gptClient *openai.Client) *GPTCommitImpactReport {
	return &GPTCommitImpactReport{
		Commits:   commits,
		Impact:    common.CommitMap{},
		gptClient: gptClient,
	}
}

func (cir *GPTCommitImpactReport) randomCommits() []*common.Commit {
	numCommits := len(cir.Commits)
	numRandomCommits := common.MinInt(commitLimit, numCommits)

	log.Printf("Choosing %v commits to send to openai for analysis.", numRandomCommits)

	commitSlice := make([]*common.Commit, numCommits)
	randomCommitSlice := make([]*common.Commit, numRandomCommits)

	i := 0
	for _, commit := range cir.Commits {
		commitSlice[i] = commit
		i++
	}

	for i := 0; i < numRandomCommits; i++ {
		randCommitIndex := rand.Intn(numCommits - 1)
		randomCommitSlice[i] = commitSlice[randCommitIndex]
	}

	return randomCommitSlice
}

func (cir *GPTCommitImpactReport) buildPromptString(commits []*common.Commit) string {
	promptString := "The following commit data is formatted as an array of JSON objects. "

	promptString += "Commits containing bodies or subjects which describe new features are highly impactful. "
	promptString += "Commits containing bodies or subjects which describe bug fixes are impactful. "
	promptString += "Commits containing bodies or subjects which describe changing wording or renaming variables are not impactful at all. "
	promptString += "Commits containing bodies or subjects which describe adding test data are not impactful at all. "

	promptString += "Commits containing longer bodies may be more impactful. "

	promptString += "Commits with higher numbers of insertions and deletions tend to be more impactful. "

	promptString += "Commits authored by bots are not very impactful."

	promptString += "Using this information, estimate the impact of each commit on a scale of 0 to 1. "
	promptString += "Do not introduce or explain your answer. "
	promptString += "Only return a JSON array containing each impact value as a floating point number.\n\n"

	promptString += "["
	for _, commit := range commits {
		marshalledCommit, err := json.Marshal(commit)
		if err != nil {
			log.Fatalf("Could not marshal commit %+v, won't be used for impact analysis due to error: %s", commit, err)
		}

		promptString += string(marshalledCommit) + ","
	}
	promptString = strings.TrimSuffix(promptString, ",")
	promptString += "]"

	return promptString
}

func (cir *GPTCommitImpactReport) generatePrompt() string {
	log.Printf("Generating prompt for openai request.")

	promptCommits := cir.randomCommits()
	promptString := cir.buildPromptString(promptCommits)

	return promptString
}

func (cir *GPTCommitImpactReport) Generate() {
	response, err := cir.gptClient.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model:     gptModel,
			MaxTokens: maxTokens,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: cir.generatePrompt(),
				},
			},
		},
	)

	if err != nil {
		log.Fatalf("Error sending completion request to openai: %s", err)
		return
	}

	log.Printf("Received OpenAI response: %+v", response.Choices)
}
