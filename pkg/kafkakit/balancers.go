package kafkakit

import (
	"hash/fnv"

	"github.com/segmentio/kafka-go"
)

type SPBalancer struct{}

func (b *SPBalancer) Balance(msg kafka.Message, partitions ...int) (partition int) {
	hash := func(s string) int {
		h := fnv.New32a()
		h.Write([]byte(s))
		return int(h.Sum32())
	}

	return partitions[hash(string(msg.Key))%len(partitions)]
}

type SHBalancer struct{}

func (b *SHBalancer) Balance(msg kafka.Message, partitions ...int) (partition int) {
	hash := func(s string) int {
		h := fnv.New32a()
		h.Write([]byte(s))
		return int(h.Sum32()) % 13 * int(h.Sum32()) % 13
	}

	return partitions[hash(string(msg.Key))%len(partitions)]
}
