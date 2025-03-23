# Following the book

> please excuse my spelling mistakes.
> i'm also using markdownlint so if the notes look weird thats why.
> these are somewhat personal and fairly silly,
> but i try to write somewhat helpful info and
> it's good practice to put my thoughts on paper, so i've provided these notes
> later at the end i've compiled the more helpful stuff i've learned.
> notes appended with |ATE| are items that will be Appended To End.
> entire sections that are important will not be |ATE|

- 0.00:
no notes so far in this section.

- 02.01:
im not really sure what to put for the mod name.
all of my projects have been just public repos and not hosted.
i'll just stick with my github, github.com/Bekian/snippetbox.
i forgot about the command `go run .`
i was using `go build` then running the executable created from that.

at the end of the chapter they say:

"If you’re creating a project which can be downloaded and
used by other people and programs,
then it’s good practice for your module path to equal
the location that the code can be downloaded from.
For instance, if your package is hosted at `https://github.com/foo/bar`
then the module path for the project should be `github.com/foo/bar`."
which is excatly what i did on accident, lol.

- 02.02:
additional info section:
they showed that you can use `:http` or `:http-alt`
to use the default port that your computer uses
instead of a custom port for a specific protocol.
i selected this but i didnt know what the port was,
luckily the book says you can usually find it located inside the file `/etc/services`.
i went to the file to search through the file,
but there were too many results and i didnt know how to do this,
so i consulted chatgpt on how i can search through a file,
when i know the file contained numerous results
by using `grep "http" services` and
the terminal im using doesnt allow command scrolling by default,
so i asked if chatgpt knew how to fix this and
it said i could pipe the results into `less` and
that would give a scrollable option,
like `grep "http" services | less`, this is super helpful, will use again.
|ATE|

before starting 02.03, im going to setup the git project and publish to gh

- 02.03:
didnt know you could use `{$}` after a path
to restrict a route pattern to a specific route
like `\{$}` explicitly matches `\` instead of
matching anything after the slash that isnt defined.

they list some additional information about uri pattern matching
but i dont think that is super helpful to make note of here.

they also discuss the old `http.HandleFunc()` methods,
which use a global `http.DefaultServeMux` variable,
because it is global, a third party can register an unsafe route
and thus, should not be used.

- 02.04:
a routes wildcard pattern is denoted by
an identifier surrounded by curly brackets `{example}`

the identifier *must* fill the entire space between slashes,
`/{example}.file` or `/id_{example}` are not valid patterns.
additionally, only one wildcard per path segment can be present,
`/{example1}-{example2}/` is not a valid pattern.

to access the identifiers value, use the function `PathValue()` on the request object
like `example_value := r.PathValue("example")` see book or godocs for more info
the value is a string and should *always* be validated

in the additional info, they mention that there could be pattern overlap
so caution must be used when defining patterns to avoid overlap.
