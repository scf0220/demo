import re
import argparse
import sys

file_path_seq = "/data/scf/zkdex-wasm-poc/core/bin/okdexd/dev/seq.log"
file_path_prover = "/data/scf/zkdex-wasm-poc/core/bin/okdexd/dev/prover.log"

seq_map = {}
prover_map = {}
hash_block = {}

print(f'seq_log={file_path_seq}')
print(f'prover_log={file_path_prover}')

with open(file_path_prover, "r") as file:
    for line in file:
        if "rerun txs" in line:
            sum_trace = 0
            block_number = 0
            tx_number = 0

            re_match = re.search(r"trace:(\d+)", line)
            if re_match:
                sum_trace = int(re_match.group(1))

            re_match = re.search(r"block:(\d+)", line)
            if re_match:
                block_number = int(re_match.group(1))

            re_match = re.search(r"count:(\d+)", line)
            if re_match:
                tx_number = int(re_match.group(1))

            if block_number in prover_map:
                hash_block[block_number] = True
            else:
                prover_map[block_number] = {
                    "block_number": block_number,
                    "tx_number": tx_number,
                    "sum_trace": sum_trace
                }

with open(file_path_seq, "r") as file:
    for line in file:
        if "propose_controller - submit block" in line:
            sum_trace = 0
            block_number = 0
            tx_number = 0

            re_match = re.search(r"sum:(\d+)", line)
            if re_match:
                sum_trace = int(re_match.group(1))

            re_match = re.search(r"block:(\d+)", line)
            if re_match:
                block_number = int(re_match.group(1))

            re_match = re.search(r"count:(\d+)", line)
            if re_match:
                tx_number = int(re_match.group(1))

            seq_map[block_number] = {
                "block_number": block_number,
                "tx_number": tx_number,
                "sum_trace": sum_trace
            }

print("hash_block:", hash_block)

max_rounded_quotient = 0.0
max_prover_trace = 0
min_prover_trace = 1000000000
max_tx_count = 0
min_tx_count = 10000000000
seq_size = len(seq_map)

for b_num in range(1, seq_size + 1):
    if b_num in hash_block:
        continue

    seq_log = seq_map.get(b_num)
    prover_log = prover_map.get(b_num)

    if not seq_log:
        continue

    if not prover_log:
        continue

    if prover_log["sum_trace"] >= max_prover_trace:
        print("update max prover trace:", b_num, seq_log["tx_number"], seq_log["sum_trace"], prover_log["sum_trace"])
        max_prover_trace = prover_log["sum_trace"]

    if prover_log["sum_trace"] <= min_prover_trace:
        min_prover_trace = prover_log["sum_trace"]
        print("update min prover trace:", b_num, seq_log["tx_number"], seq_log["sum_trace"], prover_log["sum_trace"])

    if max_tx_count < seq_log["tx_number"]:
        max_tx_count = seq_log["tx_number"]

    if min_tx_count > seq_log["tx_number"]:
        min_tx_count = seq_log["tx_number"]

    diff = seq_log["sum_trace"] - prover_log["sum_trace"]
    if seq_log["sum_trace"] < prover_log["sum_trace"]:
        diff = prover_log["sum_trace"] - seq_log["sum_trace"]

    quotient = diff / prover_log["sum_trace"]

    if quotient > max_rounded_quotient:
        max_rounded_quotient = quotient
        print("update max diff rate:", b_num, seq_log["tx_number"], seq_log["sum_trace"], prover_log["sum_trace"],
              max_rounded_quotient)

print(f'seqLen={len(seq_map)} proverLen={len(prover_map)} max_trace={max_prover_trace} min_trace={min_prover_trace}')
if max_rounded_quotient > 0.03:
    print(f'ERROR max-diff={max_rounded_quotient}')
else:
    print(f'max diff={max_rounded_quotient}')
