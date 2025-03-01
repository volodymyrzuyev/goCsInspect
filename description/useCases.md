# Use Cases:
1. Fetch item details:
    * Description:
        Get item float, pattern, stickers, name based on "M" "A" "D" "S"
    * Preconditions:
        1. Internet connection must exist
        2. DB must be up
        3. At least 1 account has to be logged into the client.
        4. Valid connection with GC has to exist
    * Postcondition:
        1. Account used for inspect goes on cooldown
        2. Next account should be ready to used
    * Success Scenario:
        Item details obtained and decoded, and returned to the requester
    * Failure scenario:
        1. Invalid items:
            * Returns an error code to the user
        2. GC cooldown not passed on any accounts:
            * Returns error to the user
    * Status:
        WIP

2. Log into account:
    * Description:
        Log into steam client with username + password + (2FA code || secret hash)
    * Preconditions:
        1. Internet connection must exist
        2. DB must be up
    * Postconditions:
        1. Account logged into steam
    * Success Scenario:
        Account is ready to handle inspect requests
    * Failure scenario:
        1. Invalid username/password/2FA:
            * Return error to user
        2. Steam blocking new logins:
            * Inform the user, exit program
    * Status:
        WIP

3. Decode "M" "A" "D" "S" params from inspect link:
    * Description:
        Decode params needed for fetching of the skin info
    * Preconditions:
        1. NONE
    * Postconditions:
        1. NONE
    * Success Scenario:
        Item params are returned as a struct
    * Failure scenario:
        1. Invalid link:
            * Return error to the user
    * Status:
        WIP

4. Add item to DB:
    * Description:
        Adds inspected item to database
    * Preconditions:
        1. Internet connection must exist
        2. DB must be up
        3. "Fetch item details" must successfully run
    * Postconditions:
        1. Item added to the database for caching
    * Success Scenario:
        Item added to DB with no errors
    * Failure Scenario:
        1. Invalid item:
            Return error code
        2. Duplicate Item:
            Return error code
    * Status:
        WIP

