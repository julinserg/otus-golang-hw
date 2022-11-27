# file: features/calendar-notify.feature

# http://localhost:8080/
# http://calendar_service:8080/

Feature: Отправка уведомлений о наступления события в календаре

	Scenario: Доступность сервиса календаря
		When I send "GET" request to "http://calendar_service:8080/"
		Then The response code should be 200
		And The response should match text "This is my calendar!"

	Scenario: Добавление "события 100" требующего уведомления
		When I send "POST" request to "http://calendar_service:8080/add" with "application/json" data:
		"""
		{
			"id": "100",
			"title": "event100",
			"description": "testDescription",
			"time_start": "2022-11-22T01:10:30.00Z",
			"time_notify": 5000,
			"user_id": "12345"
		}
		"""
		Then The response code should be 200
		And I receive event with json: 
		"""
		{
			"id": "100",			
			"user_id": "12345"			
		}
		"""
	