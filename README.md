# CapitalData Application Test

# Requirements

1 - Packages :
  Run `go get golang.org/x/net/context`to load net context before calling the Application
2 - Data Base :
  1 ° Install sqlite3

  2 ° Create database into the app directory : `sqlit3 contacts.db`

  3 ° Create table contacts :

    ```
    CREATE TABLE `contacts` (
      `uid` VARCHAR(64) NULL,
      `email` VARCHAR(64),
      `cookie` VARCHAR(64),
      UNIQUE(email,uid,cookie) ON CONFLICT REPLACE,
      UNIQUE(email,uid) ON CONFLICT REPLACE,
      UNIQUE(email,cookie) ON CONFLICT REPLACE,
      UNIQUE(email) ON CONFLICT REPLACE
      );
    ```

3 - Build the app
  From the CapitalDataApplicationContact source directory run `go install`

# Run Application

1 - Launch the app : `CapitalDataContacts`

2 - Fill an option to the app following the apps commands

# Tests

No tests provided in this version

# Author

Alexandre Carle

# Github

https://github.com/XCarle/capitalDataContacts
