# Singleton in Go

## Description

Singleton is a creational design pattern, which ensures that only one object of its kind exists and provides a single point of access to it for any other code.

Singleton has almost the same pros and cons as global variables. Although they’re super-handy, they break the modularity of your code.

You can’t just use a class that depends on a Singleton in some other context, without carrying over the Singleton to the other context. Most of the time, this limitation comes up during the creation of unit tests.

## References

* [Conceptual Example](./example0/readme.md)
* [Another Example](./example1/readme.md)
* [Learn More](https://refactoring.guru/design-patterns/singleton)
* [Singleton in Go](https://refactoring.guru/design-patterns/singleton/go/example)
* [Conceptual Example](https://refactoring.guru/design-patterns/singleton/go/example#example-0)
* [Another Example](https://refactoring.guru/design-patterns/singleton/go/example#example-1)
