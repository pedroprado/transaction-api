Feature: Create Transaction


    Background:
        * url apiURL

        # get intermediary account
        Given url apiURL
        And path "/account/12345"
        When method GET
        Then assert responseStatus == 200

        * def intermediaryAccount = response
        * def intermediaryAccountBalanceBefore = intermediaryAccount.balance

        # create origin account
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

        * def originAccount = response
        * def originAccountBalanceBefore = response.balance

        # create destination account
        * def account =
        """
        {
            "bank": "22",
            "number": "22",
            "agency": "2",
            "balance": 100
        }    
        """

        Given url apiURL
        And path "/account"
        And request account
        When method POST
        Then assert responseStatus == 201

        * def destinationAccount = response
        * def destinationAccountBalanceBefore = response.balance
    
    @CreateTransaction
    Scenario: should create transaction sucessfully

        # BDD
        # Given an existing origin account
        # And an exiting destination account
        # When a transaction is created with value X
        # Then orign account must have X of its valance blocked (provisioneg)
        # And transaction must be in status OPEN

        # create origin account
        * def transactionAmount = 10
        * def transaction =
        """
        {
            "transaction_type": "PIX_OUT",
            "origin_account_id": "#(originAccount.account_id)",
            "destination_account_id": "#(destinationAccount.account_id)",
            "value": 10
        }    
        """

        Given url apiURL
        And path "/transactions"
        And request transaction
        When method POST
        Then assert responseStatus == 201

        * def transactionID = response.transaction_id

        Given url apiURL
        And path "/transaction_status"
        And param transaction_id = transactionID
        When method GET
        Then assert responseStatus == 200
        And assert response.status == "OPEN"

        Given url apiURL
        And path "/balance_provisions"
        And param transaction_id = transactionID
        When method GET
        Then assert responseStatus == 200
        And assert response.length == 1
        And assert response[0].type == "ADD"
        And assert response[0].status == "OPEN"

        # check amount was subtracted (blocked) from origin account
        Given url apiURL
        And path "/account/" + originAccount.account_id
        When method GET
        Then assert responseStatus == 200
        And assert response.balance ==  (originAccountBalanceBefore - transactionAmount)

        #  chec amount was add to intermediary account
        Given url apiURL
        And path "/account/12345"
        When method GET
        Then assert responseStatus == 200
        And assert response.balance ==  (intermediaryAccountBalanceBefore + transactionAmount)