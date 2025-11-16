# Another Example

There are other methods of creating a singleton instance in Go:

1. init function

    We can create a single instance inside the init function. This is only applicable if the early initialization of the instance is ok. The init function is only called once per file in a package, so we can be sure that only a single instance will be created.

2. sync.Once

    The sync.Once will only perform the operation once.