package devstreams

import (
	"bias/models"
	"bias/pkg/database/postgres"
	"bias/pkg/database/redis"
	"bias/pkg/twitch"
	"bias/store"
	"context"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/nicklaw5/helix"
	"github.com/sirupsen/logrus"
	"os"
	"strings"
	"sync"
)

func getTags(tagStore *store.TagStore) ([]*models.TagModel, error) {
	tags, err := tagStore.GetAllTags()
	if err != nil {
		return nil, err
	}
	return tags, nil
}

func getCharacterSequenceMatches(keyword, title string) (match bool, matchIndex int) {
	tokens := strings.Fields(strings.ToLower(title))
	for index, token := range tokens {
		if token == keyword {
			return true, index
		}
	}
	return false, -1
}

func checkIfMatchIsDetached(matchIndex int, title string) bool {
	if matchIndex < 0 || matchIndex >= len(title) {
		return false // matchIndex is out of range
	}

	isAtStart := matchIndex == 0 || title[matchIndex-1] == ' '
	isAtEnd := matchIndex == len(title)-1 || title[matchIndex+1] == ' '

	return isAtStart && isAtEnd
}

func extractFullSequence(keyword string, title string, matchIndex int) string {
	startIndex := matchIndex
	endIndex := matchIndex + len(keyword)
	if startIndex > 0 && title[startIndex-1] != ' ' {
		startIndex--
	}
	if endIndex < len(title)-1 && title[endIndex] != ' ' {
		endIndex++
	}
	return title[startIndex:endIndex]
}

func getDistanceBetweenTokens(sequence1, sequence2 string) int {
	sequence1 = strings.ToLower(sequence1)
	sequence2 = strings.ToLower(sequence2)
	if sequence1 == sequence2 {
		return 0
	}

	rows := len(sequence1) + 1
	cols := len(sequence2) + 1
	dp := make([][]int, rows)
	for i := range dp {
		dp[i] = make([]int, cols)
		dp[i][0] = i
	}
	for j := 1; j < cols; j++ {
		dp[0][j] = j
	}

	for i := 1; i < rows; i++ {
		for j := 1; j < cols; j++ {
			insertCost := dp[i][j-1] + 1
			deleteCost := dp[i-1][j] + 1
			replaceCost := dp[i-1][j-1]
			if sequence1[i-1] != sequence2[j-1] {
				replaceCost++
			}
			dp[i][j] = min(insertCost, deleteCost, replaceCost)
		}
	}

	return dp[rows-1][cols-1]
}

func getMatchingTagsByTitle(title string, tags []*models.TagModel) []*models.TagModel {
	var matchingTags []*models.TagModel

	for _, tag := range tags {
		for _, keyword := range tag.Keywords {
			keywordStr := keyword
			match, matchIndex := getCharacterSequenceMatches(keywordStr, title)

			if !match {
				break
			}

			// Checking if a match is detached prevents us from assuming
			// that a match is valid if it is a substring of another word.
			//
			// For example, if we are looking for the keyword "go" and
			// the title is "gopher", we do not want to match "go"
			isDetached := checkIfMatchIsDetached(matchIndex, title)
			if !isDetached {
				break
			}

			matchedToken := extractFullSequence(keywordStr, title, matchIndex)
			charDistance := getDistanceBetweenTokens(keywordStr, matchedToken)

			if charDistance <= 4 {
				matchingTag := tag
				matchingTags = append(matchingTags, matchingTag)

				// Add the parent tag if it exists
				if tag.Parent != nil {
					matchingTags = append(matchingTags, tag.Parent)
				}
			}
		}
	}
	return matchingTags
}

func CollectStreams(clientID string, clientSecret string) error {
	godotenv.Load()

	twitchService := twitch.NewTwitchService(clientID, clientSecret)
	if twitchService == nil {
		return fmt.Errorf("failed to initialize Twitch service")
	}

	categoryIDs := []string{"1469308723", "509670"}

	streams, err := twitchService.GetStreamsByCategory(categoryIDs)
	if err != nil {
		return fmt.Errorf("error getting streams from Twitch: %w", err)
	}

	redisService, err := redis.NewRedisClientService(os.Getenv("REDIS_ADDR"), os.Getenv("REDIS_PASSWORD"), 0)
	if err != nil {
		return fmt.Errorf("error initializing redis service: %w", err)
	}

	connectionString := fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable", os.Getenv("POSTGRES_HOST"), os.Getenv("POSTGRES_USER"), os.Getenv("POSTGRES_PASSWORD"), os.Getenv("POSTGRES_DB"))
	postgresClient, err := postgres.NewPostgresDB(connectionString)

	if err != nil {
		return fmt.Errorf("error initializing postgres service: %w", err)
	}

	redisClient := redisService.GetClient()

	ctx := context.Background()
	fmt.Println("Context: ", ctx)

	streamStore := store.NewStreamStore(postgresClient, redisClient, ctx)
	if streamStore == nil {
		return fmt.Errorf("failed to initialize StreamStore")
	}

	tagStore := store.NewTagStore(postgresClient, redisClient, ctx)
	if tagStore == nil {
		return fmt.Errorf("failed to initialize TagStore")
	}

	tags, err := getTags(tagStore)
	if err != nil {
		return fmt.Errorf("error getting tags: %w", err)
	}

	var wg sync.WaitGroup
	var streamsSaved int

	for _, stream := range streams {
		wg.Add(1)
		go func(s helix.Stream) {
			defer wg.Done()

			newStream := models.StreamModel{
				UserID:      s.UserID,
				UserName:    s.UserName,
				GameID:      s.GameID,
				Title:       s.Title,
				ViewerCount: s.ViewerCount,
				StartedAt:   s.StartedAt,
			}

			matchingTags := getMatchingTagsByTitle(s.Title, tags)
			if matchingTags != nil {
				for _, tag := range matchingTags {
					newStream.Tags = append(newStream.Tags, *tag)
				}
			}

			if err := streamStore.CreateStream(&newStream); err != nil {
				errorMsg := fmt.Sprintf("error storing stream for category %s: %v", categoryIDs, err)
				logrus.Errorf(errorMsg)
			} else {
				streamsSaved++
			}

		}(stream)
	}

	wg.Wait()
	logrus.Infof("Completed stream collection for %d categories: %s", len(categoryIDs), strings.Join(categoryIDs, ", "))
	logrus.Infof("Total streams saved: %d", streamsSaved)

	return nil
}
