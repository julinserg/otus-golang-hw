# file: features/notification.feature

# http://localhost:8088/
# http://reg_service:8088/

Feature: Email notification sending
	As API client of registration service
	In order to understand that the user was informed about registration
	I want to receive event from notifications queue

	Scenario: Registration service is available
		When I send "GET" request to "http://localhost:8080/"
		Then The response code should be 200
		And The response should match text "This is my calendar!"
