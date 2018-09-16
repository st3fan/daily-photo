# daily-photo
*Stefan Arentz, September 2018*

When I asked Instagram for a full data dump of my 3000+ private photos I realized than many photos would be fun to share publicly.

So here is a script that does that. It picks a random photo from a directory and posts it to Twitter. The caption comes from a text file that follows the same naming as the photo. I run it from a daily *cron job* on an OpenBSD machine that is provided by [openbsd.amsterdam](https://openbsd.amsterdam).

It needs a Twitter API account, which you can setup in their developer portal. Just fill in the API keys/secrets in the `daily-photo.sh` script.

I manage/edit the captions with another hack, which you can find at [github.com/st3fan/daily-photo-editor](https://github.com/st3fan/daily-photo-editor). 

This is pretty much written for me, but if you think this is useful to you, feel free to fork or submit a pull request to make things more generic.
