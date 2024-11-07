package topic_test

import (
	"carousel/infrastructure/broker/topic"
	"testing"
)

func TestTopicAppend(t *testing.T) {
	tt := "/hello"
	expect := "/hello/world"
	tp := topic.New(tt)
	tp.Appned("world")
	if tp.Get() != expect {
		t.Errorf("Expected '%s', but got '%s'", expect, tp.Get())
	}
}

func TestTopicParent(t *testing.T) {
	initial := "/hello/world/topic"
	expect := "/hello/world"
	tp := topic.New(initial)
	parent := tp.Parent()
	if parent != expect {
		t.Errorf("Expected '%s', but got '%s'", expect, parent)
	}
}

func TestTopicSibscribable(t *testing.T) {
	initial := "/hello/world/topic"
	expect := "/hello/world/topic/#"
	tp := topic.New(initial)
	subscribable := tp.Subscribable()
	if subscribable != expect {
		t.Errorf("Expected '%s', but got '%s'", expect, subscribable)
	}
}

func TestTopicPartOf(t *testing.T) {
	full := "/hello/world/topic"
	part := "/hello/world"
	tp := topic.New(full)
	if tp.PartOf(part) {
		t.Errorf("Expected True, in contains comparison of '%s' and '%s'", full, part)
	}
}

func TestTopicContains(t *testing.T) {
	full := "/hello/world/topic"
	part := "/hello/world"
	tp := topic.New(part)
	if tp.Contains(full) {
		t.Errorf("Expected True, in contains comparison of '%s' and '%s'", full, part)
	}
}
