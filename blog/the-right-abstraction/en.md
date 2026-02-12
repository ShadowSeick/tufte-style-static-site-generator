# The right abstraction
[^sub-header](2026-01-28)
In programming, we focus a lot in the right abstraction, data structure, algorithm or design pattern. We think a lot about the API endpoints or an interface we expose to the end user, but do we think about the implication of what we are building?

We are giving to the users a set of interfaces, visual or programmatic ones, so they can interact with our data, our own abstraction of the business or the real world to perform some action. This data needs to be stored in some place, whether using a cache layer, database, in memory or just a file. Those persistent or not so persistent layers gives us an interface to interact with their abstraction of disk writes and reads which at the same time they are an abstraction of ones and zeroes; it’s just data all the way through.

There is no easy way to create the right shape for the struct, object or class, or the right functions to interact with it, but we are the ones working with the code where this solution should live everyday; we are the ones we should create those solutions.

Maybe we are dealing with existing code and we have to work around it or we need to get out that feature so we cannot focus on what’s really important, data. Most of the time, if we have a good data representation, problems and features become easier to solve and the application becomes simpler and easier to maintain. Technical debt can be done by hurry that feature or not thinking twice about the solution we think is the correct one. It also can be develop by changing the needs of the business. Solutions are a matter of time and space, what was once good could no longer be. Refactors exist for a reason and we should do them from time to time, but sometimes if we spend a little more time thinking on the big picture we could make not so expensive refactors nor making features take longer than should be.

I am not particularly good at doing this, nor in programming in general. This is a craft we need to take care of and practice it, programming is just problem solving in digital paper. Next time you are trying to solve something ask yourself, is this the right abstraction?
