package store

import "testing"

func TestPathTransform(t *testing.T) {
	key := "somepicture"
	path := hashPathTransform(key)
	expectedPathName := "../storage_content/0d1a9/04d68/8adeb/3eab5/8fa69/46e5d/0a5dd/90b0f"
	expectedOriginal := "0d1a904d688adeb3eab58fa6946e5d0a5dd90b0f"

	if path.PathName != expectedPathName {
		t.Errorf("expected %s\ngot %s", expectedPathName, path.PathName)
	}
	if path.FileName != expectedOriginal {
		t.Errorf("expected %s\ngot %s", expectedOriginal, path.FileName)
	}
}
