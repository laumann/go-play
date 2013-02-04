package gocloth

import "testing"

var testInputExpected = map[string]string {
	"Simple text": "<p>Simple text</p>",
}

func TestLexSimple(t *testing.T) {
	l := lex("", "Simple text")
	go l.run()
	output := <-l.items
	
	if output.val != "<p>Simple text</p>" {
		t.Logf("\nGot: '%s'\nExpected: '%s'\n", output.val, "<p>Simple text</p>")
		t.FailNow()
	}

}
