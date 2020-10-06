# Amortize the cost of importing pandas.
from pandas import read_csv


def read_rows(in_stream):
    for line in in_stream:
        yield line.split('\t')


def count_gen():
    with open('data/cities500.txt') as in_stream:
        population_south = 0
        population_total = 0
        row_count = 0
        for row in read_rows(in_stream):
            latitude = float(row[4])
            population = int(row[14])
            population_total += population
            if latitude < 0:
                population_south += population
            row_count += 1
        return dict(
            population_south=population_south,
            population_total=population_total,
            row_count=row_count)


def count_pandas():
    # from pandas import read_csv
    fields = dict(lat=4, pop=14)
    cities = read_csv(
        'data/cities500.txt', header=None, sep='\t',
        names=fields.keys(),
        usecols=fields.values())
    return dict(
        population_south=cities[cities['lat'] < 0]['pop'].sum(),
        population_total=cities['pop'].sum(),
        row_count=cities.shape[0])


def main():
    from sys import argv
    mode = dict(gen=count_gen, pandas=count_pandas)
    result = mode.get(argv[1])()
    print(f'South population: {result["population_south"]}')
    print(f'Total population: {result["population_total"]}')
    print(f'# of rows: {result["row_count"]}')


if __name__ == "__main__":
    main()
