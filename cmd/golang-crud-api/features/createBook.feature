Feature: getProductExportV2

  Background: Clean database
    Given SQL command
    """
    DELETE FROM myschema.books;
    """

  Scenario: Create a new book successfully
    When I setup a mock server for "GET" "https://api.mybiz.com/isbn/0-061-96436-1" with response 200 and body
    """json
    {
       "id": "0-061-96436-1"
    }
    """
    When API "POST" request is sent to "/api/v1/createBook" with payload
    """json
    {
      "isbn": "0-061-96436-1",
      "title": "The Art of Computer Programming"
    }
    """
    Then response status code is 200
    Then response body is
    """json
 {
    "isbn": "0-061-96436-1",
    "title": "The Art of Computer Programming",
    "total_pages": 0
}

    """