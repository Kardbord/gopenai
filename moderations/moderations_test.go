package moderations_test

import (
	"testing"

	"github.com/TannerKvarfordt/gopenai/moderations"
)

func TestModerations(t *testing.T) {
	resp, err := moderations.MakeRequest(&moderations.Request{
		Input: "So long, and thanks for all the fish.",
		Model: moderations.ModelStable,
	}, nil)
	if err != nil {
		t.Fatal(err)
	}
	if resp == nil {
		t.Fatal("nil response received")
	}
}
