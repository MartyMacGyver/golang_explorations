
package main

import (
//	"time"
	"fmt"
	"math"
	"sort"
)

func findprimes(cid, min_prime, max_prime int, c *[]chan int) {
    println(fmt.Sprintf("{%d} --> checking (%d, %d]", cid, min_prime, max_prime))
    for j := min_prime; j < max_prime; j++ {
    	if j < 2 {
    		continue
    	}
	    ISPRIME := true
	    for i := 2; i <= int(math.Sqrt(float64(j))); i++ {
	    	if j%i == 0 {
	   			ISPRIME = false
	   			break
	    	}
	    }
	    if ISPRIME {
	   		//amsg := fmt.Sprintf("{%d} --> %d is  prime\n", cid, j)
	    	(*c)[cid] <- j
	    }
    }
    close((*c)[cid])
}

func main() {
	num_gos := 4

    var c []chan int = make([](chan int), num_gos)
	for i := range c {
	   c[i] = make(chan int)
	}
	go findprimes(0,       0, 2000000, &c)
	go findprimes(1, 4000000, 6000000, &c)
	go findprimes(2, 2000000, 4000000, &c)
	go findprimes(3, 6000000, 8000000, &c)
	
	var ok bool
	var msg int
	var primes = make(map[int]bool)
    isRunning := true
    for isRunning {
        select {
            case msg, ok = <-c[0]:
                if ok {
			        //fmt.Printf("%v\n", msg)
			        primes[msg] = true
                } else {
                	println("{0} --> DONE!")
                	c[0] = nil
                }
            case msg, ok = <-c[1]:
                if ok {
			        //fmt.Printf("%v\n", msg)
			        primes[msg] = true
                } else {
                	println("{1} --> DONE!")
                	c[1] = nil
                }
            case msg, ok = <-c[2]:
                if ok {
			        //fmt.Printf("%v\n", msg)
			        primes[msg] = true
                } else {
                	println("{2} --> DONE!")
                	c[2] = nil
                }
            case msg, ok = <-c[3]:
                if ok {
			        //fmt.Printf("%v\n", msg)
			        primes[msg] = true
                } else {
                	println("{3} --> DONE!")
                	c[3] = nil
                }
        }
        if c[0] == nil && c[1] == nil && c[2] == nil && c[3] == nil {
        	isRunning = false
        }
	}

	println("Sorting output")
	var keys []int
	for k := range primes {
	    keys = append(keys, k)
	}
	sort.Ints(keys)
	for _, k := range keys {
	    //fmt.Println("Key:", k, "Value:", primes[k])
	    //fmt.Printf("%d ", k)
	    _=k
	}
	fmt.Printf("Found %d primes\n", len(keys))
}
