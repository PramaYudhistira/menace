#!/usr/bin/env python3

import random

def generate_random_number():
    """Return a random float between 0 and 1."""
    return random.random()

if __name__ == "__main__":
    print(generate_random_number())
