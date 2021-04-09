# Memory Hierarchy

## 1. Does Loop Order Matter?

Option 2 takes almost 5x as long to run as Option 1. Despite the fact that
cachegrind shows them executing the same number of instructions. Examining the
cachegrind output, one can see a dramatic difference in the L1 write misses.
Over 14% of the Option 2 L1 refs are misses, almost 100% of the writes, compared
to less than 1% for Option 1. Summarized into a table:

| Function | L1 wr refs | L1 wr misses | L1 wr miss rate |
| -------- | ---------: | -----------: | --------------: |
| Option 1 | 16,015,088 | 1,000,570    | 0.9%            |
| Option 2 | 16,015,088 | 16,000,570   | 99.9%           |

This can be explained by the pointing out that since a 2D array is laid out in
contiguous rows in memory, subsequent accesses of the same row are "free".
However, filling in the array by column requires jumping over each row and thus
outside of the cache line on each iteration. Of the 16,000,000 elements in the
4000x4000 array, all required a trip to main memory.

Optimizing Option 2 with both `-O1` or `-O0` dramatically improve cache
performance, brining it closer (but not exactly on par) with Option 1
un-optimized. This is largely due to the vastly reduced number of instructions
executed. But also the D1 miss rate falls to 4.7% from the 14.3% un-optimized.
Option 1 optimized is identical to Option 2 optimized, so the compiler (gcc in
this case) was able to figure out they function identically.

The un-optimized cachegrind output and timings follow below.

*Option 1*
```sh
$ time ./option1

real    0m0.087s
user    0m0.056s
sys     0m0.026s
```
```sh
==39210== Command: ./option1
==39210==
--39210-- warning: L3 cache found, using its data for the LL simulation.
--39210-- warning: specified LL cache: line_size 64  assoc 20  total_size 31,457,280
--39210-- warning: simulated LL cache: line_size 64  assoc 30  total_size 31,457,280
==39210==
==39210== I   refs:      240,142,899
==39210== I1  misses:            908
==39210== LLi misses:            897
==39210== I1  miss rate:        0.00%
==39210== LLi miss rate:        0.00%
==39210==
==39210== D   refs:      112,055,362  (96,040,274 rd   + 16,015,088 wr)
==39210== D1  misses:      1,001,895  (     1,325 rd   +  1,000,570 wr)
==39210== LLd misses:      1,001,716  (     1,164 rd   +  1,000,552 wr)
==39210== D1  miss rate:         0.9% (       0.0%     +        6.2%  )
==39210== LLd miss rate:         0.9% (       0.0%     +        6.2%  )
==39210==
==39210== LL refs:         1,002,803  (     2,233 rd   +  1,000,570 wr)
==39210== LL misses:       1,002,613  (     2,061 rd   +  1,000,552 wr)
==39210== LL miss rate:          0.3% (       0.0%     +        6.2%  )
```
*Option 2*
```sh
$ time ./option2

real    0m0.412s
user    0m0.370s
sys     0m0.040s
```
```sh
==39244== Command: ./option2
==39244==
--39244-- warning: L3 cache found, using its data for the LL simulation.
--39244-- warning: specified LL cache: line_size 64  assoc 20  total_size 31,457,280
--39244-- warning: simulated LL cache: line_size 64  assoc 30  total_size 31,457,280
==39244==
==39244== I   refs:      240,142,899
==39244== I1  misses:            907
==39244== LLi misses:            896
==39244== I1  miss rate:        0.00%
==39244== LLi miss rate:        0.00%
==39244==
==39244== D   refs:      112,055,362  (96,040,274 rd   + 16,015,088 wr)
==39244== D1  misses:     16,001,895  (     1,325 rd   + 16,000,570 wr)
==39244== LLd misses:      1,001,716  (     1,164 rd   +  1,000,552 wr)
==39244== D1  miss rate:        14.3% (       0.0%     +       99.9%  )
==39244== LLd miss rate:         0.9% (       0.0%     +        6.2%  )
==39244==
==39244== LL refs:        16,002,802  (     2,232 rd   + 16,000,570 wr)
==39244== LL misses:       1,002,612  (     2,060 rd   +  1,000,552 wr)
==39244== LL miss rate:          0.3% (       0.0%     +        6.2%  )
```
