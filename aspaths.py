import mrtparse as mp
import sys

# generator, yield a path each time.
def get_as_paths(ent):
    if not (13 in ent["type"]):
        return

    if not ((4 in ent["subtype"]) or (10 in ent["subtype"])):
        return

    for rib_ent in ent['rib_entries']:
      for path_attr in rib_ent["path_attributes"]:
          if not (2 in path_attr["type"]):
              continue
          for path_val in path_attr["value"]:
              if not (2 in path_val["type"]):
                  continue
              yield path_val['value']

def main():
    path = sys.argv[1]
    if path == "":
        print("no mrt path provided")
        return

    limit=-1
    if len(sys.argv) >= 3 and sys.argv[2] != "":
        limit=int(sys.argv[2])

    reader = mp.Reader(path)
    pathset = set()
    for chunk in reader:
        entry = chunk.data
        for asp in get_as_paths(entry):
            path_key = ",".join(asp)
            if path_key in pathset:
                continue
            pathset.add(path_key)
            print(path_key)
            limit = limit - 1
            if limit == 0:
                return

if __name__ == '__main__':
    main()
