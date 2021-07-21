# compressor

## example
```
// add pack
pack := compressor.NewPack(".")
```

```
// add file or directory into pack
// if enter a directory, it will adds all sub-items of a directory.
pack.AddItem("verybigdata.txt")
```

```
tgz, _ := os.Create("verybigdata.txt.tgz")
gzw := gzip.NewWriter(tgz)

// archive the files to "gzw" in the pack.
err := pack.WriteTar(gzw)
```

## TODO::
zip   
multicore archiving