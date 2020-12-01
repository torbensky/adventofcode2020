# Day 1

## Challenge

### Part 1

Find two numbers that exist in the input list that sum to 2020. The answer to the challenge is the product of these two numbers

### Part 2

Find three numbers in the input list that sum to 2020. The answer to the challenge is the product of these 3 numbers.

## Running

From this directory, `go run cmd/main.go ./input.txt`

## Implementation Notes

- Sorting makes the algorith implementation very easy because we can stop iterating the slice when we exceed the target value `2020`
- The algorithm is basically a permutation generator. Since order doesn't matter, we can solve using some straightforward nested loops.
