package service

import (
	"reflect"
	"strings"
	"testing"
)

func TestParseCaptionResponse_StrictFormat(t *testing.T) {
	raw := "CAPTION: Insane clutch moment!\nHASHTAGS: #Valorant #Clutch"
	caption, hashtags, ok := parseCaptionResponse(raw)

	if !ok {
		t.Fatal("expected ok=true")
	}
	if caption != "Insane clutch moment!" {
		t.Errorf("caption = %q", caption)
	}
	if !reflect.DeepEqual(hashtags, []string{"#Valorant", "#Clutch"}) {
		t.Errorf("hashtags = %v", hashtags)
	}
}

func TestParseCaptionResponse_FreeFormWithInlineHashtags(t *testing.T) {
	// Real response observed from llama3.2:3b during live testing — it
	// ignored the requested CAPTION:/HASHTAGS: format entirely.
	raw := "\"DER ULTIMATE 1V5 CLUTCH! \n Wer hat es geschafft, gegen die Übermacht der Gegner zu kämpfen und den Sieg zu erringen? Watch bis zum Ende zu sehen!\" \n\n#Valorant #GamingHighlights #TwitchClips"

	caption, hashtags, ok := parseCaptionResponse(raw)

	if !ok {
		t.Fatal("expected ok=true even without the strict format")
	}
	if caption == "" {
		t.Error("expected a non-empty caption")
	}
	for _, tag := range hashtags {
		if strings.Contains(caption, tag) {
			t.Errorf("hashtag %q should have been stripped out of the caption body, got caption=%q", tag, caption)
		}
	}
	if !reflect.DeepEqual(hashtags, []string{"#Valorant", "#GamingHighlights", "#TwitchClips"}) {
		t.Errorf("hashtags = %v", hashtags)
	}
}

func TestParseCaptionResponse_Empty(t *testing.T) {
	_, _, ok := parseCaptionResponse("   ")
	if ok {
		t.Error("expected ok=false for an empty response")
	}
}
