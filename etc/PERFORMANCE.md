# Performance

Currently I'm using dstask to manage thousands of tasks and the interface still
appears instant.

Dstask currently loads and parses every non-resolved task, each task being a
single file. This may sound wasteful, but it allows git to track history
natively and is actually performant thanks to modern OS disk caches and SSDs.

If it starts to slow down as my number of non-resolved tasks increases, I'll
look into indexing and other optimisations such as archiving really old tasks.
I don't believe that this will be necessary, as the number of open tasks is
(hopefully) bounded.


TODO add `hyperfine` benchmarks.
