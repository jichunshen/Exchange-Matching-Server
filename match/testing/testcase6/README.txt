test match order when one seller and multiple buyer
multiple buyer are(ordered by time sequence):
account   share    price   symbol     order_id
1         10       5       Protoss    1
1         10       10      Protoss    2 
1         10       15      Protoss    3
2         20       10      Protoss    4
2         20       15      Protoss    5
2         20       20      Protoss    6

one seller is :
account   share    price   symbol    order_id
3         60       10      Protoss   7


the result will be:
order 6 will match for highest price
order 3,5 will match
order 2 will match for it was generated earlier than 4
