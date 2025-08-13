package savefile

import (
	"os"
	"testing"
)

func TestSaveLoadFile(t *testing.T) {
	saver, err := NewLimit("./savefolder/gob", GobCodec{}, 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	data := [3]string{"test", "case", "1"}
	err = saver.Save(data)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var retrievedData [3]string
	err = saver.LoadLatest(&retrievedData)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestSaveLoadFileMultiple(t *testing.T) {
	saver, err := NewLimit("./savefolder/", JSONCodec{}, 2)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	data1 := [3]string{"another", "test", "case"}
	err = saver.Save(data1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	files, err := os.ReadDir(saver.dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(files) > 3 {
		t.Fatalf("exceeded the max number of files")
	}

	var retrievedData [3]string
	err = saver.LoadLatest(&retrievedData)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if retrievedData != data1 {
		t.Fatalf("the loaded data is not the same as the saved data")
	}
}
