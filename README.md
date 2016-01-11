# Restreamer docs and website

## Setup

1. Clone the Restreamer repository:   
   
        $ git clone git@github.com:datarhei/restreamer.git gh-pages

2. Change into the directory where docs live

        $ cd /path/to/gh-pages

3. Start the web server

        $ docker run -d -p 4000:4000 -v $PWD:/restreamer datarhei/jekyll:latest

4. Visit the site at
   http://docker-host-ip:4000/restreamer/

## Deployment

1. Clone a separate copy of the Restreamer repo as a sibling of your normal
   Restreamer project directory and name it "gh-pages".

        $ git clone git@github.com:datarhei/restreamer.git gh-pages

2. Check out the "gh-pages" branch.

        $ cd /path/to/gh-pages
        $ git checkout gh-pages

3. Commit and push the changes

        $ git commit . -m "Syncing with gh-pages branch"
        $ git push origin gh-pages
