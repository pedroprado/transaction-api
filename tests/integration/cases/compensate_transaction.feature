Feature: Compensate transaction

    Background:
        * url apiURL

    
    Scenario: should compensate a failed transacti'on

        # BDD
        # Given a transaction FAILED'
        # A an void balance provision OPEN
        # When it is requested to compensate
        # And the provision should be closed
        # And the balance should be added to the origin account
        # And the balance should be subtracted from the intermediary account

        #  starting from cenario where we create a failed transaction 
        * call read('classpath:cases/complete_transaction.feature@GenerateCompensation')

        # add funds from intermediary account
        * def patchRequest = 
        """
        {
            "account_id": "12345",
            "balance": 1000
        }    
        """

        Given url apiURL
        And path "/account"
        And request patchRequest
        When method PATCH
        Then assert responseStatus == 204

        Given url apiURL
        And path "/transaction/" + transactionID + "/compensate"
        And request {}
        When method POST
        Then assert responseStatus == 202


        Given url apiURL
        And path "/transaction_status"
        And param transaction_id = transactionID
        When method GET
        Then assert responseStatus == 200
        And assert response.status == "FAILED"

        Given url apiURL
        And path "/balance_provisions"
        And param transaction_id = transactionID
        When method GET
        Then assert responseStatus == 200
        And assert response.length == 2
        And assert response[0].type == "ADD"
        And assert response[0].status == "CLOSED"
        And assert response[1].type == "VOID"
        And assert response[1].status == "CLOSED"


        # check amount has returned to origin account
        Given url apiURL
        And path "/account/" + originAccount.account_id
        When method GET
        Then assert responseStatus == 200
        And assert response.balance ==  originAccountBalanceBefore