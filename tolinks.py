import sys

# generator, yield a path each time.
def get_as_paths(line):
    path = line.strip().split(',')
    for i in range(len(path)-1):
        yield [path[i].strip(), path[i+1].strip()]

def main():
    limit=-1
    if len(sys.argv) >= 2 and sys.argv[1] != "":
        limit=int(sys.argv[1])

    pathset = set()
    for line in sys.stdin:
        for path in get_as_paths(line):
            key = ','.join(path)
            if key in pathset:
                continue
            pathset.add(key)
            print(key)
            limit -= 1
            if limit == 0:
                return

if __name__ == '__main__':
    main()
