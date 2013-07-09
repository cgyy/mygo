package main

import (
    "bufio"
    "fmt"
    "gomapreduce"
    "io"
    "os"
    "regexp"
    "sort"
)

func mapper(filename interface{}) interface{} {
    results := map[string]int{}
    word_regexp := regexp.MustCompile(`[A-Za-z0-9_]+`)
    file, err := os.Open(filename.(string))
    if err != nil {
        return results
    }
    defer file.Close()
    reader := bufio.NewReader(file)
    for {
        line, err := reader.ReadString('\n')
        if err == io.EOF {
            break
        }
        for _, match := range word_regexp.FindAllString(line, -1) {
            results[match]++
        }
    }
    return results
}

func reducer(input chan interface{}) interface{} {
    results := map[string]int{}
    for {
        new_matches := <-input
        if new_matches == nil {
            break
        }
        for key, value := range new_matches.(map[string]int) {
            previous_count, exists := results[key]
            if !exists {
                results[key] = value
            } else {
                results[key] = previous_count + value
            }
        }
    }
    return results
}

type MapPair struct {
    Key string
    Val int
}

type MapPairList []MapPair

func (mpl MapPairList) Len() int {
    return len(mpl)
}

func (mpl MapPairList) Less(i, j int) bool {
    return mpl[i].Val > mpl[j].Val
}

func (mpl MapPairList) Swap(i, j int) {
    mpl[i], mpl[j] = mpl[j], mpl[i]
}

func sortMapByValue(m map[string]int) MapPairList {
    mpl := make(MapPairList, len(m))
    i := 0
    for k, v := range m {
        mpl[i] = MapPair{k, v}
        i++
    }
    sort.Sort(mpl)
    return mpl
}

func main() {
    if len(os.Args[1:]) == 0 {
        fmt.Printf("Usage: wc [<files>]\n")
        fmt.Printf("\n")
        return
    }

    input := make(chan interface{}, len(os.Args[1:]))
    for _, value := range os.Args[1:] {
        input <- value
    }
    close(input)

    a := gomapreduce.MapReduce(mapper, reducer, input, len(input))
    mpl := sortMapByValue(a.(map[string]int))
    for _, value := range mpl {
        fmt.Printf("%s: %d\n", value.Key, value.Val)
    }
}
