# golang_tasks


Micro-service Tasks is the part of the end-to-end project of the task matching system (uses [auth mucro-service](https://github.com/yuliazhelt/golang_auth) for authorization). The project is built according to the hexagonal architecture rules.

System users create tasks and send them for approval. The system monitors the process of coordination of tasks and sends a notification to participants about the progress of coordination.

Following methods of the service are available to users via http:
- POST method to create tasks, returning the task id
- GET method to get task by id
- POST method for approving tasks by a person who agrees with tasks
- POST method to reject a task by a Coordinator
- GET method for getting list of tasks, created by user