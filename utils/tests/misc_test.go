package tests

import (
	"fmt"
	"testing"

	"github.com/andefined/twtarget/utils"
)

func Test_CleanText(t *testing.T) {
	var tests = []struct {
		in  string
		out string
	}{
		{`Διέρρευσε εσωτερικό έγγραφο της ΝΔ για τον #Κουφοντινας 

		#antireport #kou`, "Διέρρευσε εσωτερικό έγγραφο της ΝΔ για τον #Κουφοντινας #antireport #kou"},
		{`See % for your country ↓ 
		🇨🇿 2.2
		🇩🇪 3.1
		🇵🇱 3
		`, "See % for your country ↓ 🇨🇿 2.2 🇩🇪 3.1 🇵🇱 3"},
	}

	for _, test := range tests {
		str := utils.CleanText(test.in)
		if str != test.out {
			fmt.Println(t)
			t.Errorf("Failed on `%s`. Texts don't match (%s)", str, test.out)
		}
	}
}
