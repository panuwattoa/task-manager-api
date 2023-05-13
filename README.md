# Project Name: Task Manager

Task Manager is a api backend that helps users manage their tasks. It allows users to post tasks, update their status, archive tasks, and add comments to tasks. The database used for this project is MongoDB, and it contains three collections: "comments," "profiles," and "task."

## Requirements

In order to run the application, you will need to have Docker and Docker Compose installed on your system.

## Installation

To set up the application, clone the repository and navigate to the project directory in your terminal. Then run the following commands:

```
docker-compose build
docker-compose up -d
```

This will build the necessary Docker images and start the containers. The initialization script `init-mongo.sh` will run the necessary database migrations.

## API Documentation

For API documentation and usage examples, please refer to the following links:

- [Task Manager API Postman Workspace](https://www.postman.com/dgtess/workspace/task-manager-api/documentation/13230697-6a3a2aee-2fe0-47f1-be9f-cdf4ccac18af)
- [Task Manager API Documentation](https://documenter.getpostman.com/view/13230697/2s93ecwqUS#d82ededf-dba3-43f8-9ca5-d2ea5a8319d2)

These resources will help you get started with using the API and integrating it with your own applications.

Thank you for using Task Manager! If you have any questions or feedback, please don't hesitate to contact us.
