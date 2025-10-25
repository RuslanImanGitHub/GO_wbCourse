SERVER_PATH = .\backend\cmd\main.go
SERVICE_PATH = .\services\kafkaService.go
ZOOKEPER_START = C:\Kafka\kafka_2.13-3.9.1\bin\windows\zookeeper-server-start.bat .\config\zookeeper.properties. 
KAFKA_START = C:\Kafka\kafka_2.13-3.9.1\bin\windows\kafka-server-start.bat .\config\server.properties 

ZOOKEPER_STOP = C:\Kafka\kafka_2.13-3.9.1\bin\windows\zookeeper-server-stop.bat
KAFKA_STOP = C:\Kafka\kafka_2.13-3.9.1\bin\windows\kafka-server-stop.bat

start:
	cmd /C "${ZOOKEPER_START}"
	@echo "Запустился zookeeper"
	cmd /C "${KAFKA_START}"
	@echo "Запустилась kafka"

	go run $(SERVER_PATH)
	@echo "Запустился сервер"
	go run $(SERVICE_PATH)
	@echo "Запустился сервис"

stop:
	cmd /C "${ZOOKEPER_STOP}"
	@echo "Запустился zookeeper"
	cmd /C "${KAFKA_STOP}"
	@echo "Запустилась kafka"