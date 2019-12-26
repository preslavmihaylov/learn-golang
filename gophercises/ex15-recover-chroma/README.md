# Exercise #15: Development Panic/Recover Middleware with Chroma

[![exercise status: released](https://img.shields.io/badge/exercise%20status-released-green.svg?style=for-the-badge)](https://gophercises.com/exercises/recover_chroma)


## Exercise details

In the [recover](https://gophercises.com/exercises/recover) exercise we learned how to create some HTTP middleware that recovers from any panics in our application and renders a stack trace if we are in a local development environment. In this exercise we will be taking that code a step further; we will be adding in the ability to navigate to any source file in the panic stack trace in order to make it easier to debug issues when they arise in a development environment.

Given the web server and the recovery middleware in `main.go`, add the following to the application:

#### 1. An HTTP handler that will render source files in the browser

This is left intentionally vague, but our primary goal is to update our application so that we can view the source code of Go files in our browser.

If you need help getting started, I would first focus on writing this handler independent of our existing application. After that, try adding it into the `devMw` function, having it only render pages if a specific path prefix (like `/debug/`) is used in the path.

#### 2. Add syntax highlighting to the source file rendering

Once you have source files displayed in the browser, try using the [chroma](https://github.com/alecthomas/chroma) package to add syntax highlighting to your source code.

#### 3. Parse the stack trace & creating links

With a source code handler in place you should now be ready to start parsing the stack trace in order to figure out which source files are mentioned in it.

For now just try to parse out the path to any source files, as well as the line number in the file where the stack trace points to. In my code I use characters like a tab and the colon (`:`) to match my string, but you can choose whatever you think is best.

Using the code from your parsing, output links in your `devMw` handler's stack trace that link to the correct path where a source file can be viewed. For this step you may need to use something like [url.Values](https://golang.org/pkg/net/url/#Values) to encode the path to the source file.

#### 4. Add line highlighting

Chroma supports [line highlighting](https://github.com/alecthomas/chroma#the-html-formatter) and our stack trace has the line number where the panic occurred. Update your source code rendering handler to accept a `line` query parameter, and add that to your links created in step 4.

## Useful links

- <https://gophercises.com/exercises/recover>
- <https://github.com/alecthomas/chroma>
- <https://regex-golang.appspot.com>
- <https://golang.org/pkg/regexp/>
