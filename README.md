# Go Telemetry CEP Temperature

This project consists of an integrated solution of two services for accessing detailed weather information, using Postal Addressing Codes (CEPs) as a query parameter. Service A allows the insertion of the CEP through a POST request, with the body `{"cep":"01001000"}`, while Service B provides the weather data through a GET request, accessible by the URL `/?cep=01001000`. Upon receiving a request, Service A processes and redirects the query to Service B to obtain the desired weather information.

This project is an integrated solution composed of two services, Service A and Service B, which offers access to detailed weather information, using Postal Addressing Codes (CEPs) as a parameter for consultation. There is the implementation of OpenTelemetry (OTEL) in conjunction with Zipkin for distributed tracing, allowing the visualization of the journey of a request between Service A and Service B. Service A accepts a CEP via a POST request, with the body `{"cep":"01001000"}`, and consults Service B, which provides the weather data through a GET request at the URL `/?cep=01001000`. The integration of OTEL with Zipkin facilitates monitoring and analysis of the response time both for CEP search and for weather information search.

## Features

- **Two Integrated Services**: The project consists of two distinct services, Service A, which receives the CEP via POST, and Service B, which provides weather information via GET, facilitating access to accurate data based on the provided CEP.
- **Direct CEP Query in Service B**: Service B allows direct access to specific weather information of a location, using the CEP as a query key.
- **Rigorous CEP Validation in Service A**: Service A implements rigorous validation of the entered CEP, ensuring it is in the correct format and consists only of numbers, with exactly 8 characters, before redirecting the query to Service B.
- **Free Authentication**: Both services have been designed to be accessible without the need for authentication, simplifying the process of querying weather information.
- **Responses in JSON Format**: Weather information is provided in JSON format by Service B, facilitating integration with other applications and the manipulation of the received data.
- **Support for Multiple Temperature Units**: Service B offers temperature information in Celsius, Fahrenheit, and Kelvin, catering to the diverse preferences and needs of users.
- **Integration with OTEL + Zipkin for Distributed Tracing**: The implementation of OpenTelemetry along with Zipkin provides effective distributed tracing between Service A and Service B.
- **Spans to Measure Response Times**: Specific spans are created to measure the response time of CEP search operations in Service A and weather information search in Service B.

## Usage Example

To consult weather information through the command line, you can use `curl`, a powerful tool available on most operating systems for making HTTP requests. Below are practical examples of how to use curl to obtain the temperature based on a specific CEP.

### Performing a Query

To make a query, simply replace CEP with the desired postal code in the URL. Here are some examples:

```bash
curl -X POST http://localhost:3000/ -H "Content-Type: application/json" -d '{"cep":"01001000"}'
```

Expected return:

```json
{"city":"São Paulo","temp_C":22.4,"temp_F":72.32,"temp_K":295.55}
```

In this example, the request returns the temperature for the CEP 01001000 (a São Paulo CEP), showing the temperature in Celsius (temp_C), Fahrenheit (temp_F), and Kelvin (temp_K) and the city (city).

## How Data is Returned

Data is returned in JSON format. Each field in the JSON represents a different temperature measure:

- `city`: Name of the city.
- `temp_C`: Temperature in degrees Celsius.
- `temp_F`: Temperature in degrees Fahrenheit.
- `temp_K`: Temperature in Kelvin.

## Development

In the development of this project, I focused on creating a solution composed of two interconnected services that use external APIs to provide accurate weather information, based on a Postal Addressing Code (CEP) provided. Service A is responsible for receiving the CEP through a POST request and then communicating with Service B, which performs the queries to the external APIs and returns the weather data. Below, I describe the steps involved and how each service and API are employed, including the implementation of OpenTelemetry (OTEL) and Zipkin for distributed tracing.

### Address Search by CEP with viacep.com.br (Service B)

