# Golang Explorations
## Prime number identification

Here we use Go to identify prime numbers in a given range using language features such as channels and goroutines to divide the workload.

While there is a progression of functionality (e.g., command-line argument handling) and some optimization, this is intended as a demonstration of the language and its features, as well as the pros and cons involved as we progress through multiple iterations of the program.

The optimization between primer_5 and primer_6 is of particular interest here - while the number of cycles used is much better (an approximately 6-fold reduction) the timings when spread over 8 cores were not noticeably improved. The timings over two cores on a different machine were, however, also improved 6-fold. Why is this?

