import sys

# generator, yield a segment each time.
def get_as_paths(line):
    path = line.strip().split(',')
    for seg in path:
        yield seg.strip()

def main():
    limit=-1
    if len(sys.argv) >= 2 and sys.argv[1] != "":
        limit=int(sys.argv[1])

    pathset = set()
    for line in sys.stdin:
        for seg in get_as_paths(line):
            if seg in pathset:
                continue
            pathset.add(seg)
            print(seg)
            limit -= 1
            if limit == 0:
                return

if __name__ == '__main__':
    main()
