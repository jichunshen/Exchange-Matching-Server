test match order when one buyer and multiple seller
multiple seller are(ordered by time sequence):
account   share    price   symbol     order_id
2         10       5       Protoss    1
2         10       10      Protoss    2
2         10       15      Protoss    3
3         20       10      Protoss    4
3         20       15      Protoss    5
3         20       5       Protoss    6

one buyer is :
account   share    price   symbol    order_id
3         40       10      Protoss   7


the result will be:
order 1,6 will match for lowest price
order 2 will match for it was generated earlier than 4
