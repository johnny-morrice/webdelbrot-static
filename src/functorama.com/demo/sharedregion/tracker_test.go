package sharedregion

import (
	"image"
	"testing"
	"time"
	"functorama.com/demo/base"
	"functorama.com/demo/draw"
)

func TestNewRenderTracker(t *testing.T) {
	const jobCount = 5
	const expectWorkers = jobCount - 1
	mock := &MockRenderApplication{}
	mock.SharedConfig.Jobs = uint16(jobCount)
	mock.SharedFactory = &MockFactory{}
	tracker := NewRenderTracker(mock)

	if !(mock.TSharedRegionConfig && mock.TDrawingContext) {
		t.Error("Expected methods not called on mock", mock)
	}

	if tracker == nil {
		t.Error("Expected tracker to be non-nil")
	}

	actualWorkers := len(tracker.workers)
	if actualWorkers != expectWorkers {
		// Expect one less worker as one job is for the tracker
		t.Error("Expected", expectWorkers, "workers but received", actualWorkers)
	}
}

func TestTrackerDraw(t *testing.T) {
	const iterateLimit = 255
	uniform := uniformer()
	point := base.PixelMember{I: 1, J: 2, Member: base.BaseMandelbrot{}}
	context := &draw.MockDrawingContext{
		Pic: image.NewNRGBA(image.ZR),
		Col: draw.NewRedscalePalette(iterateLimit),
	}
	tracker := RenderTracker{
		workerOutput: RenderOutput{
			UniformRegions: make(chan SharedRegionNumerics),
			Children: make(chan SharedRegionNumerics),
			Members: make(chan base.PixelMember),
		},
		context:    context,
	}

	go func() {
		tracker.workerOutput.UniformRegions<- uniform
		close(tracker.workerOutput.UniformRegions)
	}()
	go func() {
		tracker.workerOutput.Members<- point
		close(tracker.workerOutput.Members)
	}()

	packets := tracker.syncDrawing()
	tracker.draw(packets)

	if !(uniform.TRect && uniform.TRegionMember) {
		t.Error("Expected method not called on uniform region:", *uniform)
	}

	if !(context.TPicture && context.TColors) {
		t.Error("Expected method not called on drawing context")
	}
}

func TestTrackerCirculate(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping in short mode")
	}

	tracker := &RenderTracker{
		workerOutput: RenderOutput{
			UniformRegions: make(chan SharedRegionNumerics),
			Children: make(chan SharedRegionNumerics),
			Members: make(chan base.PixelMember),
		},
		workersDone: make(chan bool),
		schedule: make(chan chan<- RenderInput),
		running: true,
	}

	expectedRegion := &MockNumerics{}

	// Feed input
	workerInput := make(chan RenderInput)
	go func() {
		tracker.schedule<- workerInput
	}()
	go func() {
		tracker.workerOutput.Children<- expectedRegion
	}()

	done := make(chan bool, 1)
	go func() {
		tracker.circulate()
		done<- true
	}()

	// Test input
	actualInput := <-workerInput
	abstractNumerics := actualInput.Region
	actualRegion := abstractNumerics.(*MockNumerics)
	if actualRegion != expectedRegion {
		t.Error("Expected", expectedRegion,
			"but received", actualRegion)
	}

	// Test shutdown
	go func() {
		tracker.workersDone<- true
		timeout(t, func() <-chan bool { return done })
	}()
}

func TestTrackerScheduleWorkers(t *testing.T) {
	const jobCount = 2
	tracker := &RenderTracker{
		workers: make([]*Worker, jobCount),
		schedule: make(chan chan<- RenderInput),
		stateChan: make(chan workerState),
		workerOutput: RenderOutput{
			UniformRegions: make(chan SharedRegionNumerics),
			Children: make(chan SharedRegionNumerics),
			Members: make(chan base.PixelMember),
		},
		running: true,
	}

	app := &MockRenderApplication{}
	factory := NewWorkerFactory(app, tracker.workerOutput)

	workerA := factory.Build()
	workerB := factory.Build()

	tracker.workers = []*Worker{workerA, workerB}

	// Run schedule process
	go tracker.scheduleWorkers()

	// Test input scheduling
	go func() {
		workerA.WaitingChan<- true
		close(workerA.WaitingChan)
	}()
	actualA := <-tracker.schedule

	if actualA != workerA.InputChan {
		t.Error("Expected", workerA.InputChan,
			"but received", actualA)
	}

	stateA := <-tracker.stateChan
	expectStateA := workerState{0, true}
	if stateA != expectStateA {
		t.Error("Expected", expectStateA,
			"but received", stateA)
	}

	go func() {
		workerB.WaitingChan<- false
		close(workerB.WaitingChan)
	}()

	stateB := <-tracker.stateChan
	expectStateB := workerState{1, false}
	if stateB != expectStateB {
		t.Error("Expected", expectStateB,
			"but received", stateB)
	}
}

func TestTrackerDetectEnd(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping in short mode")
	}

	const jobCount int = 2

	tracker := &RenderTracker{
		jobs: uint16(jobCount),
		stateChan: make(chan workerState),
		workersDone: make(chan bool),
		workers: make([]*Worker, jobCount),
		workerOutput: RenderOutput{
			UniformRegions: make(chan SharedRegionNumerics),
			Children: make(chan SharedRegionNumerics),
			Members: make(chan base.PixelMember),
		},
	}

	app := &MockRenderApplication{}
	workerFactory := NewWorkerFactory(app, tracker.workerOutput)
	for i := 0; i < jobCount; i++ {
		tracker.workers[i] = workerFactory.Build()
	}

	go func() {
		tracker.stateChan<- workerState{0, true}
	}()

	go func() {
		tracker.stateChan<- workerState{1, true}
	}()

	go tracker.detectEnd()

	timeout(t, func() <-chan bool { return tracker.workersDone })
	close(tracker.stateChan)
}

// todo find a library that does this already
func timeout(t *testing.T, f func() <-chan bool) {
	timer := make(chan bool, 1)
	done := f()
	go func() {
		time.Sleep(1 * time.Second)
		timer <- true
	}()

	select {
	case <-done:
		return
	case <-timer:
		t.Error("Timed out")
	}
}