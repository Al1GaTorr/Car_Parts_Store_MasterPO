package main

import "context"

func StartLowStockWorker(r *Repo) {
	go func() {
		for a := range r.lowStockCh {
			_ = r.InsertAlert(context.Background(), a)
		}
	}()
}
