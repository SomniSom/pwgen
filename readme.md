# Password Reusable Generator

If you're tired of storing passwords in browsers where they vanish or struggle to come up with strong passwords, this application is for you. You input a domain and a master password
(known only to you), and receive a site-specific password. If it's time to change the password, set a version number and get a new one.

Lost access to browser passwords due to system failure, etc., install the app, enter the domain and password, and you can log back into your favorite sites again.

## Basic Functionality

1. **Password Generation by Domain and Master Password**
2. **Storage of Additional Encrypted Information in DB** (can be lost or publicly stored; passwords are not stored there)
   * Domain
   * Login (optional, can be omitted)
   * Version (optional, used for generating a new password for an old site)
3. **Search by Domain** (substring matching)
4. **Search by Login** (substring matching)
5. **Application Version**

## Password Generation

Uses the ChaCha8 algorithm for randomness. A combination of the password and domain is used to ensure uniqueness of the algorithm.

## Commands

* `pw.gen gen -d <domain> -p <password> [-l <login> -v <version:int>]`
* `pw.gen logins -p <password> [-d <domain>]`
* `pw.gen domains -p <password> [-l login]`
* `ps.gen version`

## Roadmap

* _Web UI_ for viewing the database of domains and logins, as well as retrieving passwords
* _Chrome Extension_ for convenient password input
* _Firefox Extension_ (see above)
