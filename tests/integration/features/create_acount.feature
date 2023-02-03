Feature: Create account

    Background:
        * url apiURL




    Scenario: should create account

        * def account =
        """
        {
            "bank": "11",
            "number": "11",
            "agency": "1",
            "balance": 100
        }    
        """

        Given url apiURL
        And path "/account"
        And request account
        When method POST
        Then assert responseStatus == 201
        
        * def created = response

        Given url apiURL
        And path "/account/" + created.account_id
        When method GET
        Then assert responseStatus == 200
        And match response == created
        