package functions

import (
    "sort"
)

func init() {
    Functions = append(Functions, &Function{
        Func: sortmain,
        Path: "/sort",
    })
}

/*
 * Compare two arbitrary go objects.
 * Return true if a is bigger than b, or false otherwise.
 */
func cmp(a interface{}, b interface{}) bool {
    switch aa := a.(type) {
    case string:
        switch bb := b.(type) {
        case string:
            return aa > bb
        case float64:
            panic("cannot sort two difference types")
        }
    case float64:
        switch bb := b.(type) {
        case string:
            panic("cannot sort two difference types")
        case float64:
            return aa > bb
        }
    }
    panic("Unsupported type")
}

// ListSorter
type listSorter []interface{}

func (ls listSorter) Len() int {
    return len(ls)
}

func (ls listSorter) Less(i, j int) bool {
    return cmp(ls[i], ls[j])
}

func (ls listSorter) Swap(i, j int) {
    ls[i], ls[j] = ls[j], ls[i]
}


// MapSorter
type mapSorter []mapPair

type mapPair struct {
    Key string
    Val interface{}
}

func (ms mapSorter) Len() int {
    return len(ms)
}

func (ms mapSorter) Less(i, j int) bool {
    return cmp(ms[i].Val, ms[j].Val)
}

func (ms mapSorter) Swap(i, j int) {
    ms[i], ms[j] = ms[j], ms[i]
}

func sortMapByValue(m map[string]interface{}) mapSorter {
    ms := make(mapSorter, len(m))
    i := 0
    for k, v := range m { 
        ms[i] = mapPair{k, v}
        i++ 
    }
    sort.Sort(ms)
    return ms
}

func sortmain(input interface{}) interface{} {
    output := listSorter(input.([]interface{}))
    sort.Sort(output)
    return output
}
