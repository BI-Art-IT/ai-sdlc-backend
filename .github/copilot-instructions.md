All code shpoud refer the tech design in https://github.com/BI-Art-IT/AirLines/commit/67ad271c77fd8bb5cb3e9e713603c3ecd4704582

The codee standard:>
 - Use go as programming language
 - Use PostgresSQL
 - No ORM framework
 - Use REST API as a communication protocol
 - Use GOLANG-migrate and GOLANG-migrate for pgx5 for any table DDL or data setup, put in a separate SQL file. Rule: One table per SQL file.

Working process
- Ensure to use clean codee approach
- Always create unit tests whenever possible
- Ensure Unit ytests passes before sumbitting pull request
- Use conventional commit on pull request title
