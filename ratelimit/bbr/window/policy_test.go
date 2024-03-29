package window

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func GetRollingPolicy() *RollingPolicy {
	w := NewWindow(Options{Size: 3})
	return NewRollingPolicy(w, RollingPolicyOpts{BucketDuration: 100 * time.Millisecond})
}

func TestRollingPolicy_Add(t *testing.T) {
	// test func timespan return real span
	tests := []struct {
		timeSleep []int
		offset    []int
		points    []int
	}{
		{
			timeSleep: []int{150, 51},
			offset:    []int{1, 2},
			points:    []int{1, 1},
		},
		{
			timeSleep: []int{94, 250},
			offset:    []int{0, 0},
			points:    []int{1, 1},
		},
		{
			timeSleep: []int{150, 300, 600},
			offset:    []int{1, 1, 1},
			points:    []int{1, 1, 1},
		},
	}

	for _, test := range tests {
		t.Run("test policy add", func(t *testing.T) {
			var totalTS, lastOffset int
			timeSleep := test.timeSleep
			policy := GetRollingPolicy()
			for i, n := range timeSleep {
				totalTS += n
				time.Sleep(time.Duration(n) * time.Millisecond)
				policy.Add(float64(test.points[i]))
				offset, points := test.offset[i], test.points[i]

				assert.Equal(t, points, int(policy.window.buckets[offset].Points[0]),
					fmt.Sprintf("error, time since last append: %vms, last offset: %v", totalTS, lastOffset))
				lastOffset = offset
			}
		})
	}
}

func TestRollingPolicy_AddWithTimespan(t *testing.T) {
	t.Run("timespan < bucket number", func(t *testing.T) {
		policy := GetRollingPolicy()
		// bucket 0
		policy.Add(0)
		// bucket 1
		time.Sleep(101 * time.Millisecond)
		policy.Add(1)
		// bucket 2
		time.Sleep(101 * time.Millisecond)
		policy.Add(2)
		// bucket 1
		time.Sleep(201 * time.Millisecond)
		policy.Add(4)

		for _, bkt := range policy.window.buckets {
			t.Logf("%+v", bkt)
		}

		assert.Equal(t, 0, len(policy.window.buckets[0].Points))
		assert.Equal(t, 4, int(policy.window.buckets[1].Points[0]))
		assert.Equal(t, 2, int(policy.window.buckets[2].Points[0]))
	})

	t.Run("timespan > bucket number", func(t *testing.T) {
		policy := GetRollingPolicy()

		// bucket 0
		policy.Add(0)
		// bucket 1
		time.Sleep(101 * time.Millisecond)
		policy.Add(1)
		// bucket 2
		time.Sleep(101 * time.Millisecond)
		policy.Add(2)
		// bucket 1
		time.Sleep(501 * time.Millisecond)
		policy.Add(4)

		for _, bkt := range policy.window.buckets {
			t.Logf("%+v", bkt)
		}

		assert.Equal(t, 0, len(policy.window.buckets[0].Points))
		assert.Equal(t, 4, int(policy.window.buckets[1].Points[0]))
		assert.Equal(t, 0, len(policy.window.buckets[2].Points))
	})
}
