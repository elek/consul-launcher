package consullauncher

import (
	"testing"
	"os"
)

func TestProcessContent(testing *testing.T) {
	content := "{{ env \"YYY\"}}\ntest"
	result := string(templatePlugin.ProcessContent([]byte(content), nil))
	os.ExpandEnv("")
	if (result != "\ntest") {
		testing.Error("result is " + result)
	}

}