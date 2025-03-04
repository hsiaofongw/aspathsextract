#!/usr/bin/env python

import json
import sys
from mrtparse import *

def main():
    mrtFilePath = ""
    if len(sys.argv) >= 2 and sys.argv[1] != "":
        mrtFilePath = sys.argv[1]
    
    if mrtFilePath == "":
        print("no mrtdump file path given")
        sys.exit(1)

    limit=-1
    if len(sys.argv) >= 3 and sys.argv[2] != "":
        limit=int(sys.argv[2])

    entries = []
    for entry in Reader(sys.argv[1]):
        entries.append(entry.data)
        limit -= 1
        if limit == 0:
            break
    sys.stdout.write(json.dumps(entries, indent=2))

if __name__ == '__main__':
    main()
