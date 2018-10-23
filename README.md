# Restreamer docs and website

## Setup, Develop, and Deploy

1. Clone the Restreamer repository:   
   
        $ git clone git@github.com:datarhei/restreamer.git gh-pages

2. Check out the "gh-pages" branch.

        $ cd /path/to/gh-pages
        $ git checkout gh-pages

3. Build the docker image with Jekyll

        $ docker build -t jekyll:restreamer .

4. Start the web server

        $ docker run -d -p 4000:4000 -v ${PWD}:/gh-pages jekyll:restreamer

5. Visit the site at http://localhost:4000/restreamer/

6. Edit the pages. Pages get automatically rebuild as soon as you save them.

7. Commit your changes

        $ git add [your changes files]
        $ git commit -m "[describe your changes]"

8. Push your changes

        $ git push origin gh-pages

9. The page is now available on https://datarhei.github.io/restreamer/
