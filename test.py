import csv
import datetime
import itertools

n = 1000000
with open('test.csv') as f:
    r = csv.DictReader(f)

    start = datetime.datetime.now()
    for line in itertools.islice(r, n):
        pass
    end = datetime.datetime.now()

print('%.0f ns / op' % (1e9 * (end - start).total_seconds() / n))
