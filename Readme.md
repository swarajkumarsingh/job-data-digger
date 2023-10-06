# Job Data Digger

![GitHub License](https://img.shields.io/github/license/swarajkumarsingh/job-data-digger)
![GitHub Stars](https://img.shields.io/github/stars/swarajkumarsingh/job-data-digger)
![GitHub Issues](https://img.shields.io/github/issues/swarajkumarsingh/job-data-digger)
![GitHub Forks](https://img.shields.io/github/forks/swarajkumarsingh/job-data-digger)

Job Data Scraper is a Go-based web scraper application that collects job listings from the Google Careers page and stores them in Redis. It is designed to run in a Docker container and automatically refreshes the data every 24 hours by re-scraping the Google Careers page.

## Table of Contents

- [Features](#features)
- [Prerequisites](#prerequisites)
- [Getting Started](#getting-started)
  - [Installation](#installation)
  - [Configuration](#configuration)
- [Usage](#usage)
- [Docker](#docker)
- [Contributing](#contributing)
- [License](#license)

## Features

- **Google Careers Page Scraping**: Automatically scrapes job listings from the Google Careers page.
- **Data Storage**: Stores scraped job data in a Redis database.
- **Scheduled Refresh**: Refreshes job data by re-scraping the Google Careers page every 24 hours.
- **Dockerized**: Easily deploy the application using Docker.

## Realtime images
![image](https://github.com/gin-gonic/gin/assets/89764448/f7aa1778-56e1-4f79-8b6a-19b58edb9341)

## Prerequisites

Before you begin, ensure you have met the following requirements:

- Docker installed on your system

**NOTE: if you don't have docker installed then install the following programs**

- Go (Golang) installed on your system.
- Redis server up and running.

## Getting Started

### Installation

To get started with the Job Data Scraper, follow these steps:

1. Clone the repository:

   ```bash
   git clone https://github.com/swarajkumarsingh/job-data-digger.git
   cd job-data-digger
   ```

2. Run the Go application:
    ```
    make compose
    ```

### Configuration
No configuration needed for this project

### Usage
1. For development(make sure a redis container is running)
```bash
./dev.sh
```

2. For running in production
```bash
./run.sh
```

### Docker
You can also run the Job Data Scraper in a Docker container. To build the Docker image, use the following command:
```bash
docker build -t job-data-scraper . && docker run -p 8080:8080 job-data-scraper
```

or

```bash
make run
```

### Contributing
Contributions are welcome! If you'd like to contribute to this project, please follow these guidelines:

1. Fork the repository.
2. Create a new branch for your feature or bug fix.
3. Make your changes and test thoroughly.
4. Commit your changes with clear commit messages.
5. Create a pull request against the main branch.

### License
This project is licensed under the MIT License. See the LICENSE file for details.