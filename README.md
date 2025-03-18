
# MongoDB Migrate

MongoDB Migrate is a tool for migrating databases between two different MongoDB instances. This project is designed to facilitate the transfer of data from one MongoDB instance to another, ensuring a smooth and reliable migration process.

## Features

- Migration of databases between two MongoDB instances
- Support for migrating collections and documents
- Management of MongoDB connections via cli options
- Detailed logs to monitor the migration process

## Installation

To install MongoDB Migrate, you need to have Go installed on your system. You can install Go by following the instructions on the official [Go](https://golang.org/)) website.

Clone the repository:

```bash
git clone https://github.com/bizmate-oss/mongodb-migrate.git
cd mongodb-migrate
```

Build the project:

```bash
go build
```

## Usage

Run the migration:

```bash
./mongodb-migrate -src mongodb://<host>:<port>/<database> -dst mongodb://<host>:<port>/<database>
```

## Contributing

If you would like to contribute to this project, feel free to fork the repository and submit a pull request. Make sure to follow the contribution guidelines and thoroughly test your changes.

## License

This project is licensed under the Apache 2.0 License. See the LICENSE file for more details.

## Contact

For questions or support, you can contact the project maintainer through the issues page.