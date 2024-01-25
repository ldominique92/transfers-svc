<h1>Transfer Service</h1>	
This app was developed as an practical exercise proposed by Qonto.

The goal is to perform a batch transfer for a costumer account.

<h2> How to run </h2>
To run this app please execute in command line:

```sh
   make run
   ```

The service will expose a POST enpoint:

```http://localhost:8081/transfers```

That accepts json requests with the provided format.

<h2>How I approached the problem</h2>
From the problem description I could deduce that the solution should be sychronous and allow concurrence, which I solved determining the level of the transaction isolation. 

<h3>Useful links I bumped into</h3>
- https://stackoverflow.com/questions/129329/optimistic-vs-pessimistic-locking

- https://stackoverflow.com/questions/129329/optimistic-vs-pessimistic-locking

<h2>What I would improve</h2>

- I could have used a library for validation like `go-playground/validator`
- It would be nice to have a concurrency test to check transactions for the same user
- I would have a test database and set my connections using env vars
- Some configuration would also be good, like the port number
- Unit tests for use case, ensuring that repository unexpected behaviours would be treated correctly

<h2>Feedback</h2>

1. How much time did you spend completing this test? 5h to 6h

2. How proud are you of your work? Fairly proud