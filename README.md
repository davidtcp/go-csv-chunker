# go-csv-chunker
CSV Chunker

Split big CSV file to many small chunks with safe header

```bash
$ csv-chunker big_file.csv 42mb
```

big_csv.csv
```csv
id,name,age
1,alice,42
2,bob,24
```

after chunker worked

chunk_1.csv
```
id,name,age
1,alice,42
```

chunk_2.csv
```
id,name,age
2,bob,24
```
