#!/usr/bin/env python3
# merge_sorted_chunks.py
# Usage: python3 merge_sorted_chunks.py <chunk_dir> <out_file>
# Merges sorted chunk files (ending with .sorted) from chunk_dir in a memory-efficient way,
# writing deduplicated lines to out_file.

import sys
import os
import heapq

def main():
    if len(sys.argv) != 3:
        print("Usage: merge_sorted_chunks.py <chunk_dir> <out_file>", file=sys.stderr)
        sys.exit(2)
    chunk_dir = sys.argv[1]
    out_file = sys.argv[2]

    if not os.path.isdir(chunk_dir):
        print(f"chunk_dir not found: {chunk_dir}", file=sys.stderr)
        sys.exit(3)

    files = sorted(
        os.path.join(chunk_dir, f)
        for f in os.listdir(chunk_dir)
        if f.endswith(".sorted")
    )

    if not files:
        print("No .sorted chunk files found in " + chunk_dir, file=sys.stderr)
        sys.exit(4)

    fds = []
    try:
        for f in files:
            # open in text mode; errors='replace' to be robust to encoding issues
            fds.append(open(f, "r", encoding="utf-8", errors="replace"))

        merged = heapq.merge(*fds)
        with open(out_file, "w", encoding="utf-8") as out:
            last = None
            for line in merged:
                line = line.rstrip("\n")
                if not line:
                    continue
                if line == last:
                    continue
                out.write(line + "\n")
                last = line
    finally:
        for fd in fds:
            try:
                fd.close()
            except Exception:
                pass

if __name__ == "__main__":
    main()
