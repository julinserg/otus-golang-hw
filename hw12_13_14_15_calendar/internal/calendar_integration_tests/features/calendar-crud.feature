# file: features/calendar-crud.feature

# http://localhost:8080/
# http://calendar_service:8080/

Feature: Добавление/удаление/изменение событий в календаре

	Scenario: Доступность сервиса календаря
		When I send "GET" request to "http://calendar_service:8080/"
		Then The response code should be 200
		And The response should match text "This is my calendar!"

	Scenario: Добавление корректного "события 1" в календарь
		When I send "POST" request to "http://calendar_service:8080/add" with "application/json" data:
		"""
		{
			"id": "1",
			"title": "event1",
			"description": "testDescription",
			"time_start": "2021-02-18T21:54:42.123Z"
		}
		"""
		Then The response code should be 200
		And The response should match text ""

	Scenario: Добавление корректного "события 2" в календарь
		When I send "POST" request to "http://calendar_service:8080/add" with "application/json" data:
		"""
		{
			"id": "2",
			"title": "event2",
			"description": "testDescription",
			"time_start": "2022-02-18T21:54:42.123Z"
		}
		"""
		Then The response code should be 200
		And The response should match text ""

	Scenario: Добавление НЕ корректного "события 3" в календарь(время начала события уже занято)
		When I send "POST" request to "http://calendar_service:8080/add" with "application/json" data:
		"""
		{
			"id": "3",
			"title": "event3",
			"description": "testDescription",
			"time_start": "2022-02-18T21:54:42.123Z"
		}
		"""
		Then The response code should be 500
		And The response should match error "time event is busy"
	
	Scenario: Добавление НЕ корректного "события 4" в календарь(невалидный json)
		When I send "POST" request to "http://calendar_service:8080/add" with "application/json" data:
		"""
		{
			"id": "4",
			"title": "test
		}
		"""
		Then The response code should be 400
		And The response should match error "unexpected end of JSON input"
	
	Scenario: Удаление "события 1"
		When I send "POST" request to "http://calendar_service:8080/remove" with "application/json" data:
		"""
		{
			"id": "1"
		}
		"""
		Then The response code should be 200
		And The response should match text ""
	
	Scenario: Повторное удаление "события 1" (ошибка)
		When I send "POST" request to "http://calendar_service:8080/remove" with "application/json" data:
		"""
		{
			"id": "1"
		}
		"""
		Then The response code should be 500
		And The response should match error "Event ID not exist"

	Scenario: Обновление "события 2"
		When I send "POST" request to "http://calendar_service:8080/update" with "application/json" data:
		"""
		{
			"id": "2", 
			"title": "testTitle", 
			"description": "testDescription", 
			"time_start": "2021-02-18T21:54:42.123Z"
		}
		"""
		Then The response code should be 200
		And The response should match text ""

	Scenario: Обновление несуществующего "события 1" (ошибка)
		When I send "POST" request to "http://calendar_service:8080/update" with "application/json" data:
		"""
		{
			"id": "1", 
			"title": "testTitle", 
			"description": "testDescription", 
			"time_start": "2021-02-18T21:54:42.123Z"
		}
		"""
		Then The response code should be 500
		And The response should match error "Event ID not exist"


