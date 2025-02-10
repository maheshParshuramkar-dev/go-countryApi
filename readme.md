// README.md
# Go GET COUNTRY DETAILS

### Introduction
Project Go GET COUNTRY DETAILS is a Golang program to get details of country from external api and store it into cache.

### Project Go Get Country Details Features
* User can get Detais of country by giving it's country name.

### Prerequisites
Before running the application, ensure you have the following installed
* Go: Version 1.20 or later

### Installation Guide
* Clone the repository.
* Run go mod tidy to install all dependencies.

### Usage
* Run go run . or go run main.go to start the application.
* Connect to the API using Postman on port 3026.

### API Endpoints
| HTTP Verbs | Endpoints               | Action                                      |
|------------|-------------------------|---------------------------------------------|
| GET        | /api/countries/ping               | To check the status of Server               |
| POST       | /api/countries/search?name={countryName}         | To get the country details               |


### Technologies Used
* [Golang](https://go.dev/) This is an open-source programming language supported by Google, Easy to learn and great for teams ,Built-in concurrency and a robust standard library ,Large ecosystem of partners, communities, and tools.
* An in-memory data structure store used for caching.

### Authors
* Mahesh Parshuramkar
* [Github](https://github.com/maheshParshuramkar-dev)