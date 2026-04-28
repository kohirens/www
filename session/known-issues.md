# Known Issues

* If you know that someone is currently logged in, and have access to get their
  session ID. Then you can easily hijack their session. Simply go to the site
  and set the session ID cookie to the value of the session you want to hijack.
  There is currently no way to validate if a session actually belongs to a
  specific browser.
  * Possible mitigations:
    1. Tie a session ID to geo-location data.
