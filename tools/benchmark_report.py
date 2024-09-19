import os
import sys
import pandas as pd

def get_test_name_pretty(input_string):
    return input_string.replace("_", " ").title()

def parse_benchmark_file(filepath):
    """Parses the benchmark file to extract the p50 time in milliseconds."""
    with open(filepath, 'r') as file:
        for line in file:
            if line.startswith("p50"):
                parts = line.strip().split(" ")
                # Assume format is "p50 - <time_in_ms> ms"
                time_ms = float(parts[2])
                return time_ms
    return None

def dataframe_to_markdown(df, output_file):
    """Converts a DataFrame to a Markdown table and writes to a file."""
    with open(output_file, 'w') as f:
        f.write('| ' + ' | '.join(df.columns) + ' |\n')
        f.write('|' + '---|' * len(df.columns) + '\n')
        for _, row in df.iterrows():
            f.write('| ' + ' | '.join(f'{x:.4f}' if isinstance(x, float) else str(x) for x in row) + ' ms |\n')

def main(folder_path):
    benchmarks = {}

    for filename in os.listdir(folder_path):
        if filename.startswith("benchmark_") and filename.endswith(".txt"):
            parts = filename.split("_")
            test_name = get_test_name_pretty("_".join(parts[1:-2]))
            with_aikido = "with_aikido" in filename

            filepath = os.path.join(folder_path, filename)
            time_ms = parse_benchmark_file(filepath)

            if test_name not in benchmarks:
                benchmarks[test_name] = {}
            if with_aikido:
                benchmarks[test_name]['with_aikido'] = time_ms
            else:
                benchmarks[test_name]['without_aikido'] = time_ms

    rows = []
    for test_name, times in benchmarks.items():
        without_aikido = times.get('without_aikido', 0)
        with_aikido = times.get('with_aikido', 0)
        difference_ms = with_aikido - without_aikido
        difference_pct = (difference_ms / without_aikido) * 100 if without_aikido != 0 else 0
        rows.append({
            'Benchmark': test_name,
            'Avg. time w/o Zen': without_aikido,
            'Avg. time w/ Zen': with_aikido,
            'Delta': difference_ms,
            'Delta in %': difference_pct,
        })

    df = pd.DataFrame(rows)

    output_file = os.path.join(folder_path, "benchmark_results.md")
    dataframe_to_markdown(df, output_file)
    print(f"Markdown table has been saved to {output_file}")

if __name__ == "__main__":
    if len(sys.argv) < 2:
        print("Usage: python script.py <folder_path>")
        sys.exit(1)
    
    folder_path = sys.argv[1]
    if not os.path.isdir(folder_path):
        print(f"Error: The path {folder_path} is not a directory or does not exist.")
        sys.exit(1)

    main(folder_path)
