# metadata-parser

A Chef cookbook metadata dependency parser lib for Golang


## What, WHY?

Because although the chef-server understands and spews metadata as JSON all of the userspace tooling speaks the ruby form (see knife).
That ruby form is what everyone puts in cookbooks, and in almost all cases the .json form isn't checked into repos :|

I needed a way to dependably parse name,version, and dependency info from cookbook metadata. I could have regexed, but I was interested in writing a toke-parser in Golang

## Usage
Using this is pretty simple, just create a new parser from any io.Reader, and `Parse()`.

```
meta, err := metadata.NewParser(strings.NewReader(metadata)).Parse()
if err != nil {
  fmt.Println(" Error parsing: ", err)
  t.Fail()
}

fmt.Println("Name: ", meta.Name)
fmt.Println("Version:", meta.Version)
for i, dep := range(meta.Depends) {
  fmt.Printf("\tdepends: %s@%s\n", dep.Name, dep.Constraint)
}

```

## Contributing

Please contribute and help improve this project!

- Fork the repo
- Make sure the tests pass
- Improve the code
- Make sure your feature has test coverage
- Make sure the tests pass
- Submit a pull request

## Thanks
Thanks to the InfluxDB team for MIT licensing their sql parser. Much of this code is derived from their impelemtnation: https://github.com/influxdb/influxdb
Also thanks to GopherAcademy for this awesome article on writing lexers and parsers in Go: http://blog.gopheracademy.com/advent-2014/parsers-lexers/
