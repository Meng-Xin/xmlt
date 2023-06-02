package utils

import (
	"errors"
	"sync"
	"time"
)

const (
	//机器ID
	workerIDBits = uint64(5)
	//数据中心ID
	dataCenterIDBits = uint64(5)
	//序列号
	sequenceBits = uint64(12)
	//最大机器ID  以8位的有符号整数为例，-1的二进制11111111 其右移workerIDBits(5)位 结果为11100000 ，
	//与11111111异或得00011111 即2^5-1 为最大机器ID
	maxWorkerID = int64(-1) ^ (int64(-1) << workerIDBits)

	maxDataCenterID = int64(-1) ^ (int64(-1) << dataCenterIDBits)
	maxSequence     = int64(-1) ^ (int64(-1) << sequenceBits)
	//分别代表偏移量，以时间戳为例，它离最低位差22位，在分别生成序列号、机器id，时间撮之后需要进行右移，在进行
	//或运算即可得到ID
	timeLeft = uint8(22)
	dataLeft = uint8(17)
	workLeft = uint8(12)

	twepoch = int64(1589923200000)
)

type Worker struct {
	mu           sync.Mutex //锁
	LastStamp    int64      //上一个时间戳
	WorkerID     int64      //机器ID
	DataCenterId int64      //数据中心ID
	Sequence     int64      //序列号
}

func NewWorker(WorkerID, dataCenterID int64) *Worker {
	return &Worker{
		WorkerID:     WorkerID,
		LastStamp:    0,
		Sequence:     0,
		DataCenterId: dataCenterID,
	}
}
func (w *Worker) getMilliSeconds() int64 {
	return time.Now().Unix() / 1e6
}
func (w *Worker) NextID() (uint64, error) {
	w.mu.Lock()
	defer w.mu.Unlock()
	return w.nextID()
}
func (w *Worker) nextID() (uint64, error) {
	timeStamp := w.getMilliSeconds()
	if timeStamp < w.LastStamp {
		return 0, errors.New("time is moving backwards,waiting until")
	}
	if w.LastStamp == timeStamp {
		w.Sequence = (w.Sequence + 1) & maxSequence
		if w.Sequence == 0 {
			//保证时间戳递增
			for timeStamp <= w.LastStamp {
				timeStamp = w.getMilliSeconds()
			}
		}

	} else {
		w.Sequence = 0
	}
	w.LastStamp = timeStamp
	id := ((timeStamp - twepoch) << timeLeft) |
		(w.DataCenterId << dataLeft) |
		(w.WorkerID << workLeft) |
		w.Sequence
	return uint64(id), nil

}
