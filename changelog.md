# Changelog

- 0.0.1 basic implementation.

  **Performance**
  Each record has between 1 and 5 groupings of metrics, and each grouping having between 1 and 5 metrics. Each metric being a floating point number (cast to an int and held as a string) between 0 and 99.  Each search was done on 1 metric limiting the values to between 0 inclusive and 90 exclusive.

  |Record Count|Initialization|Search
  |--------:|----------:|--------:
  | 1       | 2 ms      | 1 ms    
  | 100     | 47 ms     | 2 ms    
  | 1,000   | 528 ms    | 15 ms  
  | 10,000  | 2,776 ms  | 42 ms
  | 100,000 | 34,553 ms | 544 ms 

