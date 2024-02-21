# Using PASS Utility

Pass is a simple Linux Password manager which stores the data in gpg2 encrypted files. It is used to manipulate password
store like adding/removing password, encrypting a piece of data, decrypting an encrypted file or storing passwords.

## Installing PASS

### Installing in Ubuntu

You can install pass using apt utility
`sudo apt-get install pass`

### Installing in RHEL/CentOS/Fedora

In Red Hat, Pass is provided through EPEL Repository. You can enable epel repo using following command

```
sudo dnf install https://dl.fedoraproject.org/pub/epel/epel-release-latest-8.noarch.rpm
```

Once installed, you can install pass utility using dnf
```sudo dnf install pass```

## Setting up a gpgkey

GPG Stands for GNU PRivacy Guard and it is an encryption and signing tool. It is based on OpenPGP method. PASS Utility
uses a gpg key to encrypt and/or decrypt data. In order to to use pass, first a gpg key needs to be created. You can
create a gpg key using `gpg` utility in linux.

* To create a new gpg key, simply run the following command

  ```
  gpg --quick-gen-key
  ```

  or for a more customized key

  ```
  gpg --gen-key
  ```

* To list available keys

  ```
  gpg --list-keys
  ```

* To delete a key

  ```
  gpg --delete-keys key_name
  ```

* To encrypt data

  ```
  gpg --encrypt data_or_file
  ```

* To decrypt data

  ```
  gpg --decrypt encrypted_file
  ```

## Setting up PASS password store

You can use the provided script `init-pass.sh` to setup the password store automatically. If you want to manually create
a password store, you can use

  ```
  pass init "pass_store_name"
  ```

Onomy chain stores everything in the pass store under `keyring-onomy` name.

* List password stores To list contents of password store, you can use
  `pass ls` or `pass show` command. It will list all the encrypted files in all the password store.

  These are some commands used to manipulate password store

* Insert a new key

  ```
  pass insert keyring/newKey
  ```

* Decrypt and print a stored key

  ```
  pass show keyring/StoredKey
  ```

* Edit already stored key

  ```
  pass edit keyring/StoredKey
  ```

* Find a specific phrase the decrypted files

  ```
  pass grep searchString
  ```

* Find a stored key

  ```
  pass find key_name
  ```

  or

  ```
  pass search key_name
  ```
