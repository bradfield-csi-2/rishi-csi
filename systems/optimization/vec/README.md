# Optimization Notes

The functions are numbered based on the order I tried them , roughly following
the same procedure Bryant and O'Halloran show in CS:APP (3e). The profiling was
done very simply using the `clock()` function from `time.h`, as was shown in the
pagecount exercise. I also copied the general methodology from that exercise.

```
Function        Description       Runs    Ticks   Time(s) Avg(s)  Vec Len
dotproduct1     Original          100     823116  0.823   0.0082  1000000
dotproduct2     Move length calc  100     710064  0.710   0.0071  1000000
dotproduct3     No bounds check   100     201389  0.201   0.0020  1000000
dotproduct4     Unroll 2x1        100     170122  0.170   0.0017  1000000
dotproduct      Unroll 3x1        100     145704  0.146   0.0015  1000000
```
