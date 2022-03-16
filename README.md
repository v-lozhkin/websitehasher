# MD5 Website hasher 1.0.0
___
This CLI tool can get any internet resource and print its MD5 hash,


## Features:
___
Get resource, print the md5 hash

## Get started:
___
### build:
* go build -o myhttp main.go 

### run:
./myhttp [-parallel *max_number_of_parallelism*] [paths]
>parallel: limit of parallel requests (default 10)

>paths: 0..inf resources to parse

### example:
./myhttp -parallel 3 google.com yandex.ru http://mail.ru


## Author:
___
Lozhkin Vladimir. 2022



