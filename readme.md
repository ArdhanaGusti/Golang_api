# Golang API

## Install Tutorial
1. Clone Repository
```bash
git clone https://github.com/ArdhanaGusti/Golang_api.git
```

2. Navigate to the directory
```bash
cd Golang_api
```

3. Install dependencies
```bash
go mod init github.com/ArdhanaGusti/Golang_api
go mod tidy
```
4. Copy .env.example and rename to .env

5. Fill the DB connetion mysql and client credentials (Optional) in .env

6. Run Project
```bash
go run server.go
```

## Reason Why Using MVC

I'm using MVC Design Pattern for this project. MVC design pattern has 3 main module such as model, view & controller. Model is for structure of table in database, controller is for main logic and view is for user interface. The main reason why i'm using this design pattern is easily to understand for next programmer that want to clone this project, and also this design pattern is implement the clean architecture. Clean architecture can make code easily to understand because it separated by it's function, make it neater.