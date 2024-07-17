# dupsMerger

It will search file with same hash and remove the file to filter duplicate files.
For eg: 
checker.com_checker.com_Overall_subdomain.txt and not_checker.com_checker.com_Overall_subdomain.txt if both file has same content then their hashes will be same. This will remove one file to avoid duplication.

Separating Prefix underscore match to megre duplicates.
For eg: checker.com_checker.com_Overall_subdomain.txt and checker.com_Overall_subdomain.txt files are present in a directory. This will filter and treat this two file as a same entity and merge the content in one file and remove the duplicate lines.

Usage: 
``` go run main.go /path/to/directory```
