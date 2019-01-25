## File download and concatenation benchmarking

It is assumed that the bucket contains files `one.mp3`, `two.mp3` and `three.mp3`.

To benchmark two files being concatenated:

```
BUCKET="https://s3-us-west-2.amazonaws.com/somebucket/optional-path" go test -bench=BenchmarkConcat2 -count=10
```

To benchmark three files being concatenated:

```
BUCKET="https://s3-us-west-2.amazonaws.com/somebucket/optional-path" go test -bench=BenchmarkConcat3 -count=10
```

Tests can be run in the usual way:


```
BUCKET="https://s3-us-west-2.amazonaws.com/somebucket/optional-path" go test ./...
```
