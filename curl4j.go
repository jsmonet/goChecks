package main

import (
  "fmt"
  "strings"
  "flag"
  "os"
  // "io/ioutil"
  // "net/http"

)

// this curler is currently purpose-built to hit neo4j hosts and figure out
// if the host is master/slave, or up at all.
// I'll expand it's utility just a bit later.
// Also, this is ONLY functional with HA neo4j, not causal clustering or
// singletons at the moment. I'll fix that later. Right now I'm all POC all the time
// like, you know, POC to prod without any changes ftw

func main() {
  rawHostAddress := flag.String("host", "", "enter a host address or IP") // using fqdn relies on faithful DNS resolution
  rawNeo4jRole := flag.String("role", "master", "enter master or slave") // Neo4j HA cluster has 3 roles: master, slave, arbiter. We won't test for the latter

  flag.Parse() // parse those flags!

  hostAddress := strings.ToLower(*rawHostAddress)
  neo4jRole := strings.ToLower(*rawNeo4jRole)
  fmt.Println("debug: printing hostAddress", hostAddress)
  fmt.Println("debug: printing neo4jRole", neo4jRole)

  // make sure neo4jRole is master or slave
  if neo4jRole != "master" && neo4jRole != "slave" {
    fmt.Println("Please ONLY use master or slave with the -role flag")
    os.Exit(1)
  }
  // rawHostAddress/hostAddress cannot be empty
  if hostAddress == "" {
    fmt.Println("Enter an address for your curl target")
    os.Exit(1)
  }

  curlTarget := fmt.Sprintf("http://%v:7474/db/manage/server/ha/%v", hostAddress, neo4jRole)

  fmt.Println(curlTarget)

}
