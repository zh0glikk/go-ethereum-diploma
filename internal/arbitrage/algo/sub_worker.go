package algo

func runSubWorker(ch chan *SubWorkerJob, chResp chan *SubWorkerJobResponse, finishCh chan struct{}) {
	for {
		select {
		case jobToProc := <-ch:
			results, profit := jobToProc.F(
				jobToProc.Point,
			)

			jobResponse := SubWorkerJobResponse{
				Point:     jobToProc.Point,
				Index:     jobToProc.Index,
				Profit:    profit,
				Execution: results,
			}

			chResp <- &jobResponse
		case _, ok := <-finishCh:
			if ok {
				continue
			}
			return
		}
	}
}
