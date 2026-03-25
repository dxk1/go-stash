/**
 * Package store
 * @Author: tbb
 * @Date: 2024/10/23 21:18
 */
package store

import "time"

type DataPoint struct {
	FileKey string
	Value   int64
}

type TimeWindow struct {
	Data       map[string]int64
	WindowSize time.Duration
	Duration   time.Duration
	Ticker     *time.Ticker
}

// NewTimeWindow initializes a new TimeWindow with a specified window size and duration.
func NewTimeWindow(windowSize, duration time.Duration) *TimeWindow {
	return &TimeWindow{
		Data:       make(map[string]int64),
		WindowSize: windowSize,
		Duration:   duration,
		Ticker:     time.NewTicker(windowSize),
	}
}

// AddDataPoint adds a new data point to the window.
func (tw *TimeWindow) AddDataPoint(dp DataPoint) {
	tw.Data[dp.FileKey] += dp.Value
}

// Start starts the window ticker to clear old data based on the defined window size and duration.
func (tw *TimeWindow) Start() {
	go func() {
		for range tw.Ticker.C {
			tw.clearOldData()
		}
	}()
}

// Stop stops the ticker.
func (tw *TimeWindow) Stop() {
	tw.Ticker.Stop()
}

// clearOldData clears data that is outside the duration window.
func (tw *TimeWindow) clearOldData() {
	// Add your logic here to clear old data
}

// Example of usage
func main() {
	tenSecWindow := NewTimeWindow(10*time.Second, 30*time.Minute)
	oneMinWindow := NewTimeWindow(1*time.Minute, 1*time.Hour)

	tenSecWindow.Start()
	oneMinWindow.Start()

	// Example data point addition
	tenSecWindow.AddDataPoint(DataPoint{FileKey: "file1", Value: 100})
	oneMinWindow.AddDataPoint(DataPoint{FileKey: "file1", Value: 200})

	// Remember to stop the windows when done
	defer tenSecWindow.Stop()
	defer oneMinWindow.Stop()
}
