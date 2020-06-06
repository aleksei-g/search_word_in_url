# Word count in URLs

This script receives a list of URLs from stdin and counts the number of occurrences of the word "go" in the response body of each URL.

## Local installing
Create a directory and go to it. Clone project into this folder:

```
$ mkdir project_folder
$ cd project_folder
$ git clone https://github.com/aleksei-g/search_word_in_url.git .
```

## Run script
To run script, use the following command:
```
$ go run main.go
```
and pass URL to count the number of occurrences of the word "go":
```
https://golang.org
Count for https://golang.org: 9
^C
Total: 9
```
Or pass the list of URLs when running the script:
```
$ echo -e 'https://golang.org\nhttps://golang.org' | go run main.go
```
result:
```
Count for https://golang.org: 9
Count for https://golang.org: 9
Total: 18
```
