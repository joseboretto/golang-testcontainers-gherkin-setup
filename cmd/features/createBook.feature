Feature: Create book

  Background: Clean database
    Given SQL command
    """
    DELETE FROM myschema.books;
    """
    And reset mock server

  Scenario: Create a new book successfully
    Given a mock server request with method: "GET" and url: "https://api.isbncheck.com/isbn/0-061-96436-1"
    And a mock server response with status 200 and body
    """json
    {
       "id": "0-061-96436-1"
    }
    """
    And a mock server request with method: "POST" and url: "https://api.gmail.com/send-email" and body
    """json
    {
      "email" : "helloworld@gmail.com",
      "book" : {
        "isbn" : "0-061-96436-1",
        "title" : "The Art of Computer Programming"
      }
    }
    """
    And a mock server response with status 200 and body
    """json
    {
       "status": "OK"
    }
    """
    When API "POST" request is sent to "/api/v1/createBook" with payload
    """json
    {
      "isbn": "0-061-96436-1",
      "title": "The Art of Computer Programming"
    }
    """
    Then API response status code is 200 and payload is
    """json
    {
        "isbn": "0-061-96436-1",
        "title": "The Art of Computer Programming"
    }
    """