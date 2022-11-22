# file: features/calendar-find.feature

# http://localhost:8080/
# http://calendar_service:8080/

Feature: Поиск событий в календаре по дате

	Scenario: Доступность сервиса календаря
		When I send "GET" request to "http://localhost:8080/"
		Then The response code should be 200
		And The response should match text "This is my calendar!"

	Scenario: Добавление корректного "события 10" в календарь
		When I send "POST" request to "http://localhost:8080/add" with "application/json" data:
		"""
		{
			"id": "10",
			"title": "event10",
			"description": "testDescription",
			"time_start": "2022-01-01T01:10:30.00Z"
		}
		"""
		Then The response code should be 200
		And The response should match text ""

	Scenario: Добавление корректного "события 20" в календарь
		When I send "POST" request to "http://localhost:8080/add" with "application/json" data:
		"""
		{
			"id": "20",
			"title": "event20",
			"description": "testDescription",
			"time_start": "2022-01-03T01:10:30.00Z"
		}
		"""
		Then The response code should be 200
		And The response should match text ""

	Scenario: Добавление корректного "события 30" в календарь
		When I send "POST" request to "http://localhost:8080/add" with "application/json" data:
		"""
		{
			"id": "30",
			"title": "event30",
			"description": "testDescription",
			"time_start": "2022-01-14T01:10:30.00Z"
		}
		"""
		Then The response code should be 200
		And The response should match text ""
	
	Scenario: Получение списка событий за день - 2022-01-01
		When I send "GET" request to "http://localhost:8080/get_by_day?time=2022-01-01T00:00:00Z"
		Then The response code should be 200
		And The response should match json:
		"""
		{
			"data": [
				{
					"id": "10"					
				}
			]					
		}
		"""

	Scenario: Получение списка событий за день - 2023-01-01 (таких событий нет)
		When I send "GET" request to "http://localhost:8080/get_by_day?time=2023-01-01T00:00:00Z"
		Then The response code should be 200
		And The response should match json:
		"""
		{
			"data": []					
		}
		"""
	
	Scenario: Получение списка событий за неделю начиная с 2022-01-01
		When I send "GET" request to "http://localhost:8080/get_by_week?time=2022-01-01T00:00:00Z"
		Then The response code should be 200
		And The response should match json:
		"""
		{
			"data": [
				{
					"id": "10"					
				},
				{
					"id": "20"					
				}
			]					
		}
		"""

	Scenario: Получение списка событий за месяц начиная с 2022-01-01
		When I send "GET" request to "http://localhost:8080/get_by_month?time=2022-01-01T00:00:00Z"
		Then The response code should be 200
		And The response should match json:
		"""
		{
			"data": [
				{
					"id": "10"					
				},
				{
					"id": "20"					
				},
				{
					"id": "30"					
				}
			]					
		}
		"""