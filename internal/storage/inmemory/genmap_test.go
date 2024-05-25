package inmemory_test

import (
	"github.com/rmntim/ozon-task/internal/storage/inmemory"
	"testing"
)

func TestMap_Delete(t *testing.T) {
	m := inmemory.Map[string, int]{}
	m.Store("1", 2)
	m.Delete("1")

	if _, ok := m.Load("1"); ok {
		t.Error("key should be deleted")
	}
}

func TestMap_Load(t *testing.T) {
	m := inmemory.Map[string, int]{}
	m.Store("1", 2)
	if v, ok := m.Load("1"); !ok || v != 2 {
		t.Error("key should be loaded")
	}
}

func TestMap_LoadAndDelete(t *testing.T) {
	m := inmemory.Map[string, int]{}
	m.Store("1", 2)
	if v, ok := m.LoadAndDelete("1"); !ok || v != 2 {
		t.Error("key should be loaded and deleted")
	}

	if _, ok := m.Load("1"); ok {
		t.Error("key should be deleted")
	}
}

func TestMap_LoadOrStore(t *testing.T) {
	m := inmemory.Map[string, int]{}
	if v, loaded := m.LoadOrStore("1", 2); !loaded || v != 2 {
		t.Error("key should be loaded or stored")
	}

	if v, loaded := m.LoadOrStore("1", 3); loaded || v != 2 {
		t.Error("key should not be loaded or stored")
	}
}

func TestMap_Range(t *testing.T) {
	m := inmemory.Map[string, int]{}
	m.Store("1", 2)
	m.Store("2", 3)

	m.Range(func(key string, value int) bool {
		if key == "1" && value != 2 {
			t.Error("key should be loaded")
		}
		if key == "2" && value != 3 {
			t.Error("key should be loaded")
		}
		return true
	})
}

func TestMap_Store(t *testing.T) {
	m := inmemory.Map[string, int]{}
	m.Store("1", 2)

	if _, ok := m.Load("1"); !ok {
		t.Error("key should be stored")
	}
}
