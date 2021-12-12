# palette
Just colors. Go.


## Contributing
I welcome palettes of all kinds but just to be sure file an issue with the palette(s) you wish to add.

Add your colors in palettes.txt. If you wish to add your colors in a new format edit 
`parseGenericColor` in [`parse.go`](./generate_palettes/parse.go).

Then run `go generate` in base directory and make sure `palettes.go` is gofmt'ed.