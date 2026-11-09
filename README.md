# What's Go Eat Lua

This is a Go project that interfaces with C code through CGO,
with the goal of being able to define and execute Lua Functions
through the C API offered by the Lua Interpreter.

Turns out, this wasn't as easy as I initially thought it would be.
Right now, the whole project is just a minimum REPL for Lua,
with a pretty bare bones interface to interact with the Lua interpreter.

# Usage
```bash
# generate and build and run and whatever else we can do with this thing
make

# just run the c_importer tool and build the go code:
make build

# just run the c_importer tool
make generate
```


# C Importer
Ah, and because I really disliked the way Go expects the C code to be
in a single comment block, a part of this project is a sort of CImporter 
that uses the 'go generate' tool to generate a '*.c.go' file by merging a
'*.c' and a '*.go' file.

If by some horrible reason you want (or even worse, *needs*) to use this 
c_importer thingy, notice the following "convention":
> if your C code needs specfic flags for compiling correctly,
> those flags should be in the very first line of your Go file that contains
> the '//go:generate' directive.
> check the first line of [_interp/lua.go] for an example.


## How to Use C Importer

The c_importer is just another go file that you can find in [_cimport/c_importer.go],
and using it pretty much boils down to putting a comment in the form of
```go
//go:generate go run ../cimport/c_importer.go <filename_to_be_imported.c>
```
anywhere in your go file.

The way CImporter is written was 100% focused on making it work within the context
of this current project. Because of that, it expects the environment variable GOFILE to be set,
since the 'go generate' tool sets that variable to be the name of the file where it found the
'//go:generate' comment, which CImporter reads to duplicate the contents on the output.

I suppose you could do something like:
```bash
export GOFILE="some_go_file_in_which_to_dump_the_contents_of_some_C_file.go"
go run ./cimport/c_importer.go some_c_file_that_will_be_dumped_as_a_comment_in_the_go_file.c
```
but I didn't try it and I don't know why you are even reading this, but I hope
everything goes well for you.

