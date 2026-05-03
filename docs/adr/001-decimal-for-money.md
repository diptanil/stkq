# Never use `float64` for money

We know that 0.1 + 0.2 ≠ 0.3 in floating point. 
You may have shrugged and moved on. 
**Don't shrug for money.** 
Real-world losses from this:

A 2010 [Pittsburgh election](https://en.wikipedia.org/wiki/Patriot_missile)
was decided wrong because a clock counter accumulated float drift over 100
hours. The Patriot missile system did the same thing in 1991 and missed an
incoming Scud.

Banks and exchanges *never* store currency as float. Stripe, Square, every
US brokerage — all integer cents or arbitrary-precision decimals.

We use [`shopspring/decimal`](https://github.com/shopspring/decimal),
a popular Go arbitrary-precision decimal type.