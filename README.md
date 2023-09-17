# hkgi-go

Rewrite of Ced's [hkgi](https://github.com/cedric-h/hkgi) in Go and Postgres.


The goal of this project is stability and persistence. One of the biggest issues
with the previous iteration of hkgi was how it stored user data in memory and in
a constantly-overwritten JSON file. This project aims to solve that issue by
using a persistent RDBMS (Postgres).
