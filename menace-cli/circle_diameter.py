#!/usr/bin/env python3

import argparse

def diameter(radius):
    """Return the diameter of a circle given its radius."""
    return 2 * radius

def main():
    parser = argparse.ArgumentParser(description="Compute the diameter of a circle given its radius.")
    parser.add_argument("radius", type=float, help="Radius of the circle")
    args = parser.parse_args()
    result = diameter(args.radius)
    print(f"Diameter: {result}")

if __name__ == "__main__":
    main()
