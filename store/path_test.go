package store

import "testing"

func TestPathTransform(t *testing.T) {
	root := "../storage_content"
	key := "somepicture"
	path := hashPathTransform(root, key)
	expectedPathName := "../storage_content/0d1a9/04d68/8adeb/3eab5/8fa69/46e5d/0a5dd/90b0f"
	expectedFileName := "0d1a904d688adeb3eab58fa6946e5d0a5dd90b0f"
	expectedFullPath := expectedPathName + "/" + expectedFileName
	if path.PathName != expectedPathName {
		t.Errorf("expected %s\ngot %s", expectedPathName, path.PathName)
	}
	if path.FileName != expectedFileName {
		t.Errorf("expected %s\ngot %s", expectedFileName, path.FileName)
	}
	if path.FullPath() != expectedFullPath {
		t.Errorf("expected %s\ngot %s", expectedFullPath, path.FullPath())
	}
}
