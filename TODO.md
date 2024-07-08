Generalized Media Scraper
Imagine media ( audio, video, images, pdf,..) are being stored in some websites. We need to create a program such that we can scrape out the entire website targetting ( read more below ) specific set of media, and downloading all of them in the form of ( original_url , actual_stored_file, metadata_text )

Targetting can be done via starting with a single URL or can be done with url pattern matching.

The program should be such that:
one should be able to add websites into it with ease - i.e. almost no code required to scrape through different websites

Automated retries on failure - on full failire, put the failure into error logs

In case of too many failures - abort. Too many failure is an absolute or relative number which are to come from configuration.

Should be able to do it very fast, fastet possible.

There would be server throttling, code against it.
