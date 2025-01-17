package topic

import (
	"fmt"
	"strings"
)

type Topic struct {
	topic string
}

func New(root string) *Topic {
	return &Topic{topic: root}
}

func (t *Topic) Get() string {
	return t.topic
}

func (t *Topic) PartOf(topic string) bool {
	return strings.Index(topic, t.topic) != -1
}

func (t *Topic) Contains(topic string) bool {
	return strings.Index(t.topic, topic) != -1
}

func (t *Topic) Subscribable() string {
	return fmt.Sprintf("%s/#", t.topic)
}

func (t *Topic) Parent() string {
	last_slash := strings.LastIndex(t.topic, "/")
	return t.topic[0:last_slash]
}

func (t *Topic) Appned(node string) {
	t.topic = fmt.Sprintf("%s/%s", t.topic, node)
}
