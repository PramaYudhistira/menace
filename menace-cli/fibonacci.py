#!/usr/bin/env python3

import sys

def fib(n):
    a, b = 0, 1
    seq = []
    for _ in range(n):
        seq.append(a)
        a, b = b, a + b
    return seq

def main():
    if len(sys.argv) < 2:
        n = 10
    else:
        try:
            n = int(sys.argv[1])
        except ValueError:
            print("Please provide an integer for the number of terms.")
            sys.exit(1)
    sequence = fib(n)
    print(", ".join(str(num) for num in sequence))

if __name__ == "__main__":
    main()
