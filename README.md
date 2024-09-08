
#  🚂Train Route Planner

Welcome to the Train Route Planner project! This project helps you plan train routes between stations, find all possible paths, and determine the optimal route for a set number of trains. It reads a network map file of stations and connections and generates routes based on the specified start and end stations.
### The project performs the following key functions:

#### 1. Data Reading and Validation:
- Reads a file with the map of stations and connections between them.
- Validates the correctness of the data (duplicates, correct format, existence of stations and their connections).
#### 2. Route Search Algorithm:
- Generates all possible routes between two stations.
- Searches for all route combinations without intersections, allowing trains to move without conflicts.
#### 3. Optimization:
- Finds the optimal train distribution along routes based on the minimum number of moves.
- Calculates the number of trips required for each train to complete all routes.
#### 4. Displaying Results:
- Simulates train movement along routes with real-time updates on train status and locations.
- The program dynamically updates information about stations and trains, displaying all movements at each step.

## Table of Contents
1. [Project Overview](#1-project-overview)
2. [Features](#2-features)
3. [How It Works](#3-how-it-works)
4. [Installation and Setup](#4-installation-and-setup)
5. [Checking Total Turns](#5-checking-total-turns)
6. [Usage](#6-usage)
7. [Command Line Arguments](#7-command-line-arguments)
8. [Detailed Process Flow](#8-detailed-process-flow)
9. [Error Handling](#9-error-handling)
10. [Program Structure](#10-program-structure)
11. [Conclusion](#11-conclusion)


## 1. Project Overview
The Train Route Planner is a command-line tool that takes a network map file as input and outputs the best route plan for a given number of trains. It ensures that no stations are overused by multiple trains simultaneously and provides an optimal sequence of train movements.

## 2. Features
 - Station and Connection Parsing: Reads a network map and identifies stations and connections between them.
 - Route Generation: Finds all possible routes from a start station to an end station.
 - Route Optimization: Determines the best combination of routes for the given number of trains.
 - Simulated Train Movements: Displays how each train moves across the stations in the most efficient way.
 - Error Handling: Provides detailed error messages when invalid inputs are provided or no routes are found.

## 3. How It Works
 - Input: The program takes a network map file and user inputs for the start station, end station, and the number of trains.
 - Parsing: The network map file is parsed to get a list of stations and their connections.
 - Route Finding: The tool finds all valid routes from the start to the end station.
 - Optimization: It evaluates different route combinations to find the best route for the trains.
 - Output: Finally, the train movements are displayed, showing how trains are allocated to routes and move across stations.

## 4. Installation and Setup
Clone the repository:

```
git clone https://gitea.koodsisu.fi/irynazaporozhets/stations.git
cd stations
```
Install Go: If you don't have Go installed, follow this guide.

Run the program: Once the Go environment is set up, run the program as follows:

`go run . tests/<network-map-file> <start-station> <end-station> <number-of-trains>`


## 5. Checking Total Turns

To check the total number of turns used for the fastest route, You can use the following command:  
`go run . tests/<network-map-file> <start-station> <end-station> <number-of-trains> wc -l`

#### Example:

```
$ go run . tests/londonNetwork.map waterloo st_pancras 2| wc -l
8
```

### 6. Usage
Here is how you can use the tool via the command line:

```
go run . tests/londonNetwork.map waterloo st_pancras 2
```

Where:

- londonNetwork.map is the file containing the network of stations and connections.
- waterloo is the starting station.
- st_pancras is the destination station.
- 2 is the number of trains you want to plan for.

#### Example:
Given a map file like:

stations:
waterloo,3,1
victoria,6,7
euston,11,23
st_pancras,5,15

connections:
waterloo-victoria
waterloo-euston
st_pancras-euston
victoria-st_pancras

Running:

```
go run . tests/londonNetwork.map waterloo st_pancras 2
```

Would result in the tool finding routes between waterloo st_pancras for 2 trains:
```
 T1-victoria T2-euston 
 T1-st_pancras T2-st_pancras
```


## 7. Command Line Arguments
Required:
Network Map File: A .txt file containing station details and connections between them.
Start Station: The station where trains start.
End Station: The destination station for the trains.
Number of Trains: A positive integer specifying the number of trains to be planned.
Optional:
Extra arguments: Additional options such as "extra" or "bonus" (e.g., extra, bonus).

## 8. Detailed Process Flow
- Argument Validation: The program ensures there are enough command-line arguments and that the number of trains is valid.
 - Network Map Parsing:
The map file is parsed to extract a list of stations and their connections.
Validates that the file includes both the stations and connections sections.
Ensures there are no duplicate stations or routes.
 - Route Generation: All possible routes between the start and end station are found. Routes are sorted by length.
 - Route Optimization: The best route combination for the given number of trains is calculated to minimize turns.
 - Train Movement Simulation: Displays how each train moves across its assigned route.

## 9. Error Handling
### The program handles errors such as:

 - Too few or too many arguments: Prints an error message and exits.
 - Invalid train count: Ensures the number of trains is a valid positive integer.
 - Missing or invalid stations: Checks that the start and end stations exist in the network.
 - nvalid routes: Ensures that there is a valid route between the start and end stations.

 
## 10. Program Structure
#### main.go File
This is the main file of the program that manages the train movement simulation process. It:

- Reads command line parameters: path to the map file, start and end stations, and the number of trains.
- Verifies the correctness of the input data (e.g., ensuring the number of trains is positive).
- Calls functions to parse the map of stations and roads.
- Checks for the existence of the start and end stations on the map.
- Finds all possible routes and determines the optimal route for the specified number of trains.
- Starts the train movement simulation and displays the results on the screen.

#### train Package
The train package contains the main functions that implement the program's logic:

- ParseNetworkMap — loads the map of stations and roads from a file, checks for the presence of all necessary sections, and validates the data format.
- ValidateStationExistence — checks that the specified start and end stations exist on the map.
- BuildConnectionMap — builds a connection map between stations in the form of a dictionary, where each element corresponds to its neighbors.
- FindAllPossibleRoutes — finds all possible routes between two stations using recursive pathfinding.
- FindOptimalRoute — determines the most efficient route for the trains based on the number of trains and the length of the routes.
- DisplayTrainMovements — simulates and displays the train movements along the selected routes.
- Auxiliary functions — such as validation of data (number of trains, coordinate correctness, avoidance of duplicate routes and stations).

#### Sample Files:
- main.go: Responsible for organizing the main execution process of the program.
- train.go: Implements the core algorithms and logic for working with the network map, route searching, and train movement simulation.

## 11. Conclusion
`This project enables the solution of the task of efficiently organizing train movement within a complex network of stations and routes. During development, methods for working with map data, searching for optimal routes, and simulating train movements were implemented. The program is quite flexible and can be used to optimize any transportation systems with fixed routes and a limited number of trains.`

`This project demonstrates the importance of data processing, validating input data, and applying efficient algorithms to solve planning and optimization tasks in real-world conditions.`
 