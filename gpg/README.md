# GPG

This package is used to work with

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