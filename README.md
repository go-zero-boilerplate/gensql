# gensql
SQL generator for golang (generating static go code instead of using reflection)


## Features
- Entity - golang struct
- Iterator (using the **paginator** of [databases](https://github.com/go-zero-boilerplate/databases)) to iterate through entities in the Database by using Offset and Limit in the background to "page" through your entries but in your code all you need is to call `Next()`
- Repository - had some inspiration from Entity Framework of C#. But the repository has your common CRUD calls like `GetByPk` , `List` , `Add` , `Delete` , `Save`
- Generates the SQL `CREATE` sql statements

## Usage example (quick start)

Steps are:

- Install gensql - `go get -u github.com/go-zero-boilerplate/gensql`
- Create YAML setup file (see example below)
- Run `gensql -in="YOUR_YAML_FILE" -out="YOUR_OUT_DIR"`



## Example YAML setup

```
blog:
  dialect: mysql
  fields:
    - id              int64 pk auto
    - name            string
    - description     text
    - created
    - updated
  uniques:
    - [name]
```

Regarding the above yaml setup:
- The `created` field "implies" the expected behavior - ie. to use `CURRENT_TIMESTAMP` as the `DEFAULT` value and a field type `DATETIME` in the MySql dialect
- The `updated` field "implies" the expected behavior - ie. to use `CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP` as the `DEFAULT` value for the field in the MySql dialect. This means in mysql terms it is the trigger
- The number of spaces after the field name and the type with arguments does not matter as long as there is one. So the above could have been `id int64   pk auto`
- By default all fields are `NOT NULL`, use argument `nullable` to make it nullable
- If any of args are a number it implies to apply a size. Currently if no size is specified for a `string` (`VARCHAR`) field it defaults to 200 (due to mysql `INDEX` length limitations). To obtain no size limit for a field use `text` instead of `size`
- Arguments `pk` results in `PRIMARY KEY` and `auto` in `AUTO_INCREMENT` for mysql
- Default values can be specified like `default:123`


## Acknowledgments

Inspiration from https://github.com/drone/sqlgen and C# Entity Framework.