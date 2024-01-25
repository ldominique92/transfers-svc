```sh
   go get -u github.com/gocraft/dbr/v2 
   go get github.com/mattn/go-sqlite3
   ```


Tests:
- ensure transfered amount
- write this readme


Improvements:
- in a RL scenario, I would write assertions for all the validation messages
- also, I would have a structure to consolidate all my validation messages to improve UX
- It would be nice to have a concurrency test to check transactions for the same user
- I would have a test database and set my connections using env vars
- Unit tests for use case, ensuring that repository unexpected behaviours would be treated correctly
- Could have used go-playground/validator for structs