The journey begins when Service B collects detailed address information using the CEP provided by Service A. For this, it consults the ViaCEP API, which returns data such as street, neighborhood, city, and state. These details are crucial for identifying the precise geographical location for subsequent weather queries.

### Longitude and Latitude Search with nominatim.openstreetmap.org (Service B)

With the address data in hand, Service B then converts this information into geographical coordinates (latitude and longitude) through the Nominatim API, part of the OpenStreetMap project. This conversion is essential to ensure the accuracy of the weather queries that depend on geographical coordinates.

### Temperature Search (Service B)

With the geographical coordinates available, Service B performs the query for the current weather conditions. Depending on the availability of data, it can use:

- The Open-Meteo API, for detailed weather queries based on coordinates, providing accurate temperature information for the specified location.
- The wttr.in API, for weather information based on location names, which although may not be as precise as the query by coordinates, still provides a valid estimate of the weather conditions.

### Integration with OTEL + Zipkin

The integration with OpenTelemetry (OTEL) and Zipkin adds a layer of observability to the project, allowing for distributed tracing between Service A and Service B. This functionality enables the monitoring of the complete journey of a request, including measuring the response time for CEP search and temperature search, facilitating the identification and resolution of possible bottlenecks or performance issues.

## Error Handling

I implemented error handling at each stage to ensure that the system can appropriately handle scenarios such as invalid CEPs, failures in obtaining coordinates, or errors in API responses.

## Unit Tests

A part of the development of this project involves the implementation of comprehensive unit tests, ensuring the reliability and robustness of each functionality offered by the application. The approach adopted for the tests follows best software development practices, focusing on validating each component in isolation to ensure its correct operation in various scenarios.

### Test Coverage

The unit tests cover a wide range of use cases and error scenarios, including, but not limited to:

- CEP Validation: Tests to ensure that only valid CEPs in the correct format are accepted, and that appropriate error messages are returned for invalid or improperly formatted CEPs.
- External API Queries: Tests to verify the correct interaction with the external APIs used to obtain address information, geographical coordinates, and weather data. This includes simulating API responses to test proper handling of data and errors.
- Temperature Unit Conversion: Tests that validate the accuracy of temperature conversions between Celsius, Fahrenheit, and Kelvin, ensuring that the calculations are correct.
- Error Handling: Specific tests to verify the system's robustness in facing errors during information querying, including network failures, errors in external APIs, and unexpected data.

## Makefile

This project includes a Makefile designed to offer an efficient and simplified interface for managing development and production environments, as well as executing automated tests. The commands provided allow optimizing and streamlining the development workflow, testing, and project maintenance, ensuring a more effective and organized management.

### Development Commands

### `make dev-start`

Starts the services defined in the `docker-compose.dev.yml` file for the development environment in detached mode (in the background). This allows the services to run in the background without occupying the terminal.

### `make dev-stop`

Stops the services that are running in the background for the development environment. This does not remove the containers, networks, or volumes created by `docker compose up`.

### `make dev-down`

Shuts down the development environment services and removes the containers, networks, and volumes associated created by `docker compose up`. Use this command to clean up resources after development.

### `dev-run-service-a`

Starts the execution of Service A within the development environment, using Docker Compose to execute the `go run` command in the `/cmd/input_server/main.go` file. It is ideal for quickly starting the project server in development mode.

### `dev-run-service-b`

Starts the execution of Service B within the development environment, using Docker Compose to execute the `go run` command in the `/cmd/temperature_server/main.go` file. It is ideal for quickly starting the project server in development mode.

### `make dev-run-tests`

Executes all Go tests within the development environment, showing verbose details of each test. This command is useful for running the project's test suite and checking if everything is functioning as expected.

### Production Commands

### `make prod-start`

Starts the services defined in the `docker-compose.prod.yml` file for the production environment in detached mode. This is useful for running the project in an environment that simulates production.

### `make prod-stop`

Stops the production environment services that are running in the background, without removing the associated containers, networks, or volumes.

