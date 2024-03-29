package storages

import (
	"github.com/jitsucom/jitsu/server/adapters"
	"github.com/jitsucom/jitsu/server/appconfig"
	"github.com/jitsucom/jitsu/server/caching"
	"github.com/jitsucom/jitsu/server/counters"
	"github.com/jitsucom/jitsu/server/events"
	"github.com/jitsucom/jitsu/server/logging"
	"github.com/jitsucom/jitsu/server/metrics"
	"github.com/jitsucom/jitsu/server/safego"
	"github.com/jitsucom/jitsu/server/schema"
	"github.com/jitsucom/jitsu/server/telemetry"
	"math/rand"
	"strings"
	"time"
)

type StreamingStorage interface {
	Storage
	Insert(dataSchema *adapters.Table, event events.Event) (err error)
}

//StreamingWorker reads events from queue and using events.StreamingStorage writes them
type StreamingWorker struct {
	eventQueue       *events.PersistentQueue
	processor        *schema.Processor
	streamingStorage StreamingStorage
	eventsCache      *caching.EventsCache
	archiveLogger    *logging.AsyncLogger
	tableHelper      []*TableHelper

	closed bool
}

func newStreamingWorker(eventQueue *events.PersistentQueue, processor *schema.Processor, streamingStorage StreamingStorage,
	eventsCache *caching.EventsCache, archiveLogger *logging.AsyncLogger, tableHelper ...*TableHelper) *StreamingWorker {
	return &StreamingWorker{
		eventQueue:       eventQueue,
		processor:        processor,
		streamingStorage: streamingStorage,
		eventsCache:      eventsCache,
		archiveLogger:    archiveLogger,
		tableHelper:      tableHelper,
	}
}

//Run goroutine to:
//1. read from queue
//2. Insert in events.StreamingStorage
func (sw *StreamingWorker) start() {
	safego.RunWithRestart(func() {
		for {
			if sw.streamingStorage.IsStaging() {
				break
			}
			if sw.closed {
				break
			}

			fact, dequeuedTime, tokenID, err := sw.eventQueue.DequeueBlock()
			if err != nil {
				if err == events.ErrQueueClosed && sw.closed {
					continue
				}
				logging.SystemErrorf("[%s] Error reading event from queue: %v", sw.streamingStorage.ID(), err)
				continue
			}

			//dequeued event was from retry call and retry timeout hasn't come
			if time.Now().Before(dequeuedTime) {
				sw.eventQueue.ConsumeTimed(fact, dequeuedTime, tokenID)
				continue
			}

			batchHeader, flattenObject, err := sw.processor.ProcessEvent(fact)
			if err != nil {
				if err == schema.ErrSkipObject {
					if !appconfig.Instance.DisableSkipEventsWarn {
						logging.Warnf("[%s] Event [%s]: %v", sw.streamingStorage.ID(), sw.streamingStorage.GetUniqueIDField().Extract(fact), err)
					}

					counters.SkipEvents(sw.streamingStorage.ID(), 1)
				} else {
					serialized := fact.Serialize()
					logging.Errorf("[%s] Unable to process object %s: %v", sw.streamingStorage.ID(), serialized, err)

					metrics.ErrorTokenEvent(tokenID, sw.streamingStorage.ID())
					counters.ErrorEvents(sw.streamingStorage.ID(), 1)
					telemetry.Error(tokenID, sw.streamingStorage.ID(), events.ExtractSrc(fact), 1)

					sw.streamingStorage.Fallback(&events.FailedEvent{
						Event:   []byte(serialized),
						Error:   err.Error(),
						EventID: sw.streamingStorage.GetUniqueIDField().Extract(fact),
					})
				}

				//cache
				sw.eventsCache.Error(sw.streamingStorage.ID(), sw.streamingStorage.GetUniqueIDField().Extract(fact), err.Error())

				continue
			}

			//don't process empty object
			if !batchHeader.Exists() {
				continue
			}

			table := sw.getTableHelper().MapTableSchema(batchHeader)

			if err := sw.streamingStorage.Insert(table, flattenObject); err != nil {
				logging.Errorf("[%s] Error inserting object %s to table [%s]: %v", sw.streamingStorage.ID(), flattenObject.Serialize(), table.Name, err)
				if strings.Contains(err.Error(), "connection refused") ||
					strings.Contains(err.Error(), "EOF") ||
					strings.Contains(err.Error(), "write: broken pipe") ||
					strings.Contains(err.Error(), "context deadline exceeded") ||
					strings.Contains(err.Error(), "connection reset by peer") {
					sw.eventQueue.ConsumeTimed(fact, time.Now().Add(20*time.Second), tokenID)
				} else {
					sw.streamingStorage.Fallback(&events.FailedEvent{
						Event:   []byte(fact.Serialize()),
						Error:   err.Error(),
						EventID: sw.streamingStorage.GetUniqueIDField().Extract(flattenObject),
					})
				}

				counters.ErrorEvents(sw.streamingStorage.ID(), 1)
				telemetry.Error(tokenID, sw.streamingStorage.ID(), events.ExtractSrc(fact), 1)

				//cache
				sw.eventsCache.Error(sw.streamingStorage.ID(), sw.streamingStorage.GetUniqueIDField().Extract(fact), err.Error())

				metrics.ErrorTokenEvent(tokenID, sw.streamingStorage.ID())
				continue
			}

			counters.SuccessEvents(sw.streamingStorage.ID(), 1)
			telemetry.Event(tokenID, sw.streamingStorage.ID(), events.ExtractSrc(fact), 1)

			//cache
			sw.eventsCache.Succeed(sw.streamingStorage.ID(), sw.streamingStorage.GetUniqueIDField().Extract(fact), flattenObject, table)

			metrics.SuccessTokenEvent(tokenID, sw.streamingStorage.ID())

			//archive
			sw.archiveLogger.Consume(fact, tokenID)
		}
	})
}

func (sw *StreamingWorker) Close() error {
	sw.closed = true

	return sw.archiveLogger.Close()
}

func (sw *StreamingWorker) getTableHelper() *TableHelper {
	num := rand.Intn(len(sw.tableHelper))
	return sw.tableHelper[num]
}
