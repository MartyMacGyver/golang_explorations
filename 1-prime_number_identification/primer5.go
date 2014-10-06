
package main

// In the range [0, 2^24-1] there are  1,077,871 primes (sum = 8729068693022)
// In the range [0, 2^28-1] there are 14,630,843 primes (sum = ?)

import (
	"fmt"
	"runtime"
	"reflect"
	"time"
	"math"
	"sort"
	"flag"
)

type Uint64Sorter []uint64
func (s Uint64Sorter) Len() int {
    return len(s)
}
func (s Uint64Sorter) Swap(i, j int) {
    s[i], s[j] = s[j], s[i]
}
func (s Uint64Sorter) Less(i, j int) bool {
    return s[i] < s[j]
}

type WorkerInput struct {
    stepNum uint64
    stepIncomplete uint64
    minVal uint64
    maxVal uint64
}

type WorkerOutput struct {
    stepNum uint64
    primes []uint64
    statsCycles uint64
    statsPasses uint64
}

func checkIfPrime(val uint64) (IsPrime bool, cycles uint64) {
	IsPrime = true
	cycles++
	if val < 2 {
		IsPrime = false
		return
	}
    for i := uint64(2); i <= uint64(math.Sqrt(float64(val))); i++ {
    	cycles++
    	if val % i == 0 {
   			IsPrime = false
   			break
    	}
    }
    return
}

func myWorker(cid int, cIn chan WorkerInput, cOut chan WorkerOutput) {
	fmt.Printf("{%d} ++ myWorker initialized\n", cid)
    for {
        select {
            case workUnit, ok := <- cIn:
                if ok {
                	primes := make([]uint64, 0)
				    statsPasses := uint64(0)
				    statsCycles := uint64(0)
                	for i := workUnit.minVal; i <= workUnit.maxVal; i++ {
					    IsPrime, statsCyclesTemp := checkIfPrime(i)
					    if IsPrime {
					    	primes = append(primes, i)
					    }
					    statsPasses++
					    statsCycles += statsCyclesTemp
					}
				    cOut <- WorkerOutput{workUnit.stepNum, primes, statsCycles, statsPasses}
                } else {
               		fmt.Printf("{%d} --> cIn not OK!\n", cid)
                }
        }
    }
    close(cOut)
	fmt.Printf("{%d} xx myWorker exiting\n", cid)
}

func main() {
	numWorkers := runtime.NumCPU()
	runtime.GOMAXPROCS(numWorkers)

 	valMin  := uint64(0)
	valMax  := uint64(0x01000000-1)
	valStep := uint64(numWorkers*256)

	flag.Uint64Var(&valMin,  "min",  valMin,  "Minimum value")
	flag.Uint64Var(&valMax,  "max",  valMax,  "Maximum value")
	flag.Uint64Var(&valStep, "step", valStep, "Values per chunk")
	flag.Parse()
	
	steps   := uint64((valMax-valMin+1)/valStep)+1  // Boundary case protection
	stepCur := uint64(0)
	stepSlice := make([]bool, steps)
	stepIncomplete := uint64(0)

    fmt.Println("Min, max, step, steps = ", valMin, valMax, valStep, steps)

	timeStart := time.Now()

	gCycles := make([]uint64, numWorkers)
	gPasses := make([]uint64, numWorkers)

    var downChans []chan WorkerOutput = make([](chan WorkerOutput), numWorkers)
	for i := range downChans {
		downChans[i] = make(chan WorkerOutput)
	}
	upChan := make(chan WorkerInput)

	for i := 0; i < numWorkers; i++ {
		fmt.Printf("Creating worker %d\n", i)
		go myWorker(i, upChan, downChans[i])
	}

	cases := make([]reflect.SelectCase, numWorkers+1) // Add room for default
	for i, ch := range downChans {
	    cases[i] = reflect.SelectCase{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(ch)}
	}
	defaultCid := numWorkers
    cases[defaultCid] = reflect.SelectCase{Dir: reflect.SelectDefault}

	var primes = make(map[uint64]bool)

	activeWorkers := 0
	poolMin:=valMin
	for i := 0; i < numWorkers && poolMin <= valMax; i++ {
		poolMax := poolMin+valStep-1
		if poolMax > valMax {
			poolMax = valMax
		}
		upChan <- WorkerInput{stepCur, stepIncomplete, poolMin, poolMax}
		poolMin = poolMax+1
		stepCur++
		activeWorkers++
	}

	for activeWorkers > 0 {
		var ok bool
		var msg reflect.Value
		var cid int
		//time.Sleep(time.Millisecond*00)
		cid, msg, ok = reflect.Select(cases)
		if cid == defaultCid {
			//fmt.Printf("{%d} <-- DEFAULT\n", cid)
			continue
		} else if ok {
			newPrimes := msg.FieldByName("primes")
			stepNum := msg.FieldByName("stepNum").Uint()

			//fmt.Printf("Stepnum = %d\n", stepNum)

			for i := 0; i < newPrimes.Len(); i++ {
				value := newPrimes.Index(i).Uint()
				//fmt.Printf("\t%d --> %d\n", i, value)
				primes[value] = true
			}

			stepSlice[stepNum] = true
			for i := range(stepSlice) {
				if !stepSlice[i] {
					//fmt.Printf("Ping on %d\n", i)
					stepIncomplete = uint64(i)
					_ = stepIncomplete
					break
				}
			}

			statsCycles := msg.FieldByName("statsCycles").Uint()
			statsPasses := msg.FieldByName("statsPasses").Uint()
		    gCycles[cid]+=statsCycles
		    gPasses[cid]+=statsPasses
		    
			if (poolMin <= valMax) {
	    		poolMax := poolMin+valStep-1
				if poolMax > valMax {
					poolMax = valMax
				}
				upChan <- WorkerInput{stepCur, stepIncomplete, poolMin, poolMax}
				poolMin = poolMax+1
				stepCur++
			} else {
				activeWorkers--
				fmt.Printf("activeWorkers channels: %d\n", activeWorkers)
			}
		} else {
			cases[cid].Chan = reflect.ValueOf(nil) // Channel closed
			fmt.Printf("{%d} --> DONE!\n", cid)
			activeWorkers -= 1
			continue
		}
		//fmt.Printf("Read from channel %#v and received %d\n", downChans[cid], msg.Uint())
	}

	timeEnd := time.Now()
	timeDelta := timeEnd.Sub(timeStart)

	time.Sleep(time.Millisecond*1)

	fmt.Printf("Sorting output\n")
	var keys []uint64
	for k := range primes {
	    keys = append(keys, uint64(k))
	}
	sort.Sort(Uint64Sorter(keys))
	sumOfPrimes := uint64(0)
	for _, k := range keys {
	    sumOfPrimes+=k
	}

	totalCycles := uint64(0)
	totalPasses := uint64(0)
	for i := 0; i < numWorkers; i++ {
		totalCycles+=gCycles[i]
		totalPasses+=gPasses[i]
	}
	fmt.Printf("%d workers ran %12d cycles          for %11d values (%5.2f%% = ideal)\n", numWorkers,
		totalCycles, totalPasses, 100.0/float64(numWorkers))
	for i := 0; i < numWorkers; i++ {
		fmt.Printf("Worker #%d ran %12d cycles (%5.2f%%) for %11d values (%5.2f%%)\n", i,
			gCycles[i], float64(gCycles[i])/float64(totalCycles)*100.0,
			gPasses[i], float64(gPasses[i])/float64(totalPasses)*100.0,
		)
	}
	fmt.Printf("Found %d primes (sum = %d) in %s\n", len(keys), sumOfPrimes, timeDelta)
}
