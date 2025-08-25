package savefile

import (
	"os"
	"testing"
)

func TestSaveLoadFile(t *testing.T) {
	saver, err := NewLimit("./savefoldertesting/gob", GobCodec{}, 1)
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
	saver, err := NewLimit("./savefoldertesting/", JSONCodec{}, 2)
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
	if len(files) > 4 {
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

func TestManuallyDelete(t *testing.T) {
	saver, err := New("./savefoldertesting/testdeletingfolder", JSONCodec{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	err = saver.Save("hello testworld")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	file, err := os.ReadDir(saver.dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	saver.Delete(file[0].Name())
	_, err = os.Stat(file[0].Name())
	if err == nil {
		t.Fatalf("Expected error of type [*PathError]")
	}
}

func TestDeleteOldNew(t *testing.T) {
	saver, err := New("./savefoldertesting/testdeletingfolder", JSONCodec{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	err = saver.Save("Do you also like Jurrasic Park?")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	file, err := os.ReadDir(saver.dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	saver.DeleteOld()
	_, err = os.Stat(file[0].Name())
	if err == nil {
		t.Fatalf("Expected error of type [*PathError]")
	}
}
