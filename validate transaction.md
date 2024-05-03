For a bitcoin node, one of its major task is to velidate a transaction, there are several steps to take for it, the first thing is to check the output can match to the transaction. For example if a transaction
is about "jim using 10 dollars to by a cup of coffee with price of 3 dollars", then we need to check :

1, jim really has 10 dollars

2, the amount left after buying the coffee should be 7 dollars

If the transaction is honest, then the input of the transaction(10 dollars) should greater than the output of the transaction(7 dollars), that is when we use the amount of input minus the amount of the output
the result should be positive, if the result is negative, then the transaction is "dishonest" it want to fake money from air. We use following code to compare the input amount and output amont:
```g
```
