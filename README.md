# Restreamer docs and website

## Develop and Deploy

1. Clone the Restreamer repository:   
   
```sh
git clone git@github.com:datarhei/restreamer.git gh-pages
```

2. Check out the "gh-pages" branch.

```sh
cd gh-pages
git checkout gh-pages
```

3. Build the docker image with Jekyll

```sh
docker build -t jekyll:restreamer .
```

4. Start the web server

```sh
docker run -d --rm -p 4000:4000 -v ${PWD}:/gh-pages jekyll:restreamer

# or

docker run -it --rm -p 4000:4000 -v ${PWD}:/gh-pages jekyll:restreamer
```

5. Visit the site at http://localhost:4000/restreamer/

6. Edit the pages. Pages get automatically rebuild as soon as you save them.

7. Commit your changes

```sh
git add [your changed files]
git commit -m "[describe your changes]"
```

8. Push your changes. This will automatically deploy the changes and makes them public!

```sh
git push origin gh-pages
```

9. The page is now available on https://datarhei.github.io/restreamer/
