package internal

import (
	"fmt"
	"strings"
	"time"

	"github.com/gookit/color"
)

type Mood = string

const (
	MoodHappy    Mood = "happy"
	MoodSad      Mood = "sad"
	MoodAngry    Mood = "angry"
	MoodAnxious  Mood = "anxious"
	MoodExcited  Mood = "excited"
	MoodCalm     Mood = "calm"
	MoodStressed Mood = "stressed"
	MoodTired    Mood = "tired"
	MoodNeutral  Mood = "neutral"
)

type Entry struct {
	ID        int       `json:"id" db:"id"`
	Intensity int8      `json:"intensity" db:"intensity"`
	Mood      Mood      `json:"mood" db:"mood"`
	Message   string    `json:"message" db:"message"`
	Tags      []string  `json:"tags" db:"tags"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

func (e *Entry) String() string {
	moodColor := e.getMoodColor()
	intensityStyle := e.getIntensityStyle(moodColor)

	coloredMood := intensityStyle.Sprint(string(e.Mood))

	parts := []string{
		e.CreatedAt.Format("2006/01/02 15:04"),
		fmt.Sprintf("%02d/10 %s", e.Intensity, coloredMood),
	}

	if e.Message != "" {
		parts = append(parts, fmt.Sprintf("\"%s\"", e.Message))
	}

	if len(e.Tags) > 0 {
		parts = append(parts, fmt.Sprintf("[%s]", strings.Join(e.Tags, ", ")))
	}

	return strings.Join(parts, " | ")
}

func (e *Entry) getMoodColor() color.Color {
	switch e.Mood {
	case MoodHappy:
		return color.Green
	case MoodSad:
		return color.Blue
	case MoodAngry:
		return color.Red
	case MoodAnxious:
		return color.Yellow
	case MoodExcited:
		return color.Magenta
	case MoodCalm:
		return color.Cyan
	case MoodStressed:
		return color.LightRed
	case MoodTired:
		return color.Gray
	case MoodNeutral:
		return color.White
	default:
		return color.White
	}
}

func (e *Entry) getIntensityStyle(baseColor color.Color) color.Style {
	intensity := e.Intensity

	switch {
	case intensity >= 8:
		return color.New(baseColor, color.OpBold, color.OpItalic)
	case intensity >= 6:
		return color.New(baseColor, color.OpBold)
	case intensity >= 4:
		return color.New(baseColor)
	case intensity >= 2:
		return color.New(baseColor, color.OpFuzzy)
	default:
		return color.New(baseColor, color.OpFuzzy, color.OpItalic)
	}
}
