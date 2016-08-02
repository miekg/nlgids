# NLGIDS

A public-private repo. You may fork and whatever, but this is some custom Caddy middleware for a
website. This interacts with Google calendar and is capable of sending email, creating pdf, by
filling out a text.Template and calling LaTeX and then sending that email.

It is also capable of creating invoices for tours.

See <https://nlgids.london> for the result.

## BUGS

Multiday events are not parsed, if you want to set a period as busy, you'll have to do this on a day
to day basis.
