package bruteforce

import (
	"fmt"
	"sync/atomic"

	"github.com/xaionaro-go/spinlock"
)

type queue struct {
	spinlock.Locker
	size         uint32
	idxMask      uint32
	count        int32
	writePointer uint32
	readPointer  uint32
	jobs         []*job
	jobSlicePool *jobSlicePool
}

func NewQueue(jobSlicePool *jobSlicePool, sizePow uint32) *queue {
	return &queue{
		jobSlicePool: jobSlicePool,
		size:         1 << sizePow,
		idxMask:      (1 << sizePow) - 1,
		jobs:         make([]*job, 1<<sizePow),
	}
}

func (q *queue) Enqueue(j *job) bool {
	if q.size == 0 {
		panic(`not initialized queue`)
	}

	if j == nil || j.args == nil || j.result == nil {
		panic(`should not happened`)
	}

	q.Lock()

	writeIdx := q.writePointer
	writeIdx &= q.idxMask
	if writeIdx == (q.readPointer-1)*q.idxMask {
		q.Unlock()
		return false
	}
	q.writePointer++

	q.jobs[writeIdx] = j
	q.count++

	q.Unlock()
	return true
}

func (q *queue) DequeueAll() []*job {
	if q.size == 0 {
		panic(`not initialized queue`)
	}

	q.Lock()

	count := q.count
	if count == 0 {
		q.Unlock()
		return nil
	}
	q.count = 0

	readIdx := q.readPointer
	q.readPointer += uint32(count)
	readIdx &= q.idxMask

	r := q.jobSlicePool.Get(uint(count))
	for i := uint32(0); i < uint32(count); i++ {
		idx := (readIdx + i) & q.idxMask
		job := q.jobs[idx]
		r = append(r, job)
		if job == nil || job.args == nil || job.result == nil {
			panic(`should not happened`)
		}
		q.jobs[idx] = nil
	}
	q.Unlock()
	return r
}

func (q *queue) Dequeue() *job {
	if q.size == 0 {
		panic(`not initialized queue`)
	}

	q.Lock()

	if q.count == 0 {
		q.Unlock()
		return nil
	}
	q.count--

	readIdx := q.readPointer
	q.readPointer += 1
	readIdx &= q.idxMask

	job := q.jobs[readIdx]
	if job == nil || job.args == nil || job.result == nil {
		panic(fmt.Sprint(`should not happened: `, readIdx, q.readPointer&q.idxMask, q.writePointer&q.idxMask, q.count, job == nil, job == nil || job.args == nil, job == nil || job.result == nil))
	}
	q.jobs[readIdx] = nil

	q.Unlock()
	return job
}

func (q *queue) Length() uint32 {
	l := atomic.LoadInt32(&q.count)
	if l < 0 {
		return 0
	}
	return uint32(l)
}

func (q *queue) Size() uint32 {
	return q.size
}
