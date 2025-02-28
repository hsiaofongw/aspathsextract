import mrtparse as mp
import sys

def get_as_paths(ent, as_paths):
    if not (13 in ent["type"]):
        return

    if not (4 in ent["subtype"]):
        return

    for rib_ent in ent['rib_entries']:
      for path_attr in rib_ent["path_attributes"]:
          if not (2 in path_attr["type"]):
              continue
          for path_val in path_attr["value"]:
              if not (2 in path_val["type"]):
                  continue
              as_paths.append(path_val['value'])

def main():
    path = sys.argv[1]
    if path == "":
        print("no mrt path provided")
        return

    limit=-1
    if len(sys.argv) >= 3 and sys.argv[2] != "":
        limit=int(sys.argv[2])

    reader = mp.Reader(path)
    n=0
    for chunk in reader:
        entry = chunk.data
        as_paths = []
        get_as_paths(entry, as_paths)
        for asp in as_paths:
            print(",".join(asp))
            n=n+1
            if limit != -1 and n>=limit:
                return

if __name__ == '__main__':
    main()
