# Autify CLI 

## Overview
Autify CLI helps download web pages and saves them to disk for future retrieval.
### Demo



https://user-images.githubusercontent.com/52675260/196665785-2400d735-70b7-41dd-b3bc-70acf0c88017.mp4


## How to use

#### create binary executable file 

```
 go build -o  fetch
```

#### Download web page(s)

```
$ ./fetch https://www.google.com

Downloading: https://www.google.com
Download successful: 1 page(s) downloaded
Execution time: 0.910034 seconds%     

```

View downloaded file
```
$ ls
Dockerfile              README.md               fetch                   go.sum                  meta.json               pkg
LICENSE                 cmd                     go.mod                  main.go                 models                  ***www.google.com.html***

```

You can also download more than one web page
```
$ ./fetch https://www.google.com  https://www.mongodb.com

Downloading: https://www.mongodb.com
Download successful: 2 page(s) downloaded
Execution time: 3.817177 seconds%      

```

#### View page metadata

```
$ ./fetch --metadata  https://www.google.com  

site: https://www.google.com
num_links: 22
images: 1
last_fetch: Wed Oct 10 2022 10:49:24 UTC
Execution time: 0.000207 seconds%  
```

### TODO
Write Unit and Integration test
