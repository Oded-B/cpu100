# cpu100

Creates 100% CPU utilization for testing.

Helps to detect how much hashes or signatures can be computed maximum per second. Requests to services uses JWT for authorized access, so test can detect limit of empty authorized requests without any overhead that service can receive.

Command line parameters:

    -a string
        hash or signature algorithm, can be: md5, sha1, sha224, sha256, sha384, sha512, sha512/224, sha512/256, ecdsa, ed25519 (default "sha256")
    -b int
        length of random bytes block to calculate for each hash (default 1024)
    -d duration
        duration of program working (in format '1d8h15m30s') (default 1h30m0s)
    -n int
        number of threads to start (default 8)

(c) schwarzlichtbezirk, 2021.
