# Exercise #8: Phone Number Normalizer

[![exercise status: released](https://img.shields.io/badge/exercise%20status-released-green.svg?style=for-the-badge)](https://gophercises.com/exercises/phone)

## Exercise details

This exercise is fairly straight-forward - we are going to be writing a program that will iterate through a database and normalize all of the phone numbers in the DB. After normalizing all of the data we might find that there are duplicates, so we will then remove those duplicates keeping just one entry in our database.

While the primary goal is to create a program that normalizes the database, we will also be looking at how to setup the database and write entries to it using Go as well. This is done intentionally to try to teach how to use Go when interacting with an SQL database.

There are many ways to use SQL in Go. Rather than just picking one, I am going to try to cover a few. If you would like to see any additional libraries covered feel free to reach out and I'll try to add it. For now, here are the libraries I intend to cover:

- Writing raw SQL and using the [database/sql](https://golang.org/pkg/database/sql/) package in the standard library
- Using the very popular [sqlx](https://github.com/jmoiron/sqlx) third party package, which is basically an extension of Go's sql package.
- Using a relatively minimalistic ORM (I will be using [gorm](https://github.com/jinzhu/gorm))

We will also need to explore some basic string manipulation techniques to normalize our phone numbers, but the primary focus here is going to be on the SQL so I'll try to keep that part of the code relatively simple.

I am intending on using SQLite and Postgres for the videos, but if you would like to see a MySQL example (it should be nearly identical to the other two) let me know and I'll try to add it once the rest of the Gophercises are done.

On to the exercise - we will start by creating a database along with a `phone_numbers` table. Inside that table we want to add the following entries (yes, I know there are duplicates):

```
1234567890
123 456 7891
(123) 456 7892
(123) 456-7893
123-456-7894
123-456-7890
1234567892
(123)456-7892
```

You can organize your table however you want, and you may add whatever extra fields you want. My tables will likely vary depending on which of the libraries I'm using, as ORMs like GORM will often add a few fields for us automatically. You may also create the table manually (via raw SQL) if you want, but try to insert the entries using Go code.

Once you have the entries created, our next step is to learn how to iterate over entries in the database using Go code. With this we should be able to retrieve every number so we can start normalizing its contents.

Once you have all the data in the DB, our next step is to normalize the phone number. We are going to update all of our numbers so that they match the format:

```
##########
```

That is, we are going to remove all formatting and only store the digits. When we want to display numbers later we can always format them, but for now we only need the digits.

In the list of numbers provided, the first entry, along with the second to last entry, match this format. All of the others do not and will need to be reformatted. There are also some duplicates that will show up once we have reformatted all the numbers, and those will need removed form the database but don't worry about that for now.

Once you written code that will successfully take a number in with any format and return the same number in the proper format we are going to use an `UPDATE` to alter the entries in the database. If the value we are inserting into our database already exists (it is a duplicate), we will instead be deleting the original entry.

When your program is done your database entries should look like this (the order is irrelevant, but duplicates should be removed):

```
1234567890
1234567891
1234567892
1234567893
1234567894
```


## Bonus

There isn't a concrete bonus for this exercise. If you want you can explore other third party SQL libraries and other SQL databases (MySQL, etc), but that is up to you.

## Additional Resources

I have written about Go and PostgreSQL a good bit here - https://www.calhoun.io/using-postgresql-with-golang/  
While most articles are about PostgreSQL, using other variants of SQL tends to be nearly identical with a few minor exceptions (like when connecting to the database or using DB-specific features).

Not all of the articles are complete, but if you get stuck check them out for some help. Most notably, if you are using Postgres and getting an error like `pq: SSL is not enabled on the server` when trying to connect I recommend looking at this article - https://www.calhoun.io/connecting-to-a-postgresql-database-with-gos-database-sql-package/
