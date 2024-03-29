package window

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestRollingCounterAdd(t *testing.T) {
	size := 3
	bucketDuration := time.Second
	opts := RollingCounterOpts{
		Size:           size,
		BucketDuration: bucketDuration,
	}
	r := NewRollingCounter(opts)
	listBuckets := func() [][]float64 {
		buckets := make([][]float64, 0)
		r.Reduce(func(i Iterator) float64 {
			for i.Next() {
				bucket := i.Bucket()
				buckets = append(buckets, bucket.Points)
			}
			return 0.0
		})
		return buckets
	}
	assert.Equal(t, [][]float64{{}, {}, {}}, listBuckets())
	r.Add(1)
	assert.Equal(t, [][]float64{{}, {}, {1}}, listBuckets())
	time.Sleep(time.Second)
	r.Add(2)
	r.Add(3)
	assert.Equal(t, [][]float64{{}, {1}, {5}}, listBuckets())
	time.Sleep(time.Second)
	r.Add(4)
	r.Add(5)
	r.Add(6)
	assert.Equal(t, [][]float64{{1}, {5}, {15}}, listBuckets())
	time.Sleep(time.Second)
	r.Add(7)
	assert.Equal(t, [][]float64{{5}, {15}, {7}}, listBuckets())
}

func TestRollingCounterReduce(t *testing.T) {
	size := 3
	bucketDuration := time.Second
	opts := RollingCounterOpts{
		Size:           size,
		BucketDuration: bucketDuration,
	}
	r := NewRollingCounter(opts)
	for x := 0; x < size; x++ {
		for i := 0; i <= x; i++ {
			r.Add(1)
		}
		if x < size-1 {
			time.Sleep(bucketDuration)
		}
	}
	var result = r.Reduce(func(iterator Iterator) float64 {
		var result float64
		for iterator.Next() {
			bucket := iterator.Bucket()
			result += bucket.Points[0]
		}
		return result
	})
	if result != 6.0 {
		t.Fatalf("Validate sum of points. result: %f", result)
	}
}

func TestRollingCounterDataRace(t *testing.T) {
	size := 3
	bucketDuration := time.Millisecond * 10
	opts := RollingCounterOpts{
		Size:           size,
		BucketDuration: bucketDuration,
	}
	r := NewRollingCounter(opts)
	var stop = make(chan bool)
	go func() {
		for {
			select {
			case <-stop:
				return
			default:
				r.Add(1)
				time.Sleep(time.Millisecond * 5)
			}
		}
	}()
	go func() {
		for {
			select {
			case <-stop:
				return
			default:
				_ = r.Reduce(func(i Iterator) float64 {
					for i.Next() {
						bucket := i.Bucket()
						for range bucket.Points {
							continue
						}
					}
					return 0
				})
			}
		}
	}()
	time.Sleep(time.Second * 3)
	close(stop)
	assert.Equal(t, nil, nil)
}

func BenchmarkRollingCounterIncr(b *testing.B) {
	size := 3
	bucketDuration := time.Millisecond * 100
	opts := RollingCounterOpts{
		Size:           size,
		BucketDuration: bucketDuration,
	}
	r := NewRollingCounter(opts)
	b.ResetTimer()
	for i := 0; i <= b.N; i++ {
		r.Add(1)
	}
}

func BenchmarkRollingCounterReduce(b *testing.B) {
	size := 3
	bucketDuration := time.Second
	opts := RollingCounterOpts{
		Size:           size,
		BucketDuration: bucketDuration,
	}
	r := NewRollingCounter(opts)
	for i := 0; i <= 10; i++ {
		r.Add(1)
		time.Sleep(time.Millisecond * 500)
	}
	b.ResetTimer()
	for i := 0; i <= b.N; i++ {
		var _ = r.Reduce(func(i Iterator) float64 {
			var result float64
			for i.Next() {
				bucket := i.Bucket()
				if len(bucket.Points) != 0 {
					result += bucket.Points[0]
				}
			}
			return result
		})
	}
}
