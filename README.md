# ReStreamer Docs and Website

## Setup

1. Clone the ReStreamer repository

2. Change into the "docs" directory where docs live

        $ cd docs

3. Build the Dockerfile

        $ docker build -t rs-docs .

4. Start the web server

        $ docker run -d -p 4000:4000 -v $PWD:/restreamer rs-docs

5. Visit the site at
   http://docker-host-ip:4000/restreamer/

## Deployment

1. Clone a separate copy of the ReStreamer repo as a sibling of your normal
   ReStreamer project directory and name it "restreamer-gh-pages".

        $ git clone git@github.com:datarhei/restreamer.git restreamer-gh-pages

2. Check out the "gh-pages" branch.

        $ cd /path/to/restreamer-gh-pages
        $ git checkout gh-pages

3. Copy the contents of the "docs" directory in master to the root of your
   restreamer-gh-pages directory.

        $ cd /path/to/restreamer
        $ cp -r docs/** ../restreamer-gh-pages

4. Change to the small-restreamer-gh-pages directory, commit, and push the changes

        $ cd /path/to/restreamer-gh-pages
        $ git commit . -m "Syncing docs with master branch"
        $ git push
