import csv
import datetime
import itertools
import contextlib
import mmap

@contextlib.contextmanager
def csv_reader(filename):
    with open('test.csv') as f:
        data = mmap.mmap(f.fileno(), 0, prot=mmap.PROT_READ)
        try:
            def generator():
                start = 0

                while True:
                    eol = data.find(b'\n', start)
                    if eol == -1:
                        break
                    yield data[start:eol].split(b',')
                    start = eol + 1

            yield generator()
        finally:
            data.close()

n = 1000000
with open('test.csv') as f:
    r = csv.DictReader(f)

    start = datetime.datetime.now()
    for line in itertools.islice(r, n):
        pass
    end = datetime.datetime.now()
print('python csv\t%.0f ns / op' % (1e9 * (end - start).total_seconds() / n))

n = 1000000
with csv_reader('test.csv') as r:
    start = datetime.datetime.now()
    for line in itertools.islice(r, n):
        pass
    end = datetime.datetime.now()
print('python mmap-csv\t%.0f ns / op' % (1e9 * (end - start).total_seconds() / n))
