# Distributed Union-Find
Distributed Systems – TD9 : The Contest

These are the subjects of the M1: Mid-term exam - Distributed System by Eddy Caron.

The goal is the same for each question, that is to optimize your algorithms both in number of messages and message size (a way to have the knowledge of all the platform in a data structure is irrelevant).

# Problems

## From same ID to unique ID

### Everybody is different in the tree

Considering a distributed system with a tree structure where each node has the same ID (the identity of the node like a number). Write two algorithms to ensure that each node has a unique ID.
1. with one initiator
2. with n initiators (1 < n ≤ N where N is the number of nodes)

### Everybody is different anywhere

Considering a distributed system with a general graph structure where each node has the same ID. Write two algorithms to ensure that each node has a unique ID.
1. with one initiator
2. with n initiators (1 < n ≤ N where N is the number of nodes)

## From random ID to unique ID

Consider a distributed system with a general graph structure where each node has a random ID (probably a few nodes can have the same ID, but not all). Write an algorithm to ensure that each node has a unique ID even with many initiators.

# Solution

A Distributed implementation of Disjoint Sets Union a.k.a Union-Find.

**TODO**: 
   1. More detailed description of the algorithm
   2. Grab component: master (Done), slave
   3. Try to grab Neighbors: master (Done), slave
   4. Distribute IDs to slaves: master, slave
   4. Unite Edge masters (Done)

## Install

``` bash
go build
```

## Setting up

Before running the App replica, is necessary to config some attributes in the ```appconfig.json``` file.
``` bash
{
  "MyIP": "localhost",
  "Port": 9002,
  "Initiator": true,
  "Neighborhood": ["http://localhost:9000", "http://localhost:9001"]
}
```

## Run

Let's play.
```bash
./distributed-union-find
```
