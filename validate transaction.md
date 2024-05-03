For now we are still "off the hood" which means we still playing by our own, we need to "on the chain" which means we need to meangle with nodes for the bitcoin blockchain and act as one of them. One major task 
for nodes on the chain is to validate transaction. it needs two do the following things:

1, The inputs of transaction are previously unspent. This is used to prevent double spending, you can use the money you already spent.

2, The sum of the inpupts is greater than or equal to the sum of outputs. This is used to make sure account balance, if you account has 100 dollars, then you can used at most 100, not any more

3, The scriptSig successfully unlocks the ScriptPubKey, that is the script that are created by combining the two of them can evaluate successfully. This is used to validate the receiver of is the legitimate
one (The signature can be verified).

Let's see how to satisfy the first point. In bitcoin all transaction that contains money that are not be spent are called UTXOs (unspent transaction outputs), As we have seen in previous section, the money 
that can be spent in this transaction is from the output of previous transaction. Then we need to check the money from the outputs of previous transaction is still valid and the current transaction can spend
those money.

A transaction is just like you spend some amount of money to buy some products or services, after completing the transaction, the amont of  money left in your pocket should less than the amount before the 
transaction, therefore when we validate a transaction, we need to make sure the amount in the output of a transaction is less than the total amount of inputs, if the amount of output is bigger than the input
which means the transaction is trying to fake some money out from the air, then we should sound the alarm.

Let's see how we can use code to check the amount after the transaction, in transaction.go add the following code:
```g

```
