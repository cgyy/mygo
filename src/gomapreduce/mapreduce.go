package gomapreduce

import (
    //    "fmt"
    //    "os"
    "runtime"
    "sync"
)

func MapReduce(mapper func(interface{}) interface{},
    reducer func(chan interface{}) interface{},
    input chan interface{},
    parallel_num int) interface{} {

    // !important for current go version
    runtime.GOMAXPROCS(runtime.NumCPU())

    reduce_input := make(chan interface{})

    // parallel scheduler
    go func() {
        var wg sync.WaitGroup
        task_queue := make(chan interface{}, parallel_num)

        for {
            item := <-input
            if item == nil {
                break
            }

            // block here if task_queue is full
            task_queue <- item
            wg.Add(1)

            // start a goroutine for current task
            go func() {
                reduce_input <- mapper(item)
                item = <-task_queue
                wg.Done()
            }()
        }

        wg.Wait()

        // all input is processed, close reduce input
        close(reduce_input)
    }()

    return reducer(reduce_input)
}
