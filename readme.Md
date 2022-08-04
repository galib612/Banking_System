several way to create CURD REestFul api


    1.) database/sql (low-level)
        very fast & straightforward
        manual mapping SQL fields to variables
        easy to make mistakes, not caught until runtime.

    2.) GORM (high-level)
        CRUD functions alreadby implemented
        very short production code
        must learn to write queries using gorm's functions
        run slowly on high load

    3.) SQLX (mid level)
        quite fast & easy to use
        fields mapping via query text & struct tags
        failure won't occur until runtime

    4.) SQLC
        very fast & easy to use
        Automatic code generation
        catch SQL query errors before generating codes
        Full support Postgres. Mysql is experimental {so if u are using postgres SQLC is right tool}




        ...............................................................................

        What are "Read Phenomena"? Well, this comes straight from the SQL-92 standard in black and white.

        Dirty Read: A transaction reads data writeen by other concurrent uncomitted transaction.

        Non-Repeatable Read: a transaction reads the same row twice and sees different value because
            it has been modified by other committed transaction.

        Phantom Read: A transaction re-executed a query to find rows that satisfy a condition and
           sees a different set of rows, due to chnages by other committed transaction.

        Serialization Anomaly: The result of group of concurrent committed transaction is impossible
            to achieve if we try them sequentially in any order without overlapping.


        4 Standard Isolation Levels (Ansi)

        ...........1........................2...........................3.......................4...>
        Read Uncommitted          Read Committed             Repeatable Read           Serializable

        Read Uncommitted: Can see data written by uncommitted transaction.

        Read Committed: Only see data written by committed transaction.

        Repeatable Read: Same read query always returns same result.

        Serializable: Can achieve same result if execute transactions serially in some order instead
            of concurrently.



    Isolation Level	    Dirty Read	    Nonrepeatable Read	    Phantom Read	 Serialization Anomaly
    Read uncommitted  Allowed, but not in PG	Possible	    Possible	    Possible
    Read committed	    Not possible	Possible	            Possible	    Possible
    Repeatable read	    Not possible	Not possible	Allowed, but not in PG	Possible
    Serializable	    Not possible	Not possible	        Not possible	Not possible

        .......              ..............               .............           ...............
        
        Some Query for isolation level.
            1. show transaction isolation level;
            2. set session transaction isolation level <isolation level>;
        
    Note:- In postgres Read Uncommitted behave likes Read Committed.
           MySql use locking mechanism to achieve serializable isolation level but postgres not use locking mechanism instead don't reexecute same query and throw errors and hints(Dependencies Detection).

           By Default Read Committed in postgres.
           By Default Repeatable read in MySql.


           ..........................................................................................

           Popular golang web frameworks with github star

           1. Gin      42.4k               5. Martini        11.1k
           2. beego    25.1k               6. Fiber          8.6k
           3. Echo     18.3k               7. Buffalo        5.9l
           4. Revel    11.9k

           Popular Http routers

           1. FastHttp       13.6k
           2. Gorill Mux     12.9k        
           3. HttpRouter     11.9k  
           4. Chi            8.3k




           