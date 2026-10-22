#!/usr/bin/env python3
"""Exercise all core snapshots against one isolated PostgreSQL instance."""
from pathlib import Path
import coursecheck
ROOT=Path(__file__).resolve().parents[1]
def main():
    with coursecheck.infrastructure(24) as env:
        for n in range(1,25):coursecheck.check(n,ROOT/'reference'/f'{n:02}',env)
if __name__=='__main__':main()