### `make prod-down`

Shuts down the production environment services and removes the associated containers, networks, and volumes, cleaning up resources after use in production.

## Prerequisites

Before starting, make sure you have Docker and Docker Compose installed on your machine. If not, you can download and install from the following links:

- Docker: https://docs.docker.com/get-docker/

### Clone the Repository

First, clone the project repository to your local machine. Open a terminal and execute the command:

```bash
git clone https://github.com/aronkst/go-telemetry-cep-temperature.git
```

### Navigate to the Project Directory

After cloning the repository, navigate to the project directory using the cd command:

```bash
cd go-telemetry-cep-temperature
```

## Development Environment

### Build the Project with Docker Compose

In the project directory, execute the following command to build and start the project using Docker Compose:

```bash
docker compose -f docker-compose.dev.yml up --build
```

Or using the Makefile:

```bash
make dev-start
```

This command will build the Docker image of the project and start the container.

### Run the Project with Docker Compose

To start the main service of your project in development mode, you can use the direct commands from Docker Compose:

```bash
docker compose -f docker-compose.dev.yml exec dev go run cmd/input_server/main.go
```

```bash
docker compose -f docker-compose.dev.yml exec dev go run cmd/temperature_server/main.go
```

Or using the Makefile:

```bash
make dev-run-service-a
```

```bash
make dev-run-service-b
```

### Access the Project

With the container running, you can access the project through the browser or using tools like curl, pointing to http://localhost:3000/, replacing CEP with the desired postal code.

### curl Command Example

To test if the project is running correctly, you can use the following curl command in a new terminal:

```bash
curl -X POST http://localhost:3000/ -H "Content-Type: application/json" -d '{"cep":"01001000"}'
```

You should receive a JSON response with temperatures in Celsius, Fahrenheit, Kelvin, and the city.

### Visualizing Telemetry

Open a browser and access http://localhost:9411/zipkin/. This URL will take you to the Zipkin user interface, where you can start visualizing the telemetry of your services.

1. **Search for Traces**: In the Zipkin interface, you can search for traces in various ways, such as by service, operation name, annotations, and tags.

2. **Analyze Traces**: After finding a specific trace, you can click on it to see the details. This includes information like the duration of the trace.

### Ending the Project

To end the project and stop the Docker container, go back to the terminal where Docker Compose is running and press Ctrl+C. To remove the containers created by Docker Compose, execute:

```bash
docker compose -f docker-compose.dev.yml down
```

Or using the Makefile:

```bash
make dev-down
```

## Production Environment

### Build and Run the Project with Docker Compose

In the project directory, execute the following command to build and start the project in the production environment using Docker Compose:

```bash
docker compose -f docker-compose.prod.yml up --build
```

Or using the Makefile:

```bash
make prod-start
```

This command will build the Docker image of the project for production and start the containers.

## curl Command Example

To check if the production project is operational, use the following curl command, adjusting the address according to your configuration:

```bash
curl -X POST http://localhost:3000/ -H "Content-Type: application/json" -d '{"cep":"01001000"}'
```

You should receive a JSON response with the requested information, such as temperatures in Celsius, Fahrenheit, Kelvin, and the city.

### Visualizing Telemetry

Open a browser and access http://localhost:9411/zipkin/. This URL will take you to the Zipkin user interface, where you can start visualizing the telemetry of your services.

1. **Search for Traces**: In the Zipkin interface, you can search for traces in various ways, such as by service, operation name, annotations, and tags.

2. **Analyze Traces**: After finding a specific trace, you can click on it to see the details. This includes information like the duration of the trace.

### Ending the Project

To end the project and stop the production containers, use the following command:

```bash
docker compose -f docker-compose.prod.yml down
```

Or using the Makefile:

```bash
make prod-down
```

This command ends all production services and removes the associated containers, networks, and volumes, ensuring that the production environment is cleaned up after use.
