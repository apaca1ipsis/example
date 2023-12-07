package kafkakit_test

import (
	"kafka-parallel-queues/pkg/kafkakit"

	"testing"

	"github.com/segmentio/kafka-go"
	"github.com/stretchr/testify/assert"
)

func TestSPBalancer_Balance(t *testing.T) {
	tests := []struct {
		msg        kafka.Message
		partitions []int
		want       int
	}{
		{msg: kafka.Message{Key: []byte("A")}, partitions: []int{0, 1}, want: 0},
		{msg: kafka.Message{Key: []byte("B")}, partitions: []int{0, 1}, want: 1},
		{msg: kafka.Message{Key: []byte("C")}, partitions: []int{0, 1}, want: 0},
	}

	for _, tt := range tests {
		balancer := &kafkakit.SPBalancer{}

		actualPartition := balancer.Balance(tt.msg, tt.partitions...)

		assert.Equal(t, tt.want, actualPartition)
	}
}

func TestSHBalancer_Balance(t *testing.T) {
	tests := []struct {
		msg        kafka.Message
		partitions []int
		want       int
	}{
		{msg: kafka.Message{Key: []byte("A")}, partitions: []int{0, 1}, want: 0},
		{msg: kafka.Message{Key: []byte("C")}, partitions: []int{0, 1}, want: 1},
	}

	for _, tt := range tests {
		balancer := &kafkakit.SHBalancer{}

		actualPartition := balancer.Balance(tt.msg, tt.partitions...)

		assert.Equal(t, tt.want, actualPartition)
	}
}
