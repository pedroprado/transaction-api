Feature: Complete Transaction


    Background:
        * url apiURL

    Scenario: should complete transaction

        # BDD
        # Given a transaction OPEN
        # When it is requested to complete
        # Then the transaction goes to BOOKED status
        # And the provision should be closed
        # And the balance should be added to the destination account
        # And the balance should be subtracted from the intermediary account

        #  starting from cenario where we create a transaction (produces a OPEN transaction)
        * call read('classpath:features/create_transaction.feature@CreateTransaction')

        Given url apiURL
        And path "/transaction/" + transactionID + "/complete"
        And request {}
        When method POST
        Then assert responseStatus == 202
        
        Given url apiURL
        And path "/transaction_status"
        And param transaction_id = transactionID
        When method GET
        Then assert responseStatus == 200
        And assert response.status == "BOOKED"

        Given url apiURL
        And path "/balance_provisions"
        And param transaction_id = transactionID
        When method GET
        Then assert responseStatus == 200
        And assert response.length == 1
        And assert response[0].type == "ADD"
        And assert response[0].status == "CLOSED"

        # check amount was added to destination account
        Given url apiURL
        And path "/account/" + destinationAccount.account_id
        When method GET
        Then assert responseStatus == 200
        And assert response.balance ==  (destinationAccountBalanceBefore + transactionAmount)

        # check the intermediary account have the original balance
        Given url apiURL
        And path "/account/12345"
        When method GET
        Then assert responseStatus == 200
        And assert response.balance ==  intermediaryAccountBalanceBefore

    Scenario: should failure completing transaction

        
