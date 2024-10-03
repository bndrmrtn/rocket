# Rocket - A Query Language Helper

## Installation

To get started, you need to install the Rocket generator package, which provides the `rck` command:

```bash
go install github.com/bndrmrtn/rocket@latest
```

### Usage

Rocket (`rck`) is a command-line interface (CLI) for parsing files with a `.rocket` extension.

#### Parsing a File

Use the following command to parse a Rocket file:

```bash
rocket generate --file="path/to/data.rocket" --language="go" --database="mysql"
```

You can also include the `--output` flag to specify the output folder. The default output is `--output="*.{ext}"`.
**Note:** Remember to include `{ext}` since Rocket may generate files with multiple extensions.

## How to Write Rocket Files

Rocket files are straightforward to create. Each Rocket file consists of `blocks`, and each `block` functions differently.

The predefined blocks for Rocket include: `schema`, `enum`, `model`, `hashing`, and `query`.

### Schema Block

The schema block acts like a model, serving as a wrapper that can be imported into other models.

```
schema Base {
  id number primary increment
  created_at datetime @default(now())
}
```

### Enum Block

You can define enumerations like this:

```
enum Role {
  ADMIN = "admin"
  USER = "user"
}
```

### Model Block

The model block describes a database table. Each line within a block represents a column. Each row begins with a column name followed by a `datatype`. The `&` symbol can be used to import a defined schema.

Attributes and annotations follow the datatype. An attribute does not have a value, but annotations can. An annotation begins with `@` and may include data within parentheses.

In the example below, `@sensitive`, `@hash("password")`, and `@default(&Role.USER)` are annotations, while `nullable` is an attribute. Note that `@sensitive` and `@hash` do not affect the generated SQL; they assist the program in creating and applying hashes automatically. The `@default` annotation sets the default value for the column.

```
model User {
  &Base
  first_name string
  password string @sensitive @hash("Password")
  phone string nullable
  role &Role @default(&Role.USER)
}
```

### Hashing Block

The hashing block defines hash configurations for sensitive fields. Rocket currently supports several hashing algorithms: `bcrypt`, `sha256`, `sha512`, `md5`, and `sha1`. The function arguments differ for each algorithm.

```
hashing Password {
  algo bcrypt(15)
}
```

### Query Block

The query block contains a Rocket Query **(currently no release)**.

Queries can be simple, such as retrieving all users. The `[]` indicates that the result should be an array containing multiple users. The `get` keyword signifies a selection from the database, and `User` specifies the model to query.

```
query getAllUser {
  []get User
}
```

You can also create more specific queries, as shown below. The fields within `{}` indicate the database columns to select. In this case, since the `[]` token is omitted, Rocket understands that it should return **only one** result.

```
query findUserName(_id number) {
  get{User.first_name} User.Where(id == _id)
}
```
