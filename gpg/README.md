# GPG

A tool that provides a few functions to use existing keys to encrypt or decrypt
messages.

## Summary

This tool can be applied to any project that needs to handle GPG-encrypted
messages. It provides functions to use existing keys to encrypt or decrypt
messages. It uses [gopenpgp] library for all GPG processing.

# Generate Test Key

The following process was used to generate a key for using with tests.

```text
gpg --full-generate-key
Real Name: GPG Test
Email: gpgtest@example.com
Comment: Login package GPG test
Passphase: test1234

$Env:GPG_UID="GPG Test (Login package GPG test) <gpgtest@example.com>"
$Env:GPG_FILE_PREFIX="gpg-test"
gpg --export --armor "${Env:GPG_UID}" > "${Env:GPG_FILE_PREFIX}.public.asc"
gpg --export-secret-keys --armor "${Env:GPG_UID}" > "${Env:GPG_FILE_PREFIX}.private.asc"
gpg --export-secret-subkeys --armor "${Env:GPG_UID}" > "${Env:GPG_FILE_PREFIX}.sub_private.asc"
```

---

[gopenpgp]: https://pkg.go.dev/github.com/ProtonMail/gopenpgp/v3#section-readme
